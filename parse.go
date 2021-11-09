//
// AST parser
//
package main

import "fmt"

// the types of AST node
type NodeKind int

const (
	ND_ADD       NodeKind = iota // 0: +
	ND_SUB                       // 1: -
	ND_MUL                       // 2: *
	ND_DIV                       // 3: /
	ND_EQ                        // 4: ==
	ND_NE                        // 5: !=
	ND_LT                        // 6: <
	ND_LE                        // 7: <=
	ND_ASSIGN                    // 8: =
	ND_PRE_INC                   // 9: pre ++
	ND_PRE_DEC                   // 10: pre --
	ND_POST_INC                  // 11: post ++
	ND_POST_DEC                  // 12: post --
	ND_COMMA                     // 13: ,
	ND_VAR                       // 14: local or global variables
	ND_NUM                       // 15: integer
	ND_RETURN                    // 16: 'return'
	ND_IF                        // 17: "if"
	ND_WHILE                     // 18: "while"
	ND_FOR                       // 19: "for"
	ND_BLOCK                     // 20: {...}
	ND_FUNCCALL                  // 21: function call
	ND_MEMBER                    // 22: . (struct member access)
	ND_ADDR                      // 23: unary &
	ND_DEREF                     // 24: unary *
	ND_EXPR_STMT                 // 25: expression statement
	ND_STMT_EXPR                 // 26: statement expression
	ND_CAST                      // 27: type cast
	ND_NULL                      // 28: empty statement
	ND_SIZEOF                    // 29: "sizeof" operator
)

// define AST node
type Node struct {
	Kind NodeKind // the type of node
	Next *Node    // the next node
	Ty   *Type    // the data type
	Tok  *Token   // current token

	Lhs *Node // the left branch
	Rhs *Node // the right branch

	// "if" or "while" of "for" statement
	Cond *Node
	Then *Node
	Els  *Node
	Init *Node
	Inc  *Node

	// block or statement expression
	Body *Node

	// struct member access
	MemName string
	Mem     *Member

	// for function call
	FuncName string
	Args     *Node

	Val int64 // it would be used when 'Kind' is 'ND_NUM'
	Var *Var  // it would be used when 'Kind' is 'ND_VAR'
}

func newNode(kind NodeKind, lhs *Node, rhs *Node, tok *Token) *Node {
	return &Node{
		Kind: kind,
		Lhs:  lhs,
		Rhs:  rhs,
		Tok:  tok,
	}
}

func newNodeNum(val int64, tok *Token) *Node {
	return &Node{
		Kind: ND_NUM,
		Val:  val,
		Tok:  tok,
	}
}

func newUnary(kind NodeKind, expr *Node, tok *Token) *Node {
	return &Node{Kind: kind, Tok: tok, Lhs: expr}
}

// the type of local variables
type Var struct {
	Name    string // the name of the variable
	Ty      *Type  // the data type
	IsLocal bool   // local or global

	// local variables
	Offset int // the offset from RBP

	// global vaiables
	Contents []rune
	ContLen  int
}

type VarList struct {
	Next *VarList
	Var  *Var
}

// scope for local variables, global variables or typedefs
type VarScope struct {
	Next    *VarScope
	Name    string
	Var     *Var
	TyDef   *Type
	EnumTy  *Type
	EnumVal int
}

// scope for struct tags
type TagScope struct {
	Next *TagScope
	Name string
	Ty   *Type
}

// local variables
var locals *VarList
var globals *VarList

var varScope *VarScope
var tagScope *TagScope

// findVar searchs a variable by name.
// if it wasn't find, return nil.
func findVar(tok *Token) *VarScope {
	for sc := varScope; sc != nil; sc = sc.Next {
		if len(sc.Name) == tok.Len && tok.Str == sc.Name {
			return sc
		}
	}
	return nil
}

func findTag(tok *Token) *TagScope {
	for sc := tagScope; sc != nil; sc = sc.Next {
		if len(sc.Name) == tok.Len && tok.Str == sc.Name {
			return sc
		}
	}
	return nil
}

func newVar(lvar *Var, tok *Token) *Node {
	return &Node{Kind: ND_VAR, Tok: tok, Var: lvar}
}

func pushScope(name string) *VarScope {
	sc := &VarScope{
		Name: name,
		Next: varScope,
	}
	varScope = sc
	return sc
}

