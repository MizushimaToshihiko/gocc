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

var labelseq int

// Pushes the given node's address to the stack
func (c *codeWriter) genAddr(node *Node) {
	if node.Kind == ND_VAR {
		c.printf("	lea rax, [rbp-%d]\n", node.Var.Offset)
		c.printf("	push rax\n")
		return
	}

	c.err = fmt.Errorf("not an lvalue")
}

func (c *codeWriter) load() {
	c.printf("	pop rax\n")
	c.printf("	mov rax, [rax]\n")
	c.printf("	push rax\n")
}

func (c *codeWriter) store() {
	c.printf("	pop rdi\n")
	c.printf("	pop rax\n")
	c.printf("	mov [rax], rdi\n")
	c.printf("	push rdi\n")
}

func (c *codeWriter) gen(node *Node) (err error) {
	if c.err != nil {
		return
	}

	switch node.Kind {
	case ND_NUM:
		c.printf("	push %d\n", node.Val)
		return
	case ND_EXPR_STMT:
		c.gen(node.Lhs)
		c.printf("	add rsp, 8\n")
		return
	case ND_VAR:
		c.genAddr(node)
		c.load()
		return
	case ND_ASSIGN:
		c.genAddr(node.Lhs)
		c.gen(node.Rhs)
		c.store()
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
		c.printf("	call %s\n", node.FuncName)
		c.printf("	push rax\n")
		return
	case ND_RETURN:
		c.gen(node.Lhs)
		c.printf("	pop rax\n")
		c.printf("	jmp .Lreturn\n")
		return
	}

	c.gen(node.Lhs)
	c.gen(node.Rhs)

	c.printf("	pop rdi\n")
	c.printf("	pop rax\n")

	switch node.Kind {
	case ND_ADD:
		c.printf("	add rax, rdi\n")
	case ND_SUB:
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

func codegen(prog *Program, w io.Writer) error {
	c := &codeWriter{w: w}
	// output the former 3 lines of the assembly
	c.printf(".intel_syntax noprefix\n.globl main\nmain:\n")

	// Prologue
	c.printf("	push rbp\n")
	c.printf("	mov rbp, rsp\n")
	c.printf("	sub rsp, %d\n", prog.StackSz)

	for n := prog.Node; n != nil; n = n.Next {
		c.gen(n)
	}

	// Epilogue
	c.printf(".Lreturn:\n")
	c.printf("	mov rsp, rbp\n")
	c.printf("	pop rbp\n")
	c.printf("	ret\n")
	return c.err
}
