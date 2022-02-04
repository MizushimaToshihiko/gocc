package test

func assert(want int, act int, code string)
func println(format string)

func main() {
	assert(131585, int(8590066177), "int(8590066177)")
	assert(513, int16(8590066177), "int16(8590066177)")
	assert(1, byte(8590066177), "byte(8590066177)")
	assert(1, int64(1), "int64(1)")
	// var x int=512; *(*byte)(&x)=1;
	// assert(513, x, "var x int=512; *(*byte)(&x)=1; x");
	// assert(5, ({ int x=5; long y=(long)&x; *(int*)y; }));

	// (void)1;

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
	assert(254, int8(127)+int8(127), int8(127)+int8(127))
	assert(65534, int16(32767)+int16(32767), "int16(32767)+int16(32767)")
	assert(-1, -1>>1, "-1>>1")
	assert(-1, uint64(-1), "uint64(-1)")
	assert(2147483647, (uint(-1))>>1, "(uint(-1))>>1")
	// assert(-50, (-100)/2);
	// assert(2147483598, ((unsigned)-100)/2);
	// assert(9223372036854775758, ((unsigned long)-100)/2);
	// assert(0, ((long)-1)/(unsigned)100);
	// assert(-2, (-100)%7);
	// assert(2, ((unsigned)-100)%7);
	// assert(6, ((unsigned long)-100)%9);

	// assert(65535, (int)(unsigned short)65535);
	// assert(65535, ({ unsigned short x = 65535; x; }));
	// assert(65535, ({ unsigned short x = 65535; (int)x; }));

	// assert(-1, ({ typedef short T; T x = 65535; (int)x; }));
	// assert(65535, ({ typedef unsigned short T; T x = 65535; (int)x; }));

	println("OK")
}
