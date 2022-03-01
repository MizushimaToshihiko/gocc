package test_variable

func assert(want int, act int, code string)
func println(format string)

var g1 int
var g2 [4]int

func main() {
	var a1 int
	a1 = 3
	assert(3, a1, "var a1 int; a1=3; a1")
	var a2 int = 3
	assert(3, a2, "var a2 int=3; a2")
	var a3 int = 3
	var z3 int = 5
	assert(8, a3+z3, "var a3 int=3;var z3 int=5; a3+z3")
	var a4 int;var b4 int; a4=b4=3; // a4=b4=3 is not supported in Go.
	assert(6, a4+b4, "var a4 int;var b4 int; a4=b4=3; a4+b4")
	var x5 int;
	assert(4, Sizeof(x5), "var x5 int; Sizeof(x5)")
	var x6 *int;
	assert(8, Sizeof(x6), "var x6 *int; Sizeof(x6)")
	var x7 [4]int;
	assert(16, Sizeof(x7), "var x7 [4]int; Sizeof(x7)")
	var x8 [3][4]int;
	assert(48, Sizeof(x8), "var x8 [3][4]int; Sizeof(x8)")
	assert(16, Sizeof(*x8), "var x8 [3][4]int; Sizeof(*x8)")
	assert(4, Sizeof(**x8), "var x8 [3][4]int; Sizeof(**x8)")
	assert(4, Sizeof(**x8+1), "var x8 [3][4]int; Sizeof(**x8+1)")
	var x9 int=1;
	assert(4, Sizeof(x9=2), "var x9 int=1; Sizeof(x9=2)")
	assert(1, x9, "var x9 int=1; Sizeof(x9=2); x9")

	assert(0, g1, "g1")
	g1=3
	assert(3, g1, "g1=3; g1");
	g2[0]=0; g2[1]=1; g2[2]=2; g2[3]=3;
	assert(0, g2[0], "g2[0]=0; g2[1]=1; g2[2]=2; g2[3]=3; g2[0]");
	assert(1, g2[1], "g2[0]=0; g2[1]=1; g2[2]=2; g2[3]=3; g2[1]");
	assert(2, g2[2], "g2[0]=0; g2[1]=1; g2[2]=2; g2[3]=3; g2[2]");
	assert(3, g2[3], "g2[0]=0; g2[1]=1; g2[2]=2; g2[3]=3; g2[3]");

	assert(4, Sizeof(g1), "Sizeof(g1)");
	assert(16, Sizeof(g2), "Sizeof(g2)");

	var x10 byte=1;
	assert(1, x10, "var x10 byte=1; x10");
	var x11 byte=1;var y11 byte=2;
	assert(1, x11, "var x11 byte=1;var y11 byte=2; x11");
	assert(2, y11, "var x11 byte=1;var y11 byte=2; y11");
	var x12 byte;
	assert(1, Sizeof(x12), "var x12 byte; Sizeof(x12)");
	var x13 [10]byte;
	assert(10, Sizeof(x13), "var x13 [10]byte; Sizeof(x13)");
	var x14 int=2; {var x14 int=3;};
	assert(2, x14, "var x14 int=2; {var x14 int=3;}; x14");
	var x15 int=2; {var x15 int=3;};var y15 int=4;
	assert(2, x15, "var x15 int=2; {var x15 int=3;};var y15 int=4; x15");
	var x16 int=2; {x16=3;};
	assert(3, x16, "var x16 int=2; { x16=3; }; x16");
	var x17 int;var y17 int;var z17 byte;var a17 *byte=&y17;var b17 *byte=&z17;
	assert(7, b17-a17, "var x17 int;var y17 int;var z17 byte;var a17 *byte=&y17;var b17 *byte=&z17; b17-a17");
	var x18 int;var y18 byte;var z18 int;var a18 *byte=&y18;var b18 *byte=&z18;
	assert(1, b18-a18, "var x18 int;var y18 byte;var z18 int;var a18 *byte=&y18;var b18 *byte=&z18; b18-a18");
	var x19 int64;
	assert(8, Sizeof(x19), "var x19 int64; Sizeof(x19)");
	var x20 int16;
	assert(2, Sizeof(x20), "var x20 int16; Sizeof(x20)");

	var x21 [3]*byte;
	assert(24, Sizeof(x21), "var x21 [3]*byte; Sizeof(x21)");
	var x22 *[3]byte;
	assert(8, Sizeof(x22), "var x22 *[3]byte; Sizeof(x22)");

	// Belows is not supported yet.
	// var (x23 byte);
	// assert(1, Sizeof(x23), "var (x23) byte; Sizeof(x23)");
	// assert(3, ({ char (x)[3]; sizeof(x); }));
	// assert(12, ({ char (x[3])[4]; sizeof(x); }));
	// assert(4, ({ char (x[3])[4]; sizeof(x[0]); }));
	// assert(3, ({ char *x[3]; char y; x[0]=&y; y=3; x[0][0]; }));
	// assert(4, ({ char x[3]; char (*y)[3]=x; y[0][0]=4; y[0][0]; }));

	println("OK");
}
