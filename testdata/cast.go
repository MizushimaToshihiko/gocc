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

	println("OK")
}
