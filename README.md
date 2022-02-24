### The toy compiler for Go
- It's based on rui314's [chibicc](https://github.com/rui314/chibicc/tree/reference) and https://www.sigbus.info/compilerbook . 
- This is for my own learning.

### EBNF
```ebnf
Production  = name "=" [ Expression ] "." .
Expression  = Alternative { "|" Alternative } .
Alternative = Term { Term } .
Term        = name | token [ "…" token ] | Group | Option | Repetition .
Group       = "(" Expression ")" .
Option      = "[" Expression "]" .
Repetition  = "{" Expression "}" .
```
From : [golang.org/x/exp](https://golang.org/ref/spec#Notation)
