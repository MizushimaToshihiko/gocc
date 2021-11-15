package main

import (
	"fmt"
	"runtime"
)

// for printTokens function, the pointer of the head token
// stored in 'headTok'.
var headTok *Token

func printTokens() {
	fmt.Print("# Tokens:\n")
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

// walk AST in in-order
func walkInOrder(node *Node) {
	fmt.Print("# Nodes in-order: ")
	inOrder(node)
	fmt.Println()
}

func inOrder(node *Node) {
	if node == nil {
		return
	}
	inOrder(node.Lhs)
	switch node.Kind {
	case ND_NUM:
		if isLeaf(node) {
			fmt.Printf(" '%s': %d: leaf ", "ND_NUM", node.Val)
		} else {
			fmt.Printf(" '%s': %d: ", "ND_NUM", node.Val)
		}
	case ND_ADD:
		if isLeaf(node) {
			fmt.Printf(" '%s': leaf ", "ND_ADD: +")
		} else {
			fmt.Printf(" '%s': ", "ND_ADD: +")
		}
	case ND_SUB:
		if isLeaf(node) {
			fmt.Printf(" '%s': leaf ", "ND_SUB: -")
		} else {
			fmt.Printf(" '%s': ", "ND_SUB: -")
		}
	case ND_MUL:
		if isLeaf(node) {
			fmt.Printf(" '%s': leaf ", "ND_MUL: *")
		} else {
			fmt.Printf(" '%s': ", "ND_MUL: *")
		}
	case ND_DIV:
		if isLeaf(node) {
			fmt.Printf(" '%s': leaf ", "ND_DIV: /")
		} else {
			fmt.Printf(" '%s': ", "ND_DIV: /")
		}
	case ND_EQ:
		if isLeaf(node) {
			fmt.Printf(" '%s': leaf ", "ND_EQ: ==")
		} else {
			fmt.Printf(" '%s': ", "ND_EQ: ==")
		}
	case ND_NE:
		if isLeaf(node) {
			fmt.Printf(" '%s': leaf ", "ND_NE: !=")
		} else {
			fmt.Printf(" '%s': ", "ND_NE: !=")
		}
	case ND_LT:
		if isLeaf(node) {
			fmt.Printf(" '%s': leaf ", "ND_LT: <")
		} else {
			fmt.Printf(" '%s': ", "ND_LT: <")
		}
	case ND_LE:
		if isLeaf(node) {
			fmt.Printf(" '%s': leaf ", "ND_LE: <=")
		} else {
			fmt.Printf(" '%s': ", "ND_LE: <=")
		}
	}
	inOrder(node.Rhs)
}

// walk AST in pre-order
func walkPreOrder(node *Node) {
	fmt.Print("# Nodes pre-order: ")
	preOrder(node)
	fmt.Println()
}

func preOrder(node *Node) {
	if node == nil {
		return
	}
	switch node.Kind {
	case ND_NUM:
		if isLeaf(node) {
			fmt.Printf(" '%s': %d: leaf ", "ND_NUM", node.Val)
		} else {
			fmt.Printf(" '%s': %d: ", "ND_NUM", node.Val)
		}
	case ND_ADD:
		if isLeaf(node) {
			fmt.Printf(" '%s': leaf ", "ND_ADD: +")
		} else {
			fmt.Printf(" '%s': ", "ND_ADD: +")
		}
	case ND_SUB:
		if isLeaf(node) {
			fmt.Printf(" '%s': leaf ", "ND_SUB: -")
		} else {
			fmt.Printf(" '%s': ", "ND_SUB: -")
		}
	case ND_MUL:
		if isLeaf(node) {
			fmt.Printf(" '%s': leaf ", "ND_MUL: *")
		} else {
			fmt.Printf(" '%s': ", "ND_MUL: *")
		}
	case ND_DIV:
		if isLeaf(node) {
			fmt.Printf(" '%s': leaf ", "ND_DIV: /")
		} else {
			fmt.Printf(" '%s': ", "ND_DIV: /")
		}
	case ND_EQ:
		if isLeaf(node) {
			fmt.Printf(" '%s': leaf ", "ND_EQ: ==")
		} else {
			fmt.Printf(" '%s': ", "ND_EQ: ==")
		}
	case ND_NE:
		if isLeaf(node) {
			fmt.Printf(" '%s': leaf ", "ND_NE: !=")
		} else {
			fmt.Printf(" '%s': ", "ND_NE: !=")
		}
	case ND_LT:
		if isLeaf(node) {
			fmt.Printf(" '%s': leaf ", "ND_LT: <")
		} else {
			fmt.Printf(" '%s': ", "ND_LT: <")
		}
	case ND_LE:
		if isLeaf(node) {
			fmt.Printf(" '%s': leaf ", "ND_LE: <=")
		} else {
			fmt.Printf(" '%s': ", "ND_LE: <=")
		}
	}
	preOrder(node.Lhs)
	preOrder(node.Rhs)
}

func isLeaf(node *Node) bool {
	return node.Lhs == nil && node.Rhs == nil
}
