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

	// assert(4, ({ typedef char T[]; T x="foo"; T y="x"; sizeof(x); }));
	// assert(2, ({ typedef char T[]; T x="foo"; T y="x"; sizeof(y); }));
	// assert(2, ({ typedef char T[]; T x="x"; T y="foo"; sizeof(x); }));
	// assert(4, ({ typedef char T[]; T x="x"; T y="foo"; sizeof(y); }));

	// assert(1, ({ struct {int a; int b; int c;} x={1,2,3}; x.a; }));
	// assert(2, ({ struct {int a; int b; int c;} x={1,2,3}; x.b; }));
	// assert(3, ({ struct {int a; int b; int c;} x={1,2,3}; x.c; }));
	// assert(1, ({ struct {int a; int b; int c;} x={1}; x.a; }));
	// assert(0, ({ struct {int a; int b; int c;} x={1}; x.b; }));
	// assert(0, ({ struct {int a; int b; int c;} x={1}; x.c; }));

	// assert(1, ({ struct {int a; int b;} x[2]={{1,2},{3,4}}; x[0].a; }));
	// assert(2, ({ struct {int a; int b;} x[2]={{1,2},{3,4}}; x[0].b; }));
	// assert(3, ({ struct {int a; int b;} x[2]={{1,2},{3,4}}; x[1].a; }));
	// assert(4, ({ struct {int a; int b;} x[2]={{1,2},{3,4}}; x[1].b; }));

	// assert(0, ({ struct {int a; int b;} x[2]={{1,2}}; x[1].b; }));

	// assert(0, ({ struct {int a; int b;} x={}; x.a; }));
	// assert(0, ({ struct {int a; int b;} x={}; x.b; }));

	// assert(5, ({ typedef struct {int a,b,c,d,e,f;} T; T x={1,2,3,4,5,6}; T y; y=x; y.e; }));
	// assert(2, ({ typedef struct {int a,b;} T; T x={1,2}; T y, z; z=y=x; z.b; }));

	// assert(1, ({ typedef struct {int a,b;} T; T x={1,2}; T y=x; y.a; }));

	// assert(4, ({ union { int a; char b[4]; } x={0x01020304}; x.b[0]; }));
	// assert(3, ({ union { int a; char b[4]; } x={0x01020304}; x.b[1]; }));

	// assert(0x01020304, ({ union { struct { char a,b,c,d; } e; int f; } x={{4,3,2,1}}; x.f; }));

	// assert(3, g3);
	// assert(4, g4);
	// assert(5, g5);
	// assert(6, g6);

	println("\nOK")
}
