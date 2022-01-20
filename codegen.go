//
// code generator
//
package main

import (
	"fmt"
	"io"
	"os"
)

type codeWriter struct {
	w   io.Writer
	err error
}

func (c *codeWriter) println(frmt string, a ...interface{}) {
	if c.err != nil {
		return
	}
	_, c.err = fmt.Fprintf(c.w, frmt, a...)
	_, c.err = fmt.Fprintln(c.w)
}

var i int = 1

func count() int {
	i++
	return i
}

var depth int
var argreg8 = []string{"%dil", "%sil", "%dl", "%cl", "%r8b", "%r9b"}
var argreg16 = []string{"%di", "%si", "%dx", "%cx", "%r8w", "%r9w"}
var argreg32 = []string{"%edi", "%esi", "%edx", "%ecx", "%r8d", "%r9d"}
var argreg64 = []string{"%rdi", "%rsi", "%rdx", "%rcx", "%r8", "%r9"}

var curFnInGen *Obj

func (c *codeWriter) push() {
	if c.err != nil {
		return
	}

	c.println("	push %%rax")
	depth++
}

func (c *codeWriter) pop(arg string) {
	if c.err != nil {
		return
	}

	c.println("	pop %s", arg)
	depth--
}

func alignTo(n, align int) int {
	return (n + align - 1) / align * align
}

// Pushes the given node's address to the stack
func (c *codeWriter) genAddr(node *Node) {
	if c.err != nil {
		return
	}

	switch node.Kind {
	case ND_VAR:
		if node.Obj.IsLocal {
			//  Local variable
			c.println("	lea %d(%%rbp), %%rax", node.Obj.Offset)
			return
		}
		// Global variable
		c.println("	lea %s(%%rip), %%rax", node.Obj.Name)
		return
	case ND_DEREF:
		c.genExpr(node.Lhs)
		return
	case ND_COMMA:
		c.genExpr(node.Lhs)
		c.genAddr(node.Rhs)
		return
	case ND_MEMBER:
		c.genAddr(node.Lhs)
		c.println("	add $%d, %%rax", node.Mem.Offset)
		return
	default:
		fmt.Fprintf(os.Stderr, "\nnode: %#v\n\n", node)
		fmt.Fprintf(os.Stderr, "node.Lhs: %#v\n\n", node.Lhs)
		c.err = fmt.Errorf(errorTok(node.Tok, "not an lvalue"))
	}

}

func (c *codeWriter) load(ty *Type) {
	if c.err != nil {
		return
	}

	if ty.Kind == TY_ARRAY || ty.Kind == TY_STRUCT {
		// If it is an array, do not attempt to load a value to the
		// register because in general we can't load an entire array to a
		// register. As a result, the result of an evaluation of an array
		// become not the array itself but the address of the array.
		// This is where "array is automatically converted to a pointer to
		// the first element of the array in C" occurs.
		return
	}

	// When we load a char or a short value to a register, we always
	// extend them to the size of int, so we can assume the lower half of a
	// register for char, short and int may contain garbage. When we load
	// a long value to a register, it simply occupies the entire register.
	switch ty.Sz {
	case 1:
		c.println("	movsbl (%%rax), %%eax")
	case 2:
		c.println("	movswl (%%rax), %%eax")
	case 4:
		c.println("	movsxd (%%rax), %%rax")
	case 8:
		c.println("	mov (%%rax), %%rax")
	default:
		c.err = fmt.Errorf("invalid size")
		return
	}
}

func (c *codeWriter) store(ty *Type) {
	if c.err != nil {
		return
	}

	c.pop("%rdi")

	if ty.Kind == TY_STRUCT {
		for i := 0; i < ty.Sz; i++ {
			c.println("	mov %d(%%rax), %%r8b", i)
			c.println("	mov %%r8b, %d(%%rdi)", i)
		}
		return
	}

	switch ty.Sz {
	case 1:
		c.println("	mov %%al, (%%rdi)")
	case 2:
		c.println("	mov %%ax, (%%rdi)")
	case 4:
		c.println("	mov %%eax, (%%rdi)")
	case 8:
		c.println("	mov %%rax, (%%rdi)")
	default:
		c.err = fmt.Errorf("invalid size")
	}
}

func (c *codeWriter) cmpZero(ty *Type) {
	if c.err != nil {
		return
	}

	if isInteger(ty) && ty.Sz <= 4 {
		c.println("	cmp $0, %%eax")
	} else {
		c.println("	cmp $0, %%rax")
	}
}

