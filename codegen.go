//
// code generator
//
package main

import (
	"errors"
	"fmt"
	"io"
)

// struct errWriter is for the error handling
// it's based on:
// https://jxck.hatenablog.com/entry/golang-error-handling-lesson-by-rob-pike
type codeWriter struct {
	w   io.Writer
	err error
}

func (c *codeWriter) Fprintf(format string, a ...interface{}) {
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
			c.Fprintf("	lea rax, [rbp-%d]\n", node.Var.Offset)
			c.Fprintf("	push rax\n")
		} else {
			c.Fprintf("	push offset %s\n", node.Var.Name)
		}
		return
	case ND_DEREF:
		c.gen(node.Lhs)
		return
	default:
		c.err = fmt.Errorf(
			"c.genLval(): err:\n%s",
			errorTok(node.Tok, "the left value is not a variable"),
		)
	}
}

func (c *codeWriter) genLval(node *Node) {
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

	c.Fprintf("	pop rax\n")
	if sizeOf(ty) == 1 {
		c.Fprintf("	movsx rax, byte ptr [rax]\n")
	} else {
		c.Fprintf("	mov rax, [rax]\n")
	}
	c.Fprintf("	push rax\n")
}

func (c *codeWriter) store(ty *Type) {
	if c.err != nil {
		return // do nothing
	}

	c.Fprintf("	pop rdi\n")
	c.Fprintf("	pop rax\n")

	if sizeOf(ty) == 1 {
		c.Fprintf("	mov [rax], dil\n")
	} else {
		c.Fprintf("	mov [rax], rdi\n")
	}

	c.Fprintf("	push rdi\n")
}

func (c *codeWriter) gen(node *Node) {
	if c.err != nil {
		return // do nothing
	}

	switch node.Kind {
	case ND_NULL:
		return
	case ND_NUM:
		c.Fprintf("	push %d\n", node.Val)
		return
	case ND_EXPR_STMT:
		c.gen(node.Lhs)
		c.Fprintf("	add rsp, 8\n")
		return
	case ND_VAR:
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
		c.Fprintf("	pop rax\n")
		c.Fprintf("	cmp rax, 0\n")

		labelNo++
		if node.Els != nil {
			c.Fprintf("	je .Lelse%03d\n", labelNo)
		} else {
			c.Fprintf("	je .Lend%03d\n", labelNo)
		}

		c.gen(node.Then)

		if node.Els != nil {
			c.Fprintf(" jmp .Lend%03d\n", labelNo)
			c.Fprintf(".Lelse%03d:\n", labelNo)
			c.gen(node.Els)
		}

		c.Fprintf(".Lend%03d:\n", labelNo)
		return

	case ND_WHILE:
		labelNo++
		c.Fprintf(".Lbegin%03d:\n", labelNo)
		c.gen(node.Cond)
		c.Fprintf("	pop rax\n")
		c.Fprintf("	cmp rax, 0\n")
		c.Fprintf("	je .Lend%03d\n", labelNo)

		c.gen(node.Then)

		c.Fprintf("	jmp .Lbegin%03d\n", labelNo)
		c.Fprintf(".Lend%03d:\n", labelNo)
		return

	case ND_FOR:
		if node.Init != nil {
			c.gen(node.Init)
		}

		labelNo++
		c.Fprintf(".Lbegin%03d:\n", labelNo)

		if node.Cond != nil {
			c.gen(node.Cond)
			c.Fprintf("	pop rax\n")
			c.Fprintf("	cmp rax, 0\n")
			c.Fprintf("	je .Lend%03d\n", labelNo)
		}

		c.gen(node.Then)

		if node.Inc != nil {
			c.gen(node.Inc)
		}
		c.Fprintf("	jmp .Lbegin%03d\n", labelNo)
		c.Fprintf(".Lend%03d:\n", labelNo)
		return

	case ND_FUNCCALL:
		numArgs := 0
		for arg := node.Args; arg != nil; arg = arg.Next {
			c.gen(arg)
			numArgs++
		}

		for i := numArgs - 1; i >= 0; i-- {
			c.Fprintf("	pop %s\n", argReg8[i])
		}

		labelNo++
		c.Fprintf("	mov rax, rsp\n")            // move rsp to rax
		c.Fprintf("	and rax, 15\n")             // calculate rax & 15, when rax == 16, rax is 0b10000, and 15(0b1110) & 0b10000, ZF become 0.
		c.Fprintf("	jnz .Lcall%03d\n", labelNo) // if ZF is 0, jamp to Lcall???.
		c.Fprintf("	mov rax, 0\n")              // remove rax
		c.Fprintf("	call %s\n", node.FuncName)
		c.Fprintf("	jmp .Lend%03d\n", labelNo)
		c.Fprintf(".Lcall%03d:\n", labelNo)
		c.Fprintf("	sub rsp, 8\n") // rspは8の倍数なので16の倍数にするために8を引く
		c.Fprintf("	mov rax, 0\n")
		c.Fprintf("	call %s\n", node.FuncName)
		c.Fprintf("	add rsp, 8\n")
		c.Fprintf(".Lend%03d:\n", labelNo)
		c.Fprintf("	push rax\n")
		return

	case ND_BLOCK:
		for n := node.Body; n != nil; n = n.Next {
			c.gen(n)
		}
		return

	case ND_RETURN:
		c.gen(node.Lhs)
		c.Fprintf("	pop rax\n")
		c.Fprintf("	jmp .Lreturn.%s\n", funcName)
		return
	}

	c.gen(node.Lhs)
	c.gen(node.Rhs)

	c.Fprintf("	pop rdi\n")
	c.Fprintf("	pop rax\n")

	switch node.Kind {
	case ND_ADD:
		if node.Ty.PtrTo != nil {
			c.Fprintf("	imul rdi, %d\n", sizeOf(node.Ty.PtrTo))
		}
		c.Fprintf("	add rax, rdi\n")
	case ND_SUB:
		if node.Ty.PtrTo != nil {
			c.Fprintf("	imul rdi, %d\n", sizeOf(node.Ty.PtrTo))
		}
		c.Fprintf("	sub rax, rdi\n")
	case ND_MUL:
		c.Fprintf("	imul rax, rdi\n")
	case ND_DIV:
		c.Fprintf("	cqo\n")
		c.Fprintf("	idiv rdi\n")
	case ND_EQ:
		c.Fprintf("	cmp rax, rdi\n")
		c.Fprintf("	sete al\n")
		c.Fprintf("	movzb rax, al\n")
	case ND_NE:
		c.Fprintf("	cmp rax, rdi\n")
		c.Fprintf("	setne al\n")
		c.Fprintf("	movzb rax, al\n")
	case ND_LT:
		c.Fprintf("	cmp rax, rdi\n")
		c.Fprintf("	setl al\n")
		c.Fprintf("	movzb rax, al\n")
	case ND_LE:
		c.Fprintf("	cmp rax, rdi\n")
		c.Fprintf("	setle al\n")
		c.Fprintf("	movzb rax, al\n")
	}

	c.Fprintf("	push rax\n")
}

