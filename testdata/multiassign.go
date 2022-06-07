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
	a int64
	b string
}

func multiRetStruct(a int64, b string) (int64, gT01, string) {
	var g = gT01{
		a: a,
		b: b,
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

func multiRet2Struct(a int, b string) (int, gT01, gT01, string) {
	var g1 = gT01{
		a: a,
		b: b,
	}
	var g2 = gT01{
		a: a + 1,
		b: "bbb",
	}
	return g1.a, g1, g2, g1.b
}

type gT02 struct {
	a [20]int
}

func multiRet2BigStruct(a int, b string) (int, gT01, gT02, string) {
	var g1 = gT01{
		a: int64(a + 2),
		b: b,
	}
	var g2 = gT02{
		a: [20]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
	}
	return int(g1.a), g1, g2, g1.b
}

func multiRetSlice() ([]int, []string, []float64) {
	return []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		[]string{"aaa", "bbb", "ccc", "ddd", "eee", "fff"},
		[]float64{1., 2., 3.}
}

func multiRetArged(x int, y string) (gT01, gT02) {
	var g1 = gT01{
		a: x + 1,
		b: y,
	}
	var g2 = gT02{
		a: [20]int{
			x + 1,
			x + 2,
			x + 3,
			x + 4,
			x + 5,
			x + 6,
			x + 7,
			x + 8,
			x + 9,
			x + 10,
			x + 11,
			x + 12,
			x + 13,
			x + 14,
			x + 15,
			x + 16,
			x + 17,
			x + 18,
			x + 19,
			x + 20,
		},
	}
	return g1, g2
}

func multiRetArged4(x int, y string) (gT01, gT01, gT01, gT02) {
	var g1 = gT01{a: x + 1, b: y}
	var g2 = gT01{a: x + 2, b: y}
	var g3 = gT01{a: x + 3, b: y}
	var g4 = gT02{
		a: [20]int{
			x + 1,
			x + 2,
			x + 3,
			x + 4,
			x + 5,
			x + 6,
			x + 7,
			x + 8,
			x + 9,
			x + 10,
			x + 11,
			x + 12,
			x + 13,
			x + 14,
			x + 15,
			x + 16,
			x + 17,
			x + 18,
			x + 19,
			x + 20,
		},
	}
	return g1, g2, g3, g4
}

func multiRetArr() ([3]int64, [4]int64, [5]int64) {
	var a, b, c, d, e int64 = 1, 2, 3, 4, 5
	return [3]int64{a, b, c}, [4]int64{a, b, c, d}, [5]int64{a, b, c, d, e}
}

func multiRetFloatArr() ([3]float64, [4]float64, [5]float64) {
	var a, b, c, d, e float64 = 1.1, 2.2, 3.3, 4.4, 5.5
	return [3]float64{a, b, c}, [4]float64{a, b, c, d}, [5]float64{a, b, c, d, e}
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

	var c041 gT01
	c041 = retStruct()
	assert(1, int(c041.a), "c041.a")
	assert(0, strcmp(c041.b, "aaa"), "strcmp(c041.b, \"aaa\")")

	var a03, b03 float64
	a03, b03 = multiRetFloat()
	println("%lf", a03)
	println("%lf", b03)
	assert(1, a03 == 0.1, "a03==0.1")
	assert(1, b03 == 0.2, "b03==0.2")

	var a04 int64
	var b04 string
	var c04 gT01
	a04, c04, b04 = multiRetStruct(1, "aaa")
	assert(1, a04, "a04")
	assert(0, strcmp(b04, "aaa"), "strcmp(b04, \"aaa\")")
	assert(0, strcmp(c04.b, "aaa"), "strcmp(c04.b, \"aaa\")")
	assert(1, c04.a, "c04.a")

	var a042 int
	var b042 string
	var c042 gT01
	var d042 gT01
	a042, c042, d042, b042 = multiRet2Struct(1, "aaa")
	assert(1, a042, "a042")
	assert(0, strcmp(b042, "aaa"), "strcmp(b042, \"aaa\")")
	assert(1, c042.a, "c042.a")
	assert(0, strcmp(c042.b, "aaa"), "strcmp(c042.b, \"aaa\")")
	assert(2, d042.a, "d042.a")
	assert(0, strcmp(d042.b, "bbb"), "strcmp(d042.b, \"bbb\")")

	var a043 int
	var b043 string
	var c043 gT01
	var d043 gT02
	a043, c043, d043, b043 = multiRet2BigStruct(1, "ccc")
	assert(3, a043, "a043")
	assert(0, strcmp(b043, "ccc"), "strcmp(b043, \"ccc\")")
	assert(3, c043.a, "c043.a")
	assert(0, strcmp(c043.b, "ccc"), "strcmp(c043.b, \"ccc\")")
	assert(1, d043.a[0], "d043.a[0]")
	assert(2, d043.a[1], "d043.a[1]")
	assert(3, d043.a[2], "d043.a[2]")
	assert(4, d043.a[3], "d043.a[3]")
	assert(5, d043.a[4], "d043.a[4]")
	assert(6, d043.a[5], "d043.a[5]")
	assert(7, d043.a[6], "d043.a[6]")
	assert(8, d043.a[7], "d043.a[7]")
	assert(9, d043.a[8], "d043.a[8]")
	assert(10, d043.a[9], "d043.a[9]")
	assert(11, d043.a[10], "d043.a[10]")
	assert(12, d043.a[11], "d043.a[11]")
	assert(13, d043.a[12], "d043.a[12]")
	assert(14, d043.a[13], "d043.a[13]")
	assert(15, d043.a[14], "d043.a[14]")
	assert(16, d043.a[15], "d043.a[15]")
	assert(17, d043.a[16], "d043.a[16]")
	assert(18, d043.a[17], "d043.a[17]")
	assert(19, d043.a[18], "d043.a[18]")
	assert(20, d043.a[19], "d043.a[19]")

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

	var a08 []int
	var b08 []string
	var c08 []float64
	a08, b08, c08 = multiRetSlice()
	assert(1, a08[0], "a08[0]")
	assert(11, a08[10], "a08[10]")
	assert(20, a08[19], "a08[19]")
	assert(0, strcmp(b08[0], "aaa"), "strcmp(b08[0], \"aaa\")")
	assert(0, strcmp(b08[3], "ddd"), "strcmp(b08[3], \"ddd\")")
	assert(0, strcmp(b08[5], "fff"), "strcmp(b08[5], \"fff\")")
	assert(1, c08[0], "c08[0]")
	assert(2, c08[1], "c08[2]")
	assert(3, c08[2], "c08[3]")

	var a09 gT01
	var b09 gT02
	a09, b09 = multiRetArged(100, "abc")
	assert(101, a09.a, "a09.a")
	assert(0, strcmp(a09.b, "abc"), "strcmp(a09.b, \"abc\")")
	assert(101, b09.a[0], "b09.a[0]")
	assert(102, b09.a[1], "b09.a[1]")
	assert(103, b09.a[2], "b09.a[2]")
	assert(104, b09.a[3], "b09.a[3]")
	assert(105, b09.a[4], "b09.a[4]")
	assert(106, b09.a[5], "b09.a[5]")
	assert(110, b09.a[9], "b09.a[9]")
	assert(115, b09.a[14], "b09.a[14]")
	assert(120, b09.a[19], "b09.a[19]")

	var a10 gT01
	var b10 gT01
	var c10 gT01
	var d10 gT02
	a10, b10, c10, d10 = multiRetArged4(200, "abc")
	assert(201, a10.a, "a10.a")
	assert(0, strcmp(a10.b, "abc"), "strcmp(a10.b, \"abc\")")
	assert(202, b10.a, "b10.a")
	assert(0, strcmp(b10.b, "abc"), "strcmp(b10.b, \"abc\")")
	assert(203, c10.a, "c10.a")
	assert(0, strcmp(c10.b, "abc"), "strcmp(c10.b, \"abc\")")
	assert(201, d10.a[0], "d10.a[0]")
	assert(202, d10.a[1], "d10.a[1]")
	assert(203, d10.a[2], "d10.a[2]")
	assert(204, d10.a[3], "d10.a[3]")
	assert(205, d10.a[4], "d10.a[4]")
	assert(206, d10.a[5], "d10.a[5]")
	assert(210, d10.a[9], "d10.a[9]")
	assert(215, d10.a[14], "d10.a[14]")
	assert(220, d10.a[19], "d10.a[19]")

	var a11 [3]int64
	var b11 [4]int64
	var c11 [5]int64
	a11, b11, c11 = multiRetArr()
	assert(1, a11[0], "a11[0]")
	assert(2, a11[1], "a11[1]")
	assert(3, b11[2], "b11[2]")
	assert(4, b11[3], "b11[3]")
	assert(5, c11[4], "c11[4]")

	var a12 [3]float64
	var b12 [4]float64
	var c12 [5]float64
	a12, b12, c12 = multiRetFloatArr()
	assert(1, a12[0] == 1.1, "a12[0]==1.1")
	assert(1, a12[2] == 3.3, "a12[2]==3.3")
	assert(1, b12[0] == 1.1, "b12[0]==1.1")
	assert(1, b12[3] == 4.4, "b12[3]==4.4")
	assert(1, c12[4] == 5.5, "c12[4]==5.5")

	var a13, b13, c13, d13, e13 int
	a13, b13, c13, _, d13, e13 = multiRet()
	assert(1, a13, "a13")
	assert(2, b13, "b13")
	assert(3, c13, "c13")
	assert(5, d13, "d13")
	assert(6, e13, "e13")

	println("OK")
}