const (
	I8 = iota
	I16
	I32
	I64
)

func getTypeId(ty *Type) int {
	switch ty.Kind {
	case TY_BYTE:
		return I8
	case TY_SHORT:
		return I16
	case TY_INT:
		return I32
	}
	return I64
}

// The table for type casts
const i32i8 string = "movsbl %al, %eax"
const i32i16 string = "movswl %ax, %eax"
const i32i64 string = "movsxd %eax, %rax"

var castTable [4][4]string = [4][4]string{
	{"", "", "", i32i64},        // i8
	{i32i8, "", "", i32i64},     // i16
	{i32i8, i32i16, "", i32i64}, // i32
	{i32i8, i32i16, "", ""},     // i64
}

func (c *codeWriter) cast(from *Type, to *Type) {
	if c.err != nil {
		return
	}

	if to.Kind == TY_VOID {
		return
	}

	if to.Kind == TY_BOOL {
		c.cmpZero(from)
		c.println("	setne %%al")
		c.println("	movzx %%al, %%eax")
		return
	}

	t1 := getTypeId(from)
	t2 := getTypeId(to)
	if castTable[t1][t2] != "" {
		c.println("	%s", castTable[t1][t2])
	}
}

func (c *codeWriter) genExpr(node *Node) {
	if c.err != nil {
		return
	}

	c.println("	.loc 1 %d", node.Tok.LineNo)

	switch node.Kind {
	case ND_NULL_EXPR:
		return
	case ND_NUM:
		c.println("	mov $%d, %%rax", node.Val)
		return
	case ND_NEG:
		c.genAddr(node.Lhs)
		c.println("	neg %%rax")
		return
	case ND_VAR, ND_MEMBER:
		c.genAddr(node)
		c.load(node.Ty)
		return
	case ND_DEREF:
		c.genExpr(node.Lhs)
		c.load(node.Ty)
		return
	case ND_ADDR:
		c.genAddr(node.Lhs)
		return
	case ND_ASSIGN:
		c.genAddr(node.Lhs)
		c.push()
		c.genExpr(node.Rhs)
		c.store(node.Ty)
		return
	case ND_STMT_EXPR:
		for n := node.Body; n != nil; n = n.Next {
			c.genStmt(n)
		}
		return
	case ND_COMMA:
		c.genExpr(node.Lhs)
		c.genExpr(node.Rhs)
		return
	case ND_CAST:
		c.genExpr(node.Lhs)
		c.cast(node.Lhs.Ty, node.Ty)
		return
	case ND_MEMZERO:
		// `rep stosb` is equivalent to `memset(%rdi, %al, %rcx)`.
		c.println("	mov $%d, %%rcx", node.Obj.Ty.Sz)
		c.println("	lea %d(%%rdp), %%rdi", node.Obj.Offset)
		c.println("	mov $0, %%al")
		c.println("	rep stosb")
		return
	case ND_COND:
		cnt := count()
		c.genExpr(node.Cond)
		c.println("	cmp $0, %%rax")
		c.println("	je .L.else.%d", cnt)
		c.genExpr(node.Then)
		c.println("	jmp .L.end.%d", cnt)
		c.println(".L.else.%d", cnt)
		c.genExpr(node.Els)
		c.println(".L.end.%d", cnt)
		return
	case ND_NOT:
		c.genExpr(node.Lhs)
		c.println("	cmp $0, %%rax")
		c.println("	sete %%al")
		c.println("	movzx %%al, %%rax")
		return
	case ND_BITNOT:
		c.genExpr(node.Lhs)
		c.println("	not %%rax")
		return
	case ND_LOGAND:
		cnt := count()
		c.genExpr(node.Lhs)
		c.println("	cmp $0, %%rax")
		c.println("	je  .L.false.%d", cnt)
		c.genExpr(node.Rhs)
		c.println("	cmp $0, %%rax")
		c.println("	je  .L.false.%d", cnt)
		c.println("	mov $1, %%rax")
		c.println("	jmp  .L.end.%d", cnt)
		c.println(".L.false.%d:", cnt)
		c.println("	mov $0, %%rax")
		c.println(".L.end.%d:", cnt)
		return
	case ND_LOGOR:
		cnt := count()
		c.genExpr(node.Lhs)
		c.println("	cmp $0, %%rax")
		c.println("	jne  .L.true.%d", cnt)
		c.genExpr(node.Rhs)
		c.println("	cmp $0, %%rax")
		c.println("	jne  .L.true.%d", cnt)
		c.println("	mov $0, %%rax")
		c.println("	jmp  .L.end.%d", cnt)
		c.println(".L.true.%d:", cnt)
		c.println("	mov $1, %%rax")
		c.println(".L.end.%d:", cnt)
		return
	case ND_FUNCALL:
		nargs := 0
		for arg := node.Args; arg != nil; arg = arg.Next {
			c.genExpr(arg)
			c.push()
			nargs++
		}

		for i := nargs - 1; i >= 0; i-- {
			c.pop(argreg64[i])
		}

		c.println("	mov $0, %%rax")
		c.println("	call %s", node.FuncName)
		return
	}

	c.genExpr(node.Rhs)
	c.push()
	c.genExpr(node.Lhs)
	c.pop("%rdi")

	var ax, di string

	if node.Lhs.Ty.Kind == TY_LONG || node.Lhs.Ty.Base != nil {
		ax = "%rax"
		di = "%rdi"
	} else {
		ax = "%eax"
		di = "%edi"
	}

	switch node.Kind {
	case ND_ADD:
		c.println("	add %s, %s", di, ax)
		return
	case ND_SUB:
		c.println("	sub %s, %s", di, ax)
		return
	case ND_MUL:
		c.println("	imul %s, %s", di, ax)
		return
	case ND_DIV, ND_MOD:
		if node.Lhs.Ty.Sz == 8 {
			c.println("	cqo")
		} else {
			c.println("	cdq")
		}
		c.println("	idiv %s", di)

		if node.Kind == ND_MOD {
			c.println("	mov %%rdx, %%rax")
		}
		return
	case ND_BITAND:
		c.println("	and %%rdi, %%rax")
		return
	case ND_BITOR:
		c.println("	or %%rdi, %%rax")
		return
	case ND_BITXOR:
		c.println("	xor %%rdi, %%rax")
		return
	case ND_EQ, ND_NE, ND_LT, ND_LE:
		c.println("	cmp %s, %s", di, ax)

		switch node.Kind {
		case ND_EQ:
			c.println("	sete %%al")
		case ND_NE:
			c.println("	setne %%al")
		case ND_LT:
			c.println("	setl %%al")
		case ND_LE:
			c.println("	setle %%al")
		}

		c.println("	movzb %%al, %%rax")
		return
	case ND_SHL:
		c.println("	mov %%rdi, %%rcx")
		c.println("	shl %%cl, %s", ax)
		return
	case ND_SHR:
		c.println("	mov %%rdi, %%rcx")
		if node.Ty.Sz == 8 {
			c.println("	sar %%cl, %s", ax)
		} else {
			c.println("	sar %%cl, %s", ax)
		}
		return
	}

	c.err = fmt.Errorf(errorTok(node.Tok, "invalid expression"))
}

