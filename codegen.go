//
// code generator
//
package main

import (
	"fmt"
	"io"
	"os"
	"unsafe"
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
		if c.err == nil {
			c.err = fmt.Errorf(errorTok(node.Tok, "not an lvalue"))
		} else {
			c.err = fmt.Errorf(c.err.Error() + "\n" + errorTok(node.Tok, "not an lvalue"))
		}
	}

}

func (c *codeWriter) load(ty *Type) {
	if c.err != nil {
		return
	}

	switch ty.Kind {
	case TY_ARRAY, TY_STRUCT:
		// If it is an array, do not attempt to load a value to the
		// register because in general we can't load an entire array to a
		// register. As a result, the result of an evaluation of an array
		// become not the array itself but the address of the array.
		// This is where "array is automatically converted to a pointer to
		// the first element of the array in C" occurs.
		return
	case TY_FLOAT:
		c.println("	movss (%%rax), %%xmm0")
		return
	case TY_DOUBLE:
		c.println("	movsd (%%rax), %%xmm0")
		return
	}

	var insn string
	if ty.IsUnsigned {
		insn = "movz"
	} else {
		insn = "movs"
	}

	// When we load a char or a short value to a register, we always
	// extend them to the size of int, so we can assume the lower half of a
	// register for char, short and int may contain garbage. When we load
	// a long value to a register, it simply occupies the entire register.
	switch ty.Sz {
	case 1:
		c.println("	%sbl (%%rax), %%eax", insn)
	case 2:
		c.println("	%swl (%%rax), %%eax", insn)
	case 4:
		c.println("	movsxd (%%rax), %%rax")
	case 8:
		c.println("	mov (%%rax), %%rax")
	default:
		if c.err == nil {
			c.err = fmt.Errorf("invalid size")
		} else {
			c.err = fmt.Errorf(c.err.Error() + "\ninvalid size")
		}
		return
	}
}

