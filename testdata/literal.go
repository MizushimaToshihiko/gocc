package test

func assert(want int, act int, code string)
func println(format string)

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

	println("OK")
}
