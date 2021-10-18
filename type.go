// add 'Type' to AST
package main

import (
	"fmt"
	"os"
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

func intType() *Type {
	return &Type{Kind: TY_INT}
}

func pointerTo(base *Type) *Type {
	return &Type{Kind: TY_PTR, Base: base}
}

func visit(node *Node) {
	if node == nil {
		return
	}

	visit(node.Lhs)
	visit(node.Rhs)
	visit(node.Cond)
	visit(node.Then)
	visit(node.Els)
	visit(node.Init)
	visit(node.Inc)

	for n := node.Body; n != nil; n = n.Next {
		visit(n)
	}
	for n := node.Args; n != nil; n = n.Next {
		visit(n)
	}

	switch node.Kind {
	case ND_MUL:
	case ND_DIV:
	case ND_EQ:
	case ND_NE:
	case ND_LT:
	case ND_LE:
	case ND_FUNCCALL:
	case ND_NUM:
		node.Ty = intType()
		return
	case ND_LVAR:
		node.Ty = node.Var.Ty
		return
	case ND_ADD:
		fmt.Printf("%#v\n'%s'\n\n", node, node.Tok.Str)
		fmt.Printf("Lhs: %#v\n'%s'\n\n", node.Lhs, node.Lhs.Tok.Str)
		fmt.Printf("Rhs: %#v\n'%s'\n\n", node.Rhs, node.Rhs.Tok.Str)

		if node.Rhs.Ty.Kind == TY_PTR {
			tmp := node.Lhs
			node.Lhs = node.Rhs
			node.Rhs = tmp
		}
		if node.Rhs.Ty.Kind == TY_PTR {
			errorTok(os.Stderr, node.Tok, "invalid pointer arithmetic operands")
		}
		node.Ty = node.Lhs.Ty
		return
	case ND_SUB:
		if node.Rhs.Ty.Kind == TY_PTR {
			errorTok(os.Stderr, node.Tok, "invalid pointer arithmetic operands")
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
			errorTok(os.Stderr, node.Tok, "invalid pointer dereference")
		}
		node.Ty = node.Lhs.Ty.Base
		return
	}
}

func addType(prog *Function) {
	for fn := prog; fn != nil; fn = fn.Next {
		for node := fn.Node; node != nil; node = node.Next {
			visit(node)
		}
	}
}
