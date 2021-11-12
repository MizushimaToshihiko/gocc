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

func (c *codeWriter) gen(node *Node) (err error) {
	if c.err != nil {
		return
	}

	switch node.Kind {
	case ND_NUM:
		c.printf("	push %d\n", node.Val)
		return
	case ND_RETURN:
		c.gen(node.Lhs)
		c.printf("	pop rax\n")
		c.printf("	ret\n")
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

func codegen(node *Node, w io.Writer) error {
	c := &codeWriter{w: w}
	// output the former 3 lines of the assembly
	c.printf(".intel_syntax noprefix\n.globl main\nmain:\n")

	for n := node; n != nil; n = n.Next {
		c.gen(n)
		c.printf("	pop rax\n")
	}

	c.printf("	ret\n")
	return c.err
}
