package test_multiassign

func assert(want int, act int, code string)
func println(format ...string)

func strcmp(s1 string, s2 string) int

#include "test.h"

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

type gT03 struct {
	a float64
}

// 小数点を含む16 bytes以下のstructを２つ以上返す関数
func multiRet2StructFlonum() (gT03, gT03, gT03, gT03, gT03, gT03, gT03) {
	var g1 = gT03{
		a: 1.1,
	}
	var g2 = gT03{
		a: 2.2,
	}
	return g1, g2, gT03{3.3}, gT03{4.4}, gT03{5.5}, gT03{6.6}, gT03{7.7}
}

type gT04 struct {
	a float64
	b float32
}

func multiRet2StructFlonum2() (gT04, gT04, gT04, gT04, gT04, gT04, gT04) {
	var g1 = gT04{
		a: 1.1,
		b: 2.2,
	}
	var g2 = gT04{
		a: 3.3,
		b: 4.4,
	}
	return g1, g2, gT04{5.5, 6.6}, gT04{7.7, 8.8}, gT04{9.9, 11.11}, gT04{12.12, 13.13}, gT04{14.14, 15.15}
}

type gT05 struct {
	a float32
	b int32
}

func multiRet2StructFlonum3() (gT05, gT05, gT05, gT05, gT05, gT05, gT05, gT05, gT05) {
	var g1 = gT05{
		a: 1.1,
		b: 2,
	}
	var g2 = gT05{
		a: 3.3,
		b: 4,
	}
	return g1, g2, gT05{5.5, 5}, gT05{6.6, 6}, gT05{7.7, 7}, gT05{8.8, 8}, gT05{9.9, 9}, gT05{11.11, 11}, gT05{12.12, 12}
}

type gT06 struct {
	a int
	b float64
}

func multiRet2StructFlonum4() (gT06, gT06, gT06, gT06, gT06, gT06, gT06, gT06, gT06, gT06) {
	var g1 = gT06{
		a: 1,
		b: 1.1,
	}
	var g2 = gT06{
		a: 2,
		b: 2.2,
	}
	return g1, g2, gT06{3, 3.3}, gT06{4, 4.4}, gT06{5, 5.5}, gT06{6, 6.6}, gT06{7, 7.7}, gT06{8, 8.8}, gT06{9, 9.9}, gT06{11, 11.11}
}

type gT07 struct {
	a float64
	b int
}

func multiRet2StructFlonum5() (gT07, gT07, gT07, gT07, gT07, gT07, gT07) {
	var g1 = gT07{
		a: 1.1,
		b: 1,
	}
	var g2 = gT07{
		a: 2.2,
		b: 2,
	}
	return g1, g2, gT07{3.3, 3}, gT07{4.4, 4}, gT07{5.5, 5}, gT07{6.6, 6}, gT07{7.7, 7}
}

