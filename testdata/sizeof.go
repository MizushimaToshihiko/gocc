package test

func main() {
	// assert(1, Sizeof(byte), "Sizeof(byte)")
	// assert(2, Sizeof(int16), "Sizeof(int16)")
	// assert(4, Sizeof(int), "Sizeof(int)")
	// assert(8, Sizeof(int64), "Sizeof(int64)")
	// type T1 struct {
	// 	a int
	// 	b int
	// }
	// assert(8, Sizeof(T1), "Sizeof(type T1 struct {a int; b int;};)")

	// var i int = 0
	// assert(1, Sizeof(i+1), "Sizeof(i+1)")
	assert(8, Sizeof(-10+int64(5)), "Sizeof(-10+int64(5))")
	// assert(8, Sizeof(-10 - (long)5));
	// assert(8, Sizeof(-10 * (long)5));
	// assert(8, Sizeof(-10 / (long)5));
	// assert(8, Sizeof((long)-10 + 5));
	// assert(8, Sizeof((long)-10 - 5));
	// assert(8, Sizeof((long)-10 * 5));
	// assert(8, Sizeof((long)-10 / 5));

	// assert(1, ({ char i; Sizeof(++i); }));
	// assert(1, ({ char i; Sizeof(i++); }));

	// assert(8, Sizeof(int(*)[10]));
	// assert(8, Sizeof(int(*)[][10]));

	println("OK")
}
