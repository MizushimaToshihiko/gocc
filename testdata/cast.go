package test

func assert(want int, act int, code string)
func println(format string)

func main() {
	assert(131585, int(8590066177), "int(8590066177)")
	assert(513, int16(8590066177), "int16(8590066177)")
	assert(1, int8(8590066177), "int8(8590066177)")
	assert(1, int64(1), "int64(1)")
	// var x int=512; *(*byte)(&x)=1;
	// assert(513, x, "var x int=512; *(*byte)(&x)=1; x");
	// assert(5, ({ int x=5; long y=(long)&x; *(int*)y; }));

	assert(-1, int8(255), "int8(255)")
	assert(255, uint8(255), "uint8(255)")
	assert(-1, int16(65535), "int16(65535)")
	assert(65535, uint16(65535), "uint16(65535)")
	assert(-1, int(0xffffffff), "int(0xffffffff)")
	assert(-1, int32(0xffffffff), "int32(0xffffffff)")
	assert(0xffffffff, uint(0xffffffff), "uint(0xffffffff)")
	assert(0xffffffff, uint32(0xffffffff), "uint32(0xffffffff)")

	assert(1, -1 < 1, "-1<1")
	assert(0, -1 < uint(1), "-1<uint(1)")
	assert(254, int8(127)+int8(127), "int8(127)+int8(127)")
	assert(65534, int16(32767)+int16(32767), "int16(32767)+int16(32767)")
	assert(-1, -1>>1, "-1>>1")
	assert(-1, uint64(-1), "uint64(-1)")
	assert(2147483647, uint(-1)>>1, "uint(-1)>>1")
	assert(-50, (-100)/2, "(-100)/2")
	assert(2147483598, uint(-100)/2, "uint(-100)/2")
	// Floating point exception???
	// assert(9223372036854775758, uint64(-100)/2, "uint64(-100)/2")
	assert(0, int64(-1)/uint(100), "int64(-1)/uint(100)")
	assert(-2, (-100)%7, "(-100)%7")
	assert(2, uint(-100)%7, "uint(-100)%7")
	// Floating point exception???
	// assert(6, uint64(-100)%9, "uint64(-100)%9")

	assert(65535, (int(uint16(65535))), "(int(uint16(65535))")
	var x uint16 = 65535
	assert(65535, x, "var x uint16=65535;x")
	var x uint16 = 65535
	assert(65535, int(x), "var x uint16=65535;int(x)")

	type T1 int16
	var x T1 = 65535
	assert(-1, int(x), "type T1 int16;var x T1=65535;int(x)")
	type T2 uint16
	var x T2 = 65535
	assert(65535, int(x), "type T2 uint16;var x T2=65535;int(x)")

	assert(0, bool(0.0), "bool(0.0)")
	assert(1, bool(0.1), "bool(0.1)")
	assert(3, int8(3.0), "int8(3.0)")
	assert(1000, int16(1000.3), "int16(1000.3)")
	assert(3, int(3.99), "int(3.99)")
	assert(2000000000000000, int64(2e15), "int64(2e15)")
	assert(3, float32(3.5), "float32(3.5)")
	assert(5, float64(float32(5.5)), "float64(float32(5.5))")
	assert(3, float32(3), "float32(3)")
	assert(3, float64(3), "float64(3)")

	println("OK")
}