func (c *codeWriter) store(ty *Type) {
	if c.err != nil {
		return
	}

	c.pop("%rdi")

	switch ty.Kind {
	case TY_STRUCT:
		for i := 0; i < ty.Sz; i++ {
			c.println("	mov %d(%%rax), %%r8b", i)
			c.println("	mov %%r8b, %d(%%rdi)", i)
		}
		return
	case TY_FLOAT:
		c.println("	movss %%xmm0, (%%rdi)")
		return
	case TY_DOUBLE:
		c.println("	movsd %%xmm0, (%%rdi)")
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
		if c.err == nil {
			c.err = fmt.Errorf("invalid size")
		} else {
			c.err = fmt.Errorf(c.err.Error() + "\ninvalid size")
		}
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
	U8
	U16
	U32
	U64
	F32
	F64
)

func (c *codeWriter) getTypeId(ty *Type) int {
	switch ty.Kind {
	case TY_BYTE:
		if ty.IsUnsigned {
			return U8
		}
		return I8
	case TY_SHORT:
		if ty.IsUnsigned {
			return U16
		}
		return I16
	case TY_INT:
		if ty.IsUnsigned {
			return U32
		}
		return I32
	case TY_LONG:
		if ty.IsUnsigned {
			return U64
		}
		return I64
	case TY_FLOAT:
		return F32
	case TY_DOUBLE:
		return F64
	default:
		return U64
	}
}

// The table for type casts
const (
	i32i8  string = "movsbl %al, %eax"
	i32u8  string = "movzbl %al, %eax"
	i32i16 string = "movswl %ax, %eax"
	i32u16 string = "movzwl %ax, %eax"
	i32f32 string = "cvtsi2ssl %eax, %xmm0"
	i32i64 string = "movsxd %eax, %rax"
	i32f64 string = "cvtsi2sdl %eax, %xmm0"

	u32f32 string = "mov %eax, %eax; cvtsi2ssq %rax, %xmm0"
	u32i64 string = "mov %eax, %eax"
	u32f64 string = "mov %eax, %eax; cvtsi2sdq %rax, %xmm0"

	i64f32 string = "cvtsi2ssq %rax, %xmm0"
	i64f64 string = "cvtsi2sdq %rax, %xmm0"

	u64f32 string = "cvtsi2ssq %rax, %xmm0"
	u64f64 string = "test %rax,%rax; js 1f; pxor %xmm0,%xmm0; cvtsi2sd %rax,%xmm0; jmp 2f; " +
		"1: mov %rax,%rdi; and $1,%eax; pxor %xmm0,%xmm0; shr %rdi; " +
		"or %rax,%rdi; cvtsi2sd %rdi,%xmm0; addsd %xmm0,xmm0; 2:"

	f32i8  string = "cvttss2sil %xmm0, %eax; movsbl %al, %eax"
	f32u8  string = "cvttss2sil %xmm0, %eax; movzbl %al, %eax"
	f32i16 string = "cvttss2sil %xmm0, %eax; movswl %ax, %eax"
	f32u16 string = "cvttss2sil %xmm0, %eax; movzwl %ax, %eax"
	f32i32 string = "cvttss2sil %xmm0, %eax"
	f32u32 string = "cvttss2siq %xmm0, %rax"
	f32i64 string = "cvttss2siq %xmm0, %rax"
	f32u64 string = "cvttss2siq %xmm0, %rax"
	f32f64 string = "cvtss2sd %xmm0, %xmm0"

	f64i8  string = "cvttsd2sil %xmm0, %eax; movsbl %al, %eax"
	f64u8  string = "cvttsd2sil %xmm0, %eax; movzbl %al, %eax"
	f64i16 string = "cvttsd2sil %xmm0, %eax; movswl %ax, %eax"
	f64u16 string = "cvttsd2sil %xmm0, %eax; movzwl %ax, %eax"
	f64i32 string = "cvttsd2sil %xmm0, %eax"
	f64u32 string = "cvttsd2siq %xmm0, %rax"
	f64f32 string = "cvtsd2ss %xmm0, %xmm0"
	f64i64 string = "cvttsd2siq %xmm0, %rax"
	f64u64 string = "cvttsd2siq %xmm0, %rax"
)

var castTable = [10][10]string{
	//i8    i16     i32    i64     u8     u16     u32     u64     f32     f64
	{"nil", "null", "nil", i32i64, i32u8, i32u16, "null", i32i64, i32f32, i32f64}, // i8
	{i32i8, "null", "nil", i32i64, i32u8, i32u16, "null", i32i64, i32f32, i32f64}, // i16
	{i32i8, i32i16, "nil", i32i64, i32u8, i32u16, "null", i32i64, i32f32, i32f64}, // i32
	{i32i8, i32i16, "nil", "null", i32u8, i32u16, "null", "null", i64f32, i64f64}, // i64

	{i32i8, "null", "nil", i32i64, "nil", "null", "null", i32i64, i32f32, i32f64}, // u8
	{i32i8, i32i16, "nil", i32i64, i32u8, "null", "null", i32i64, i32f32, i32f64}, // u16
	{i32i8, i32i16, "nil", u32i64, i32u8, i32u16, "null", u32i64, u32f32, u32f64}, // u32
	{i32i8, i32i16, "nil", "null", i32u8, i32u16, "null", "null", u64f32, u64f64}, // u64

	{f32i8, f32i16, f32i32, f32i64, f32u8, f32u16, f32u32, f32u64, "null", f32f64}, // f32
	{f64i8, f64i16, f64i32, f64i64, f64u8, f64u16, f64i32, f64u64, f64f32, "null"}, // f64
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

	t1 := c.getTypeId(from)
	t2 := c.getTypeId(to)
	if castTable[t1][t2] != "nil" && castTable[t1][t2] != "null" {
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
		switch node.Ty.Kind {
		case TY_FLOAT:
			f32 := node.FVal
			c.println("	mov $%d, %%eax  # float %f", *(*uint32)(unsafe.Pointer(&f32)), f32)
			c.println("	movq %%rax, %%xmm0")
			return
		case TY_DOUBLE:
			f64 := node.FVal
			c.println("	mov $%d, %%rax  # double %f", *(*uint64)(unsafe.Pointer(&f64)), f64)
			c.println("	movq %%rax, %%xmm0")
			return
		}

		c.println("	mov $%d, %%rax", node.Val)
		return
	case ND_NEG:
		c.genExpr(node.Lhs)
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
		c.println("	lea %d(%%rbp), %%rdi", node.Obj.Offset)
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

		if depth%2 == 0 {
			c.println("	call %s", node.FuncName)
		} else {
			c.println("	sub $8, %%rsp")
			c.println("	call %s", node.FuncName)
			c.println("	add $8, %%rsp")
		}

		// It looks like the most significant 48 or 56 bits int RAX may
		// contain garbage if a function return type is short or bool/char,
		// respectively. We clear the upper bits here.
		switch node.Ty.Kind {
		case TY_BOOL:
			c.println("	movzx %%al, %%eax")
			return
		case TY_BYTE:
			if node.Ty.IsUnsigned {
				c.println("	movzbl %%al, %%eax")
				return
			}
			c.println("	movsbl %%al, %%eax")
			return
		case TY_SHORT:
			if node.Ty.IsUnsigned {
				c.println("	movzwl %%ax, %%eax")
				return
			}
			c.println("	movswl %%ax, %%eax")
			return
		}
		return
	}

	c.genExpr(node.Rhs)
	c.push()
	c.genExpr(node.Lhs)
	c.pop("%rdi")

	var ax, di, dx string

	if node.Lhs.Ty.Kind == TY_LONG || node.Lhs.Ty.Base != nil {
		ax = "%rax"
		di = "%rdi"
		dx = "%rdi"
	} else {
		ax = "%eax"
		di = "%edi"
		dx = "%edx"
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
		if node.Ty.IsUnsigned {
			c.println("	mov $0, %s", dx)
			c.println("	div %s", di)
		} else {
			if node.Lhs.Ty.Sz == 8 {
				c.println("	cqo")
			} else {
				c.println("	cdq")
			}
			c.println("	idiv %s", di)
		}

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
			if node.Lhs.Ty.IsUnsigned {
				c.println("	setb %%al")
			} else {
				c.println("	setl %%al")
			}
		case ND_LE:
			if node.Lhs.Ty.IsUnsigned {
				c.println("	setbe %%al")
			} else {
				c.println("	setle %%al")
			}
		}

		c.println("	movzb %%al, %%rax")
		return
	case ND_SHL:
		c.println("	mov %%rdi, %%rcx")
		c.println("	shl %%cl, %s", ax)
		return
	case ND_SHR:
		c.println("	mov %%rdi, %%rcx")
		if node.Lhs.Ty.IsUnsigned {
			c.println("	shr %%cl, %s", ax)
		} else {
			c.println("	sar %%cl, %s", ax)
		}
		return
	}

	if c.err == nil {
		c.err = fmt.Errorf("invalid expression")
	} else {
		c.err = fmt.Errorf(c.err.Error() + "\ninvalid expression")
	}
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
			n.CaseEndLbl = node.BrkLabel
		}

		if node.DefCase != nil {
			node.DefCase.CaseEndLbl = node.BrkLabel
			c.println("	jmp %s", node.DefCase.Lbl)
		}

		c.println("	jmp %s", node.BrkLabel)
		c.genStmt(node.Then)
		c.println("%s:", node.BrkLabel)
		return
	case ND_CASE:
		c.println("%s:", node.Lbl)
		c.genStmt(node.Lhs)
		c.println("	jmp %s", node.CaseEndLbl)
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
		if node.Lhs != nil {
			c.genExpr(node.Lhs)
		}
		c.println("	jmp .L.return.%s", curFnInGen.Name)
		return
	case ND_EXPR_STMT:
		c.genExpr(node.Lhs)
		return
	}
	if c.err == nil {
		c.err = fmt.Errorf(errorTok(node.Tok, "invalid statement"))
	} else {
		c.err = fmt.Errorf(c.err.Error() + "\n" + errorTok(node.Tok, "invalid statement"))
	}
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
			offset = alignTo(offset, v.Align)
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
		if v.IsFunc || !v.IsDef {
			continue
		}

		c.println("	.globl %s", v.Name)
		c.println("	.align %d", v.Align)

		if v.InitData != nil {
			c.println("	.data")
			c.println("%s:", v.Name)

			rel := v.Rel
			pos := 0
			for pos < v.Ty.Sz {
				if rel != nil && rel.Offset == pos {
					c.println("	.quad %s%+d", rel.Lbl, rel.Addend)
					rel = rel.Next
					pos += 8
					continue
				} else {
					c.println("	.byte %d", v.InitData[pos])
					pos++
				}
			}
			continue
		}

		c.println("	.bss")
		c.println("%s:", v.Name)
		c.println("	.zero %d", v.Ty.Sz)
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
		if c.err == nil {
			c.err = fmt.Errorf("internal error")
		} else {
			c.err = fmt.Errorf(c.err.Error() + "\ninternal error")
		}
	}
}

func (c *codeWriter) emitText(prog *Obj) {
	if c.err != nil {
		return
	}

	for fn := prog; fn != nil; fn = fn.Next {
		if !fn.IsFunc || !fn.IsDef {
			continue
		}

		if fn.Name != "main" && fn.IsStatic {
			c.println("	.local %s", fn.Name)
		} else {
			c.println("	.globl %s", fn.Name)
		}

		c.println("	.text")
		c.println("%s:", fn.Name)
		curFnInGen = fn

		// Prologue
		c.println("	push %%rbp")
		c.println("	mov %%rsp, %%rbp")
		c.println("	sub $%d, %%rsp", int(fn.StackSz))

		// Push arguments to the stack
		i := 0
		for v := fn.Params; v != nil; v = v.Next {
			c.storeGp(i, v.Offset, v.Ty.Sz)
			i++
		}

		// Emit code
		c.genStmt(fn.Body)
		if depth != 0 {
			if c.err == nil {
				c.err = fmt.Errorf("expected depth is 0, but %d", depth)
			} else {
				c.err = fmt.Errorf(c.err.Error()+"\nexpected depth is 0, but %d", depth)
			}
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
