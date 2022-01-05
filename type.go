package main

import (
	"fmt"
	"math"
	"strconv"
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
	Sz    int // Sizeof() value
	Align int // alignment

	Base *Type // pointer or array

	// Declaration
	Name *Token

	TyName string

	ArrSz int     // Array size
	Mems  *Member // struct

	// function
	RetTy  *Type
	Params *Type
	Next   *Type
}

type Member struct {
	Next   *Member
	Ty     *Type
	Tok    *Token
	Name   string
	Idx    int
	Offset int
}

var ty_void *Type = &Type{Kind: TY_VOID, Sz: 1, Align: 1, TyName: "void"}
var ty_bool *Type = &Type{Kind: TY_BOOL, Sz: 1, Align: 1, TyName: "bool"}

var ty_char *Type = &Type{Kind: TY_BYTE, Sz: 1, Align: 1, TyName: "byte"}
var ty_short *Type = &Type{Kind: TY_SHORT, Sz: 2, Align: 2, TyName: "int16"}
var ty_int *Type = &Type{Kind: TY_INT, Sz: 4, Align: 4, TyName: "int"}
var ty_long *Type = &Type{Kind: TY_LONG, Sz: 8, Align: 8, TyName: "int64"}

func alignTo(n, align int) int {
	return (n + align - 1) / align * align
}

func newType(kind TypeKind, size, align int, name string) *Type {
	return &Type{Kind: kind, Align: align, TyName: name}
}

func isInteger(ty *Type) bool {
	k := ty.Kind
	return k == TY_BOOL || k == TY_BYTE || k == TY_SHORT ||
		k == TY_INT || k == TY_LONG
}

func copyType(ty *Type) *Type {
	ret := &Type{}
	ret = ty
	return ret
}

func voidType() *Type {
	return newType(TY_VOID, 1, 1, "void")
}

func boolType() *Type {
	return newType(TY_BOOL, 1, 1, "bool")
}

func charType() *Type {
	return newType(TY_BYTE, 1, 1, "byte")
}

func shortType() *Type {
	return newType(TY_SHORT, 2, 2, "int16")
}

func intType() *Type {
	return newType(TY_INT, 4, 4, "int")
}

func longType() *Type {
	return newType(TY_LONG, 8, 8, "int64")
}

func funcType(retTy *Type) *Type {
	return &Type{Kind: TY_FUNC, Align: 1, RetTy: retTy, TyName: "func"}
}

func stringType() *Type {
	return &Type{Kind: TY_PTR, Base: charType(), Align: 8, TyName: "string"}
}

func pointerTo(base *Type) *Type {
	return &Type{Kind: TY_PTR, Base: base, Align: 8, TyName: "pointer"}
}

func arrayOf(base *Type, size int) *Type {
	return &Type{
		Kind:   TY_ARRAY,
		Align:  base.Align,
		Base:   base,
		ArrSz:  size,
		TyName: "[" + strconv.Itoa(size) + "]" + base.TyName}
}

func structType() *Type {
	return newType(TY_STRUCT, 0, 1, "struct")
}

func sizeOf(ty *Type, tok *Token) int {
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
		return sizeOf(ty.Base, tok) * ty.ArrSz
	case TY_STRUCT:
		mem := ty.Mems
		for mem.Next != nil {
			mem = mem.Next
		}
		end := mem.Offset + sizeOf(mem.Ty, mem.Tok)
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

func getCommonType(ty1, ty2 *Type) *Type {
	if ty1.Base != nil {
		return pointerTo(ty1.Base)
	}
	if ty1.Sz == 8 || ty2.Sz == 8 {
		return ty_long
	}
	return ty_int
}

// For many binary operators, we implicitly promote operands sp that
// both operands have the same type. Any integral type smaller than
// int is always promoted to int. If the type of one operand is larger
// than the other's (e.g. "long" vs. "int"), the smaller operand will
// be promoted to match with the other.
//
// This operation is called the "usual arithmetic conversion".
func usualArithConv(lhs **Node, rhs **Node) {
	ty := getCommonType((*lhs).Ty, (*rhs).Ty)
	*lhs = newCast(*lhs, ty)
	*rhs = newCast(*rhs, ty)
}

func (e *errWriter) visit(node *Node) {
	if e.err != nil {
		return
	}

	if node == nil || node.Ty != nil {
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
		node.Ty = node.Obj.Ty
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
		ND_BITNOT:
		node.Ty = node.Lhs.Ty
		return
	case ND_COMMA:
		node.Ty = node.Rhs.Ty
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
		node.Val = int64(sizeOf(node.Lhs.Ty, node.Tok))
		node.Lhs = nil
		return
	}
}

func addType(node *Node) error {
	e := &errWriter{}
	e.visit(node)
	return e.err
}
