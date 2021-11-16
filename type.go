package main

import "fmt"

type TypeKind int

// errWriter is struct for error handling
// it's based on:
// https://jxck.hatenablog.com/entry/golang-error-handling-lesson-by-rob-pike
type errWriter struct {
	err error
}

const (
	TY_INT TypeKind = iota
	TY_PTR
)

type Type struct {
	Kind TypeKind
	Base *Type
}

func intType() *Type {
	return &Type{Kind: TY_INT}
}

func pointerTo(base *Type) *Type {
	return &Type{Kind: TY_PTR, Base: base}
}

func (e *errWriter) visit(node *Node) {
	if e.err != nil {
		return
	}

	if node == nil {
		return
	}

	e.visit(node.Lhs)
	e.visit(node.Rhs)
	e.visit(node.Cond)
	e.visit(node.Then)
	e.visit(node.Els)
	e.visit(node.Init)
	e.visit(node.Inc)

	for n := node.Body; n != nil; n = n.Next {
		e.visit(n)
	}
	for n := node.Args; n != nil; n = n.Next {
		e.visit(n)
	}

	switch node.Kind {
	case ND_MUL, ND_DIV, ND_EQ, ND_NE, ND_LT, ND_LE, ND_VAR, ND_FUNCALL, ND_NUM:
		node.Ty = intType()
		return
	case ND_ADD:
		if node.Rhs.Ty.Kind == TY_PTR {
			tmp := node.Lhs
			node.Lhs = node.Rhs
			node.Rhs = tmp
		}
		if node.Rhs.Ty.Kind == TY_PTR {
			e.err = fmt.Errorf(errorTok(node.Tok, "invalid pointer arithmetic operands"))
		}
		node.Ty = node.Lhs.Ty
		return
	case ND_SUB:
		if node.Rhs.Ty.Kind == TY_PTR {
			tmp := node.Lhs
			node.Lhs = node.Rhs
			node.Rhs = tmp
		}
		if node.Rhs.Ty.Kind == TY_PTR {
			e.err = fmt.Errorf(errorTok(node.Tok, "invalid pointer arithmetic operands"))
		}
		node.Ty = node.Lhs.Ty
		return
	case ND_ASSIGN:
		node.Ty = node.Lhs.Ty
		return
	case ND_ADDR:
		node.Ty = pointerTo(node.Lhs.Ty)
		return
	case ND_DEREF:
		node.Ty = intType()
		return
	}
}

func addType(prog *Function) {
	e := &errWriter{}

	for fn := prog; fn != nil; fn = fn.Next {
		for node := fn.Node; node != nil; node = node.Next {
			e.visit(node)
		}
	}
}