func (c *codeWriter) genStmt(node *Node) {
	if c.err != nil {
		return
	}

	c.println("	.loc 1 %d", node.Tok.LineNo)

	switch node.Kind {
	case ND_IF:
		cnt := count()
		c.genExpr(node.Cond)
		c.println("	cmp $0, %%rax")
		c.println("	je .L.else.%d", cnt)
		c.genStmt(node.Then)
		c.println("	jmp .L.end.%d", cnt)
		c.println(".L.else.%d:", cnt)
		if node.Els != nil {
			c.genStmt(node.Els)
		}
		c.println(".L.end.%d:", cnt)
		return
	case ND_FOR:
		cnt := count()
		if node.Init != nil {
			c.genStmt(node.Init)
		}
		c.println(".L.begin.%d:", cnt)
		if node.Cond != nil {
			c.genExpr(node.Cond)
			c.println("	cmp $0, %%rax")
			c.println("	je %s", node.BrkLabel)
		}
		c.genStmt(node.Then)
		c.println("%s:", node.ContLabel)
		if node.Inc != nil {
			c.genExpr(node.Inc)
		}
		c.println("	jmp .L.begin.%d", cnt)
		c.println("%s:", node.BrkLabel)
		return
	case ND_SWITCH:
		c.genExpr(node.Cond)

		for n := node.CaseNext; n != nil; n = n.CaseNext {
			var reg string
			if node.Cond.Ty.Sz == 8 {
				reg = "%rax"
			} else {
				reg = "%eax"
			}
			c.println("	cmp $%d, %s", n.Val, reg)
			c.println("	je %s", n.Lbl)
		}

		if node.DefCase != nil {
			c.println("	jmp %s", node.DefCase.Lbl)
		}

		c.println("	jmp %s", node.BrkLabel)
		c.genStmt(node.Then)
		c.println("%s:", node.BrkLabel)
		return
	case ND_CASE:
		c.println("%s:", node.Lbl)
		c.genStmt(node.Lhs)
		return
	case ND_BLOCK:
		for n := node.Body; n != nil; n = n.Next {
			c.genStmt(n)
		}
		return
	case ND_GOTO:
		c.println("	jmp %s", node.UniqueLbl)
		return
	case ND_LABEL:
		c.println("%s:", node.UniqueLbl)
		c.genStmt(node.Lhs)
		return
	case ND_RETURN:
		c.genExpr(node.Lhs)
		c.println("	jmp .L.return.%s", curFnInGen.Name)
		return
	case ND_EXPR_STMT:
		c.genExpr(node.Lhs)
		return
	}

	c.err = fmt.Errorf(errorTok(node.Tok, "invalid statement"))
}