func pushVar(name string, ty *Type, isLocal bool) *Var {
	lvar := &Var{
		Name:    name,
		Ty:      ty,
		IsLocal: isLocal,
	}

	vl := &VarList{Var: lvar}

	if isLocal {
		vl.Next = locals
		locals = vl
	} else if ty.Kind != TY_FUNC {
		vl.Next = globals
		globals = vl
	}

	return lvar
}

func findTypedef(tok *Token) *Type {
	if tok.Kind == TK_IDENT {
		sc := findVar(token)
		if sc != nil {
			return sc.TyDef
		}
	}
	return nil
}

// for newLabel function
var cnt int

func newLabel() string {
	res := fmt.Sprintf(".L.data.%d", cnt)
	cnt++
	return res
}

type Function struct {
	Next   *Function
	Name   string
	Params *VarList

	Node    *Node
	Locals  *VarList
	StackSz int
}

type Program struct {
	Globals *VarList
	Fns     *Function
}

func isFunction() bool {
	tok := token

	ty := typeSpecifier()
	var name string
	declarator(ty, &name)
	isFunc := name != "" && consume("(") != nil

	token = tok
	return isFunc
}

// program = (global-var | function*)
func program() *Program {
	// printCurTok()
	// printCurFunc()
	cur := &Function{}
	head := cur
	globals = nil

	for !atEof() {
		if isFunction() {
			fn := function()
			if fn == nil {
				continue
			}
			cur.Next = fn
			cur = cur.Next
			continue
		}

		globalVar()
	}

	prog := &Program{Globals: globals, Fns: head.Next}
	return prog
}

// type-specifier = builtin-type | struct-decl | typedef-name | enum-specifier
//
// builtin-type   = "void"
//                | "_Bool"
//                | "char"
//                | "short" | "short" "int" | "int" "short"
//                | "int"
//                | "long" | "long" "int" | "int" "long"
//
// node that "typedef" and "static" can appear anywhere in a type-specifier
func typeSpecifier() *Type {
	if !isTypename() {
		panic("\n" + errorTok(token, "typename expected"))
	}

	var ty *Type = nil

	const (
		VOID  = 1 << 1
		BOOL  = 1 << 3
		CHAR  = 1 << 5
		SHORT = 1 << 7
		INT   = 1 << 9
		LONG  = 1 << 11
	)

	baseTy := 0
	var userTy *Type

	isTypedef := false
	isStatic := false

	for {
		// read one token at a time.
		tok := token
		if consume("typedef") != nil {
			isTypedef = true
		} else if consume("static") != nil {
			isStatic = true
		} else if consume("void") != nil {
			baseTy += VOID
		} else if consume("_Bool") != nil {
			baseTy += BOOL
		} else if consume("char") != nil {
			baseTy += CHAR
		} else if consume("short") != nil {
			baseTy += SHORT
		} else if consume("int") != nil {
			baseTy += INT
		} else if consume("long") != nil {
			baseTy += LONG
		} else if peek("struct") != nil {
			if baseTy != 0 || userTy != nil {
				break
			}
			userTy = structDecl()
		} else if peek("enum") != nil {
			if baseTy != 0 || userTy != nil {
				break
			}
			userTy = enumSpecifier()
		} else {
			if baseTy != 0 || userTy != nil {
				break
			}
			ty_ := findTypedef(token)
			if ty_ == nil {
				break
			}
			token = token.Next
			userTy = ty_
		}

		switch baseTy {
		case VOID:
			ty = voidType()
		case BOOL:
			ty = boolType()
		case CHAR:
			ty = charType()
		case SHORT, SHORT + INT:
			ty = shortType()
		case INT:
			ty = intType()
		case LONG, LONG + INT:
			ty = longType()
		case 0:
			// if there's no type specifier, it becomes int.
			// for expample, 'typedef x' defines x as an alias for int.
			if userTy != nil {
				ty = userTy
			} else {
				ty = intType()
			}
		default:
			panic("\n" + errorTok(tok, "invalid type"))
		}
	}
	ty.IsTypedef = isTypedef
	ty.IsStatic = isStatic
	return ty
}

// declarator = "*" ("(" declarator ")") | ident) type-suffix
func declarator(ty *Type, name *string) *Type {
	for consume("*") != nil {
		ty = pointerTo(ty)
	}

	if consume("(") != nil {
		placeholder := &Type{}
		newTy := declarator(placeholder, name)
		expect(")")
		*placeholder = *typeSuffix(ty)
		return newTy
	}

	*name = expectIdent()
	return typeSuffix(ty)
}

