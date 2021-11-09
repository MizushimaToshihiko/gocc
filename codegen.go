//
// code generator
//
package main

import (
	"errors"
	"fmt"
	"io"
	"math"
)

var ErrInvalidSize error = errors.New("invalid size")

type codeWriter struct {
	w   io.Writer
	err error
}

func (c *codeWriter) printf(format string, a ...interface{}) {
	if c.err != nil {
		return // do nothing
	}
	// Almost all error is from 'fmt.Fprintf'
	_, c.err = fmt.Fprintf(c.w, format, a...)
}

var labelNo int = 1
var brkseq int
var contseq int
var argReg1 = []string{"dil", "sil", "dl", "cl", "r8b", "r9b"}
var argReg2 = []string{"di", "si", "dx", "cx", "r8w", "r9w"}
var argReg4 = []string{"edi", "esi", "edx", "ecx", "r8d", "r9d"}
var argReg8 = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}
var funcName string

func (c *codeWriter) genAddr(node *Node) {
	if c.err != nil {
		return // do nothing
	}

	switch node.Kind {
	case ND_VAR:
		if node.Var.IsLocal {
			c.printf("	lea rax, [rbp-%d]\n", node.Var.Offset)
			c.printf("	push rax\n")
		} else {
			c.printf("	push offset %s\n", node.Var.Name)
		}
		return
	case ND_DEREF:
		c.gen(node.Lhs)
		return
	case ND_MEMBER:
		c.genAddr(node.Lhs)
		c.printf("	pop rax\n")
		c.printf("	add rax, %d\n", node.Mem.Offset)
		c.printf("	push rax\n")
		return
	default:
		c.err = fmt.Errorf(
			"c.genLval(): err:\n%s",
			errorTok(node.Tok, "the left value is not a variable"),
		)
	}
}

func (c *codeWriter) genLval(node *Node) {
	if c.err != nil {
		return // do nothing
	}

	if node.Ty.Kind == TY_ARRAY {
		c.err = fmt.Errorf(
			"c.genLval(): err:\n%s",
			errorTok(node.Tok, "the left value is not a variable"),
		)
	}
	c.genAddr(node)
}

func (c *codeWriter) load(ty *Type) {
	if c.err != nil {
		return // do nothing
	}

	c.printf("	pop rax\n")

	switch sizeOf(ty, nil) {
	case 1:
		c.printf("	movsx rax, byte ptr [rax]\n")
	case 2:
		c.printf("	movsx rax, word ptr [rax]\n")
	case 4:
		c.printf("	movsxd rax, dword ptr [rax]\n")
	case 8:
		c.printf("	mov rax, [rax]\n")
	default:
		c.err = ErrInvalidSize
	}

	c.printf("	push rax\n")
}

func (c *codeWriter) store(ty *Type) {
	if c.err != nil {
		return // do nothing
	}

	c.printf("	pop rdi\n")
	c.printf("	pop rax\n")

	if ty.Kind == TY_BOOL {
		c.printf("	cmp rdi, 0\n")
		c.printf("	setne dil\n")
		c.printf("	movzb rdi, dil\n")
	}

	switch sizeOf(ty, nil) {
	case 1:
		c.printf("	mov [rax], dil\n")
	case 2:
		c.printf("	mov [rax], di\n")
	case 4:
		c.printf("	mov [rax], edi\n")
	case 8:
		c.printf("	mov [rax], rdi\n")
	default:
		c.err = ErrInvalidSize
	}

	c.printf("	push rdi\n")
}

func (c *codeWriter) truncate(ty *Type) {
	c.printf("	pop rax\n")

	if ty.Kind == TY_BOOL {
		c.printf("	cmp rax, 0\n")
		c.printf("	setne al\n")
	}

	switch sizeOf(ty, nil) {
	case 1:
		c.printf("	movsx rax, al\n")
	case 2:
		c.printf("	movsx rax, ax\n")
	case 4:
		c.printf("	movsxd rax, eax\n")
	}

	c.printf("	push rax\n")
}

func (c *codeWriter) inc(node *Node) {
	var sz int = 1
	if node.Ty.PtrTo != nil {
		sz = sizeOf(node.Ty.PtrTo, node.Tok)
	}
	c.printf("	pop rax\n")
	c.printf("	add rax, %d\n", sz)
	c.printf("	push rax\n")
}

func (c *codeWriter) dec(node *Node) {
	var sz int = 1
	if node.Ty.PtrTo != nil {
		sz = sizeOf(node.Ty.PtrTo, node.Tok)
	}
	c.printf("	pop rax\n")
	c.printf("	sub rax, %d\n", sz)
	c.printf("	push rax\n")
}

