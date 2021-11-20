//
// code generator
//
package main

import (
	"fmt"
	"io"
)

type codeWriter struct {
	w   io.Writer
	err error
}

func (c *codeWriter) printf(frmt string, a ...interface{}) {
	if c.err != nil {
		return
	}
	_, c.err = fmt.Fprintf(c.w, frmt, a...)
}

var argreg1 = []string{"dil", "sil", "dl", "cl", "r8b", "r9b"}
var argreg8 = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

var labelseq int
var funcname string

// Pushes the given node's address to the stack
func (c *codeWriter) genAddr(node *Node) {
	switch node.Kind {
	case ND_VAR:
		if node.Var.IsLocal {
			c.printf("	lea rax, [rbp-%d]\n", node.Var.Offset)
			c.printf("	push rax\n")
			return
		}
		c.printf("	push offset %s\n", node.Var.Name)
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
	}

	c.err = fmt.Errorf(errorTok(node.Tok, "not an lvalue"))
}

func (c *codeWriter) genLval(node *Node) {
	if node.Ty.Kind == TY_ARRAY {
		c.err = fmt.Errorf(errorTok(node.Tok, "not an lvalue"))
	}
	c.genAddr(node)
}

func (c *codeWriter) load(ty *Type) {
	c.printf("	pop rax\n")
	if sizeOf(ty) == 1 {
		c.printf("	movsx rax, byte ptr [rax]\n")
	} else {
		c.printf("	mov rax, [rax]\n")
	}
	c.printf("	push rax\n")
}

func (c *codeWriter) store(ty *Type) {
	c.printf("	pop rdi\n")
	c.printf("	pop rax\n")
	if sizeOf(ty) == 1 {
		c.printf("	mov [rax], dil\n")
	} else {
		c.printf("	mov [rax], rdi\n")
	}
	c.printf("	push rdi\n")
}

