package main

import (
	"fmt"
)

type TypeKind int

// errWriter is struct for error handling
// it's based on:
// https://jxck.hatenablog.com/entry/golang-error-handling-lesson-by-rob-pike
type errWriter struct {
	err error
}

const (
	TY_VOID   TypeKind = iota // void type
	TY_BOOL                   // bool type
	TY_BYTE                   // char(int8) type
	TY_SHORT                  // int16 type
	TY_INT                    // int32 type
	TY_LONG                   // int64 type
	TY_FLOAT                  // float32 type
	TY_DOUBLE                 // float64 type
	TY_PTR                    // pointer type
	TY_ARRAY
	TY_STRUCT
	TY_FUNC
	TY_SLICE
)

type Type struct {
	Kind       TypeKind
	Sz         int  // Sizeof() value
	Align      int  // alignment
	IsUnsigned bool // unsigned or signed

	Base *Type // pointer or array

	// Declaration
	Name    *Token
	NamePos *Token

	TyName string

	ArrSz  int     // Array size
	Mems   *Member // struct
	IsFlex bool    // struct

	// function
	RetTy      *Type
	Params     *Type
	IsVariadic bool
	Next       *Type

	// slice
	Len      int
	Cap      int
	UArr     *Obj  // It's underlying array
	UArrIdx  int64 // Underlying array's index
	UArrNode *Node // Underlying array node

	Init *Initializer
}

type Member struct {
	Next   *Member
	Ty     *Type
	Tok    *Token // for error message
	Name   *Token
	Idx    int
	Align  int
	Offset int
}

var ty_void = &Type{Kind: TY_VOID, Sz: 1, Align: 1, TyName: "void"}
var ty_bool = &Type{Kind: TY_BOOL, Sz: 1, Align: 1, TyName: "bool"}

var ty_char = &Type{Kind: TY_BYTE, Sz: 1, Align: 1, TyName: "int8"}
var ty_short = &Type{Kind: TY_SHORT, Sz: 2, Align: 2, TyName: "int16"}
var ty_int = &Type{Kind: TY_INT, Sz: 4, Align: 4, TyName: "int"}
var ty_long = &Type{Kind: TY_LONG, Sz: 8, Align: 8, TyName: "int64"}

var ty_uchar = &Type{Kind: TY_BYTE, Sz: 1, Align: 1, TyName: "uint8", IsUnsigned: true}
var ty_ushort = &Type{Kind: TY_SHORT, Sz: 2, Align: 2, TyName: "uint16", IsUnsigned: true}
var ty_uint = &Type{Kind: TY_INT, Sz: 4, Align: 4, TyName: "uint", IsUnsigned: true}
var ty_ulong = &Type{Kind: TY_LONG, Sz: 8, Align: 8, TyName: "uint32", IsUnsigned: true}

var ty_float = &Type{Kind: TY_FLOAT, Sz: 4, Align: 4, TyName: "float32"}
var ty_double = &Type{Kind: TY_DOUBLE, Sz: 8, Align: 8, TyName: "float64"}

func newType(kind TypeKind, size, align int, name string, isUnsign bool) *Type {
	return &Type{Kind: kind, Sz: size, Align: align, TyName: name}
}

func isInteger(ty *Type) bool {
	k := ty.Kind
	return k == TY_BOOL || k == TY_BYTE || k == TY_SHORT ||
		k == TY_INT || k == TY_LONG
}

func isFlonum(ty *Type) bool {
	return ty.Kind == TY_FLOAT || ty.Kind == TY_DOUBLE
}

func isNumeric(ty *Type) bool {
	return isInteger(ty) || isFlonum(ty)
}

func copyType(ty *Type) *Type {
	ret := &Type{
		Kind:       ty.Kind,
		Sz:         ty.Sz,
		Align:      ty.Align,
		IsUnsigned: ty.IsUnsigned,
		Base:       ty.Base,
		Name:       ty.Name,
		TyName:     ty.TyName,
		ArrSz:      ty.ArrSz,
		Mems:       ty.Mems,
		RetTy:      ty.RetTy,
		Params:     ty.Params,
		IsVariadic: ty.IsVariadic,
		Next:       ty.Next,
	}
	return ret
}

func charType() *Type {
	return newType(TY_BYTE, 1, 1, "byte", true)
}

func pointerTo(base *Type) *Type {
	tyname := "*"
	for b := base; b != nil; b = b.Base {
		if b.TyName == "string" {
			tyname += "string"
			break
		} else if b.Kind == TY_PTR {
			tyname += "*"
		} else {
			tyname += b.TyName
			break
		}
	}
	return &Type{
		Kind:       TY_PTR,
		Base:       base,
		Sz:         8,
		Align:      8,
		TyName:     tyname,
		IsUnsigned: true,
	}
}

func funcType(retTy *Type, params *Type) *Type {
	head := params
	cur := head
	var typarams string
	for ; cur != nil; cur = cur.Next {
		typarams += cur.TyName + ","
	}

	var retty string
	if retTy != nil {
		retty = retTy.TyName
	}

	return &Type{
		Kind:   TY_FUNC,
		RetTy:  retTy,
		Params: head,
		TyName: "func(" + typarams + ")" + retty,
	}
}

