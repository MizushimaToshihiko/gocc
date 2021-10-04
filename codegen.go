//
// code generator
//
package main

import (
	"fmt"
	"io"
)

func gen(node *Node, w io.Writer) (err error) {
	if node.Kind == ND_NUM {
		_, err = fmt.Fprintf(w, "	push %d\n", node.Val)
		return
	}

	err = gen(node.Lhs, w)
	if err != nil {
		return
	}
	err = gen(node.Rhs, w)
	if err != nil {
		return
	}

	fmt.Fprintln(w, "	pop rdi")
	fmt.Fprintln(w, "	pop rax")

	switch node.Kind {
	case ND_ADD:
		fmt.Fprintln(w, "	add rax, rdi")
	case ND_SUB:
		fmt.Fprintln(w, "	sub rax, rdi")
	case ND_MUL:
		fmt.Fprintln(w, "	imul rax, rdi")
	case ND_DIV:
		fmt.Fprintln(w, "	cqo")
		fmt.Fprintln(w, "	idiv rdi")
	case ND_EQ:
		fmt.Fprintln(w, "	cmp rax, rdi")
		fmt.Fprintln(w, "	sete al")
		fmt.Fprintln(w, "	movzb rax, al")
	case ND_NE:
		fmt.Fprintln(w, "	cmp rax, rdi")
		fmt.Fprintln(w, "	setne al")
		fmt.Fprintln(w, "	movzb rax, al")
	case ND_LT:
		fmt.Fprintln(w, "	cmp rax, rdi")
		fmt.Fprintln(w, "	setl al")
		fmt.Fprintln(w, "	movzb rax, al")
	case ND_LE:
		fmt.Fprintln(w, "	cmp rax, rdi")
		fmt.Fprintln(w, "	setle al")
		fmt.Fprintln(w, "	movzb rax, al")
	}

	fmt.Fprintln(w, "	push rax")
	return
}