// Assign offsets to local variables
func (c *codeWriter) assignLvarOffsets(prog *Obj) {
	if c.err != nil {
		return
	}

	for fn := prog; fn != nil; fn = fn.Next {
		if !fn.IsFunc {
			continue
		}

		offset := 0
		for v := fn.Locals; v != nil; v = v.Next {
			offset += v.Ty.Sz
			offset = alignTo(offset, v.Ty.Align)
			v.Offset = -offset
		}
		fn.StackSz = alignTo(offset, 16)
	}
}

func (c *codeWriter) emitData(prog *Obj) {
	if c.err != nil {
		return
	}

	for v := prog; v != nil; v = v.Next {
		if v.IsFunc {
			continue
		}

		c.println("	.data")
		c.println("	.globl %s", v.Name)
		c.println("%s:", v.Name)

		if v.InitData != nil {
			for i := 0; i < len(v.InitData); i++ {
				c.println("	.byte %d", v.InitData[i])
			}
		} else {
			c.println("	.zero %d", v.Ty.Sz)
		}
	}
}

func (c *codeWriter) storeGp(r, offset, sz int) {
	if c.err != nil {
		return
	}

	switch sz {
	case 1:
		c.println("	mov %s, %d(%%rbp)", argreg8[r], offset)
		return
	case 2:
		c.println("	mov %s, %d(%%rbp)", argreg16[r], offset)
		return
	case 4:
		c.println("	mov %s, %d(%%rbp)", argreg32[r], offset)
		return
	case 8:
		c.println("	mov %s, %d(%%rbp)", argreg64[r], offset)
		return
	default:
		c.err = fmt.Errorf("internal error")
	}
}

func isImplicitFn(fnName string) bool {
	var imFn []string = []string{"printf", "exit", "assert"}
	for _, f := range imFn {
		if fnName == f {
			return true
		}
	}
	return false
}

func (c *codeWriter) emitText(prog *Obj) {
	if c.err != nil {
		return
	}

	for fn := prog; fn != nil; fn = fn.Next {
		if !fn.IsFunc || !fn.IsDef || isImplicitFn(fn.Name) {
			continue
		}

		c.println("	.globl %s", fn.Name)
		c.println("	.text")
		c.println("%s:", fn.Name)
		curFnInGen = fn

		// Prologue
		c.println("	push %%rbp")
		c.println("	mov %%rsp, %%rbp")
		c.println("	sub $%d, %%rsp", fn.StackSz)

		// Push arguments to the stack
		i := 0
		for v := fn.Params; v != nil; v = v.Next {
			c.storeGp(i, v.Offset, v.Ty.Sz)
			i++
		}

		// Emit code
		c.genStmt(fn.Body)
		if depth != 0 {
			c.err = fmt.Errorf("depth is not 0")
			return
		}

		// 'main' function returns implicitly 0.
		// [https://www.sigbus.info/n1570#5.1.2.2.3p1] The C spec defines
		// a special rule for the main function. Reaching the end of the
		// main function is equivalent to returning 0, evan though the
		// behavior is undifined for the other functions.
		if fn.Name == "main" {
			c.println("	mov $0, %%rax")
		}

		// Epilogue
		c.println(".L.return.%s:", fn.Name)
		c.println("	mov %%rbp, %%rsp")
		c.println("	pop %%rbp")
		c.println("	ret")
	}
}

func codegen(w io.Writer, prog *Obj) error {
	c := &codeWriter{w: w}

	c.assignLvarOffsets(prog)
	c.emitData(prog)
	c.emitText(prog)

	return c.err
}