func stringType() *Type {
	return &Type{Kind: TY_PTR, Base: charType(), Sz: 8, Align: 8, TyName: "string"}
}

func itoa(i int) string {
	var tmp = make([]rune, 0)
	var base int = 10
	for i > 0 {
		r := i%base + '0'
		tmp = append(tmp, rune(r))
		i /= base
	}
	var ret = make([]rune, 0)
	for j := len(tmp) - 1; j >= 0; j-- {
		ret = append(ret, tmp[j])
	}
	return string(ret)
}

func arrayOf(base *Type, len int) *Type {
	return &Type{
		Kind:   TY_ARRAY,
		Sz:     base.Sz * len,
		Align:  base.Align,
		Base:   base,
		ArrSz:  len,
		TyName: "[" + itoa(len) + "]" + base.TyName}
}

func structType() *Type {
	return newType(TY_STRUCT, 0, 1, "struct", false)
}

func sliceType(base *Type, len int, cap int) *Type {
	return &Type{
		Kind:   TY_SLICE,
		Sz:     8,
		Align:  8,
		Base:   base,
		ArrSz:  len,
		Len:    len,
		Cap:    cap,
		TyName: "[]" + base.TyName,
		IsFlex: true,
	}
}

func getCommonType(ty1, ty2 *Type) *Type {
	if ty1.Base != nil {
		return pointerTo(ty1.Base)
	}

	if ty1.Kind == TY_FUNC {
		return pointerTo(ty1)
	}
	if ty2.Kind == TY_FUNC {
		return pointerTo(ty2)
	}

	if ty1.Kind == TY_DOUBLE || ty2.Kind == TY_DOUBLE {
		return ty_double
	}
	if ty1.Kind == TY_FLOAT || ty2.Kind == TY_FLOAT {
		return ty_float
	}

	if ty1.Sz < 4 {
		ty1 = ty_int
	}
	if ty2.Sz < 4 {
		ty2 = ty_int
	}

	if ty1.Sz != ty2.Sz {
		if ty1.Sz < ty2.Sz {
			return ty2
		} else {
			return ty1
		}
	}

	if ty2.IsUnsigned {
		return ty2
	}
	return ty1
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
	case ND_NUM:
		node.Ty = ty_int
		return
	case ND_ADD,
		ND_SUB,
		ND_MUL,
		ND_DIV,
		ND_MOD,
		ND_BITAND,
		ND_BITOR,
		ND_BITXOR:
		usualArithConv(&node.Lhs, &node.Rhs)
		node.Ty = node.Lhs.Ty
		return
	case ND_NEG:
		ty := getCommonType(ty_int, node.Lhs.Ty)
		node.Lhs = newCast(node.Lhs, ty)
		node.Ty = ty
		return
	case ND_ASSIGN:
		if node.Lhs.Ty.Kind == TY_ARRAY {
			e.err = fmt.Errorf(errorTok(node.Lhs.Tok, "not an lvalue"))
		}
		if node.Lhs.Ty.Kind != TY_STRUCT {
			node.Rhs = newCast(node.Rhs, node.Lhs.Ty)
		}
		node.Ty = node.Lhs.Ty
		return
	case ND_EQ,
		ND_NE,
		ND_LT,
		ND_LE:
		usualArithConv(&node.Lhs, &node.Rhs)
		node.Ty = ty_bool
		return
	case ND_FUNCALL:
		node.Ty = ty_long
		return
	case ND_BITNOT,
		ND_SHL,
		ND_SHR:
		node.Ty = node.Lhs.Ty
		return
	case ND_NOT,
		ND_LOGOR,
		ND_LOGAND:
		node.Ty = ty_int
		return
	case ND_VAR:
		node.Ty = node.Obj.Ty
		return
	case ND_COND:
		if node.Then.Ty.Kind == TY_VOID || node.Els.Ty.Kind == TY_VOID {
			node.Ty = ty_void
			return
		}
		usualArithConv(&node.Then, &node.Els)
		node.Ty = node.Then.Ty
		return
	case ND_COMMA:
		node.Ty = node.Rhs.Ty
		return
	case ND_MEMBER:
		node.Ty = node.Mem.Ty
		return
	case ND_ADDR:
		ty := node.Lhs.Ty
		if ty.Kind == TY_ARRAY {
			node.Ty = pointerTo(node.Lhs.Ty.Base)
			return
		}
		node.Ty = pointerTo(ty)
		return
	case ND_DEREF:
		if node.Lhs.Ty.Base == nil {
			e.err = fmt.Errorf(errorTok(node.Tok, "invalid pointer dereference"))
			return
		}
		if node.Lhs.Ty.Base.Kind == TY_VOID {
			e.err = fmt.Errorf(errorTok(node.Tok, "dereference a void pointer"))
			return
		}

		node.Ty = node.Lhs.Ty.Base
		return
	case ND_STMT_EXPR:
		if node.Body != nil {
			stmt := node.Body
			for stmt.Next != nil {
				stmt = stmt.Next
			}
			if stmt.Kind == ND_EXPR_STMT {
				node.Ty = stmt.Lhs.Ty
				return
			}
		}
		e.err = fmt.Errorf(errorTok(node.Tok,
			"statement expressionreturning void is not supported"))
		return
	}
}

func addType(node *Node) error {
	e := &errWriter{}
	e.visit(node)
	return e.err
}