// abstract-declarator = "*"* ("(" abstract-declarator ")")? type-suffix
func abstractDeclarator(ty *Type) *Type {
	for consume("*") != nil {
		ty = pointerTo(ty)
	}

	if consume("(") != nil {
		placeholder := &Type{}
		newTy := abstractDeclarator(placeholder)
		expect(")")
		*placeholder = *typeSuffix(ty)
		return newTy
	}
	return typeSuffix(ty)
}

// type-suffix = ("[" num "]" type-suffix)?
func typeSuffix(ty *Type) *Type {
	if consume("[") == nil {
		return ty
	}
	sz := expectNumber()
	expect("]")
	ty = typeSuffix(ty)
	return arrayOf(ty, uint16(sz))
}

func typeName() *Type {
	ty := typeSpecifier()
	ty = abstractDeclarator(ty)
	return typeSuffix(ty)
}

func pushTagScope(tok *Token, ty *Type) {
	sc := &TagScope{
		Next: tagScope,
		Name: tok.Str,
		Ty:   ty,
	}
	tagScope = sc
}

// struct-decl = "struct" ident
//             | "struct" ident? "{" struct-member "}"
func structDecl() *Type {

	// read struct tag.
	expect("struct")
	tag := consumeIdent()
	if tag != nil && peek("{") == nil {
		sc := findTag(tag)
		if sc == nil {
			panic("\n" + errorTok(tag, "unknown struct type"))
		}
		if sc.Ty.Kind != TY_STRUCT {
			panic("\n" + errorTok(tag, "not a struct tag"))
		}
		return sc.Ty
	}
	expect("{")

	// read struct members.
	head := &Member{}
	cur := head

	for consume("}") == nil {
		cur.Next = structMember()
		cur = cur.Next
	}

	ty := &Type{Kind: TY_STRUCT, Mems: head.Next}

	// assign offsets within the struct to members.
	offset := 0
	for mem := ty.Mems; mem != nil; mem = mem.Next {
		offset = alignTo(offset, mem.Ty.Align)
		mem.Offset = offset
		offset += sizeOf(mem.Ty)

		if ty.Align < mem.Ty.Align {
			ty.Align = mem.Ty.Align
		}
	}

	// register the struct type if a name was given.
	if tag != nil {
		pushTagScope(tag, ty)
	}
	return ty
}

// enum-specifier = "enum" ident
//                | "enum" ident? "{" enum-list? "}"
//
// enum-list = ident ("=" num)? ("," ident ("=" num)?)*
func enumSpecifier() *Type {
	expect("enum")
	ty := enumType()

	// read an enum tag
	tag := consumeIdent()
	if tag != nil && peek("{") == nil {
		sc := findTag(tag)
		if sc == nil {
			panic("\n" + errorTok(tag, "unknown enum type"))
		}
		if sc.Ty.Kind != TY_ENUM {
			panic("\n" + errorTok(tag, "not an enum tag"))
		}
		return sc.Ty
	}

	expect("{")

	// read enum-list
	cnt := 0
	for {
		name := expectIdent()
		if consume("=") != nil {
			cnt = int(expectNumber())
		}

		sc := pushScope(name)
		sc.EnumTy = ty
		sc.EnumVal = cnt
		cnt++

		if consume(",") != nil {
			if consume("}") != nil {
				break
			}
			continue
		}
		expect("}")
		break
	}

	if tag != nil {
		pushTagScope(tag, ty)
	}
	return ty
}

// struct-member = type-specifier declarator type-suffix ";"
func structMember() *Member {
	var ty *Type = typeSpecifier()
	var name string
	ty = declarator(ty, &name)
	ty = typeSuffix(ty)
	expect(";")

	mem := &Member{Ty: ty, Name: name}
	return mem
}

// param = type-specifier declarator type-suffix
func readFuncParam() *VarList {
	ty := typeSpecifier()
	var name string
	ty = declarator(ty, &name)
	ty = typeSuffix(ty)

	var_ := pushVar(name, ty, true)
	pushScope(name).Var = var_

	vl := &VarList{Var: var_}
	return vl
}

