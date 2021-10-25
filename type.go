// add 'Type' to AST
package main

import (
	"fmt"
)

type TypeKind int

const (
	TY_INT   TypeKind = iota // int
	TY_PTR                   // pointer
	TY_ARRAY                 // array type
)

type Type struct {
	Kind      TypeKind
	PtrTo     *Type
	ArraySize uint16
}

func arrayOf(base *Type, size uint16) *Type {
	return &Type{
		Kind:      TY_ARRAY,
		PtrTo:     base,
		ArraySize: size,
	}
}

func sizeOf(ty *Type) int {
	if ty.Kind == TY_INT || ty.Kind == TY_PTR {
		return 8
	}
	return sizeOf(ty.PtrTo) * int(ty.ArraySize)
}

func intType() *Type {
	return &Type{Kind: TY_INT}
}

func pointerTo(base *Type) *Type {
	return &Type{Kind: TY_PTR, PtrTo: base}
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
		if node.Rhs.Ty != nil {
			if node.Rhs.Ty.PtrTo != nil {
				tmp := node.Lhs
				node.Lhs = node.Rhs
				node.Rhs = tmp
			}
		}

		if node.Rhs.Ty != nil {
			if node.Rhs.Ty.PtrTo != nil {
				e.err = fmt.Errorf(
					"e.visit(): err: \n%s",
					errorTok(node.Tok, "invalid pointer arithmetic operands"),
				)
			}
		}
		node.Ty = node.Lhs.Ty
		return
	case ND_SUB:
		if node.Rhs.Ty != nil {
			if node.Rhs.Ty.PtrTo != nil {
				e.err = fmt.Errorf(
					"e.visit(): err: \n%s",
					errorTok(node.Tok, "invalid pointer arithmetic operands"),
				)
			}
		}
		node.Ty = node.Lhs.Ty
		return
	case ND_ASSIGN:
		node.Ty = node.Lhs.Ty
		return
	case ND_ADDR:
		if node.Lhs.Ty.Kind == TY_ARRAY {
			node.Ty = pointerTo(node.Lhs.Ty.PtrTo)
		} else {
			node.Ty = pointerTo(node.Lhs.Ty)
		}
		return
	case ND_DEREF:
		if node.Lhs.Ty.PtrTo == nil {

			// fmt.Printf("node: %#v\n'%s'\n\n", node, node.Tok.Str)
			// fmt.Printf("node.Rhs: %#v\n'%s'\n\n", node.Rhs, node.Rhs.Tok.Str)
			// fmt.Printf("node.Lhs: %#v\n'%s'\n\n", node.Lhs, node.Lhs.Tok.Str)

			e.err = fmt.Errorf(
				"e.visit(): err: \n%s",
				errorTok(node.Tok, "invalid pointer dereference"),
			)
		}
		node.Ty = node.Lhs.Ty.PtrTo
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