func (c *codeWriter) gen(node *Node) (err error) {
	if c.err != nil {
		return
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
		seq := labelseq
		labelseq++
		if node.Els != nil {
			c.gen(node.Cond)
			c.printf("	pop rax\n")
			c.printf("	cmp rax, 0\n")
			c.printf("	je .Lelse%d\n", seq)
			c.gen(node.Then)
			c.printf("	jmp .Lend%d\n", seq)
			c.printf(".Lelse%d:\n", seq)
			c.gen(node.Els)
			c.printf(".Lend%d:\n", seq)
			return
		}
		c.gen(node.Cond)
		c.printf("	pop rax\n")
		c.printf("	cmp rax, 0\n")
		c.printf("	je .Lend%d\n", seq)
		c.gen(node.Then)
		c.printf(".Lend%d:\n", seq)
		return
	case ND_WHILE:
		seq := labelseq
		labelseq++
		c.printf(".Lbegin%d:\n", seq)
		c.gen(node.Cond)
		c.printf("	pop rax\n")
		c.printf("	cmp rax, 0\n")
		c.printf("	je .Lend%d\n", seq)
		c.gen(node.Then)
		c.printf("	jmp .Lbegin%d\n", seq)
		c.printf(".Lend%d:\n", seq)
		return
	case ND_FOR:
		seq := labelseq
		labelseq++
		if node.Init != nil {
			c.gen(node.Init)
		}
		c.printf(".Lbegin%d:\n", seq)
		if node.Cond != nil {
			c.gen(node.Cond)
			c.printf("	pop rax\n")
			c.printf("	cmp rax, 0\n")
			c.printf("	je .Lend%d\n", seq)
		}
		c.gen(node.Then)
		if node.Inc != nil {
			c.gen(node.Inc)
		}
		c.printf("	jmp .Lbegin%d\n", seq)
		c.printf(".Lend%d:\n", seq)
		return
	case ND_BLOCK:
		for n := node.Body; n != nil; n = n.Next {
			c.gen(n)
		}
		return
	case ND_FUNCALL:
		nargs := 0
		for arg := node.Args; arg != nil; arg = arg.Next {
			c.gen(arg)
			nargs++
		}

		for i := nargs - 1; i >= 0; i-- {
			c.printf("	pop %s\n", argreg8[i])
		}

		// We need to align RSP to a 16 byte boundary before
		// calling a function because it is an ABI requirement.
		// RAX is set to 0 for variadic function.
		seq := labelseq
		labelseq++
		c.printf("	mov rax, rsp\n")
		c.printf("	and rax, 15\n")
		c.printf("	jnz .Lcall%d\n", seq)
		c.printf("	mov rax, 0\n")
		c.printf("	call %s\n", node.FuncName)
		c.printf("	jmp .Lend%d\n", seq)
		c.printf(".Lcall%d:\n", seq)
		c.printf("	sub rsp, 8\n")
		c.printf("	mov rax, 0\n")
		c.printf("	call %s\n", node.FuncName)
		c.printf("	add rsp, 8\n")
		c.printf(".Lend%d:\n", seq)
		c.printf("	push rax\n")
		return
	case ND_RETURN:
		c.gen(node.Lhs)
		c.printf("	pop rax\n")
		c.printf("	jmp .Lreturn.%s\n", funcname)
		return
	}

	c.gen(node.Lhs)
	c.gen(node.Rhs)

	c.printf("	pop rdi\n")
	c.printf("	pop rax\n")

	switch node.Kind {
	case ND_ADD:
		if node.Ty.Base != nil {
			c.printf("	imul rdi, %d\n", sizeOf(node.Ty.Base))
		}
		c.printf("	add rax, rdi\n")
	case ND_SUB:
		if node.Ty.Base != nil {
			c.printf("	imul rdi, %d\n", sizeOf(node.Ty.Base))
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
	return
}

func (c *codeWriter) loadArg(v *Var, idx int) {
	sz := sizeOf(v.Ty)
	if sz == 1 {
		c.printf("	mov [rbp-%d], %s\n", v.Offset, argreg1[idx])
	} else {
		if sz != 8 {
			c.err = fmt.Errorf("invalid size")
		}
		c.printf("	mov [rbp-%d], %s\n", v.Offset, argreg8[idx])
	}
}

func (c *codeWriter) emitData(prog *Program) {
	if c.err != nil {
		return
	}

	c.printf(".data\n")

	for vl := prog.Globs; vl != nil; vl = vl.Next {
		c.printf("%s:\n", vl.Var.Name)
		if vl.Var.Conts == nil {
			c.printf("	.zero %d\n", sizeOf(vl.Var.Ty))
			continue
		}

		for i := 0; i < vl.Var.ContLen; i++ {
			c.printf("	.byte %d\n", vl.Var.Conts[i])
		}
	}
}

func (c *codeWriter) emitText(prog *Program) {
	if c.err != nil {
		return
	}

	c.printf(".text\n")

	for fn := prog.Fns; fn != nil; fn = fn.Next {
		c.printf(".globl %s\n", fn.Name)
		c.printf("%s:\n", fn.Name)
		funcname = fn.Name

		// Prologue
		c.printf("	push rbp\n")
		c.printf("	mov rbp, rsp\n")
		c.printf("	sub rsp, %d\n", fn.StackSz)

		// Push arguments to the stack
		i := 0
		for vl := fn.Params; vl != nil; vl = vl.Next {
			c.loadArg(vl.Var, i)
			i++
		}

		// Emit code
		for n := fn.Node; n != nil; n = n.Next {
			c.gen(n)
		}

		// Epilogue
		c.printf(".Lreturn.%s:\n", funcname)
		c.printf("	mov rsp, rbp\n")
		c.printf("	pop rbp\n")
		c.printf("	ret\n")
	}
}

func codegen(w io.Writer, prog *Program) error {
	c := &codeWriter{w: w}

	c.printf(".intel_syntax noprefix\n")
	c.emitData(prog)
	c.emitText(prog)

	return c.err
}
