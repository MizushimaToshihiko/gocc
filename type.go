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
	TY_ARRAY
)

type Type struct {
	Kind  TypeKind
	Base  *Type
	ArrSz int // Array size
}

func intType() *Type {
	return &Type{Kind: TY_INT}
}

func pointerTo(base *Type) *Type {
	return &Type{Kind: TY_PTR, Base: base}
}

func arrayOf(base *Type, size int) *Type {
	return &Type{Kind: TY_ARRAY, Base: base, ArrSz: size}
}

func sizeOf(ty *Type) int {
	if ty.Kind == TY_INT || ty.Kind == TY_PTR {
		return 8
	}
	if ty.Kind != TY_ARRAY {
		panic("invalid type")
	}
	return sizeOf(ty.Base) * ty.ArrSz
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
	case ND_MUL, ND_DIV, ND_EQ, ND_NE, ND_LT, ND_LE, ND_FUNCALL, ND_NUM:
		node.Ty = intType()
		return
	case ND_VAR:
		node.Ty = node.Var.Ty
		return
	case ND_ADD:
		if node.Rhs.Ty.Base != nil {
			tmp := node.Lhs
			node.Lhs = node.Rhs
			node.Rhs = tmp
		}
		if node.Rhs.Ty.Base != nil {
			e.err = fmt.Errorf(errorTok(node.Tok, "invalid pointer arithmetic operands"))
		}
		node.Ty = node.Lhs.Ty
		return
	case ND_SUB:
		if node.Rhs.Ty.Base != nil {
			e.err = fmt.Errorf(errorTok(node.Tok, "invalid pointer arithmetic operands"))
		}
		node.Ty = node.Lhs.Ty
		return
	case ND_ASSIGN:
		node.Ty = node.Lhs.Ty
		return
	case ND_ADDR:
		if node.Lhs.Ty.Kind == TY_ARRAY {
			node.Ty = pointerTo(node.Lhs.Ty.Base)
			return
		}
		node.Ty = pointerTo(node.Lhs.Ty)
		return
	case ND_DEREF:
		// fmt.Printf("node: %#v\n'%s'\n\n", node, node.Tok.Str)
		// fmt.Printf("node.Lhs: %#v\n'%s'\n\n", node.Lhs, node.Lhs.Tok.Str)
		// fmt.Printf("node.Lhs.Var: %#v\n\n", node.Lhs.Var)
		if node.Lhs.Ty.Base == nil {
			e.err = fmt.Errorf(errorTok(node.Tok, "invalid pointer dereference"))
			return
		}
		node.Ty = node.Lhs.Ty.Base
		return
	}
}

func addType(prog *Function) error {
	e := &errWriter{}

	for fn := prog; fn != nil; fn = fn.Next {
		for node := fn.Node; node != nil; node = node.Next {
			e.visit(node)
		}
	}
	return e.err
}
