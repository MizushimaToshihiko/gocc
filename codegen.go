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
type errWriter struct {
	err error
}

// Almost all error is from 'fmt.Fprintf'
func (e *errWriter) Fprintf(w io.Writer, format string, a ...interface{}) {
	if e.err != nil {
		return // do nothing
	}
	_, e.err = fmt.Fprintf(w, format, a...)
}

var labelNo int
var argReg1 = []string{"dil", "sil", "dl", "cl", "r8b", "r9b"}
var argReg8 = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}
var funcName string

func (e *errWriter) genAddr(w io.Writer, node *Node) {
	if e.err != nil {
		return // do nothing
	}

	switch node.Kind {
	case ND_VAR:
		if node.Var.IsLocal {
			e.Fprintf(w, "	lea rax, [rbp-%d]\n", node.Var.Offset)
			e.Fprintf(w, "	push rax\n")
		} else {
			e.Fprintf(w, "	push offset %s\n", node.Var.Name)
		}
		return
	case ND_DEREF:
		e.gen(w, node.Lhs)
		return
	default:
		e.err = fmt.Errorf(
			"e.genLval(): err:\n%s",
			errorTok(node.Tok, "the left value is not a variable"),
		)
	}
}

func (e *errWriter) genLval(w io.Writer, node *Node) {
	if node.Ty.Kind == TY_ARRAY {
		e.err = fmt.Errorf(
			"e.genLval(): err:\n%s",
			errorTok(node.Tok, "the left value is not a variable"),
		)
	}
	e.genAddr(w, node)
}

func (e *errWriter) load(w io.Writer, ty *Type) {
	if e.err != nil {
		return // do nothing
	}

	e.Fprintf(w, "	pop rax\n")
	if sizeOf(ty) == 1 {
		e.Fprintf(w, "	movsx rax, byte ptr [rax]\n")
	} else {
		e.Fprintf(w, "	mov rax, [rax]\n")
	}
	e.Fprintf(w, "	push rax\n")
}

func (e *errWriter) store(w io.Writer, ty *Type) {
	if e.err != nil {
		return // do nothing
	}

	e.Fprintf(w, "	pop rdi\n")
	e.Fprintf(w, "	pop rax\n")

	if sizeOf(ty) == 1 {
		e.Fprintf(w, "	mov [rax], dil\n")
	} else {
		e.Fprintf(w, "	mov [rax], rdi\n")
	}

	e.Fprintf(w, "	push rdi\n")
}

