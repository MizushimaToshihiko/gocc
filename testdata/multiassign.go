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

// func multiRetFloat() (float64, float64) {
// 	return 0.1, 0.2
// }

func main() {
	var a01, b01 int
	a01, b01 = multiRet()
	assert(1, a01, "a01")
	assert(2, b01, "b01")

	var a02, b02 string
	a02, b02 = multiRetStr()
	assert(0, strcmp(a02, "abc"), "strcmp(a02, \"abc\")")
	assert(0, strcmp(b02, "def"), "strcmp(b02, \"def\")")

	// // flonumは未対応
	// var a03, b03 float64
	// a03, b03 = multiRetFloat()
	// println("%f", a03)
	// println("%f", b03)
	// assert(1, a03 == 0.1, "a03==0.1")
	// assert(1, b03 == 0.2, "b03==0.2")

	println("OK")
}
