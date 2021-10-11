//
// code generator
//
package main

import (
	"errors"
	"fmt"
	"io"
)

// struct errWriter is for error handling
// it's based on:
// https://jxck.hatenablog.com/entry/golang-error-handling-lesson-by-rob-pike
type errWriter struct {
	err error
}

func (e *errWriter) Fprintf(w io.Writer, format string, a ...interface{}) {
	if e.err != nil {
		return
	}
	_, e.err = fmt.Fprintf(w, format, a...)
}

var labelNo int
var argReg = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

func genLval(w io.Writer, node *Node) (err error) {
	e := &errWriter{}

	if node.Kind != ND_LVAR {
		err = errors.New("the left value is not a variable")
		return
	}

	e.Fprintf(w, "	mov rax, rbp\n")
	e.Fprintf(w, "	sub rax, %d\n", node.Offset)
	e.Fprintf(w, "	push rax\n")

	if e.err != nil {
		return e.err
	}
	return nil
}

func gen(w io.Writer, node *Node) (err error) {
	e := &errWriter{}

	switch node.Kind {
	case ND_NUM:
		e.Fprintf(w, "	push %d\n", node.Val)
		err = e.err
		return
	case ND_LVAR:
		err = genLval(w, node)
		if err != nil {
			return
		}
		e.Fprintf(w, "	pop rax\n")
		e.Fprintf(w, "	mov rax, [rax]\n")
		e.Fprintf(w, "	push rax\n")
		err = e.err
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

		e.Fprintf(w, "	pop rdi\n")
		e.Fprintf(w, "	pop rax\n")
		e.Fprintf(w, "	mov [rax], rdi\n")
		if err != nil {
			return
		}
		e.Fprintf(w, "	push rdi\n")
		err = e.err
		return

	case ND_IF:
		err = gen(w, node.Cond)
		if err != nil {
			return
		}
		e.Fprintf(w, "	pop rax\n")
		e.Fprintf(w, "	cmp rax, 0\n")

		labelNo++
		if node.Els != nil {
			e.Fprintf(w, "	je .Lelse%03d\n", labelNo)
		} else {
			e.Fprintf(w, "	je .Lend%03d\n", labelNo)
		}

		err = gen(w, node.Then)
		if err != nil {
			return
		}

		if node.Els != nil {
			e.Fprintf(w, " jmp .Lend%03d\n", labelNo)
			e.Fprintf(w, ".Lelse%03d:\n", labelNo)
			err = gen(w, node.Els)
			if err != nil {
				return
			}
		}

		e.Fprintf(w, ".Lend%03d:\n", labelNo)
		err = e.err
		return

	case ND_WHILE:
		labelNo++
		e.Fprintf(w, ".Lbegin%03d:\n", labelNo)
		err = gen(w, node.Cond)
		if err != nil {
			return
		}
		e.Fprintf(w, "	pop rax\n")
		e.Fprintf(w, "	cmp rax, 0\n")
		e.Fprintf(w, "	je .Lend%03d\n", labelNo)

		err = gen(w, node.Then)
		if err != nil {
			return
		}

		e.Fprintf(w, "	jmp .Lbegin%03d\n", labelNo)
		e.Fprintf(w, ".Lend%03d:\n", labelNo)
		err = e.err
		return

	case ND_FOR:
		if node.Init != nil {
			err = gen(w, node.Init)
			if err != nil {
				return
			}
		}

		labelNo++
		e.Fprintf(w, ".Lbegin%03d:\n", labelNo)

		if node.Cond != nil {
			err = gen(w, node.Cond)
			if err != nil {
				return
			}
			e.Fprintf(w, "	pop rax\n")
			e.Fprintf(w, "	cmp rax, 0\n")
			e.Fprintf(w, "	je .Lend%03d\n", labelNo)
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
		e.Fprintf(w, "	jmp .Lbegin%03d\n", labelNo)
		e.Fprintf(w, ".Lend%03d:\n", labelNo)
		err = e.err
		return

	case ND_BLOCK:
		for n := node.Body; n != nil; n = n.Next {
			err = gen(w, n)
			if err != nil {
				return
			}
		}
		return

	case ND_FUNCCALL:
		numArgs := 0
		for arg := node.Args; arg != nil; arg = arg.Next {
			err = gen(w, arg)
			if err != nil {
				return
			}
			numArgs++
		}

		for i := numArgs - 1; i >= 0; i-- {
			e.Fprintf(w, "pop %s\n", argReg[i])
		}

		e.Fprintf(w, "	call %s\n", node.FuncName)
		e.Fprintf(w, "	push rax\n")
		err = e.err
		return

	case ND_RETURN:
		err = gen(w, node.Lhs)
		if err != nil {
			return
		}
		e.Fprintf(w, "	pop rax\n")
		e.Fprintf(w, "	jmp .Lreturn\n")
		err = e.err
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

	e.Fprintf(w, "	pop rdi\n")
	e.Fprintf(w, "	pop rax\n")

	switch node.Kind {
	case ND_ADD:
		e.Fprintf(w, "	add rax, rdi\n")
	case ND_SUB:
		e.Fprintf(w, "	sub rax, rdi\n")
	case ND_MUL:
		e.Fprintf(w, "	imul rax, rdi\n")
	case ND_DIV:
		e.Fprintf(w, "	cqo\n")
		e.Fprintf(w, "	idiv rdi\n")
	case ND_EQ:
		e.Fprintf(w, "	cmp rax, rdi\n")
		e.Fprintf(w, "	sete al\n")
		e.Fprintf(w, "	movzb rax, al\n")
	case ND_NE:
		e.Fprintf(w, "	cmp rax, rdi\n")
		e.Fprintf(w, "	setne al\n")
		e.Fprintf(w, "	movzb rax, al\n")
	case ND_LT:
		e.Fprintf(w, "	cmp rax, rdi\n")
		e.Fprintf(w, "	setl al\n")
		e.Fprintf(w, "	movzb rax, al\n")
	case ND_LE:
		e.Fprintf(w, "	cmp rax, rdi\n")
		e.Fprintf(w, "	setle al\n")
		e.Fprintf(w, "	movzb rax, al\n")
	}

	e.Fprintf(w, "	push rax\n")
	err = e.err
	return
}

func codeGen(w io.Writer) (err error) {
	e := &errWriter{}
	// output the former 3 lines of the assembly
	e.Fprintf(w, ".intel_syntax noprefix\n.globl main\nmain:\n")

	// prologue
	// secure an area for 26 variables
	e.Fprintf(w, "	push rbp\n")
	e.Fprintf(w, "	mov rbp, rsp\n")
	e.Fprintf(w, "	sub rsp, 208\n")

	for _, c := range code {
		if c == nil {
			break
		}

		err = gen(w, c)
		if err != nil {
			return
		}

		// the one value shuld remain in stack,
		// so pop to keep the stack from overflowing.
		e.Fprintf(w, "	pop rax\n")
	}

	// epilogue
	// the result of the expression is in 'rax',
	// and it is the return value
	e.Fprintf(w, ".Lreturn:\n")
	e.Fprintf(w, "	mov rsp, rbp\n")
	e.Fprintf(w, "	pop rbp\n")
	e.Fprintf(w, "	ret\n")
	err = e.err
	return
}
