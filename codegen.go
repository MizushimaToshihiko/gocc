//
// code generator
//
package main

import (
	"errors"
	"fmt"
	"io"
)

var labelNo int

func genLval(w io.Writer, node *Node) (err error) {
	if node.Kind != ND_LVAR {
		err = errors.New("the left value is not a variable")
		return
	}

	if _, err = fmt.Fprintln(w, "	mov rax, rbp"); err != nil {
		return
	}
	if _, err = fmt.Fprintf(w, "	sub rax, %d\n", node.Offset); err != nil {
		return
	}
	if _, err = fmt.Fprintln(w, "	push rax"); err != nil {
		return
	}

	return nil
}

func gen(w io.Writer, node *Node) (err error) {

	switch node.Kind {
	case ND_NUM:
		_, err = fmt.Fprintf(w, "	push %d\n", node.Val)
		return
	case ND_LVAR:
		err = genLval(w, node)
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	pop rax")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	mov rax, [rax]")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	push rax")
		return
	case ND_ASSIGN:
		err = genLval(w, node.Lhs)
		if err != nil {
			return
		}
		err = gen(w, node.Rhs)
		if err != nil {
			return
		}

		_, err = fmt.Fprintln(w, "	pop rdi")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	pop rax")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	mov [rax], rdi")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	push rdi")
		return

	case ND_RETURN:
		err = gen(w, node.Lhs)
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	pop rax")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	mov rsp, rbp")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	pop rbp")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	ret")
		return

	case ND_IF:
		err = gen(w, node.Cond)
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	pop rax")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	cmp rax, 0")
		if err != nil {
			return
		}

		labelNo++
		if node.Els != nil {
			_, err = fmt.Fprintf(w, "	je .Lelse%03d\n", labelNo)
			if err != nil {
				return
			}
		} else {
			_, err = fmt.Fprintf(w, "	je .Lend%03d\n", labelNo)
			if err != nil {
				return
			}
		}

		err = gen(w, node.Then)
		if err != nil {
			return
		}

		if node.Els != nil {
			_, err = fmt.Fprintf(w, " jmp .Lend%03d\n", labelNo)
			if err != nil {
				return
			}
			_, err = fmt.Fprintf(w, ".Lelse%03d:\n", labelNo)
			if err != nil {
				return
			}
			err = gen(w, node.Els)
			if err != nil {
				return
			}
		}

		_, err = fmt.Fprintf(w, ".Lend%03d:\n", labelNo)
		return

	case ND_WHILE:
		labelNo++
		_, err = fmt.Fprintf(w, ".Lbegin%03d:\n", labelNo)
		if err != nil {
			return
		}
		err = gen(w, node.Cond)
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	pop rax")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	cmp rax, 0")
		if err != nil {
			return
		}
		_, err = fmt.Fprintf(w, "	je .Lend%03d\n", labelNo)
		if err != nil {
			return
		}
		err = gen(w, node.Then)
		if err != nil {
			return
		}
		_, err = fmt.Fprintf(w, "	jmp .Lbegin%03d\n", labelNo)
		if err != nil {
			return
		}
		_, err = fmt.Fprintf(w, ".Lend%03d:\n", labelNo)
		return
	case ND_FOR:
		if node.Init != nil {
			err = gen(w, node.Init)
			if err != nil {
				return
			}
		}
		labelNo++
		_, err = fmt.Fprintf(w, ".Lbegin%03d:\n", labelNo)
		if err != nil {
			return
		}
		if node.Cond != nil {
			err = gen(w, node.Cond)
			if err != nil {
				return
			}
		}
		_, err = fmt.Fprintln(w, "	pop rax")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	cmp rax, 0")
		if err != nil {
			return
		}
		_, err = fmt.Fprintf(w, "	je .Lend%03d\n", labelNo)
		if err != nil {
			return
		}

		err = gen(w, node.Then)
		if err != nil {
			return
		}

		if node.Inc != nil {
			err = gen(w, node.Inc)
			if err != nil {
				return
			}
		}
		_, err = fmt.Fprintf(w, "jmp .Lbegin%03d\n", labelNo)
		if err != nil {
			return
		}
		_, err = fmt.Fprintf(w, ".Lend%03d:\n", labelNo)
		return
	}

	err = gen(w, node.Lhs)
	if err != nil {
		return
	}
	err = gen(w, node.Rhs)
	if err != nil {
		return
	}

	_, err = fmt.Fprintln(w, "	pop rdi")
	if err != nil {
		return
	}
	_, err = fmt.Fprintln(w, "	pop rax")
	if err != nil {
		return
	}

	switch node.Kind {
	case ND_ADD:
		_, err = fmt.Fprintln(w, "	add rax, rdi")
		if err != nil {
			return
		}
	case ND_SUB:
		_, err = fmt.Fprintln(w, "	sub rax, rdi")
		if err != nil {
			return
		}
	case ND_MUL:
		_, err = fmt.Fprintln(w, "	imul rax, rdi")
		if err != nil {
			return
		}
	case ND_DIV:
		_, err = fmt.Fprintln(w, "	cqo")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	idiv rdi")
		if err != nil {
			return
		}
	case ND_EQ:
		_, err = fmt.Fprintln(w, "	cmp rax, rdi")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	sete al")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	movzb rax, al")
		if err != nil {
			return
		}

	case ND_NE:
		_, err = fmt.Fprintln(w, "	cmp rax, rdi")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	setne al")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	movzb rax, al")
		if err != nil {
			return
		}
	case ND_LT:
		_, err = fmt.Fprintln(w, "	cmp rax, rdi")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	setl al")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	movzb rax, al")
		if err != nil {
			return
		}
	case ND_LE:
		_, err = fmt.Fprintln(w, "	cmp rax, rdi")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	setle al")
		if err != nil {
			return
		}
		_, err = fmt.Fprintln(w, "	movzb rax, al")
		if err != nil {
			return
		}
	}

	_, err = fmt.Fprintln(w, "	push rax")
	return
}

func codeGen(w io.Writer) (err error) {
	// output the former 3 lines of the assembly
	_, err = fmt.Fprintln(w, ".intel_syntax noprefix\n.globl main\nmain:")
	if err != nil {
		return err
	}

	// prologue
	// secure an area for 26 variables
	_, err = fmt.Fprintln(w, "	push rbp")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, "	mov rbp, rsp")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, "	sub rsp, 208")
	if err != nil {
		return err
	}

	for _, c := range code {
		if c == nil {
			break
		}

		gen(w, c)

		// the one value shuld remain in stack,
		// so pop to keep the stack from overflowing.
		_, err = fmt.Fprintln(w, "	pop rax")
		if err != nil {
			return err
		}
	}

	// epilogue
	// the result of the expression is in 'rax',
	// and it is the return value
	_, err = fmt.Fprintln(w, "	mov rsp, rbp")
	if err != nil {
		return
	}
	_, err = fmt.Fprintln(w, "	pop rbp")
	if err != nil {
		return
	}
	_, err = fmt.Fprintln(w, "	ret")
	return
}