func (c *codeWriter) gen(node *Node) {
	if c.err != nil {
		return // do nothing
	}

	switch node.Kind {
	case ND_NULL:
		return
	case ND_NUM:
		if node.Val <= int64(math.MaxInt32) {
			c.printf("	push %d\n", node.Val)
		} else {
			c.printf("	movabs rax, %d\n", node.Val)
			c.printf("	push rax\n")
		}
		return
	case ND_EXPR_STMT:
		c.gen(node.Lhs)
		c.printf("	add rsp, 8\n")
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

	case ND_PRE_INC:
		c.genLval(node.Lhs)
		c.printf("	push [rsp]\n")
		c.load(node.Ty)
		c.inc(node)
		c.store(node.Ty)
		return

	case ND_PRE_DEC:
		c.genLval(node.Lhs)
		c.printf("	push [rsp]\n")
		c.load(node.Ty)
		c.dec(node)
		c.store(node.Ty)
		return

	case ND_POST_INC:
		c.genLval(node.Lhs)
		c.printf("	push [rsp]\n")
		c.load(node.Ty)
		c.inc(node)
		c.store(node.Ty)
		c.dec(node)
		return

	case ND_POST_DEC:
		c.genLval(node.Lhs)
		c.printf("	push [rsp]\n")
		c.load(node.Ty)
		c.dec(node)
		c.store(node.Ty)
		c.inc(node)
		return

	case ND_A_ADD, ND_A_SUB, ND_A_MUL, ND_A_DIV:
		c.genLval(node.Lhs)
		c.printf("	push [rsp]\n")
		c.load(node.Lhs.Ty)
		c.gen(node.Rhs)
		c.printf("	pop rdi\n")
		c.printf("	pop rax\n")

		switch node.Kind {
		case ND_A_ADD:
			if node.Ty.PtrTo != nil {
				c.printf("	imul rdi, %d\n", sizeOf(node.Ty.PtrTo, node.Tok))
			}
			c.printf("	add rax, rdi\n")
		case ND_A_SUB:
			if node.Ty.PtrTo != nil {
				c.printf("	imul rdi, %d\n", sizeOf(node.Ty.PtrTo, node.Tok))
			}
			c.printf("	sub rax, rdi\n")
		case ND_A_MUL:
			c.printf("	imul rax, rdi\n")
		case ND_A_DIV:
			c.printf("	cqo\n")
			c.printf("	idiv rdi\n")
		}

		c.printf("	push rax\n")
		c.store(node.Ty)
		return

	case ND_COMMA:
		c.gen(node.Lhs)
		c.gen(node.Rhs)
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
		c.printf("	pop rax\n")
		c.printf("	cmp rax, 0\n")
		c.printf("	sete al\n")
		c.printf("	movzb rax, al\n")
		c.printf("	push rax\n")
		return

	case ND_BITNOT:
		c.gen(node.Lhs)
		c.printf("	pop rax\n")
		c.printf("	not rax\n")
		c.printf("	push rax\n")
		return

	case ND_LOGAND:
		seq := labelNo
		labelNo++
		c.gen(node.Lhs)
		c.printf("	pop rax\n")
		c.printf("	cmp rax, 0\n")
		c.printf("	je .Lfalse%d\n", seq)
		c.gen(node.Rhs)
		c.printf("	pop rax\n")
		c.printf("	cmp rax, 0\n")
		c.printf("	je .Lfalse%d\n", seq)
		c.printf("	push 1\n")
		c.printf("	jmp .Lend%d\n", seq)
		c.printf(".Lfalse%d:\n", seq)
		c.printf("	push 0\n")
		c.printf(".Lend%d:\n", seq)
		return

	case ND_LOGOR:
		seq := labelNo
		labelNo++
		c.gen(node.Lhs)
		c.printf("	pop rax\n")
		c.printf("	cmp rax, 0\n")
		c.printf("	jne .Ltrue%d\n", seq)
		c.gen(node.Rhs)
		c.printf("	pop rax\n")
		c.printf("	cmp rax, 0\n")
		c.printf("	jne .Ltrue%d\n", seq)
		c.printf("	push 0\n")
		c.printf("	jmp .Lend%d\n", seq)
		c.printf(".Ltrue%d:\n", seq)
		c.printf("	push 1\n")
		c.printf(".Lend%d:\n", seq)
		return

	case ND_IF:
		c.gen(node.Cond)
		c.printf("	pop rax\n")
		c.printf("	cmp rax, 0\n")

		seq := labelNo
		labelNo++
		if node.Els != nil {
			c.printf("	je .Lelse%d\n", seq)
		} else {
			c.printf("	je .Lend%d\n", seq)
		}

		c.gen(node.Then)

		if node.Els != nil {
			c.printf(" jmp .Lend%d\n", seq)
			c.printf(".Lelse%d:\n", seq)
			c.gen(node.Els)
		}

		c.printf(".Lend%d:\n", seq)
		return

	case ND_WHILE:
		seq := labelNo
		labelNo++
		brk := brkseq
		cont := contseq
		brkseq = seq
		contseq = seq
		c.printf(".L.continue.%d:\n", seq)
		c.gen(node.Cond)
		c.printf("	pop rax\n")
		c.printf("	cmp rax, 0\n")
		c.printf("	je .L.break.%d\n", seq)

		c.gen(node.Then)

		c.printf("	jmp .L.continue.%d\n", seq)
		c.printf(".L.break.%d:\n", seq)

		brkseq = brk
		contseq = cont
		return

	case ND_FOR:
		seq := labelNo
		labelNo++
		brk := brkseq
		cont := contseq
		brkseq = seq
		contseq = seq

		if node.Init != nil {
			c.gen(node.Init)
		}
		c.printf(".Lbegin%d:\n", seq)

		if node.Cond != nil {
			c.gen(node.Cond)
			c.printf("	pop rax\n")
			c.printf("	cmp rax, 0\n")
			c.printf("	je .L.break.%d\n", seq)
		}

		c.gen(node.Then)
		c.printf(".L.continue.%d:\n", seq)
		if node.Inc != nil {
			c.gen(node.Inc)
		}
		c.printf("	jmp .Lbegin%d\n", seq)
		c.printf(".L.break.%d:\n", seq)

		brkseq = brk
		contseq = cont
		return

	case ND_CONTINUE:
		if contseq == 0 {
			c.err = fmt.Errorf(
				"c.gen(): err:\n%s",
				errorTok(node.Tok, "stray continue"),
			)
		}
		c.printf("	jmp .L.continue.%d\n", contseq)
		return

	case ND_GOTO:
		c.printf("	jmp .L.label.%s.%s\n", funcName, node.LabelName)
		return

	case ND_LABEL:
		c.printf(".L.label.%s.%s:\n", funcName, node.LabelName)
		c.gen(node.Lhs)
		return

	case ND_FUNCCALL:
		numArgs := 0
		for arg := node.Args; arg != nil; arg = arg.Next {
			c.gen(arg)
			numArgs++
		}

		for i := numArgs - 1; i >= 0; i-- {
			c.printf("	pop %s\n", argReg8[i])
		}

		seq := labelNo
		labelNo++
		c.printf("	mov rax, rsp\n")      // move rsp to rax
		c.printf("	and rax, 15\n")       // calculate rax & 15, when rax == 16, rax is 0b10000, and 15(0b1110) & 0b10000, ZF become 0.
		c.printf("	jnz .Lcall%d\n", seq) // if ZF is 0, jamp to Lcall???.
		c.printf("	mov rax, 0\n")        // remove rax
		c.printf("	call %s\n", node.FuncName)
		c.printf("	jmp .Lend%d\n", seq)
		c.printf(".Lcall%d:\n", seq)
		c.printf("	sub rsp, 8\n") // rspは8の倍数なので16の倍数にするために8を引く
		c.printf("	mov rax, 0\n")
		c.printf("	call %s\n", node.FuncName)
		c.printf("	add rsp, 8\n")
		c.printf(".Lend%d:\n", seq)
		c.printf("	push rax\n")

		c.truncate(node.Ty)
		return

	case ND_SWITCH:
		seq := labelNo
		labelNo++
		brk := brkseq
		brkseq = seq
		node.CaseLbl = seq

		c.gen(node.Cond)
		c.printf("	pop rax\n")

		for n := node.CaseNext; n != nil; n = n.CaseNext {
			n.CaseLbl = labelNo
			labelNo++
			n.CaseEndLbl = seq
			c.printf("	cmp rax, %d\n", n.Val)
			c.printf("	je .L.case.%d\n", n.CaseLbl)
		}

		if node.DefCase != nil {
			i := labelNo
			labelNo++
			node.DefCase.CaseEndLbl = seq
			node.DefCase.CaseLbl = i
			c.printf("	jmp .L.case.%d\n", i)
		}

		c.printf("	jmp .L.break.%d\n", seq)
		c.gen(node.Then)
		c.printf(".L.break.%d:\n", seq)

		brkseq = brk
		return

	case ND_CASE:
		c.printf(".L.case.%d:\n", node.CaseLbl)
		c.gen(node.Lhs)
		c.printf("	jmp .L.break.%d\n", node.CaseEndLbl)
		return

	case ND_BLOCK, ND_STMT_EXPR:
		for n := node.Body; n != nil; n = n.Next {
			c.gen(n)
		}
		return

	case ND_BREAK:
		if brkseq == 0 {
			c.err = fmt.Errorf(
				"c.gen(): err:\n%s",
				errorTok(node.Tok, "stray break"),
			)
		}
		c.printf("	jmp .L.break.%d\n", brkseq)
		return
	case ND_RETURN:
		c.gen(node.Lhs)
		c.printf("	pop rax\n")
		c.printf("	jmp .Lreturn.%s\n", funcName)
		return

	case ND_CAST:
		c.gen(node.Lhs)
		c.truncate(node.Ty)
		return
	}

	c.gen(node.Lhs)
	c.gen(node.Rhs)

	c.printf("	pop rdi\n")
	c.printf("	pop rax\n")

	switch node.Kind {
	case ND_ADD:
		if node.Ty.PtrTo != nil {
			c.printf("	imul rdi, %d\n", sizeOf(node.Ty.PtrTo, node.Tok))
		}
		c.printf("	add rax, rdi\n")
	case ND_SUB:
		if node.Ty.PtrTo != nil {
			c.printf("	imul rdi, %d\n", sizeOf(node.Ty.PtrTo, node.Tok))
		}
		c.printf("	sub rax, rdi\n")
	case ND_MUL:
		c.printf("	imul rax, rdi\n")
	case ND_DIV:
		c.printf("	cqo\n")
		c.printf("	idiv rdi\n")
	case ND_BITAND:
		c.printf("	and rax, rdi\n")
	case ND_BITOR:
		c.printf("	or rax, rdi\n")
	case ND_BITXOR:
		c.printf("	xor rax, rdi\n")
	case ND_EQ:
		c.printf("	cmp rax, rdi\n")
		c.printf("	sete al\n")
		c.printf("	movzb rax, al\n")
	case ND_NE:
		c.printf("	cmp rax, rdi\n")
		c.printf("	setne al\n")
		c.printf("	movzb rax, al\n")
	case ND_LT:
		c.printf("	cmp rax, rdi\n")
		c.printf("	setl al\n")
		c.printf("	movzb rax, al\n")
	case ND_LE:
		c.printf("	cmp rax, rdi\n")
		c.printf("	setle al\n")
		c.printf("	movzb rax, al\n")
	}

	c.printf("	push rax\n")
}

