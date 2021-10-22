// add 'Type' to AST
package main

import (
	"errors"
)

type TypeKind int

const (
	TY_INT TypeKind = iota // int
	TY_PTR                 // pointer
)

type Type struct {
	Kind TypeKind
	Base *Type
}

func sizeOf(ty *Type) int {
	if ty.Kind == TY_INT {
		return 4
	}
	if ty.Kind == TY_PTR {
		return 8
	}
	return sizeOf(ty.Base)
}

func intType() *Type {
	return &Type{Kind: TY_INT}
}

func pointerTo(base *Type) *Type {
	return &Type{Kind: TY_PTR, Base: base}
}

func (e *errWriter) visit(node *Node) {
	// printCalledFunc()

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
	case ND_MUL, ND_DIV, ND_EQ, ND_NE, ND_LT, ND_LE, ND_FUNCCALL, ND_NUM:
		node.Ty = intType()
		return
	case ND_LVAR:
		node.Ty = node.Var.Ty
		return
	case ND_ADD:
		if node.Rhs.Ty.Kind == TY_PTR {
			tmp := node.Lhs
			node.Lhs = node.Rhs
			node.Rhs = tmp
		}
		if node.Rhs.Ty.Kind == TY_PTR {
			e.err = errors.New("invalid pointer arithmetic operands")
		}
		node.Ty = node.Lhs.Ty
		return
	case ND_SUB:
		if node.Rhs.Ty.Kind == TY_PTR {
			e.err = errors.New("invalid pointer arithmetic operands")
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
		if node.Lhs.Ty.Kind != TY_PTR {
			e.err = errors.New("invalid pointer dereference")
		}
		node.Ty = node.Lhs.Ty.Base
		return
	case ND_SIZEOF:
		node.Kind = ND_NUM
		node.Ty = intType()
		node.Val = sizeOf(node.Lhs.Ty)
		node.Lhs = nil
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
