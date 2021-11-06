//
// code generator
//
package main

import (
	"errors"
	"fmt"
	"io"
)

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

var labelNo int
var argReg1 = []string{"dil", "sil", "dl", "cl", "r8b", "r9b"}
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
	if sizeOf(ty) == 1 {
		c.printf("	movsx rax, byte ptr [rax]\n")
	} else {
		c.printf("	mov rax, [rax]\n")
	}
	c.printf("	push rax\n")
}

func (c *codeWriter) store(ty *Type) {
	if c.err != nil {
		return // do nothing
	}

	c.printf("	pop rdi\n")
	c.printf("	pop rax\n")

	if sizeOf(ty) == 1 {
		c.printf("	mov [rax], dil\n")
	} else {
		c.printf("	mov [rax], rdi\n")
	}

	c.printf("	push rdi\n")
}

func (c *codeWriter) gen(node *Node) {
	if c.err != nil {
		return // do nothing
	}

	switch node.Kind {
	case ND_NULL:
		return
	case ND_NUM:
		c.printf("	push %d\n", node.Val)
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
		// store
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

	case ND_IF:
		c.gen(node.Cond)
		c.printf("	pop rax\n")
		c.printf("	cmp rax, 0\n")

		seq := labelNo
		labelNo++
		if node.Els != nil {
			c.printf("	je .Lelse%03d\n", seq)
		} else {
			c.printf("	je .Lend%03d\n", seq)
		}

		c.gen(node.Then)

		if node.Els != nil {
			c.printf(" jmp .Lend%03d\n", seq)
			c.printf(".Lelse%03d:\n", seq)
			c.gen(node.Els)
		}

		c.printf(".Lend%03d:\n", seq)
		return

	case ND_WHILE:
		seq := labelNo
		labelNo++
		c.printf(".Lbegin%03d:\n", seq)
		c.gen(node.Cond)
		c.printf("	pop rax\n")
		c.printf("	cmp rax, 0\n")
		c.printf("	je .Lend%03d\n", seq)

		c.gen(node.Then)

		c.printf("	jmp .Lbegin%03d\n", seq)
		c.printf(".Lend%03d:\n", seq)
		return

	case ND_FOR:
		if node.Init != nil {
			c.gen(node.Init)
		}

		seq := labelNo
		labelNo++
		c.printf(".Lbegin%03d:\n", seq)

		if node.Cond != nil {
			c.gen(node.Cond)
			c.printf("	pop rax\n")
			c.printf("	cmp rax, 0\n")
			c.printf("	je .Lend%03d\n", seq)
		}

		c.gen(node.Then)

		if node.Inc != nil {
			c.gen(node.Inc)
		}
		c.printf("	jmp .Lbegin%03d\n", seq)
		c.printf(".Lend%03d:\n", seq)
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
		c.printf("	mov rax, rsp\n")        // move rsp to rax
		c.printf("	and rax, 15\n")         // calculate rax & 15, when rax == 16, rax is 0b10000, and 15(0b1110) & 0b10000, ZF become 0.
		c.printf("	jnz .Lcall%03d\n", seq) // if ZF is 0, jamp to Lcall???.
		c.printf("	mov rax, 0\n")          // remove rax
		c.printf("	call %s\n", node.FuncName)
		c.printf("	jmp .Lend%03d\n", seq)
		c.printf(".Lcall%03d:\n", seq)
		c.printf("	sub rsp, 8\n") // rspは8の倍数なので16の倍数にするために8を引く
		c.printf("	mov rax, 0\n")
		c.printf("	call %s\n", node.FuncName)
		c.printf("	add rsp, 8\n")
		c.printf(".Lend%03d:\n", seq)
		c.printf("	push rax\n")
		return

	case ND_BLOCK, ND_STMT_EXPR:
		for n := node.Body; n != nil; n = n.Next {
			c.gen(n)
		}
		return

	case ND_RETURN:
		c.gen(node.Lhs)
		c.printf("	pop rax\n")
		c.printf("	jmp .Lreturn.%s\n", funcName)
		return
	}

	c.gen(node.Lhs)
	c.gen(node.Rhs)

	c.printf("	pop rdi\n")
	c.printf("	pop rax\n")

	switch node.Kind {
	case ND_ADD:
		if node.Ty.PtrTo != nil {
			c.printf("	imul rdi, %d\n", sizeOf(node.Ty.PtrTo))
		}
		c.printf("	add rax, rdi\n")
	case ND_SUB:
		if node.Ty.PtrTo != nil {
			c.printf("	imul rdi, %d\n", sizeOf(node.Ty.PtrTo))
		}
		c.printf("	sub rax, rdi\n")
	case ND_MUL:
		c.printf("	imul rax, rdi\n")
	case ND_DIV:
		c.printf("	cqo\n")
		c.printf("	idiv rdi\n")
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
			c.printf("	.zero %d\n", sizeOf(vl.Var.Ty))
			continue
		}

		for i := 0; i < vl.Var.ContLen; i++ {
			c.printf("	.byte %d\n", vl.Var.Contents[i])
		}
	}
}

func (c *codeWriter) loadArg(lvar *Var, idx int) {
	sz := sizeOf(lvar.Ty)
	if sz == 1 {
		c.printf("	mov [rbp-%d], %s\n", lvar.Offset, argReg1[idx])
	} else {
		if sz != 8 {
			c.err = errors.New("invalid size")
		}
		c.printf("	mov [rbp-%d], %s\n", lvar.Offset, argReg8[idx])
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
