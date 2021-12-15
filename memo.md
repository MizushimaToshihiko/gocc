
#### 【後回し】
 - Allow for-loops to define local  variables  
   => 型推論が終わってから, for-clauseのinitではShortVarDeclしか記載できない為

#### 【VarSpecの追加について】
 - EBNF:VarSpec = ident-list (type-preffix type-specifier [ "=" expr-list ] | "=" expr-list)
 ```Go:test01.go
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