func main() {
	var a01, b01, c01, d01, e01, f01 int
	a01, b01, c01, d01, e01, f01 = multiRet()
	ASSERT(1, a01)
	ASSERT(2, b01)
	ASSERT(3, c01)
	ASSERT(4, d01)
	ASSERT(5, e01)
	ASSERT(6, f01)

	var a02, b02, c02, d02, e02, f02 string
	a02, b02, c02, d02, e02, f02 = multiRetStr()
	ASSERT(0, strcmp(a02, "abc"))
	ASSERT(0, strcmp(b02, "def"))
	ASSERT(0, strcmp(c02, "ghi"))
	ASSERT(0, strcmp(d02, "jkl"))
	ASSERT(0, strcmp(e02, "mno"))
	ASSERT(0, strcmp(f02, "pqr"))

	var c041 gT01
	c041 = retStruct()
	ASSERT(1, int(c041.a))
	ASSERT(0, strcmp(c041.b, "aaa"))

	var a03, b03 float64
	a03, b03 = multiRetFloat()
	println("%lf", a03)
	println("%lf", b03)
	ASSERT(1, a03 == 0.1)
	ASSERT(1, b03 == 0.2)

	var a04 int64
	var b04 string
	var c04 gT01
	a04, c04, b04 = multiRetStruct(1, "aaa")
	ASSERT(1, a04)
	ASSERT(0, strcmp(b04, "aaa"))
	ASSERT(0, strcmp(c04.b, "aaa"))
	ASSERT(1, c04.a)

	var a042 int
	var b042 int
	var c042 int
	var d042 string
	var e042 gT01
	var f042 gT01
	a042, b042, c042, d042, e042, f042 = multiRet2Struct(1, "aaa")
	ASSERT(1, a042)
	ASSERT(2, b042)
	ASSERT(3, c042)
	ASSERT(0, strcmp(d042, "aaa"))
	ASSERT(1, int(e042.a))
	ASSERT(0, strcmp(e042.b, "aaa"))
	ASSERT(2, int(f042.a))
	ASSERT(0, strcmp(f042.b, "bbb"))

	var a043 int
	var b043 string
	var c043 gT01
	var d043 gT02
	var e043 gT02
	a043, c043, d043, e043, b043 = multiRet2BigStruct(1, "ccc")
	ASSERT(3, a043)
	ASSERT(0, strcmp(b043, "ccc"))
	ASSERT(3, c043.a)
	ASSERT(0, strcmp(c043.b, "ccc"))
	ASSERT(1, d043.a[0])
	ASSERT(2, d043.a[1])
	ASSERT(3, d043.a[2])
	ASSERT(4, d043.a[3])
	ASSERT(5, d043.a[4])
	ASSERT(6, d043.a[5])
	ASSERT(7, d043.a[6])
	ASSERT(8, d043.a[7])
	ASSERT(9, d043.a[8])
	ASSERT(10, d043.a[9])
	ASSERT(11, d043.a[10])
	ASSERT(12, d043.a[11])
	ASSERT(13, d043.a[12])
	ASSERT(14, d043.a[13])
	ASSERT(15, d043.a[14])
	ASSERT(16, d043.a[15])
	ASSERT(17, d043.a[16])
	ASSERT(18, d043.a[17])
	ASSERT(19, d043.a[18])
	ASSERT(20, d043.a[19])
	ASSERT(2, e043.a[0])
	ASSERT(3, e043.a[1])
	ASSERT(4, e043.a[2])
	ASSERT(5, e043.a[3])
	ASSERT(6, e043.a[4])
	ASSERT(7, e043.a[5])
	ASSERT(8, e043.a[6])
	ASSERT(9, e043.a[7])
	ASSERT(10, e043.a[8])
	ASSERT(11, e043.a[9])
	ASSERT(12, e043.a[10])
	ASSERT(13, e043.a[11])
	ASSERT(14, e043.a[12])
	ASSERT(15, e043.a[13])
	ASSERT(16, e043.a[14])
	ASSERT(17, e043.a[15])
	ASSERT(18, e043.a[16])
	ASSERT(19, e043.a[17])
	ASSERT(20, e043.a[18])
	ASSERT(21, e043.a[19])

	var a05, b05, c05, d05, e05, f05, g05 = 1, 2, 3, 4, 5, 6, 7
	ASSERT(1, a05)
	ASSERT(2, b05)
	ASSERT(3, c05)
	ASSERT(4, d05)
	ASSERT(5, e05)
	ASSERT(6, f05)
	ASSERT(7, g05)
	a05, b05, c05, _, e05, f05, g05 = g05, f05, e05, d05, c05, b05, a05
	ASSERT(7, a05)
	ASSERT(6, b05)
	ASSERT(5, c05)
	ASSERT(4, d05)
	ASSERT(3, e05)
	ASSERT(2, f05)
	ASSERT(1, g05)
	a05, b05, c05, d05, e05, f05, g05 = 1, 2, 3, 4, 5, 6, 7
	ASSERT(1, a05)
	ASSERT(2, b05)
	ASSERT(3, c05)
	ASSERT(4, d05)
	ASSERT(5, e05)
	ASSERT(6, f05)
	ASSERT(7, g05)

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
	ASSERT(0, strcmp(a06, "aaa"))
	ASSERT(0, strcmp(b06, "bbb"))
	ASSERT(0, strcmp(c06, "ccc"))
	ASSERT(0, strcmp(d06, "ddd"))
	ASSERT(0, strcmp(e06, "eee"))
	ASSERT(0, strcmp(f06, "fff"))
	ASSERT(0, strcmp(g06, "ggg"))
	a06, b06, c06, _, e06, f06, g06 = g06, f06, e06, d06, c06, b06, a06
	ASSERT(0, strcmp(a06, "ggg"))
	ASSERT(0, strcmp(b06, "fff"))
	ASSERT(0, strcmp(c06, "eee"))
	ASSERT(0, strcmp(d06, "ddd"))
	ASSERT(0, strcmp(e06, "ccc"))
	ASSERT(0, strcmp(f06, "bbb"))
	ASSERT(0, strcmp(g06, "aaa"))
	a06,
		b06,
		c06,
		d06,
		e06,
		f06,
		g06 = "aaa",
		"bbb", // comment
		"ccc",
		"ddd",
		"eee",
		"fff",
		"ggg"
	ASSERT(0, strcmp(a06, "aaa"))
	ASSERT(0, strcmp(b06, "bbb"))
	ASSERT(0, strcmp(c06, "ccc"))
	ASSERT(0, strcmp(d06, "ddd"))
	ASSERT(0, strcmp(e06, "eee"))
	ASSERT(0, strcmp(f06, "fff"))
	ASSERT(0, strcmp(g06, "ggg"))

	var a07, b07, c07, d07 = 0.1, 0.2, 0.3, 0.4
	ASSERT(1, a07 == 0.1)
	ASSERT(1, b07 == 0.2)
	ASSERT(1, c07 == 0.3)
	ASSERT(1, d07 == 0.4)
	a07, b07, c07, d07 = d07, c07, b07, a07
	println("a07: %lf", a07)
	ASSERT(1, a07 == 0.4)
	ASSERT(1, b07 == 0.3)
	ASSERT(1, c07 == 0.2)
	ASSERT(1, d07 == 0.1)
	a07, b07, c07, d07 = 0.1, 0.2, 0.3, 0.4
	ASSERT(1, a07 == 0.1)
	ASSERT(1, b07 == 0.2)
	ASSERT(1, c07 == 0.3)
	ASSERT(1, d07 == 0.4)

	var a08 []int
	var b08 []string
	var c08 []float64
	a08, b08, c08 = multiRetSlice()
	ASSERT(1, a08[0])
	ASSERT(11, a08[10])
	ASSERT(20, a08[19])
	ASSERT(0, strcmp(b08[0], "aaa"))
	ASSERT(0, strcmp(b08[3], "ddd"))
	ASSERT(0, strcmp(b08[5], "fff"))
	ASSERT(1, c08[0])
	ASSERT(2, c08[1])
	ASSERT(3, c08[2])

	var a09 gT01
	var b09 gT02
	a09, b09 = multiRetArged(100, "abc")
	ASSERT(101, a09.a)
	ASSERT(0, strcmp(a09.b, "abc"))
	ASSERT(101, b09.a[0])
	ASSERT(102, b09.a[1])
	ASSERT(103, b09.a[2])
	ASSERT(104, b09.a[3])
	ASSERT(105, b09.a[4])
	ASSERT(106, b09.a[5])
	ASSERT(110, b09.a[9])
	ASSERT(115, b09.a[14])
	ASSERT(120, b09.a[19])

	var a10 gT01
	var b10 gT01
	var c10 gT01
	var d10 gT02
	a10, b10, c10, d10 = multiRetArged4(200, "abc")
	ASSERT(201, a10.a)
	ASSERT(0, strcmp(a10.b, "abc"))
	ASSERT(202, b10.a)
	ASSERT(0, strcmp(b10.b, "abc"))
	ASSERT(203, c10.a)
	ASSERT(0, strcmp(c10.b, "abc"))
	ASSERT(201, d10.a[0])
	ASSERT(202, d10.a[1])
	ASSERT(203, d10.a[2])
	ASSERT(204, d10.a[3])
	ASSERT(205, d10.a[4])
	ASSERT(206, d10.a[5])
	ASSERT(210, d10.a[9])
	ASSERT(215, d10.a[14])
	ASSERT(220, d10.a[19])

	var a11 [3]int64
	var b11 [4]int64
	var c11 [5]int64
	a11, b11, c11 = multiRetArr()
	ASSERT(1, a11[0])
	ASSERT(2, a11[1])
	ASSERT(3, b11[2])
	ASSERT(4, b11[3])
	ASSERT(5, c11[4])

	var a12 [3]float64
	var b12 [4]float64
	var c12 [5]float64
	a12, b12, c12 = multiRetFloatArr()
	ASSERT(1, a12[0] == 1.1)
	ASSERT(1, a12[2] == 3.3)
	ASSERT(1, b12[0] == 1.1)
	ASSERT(1, b12[3] == 4.4)
	ASSERT(1, c12[4] == 5.5)

	var a13, b13, c13, d13, e13 int
	a13, b13, c13, _, d13, e13 = multiRet()
	ASSERT(1, a13)
	ASSERT(2, b13)
	ASSERT(3, c13)
	ASSERT(5, d13)
	ASSERT(6, e13)

	a14, _, c14 := multiRetArr()
	ASSERT(1, a14[0])
	ASSERT(2, a14[1])
	ASSERT(3, a14[2])
	ASSERT(4, c14[3])
	ASSERT(5, c14[4])

	a15, b15, c15 := multiRetSlice()
	ASSERT(1, a15[0])
	ASSERT(11, a15[10])
	ASSERT(20, a15[19])
	ASSERT(0, strcmp(b15[0], "aaa"))
	ASSERT(0, strcmp(b15[3], "ddd"))
	ASSERT(0, strcmp(b15[5], "fff"))
	ASSERT(1, c15[0])
	ASSERT(2, c15[1])
	ASSERT(3, c15[2])
	println("len(a15): %d", len(a15)) // 今のSliceのデータ構造やparserでは関数間のlen,capの受渡ができない
	println("len(b15): %d", len(b15))
	println("len(c15): %d", len(c15))

	a16, b16, c16, d16, e16, f16, g16, h16 := multiRet8Int(1)
	ASSERT(2, a16)
	ASSERT(3, b16)
	ASSERT(4, c16)
	ASSERT(5, d16)
	ASSERT(6, e16)
	ASSERT(7, f16)
	ASSERT(8, g16)
	ASSERT(9, h16)

	a17, b17, c17, d17, e17, f17, g17, h17 := multiRet8Float32(1.0)
	ASSERT(2, a17)
	ASSERT(3, b17)
	ASSERT(4, c17)
	ASSERT(5, d17)
	ASSERT(6, e17)
	ASSERT(7, f17)
	ASSERT(8, g17)
	ASSERT(9, h17)

	a18, b18, c18, d18, e18, f18, g18, h18, i18 := multiRet9Float64()
	ASSERT(1, a18 == 1.1)
	ASSERT(1, b18 == 2.2)
	ASSERT(1, c18 == 3.3)
	ASSERT(1, d18 == 4.4)
	ASSERT(1, e18 == 5.5)
	ASSERT(1, f18 == 6.6)
	ASSERT(1, g18 == 7.7)
	ASSERT(1, h18 == 8.8)
	ASSERT(1, i18 == 9.9)

	a181, b181, c181, d181, e181, f181, g181, h181 := multiRet8SmallStruct()
	ASSERT(11, a181.a)
	ASSERT(22, b181.a)
	ASSERT(33, c181.a)
	ASSERT(44, d181.a)
	ASSERT(55, e181.a)
	ASSERT(66, f181.a)
	ASSERT(77, g181.a)
	ASSERT(88, h181.a)

	a182, b182, c182, d182, e182, f182, g182, h182 := multiRet8SmallStruct2()
	ASSERT(11, a182.a)
	ASSERT(1, a182.b)
	ASSERT(22, b182.a)
	ASSERT(2, b182.b)
	ASSERT(33, c182.a)
	ASSERT(3, c182.b)
	ASSERT(44, d182.a)
	ASSERT(4, d182.b)
	ASSERT(55, e182.a)
	ASSERT(5, e182.b)
	ASSERT(66, f182.a)
	ASSERT(6, f182.b)
	ASSERT(77, g182.a)
	ASSERT(7, g182.b)
	ASSERT(88, h182.a)
	ASSERT(8, h182.b)

	// 16bytesを超えるstructを７つ以上返す関数
	a19, b19, c19, d19, e19, f19, g19, h19, i19, j19, k19, l19 := multiRet12BigStruct() //
	ASSERT(1, a19.a[0])
	ASSERT(2, b19.a[0])
	ASSERT(3, c19.a[0])
	ASSERT(4, d19.a[0])
	ASSERT(5, e19.a[0])
	ASSERT(6, f19.a[0])
	ASSERT(7, g19.a[0])
	ASSERT(8, h19.a[0])
	ASSERT(9, i19.a[0])
	ASSERT(10, j19.a[0])
	ASSERT(11, k19.a[0])
	ASSERT(12, l19.a[0])

	// hasFlonumな構造体で16bytes以下のものを2つ以上返す関数
	a20, b20, c20, d20, e20, f20, g20 := multiRet2StructFlonum()
	ASSERT(1, a20.a == 1.1)
	ASSERT(1, b20.a == 2.2)
	ASSERT(1, c20.a == 3.3)
	ASSERT(1, d20.a == 4.4)
	ASSERT(1, e20.a == 5.5)
	ASSERT(1, f20.a == 6.6)
	ASSERT(1, g20.a == 7.7)

	a21, b21, c21, d21, e21, f21, g21 := multiRet2StructFlonum2()
	ASSERT(1, a21.a == 1.1)
	ASSERT(1, a21.b == float32(2.2))
	ASSERT(1, b21.a == 3.3)
	ASSERT(1, b21.b == float32(4.4))
	ASSERT(1, c21.a == 5.5)
	ASSERT(1, c21.b == float32(6.6))
	ASSERT(1, d21.a == 7.7)
	ASSERT(1, d21.b == float32(8.8))
	ASSERT(1, e21.a == 9.9)
	ASSERT(1, e21.b == float32(11.11))
	ASSERT(1, f21.a == 12.12)
	ASSERT(1, f21.b == float32(13.13))
	ASSERT(1, g21.a == 14.14)
	ASSERT(1, g21.b == float32(15.15))

	a22, b22, c22, d22, e22, f22, g22, h22, i22 := multiRet2StructFlonum3()
	ASSERT(1, a22.a == float32(1.1))
	ASSERT(2, a22.b)
	ASSERT(1, b22.a == float32(3.3))
	ASSERT(4, b22.b)
	ASSERT(1, c22.a == float32(5.5))
	ASSERT(5, c22.b)
	ASSERT(1, d22.a == float32(6.6))
	ASSERT(6, d22.b)
	ASSERT(1, e22.a == float32(7.7))
	ASSERT(7, e22.b)
	ASSERT(1, f22.a == float32(8.8))
	ASSERT(8, f22.b)
	ASSERT(1, g22.a == float32(9.9))
	ASSERT(9, g22.b)
	ASSERT(1, h22.a == float32(11.11))
	ASSERT(11, h22.b)
	ASSERT(1, i22.a == float32(12.12))
	ASSERT(12, i22.b)

	a23, b23, c23, d23, e23, f23, g23, h23, i23, j23 := multiRet2StructFlonum4()
	ASSERT(1, a23.a)
	ASSERT(1, a23.b == 1.1)
	ASSERT(2, b23.a)
	ASSERT(1, b23.b == 2.2)
	ASSERT(3, c23.a)
	ASSERT(1, c23.b == 3.3)
	ASSERT(4, d23.a)
	ASSERT(1, d23.b == 4.4)
	ASSERT(5, e23.a)
	ASSERT(1, e23.b == 5.5)
	ASSERT(6, f23.a)
	ASSERT(1, f23.b == 6.6)
	ASSERT(7, g23.a)
	ASSERT(1, g23.b == 7.7)
	ASSERT(8, h23.a)
	ASSERT(1, h23.b == 8.8)
	ASSERT(9, i23.a)
	ASSERT(1, i23.b == 9.9)
	ASSERT(11, j23.a)
	ASSERT(1, j23.b == 11.11)

	a24, b24, c24, d24, e24, f24, g24 := multiRet2StructFlonum5()
	ASSERT(1, a24.a == 1.1)
	ASSERT(1, a24.b)
	ASSERT(1, b24.a == 2.2)
	ASSERT(2, b24.b)
	ASSERT(1, c24.a == 3.3)
	ASSERT(3, c24.b)
	ASSERT(1, d24.a == 4.4)
	ASSERT(4, d24.b)
	ASSERT(1, e24.a == 5.5)
	ASSERT(5, e24.b)
	ASSERT(1, f24.a == 6.6)
	ASSERT(6, f24.b)
	ASSERT(1, g24.a == 7.7)
	ASSERT(7, g24.b)

	println("OK")
}
