package main

import (
	"fmt"
	"math"
)

type TypeKind int

// errWriter is struct for error handling
// it's based on:
// https://jxck.hatenablog.com/entry/golang-error-handling-lesson-by-rob-pike
type errWriter struct {
	err error
}

const (
	TY_VOID  TypeKind = iota // void type
	TY_BOOL                  // bool type
	TY_BYTE                  // char type
	TY_SHORT                 // int16 type
	TY_INT                   // int32 type
	TY_LONG                  // int64 type
	TY_PTR                   // pointer type
	TY_ARRAY
	TY_STRUCT
	TY_FUNC
)

type Type struct {
	Kind  TypeKind
	Align int     // alignment
	Base  *Type   // pointer or array
	ArrSz int     // Array size
	Mems  *Member // struct
	RetTy *Type   // function
}

type Member struct {
	Next   *Member
	Ty     *Type
	Tok    *Token
	Name   string
	Offset int
}

func alignTo(n, align int) int {
	return (n + align - 1) & ^(align - 1)
}

func newType(kind TypeKind, align int) *Type {
	return &Type{Kind: kind, Align: align}
}

func voidType() *Type {
	return newType(TY_VOID, 1)
}

func boolType() *Type {
	return newType(TY_BOOL, 1)
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

func funcType(retTy *Type) *Type {
	return &Type{Kind: TY_FUNC, Align: 1, RetTy: retTy}
}

func pointerTo(base *Type) *Type {
	return &Type{Kind: TY_PTR, Base: base, Align: 8}
}

func arrayOf(base *Type, size int) *Type {
	return &Type{Kind: TY_ARRAY, Align: base.Align, Base: base, ArrSz: size}
}

func sizeOf(ty *Type) int {
	assert(ty.Kind != TY_VOID, "invalid void type")

	switch ty.Kind {
	case TY_BYTE, TY_BOOL:
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
	assert(ty.Kind == TY_STRUCT, "invalid type")

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
	case ND_MUL,
		ND_DIV,
		ND_BITAND,
		ND_BITOR,
		ND_BITXOR,
		ND_EQ,
		ND_NE,
		ND_LT,
		ND_LE,
		ND_NOT,
		ND_LOGOR,
		ND_LOGAND:
		node.Ty = intType()
		return
	case ND_NUM:
		if node.Val <= int64(math.MaxInt32) {
			node.Ty = intType()
			return
		}
		node.Ty = longType()
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
	case ND_ASSIGN,
		ND_SHL,
		ND_SHR,
		ND_INC,
		ND_DEC,
		ND_A_ADD,
		ND_A_SUB,
		ND_A_MUL,
		ND_A_DIV,
		ND_A_SHL,
		ND_A_SHR,
		ND_BITNOT:
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
		if node.Ty.Kind == TY_VOID {
			e.err = fmt.Errorf(errorTok(node.Tok, "dereference a void pointer"))
		}
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
