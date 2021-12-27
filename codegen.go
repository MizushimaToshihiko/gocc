//
// code generator
//
package main

import (
	"fmt"
	"io"
	"math"
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

var argreg1 = []string{"dil", "sil", "dl", "cl", "r8b", "r9b"}
var argreg2 = []string{"di", "si", "dx", "cx", "r8w", "r9w"}
var argreg4 = []string{"edi", "esi", "edx", "ecx", "r8d", "r9d"}
var argreg8 = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

var labelseq int = 1
var brkseq int
var contseq int
var funcname string

// Pushes the given node's address to the stack
func (c *codeWriter) genAddr(node *Node) {
	if c.err != nil {
		return
	}

	switch node.Kind {
	case ND_VAR:
		if node.Obj.IsLocal {
			c.println("	lea rax, [rbp-%d]", node.Obj.Offset)
			c.println("	push rax")
			return
		}
		c.println("	push offset %s", node.Obj.Name)
		return
	case ND_DEREF:
		c.gen(node.Lhs)
		return
	case ND_MEMBER:
		c.genAddr(node.Lhs)
		c.println("	pop rax")
		c.println("	add rax, %d", node.Mem.Offset)
		c.println("	push rax")
		return
	default:
		c.err = fmt.Errorf(errorTok(node.Tok, "not an lvalue"))
	}

}

func (c *codeWriter) genLval(node *Node) {
	if c.err != nil {
		return
	}

	if node.Ty.Kind == TY_ARRAY {
		c.err = fmt.Errorf(errorTok(node.Tok, "not an lvalue"))
	}
	c.genAddr(node)
}

func (c *codeWriter) load(ty *Type) {
	if c.err != nil {
		return
	}

	c.println("	pop rax")
	switch sizeOf(ty, nil) {
	case 1:
		c.println("	movsx rax, byte ptr [rax]")
	case 2:
		c.println("	movsx rax, word ptr [rax]")
	case 4:
		c.println("	movsxd rax, dword ptr [rax]")
	case 8:
		c.println("	mov rax, [rax]")
	default:
		c.err = fmt.Errorf("invalid size")
		return
	}

	c.println("	push rax")
}

func (c *codeWriter) store(ty *Type) {
	if c.err != nil {
		return
	}

	c.println("	pop rdi")
	c.println("	pop rax")

	if ty.Kind == TY_BOOL {
		c.println("	cmp rdi, 0")
		c.println("	setne dil")
		c.println("	movzb rdi, dil")
	}

	switch sizeOf(ty, nil) {
	case 1:
		c.println("	mov [rax], dil")
	case 2:
		c.println("	mov [rax], di")
	case 4:
		c.println("	mov [rax], edi")
	case 8:
		c.println("	mov [rax], rdi")
	default:
		c.err = fmt.Errorf("invalid size")
	}

	c.println("	push rdi")
}

func (c *codeWriter) trancate(ty *Type) {
	if c.err != nil {
		return
	}

	c.println("	pop rax")

	if ty.Kind == TY_BOOL {
		c.println("	cmp rax, 0")
		c.println("	setne al")
	}

	switch sizeOf(ty, nil) {
	case 1:
		c.println("	movsx rax, al")
	case 2:
		c.println("	movsx rax, ax")
	case 4:
		c.println("	movsxd rax, eax")
	}
	c.println("	push rax")
}

func (c *codeWriter) inc(node *Node) {
	c.println("	pop rax")
	if node.Ty.Base != nil {
		c.println("	add rax, %d", sizeOf(node.Ty.Base, node.Tok))
		c.println("	push rax")
		return
	}
	c.println("	add rax, 1")
	c.println("	push rax")
}

func (c *codeWriter) dec(node *Node) {
	c.println("	pop rax")
	if node.Ty.Base != nil {
		c.println("	sub rax, %d", sizeOf(node.Ty.Base, node.Tok))
		c.println("	push rax")
		return
	}
	c.println("	sub rax, 1")
	c.println("	push rax")
}

func (c *codeWriter) gen(node *Node) (err error) {
	if c.err != nil {
		return
	}

	c.println("	.loc 1 %d", node.Tok.LineNo)

	switch node.Kind {
	case ND_NULL:
		return
	case ND_NUM:
		if node.Val <= int64(math.MaxInt32) { // node.Val is int32
			c.println("	push %d", node.Val)
		} else { // node.Val is int64
			c.println("	movabs rax, %d", node.Val)
			c.println("	push rax")
		}
		return
	case ND_EXPR_STMT:
		c.gen(node.Lhs)
		c.println("	add rsp, 8")
		return
	case ND_VAR, ND_MEMBER:
		c.genAddr(node)
		if node.Ty.Kind != TY_ARRAY {
			c.load(node.Ty)
		}
		return
	case ND_ASSIGN:
		c.genLval(node.Lhs)
		c.gen(node.Rhs)
		c.store(node.Ty)
		return
	case ND_INC:
		c.genLval(node.Lhs)
		c.println("	push [rsp]")
		c.load(node.Ty)
		c.inc(node)
		c.store(node.Ty)
		c.dec(node)
		return
	case ND_DEC:
		c.genLval(node.Lhs)
		c.println("	push [rsp]")
		c.load(node.Ty)
		c.dec(node)
		c.store(node.Ty)
		c.inc(node)
		return
	case ND_A_ADD, ND_A_SUB, ND_A_MUL, ND_A_DIV, ND_A_SHL, ND_A_SHR:
		c.genLval(node.Lhs)
		c.println("	push [rsp]")
		c.load(node.Lhs.Ty)
		c.gen(node.Rhs)
		c.println("	pop rdi")
		c.println("	pop rax")

		switch node.Kind {
		case ND_A_ADD:
			if node.Ty.Base != nil {
				c.println("	imul rdi, %d", sizeOf(node.Ty.Base, node.Tok))
			}
			c.println("	add rax, rdi")
		case ND_A_SUB:
			if node.Ty.Base != nil {
				c.println("	imul rdi, %d", sizeOf(node.Ty.Base, node.Tok))
			}
			c.println("	sub rax, rdi")
		case ND_A_MUL:
			c.println("	imul rax, rdi")
		case ND_A_DIV:
			c.println("	cqo")
			c.println("	idiv rdi")
		case ND_A_SHL:
			c.println("	mov cl, dil")
			c.println("	shl rax, cl")
		case ND_A_SHR:
			c.println("	mov cl, dil")
			c.println("	sar rax, cl")
		}

		c.println("	push rax")
		c.store(node.Ty)
		return
	case ND_ADDR:
		c.genAddr(node.Lhs)
		return
	case ND_DEREF:
		c.gen(node.Lhs)
		if node.Ty.Kind != TY_ARRAY {
			c.load(node.Ty)
		}
		return
	case ND_NOT:
		c.gen(node.Lhs)
		c.println("	pop rax")
		c.println("	cmp rax, 0")
		c.println("	sete al")
		c.println("	movzb rax, al")
		c.println("	push rax")
		return
	case ND_BITNOT:
		c.gen(node.Lhs)
		c.println("	pop rax")
		c.println("	not rax")
		c.println("	push rax")
		return
	case ND_LOGAND:
		seq := labelseq
		labelseq++
		c.gen(node.Lhs)
		c.println("	pop rax")
		c.println("	cmp rax, 0")
		c.println("	je  .Lfalse%d", seq)
		c.gen(node.Rhs)
		c.println("	pop rax")
		c.println("	cmp rax, 0")
		c.println("	je  .Lfalse%d", seq)
		c.println("	push 1")
		c.println("	jmp  .Lend%d", seq)
		c.println(".Lfalse%d:", seq)
		c.println("	push 0")
		c.println(".Lend%d:", seq)
		return
	case ND_LOGOR:
		seq := labelseq
		labelseq++
		c.gen(node.Lhs)
		c.println("	pop rax")
		c.println("	cmp rax, 0")
		c.println("	jne  .Ltrue%d", seq)
		c.gen(node.Rhs)
		c.println("	pop rax")
		c.println("	cmp rax, 0")
		c.println("	jne  .Ltrue%d", seq)
		c.println("	push 0")
		c.println("	jmp  .Lend%d", seq)
		c.println(".Ltrue%d:", seq)
		c.println("	push 1")
		c.println(".Lend%d:", seq)
		return
	case ND_IF:
		seq := labelseq
		labelseq++
		if node.Els != nil {
			c.gen(node.Cond)
			c.println("	pop rax")
			c.println("	cmp rax, 0")
			c.println("	je .Lelse%d", seq)
			c.gen(node.Then)
			c.println("	jmp .Lend%d", seq)
			c.println(".Lelse%d:", seq)
			c.gen(node.Els)
			c.println(".Lend%d:", seq)
			return
		}
		c.gen(node.Cond)
		c.println("	pop rax")
		c.println("	cmp rax, 0")
		c.println("	je .Lend%d", seq)
		c.gen(node.Then)
		c.println(".Lend%d:", seq)
		return
	case ND_WHILE:
		seq := labelseq
		labelseq++
		brk := brkseq
		cont := contseq
		contseq = seq
		brkseq = seq

		c.println(".L.continue.%d:", seq)
		c.gen(node.Cond)
		c.println("	pop rax")
		c.println("	cmp rax, 0")
		c.println("	je .L.break.%d", seq)
		c.gen(node.Then)
		c.println("	jmp .L.continue.%d", seq)
		c.println(".L.break.%d:", seq)

		brkseq = brk
		contseq = cont
		return
	case ND_FOR:
		seq := labelseq
		labelseq++
		brk := brkseq
		cont := contseq
		contseq = seq
		brkseq = seq

		if node.Init != nil {
			c.gen(node.Init)
		}
		c.println(".Lbegin%d:", seq)
		if node.Cond != nil {
			c.gen(node.Cond)
			c.println("	pop rax")
			c.println("	cmp rax, 0")
			c.println("	je .L.break.%d", seq)
		}
		c.gen(node.Then)
		c.println(".L.continue.%d:", seq)
		if node.Inc != nil {
			c.gen(node.Inc)
		}
		c.println("	jmp .Lbegin%d", seq)
		c.println(".L.break.%d:", seq)

		brkseq = brk
		contseq = cont
		return
	case ND_SWITCH:
		seq := labelseq
		labelseq++
		brk := brkseq
		brkseq = seq
		node.CaseLbl = seq

		c.gen(node.Cond)
		c.println("	pop rax")

		for n := node.CaseNext; n != nil; n = n.CaseNext {
			n.CaseLbl = labelseq
			labelseq++
			n.CaseEndLbl = seq
			c.println("	cmp rax, %d", n.Val)
			c.println("	je .L.case.%d", n.CaseLbl)
		}

		if node.DefCase != nil {
			i := labelseq
			labelseq++
			node.DefCase.CaseEndLbl = seq
			node.DefCase.CaseLbl = i
			c.println("	jmp .L.case.%d", i)
		}

		c.println("	jmp .L.break.%d", seq)
		c.gen(node.Then)
		c.println(".L.break.%d:", seq)

		brkseq = brk
		return
	case ND_CASE:
		c.println(".L.case.%d:", node.CaseLbl)
		c.gen(node.Lhs)
		c.println("	jmp .L.break.%d", node.CaseEndLbl)
		return
	case ND_BLOCK:
		for n := node.Body; n != nil; n = n.Next {
			c.gen(n)
		}
		return
	case ND_BREAK:
		if brkseq == 0 {
			c.err = fmt.Errorf(errorTok(node.Tok, "stray break"))
		}
		c.println("	jmp .L.break.%d", brkseq)
		return
	case ND_CONTINUE:
		if contseq == 0 {
			c.err = fmt.Errorf(errorTok(node.Tok, "stray continue"))
		}
		c.println("	jmp .L.continue.%d", contseq)
		return
	case ND_GOTO:
		c.println("	jmp .L.label.%s.%s", funcname, node.LblName)
		return
	case ND_LABEL:
		c.println(".L.label.%s.%s:", funcname, node.LblName)
		c.gen(node.Lhs)
		return
	case ND_FUNCALL:
		nargs := 0
		for arg := node.Args; arg != nil; arg = arg.Next {
			c.gen(arg)
			nargs++
		}

		for i := nargs - 1; i >= 0; i-- {
			c.println("	pop %s", argreg8[i])
		}

		// We need to align RSP to a 16 byte boundary before
		// calling a function because it is an ABI requirement.
		// RAX is set to 0 for variadic function.
		seq := labelseq
		labelseq++
		c.println("	mov rax, rsp")
		c.println("	and rax, 15")
		c.println("	jnz .Lcall%d", seq)
		c.println("	mov rax, 0")
		c.println("	call %s", node.FuncName)
		c.println("	jmp .Lend%d", seq)
		c.println(".Lcall%d:", seq)
		c.println("	sub rsp, 8")
		c.println("	mov rax, 0")
		c.println("	call %s", node.FuncName)
		c.println("	add rsp, 8")
		c.println(".Lend%d:", seq)
		c.println("	push rax")

		if node.Ty.Kind != TY_VOID {
			c.trancate(node.Ty)
		}
		return
	case ND_RETURN:
		c.gen(node.Lhs)
		c.println("	pop rax")
		c.println("	jmp .Lreturn.%s", funcname)
		return
	case ND_CAST:
		c.gen(node.Lhs)
		c.trancate(node.Ty)
		return
	}

	c.gen(node.Lhs)
	c.gen(node.Rhs)

	c.println("	pop rdi")
	c.println("	pop rax")

	switch node.Kind {
	case ND_ADD:
		if node.Ty.Base != nil {
			c.println("	imul rdi, %d", sizeOf(node.Ty.Base, node.Tok))
		}
		c.println("	add rax, rdi")
	case ND_SUB:
		if node.Ty.Base != nil {
			c.println("	imul rdi, %d", sizeOf(node.Ty.Base, node.Tok))
		}
		c.println("	sub rax, rdi")
	case ND_MUL:
		c.println("	imul rax, rdi")
	case ND_DIV:
		c.println("	cqo")
		c.println("	idiv rdi")
	case ND_BITAND:
		c.println("	and rax, rdi")
	case ND_BITOR:
		c.println("	or rax, rdi")
	case ND_BITXOR:
		c.println("	xor rax, rdi")
	case ND_SHL:
		c.println("	mov cl, dil")
		c.println("	shl rax, cl")
	case ND_SHR:
		c.println("	mov cl, dil")
		c.println("	sar rax, cl")
	case ND_EQ:
		c.println("	cmp rax, rdi")
		c.println("	sete al")
		c.println("	movzb rax, al")
	case ND_NE:
		c.println("	cmp rax, rdi")
		c.println("	setne al")
		c.println("	movzb rax, al")
	case ND_LT:
		c.println("	cmp rax, rdi")
		c.println("	setl al")
		c.println("	movzb rax, al")
	case ND_LE:
		c.println("	cmp rax, rdi")
		c.println("	setle al")
		c.println("	movzb rax, al")
	}

	c.println("	push rax")
	return
}

func (c *codeWriter) loadArg(v *Obj, idx int) {
	if c.err != nil {
		return
	}

	switch sizeOf(v.Ty, v.Tok) {
	case 1:
		c.println("	mov [rbp-%d], %s", v.Offset, argreg1[idx])
	case 2:
		c.println("	mov [rbp-%d], %s", v.Offset, argreg2[idx])
	case 4:
		c.println("	mov [rbp-%d], %s", v.Offset, argreg4[idx])
	case 8:
		c.println("	mov [rbp-%d], %s", v.Offset, argreg8[idx])
	default:
		c.err = fmt.Errorf("invalid size")
	}
}

// Assign offsets to local variables
func (c *codeWriter) assignLvarOffsets(prog *Program) {
	for fn := prog.Fns; fn != nil; fn = fn.Next {
		offset := 0
		for vl := fn.Locals; vl != nil; vl = vl.Next {
			offset = alignTo(offset, vl.Obj.Ty.Align)
			offset += sizeOf(vl.Obj.Ty, vl.Obj.Tok)
			vl.Obj.Offset = offset
		}
		fn.StackSz = alignTo(offset, 8)
	}
}

func (c *codeWriter) emitData(prog *Program) {
	if c.err != nil {
		return
	}

	for vl := prog.Globs; vl != nil; vl = vl.Next {
		c.println("	.globl %s", vl.Obj.Name)
		c.println("	.align %d", vl.Obj.Ty.Align)

		if vl.Obj.Init == nil {
			c.println("	.bss")
			c.println("%s:", vl.Obj.Name)
			c.println("	.zero %d", sizeOf(vl.Obj.Ty, vl.Obj.Tok))
			continue
		}

		c.println("	.data")
		c.println("%s:", vl.Obj.Name)

		for init := vl.Obj.Init; init != nil; init = init.Next {
			if init.Lbl != "" {
				c.println("	.quad %s", init.Lbl)
				continue
			}

			if init.Sz == 1 {
				c.println("	.byte %d", init.Val)
			} else {
				c.println("	.%dbyte %d", init.Sz, init.Val)
			}
		}
	}
}

func (c *codeWriter) emitText(prog *Program) {
	if c.err != nil {
		return
	}

	c.println(".text")

	for fn := prog.Fns; fn != nil; fn = fn.Next {
		c.println(".globl %s", fn.Name)
		c.println("%s:", fn.Name)
		funcname = fn.Name

		// Prologue
		c.println("	push rbp")
		c.println("	mov rbp, rsp")
		c.println("	sub rsp, %d", fn.StackSz)

		// Push arguments to the stack
		i := 0
		for vl := fn.Params; vl != nil; vl = vl.Next {
			c.loadArg(vl.Obj, i)
			i++
		}

		// Emit code
		for n := fn.Node; n != nil; n = n.Next {
			c.gen(n)
		}

		// 'main' function returns implicitly 0.
		if fn.Name == "main" {
			c.println("	mov rax, 0")
		}

		// Epilogue
		c.println(".Lreturn.%s:", funcname)
		c.println("	mov rsp, rbp")
		c.println("	pop rbp")
		c.println("	ret")
	}
}

func codegen(w io.Writer, prog *Program) error {
	c := &codeWriter{w: w}

	c.println(".intel_syntax noprefix")
	c.assignLvarOffsets(prog)
	c.emitData(prog)
	c.emitText(prog)

	return c.err
}