func (c *codeWriter) emitData(prog *Program) {
	c.Fprintf(".data\n")

	for vl := prog.Globals; vl != nil; vl = vl.Next {
		c.Fprintf("%s:\n", vl.Var.Name)

		if vl.Var.Contents == "" {
			c.Fprintf("	.zero %d\n", sizeOf(vl.Var.Ty))
			continue
		}

		for i := 0; i < vl.Var.ContLen; i++ {
			c.Fprintf("	.byte %d\n", vl.Var.Contents[i])
		}
	}
}

func (c *codeWriter) loadArg(lvar *Var, idx int) {
	sz := sizeOf(lvar.Ty)
	if sz == 1 {
		c.Fprintf("	mov [rbp-%d], %s\n", lvar.Offset, argReg1[idx])
	} else {
		if sz != 8 {
			c.err = errors.New("invalid size")
		}
		c.Fprintf("	mov [rbp-%d], %s\n", lvar.Offset, argReg8[idx])
	}
}

func (c *codeWriter) emitText(prog *Program) {
	c.Fprintf(".text\n")

	for fn := prog.Fns; fn != nil; fn = fn.Next {
		c.Fprintf(".global %s\n", fn.Name)
		c.Fprintf("%s:\n", fn.Name)
		funcName = fn.Name

		// prologue
		// secure an area for the stack size of 'fn'
		c.Fprintf("	push rbp\n")
		c.Fprintf("	mov rbp, rsp\n")
		c.Fprintf("	sub rsp, %d\n", fn.StackSz)

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
		c.Fprintf(".Lreturn.%s:\n", funcName)
		c.Fprintf("	mov rsp, rbp\n")
		c.Fprintf("	pop rbp\n")
		c.Fprintf("	ret\n")
	}
}

func codeGen(w io.Writer, prog *Program) error {
	c := &codeWriter{w: w}
	// output the former 3 lines of the assembly
	c.Fprintf(".intel_syntax noprefix\n")
	c.emitData(prog)
	c.emitText(prog)

	return c.err
}