// params = param ("," param)*
func readFuncParams() *VarList {
	// printCurTok()
	// printCurFunc()
	if consume(")") != nil { // no argument
		return nil
	}

	head := readFuncParam()
	cur := head

	for {
		if consume(")") != nil {
			break
		}
		expect(",")
		cur.Next = readFuncParam()
		cur = cur.Next
	}

	return head
}

// function = type-specifier declarator "(" params? ")" ("{" stmt* "}" | ";")
func function() *Function {
	// printCurTok()
	// printCurFunc()
	locals = nil

	ty := typeSpecifier()
	var name string
	ty = declarator(ty, &name)

	// add a function type to the scope
	var_ := pushVar(name, funcType(ty), false)
	pushScope(name).Var = var_

	// construct a function object
	fn := &Function{Name: name}
	expect("(")
	fn.Params = readFuncParams()

	if consume(";") != nil {
		return nil
	}

	// read function body
	cur := &Node{}
	head := cur
	expect("{")

	for {
		if t := consume("}"); t != nil {
			break
		}
		cur.Next = stmt()
		cur = cur.Next
	}

	fn.Node = head.Next
	fn.Locals = locals
	return fn
}

// global-var = type-specifier declarator type-suffix ";"
func globalVar() {
	ty := typeSpecifier()
	var name string
	ty = declarator(ty, &name)
	ty = typeSuffix(ty)
	expect(";")

	var_ := pushVar(name, ty, false)
	pushScope(name).Var = var_
}

// declaration = type-specifier declarator type-suffix ("=" expr)? ";"
//             | type-specifier ";"
func declaration() *Node {
	tok := token
	ty := typeSpecifier()
	if consume(";") != nil {
		return &Node{Kind: ND_NULL, Tok: tok}
	}

	var name string
	ty = declarator(ty, &name)
	ty = typeSuffix(ty)

	if ty.IsTypedef {
		expect(";")
		ty.IsTypedef = false
		pushScope(name).TyDef = ty
		return &Node{Kind: ND_NULL, Tok: tok}
	}

	if ty.Kind == TY_VOID {
		panic("\n" + errorTok(tok, "variable declared void"))
	}

	var var_ *Var
	if ty.IsStatic {
		var_ = pushVar(newLabel(), ty, false)
	} else {
		var_ = pushVar(name, ty, true)
	}
	pushScope(name).Var = var_

	if consume(";") != nil {
		return &Node{Kind: ND_NULL, Tok: tok}
	}

	expect("=")

	lhs := newVar(var_, tok)
	rhs := expr()
	expect(";")
	node := newNode(ND_ASSIGN, lhs, rhs, tok)
	return newUnary(ND_EXPR_STMT, node, tok)
}

func readExprStmt() *Node {
	tok := token
	return &Node{Kind: ND_EXPR_STMT, Lhs: expr(), Tok: tok}
}

func isTypename() bool {
	return peek("void") != nil || peek("_Bool") != nil || peek("char") != nil ||
		peek("short") != nil || peek("int") != nil || peek("long") != nil ||
		peek("enum") != nil || peek("struct") != nil || peek("typedef") != nil ||
		peek("static") != nil || findTypedef(token) != nil
}

// stmt = "return" expr ";"
//      | "if" "(" expr ")" stmt ("else" stmt)?
//      | "while" "(" expr ")" stmt
//      | "for" "(" (expr? ";" | declaration) expr? ";" expr? ")" stmt
//      | "{" stmt* "}"
//      | declaration
//      | expr ";"
func stmt() *Node {
	// printCurTok()
	// printCurFunc()
	var node *Node

	if t := consume("return"); t != nil {

		node = &Node{Kind: ND_RETURN, Lhs: expr(), Tok: t}
		expect(";")

	} else if t := consume("if"); t != nil {

		node = &Node{Kind: ND_IF, Tok: t}
		expect("(")
		node.Cond = expr()
		expect(")")
		node.Then = stmt()

		if consume("else") != nil {
			node.Els = stmt()
		}

	} else if t := consume("while"); t != nil {

		node = &Node{Kind: ND_WHILE, Tok: t}
		expect("(")
		node.Cond = expr()
		expect(")")
		node.Then = stmt()

	} else if t := consume("for"); t != nil {

		node = &Node{Kind: ND_FOR, Tok: t}
		expect("(")

		sc1 := varScope
		sc2 := tagScope

		if consume(";") == nil {
			if isTypename() {
				node.Init = declaration()
			} else {
				node.Init = readExprStmt()
				expect(";")
			}
		}
		if consume(";") == nil {
			node.Cond = expr()
			expect(";")
		}
		if consume(")") == nil {
			node.Inc = readExprStmt()
			expect(")")
		}
		node.Then = stmt()

		varScope = sc1
		tagScope = sc2

	} else if t := consume("{"); t != nil {

		head := Node{}
		cur := &head

		sc1 := varScope
		sc2 := tagScope
		for {
			if consume("}") != nil {
				break
			}
			cur.Next = stmt()
			cur = cur.Next
		}
		varScope = sc1
		tagScope = sc2

		node = &Node{Kind: ND_BLOCK, Tok: t}
		node.Body = head.Next

	} else {

		if isTypename() {
			return declaration()
		}

		node = readExprStmt()
		expect(";")
	}

	return node
}