func (e *errWriter) gen(w io.Writer, node *Node) {
	if e.err != nil {
		return // do nothing
	}

	switch node.Kind {
	case ND_NULL:
		return
	case ND_NUM:
		e.Fprintf(w, "	push %d\n", node.Val)
		return
	case ND_EXPR_STMT:
		e.gen(w, node.Lhs)
		e.Fprintf(w, "	add rsp, 8\n")
		return
	case ND_VAR:
		e.genAddr(w, node)
		if node.Ty.Kind != TY_ARRAY {
			e.load(w, node.Ty)
		}
		return

	case ND_ASSIGN:
		e.genLval(w, node.Lhs)
		e.gen(w, node.Rhs)
		// store
		e.store(w, node.Ty)
		return

	case ND_ADDR:
		e.genAddr(w, node.Lhs)
		return

	case ND_DEREF:
		e.gen(w, node.Lhs)
		if node.Ty.Kind != TY_ARRAY {
			e.load(w, node.Ty)
		}
		return

	case ND_IF:
		e.gen(w, node.Cond)
		e.Fprintf(w, "	pop rax\n")
		e.Fprintf(w, "	cmp rax, 0\n")

		labelNo++
		if node.Els != nil {
			e.Fprintf(w, "	je .Lelse%03d\n", labelNo)
		} else {
			e.Fprintf(w, "	je .Lend%03d\n", labelNo)
		}

		e.gen(w, node.Then)

		if node.Els != nil {
			e.Fprintf(w, " jmp .Lend%03d\n", labelNo)
			e.Fprintf(w, ".Lelse%03d:\n", labelNo)
			e.gen(w, node.Els)
		}

		e.Fprintf(w, ".Lend%03d:\n", labelNo)
		return

	case ND_WHILE:
		labelNo++
		e.Fprintf(w, ".Lbegin%03d:\n", labelNo)
		e.gen(w, node.Cond)
		e.Fprintf(w, "	pop rax\n")
		e.Fprintf(w, "	cmp rax, 0\n")
		e.Fprintf(w, "	je .Lend%03d\n", labelNo)

		e.gen(w, node.Then)

		e.Fprintf(w, "	jmp .Lbegin%03d\n", labelNo)
		e.Fprintf(w, ".Lend%03d:\n", labelNo)
		return

	case ND_FOR:
		if node.Init != nil {
			e.gen(w, node.Init)
		}

		labelNo++
		e.Fprintf(w, ".Lbegin%03d:\n", labelNo)

		if node.Cond != nil {
			e.gen(w, node.Cond)
			e.Fprintf(w, "	pop rax\n")
			e.Fprintf(w, "	cmp rax, 0\n")
			e.Fprintf(w, "	je .Lend%03d\n", labelNo)
		}

		e.gen(w, node.Then)

		if node.Inc != nil {
			e.gen(w, node.Inc)
		}
		e.Fprintf(w, "	jmp .Lbegin%03d\n", labelNo)
		e.Fprintf(w, ".Lend%03d:\n", labelNo)
		return

	case ND_FUNCCALL:
		numArgs := 0
		for arg := node.Args; arg != nil; arg = arg.Next {
			e.gen(w, arg)
			numArgs++
		}

		for i := numArgs - 1; i >= 0; i-- {
			e.Fprintf(w, "	pop %s\n", argReg8[i])
		}

		labelNo++
		e.Fprintf(w, "	mov rax, rsp\n")            // move rsp to rax
		e.Fprintf(w, "	and rax, 15\n")             // calculate rax & 15, when rax == 16, rax is 0b10000, and 15(0b1110) & 0b10000, ZF become 0.
		e.Fprintf(w, "	jnz .Lcall%03d\n", labelNo) // if ZF is 0, jamp to Lcall???.
		e.Fprintf(w, "	mov rax, 0\n")              // remove rax
		e.Fprintf(w, "	call %s\n", node.FuncName)
		e.Fprintf(w, "	jmp .Lend%03d\n", labelNo)
		e.Fprintf(w, ".Lcall%03d:\n", labelNo)
		e.Fprintf(w, "	sub rsp, 8\n") // rspは8の倍数なので16の倍数にするために8を引く
		e.Fprintf(w, "	mov rax, 0\n")
		e.Fprintf(w, "	call %s\n", node.FuncName)
		e.Fprintf(w, "	add rsp, 8\n")
		e.Fprintf(w, ".Lend%03d:\n", labelNo)
		e.Fprintf(w, "	push rax\n")
		return

	case ND_BLOCK:
		for n := node.Body; n != nil; n = n.Next {
			e.gen(w, n)
		}
		return

	case ND_RETURN:
		e.gen(w, node.Lhs)
		e.Fprintf(w, "	pop rax\n")
		e.Fprintf(w, "	jmp .Lreturn.%s\n", funcName)
		return
	}

	e.gen(w, node.Lhs)
	e.gen(w, node.Rhs)

	e.Fprintf(w, "	pop rdi\n")
	e.Fprintf(w, "	pop rax\n")

	switch node.Kind {
	case ND_ADD:
		if node.Ty.PtrTo != nil {
			e.Fprintf(w, "	imul rdi, %d\n", sizeOf(node.Ty.PtrTo))
		}
		e.Fprintf(w, "	add rax, rdi\n")
	case ND_SUB:
		if node.Ty.PtrTo != nil {
			e.Fprintf(w, "	imul rdi, %d\n", sizeOf(node.Ty.PtrTo))
		}
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
}

func (e *errWriter) emitData(w io.Writer, prog *Program) {
	e.Fprintf(w, ".data\n")

	for vl := prog.Globals; vl != nil; vl = vl.Next {
		e.Fprintf(w, "%s:\n", vl.Var.Name)
		e.Fprintf(w, "	.zero %d\n", sizeOf(vl.Var.Ty))
	}
}

func (e *errWriter) loadArg(w io.Writer, lvar *Var, idx int) {
	sz := sizeOf(lvar.Ty)
	if sz == 1 {
		e.Fprintf(w, "	mov [rbp-%d], %s\n", lvar.Offset, argReg1[idx])
	} else {
		if sz != 8 {
			e.err = errors.New("invalid size")
		}
		e.Fprintf(w, "	mov [rbp-%d], %s\n", lvar.Offset, argReg8[idx])
	}
}

func (e *errWriter) emitText(w io.Writer, prog *Program) {
	e.Fprintf(w, ".text\n")

	for fn := prog.Fns; fn != nil; fn = fn.Next {
		e.Fprintf(w, ".global %s\n", fn.Name)
		e.Fprintf(w, "%s:\n", fn.Name)
		funcName = fn.Name

		// prologue
		// secure an area for the stack size of 'fn'
		e.Fprintf(w, "	push rbp\n")
		e.Fprintf(w, "	mov rbp, rsp\n")
		e.Fprintf(w, "	sub rsp, %d\n", fn.StackSz)

		// push arguments to the stack
		i := 0
		for vl := fn.Params; vl != nil; vl = vl.Next {
			e.loadArg(w, vl.Var, i)
			i++
		}

		// emit code
		for node := fn.Node; node != nil; node = node.Next {
			e.gen(w, node)
		}
		// epilogue
		// the result of the expression is in 'rax',
		// and it is the return value
		e.Fprintf(w, ".Lreturn.%s:\n", funcName)
		e.Fprintf(w, "	mov rsp, rbp\n")
		e.Fprintf(w, "	pop rbp\n")
		e.Fprintf(w, "	ret\n")
	}
}

func codeGen(w io.Writer, prog *Program) error {
	e := &errWriter{}
	// output the former 3 lines of the assembly
	e.Fprintf(w, ".intel_syntax noprefix\n")
	e.emitData(w, prog)
	e.emitText(w, prog)

	return e.err
}
