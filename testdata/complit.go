package test_complit

func assert(want int, act int, code string)
func println(format string)

type Tree struct {
	val int
	lhs *Tree
	rhs *Tree
}

var tree1 *Tree = &Tree{
	1,
	&Tree{
		2,
		&Tree{3, 0, 0},
		&Tree{4, 0, 0},
	},
	0,
}

var tree2 = &Tree{
	1,
	&Tree{
		2,
		&Tree{3, 0, 0},
		&Tree{4, 0, 0},
	},
	0,
}

func main() {

	println("\ntype Tree struct {val int;lhs *Tree;rhs *Tree};")
	println("var tree1 *Tree=&Tree{1,&Tree{2,&Tree{3,0,0},&Tree{4,0,0},},0,};\n")
	assert(1, tree1.val, "tree1.val")
	assert(2, tree1.lhs.val, "tree1.lhs.val")
	assert(3, tree1.lhs.lhs.val, "tree1.lhs.lhs.val")
	assert(4, tree1.lhs.rhs.val, "tree1.lhs.rhs.val")

	println("\ntype Tree struct {val int;lhs *Tree;rhs *Tree};\nvar tree2 *Tree=&Tree{1,&Tree{2,&Tree{3,0,0},&Tree{4,0,0},},0,};\n")
	assert(1, tree2.val, "tree2.val")
	assert(2, tree2.lhs.val, "tree2.lhs.val")
	assert(3, tree2.lhs.lhs.val, "tree2.lhs.lhs.val")
	assert(4, tree2.lhs.rhs.val, "tree2.lhs.rhs.val")

	var tree3 *Tree = &Tree{
		1,
		&Tree{
			2,
			&Tree{3, 0, 0},
			&Tree{4, 0, 0},
		},
		0,
	}

	println("\ntype Tree struct {val int;lhs *Tree;rhs *Tree};")
	println("var tree3 *Tree=&Tree{1,&Tree{2,&Tree{3,0,0},&Tree{4,0,0},},0,};\n")
	assert(1, tree3.val, "tree3.val")
	assert(2, tree3.lhs.val, "tree3.lhs.val")
	assert(3, tree3.lhs.lhs.val, "tree3.lhs.lhs.val")
	assert(4, tree3.lhs.rhs.val, "tree3.lhs.rhs.val")

	tree4 := &Tree{
		1,
		&Tree{
			2,
			&Tree{3, 0, 0},
			&Tree{4, 0, 0},
		},
		0,
	}

	println("\ntype Tree struct {val int;lhs *Tree;rhs *Tree};")
	println("var tree4 *Tree=&Tree{1,&Tree{2,&Tree{3,0,0},&Tree{4,0,0},},0,};\n")
	assert(1, tree4.val, "tree4.val")
	assert(2, tree4.lhs.val, "tree4.lhs.val")
	assert(3, tree4.lhs.lhs.val, "tree4.lhs.lhs.val")
	assert(4, tree4.lhs.rhs.val, "tree4.lhs.rhs.val")

	println("OK")
}