// expr       = assign ("," assign)*
func expr() *Node {
	// printCurTok()
	// printCurFunc()
	node := assign()
	for {
		tok := consume(",")
		if tok == nil {
			break
		}
		node = newUnary(ND_EXPR_STMT, node, node.Tok)
		node = newNode(ND_COMMA, node, assign(), tok)
	}
	return node
}

// assign     = equality ("=" assign)?
func assign() *Node {
	// printCurTok()
	// printCurFunc()
	node := equality()
	if t := consume("="); t != nil {
		node = newNode(ND_ASSIGN, node, assign(), t)
	}

	return node
}

// equality   = relational ("==" relational | "!=" relational)*
func equality() *Node {
	// printCurTok()
	// printCurFunc()
	node := relational()

	for {
		if t := consume("=="); t != nil {
			node = newNode(ND_EQ, node, relational(), t)
		} else if consume("!=") != nil {
			node = newNode(ND_NE, node, relational(), t)
		} else {
			return node
		}
	}
}

// relational = add ("<" add | "<=" add | ">" add | ">=" add)*
func relational() *Node {
	// printCurTok()
	// printCurFunc()
	node := add()

	for {
		if t := consume("<"); t != nil {
			node = newNode(ND_LT, node, add(), t)
		} else if t := consume("<="); t != nil {
			node = newNode(ND_LE, node, add(), t)
		} else if t := consume(">"); t != nil {
			node = newNode(ND_LT, add(), node, t)
		} else if t := consume(">="); t != nil {
			node = newNode(ND_LE, add(), node, t)
		} else {
			return node
		}
	}
}

// add = mul ("+" mul | "-" mul)*
func add() *Node {
	// printCurTok()
	// printCurFunc()
	node := mul()

	for {
		if t := consume("+"); t != nil {
			node = newNode(ND_ADD, node, mul(), t)
		} else if t := consume("-"); t != nil {
			node = newNode(ND_SUB, node, mul(), t)
		} else {
			return node
		}
	}
}

// mul = cast ("*" cast | "/" cast)*
func mul() *Node {
	// printCurTok()
	// printCurFunc()
	node := cast()

	for {
		if t := consume("*"); t != nil {
			node = newNode(ND_MUL, node, cast(), t)
		} else if consume("/") != nil {
			node = newNode(ND_DIV, node, cast(), t)
		} else {
			return node
		}
	}
}

// cast = "(" type-name ")" cast | unary
func cast() *Node {
	tok := token

	if consume("(") != nil {
		if isTypename() {
			ty := typeName()
			expect(")")
			node := newUnary(ND_CAST, cast(), tok)
			node.Ty = ty
			return node
		}
		token = tok
	}

	return unary()
}

// unary = ("+" | "-" | "*" | "&")? cast
//       | ("++" | "--") unary
//       | "sizeof" "(" type-name ")"
//       | "sizeof" unary
//       | postfix
func unary() *Node {
	// printCurTok()
	// printCurFunc()
	if t := consumeSizeof(); t != nil {
		if consume("(") != nil {
			if isTypename() {
				ty := typeName()
				expect(")")
				return newNodeNum(int64(sizeOf(ty)), t)
			}
			token = t.Next
		}
		return newUnary(ND_SIZEOF, unary(), t)
	}

	if t := consume("+"); t != nil {
		return cast()
	}
	if t := consume("-"); t != nil {
		return newNode(ND_SUB, newNodeNum(0, t), cast(), t)
	}
	if t := consume("&"); t != nil {
		return newUnary(ND_ADDR, cast(), t)
	}
	if t := consume("*"); t != nil {
		return newUnary(ND_DEREF, cast(), t)
	}
	if t := consume("++"); t != nil {
		return newUnary(ND_PRE_INC, unary(), t)
	}
	if t := consume("--"); t != nil {
		return newUnary(ND_PRE_DEC, unary(), t)
	}
	return postfix()
}

