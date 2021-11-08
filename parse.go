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
	ND_VAR                       // 9: local or global variables
	ND_NUM                       // 10: integer
	ND_RETURN                    // 11: 'return'
	ND_IF                        // 12: "if"
	ND_WHILE                     // 13: "while"
	ND_FOR                       // 14: "for"
	ND_BLOCK                     // 15: {...}
	ND_FUNCCALL                  // 16: function call
	ND_MEMBER                    // 17: . (struct member access)
	ND_ADDR                      // 18: unary &
	ND_DEREF                     // 19: unary *
	ND_EXPR_STMT                 // 20: expression statement
	ND_STMT_EXPR                 // 21: statement expression
	ND_NULL                      // 22: empty statement
	ND_SIZEOF                    // 23: "sizeof" operator
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

	Val int  // it would be used when 'Kind' is 'ND_NUM'
	Var *Var // it would be used when 'Kind' is 'ND_LVAR'
}

func newNode(kind NodeKind, lhs *Node, rhs *Node, tok *Token) *Node {
	return &Node{
		Kind: kind,
		Lhs:  lhs,
		Rhs:  rhs,
		Tok:  tok,
	}
}

func newNodeNum(val int, tok *Token) *Node {
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
	Next  *VarScope
	Name  string
	Var   *Var
	TyDef *Type
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

	pushScope(name).Var = lvar
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
			cur.Next = function()
			cur = cur.Next
		} else {
			globalVar()
		}
	}

	prog := &Program{Globals: globals, Fns: head.Next}
	return prog
}

// type-specifier = builtin-type | struct-decl | typedef-name
// builtin-type   = "char" | "short" | "int" | "long"
func typeSpecifier() *Type {
	if !isTypename() {
		panic("\n" + errorTok(token, "typename expected"))
	}

	if consume("char") != nil {
		return charType()
	} else if consume("short") != nil {
		return shortType()
	} else if consume("int") != nil {
		return intType()
	} else if consume("long") != nil {
		return longType()
	} else if consume("struct") != nil {
		return structDecl()
	}
	return findVar(consumeIdent()).TyDef
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
	tag := consumeIdent()
	if tag != nil && peek("{") == nil {
		sc := findTag(tag)
		if sc == nil {
			panic("\n" + errorTok(tag, "unknown struct type"))
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

// struct-member = basetype ident ("{" num "}")* ";"
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

	vl := &VarList{Var: pushVar(name, ty, true)}
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

// function = type-specifier declarator "(" params? ")" "{" stmt* "}"
func function() *Function {
	// printCurTok()
	// printCurFunc()
	locals = nil

	ty := typeSpecifier()
	var name string
	ty = declarator(ty, &name)

	// add a function type to the scope
	pushVar(name, funcType(ty), false)

	// construct a function object
	fn := &Function{Name: name}
	expect("(")
	fn.Params = readFuncParams()
	expect("{")

	// read function body
	cur := &Node{}
	head := cur

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
	pushVar(name, ty, false)
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
	lvar := pushVar(name, ty, true)

	if consume(";") != nil {
		return &Node{Kind: ND_NULL, Tok: tok}
	}

	expect("=")
	lhs := newVar(lvar, tok)
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
	return peek("char") != nil ||
		peek("short") != nil ||
		peek("int") != nil ||
		peek("long") != nil ||
		peek("struct") != nil ||
		findTypedef(token) != nil
}

// stmt = "return" expr ";"
//      | "if" "(" expr ")" stmt ("else" stmt)?
//      | "while" "(" expr ")" stmt
//      | "for" "(" expr? ";" expr? ";" expr? ")" stmt
//      | "{" stmt* "}"
//      | "typedef" type-specifier declarator type-suffix ";"
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

		if consume(";") == nil {
			node.Init = readExprStmt()
			expect(";")
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

	} else if t := consume("typedef"); t != nil {

		ty := typeSpecifier()
		var name string
		ty = declarator(ty, &name)
		ty = typeSuffix(ty)
		expect(";")

		pushScope(name).TyDef = ty
		return &Node{Kind: ND_NULL, Tok: t}

	} else {

		if isTypename() {
			return declaration()
		}

		node = readExprStmt()
		expect(";")
	}

	return node
}

// expr       = assign
func expr() *Node {
	// printCurTok()
	// printCurFunc()
	return assign()
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

// mul = unary ("*" unary | "/" unary)*
func mul() *Node {
	// printCurTok()
	// printCurFunc()
	node := unary()

	for {
		if t := consume("*"); t != nil {
			node = newNode(ND_MUL, node, unary(), t)
		} else if consume("/") != nil {
			node = newNode(ND_DIV, node, unary(), t)
		} else {
			return node
		}
	}
}

// unary = ("+" | "-" | "*" | "&")? unary
//       | "sizeof" unary
//       | postfix
func unary() *Node {
	// printCurTok()
	// printCurFunc()
	if t := consumeSizeof(); t != nil {
		return newUnary(ND_SIZEOF, unary(), t)
	}
	if t := consume("+"); t != nil {
		return unary()
	}
	if t := consume("-"); t != nil {
		return newNode(ND_SUB, newNodeNum(0, t), unary(), t)
	}
	if t := consume("&"); t != nil {
		return newUnary(ND_ADDR, unary(), t)
	}
	if t := consume("*"); t != nil {
		return newUnary(ND_DEREF, unary(), t)
	}
	return postfix()
}

// postfix = primary ("[" expr "]" | "." ident | "->" ident)*
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
		if sc != nil && sc.Var != nil {
			return newVar(sc.Var, tok)
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
