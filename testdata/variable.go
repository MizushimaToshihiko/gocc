package test_variable

func assert(want int, act int, code string)
func println(format ...string)
func printf(format ...string) int

#include "test.h"

var g1 int
var g2 [4]int

func ret3() int {
	return 3
}

var _, g3  = 1, 2

func main() {
	var a1 int
	a1 = 3
	ASSERT(3, a1)
	var a2 int = 3
	ASSERT(3, a2)
	var a3 int = 3
	var z3 int = 5
	ASSERT(8, a3+z3)
	var a4 int;var b4 int; a4=b4=3; // a4=b4=3 is not supported in Go.
	ASSERT(6, a4+b4)
	var x5 int;
	ASSERT(4, Sizeof(x5))
	var x6 *int;
	ASSERT(8, Sizeof(x6))
	var x7 [4]int;
	ASSERT(16, Sizeof(x7))
	var x8 [3][4]int;
	ASSERT(48, Sizeof(x8))
	ASSERT(16, Sizeof(*x8))
	ASSERT(4, Sizeof(**x8))
	ASSERT(4, Sizeof(**x8+1))
	var x9 int=1;
	// ASSERT(4, Sizeof(x9=2), "var x9 int=1; Sizeof(x9=2)") // It isn't supported in Go.
	ASSERT(1, x9)

	ASSERT(0, g1)
	g1=3
	ASSERT(3, g1);
	g2[0]=0; g2[1]=1; g2[2]=2; g2[3]=3;
	ASSERT(0, g2[0]);
	ASSERT(1, g2[1]);
	ASSERT(2, g2[2]);
	ASSERT(3, g2[3]);

	ASSERT(4, Sizeof(g1));
	ASSERT(16, Sizeof(g2));

	var x10 byte=1;
	ASSERT(1, x10);
	var x11 byte=1;var y11 byte=2;
	ASSERT(1, x11);
	ASSERT(2, y11);
	var x12 byte;
	ASSERT(1, Sizeof(x12));
	var x13 [10]byte;
	ASSERT(10, Sizeof(x13));
	var x14 int=2; {var x14 int=3;};
	ASSERT(2, x14);
	var x15 int=2; {var x15 int=3;};var y15 int=4;
	ASSERT(2, x15);
	var x16 int=2; {x16=3;};
	ASSERT(3, x16);
	var x17 int;var y17 int;var z17 byte;var a17 *byte=&y17;var b17 *byte=&z17;
	ASSERT(7, b17-a17);
	var x18 int;var y18 byte;var z18 int;var a18 *byte=&y18;var b18 *byte=&z18;
	ASSERT(1, b18-a18);
	var x19 int64;
	ASSERT(8, Sizeof(x19));
	var x20 int16;
	ASSERT(2, Sizeof(x20));

	var x21 [3]*byte;
	ASSERT(24, Sizeof(x21));
	var x22 *[3]byte;
	ASSERT(8, Sizeof(x22));

	// Belows is not supported yet.
	// var (x23 byte);
	// assert(1, Sizeof(x23), "var (x23) byte; Sizeof(x23)");
	// assert(3, ({ char (x)[3]; sizeof(x); }));
	// assert(12, ({ char (x[3])[4]; sizeof(x); }));
	// assert(4, ({ char (x[3])[4]; sizeof(x[0]); }));
	// assert(3, ({ char *x[3]; char y; x[0]=&y; y=3; x[0][0]; }));
	// assert(4, ({ char x[3]; char (*y)[3]=x; y[0][0]=4; y[0][0]; }));

	_ = 1
	_, g1 = 1, 2
	ASSERT(2, g1)
	var _x33 int
	_x33, _ = 1, 3
	ASSERT(1, _x33)
	_, x34 := 4, 5
	ASSERT(5, x34)
	var x35, _ = ret3()+3, printf("blank ident test\n")
	ASSERT(6, x35)
	ASSERT(2, g3)

	println("OK")
}
