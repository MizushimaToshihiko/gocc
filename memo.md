
### チラ裏的なもの
#### 【後回し】※順不同
 - Allow for-loops to define local variables  
   => 型推論が終わってから, for-clauseのinitではShortVarDeclしか記載できない為  
 - Add flexible array member  
   => とりあえず今のところはsiliceを長さ0の配列としている。後でsliceを定義してparse出来るようにする
 - tokenizer変更
   公式のtoken packageに合わせてVARトークンやFUNCトークンを作り、FUNCトークンの子としてFUNC内のstatementをtokenizeする <= parseしやすくなるかもしれないので
 - parseの順番を変える
   現状では関数の後に宣言されたグルーバル変数を参照するとparserでエラーになるので、var, type, func(変数スコープ登録のみ)のparseの後にfunc内部のparseを行うように変更
 - 型推論  
   "var x = expr"とか、"x := expr"とか
 - initializerでの型名省略
   "var x [2]T = [2]T{T{1,2},T{3,4}}"を"var x [2]T = [2]T{{1,2},{3,4}}"で可とする
 - RangeClause  
   "for x := range X"みたいなもの
 - 配列の宣言で"[...]int{1,2,3}"みたいなもの
 - 定数宣言
 - map型
 - slice
 - Typeに型の名前を持たせて、pointer型とstring型を外面上は別物にする
 - 関数戻り値の型チェック(type checkingというのかな)
 - goroutineは無理かな？
 - package
    - main package
 - import
 - built-in functions
    - new
    - make
    - len
    - println
    - cap
    - append(slice)
    - copy(slice)
    - panic
    - recover
 - "switch 変数 {"とか"switch 型 {"とか
 - blank identifiers => "_"
 - bool型でtrueやfalseを使用できるように
 - float
 - complex(複素数) いる?
 - rune(int32のエイリアス)
 - rune literal => tokenizerのchar literalを変更する?
 - method set(メソッド集合)
 - 構造体埋め込みでメソッド集合も埋め込む
 - interface
 - クロージャ
 - gc
 

