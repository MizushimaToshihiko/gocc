package test

func assert(want int, act int, code string)
func println(format string)

func main() {
	var i int = 0
	switch 3 {
	case 5 - 2 + 0*3:
		i++
	}
	assert(1, i, "var i int =0; switch(3) { case 5-2+0*3: i++; }")
	var x1 [1 + 1]int
	assert(8, Sizeof(x1), "var x [1+1]int; Sizeof(x1)")
	var x2 [8 - 2]byte
	assert(6, Sizeof(x2), "var x [8-2]byte; Sizeof(x2)")
	var x3 [2 * 3]byte
	assert(6, Sizeof(x3), "var x [2*3]byte; Sizeof(x3)")
	var x4 [12 / 4]byte
	assert(3, Sizeof(x4), "var x4 [12/4]byte; Sizeof(x4)")
	var x5 [12 % 10]byte
	assert(2, Sizeof(x5), "var x5 [12%10]byte; Sizeof(x5)")
	var x6 [0b110 & 0b101]byte
	assert(0b100, Sizeof(x6), "var x6 [0b110&0b101]byte; Sizeof(x6)")
	var x7 [0b110 | 0b101]byte
	assert(0b111, Sizeof(x7), "var x7 [0b110|0b101]byte; Sizeof(x7)")
	var x8 [0b111 ^ 0b001]byte
	assert(0b110, Sizeof(x8), "var x8 [0b111^0b001]byte; Sizeof(x8)")

	var x9 [1 << 2]byte
	assert(4, Sizeof(x9), "var x9 [1<<2]byte; Sizeof(x9)")
	var x10 [4 >> 1]byte
	assert(2, Sizeof(x10), "var x10 [4>>1]byte; Sizeof(x10)")
	var x11 [(1 == 1) + 1]byte
	assert(2, Sizeof(x11), "var x11 [(1==1)+1]byte; Sizeof(x11)")
	var x12 [(1 != 1) + 1]byte
	assert(1, Sizeof(x12), "var x12 [(1!=1)+1]byte; Sizeof(x12)")
	var x13 [(1 < 1) + 1]byte
	assert(1, Sizeof(x13), "var x13 [(1<1)+1]byte; Sizeof(x13)")
	var x14 [(1 <= 1) + 1]byte
	assert(2, Sizeof(x11), "var x14 [(1<=1)+1]byte; Sizeof(x14)")
	var x15 [!0 + 1]byte
	assert(2, Sizeof(x15), "var x15 [!0+1]byte; Sizeof(x15)")
	var x16 [!1 + 1]byte
	assert(1, Sizeof(x16), "var x16 [!1+1]byte; Sizeof(x16)")
	var x17 [^-3]byte
	assert(2, Sizeof(x17), "var x17 [^-3]byte; Sizeof(x17)")
	var x18 [(5 || 6) + 1]byte
	assert(2, Sizeof(x18), "var x18 [(5||6)+1]byte; Sizeof(x18)")
	var x19 [(0 || 0) + 1]byte
	assert(1, Sizeof(x19), "var x19 [(0||0)+1]byte; Sizeof(x19)")
	var x20 [(1 && 1) + 1]byte
	assert(2, Sizeof(x20), "var x20 [(1&&1)+1]byte; Sizeof(x20)")
	var x21 [(1 && 0) + 1]byte
	assert(1, Sizeof(x21), "var x21 [(1&&0)+1]byte; Sizeof(x21)")
	var x22 [int(3)]byte
	assert(3, Sizeof(x22), "var x22 [int(3)]byte; Sizeof(x22)")

	var x23 [(1,3)]byte
	assert(3, Sizeof(x23), "var x23 [(1,3)]byte; Sizeof(x23)")
	var x24 [byte(0xffffff0f)]byte
	assert(15, Sizeof(x24), "var x24 [byte(0xffffff0f)]byte; Sizeof(x24)")
	var x25 [int16(0xffff010f)]byte
	assert(0x10f, Sizeof(x25), "var x25 [int16(0xffff010f)]byte; Sizeof(x25)")
	
	// error occures
	// var x26 [int(0xfffffffffff)+5]byte
	// assert(4, Sizeof(x26), "var x26 [int(0xfffffffffff)+5]byte; Sizeof(x26)");

	// Below is not supported in Go.
	// var x26 [(*int)(0) + 2]byte
	// assert(8, Sizeof(x26), "var x26 [(int*)0+2]byte; Sizeof(x26)")
	// assert(12, ({ char x[(int*)16-1]; Sizeof(x); }));
	// assert(3, ({ char x[(int*)16-(int*)4]; Sizeof(x); }));

	println("OK")
}