func (c *codeWriter) emitData(prog *Program) {
	c.printf(".data\n")

	for vl := prog.Globals; vl != nil; vl = vl.Next {
		c.printf("%s:\n", vl.Var.Name)

		if vl.Var.Contents == nil {
			c.printf("	.zero %d\n", sizeOf(vl.Var.Ty, vl.Var.Tok))
			continue
		}

		for i := 0; i < vl.Var.ContLen; i++ {
			c.printf("	.byte %d\n", vl.Var.Contents[i])
		}
	}
}

func (c *codeWriter) loadArg(lvar *Var, idx int) {
	switch sizeOf(lvar.Ty, lvar.Tok) {
	case 1:
		c.printf("	mov [rbp-%d], %s\n", lvar.Offset, argReg1[idx])
	case 2:
		c.printf("	mov [rbp-%d], %s\n", lvar.Offset, argReg2[idx])
	case 4:
		c.printf("	mov [rbp-%d], %s\n", lvar.Offset, argReg4[idx])
	case 8:
		c.printf("	mov [rbp-%d], %s\n", lvar.Offset, argReg8[idx])
	default:
		c.err = ErrInvalidSize
	}
}

func (c *codeWriter) emitText(prog *Program) {
	c.printf(".text\n")

	for fn := prog.Fns; fn != nil; fn = fn.Next {
		c.printf(".global %s\n", fn.Name)
		c.printf("%s:\n", fn.Name)
		funcName = fn.Name

		// prologue
		// secure an area for the stack size of 'fn'
		c.printf("	push rbp\n")
		c.printf("	mov rbp, rsp\n")
		c.printf("	sub rsp, %d\n", fn.StackSz)

		// push arguments to the stack
		i := 0
		for vl := fn.Params; vl != nil; vl = vl.Next {
			c.loadArg(vl.Var, i)
			i++
		}

		// emit code
		for node := fn.Node; node != nil; node = node.Next {
			c.gen(node)
		}
		// epilogue
		// the result of the expression is in 'rax',
		// and it is the return value
		c.printf(".Lreturn.%s:\n", funcName)
		c.printf("	mov rsp, rbp\n")
		c.printf("	pop rbp\n")
		c.printf("	ret\n")
	}
}

func codeGen(w io.Writer, prog *Program) error {
	c := &codeWriter{w: w}
	// output the former 3 lines of the assembly
	c.printf(".intel_syntax noprefix\n")
	c.emitData(prog)
	c.emitText(prog)

	return c.err
}
