//
// code generator
//
package main

import (
	"fmt"
	"io"
)

func gen(w io.Writer, node *Node) (err error) {
	if node.Kind == ND_NUM {
		if _, err = fmt.Fprintf(w, "	push %d\n", node.Val); err != nil {
			return
		}
	}

	if err = gen(w, node.Lhs); err != nil {
		return
	}
	if err = gen(w, node.Rhs); err != nil {
		return
	}

	if _, err = fmt.Fprintln(w, "	pop rdi"); err != nil {
		return
	}
	if _, err = fmt.Fprintln(w, "	pop rax"); err != nil {
		return
	}

	switch node.Kind {
	case ND_ADD:
		if _, err = fmt.Fprintln(w, "	add rax, rdi"); err != nil {
			return
		}
	case ND_SUB:
		if _, err = fmt.Fprintln(w, "	sub rax, rdi"); err != nil {
			return
		}
	case ND_MUL:
		if _, err = fmt.Fprintln(w, "	imul rax, rdi"); err != nil {
			return
		}
	case ND_DIV:
		if _, err = fmt.Fprintln(w, "	cqo"); err != nil {
			return
		}
		if _, err = fmt.Fprintln(w, "	idiv rdi"); err != nil {
			return
		}
	case ND_EQ:
		if _, err = fmt.Fprintln(w, "	cmp rax, rdi"); err != nil {
			return
		}
		if _, err = fmt.Fprintln(w, "	sete al"); err != nil {
			return
		}
		if _, err = fmt.Fprintln(w, "	movzb rax, al"); err != nil {
			return
		}
	case ND_NE:
		if _, err = fmt.Fprintln(w, "	cmp rax, rdi"); err != nil {
			return
		}
		if _, err = fmt.Fprintln(w, "	setne al"); err != nil {
			return
		}
		if _, err = fmt.Fprintln(w, "	movzb rax, al"); err != nil {
			return
		}
	case ND_LT:
		if _, err = fmt.Fprintln(w, "	cmp rax, rdi"); err != nil {
			return
		}
		if _, err = fmt.Fprintln(w, "	setl al"); err != nil {
			return
		}
		if _, err = fmt.Fprintln(w, "	movzb rax, al"); err != nil {
			return
		}
	case ND_LE:
		if _, err = fmt.Fprintln(w, "	cmp rax, rdi"); err != nil {
			return
		}
		if _, err = fmt.Fprintln(w, "	setle al"); err != nil {
			return
		}
		if _, err = fmt.Fprintln(w, "	movzb rax, al"); err != nil {
			return
		}
	}

	if _, err = fmt.Fprintln(w, "	push rax"); err != nil {
		return
	}

	return
}

func codeGen(w io.Writer, node *Node) (err error) {
	// output the former 3 lines of the assembly
	if _, err = fmt.Fprintln(w, ".intel_syntax noprefix\n.globl main\nmain:"); err != nil {
		return
	}

	// make the asm code, down on the AST
	if err = gen(w, node); err != nil {
		return
	}

	// the value of the expression should remain on the top of 'stack',
	// so load this value into rax.
	if _, err = fmt.Fprintln(w, "	pop rax"); err != nil {
		return
	}
	if _, err = fmt.Fprintln(w, "	ret"); err != nil {
		return
	}

	return
}
