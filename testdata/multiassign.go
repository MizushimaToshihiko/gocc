package test_multiassign

func assert(want int, act int, code string)
func println(format ...string)

func multiRet() (int, int, int, int, int, int) {
	return 1, 2, 3, 4, 5, 6
}

func strcmp(s1 string, s2 string) int

func multiRetStr() (string, string, string, string, string, string) {
	return "abc", "def", "ghi", "jkl", "mno", "pqr"
}

func multiRetFloat() (float64, float64) {
	return 0.1, 0.2
}

func main() {
	var a01, b01, c01, d01, e01, f01 int
	a01, b01, c01, d01, e01, f01 = multiRet()
	assert(1, a01, "a01")
	assert(2, b01, "b01")
	assert(3, c01, "c01")
	assert(4, d01, "d01")
	assert(5, e01, "e01")
	assert(6, f01, "f01")

	var a02, b02, c02, d02, e02, f02 string
	a02, b02, c02, d02, e02, f02 = multiRetStr()
	assert(0, strcmp(a02, "abc"), "strcmp(a02, \"abc\")")
	assert(0, strcmp(b02, "def"), "strcmp(b02, \"def\")")
	assert(0, strcmp(c02, "ghi"), "strcmp(c02, \"ghi\")")
	assert(0, strcmp(d02, "jkl"), "strcmp(d02, \"jkl\")")
	assert(0, strcmp(e02, "mno"), "strcmp(e02, \"mno\")")
	assert(0, strcmp(f02, "pqr"), "strcmp(f02, \"pqr\")")

	// flonumは未対応
	var a03, b03 float64
	a03, b03 = multiRetFloat()
	println("%f", a03)
	println("%f", b03)
	assert(1, a03 == 0.1, "a03==0.1")
	assert(1, b03 == 0.2, "b03==0.2")

	println("OK")
}
