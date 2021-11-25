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
	TY_BYTE TypeKind = iota // char type
	TY_SHORT
	TY_INT
	TY_LONG
	TY_PTR
	TY_ARRAY
	TY_STRUCT
)

type Type struct {
	Kind  TypeKind
	Align int
	Base  *Type
	ArrSz int // Array size
	Mems  *Member
}

type Member struct {
	Next   *Member
	Ty     *Type
	Name   string
	Offset int
}

func alignTo(n, align int) int {
	return (n + align - 1) & ^(align - 1)
}

func newType(kind TypeKind, align int) *Type {
	return &Type{Kind: kind, Align: align}
}

func charType() *Type {
	return newType(TY_BYTE, 1)
}

func shortType() *Type {
	return newType(TY_SHORT, 2)
}

func intType() *Type {
	return newType(TY_INT, 4)
}

func longType() *Type {
	return newType(TY_LONG, 8)
}

func pointerTo(base *Type) *Type {
	return &Type{Kind: TY_PTR, Base: base, Align: 8}
}

func arrayOf(base *Type, size int) *Type {
	return &Type{Kind: TY_ARRAY, Align: base.Align, Base: base, ArrSz: size}
}

func sizeOf(ty *Type) int {
	switch ty.Kind {
	case TY_BYTE:
		return 1
	case TY_SHORT:
		return 2
	case TY_INT:
		return 4
	case TY_PTR, TY_LONG:
		return 8
	case TY_ARRAY:
		return sizeOf(ty.Base) * ty.ArrSz
	case TY_STRUCT:
		mem := ty.Mems
		for mem.Next != nil {
			mem = mem.Next
		}
		end := mem.Offset + sizeOf(mem.Ty)
		return alignTo(end, ty.Align)
	default:
		panic("invalid type")
	}
}

func findMember(ty *Type, name string) *Member {
	if ty.Kind != TY_STRUCT {
		panic("invalid type")
	}
	for mem := ty.Mems; mem != nil; mem = mem.Next {
		if mem.Name == name {
			return mem
		}
	}
	return nil
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
	case ND_MEMBER:
		if node.Lhs.Ty.Kind != TY_STRUCT {
			e.err = fmt.Errorf(errorTok(node.Tok, "not a struct"))
		}
		node.Mem = findMember(node.Lhs.Ty, node.MemName)
		if node.Mem == nil {
			e.err = fmt.Errorf(errorTok(node.Tok, "specified member does not exist"))
			return
		}
		node.Ty = node.Mem.Ty
		return
	case ND_ADDR:
		if node.Lhs.Ty.Kind == TY_ARRAY {
			node.Ty = pointerTo(node.Lhs.Ty.Base)
			return
		}
		node.Ty = pointerTo(node.Lhs.Ty)
		return
	case ND_DEREF:
		if node.Lhs.Ty.Base == nil {
			e.err = fmt.Errorf(errorTok(node.Tok, "invalid pointer dereference"))
		}
		node.Ty = node.Lhs.Ty.Base
		return
	case ND_SIZEOF:
		node.Kind = ND_NUM
		node.Ty = intType()
		node.Val = int64(sizeOf(node.Lhs.Ty))
		node.Lhs = nil
		return
	}
}

func addType(prog *Program) error {
	e := &errWriter{}

	for fn := prog.Fns; fn != nil; fn = fn.Next {
		for node := fn.Node; node != nil; node = node.Next {
			e.visit(node)
		}
	}
	return e.err
}
