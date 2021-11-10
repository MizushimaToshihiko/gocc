### 9cc in Golang
- It's based on rui314's [chibicc](https://github.com/rui314/chibicc/tree/reference) and https://www.sigbus.info/compilerbook . 
- This is for my own learning.

### EBNF
```ebnf
program        = (global-var | function*)
type-specifier = builtin-type | struct-decl | typedef-name
 node that "typedef" and "static" can appear anywhere in a type-specifier
builtin-type   = "void"
               | "_Bool"
               | "char"
               | "short" | "short" "int" | "int" "short"
               | "int"
               | "long" | "long" "int" | "int" "long"
declarator     = "*" ("(" declarator ")") | ident) type-suffix
type-suffix    = ("[" const-expr? "]" type-suffix)?
struct-decl    = "struct" ident? ("{" struct-member "}")?
struct-member  = type-specifier declarator type-suffix ";"
enum-specifier = "enum" ident
               | "enum" ident? "{" enum-list? "}"
enum-list      = enum-elem ("," enum-elem)* ","?
param          = type-specifier declarator type-suffix
params         = param ("," param)*
function       = type-specifier declarator "(" params? ")" ("{" stmt* "}" | ";")
global-var     = type-specifier declarator type-suffix ";"
lvar-initializer = assign
                 | "{" lvar-initializer ("," lvar-initializer)* ","? "}"
declaration    = type-specifier declarator type-suffix ("=" lvar-initializer)? ";"
               | type-specifier ";"
stmt           = "return" expr ";"
               | "if" "(" expr ")" stmt ("else" stmt)?
               | "switch" "(" expr ")" stmt
               | "case" const-expr ":" stmt
               | "default" ":" stmt
               | "while" "(" expr ")" stmt
               | "for" "(" expr? ";" expr? ";" expr? ")" stmt
               | "{" stmt* "}"
               | "break" ";"
               | continue" ";"
               | "goto" ident ";"
               | ident ":" stmt
               | declaration
               | expr ";"
expr           = assign
assign         = conditional (assign-op assign)?
assign-op      = "=" | "+=" | "-=" | "*=" | "/=" | "<<=" | ">>="
conditional    = logor ("?" expr ":" conditional)?
logor          = logand ("||" logand)*
logand         = bitor ("&&" bitor)*
bitor          = bitxor ("|" bitxor)*
bitxor         = bitand ("^" bitand)*
bitand         = equality ("&" equality)*
equality       = relational ("==" relational | "!=" relational)*
relational     = shift ("<" shift | "<=" shift | ">" shift | ">=" shift)*
shift          = add ("<<" add | ">>" add)*
add            = mul ("+" mul | "-" mul)*
mul            = unary ("*" unary | "/" unary)*
unary          = ("+" | "-" | "*" | "&" | "!")? unary
               | ("++" | "--") unary
               | "sizeof" "(" type-name ")"
               | "sizeof" unary
               | postfix
postfix        = primary ("[" expr "]" | "." ident | "->" ident | "++" | "--")*
stmt-expr      = "(" "{" stmt stmt* "}" ")"
func-args      = "(" (assign("," assign)*)? ")"
primary        = "(" "{" stmt-expr-tail
               | ident func-args?
               | "(" expr ")"
               | num
               | str
```