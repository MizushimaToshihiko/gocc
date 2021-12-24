package main

import (
	"fmt"
	"runtime"
)

func assert(b bool, m string) {
	if !b {
		panic(m)
	}
}

// for printTokens function, the pointer of the head token
// stored in 'headTok'.
var headTok *Token

func printTokens() {
	fmt.Print("# Tokens:\n# ")
	tok := headTok.Next
	// var kind string
	for tok.Next != nil {
		// 	switch tok.Kind {
		// 	case TK_IDENT:
		// 		kind = "IDENT"
		// 	case TK_NUM:
		// 		kind = "NUM"
		// 	case TK_RESERVED:
		// 		kind = "RESERVED"
		// 	case TK_SIZEOF:
		// 		kind = "SIZEOF"
		// 	case TK_STR:
		// 		kind = "STR"
		// 	default:
		// 		log.Fatal("unknown token kind")
		// 	}
		// 	fmt.Printf(" %s: Str:\"%s\" :%d Val:%d\n", kind, tok.Str, tok.Len, tok.Val)
		fmt.Printf(" '%s'", tok.Str)
		tok = tok.Next
	}

	if tok.Kind == TK_EOF {
		fmt.Print(" EOF ")
	}

	fmt.Print("\n\n")
}

// func printCurTokInit() {
// 	fmt.Print("# Current Token: ")
// }

func printCurTok() {
	fmt.Printf(" %d:'%s' \n", token.Kind, token.Str)
}

func printCalledFunc() {
	pc, _, line, _ := runtime.Caller(2)
	fn := runtime.FuncForPC(pc)
	fmt.Printf(" %s %d\n", fn.Name(), line)
	pc, _, line, _ = runtime.Caller(1)
	fn = runtime.FuncForPC(pc)
	fmt.Printf(" %s %d\n", fn.Name(), line)
}

var ND = map[NodeKind]string{
	0:  "ND_ADD",       // 0: +
	1:  "ND_SUB",       // 1: -
	2:  "ND_MUL",       // 2: *
	3:  "ND_DIV",       // 3: /
	4:  "ND_EQ",        // 4: ==
	5:  "ND_NE",        // 5: !=
	6:  "ND_LT",        // 6: <
	7:  "ND_LE",        // 7: <=
	8:  "ND_ASSIGN",    // 8: =
	9:  "ND_VAR",       // 9: local variables
	10: "ND_NUM",       // 10: integer
	11: "ND_RETURN",    // 11: 'return'
	12: "ND_IF",        // 12: "if"
	13: "ND_WHILE",     // 13: "while"
	14: "ND_FOR",       // 14: "for"
	15: "ND_BLOCK",     // 15: {...}
	16: "ND_FUNCALL",   // 16: function call
	17: "ND_ADDR",      // 17: unary &
	18: "ND_DEREF",     // 18: unary *
	19: "ND_EXPR_STMT", // 19: expression statement
	20: "ND_NULL",      // 20: empty statement
}

// walk AST in in-order
func walkInOrder(node *Node) {
	fmt.Print("# Nodes in-order: ")
	for n := node; n != nil; n = n.Next {
		inOrder(node)
		fmt.Println()
	}
}

func inOrder(node *Node) {
	if node == nil {
		return
	}
	inOrder(node.Lhs)
	switch node.Kind {
	case ND_NUM:
		if isLeaf(node) {
			fmt.Printf(" '%s': %d: leaf ", ND[node.Kind], node.Val)
		} else {
			fmt.Printf(" '%s': %d: ", ND[node.Kind], node.Val)
		}
	case ND_ADD:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "+")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "+")
		}
	case ND_SUB:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "-")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "-")
		}
	case ND_MUL:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "*")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "*")
		}
	case ND_DIV:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "/")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "/")
		}
	case ND_EQ:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "==")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "==")
		}
	case ND_NE:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "!=")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "!=")
		}
	case ND_LT:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "<")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "<")
		}
	case ND_LE:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "<=")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "<=")
		}
	case ND_ASSIGN:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "=")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "=")
		}
	case ND_VAR:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], node.Obj.Name)
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], node.Obj.Name)
		}
	case ND_RETURN:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "return")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "return")
		}
	case ND_IF:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "if")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "if")
		}
	case ND_WHILE:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "for-stmt")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "for-stmt")
		}
	case ND_FOR:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "for-clause")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "for-clause")
		}
	case ND_BLOCK:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "{}")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "{}")
		}
	case ND_FUNCALL:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], node.FuncName)
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], node.FuncName)
		}
	case ND_ADDR:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "&")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "&")
		}
	case ND_DEREF:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "*")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "*")
		}
	case ND_EXPR_STMT:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "ExprStmt")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "ExprStmt")
		}
	case ND_NULL:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "NULL")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "NULL")
		}
	}
	inOrder(node.Rhs)
}

// walk AST in pre-order
func walkPreOrder(node *Node) {
	fmt.Print("# Nodes pre-order: ")
	for n := node; n != nil; n = n.Next {
		preOrder(node)
		fmt.Println()
	}
}

func preOrder(node *Node) {
	if node == nil {
		return
	}
	switch node.Kind {
	case ND_NUM:
		if isLeaf(node) {
			fmt.Printf(" '%s': %d: leaf ", ND[node.Kind], node.Val)
		} else {
			fmt.Printf(" '%s': %d: ", ND[node.Kind], node.Val)
		}
	case ND_ADD:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "+")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "+")
		}
	case ND_SUB:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "-")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "-")
		}
	case ND_MUL:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "*")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "*")
		}
	case ND_DIV:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "/")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "/")
		}
	case ND_EQ:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "==")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "==")
		}
	case ND_NE:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "!=")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "!=")
		}
	case ND_LT:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "<")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "<")
		}
	case ND_LE:
		if isLeaf(node) {
			fmt.Printf(" '%s: %s': leaf ", ND[node.Kind], "<=")
		} else {
			fmt.Printf(" '%s': %s ", ND[node.Kind], "<=")
		}
	}
	preOrder(node.Lhs)
	preOrder(node.Rhs)
}

func isLeaf(node *Node) bool {
	return node.Lhs == nil && node.Rhs == nil
}
