var g1 int

// func assert(want int64, ac int64, code *byte) {
// 	if want == ac {
// 		printf("\n%s => %ld\n", code, ac)
// 	} else {
// 		printf("\n%s => %ld expeted but got %ld\n", code, want, ac)
// 		exit(1)
// 	}
// }

func ret1(a int) int {
	return 1
}

func main() {
	assert(0, g1, "g1")
	g1 = 3
	assert(3, g1, "g1")
	assert(2, ret1(1), "ret1()")

	printf("\nOK\n")
}
