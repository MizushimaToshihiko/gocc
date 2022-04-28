package test_multiassign

func assert(want int, act int, code string)
func println(format ...string)

func multiRet() (int, int) {
	return 1, 2
}

func strcmp(s1 string, s2 string) int

func multiRetStr() (string, string) {
	return "abc", "def"
}

func main() {
	var a01, b01 int
	a01, b01 = multiRet()
	assert(1, a01, "a01")
	assert(2, b01, "b01")

	var a02, b02 string
	a02, b02 = multiRetStr()
	assert(0, strcmp(a02, "abc"), "strcmp(a02, \"abc\")")
	assert(0, strcmp(b02, "def"), "strcmp(b02, \"def\")")

	println("OK")
}
