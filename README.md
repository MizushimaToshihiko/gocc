### 9cc in Golang
- It's based on rui314's [chibicc](https://github.com/rui314/chibicc/tree/reference) and https://www.sigbus.info/compilerbook .

### EBNF
```ebnf
program     = (global-var | function*)
basetype    = ("int" | "char") "*"*
param       = basetype ident
params      = param ("," param)*
function    = basetype ident "(" params? ")" "{" stmt* "}"
global-var  = basetype ident ("[" num "]")* ";"
declaration = basetype ident ("[" num "]")* ("=" expr) ";"
stmt        = "return" expr ";"
            | "if" "(" expr ")" stmt ("else" stmt)?
            | "while" "(" expr ")" stmt
            | "for" "(" expr? ";" expr? ";" expr? ")" stmt
            | "{" stmt* "}"
            | declaration
            | expr ";"
expr        = assign
assign      = equality ("=" assign)?
equality    = relational ("==" relational | "!=" relational)*
relational  = add ("<" add | "<=" add | ">" add | ">=" add)*
add         = mul ("+" mul | "-" mul)*
mul         = unary ("*" unary | "/" unary)*
unary       = ("+" | "-" | "*" | "&")? unary
            | "sizeof" unary
            | postfix
postfix    = primary ("[" expr "]")*
stmt-expr  = "(" "{" stmt stmt* "}" ")"
func-args  = "(" (assign("," assign)*)? ")"
primary    = "(" "{" stmt-expr-tail
           | ident func-args?
           | "(" expr ")"
           | num
           | str
```