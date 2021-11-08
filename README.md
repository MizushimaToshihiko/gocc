### 9cc in Golang
- It's based on rui314's [chibicc](https://github.com/rui314/chibicc/tree/reference) and https://www.sigbus.info/compilerbook .

### EBNF
```ebnf
program        = (global-var | function*)
type-specifier = builtin-type | struct-decl | typedef-name
builtin-type   = "char" | "short" | "int" | "long"
declarator     = "*" ("(" declarator ")") | ident) type-suffix
type-suffix    = ("[" num "]" type-suffix)?
struct-decl    = "struct" ident
               | "struct" ident? "{" struct-member "}"
struct-member  = type-specifier declarator type-suffix ";"
param          = type-specifier declarator type-suffix
params         = param ("," param)*
function       = type-specifier declarator "(" params? ")" "{" stmt* "}"
global-var     = type-specifier declarator type-suffix ";"
declaration    = type-specifier declarator type-suffix ("=" expr)? ";"
               | type-specifier ";"
stmt           = "return" expr ";"
               | "if" "(" expr ")" stmt ("else" stmt)?
               | "while" "(" expr ")" stmt
               | "for" "(" expr? ";" expr? ";" expr? ")" stmt
               | "{" stmt* "}"
               | "typedef" type-specifier declarator type-suffix ";"
               | declaration
               | expr ";"
expr           = assign
assign         = equality ("=" assign)?
equality       = relational ("==" relational | "!=" relational)*
relational     = add ("<" add | "<=" add | ">" add | ">=" add)*
add            = mul ("+" mul | "-" mul)*
mul            = unary ("*" unary | "/" unary)*
unary          = ("+" | "-" | "*" | "&")? unary
               | "sizeof" unary
               | postfix
postfix        = primary ("[" expr "]" | "." ident | "->" ident)*
stmt-expr      = "(" "{" stmt stmt* "}" ")"
func-args      = "(" (assign("," assign)*)? ")"
primary        = "(" "{" stmt-expr-tail
               | ident func-args?
               | "(" expr ")"
               | num
               | str
```