package test_multiassign

func assert(want int, act int, code string)
func println(format ...string)

func strcmp(s1 string, s2 string) int

type gT01 struct {
	a int64
	b string
}

type gT02 struct {
	a [20]int
}

func multiRet() (int, int, int, int, int, int) {
	return 1, 2, 3, 4, 5, 6
}

func multiRetStr() (string, string, string, string, string, string) {
	return "abc", "def", "ghi", "jkl", "mno", "pqr"
}

func multiRetFloat() (float64, float64) {
	return 0.1, 0.2
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

func multiRet2Struct(a int, b string) (int, int, int, string, gT01, gT01) {
	var g1 = gT01{
		a: int64(a),
		b: b,
	}
	var g2 = gT01{
		a: int64(a + 1),
		b: "bbb",
	}
	return int(g1.a), int(g2.a), 3, g1.b, g1, g2
}

func multiRet2BigStruct(a int, b string) (int, gT01, gT02, gT02, string) {
	var g1 = gT01{
		a: int64(a + 2),
		b: b,
	}
	var g2 = gT02{
		a: [20]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
	}
	var g3 = gT02{
		a: [20]int{a + 1, a + 2, a + 3, a + 4, a + 5, a + 6, a + 7, a + 8, a + 9, a + 10, a + 11, a + 12, a + 13, a + 14, a + 15, a + 16, a + 17, a + 18, a + 19, a + 20},
	}
	return int(g1.a), g1, g2, g3, g1.b
}

func multiRetSlice() ([]int, []string, []float64) {
	return []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		[]string{"aaa", "bbb", "ccc", "ddd", "eee", "fff"},
		[]float64{1., 2., 3.}
}

func multiRetArged(x int, y string) (gT01, gT02) {
	var g1 = gT01{
		a: int64(x + 1),
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

func multiRet8Int(a int) (int, int, int, int, int, int, int, int) {
	return a + 1, a + 2, a + 3, a + 4, a + 5, a + 6, a + 7, a + 8
}

func multiRet8Float32(a float32) (float32, float32, float32, float32, float32, float32, float32, float32) {
	return a + 1, a + 2, a + 3, a + 4, a + 5, a + 6, a + 7, a + 8
}

func multiRet9Float32() (float32, float32, float32, float32, float32, float32, float32, float32, float32) {
	return 1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9
}

func multiRet9Float64() (float64, float64, float64, float64, float64, float64, float64, float64, float64) {
	return 1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9
}

type gT011 struct {
	a int64
}

func multiRet8SmallStruct() (gT011, gT011, gT011, gT011, gT011, gT011, gT011, gT011) {
	return gT011{11}, gT011{22}, gT011{33}, gT011{44}, gT011{55}, gT011{66}, gT011{77}, gT011{88}
}

type gT012 struct {
	a int64
	b int
}

func multiRet8SmallStruct2() (gT012, gT012, gT012, gT012, gT012, gT012, gT012, gT012) {
	return gT012{a: 11, b: 1}, gT012{22, 2}, gT012{33, 3}, gT012{44, 4}, gT012{55, 5}, gT012{66, 6}, gT012{77, 7}, gT012{88, 8}
}

// 大きな(16 bytes超)structを7以上返す関数
func multiRet12BigStruct() (gT02, gT02, gT02, gT02, gT02, gT02, gT02, gT02, gT02, gT02, gT02, gT02) {
	return gT02{1}, gT02{2}, gT02{3}, gT02{4}, gT02{5}, gT02{6}, gT02{7}, gT02{8}, gT02{9}, gT02{10}, gT02{11}, gT02{12}
}

// b float64
type gT03 struct {
	a float64
}

// 以下未実装
// 小数点を含む16 bytes以下のstructを２つ以上返す関数
func multiRet2StructFlonum() (gT03, gT03) {
	// b: 2.2,
	var g1 = gT03{
		a: 1.1,
	}
	// b: 4.4,
	var g2 = gT03{
		a: 3.3,
	}
	return g1, g2
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
	var b042 int
	var c042 int
	var d042 string
	var e042 gT01
	var f042 gT01
	a042, b042, c042, d042, e042, f042 = multiRet2Struct(1, "aaa")
	assert(1, a042, "a042")
	assert(2, b042, "b042")
	assert(3, c042, "c042")
	assert(0, strcmp(d042, "aaa"), "strcmp(d042, \"aaa\")")
	assert(1, int(e042.a), "e042.a")
	assert(0, strcmp(e042.b, "aaa"), "strcmp(e042.b, \"aaa\")")
	assert(2, int(f042.a), "f042.a")
	assert(0, strcmp(f042.b, "bbb"), "strcmp(f042.b, \"bbb\")")

	var a043 int
	var b043 string
	var c043 gT01
	var d043 gT02
	var e043 gT02
	a043, c043, d043, e043, b043 = multiRet2BigStruct(1, "ccc")
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
	assert(2, e043.a[0], "e043.a[0]")
	assert(3, e043.a[1], "e043.a[1]")
	assert(4, e043.a[2], "e043.a[2]")
	assert(5, e043.a[3], "e043.a[3]")
	assert(6, e043.a[4], "e043.a[4]")
	assert(7, e043.a[5], "e043.a[5]")
	assert(8, e043.a[6], "e043.a[6]")
	assert(9, e043.a[7], "e043.a[7]")
	assert(10, e043.a[8], "e043.a[8]")
	assert(11, e043.a[9], "e043.a[9]")
	assert(12, e043.a[10], "e043.a[10]")
	assert(13, e043.a[11], "e043.a[11]")
	assert(14, e043.a[12], "e043.a[12]")
	assert(15, e043.a[13], "e043.a[13]")
	assert(16, e043.a[14], "e043.a[14]")
	assert(17, e043.a[15], "e043.a[15]")
	assert(18, e043.a[16], "e043.a[16]")
	assert(19, e043.a[17], "e043.a[17]")
	assert(20, e043.a[18], "e043.a[18]")
	assert(21, e043.a[19], "e043.a[19]")

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

	a14, _, c14 := multiRetArr()
	assert(1, a14[0], "a14[0]")
	assert(2, a14[1], "a14[1]")
	assert(3, a14[2], "a14[2]")
	assert(4, c14[3], "c14[3]")
	assert(5, c14[4], "c14[4]")

	a15, b15, c15 := multiRetSlice()
	assert(1, a15[0], "a15[0]")
	assert(11, a15[10], "a15[10]")
	assert(20, a15[19], "a15[19]")
	assert(0, strcmp(b15[0], "aaa"), "strcmp(b15[0], \"aaa\")")
	assert(0, strcmp(b15[3], "ddd"), "strcmp(b15[3], \"ddd\")")
	assert(0, strcmp(b15[5], "fff"), "strcmp(b15[5], \"fff\")")
	assert(1, c15[0], "c15[0]")
	assert(2, c15[1], "c15[2]")
	assert(3, c15[2], "c15[3]")
	println("len(a15): %d", len(a15)) // 今のSliceのデータ構造やparserでは関数間のlen,capの受渡ができない
	println("len(b15): %d", len(b15))
	println("len(c15): %d", len(c15))

	a16, b16, c16, d16, e16, f16, g16, h16 := multiRet8Int(1)
	assert(2, a16, "a16")
	assert(3, b16, "b16")
	assert(4, c16, "c16")
	assert(5, d16, "d16")
	assert(6, e16, "e16")
	assert(7, f16, "f16")
	assert(8, g16, "g16")
	assert(9, h16, "h16")

	a17, b17, c17, d17, e17, f17, g17, h17 := multiRet8Float32(1.0)
	assert(2, a17, "a17")
	assert(3, b17, "b17")
	assert(4, c17, "c17")
	assert(5, d17, "d17")
	assert(6, e17, "e17")
	assert(7, f17, "f17")
	assert(8, g17, "g17")
	assert(9, h17, "h17")

	a18, b18, c18, d18, e18, f18, g18, h18, i18 := multiRet9Float64()
	assert(1, a18 == 1.1, "a18==1.1")
	assert(1, b18 == 2.2, "b18==2.2")
	assert(1, c18 == 3.3, "c18==3.3")
	assert(1, d18 == 4.4, "d18==4.4")
	assert(1, e18 == 5.5, "e18==5.5")
	assert(1, f18 == 6.6, "f18==6.6")
	assert(1, g18 == 7.7, "g18==7.7")
	assert(1, h18 == 8.8, "h18==8.8")
	assert(1, i18 == 9.9, "i18==9.9")

	a181, b181, c181, d181, e181, f181, g181, h181 := multiRet8SmallStruct()
	assert(11, a181.a, "a181.a")
	assert(22, b181.a, "b181.a")
	assert(33, c181.a, "c181.a")
	assert(44, d181.a, "d181.a")
	assert(55, e181.a, "e181.a")
	assert(66, f181.a, "f181.a")
	assert(77, g181.a, "g181.a")
	assert(88, h181.a, "h181.a")

	a182, b182, c182, d182, e182, f182, g182, h182 := multiRet8SmallStruct2()
	assert(11, a182.a, "a182.a")
	assert(1, a182.b, "a182.b")
	assert(22, b182.a, "b182.a")
	assert(2, b182.b, "b182.b")
	assert(33, c182.a, "c182.a")
	assert(3, c182.b, "c182.b")
	assert(44, d182.a, "d182.a")
	assert(4, d182.b, "d182.b")
	assert(55, e182.a, "e182.a")
	assert(5, e182.b, "e182.b")
	assert(66, f182.a, "f182.a")
	assert(6, f182.b, "f182.b")
	assert(77, g182.a, "g182.a")
	assert(7, g182.b, "g182.b")
	assert(88, h182.a, "h182.a")
	assert(8, h182.b, "h182.b")

	// 16bytesを超えるstructを７つ以上返す関数
	a19, b19, c19, d19, e19, f19, g19, h19, i19, j19, k19, l19 := multiRet12BigStruct() //
	assert(1, a19.a[0], "a19.a[0]")
	assert(2, b19.a[0], "b19.a[0]")
	assert(3, c19.a[0], "c19.a[0]")
	assert(4, d19.a[0], "d19.a[0]")
	assert(5, e19.a[0], "e19.a[0]")
	assert(6, f19.a[0], "f19.a[0]")
	assert(7, g19.a[0], "g19.a[0]")
	assert(8, h19.a[0], "h19.a[0]")
	assert(9, i19.a[0], "i19.a[0]")
	assert(10, j19.a[0], "j19.a[0]")
	assert(11, k19.a[0], "k19.a[0]")
	assert(12, l19.a[0], "l19.a[0]")

	// 以下未実装
	// hasFlonumな構造体で16bytes以下のものを2つ以上返す関数
	a20, b20 := multiRet2StructFlonum()
	println("a20.a: %f", a20.a)
	println("b20.a: %f", b20.a)
	assert(1, a20.a == 1.1, "a20.a==1.1")
	assert(1, b20.a == 3.3, "b20.a==3.3")

	println("OK")
}
