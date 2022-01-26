package test

func main() {
	assert(97, 'a', "'a'")
	assert(10, '\n', "'\n'")

	assert(511, 0o777, "0o777")
	assert(0, 0x0, "0x0")
	assert(10, 0xa, "0xa")
	assert(10, 0xA, "0xA")
	assert(48879, 0xbeef, "0xbeef")
	assert(48879, 0xBEEF, "0xBEEF")
	assert(0, 0b0, "0b0")
	assert(1, 0b1, "0b1")
	assert(47, 0b101111, "0b101111")

	println("OK")
}