#### 【VarSpecの追加について】
 - EBNF:VarSpec = ident-list (type-preffix type-specifier [ "=" expr-list ] | "=" expr-list)
 ```Go
package main

func MerryXMas() {
	var x, y int = 1, 2
}
```
 - 上記のコードのASTを出力すると、下記のようになる。
 ```
     0  *ast.FuncDecl {
     1  .  Name: *ast.Ident {
     2  .  .  NamePos: testdata/test01.go:3:6
     3  .  .  Name: "MerryXMas"
     4  .  .  Obj: *ast.Object {
     5  .  .  .  Kind: func
     6  .  .  .  Name: "MerryXMas"
     7  .  .  .  Decl: *(obj @ 0)
     8  .  .  }
     9  .  }
    10  .  Type: *ast.FuncType {
    11  .  .  Func: testdata/test01.go:3:1
    12  .  .  Params: *ast.FieldList {
    13  .  .  .  Opening: testdata/test01.go:3:15
    14  .  .  .  Closing: testdata/test01.go:3:16
    15  .  .  }
    16  .  }
    17  .  Body: *ast.BlockStmt {
    18  .  .  Lbrace: testdata/test01.go:3:18
    19  .  .  List: []ast.Stmt (len = 1) {
    20  .  .  .  0: *ast.DeclStmt {
    21  .  .  .  .  Decl: *ast.GenDecl {
    22  .  .  .  .  .  TokPos: testdata/test01.go:4:2
    23  .  .  .  .  .  Tok: var
    24  .  .  .  .  .  Lparen: -
    25  .  .  .  .  .  Specs: []ast.Spec (len = 1) {
    26  .  .  .  .  .  .  0: *ast.ValueSpec {
    27  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 2) {
    28  .  .  .  .  .  .  .  .  0: *ast.Ident {
    29  .  .  .  .  .  .  .  .  .  NamePos: testdata/test01.go:4:6
    30  .  .  .  .  .  .  .  .  .  Name: "x"
    31  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    32  .  .  .  .  .  .  .  .  .  .  Kind: var
    33  .  .  .  .  .  .  .  .  .  .  Name: "x"
    34  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 26)
    35  .  .  .  .  .  .  .  .  .  .  Data: 0
    36  .  .  .  .  .  .  .  .  .  }
    37  .  .  .  .  .  .  .  .  }
    38  .  .  .  .  .  .  .  .  1: *ast.Ident {
    39  .  .  .  .  .  .  .  .  .  NamePos: testdata/test01.go:4:9
    40  .  .  .  .  .  .  .  .  .  Name: "y"
    41  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    42  .  .  .  .  .  .  .  .  .  .  Kind: var
    43  .  .  .  .  .  .  .  .  .  .  Name: "y"
    44  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 26)
    45  .  .  .  .  .  .  .  .  .  .  Data: 0
    46  .  .  .  .  .  .  .  .  .  }
    47  .  .  .  .  .  .  .  .  }
    48  .  .  .  .  .  .  .  }
    49  .  .  .  .  .  .  .  Type: *ast.Ident {
    50  .  .  .  .  .  .  .  .  NamePos: testdata/test01.go:4:11
    51  .  .  .  .  .  .  .  .  Name: "int"
    52  .  .  .  .  .  .  .  }
    53  .  .  .  .  .  .  .  Values: []ast.Expr (len = 2) {
    54  .  .  .  .  .  .  .  .  0: *ast.BasicLit {
    55  .  .  .  .  .  .  .  .  .  ValuePos: testdata/test01.go:4:17
    56  .  .  .  .  .  .  .  .  .  Kind: INT
    57  .  .  .  .  .  .  .  .  .  Value: "1"
    58  .  .  .  .  .  .  .  .  }
    59  .  .  .  .  .  .  .  .  1: *ast.BasicLit {
    60  .  .  .  .  .  .  .  .  .  ValuePos: testdata/test01.go:4:20
    61  .  .  .  .  .  .  .  .  .  Kind: INT
    62  .  .  .  .  .  .  .  .  .  Value: "2"
    63  .  .  .  .  .  .  .  .  }
    64  .  .  .  .  .  .  .  }
    65  .  .  .  .  .  .  }
    66  .  .  .  .  .  }
    67  .  .  .  .  .  Rparen: -
    68  .  .  .  .  }
    69  .  .  .  }
    70  .  .  }
    71  .  .  Rbrace: testdata/test01.go:5:1
    72  .  }
    73  }
  ```
  
 - FunctionBody -> BlockStmt -> DeclStmtのSpecメンバ -> x,yのValueSpecがスライスとして登録されている
 - このコンパイラでは、var x,y int = 1,2をparse.goの中でvar x int = 1; var y int = 2;としてfunction()内のstmtの後に繋げる -> declaration()とは別にvarspec()を作り、nodeをつなげたものをfunction()又はstmt()に返す?

#### 配列変数から配列変数への代入
- 現時点では配列から配列への代入ができない（not a lvalueエラーを出してしまう)
- string変数からstring変数への代入も同様にできない。stringをbase typeがbyteの配列にしているため。⇒string型をarrayType()からpointerTo()にしたら通った。

#### 型が違うので代入できないエラーを返す関数の書きかけ
 - typedefの名前の取得が現時点でできない為
 ```Go
 package main

import (
	"errors"
	"strconv"
)

// cannotAssignArr: ty1は代入される方の変数の型、ty2は代入する方の変数の型
func cannotAssignArrErr(ty1, ty2 *Type, tok *Token) error {
	var retTy1, retTy2 string = typeStr(ty1, tok), typeStr(ty2, tok)

	return errors.New(errorTok(tok,
		"cannot use %s {...} (type %s) as Type %s in array literal",
		retTy2,
		retTy2,
		retTy1))
}

func typeStr(ty *Type, tok *Token) string {
	var retTy string
	// make retTy1
	switch ty.Kind {
	case TY_ARRAY:
		retTy += "[" + strconv.Itoa(ty.ArrSz) + "]"
		for ty.Base != nil {
			retTy += typeStr(ty.Base, tok)
		}
	case TY_STRUCT:
		if findTyDef(tok) != nil {
			retTy += tok.Str
		}
	}
	return retTy
}
``` 