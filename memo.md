
### メモ
#### 【後回し】※順不同
 - [Add stage2 build](https://github.com/rui314/chibicc/commit/5d15431df1abab3a5cf596fabe0a77c030a10791)
 - バッククオート
 - parseの順番を変える
   現状では関数の後に宣言されたグルーバル変数を参照するとparserでエラーになるので、var(含初期化), type(含初期化)
 - RangeClause  
   "for x := range X"みたいなもの
 - const宣言
 - 文字列の足し算
 - *(*type-name)(unary)
 - (*var-name)[n]
 - map型
 - slice
 - Add flexible array member  
   => とりあえず今のところはsiliceを長さ0の配列としている。後でsliceを定義してparse出来るようにする
 - 配列の宣言で"[...]int{1,2,3}"みたいなもの(slice追加後)
 - 関数戻り値の型チェック(type checking)
 - goroutine
 - package
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
 - switch ident.(type) {
 - case case1,case2:
 - blank identifiers : "_"
 - bool型でtrueやfalseを使用できるように
 - float
 - complex(複素数)
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
- ~~現時点では配列変数から配列変数への代入ができない（not a lvalueエラーを出してしまう)~~  
  Type構造体の要素にInitを追加し、Obj構造体の.Ty.InitにInitializerを保存し、代入時に右辺のObjから左辺のObjに.Tyを丸ごとコピーすることで実装済み、copyType()を使った方が良いかも
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

#### 2022/01/19
- 引数付の関数定義?呼び出し?時にsegmentation faultが出る
- testdata/commonにc言語で定義するとsegmentation faultにならない
- 2022/01/20 PrologueとEpilogueのレジスタ名が間違っていたのが原因

#### 2022/02/09
 - 下記の場合、nil pointerエラーになる。原因はyのtypeが*intになっているためと思われる。
```Go
var x = [2]int{1, 2}
var y = &x
assert(2, (*y)[1], "(*y)[1]")
```
 - 下記二つのケースではエラーにならない。
```Go
var x = [2]int{1, 2}
var y *[2]int = &x
assert(2, (*y)[1], "(*y)[1]")
```
```Go
var x = [2]int{1, 2}
var y = &x
assert(2, y[1], "y[1]")
```
 - yのtypeが型推論でも *[2]int になるように、initializerを変更する必要がある。
 - 2022/02/10 struct配列へのポインタアドレス代入が可能になるよう変更
 - 上記例で(*y)[0]は動いたが代わりにy[0]がおかしくなった。でも一旦これで。

2022/02/16 Intelマニュアルの複数パラメータ受渡部分の抜粋(DQNEOさんが紹介していた部分)  

From : https://www.intel.com/content/www/us/en/developer/articles/technical/intel-sdm.html#combined  
Copyright © 1997-2021, Intel Corporation. All Rights Reserved.  

 - 6.4.3 Parameter Passing  
Parameters can be passed between procedures in any of three ways: through general-purpose registers, in an argument list, or on the stack.  
 - 6.4.3.1 Passing Parameters Through the General-Purpose Registers  
The processor does not save the state of the general-purpose registers on procedure calls. A calling procedure can thus pass up to six parameters to the called procedure by copying the parameters into any of these registers (except the ESP and EBP registers) prior to executing the CALL instruction. The called procedure can likewise pass parameters back to the calling procedure through general-purpose registers.
 - 6.4.3.2 Passing Parameters on the Stack  
To pass a large number of parameters to the called procedure, the parameters can be placed on the stack, in the stack frame for the calling procedure. Here, it is useful to use the stack-frame base pointer (in the EBP register) to make a frame boundary for easy access to the parameters.  
The stack can also be used to pass parameters back from the called procedure to the calling procedure.
 - 6.4.3.3 Passing Parameters in an Argument List  
An alternate method of passing a larger number of parameters (or a data structure) to the called procedure is to place the parameters in an argument list in one of the data segments in memory. A pointer to the argument list can then be passed to the called procedure through a general-purpose register or the stack. Parameters can also be passed back to the calling procedure in this same manner.  
 - 6.4.3パラメータの受け渡し  
パラメータは、汎用レジスタ、引数リスト、またはスタックの3つの方法のいずれかでプロシージャ間で渡すことができます。
 - 6.4.3.1汎用レジスタを介したパラメータの受け渡し  
プロセッサは、プロシージャ呼び出しで汎用レジスタの状態を保存しません。したがって、呼び出し元のプロシージャは、CALL命令を実行する前に、これらのレジスタ（ESPおよびEBPレジスタを除く）のいずれかにパラメータをコピーすることにより、呼び出されたプロシージャに最大6つのパラメータを渡すことができます。呼び出されたプロシージャも同様に、汎用レジスタを介してパラメータを呼び出し元のプロシージャに戻すことができます。
 - 6.4.3.2スタックでのパラメータの受け渡し  
呼び出されたプロシージャに多数のパラメータを渡すために、パラメータをスタックの呼び出し元のプロシージャのスタックフレームに配置できます。ここでは、（EBPレジスタ内の）スタックフレームベースポインタを使用して、パラメータに簡単にアクセスできるようにフレーム境界を作成すると便利です。
スタックを使用して、呼び出されたプロシージャから呼び出し元のプロシージャにパラメータを戻すこともできます。
 - 6.4.3.3引数リストでのパラメータの受け渡し  
呼び出されたプロシージャに多数のパラメータ（またはデータ構造）を渡す別の方法は、メモリ内のデータセグメントの1つにある引数リストにパラメータを配置することです。次に、引数リストへのポインタを、汎用レジスタまたはスタックを介して呼び出されたプロシージャに渡すことができます。これと同じ方法で、パラメータを呼び出し元のプロシージャに戻すこともできます。

#### 2022/02/17 floating-pointを扱う
 - chibiccのcodegen.cのgen_expr関数では、  
 　1. union { float f32; double f64; uint32_t u32; uint64_t u64; } u;を定義  
 　2. tokenのfvalに入っている小数の値をu.f32(floatの場合),又はu.f64(doubleの場合)に格納  
 　3. 同じ大きさのuintのメンバ(u.f32->u.u32,u.f64->u.u64)の値をprintfで%uとして整数で取り出し、raxに入れている  
 - 上記3で取り出した値は、u.f32、u.f64の値をビットで表現した時の、仮数部分を10進数にしたものと同じ？  
   例：0.1（double、10進数）の場合、下記コードで調べるとu.u64は4591870180066957722（10進数）、これを2進数にすると0b10011001100110011001100110011010になる。

```c
#include <stdio.h>
#include <stdint.h>
#include <limits.h>

void printb(unsigned int v) {
  unsigned int mask = (int)1 << (sizeof(v) * CHAR_BIT - 1);
  do putchar(mask & v ? '1' : '0');
  while (mask >>= 1);
}

void putb(unsigned int v) {
  putchar('0'), putchar('b'), printb(v), putchar('\n');
}

int main(void){
    union { float f32; double f64; uint32_t u32; uint64_t u64; } u;
    u.f64=0.1;
    printf("u.f32: %p: %f\n", &u.f32, u.f32);  // u.f32: 0x7fff2119d0d0: -0.000000
    printf("u.f64: %p: %f\n", &u.f64, u.f64);  // u.f64: 0x7fff2119d0d0: 0.100000
    printf("u.u32: %p: %u\n", &u.u32, u.u32);  // u.u32: 0x7fff2119d0d0: 2576980378
    printf("u.u64: %p: %lu\n", &u.u64, u.u64); // u.u64: 0x7fff2119d0d0: 4591870180066957722
    printf("       ");
    putb(u.u64);                               // 0b10011001100110011001100110011010
}
   ```
   また0.1を[このサイト](https://tools.m-bsys.com/calculators/ieee754.php)でIEEE754内部表現にすると下図のようになる。
   ![画像](img/20220217.png)

 - このコンパイラで再現するにはどうするか？  
   アセンブリ言語に小数は入れられない？  
   小数の仮数部分を取り出して10進数の数字を得るには?  
   例: 0.1の場合 -> 4591870180066957722を何らかの方法で算出する
   - [Goコードでの内部表現取り出し](https://go.dev/play/p/k3rD8Exk3DX)
```Go
package main

import (
	"fmt"
	"strconv"
)

func main() {
	x := 0.1
	fmt.Printf("%b\n", x)
	fmt.Printf("7205759403792794:             %b\n", 7205759403792794)
	fmt.Printf("4591870180066957722: %b\n", 4591870180066957722)
	fmt.Printf("%d\n", 0b1001100110011001100110011001100110011001100110011010)
	f := strconv.FormatFloat(7205759403792794, 'f', -1, 64)
	fmt.Println(f)
}
```
output:
```
7205759403792794p-56
7205759403792794:             11001100110011001100110011001100110011001100110011010
4591870180066957722: 11111110111001100110011001100110011001100110011001100110011010
2702159776422298
7205759403792794

Program exited.
```
 - 下記でいいみたいです。参照:https://pkg.go.dev/unsafe#Pointer , https://qiita.com/nia_tn1012/items/d26f0fc993895a09b30b#23-%E3%83%9D%E3%82%A4%E3%83%B3%E3%82%BF%E3%81%AE%E3%82%AD%E3%83%A3%E3%82%B9%E3%83%88%E5%A4%89%E6%8F%9B%E3%82%92%E5%88%A9%E7%94%A8%E3%81%97%E3%81%9F%E6%96%B9%E6%B3%95-c%E8%A8%80%E8%AA%9Ecc-
```Go
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	x := 0.1
	s := *(*uint64)(unsafe.Pointer(&x)) // ここ
	fmt.Println(s)
}
```

2022/02/18 左辺ポインタのキャスト変換&書き込みの例[URL](https://go.dev/play/p/FlES6L9lUOU)
```Go
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var x [2]int64

	var idx int = 1
	*(*float32)(unsafe.Pointer(uintptr(unsafe.Pointer(&x[0])) + idx*unsafe.Sizeof(x[0]))) = 11
	*(*float32)(unsafe.Pointer(&x[idx])) = 11 // ↑の短縮版
	// *(*float32)(&x[idx]) = 11 => c言語風の書き方はできない
	fmt.Printf("%d\n", x)
	s = *(*int32)(unsafe.Pointer(&x))
	fmt.Println(s)
}
```
