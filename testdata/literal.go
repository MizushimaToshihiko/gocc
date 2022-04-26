package test_literal

func assert(want int, act int, code string)
func println(format ...string)

func main() {
	assert(97, 'a', "'a'")
	assert(10, '\n', "'\\n'")

	assert(511, 0o777, "0o777")
	assert(0, 0x0, "0x0")
	assert(10, 0xa, "0xa")
	assert(10, 0xA, "0xA")
	assert(48879, 0xbeef, "0xbeef")
	assert(48879, 0xBEEF, "0xBEEF")
	assert(0, 0b0, "0b0")
	assert(1, 0b1, "0b1")
	assert(47, 0b101111, "0b101111")

	// '_' test
	assert(384, 0_600, "0_600")
	assert(42, 4_2, "4_2")
	assert(195951310, 0x_BadFace, "0x_BadFace")
	assert(801915078, 0x_67_7a_2f_cc_40_c6, "0x_67_7a_2f_cc_40_c6")
	assert(170141183460469, 170_141183_460469, "170_141183_460469")
	assert(1000, 1_0_0_0, "1_0_0_0")

	assert(4, Sizeof(0), "Sizeof(0)")
	assert(4, Sizeof(2147483647), "Sizeof(2147483647)")
	assert(8, Sizeof(2147483648), "Sizeof(2147483648)")

	// cannot tokenize for strconv.ParseInt function in tokenize.go
	assert(-1, 0xffffffffffffffff, "0xffffffffffffffff")
	assert(4, Sizeof(0xffffffffffffffff), "Sizeof(0xffffffffffffffff)")
	assert(-1, 0xffffffffffffffff>>63, "0xffffffffffffffff>>63")
	assert(-1, 18446744073709551615, "18446744073709551615")
	assert(4, Sizeof(18446744073709551615), "Sizeof(18446744073709551615)")
	assert(-1, 18446744073709551615>>63, "18446744073709551615>>63")
	assert(-1, 0xffffffffffffffff, "0xffffffffffffffff")
	assert(4, Sizeof(0xffffffffffffffff), "Sizeof(0xffffffffffffffff)")
	assert(-1, 0xffffffffffffffff>>63, "0xffffffffffffffff>>63")
	assert(-1, 01777777777777777777777, "01777777777777777777777")
	assert(4, Sizeof(01777777777777777777777), "Sizeof(01777777777777777777777)")
	assert(-1, 01777777777777777777777>>63, "01777777777777777777777>>63")
	assert(-1, 0b1111111111111111111111111111111111111111111111111111111111111111,
		"0b1111111111111111111111111111111111111111111111111111111111111111")
	assert(4, Sizeof(0b1111111111111111111111111111111111111111111111111111111111111111),
		"Sizeof(0b1111111111111111111111111111111111111111111111111111111111111111)")
	assert(-1, 0b1111111111111111111111111111111111111111111111111111111111111111>>63,
		"0b1111111111111111111111111111111111111111111111111111111111111111>>63")

	assert(8, Sizeof(2147483648), "Sizeof(2147483648)")
	assert(4, Sizeof(2147483647), "Sizeof(2147483647)")

	assert(8, Sizeof(0x1ffffffff), "Sizeof(0x1ffffffff)")
	assert(4, Sizeof(0x7ffffffe), "Sizeof(0xffffffff)")
	assert(1, 0xffffffff>>31, "0xffffffff>>31")

	assert(8, Sizeof(040000000000), "Sizeof(040000000000)")
	assert(4, Sizeof(017777777775), "Sizeof(017777777775)")
	assert(1, 037777777777>>31, "037777777777>>31")

	assert(8, Sizeof(0b111111111111111111111111111111111), "Sizeof(0b111111111111111111111111111111111)")
	assert(4, Sizeof(0b1111111111111111111111111111110), "Sizeof(0b11111111111111111111111111111111)")
	assert(1, 0b11111111111111111111111111111111>>31, "0b11111111111111111111111111111111>>31")

	assert(-1, 1<<31>>31, "1<<31>>31")
	assert(-1, 01<<31>>31, "01<<31>>31")
	assert(-1, 0x1<<31>>31, "0x1<<31>>31")
	assert(-1, 0b1<<31>>31, "0b1<<31>>31")

	assert(0, 0.0, "0.0")
	assert(1, 1.0, "1.0")
	assert(300000000, 3e+8, "3e+8")
	assert(16, 0x10.1p0, "0x10.1p0")
	assert(1000, .1e4, ".1e4")

	assert(16, 0x1_0.1p0, "0x1_0.1p0")
	assert(16, 0x_10.1p0, "0x_10.1p0")
	assert(348, 0x15e-2, "0x15e-2")
	assert(15, 0.15e+0_2, "0.15e+0_2")

	assert(4, Sizeof(8), "Sizeof(8)")
	assert(8, Sizeof(0.3), "Sizeof(0.3)")
	assert(8, Sizeof(0.), "Sizeof(0.)")
	assert(8, Sizeof(.0), "Sizeof(.0)")
	assert(8, Sizeof(5.), "Sizeof(5.)")
	assert(8, Sizeof(2.0), "Sizeof(2.0)")

	// assert(8, Sizeof("あいうえお"), "Sizeof(\"あいうえお\")")

	println("OK")
}
