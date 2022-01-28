package test

var g3 byte = 3
var g4 int16 = 4
var g5 int = 5
var g6 int64 = 6

func main() {
	var x1 [3]int = [3]int{1, 2, 3}
	assert(1, x1[0], "var x1 [3]int=[3]int{1,2,3}; x1[0]")
	assert(2, x1[1], "var x1 [3]int=[3]int{1,2,3}; x1[1]")
	assert(3, x1[2], "var x1 [3]int=[3]int{1,2,3}; x1[2]")

	var x2 [2][3]int = [2][3]int{{1, 2, 3}, {4, 5, 6}}
	assert(2, x2[0][1], "var x2 [2][3]int=[2][3]int{{1,2,3},{4,5,6}}; x2[0][1]")
	assert(4, x2[1][0], "var x2 [2][3]int=[2][3]int{{1,2,3},{4,5,6}}; x2[1][0]")
	assert(6, x2[1][2], "var x2 [2][3]int=[2][3]int{{1,2,3},{4,5,6}}; x2[1][2]")

	var x3 [3]int = [3]int{}
	assert(0, x3[0], "var x3 [3]int=[3]int{}; x3[0]")
	assert(0, x3[1], "var x3 [3]int=[3]int{}; x3[1]")
	assert(0, x3[2], "var x3 [3]int=[3]int{}; x3[2]")

	var x4 [2][3]int = [2][3]int{{1, 2}}
	assert(2, x4[0][1], "var x4 [2][3]int=[2][3]int{{1,2}}; x4[0][1]")
	assert(0, x4[1][0], "var x4 [2][3]int=[2][3]int{{1,2}}; x4[1][0]")
	assert(0, x4[1][2], "var x4 [2][3]int=[2][3]int{{1,2}}; x4[1][2]")
	var x5 [4]byte = "abc"
	assert('a', x5[0], "var x5 [4]byte=\"abc\"; x5[0]")
	assert('c', x5[2], "var x5 [4]byte=\"abc\"; x5[2]")
	assert(0, x5[3], "var x5 [4]byte=\"abc\"; x5[3]")

	var x6 string = "abc"
	assert('a', x6[0], "var x6 string=\"abc\"; x6[0]")
	assert('c', x6[2], "var x6 string=\"abc\"; x6[2]")
	assert(0, x6[3], "var x6 string=\"abc\"; x6[3]")

	var x7 [2][4]byte = [2][4]byte{"abc", "def"}
	assert('a', x7[0][0], "var x7 [2][4]byte={\"abc\",\"def\"}; x7[0][0]")
	assert(0, x7[0][3], "var x7 [2][4]byte={\"abc\",\"def\"}; x7[0][3]")
	assert('d', x7[1][0], "var x7 [2][4]byte={\"abc\",\"def\"}; x7[1][0]")
	assert('f', x7[1][2], "var x7 [2][4]byte={\"abc\",\"def\"}; x7[1][2]")

	var x8 [2]string = [2]string{"abc", "def"}
	assert('a', x8[0][0], "var x8 [2]string=[2]string{\"abc\",\"def\"}; x8[0][0]")
	assert(0, x8[0][3], "var x8 [2]string=[2]string{\"abc\",\"def\"}; x8[0][3]")
	assert('d', x8[1][0], "var x8 [2]string=[2]string{\"abc\",\"def\"}; x8[1][0]")
	assert('f', x8[1][2], "var x8 [2]string=[2]string{\"abc\",\"def\"}; x8[1][2]")

	// assert(4, ({ int x[]={1,2,3,4}; x[3]; }));
	// assert(16, ({ int x[]={1,2,3,4}; sizeof(x); }));
	// assert(4, ({ char x[]="foo"; sizeof(x); }));

	type T9 string
	var x9 T9 = "foo"
	var y9 T9 = "x"
	assert(8, Sizeof(x9), "type T9 string; var x9 T9=\"foo\"; var y9 T9=\"x\"; Sizeof(x9)")
	assert(8, Sizeof(y9), "type T9 string; var x9 T9=\"foo\"; var y9 T9=\"x\"; Sizeof(y9)")
	var x10 T9 = "foo"
	var y10 T9 = "x"
	assert(8, Sizeof(x10), "type T9 string; var x10 T9=\"foo\"; var y10 T9=\"x\"; Sizeof(x10)")
	assert(8, Sizeof(y10), "type T9 string; var x10 T9=\"foo\"; var y10 T9=\"x\"; Sizeof(y10)")

	// assert(4, ({ typedef char T[]; T x="foo"; T y="x"; sizeof(x); }));
	// assert(2, ({ typedef char T[]; T x="foo"; T y="x"; sizeof(y); }));
	// assert(2, ({ typedef char T[]; T x="x"; T y="foo"; sizeof(x); }));
	// assert(4, ({ typedef char T[]; T x="x"; T y="foo"; sizeof(y); }));
	type T11 struct {
		a int
		b int
		c int
	}
	var x11 T11 = T11{1, 2, 3}
	assert(1, x11.a, "type T11 struct {a int; b int; c int;}; var x11 T11=T11{1,2,3}; x11.a")
	assert(2, x11.b, "type T11 struct {a int; b int; c int;}; var x11 T11=T11{1,2,3}; x11.b")
	assert(3, x11.c, "type T11 struct {a int; b int; c int;}; var x11 T11=T11{1,2,3}; x11.c")
	var x12 T11 = T11{1}
	assert(1, x12.a, "type T11 struct {a int; b int; c int;}; var x12 T11={1}; x12.a")
	assert(0, x12.b, "type T11 struct {a int; b int; c int;}; var x12 T11={1}; x12.b")
	assert(0, x12.c, "type T11 struct {a int; b int; c int;}; var x12 T11={1}; x12.c")
	type T13 struct {
		a int
		b int
	}
	var x13 [2]T13 = [2]T13{T13{1, 2}, T13{3, 4}}
	assert(1, x13[0].a, "type T13 struct {a int; b int;}; var x [2]T13=[2]T13{{1,2},{3,4}}; x13[0].a")
	assert(2, x13[0].b, "type T13 struct {a int; b int;}; var x [2]T13=[2]T13{{1,2},{3,4}}; x13[0].b")
	assert(3, x13[1].a, "type T13 struct {a int; b int;}; var x [2]T13=[2]T13{{1,2},{3,4}}; x13[1].a")
	assert(4, x13[1].b, "type T13 struct {a int; b int;}; var x [2]T13=[2]T13{{1,2},{3,4}}; x13[1].b")
	type T14 struct {
		a int
		b int
	}
	var x14 [2]T14 = [2]T14{T14{1, 2}}
	assert(0, x14[1].b, "type T14 struct {a int; b int;}; var x14 [2]T14=[2]T14{{1,2}}; x14[1].b")
	type T15 struct {
		a int
		b int
	}
	var x15 T15 = T15{}
	assert(0, x15.a, "type T15 struct {a int; b int;}; var x15 T15=T15{}; x15.a")
	assert(0, x15.b, "type T15 struct {a int; b int;}; var x15 T15=T15{}; x15.b")
	type T16 struct {
		a int
		b int
		c int
		d int
		e int
		f int
	}
	var x16 T16 = T16{1, 2, 3, 4, 5, 6}
	var y16 T16
	y16 = x16
	assert(5, y16.e, "type T16 struct {a int;b int;c int;d int;e int;f int;}; var x16 T16=T16{1,2,3,4,5,6};var y16 T16; y16=x16; y16.e")
	type T17 struct {a int;b int;}; var x17 T17=T17{1,2};var y17 T17;var z17 T17; z17=y17=x17;
	assert(2, z17.b, "type T17 struct {a int;b int;}; var x17 T17=T17{1,2};var y17 T17,var z17 T17; z=y=x; z.b")
	type T18 struct {a int;b int;}; var x18 T18=T18{1,2};var y18 T18=x18;
	assert(1, y18.a, "type T18 struct {a int;b int;}; var x18 T18=T18{1,2};var y T18=x18; y18.a");

	assert(3, g3, "g3")
	assert(4, g4, "g4")
	assert(5, g5, "g5")
	assert(6, g6, "g6")

	println("OK")
}