// postfix = primary ("[" expr "]" | "." ident | "->" ident | "++" | "--")*
func postfix() *Node {
	node := primary()

	for {
		if tok := consume("["); tok != nil {
			// x[y] is short for *(x+y)
			exp := newNode(ND_ADD, node, expr(), tok)
			expect("]")
			node = newUnary(ND_DEREF, exp, tok)
			continue
		}

		if tok := consume("."); tok != nil {
			node = newUnary(ND_MEMBER, node, tok)
			node.MemName = expectIdent()
			continue
		}

		if tok := consume("->"); tok != nil {
			// x->y is shrot for (*x).y
			node = newUnary(ND_DEREF, node, tok)
			node = newUnary(ND_MEMBER, node, tok)
			node.MemName = expectIdent()
			continue
		}

		if tok := consume("++"); tok != nil {
			node = newUnary(ND_POST_INC, node, tok)
			continue
		}

		if tok := consume("--"); tok != nil {
			node = newUnary(ND_POST_DEC, node, tok)
			continue
		}

		return node
	}
}

// stmt-expr = "(" "{" stmt stmt* "}" ")"
//
// statement expression is a GNU extension.
func stmtExpr(tok *Token) *Node {
	sc1 := varScope
	sc2 := tagScope

	node := &Node{
		Kind: ND_STMT_EXPR,
		Tok:  tok,
		Body: stmt(),
	}
	cur := node.Body

	for {
		if consume("}") != nil {
			break
		}
		cur.Next = stmt()
		cur = cur.Next
	}
	expect(")")

	varScope = sc1
	tagScope = sc2

	if cur.Kind != ND_EXPR_STMT {
		panic("\n" +
			errorTok(cur.Tok, "stmt expr returning void is not supported"))
	}
	*cur = *cur.Lhs
	return node
}

// func-args = "(" (assign("," assign)*)? ")"
func funcArgs() *Node {
	// printCurTok()
	// printCurFunc()
	if consume(")") != nil {
		return nil
	}

	head := assign()
	cur := head
	for {
		if consume(",") == nil {
			break
		}
		cur.Next = assign()
		cur = cur.Next
	}
	expect(")")
	return head
}

// primary = "(" "{" stmt-expr-tail
//         | ident func-args?
//         | "(" expr ")"
//         | num
//         | str
func primary() *Node {
	// printCurTok()
	// printCurFunc()

	// if the next token is '(', the program must be
	// "(" expr ")"
	if t := consume("("); t != nil {
		if consume("{") != nil {
			return stmtExpr(t)
		}

		node := expr()
		expect(")")
		return node
	}

	if tok := consumeIdent(); tok != nil {
		var node *Node
		if t := consume("("); t != nil { // function call
			node = &Node{
				Kind:     ND_FUNCCALL,
				Tok:      tok,
				FuncName: tok.Str,
				Args:     funcArgs(),
			}

			sc := findVar(tok)
			if sc != nil {
				if sc.Var == nil || sc.Var.Ty.Kind != TY_FUNC {
					panic("\n" + errorTok(tok, "not a function"))
				}
				node.Ty = sc.Var.Ty.RetTy
			} else {
				node.Ty = intType()
			}
			return node
		}

		// local variables
		sc := findVar(tok)
		if sc != nil {
			if sc.Var != nil {
				return newVar(sc.Var, tok)
			}
			if sc.EnumTy != nil {
				return newNodeNum(int64(sc.EnumVal), tok)
			}
		}
		panic("\n" + errorTok(tok, "undefined variable"))
	}

	tok := token
	if tok.Kind == TK_STR {
		token = token.Next

		ty := arrayOf(charType(), uint16(tok.ContLen))
		var_ := pushVar(newLabel(), ty, false)
		var_.Contents = tok.Contents
		var_.ContLen = tok.ContLen
		return newVar(var_, tok)
	}

	if tok.Kind != TK_NUM {
		panic("\n" + errorTok(tok, "expected expression"))
	}
	// otherwise, must be integer
	return newNodeNum(expectNumber(), tok)
}
