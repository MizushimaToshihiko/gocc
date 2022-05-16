package test_multiassign

func assert(want int, act int, code string)
func println(format ...string)

func strcmp(s1 string, s2 string) int

func multiRet() (int, int, int, int, int, int) {
	return 1, 2, 3, 4, 5, 6
}

func multiRetStr() (string, string, string, string, string, string) {
	return "abc", "def", "ghi", "jkl", "mno", "pqr"
}

func multiRetFloat() (float64, float64) {
	return 0.1, 0.2
}

type gT01 struct {
	a int
	b string
}

func multiRetStruct() (int, gT01, string) {
	var g = gT01{
		a: 1,
		b: "aaa",
	}
	return g.a, g, g.b
}

func retStruct() gT01 {
	var g = gT01{
		a: 1,
		b: "aaa",
	}
	return g
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

	var a03, b03 float64
	a03, b03 = multiRetFloat()
	println("%lf", a03)
	println("%lf", b03)
	assert(1, a03 == 0.1, "a03==0.1")
	assert(1, b03 == 0.2, "b03==0.2")

	var a04 int
	var b04 string
	var c04 gT01
	a04, c04, b04 = multiRetStruct()
	assert(1, a04, "a04")
	assert(0, strcmp(b04, "aaa"), "strcmp(b04, \"aaa\")")
	assert(1, c04.a, "c04.a")
	assert(0, strcmp(c04.b, "aaa"), "strcmp(c04.b, \"aaa\")")

	var c05 gT01
	c05 = retStruct()
	assert(1, c05.a, "c05.a")
	assert(0, strcmp(c05.b, "aaa"), "strcmp(c05.b, \"aaa\")")

	var a05, b05, c05, d05, e05, f05, g05 = 1, 2, 3, 4, 5, 6, 7
	assert(1, a05, "a05")
	assert(2, b05, "b05")
	assert(3, c05, "c05")
	assert(4, d05, "d05")
	assert(5, e05, "e05")
	assert(6, f05, "f05")
	assert(7, g05, "g05")
	a05, b05, c05, _, e05, f05, g05 = g05, f05, e05, d05, c05, b05, a05
	assert(7, a05, "a05")
	assert(6, b05, "b05")
	assert(5, c05, "c05")
	assert(4, d05, "d05")
	assert(3, e05, "e05")
	assert(2, f05, "f05")
	assert(1, g05, "g05")
	a05, b05, c05, d05, e05, f05, g05 = 1, 2, 3, 4, 5, 6, 7
	assert(1, a05, "a05")
	assert(2, b05, "b05")
	assert(3, c05, "c05")
	assert(4, d05, "d05")
	assert(5, e05, "e05")
	assert(6, f05, "f05")
	assert(7, g05, "g05")

	var a06,
		b06,
		c06,
		d06,
		e06,
		f06,
		g06 = "aaa",
		"bbb",
		"ccc",
		"ddd",
		"eee",
		"fff",
		"ggg"
	assert(0, strcmp(a06, "aaa"), "strcmp(a06, \"aaa\")")
	assert(0, strcmp(b06, "bbb"), "strcmp(b06, \"bbb\")")
	assert(0, strcmp(c06, "ccc"), "strcmp(c06, \"ccc\")")
	assert(0, strcmp(d06, "ddd"), "strcmp(d06, \"ddd\")")
	assert(0, strcmp(e06, "eee"), "strcmp(e06, \"eee\")")
	assert(0, strcmp(f06, "fff"), "strcmp(f06, \"fff\")")
	assert(0, strcmp(g06, "ggg"), "strcmp(g06, \"ggg\")")
	a06, b06, c06, _, e06, f06, g06 = g06, f06, e06, d06, c06, b06, a06
	assert(0, strcmp(a06, "ggg"), "strcmp(a06, \"ggg\")")
	assert(0, strcmp(b06, "fff"), "strcmp(b06, \"fff\")")
	assert(0, strcmp(c06, "eee"), "strcmp(c06, \"eee\")")
	assert(0, strcmp(d06, "ddd"), "strcmp(d06, \"ddd\")")
	assert(0, strcmp(e06, "ccc"), "strcmp(e06, \"ccc\")")
	assert(0, strcmp(f06, "bbb"), "strcmp(f06, \"bbb\")")
	assert(0, strcmp(g06, "aaa"), "strcmp(g06, \"aaa\")")
	a06,
		b06,
		c06,
		d06,
		e06,
		f06,
		g06 = "aaa",
		"bbb",
		"ccc",
		"ddd",
		"eee",
		"fff",
		"ggg"
	assert(0, strcmp(a06, "aaa"), "strcmp(a06, \"aaa\")")
	assert(0, strcmp(b06, "bbb"), "strcmp(b06, \"bbb\")")
	assert(0, strcmp(c06, "ccc"), "strcmp(c06, \"ccc\")")
	assert(0, strcmp(d06, "ddd"), "strcmp(d06, \"ddd\")")
	assert(0, strcmp(e06, "eee"), "strcmp(e06, \"eee\")")
	assert(0, strcmp(f06, "fff"), "strcmp(f06, \"fff\")")
	assert(0, strcmp(g06, "ggg"), "strcmp(g06, \"ggg\")")

	var a07, b07, c07, d07 = 0.1, 0.2, 0.3, 0.4
	assert(1, a07 == 0.1, "a07==0.1")
	assert(1, b07 == 0.2, "b07==0.2")
	assert(1, c07 == 0.3, "c07==0.3")
	assert(1, d07 == 0.4, "d07==0.4")
	a07, b07, c07, d07 = d07, c07, b07, a07
	println("a07: %lf", a07)
	assert(1, a07 == 0.4, "a07==0.4")
	assert(1, b07 == 0.3, "b07==0.3")
	assert(1, c07 == 0.2, "c07==0.2")
	assert(1, d07 == 0.1, "d07==0.1")
	a07, b07, c07, d07 = 0.1, 0.2, 0.3, 0.4
	assert(1, a07 == 0.1, "a07==0.1")
	assert(1, b07 == 0.2, "b07==0.2")
	assert(1, c07 == 0.3, "c07==0.3")
	assert(1, d07 == 0.4, "d07==0.4")

	println("OK")
}
