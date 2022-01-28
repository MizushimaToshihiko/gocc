package test

func main() {
	assert(131585, int(8590066177), "int(8590066177)")
	assert(513, int16(8590066177), "int16(8590066177)")
	assert(1, byte(8590066177), "byte(8590066177)")
	assert(1, int64(1), "int64(1)")
	assert(0, int64(&*(*int)(0)), "int64(&*(*int)(0))")
	// assert(513, ({ int x=512; *(char *)&x=1; x; }));
	// assert(5, ({ int x=5; long y=(long)&x; *(int*)y; }));

	// (void)1;

	printf("OK\n")
}
