package test_multiassign

func assert(want int, act int, code string)
func println(format ...string)

func multiRet() (int, int) {
	return 1, 2
}

func main() {
	var a, b int
	a, b = multiRet()
	assert(1, a, "a")
	assert(2, b, "b")

	println("OK")
}
