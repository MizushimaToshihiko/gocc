package test

func assert(want int, act int, code string)
func println(format string)

// var g3 byte = 3
// var g4 int16 = 4
// var g5 int = 5
// var g6 int64 = 6
// var g9 [3]int = [3]int{0, 1, 2}

// type gT11 struct {
// 	a byte
// 	b int
// }

// var g11 [2]gT11 = [2]gT11{gT11{1, 2}, gT11{3, 4}}

// type gT12 struct{ a [2]int }

// var g12 [2]gT12 = [2]gT12{{{1, 2}}}
// var g17 string = "foobar"
// var g17_2 [7]byte = "foobar"
// var g18 string

// var g24 int = 3
// var g25 *int = &g24
// var g26 [3]int = [3]int{1, 2, 3}

// var g27 *int = g26 + 1
// var g28 *int = &g11[1].a
// var g30 = struct{ a struct{ a [3]int } }{{[3]int{1, 2, 3}}}

// var g31 *int = g30.a.a
// var g031 = g26
// var g032 = g30.a.a

// var g40 [2]struct{ a int } = [2]struct{ a int }{{1}, {3}}
// var g41 [3]struct {
// 	a int
// 	b int
// } = [3]struct {
// 	a int
// 	b int
// }{{1, 2}, {3, 4}, {5, 6}}

// var g01 = 3
// var g02 = 'a'
var g03 = [3]int{1, 2, 3}

// var g04 = [2][3]int{{1, 2, 3}, {4, 5, 6}}
// var g04_2 = g04[1]

// var g05 = "abc"
// var g06 = [2]string{"abc", "def"}

// var g07 = [2]struct {
// 	a int
// 	b int
// 	c int
// }{
// 	{1, 2, 3},
// 	{4, 5, 6},
// }

func main() {
	// comparing strings => unimplement yet.
	// assert(1, x034[0].b == "abc", "x034[0].b==\"abc\"")

	// var x01 = 3 + 1*2
	// assert(5, x01, "var x01=3;x01")
	// var x02 = [3]int{1, 2, 3}
	// assert(1, x02[0], "var x02=[3]int{1,2,3}; x02[0]")
	// assert(2, x02[1], "var x02=[3]int{1,2,3}; x02[1]")
	// assert(3, x02[2], "var x02=[3]int{1,2,3}; x02[2]")
	// var x03 = [2][3]int{{1, 2, 3}, {4, 5, 6}}
	// assert(1, x03[0][0], "var x03=[2][3]int{{1,2,3},{4,5,6}}; x03[0][1]")
	// assert(2, x03[0][1], "var x03=[2][3]int{{1,2,3},{4,5,6}}; x03[0][1]")
	// assert(3, x03[0][2], "var x03=[2][3]int{{1,2,3},{4,5,6}}; x03[0][2]")
	// assert(4, x03[1][0], "var x03=[2][3]int{{1,2,3},{4,5,6}}; x03[1][0]")
	// assert(5, x03[1][1], "var x03=[2][3]int{{1,2,3},{4,5,6}}; x03[1][1]")
	// assert(6, x03[1][2], "var x03=[2][3]int{{1,2,3},{4,5,6}}; x03[1][2]")

	// var x04 = "abc"
	// assert('a', x04[0], "var x04=\"abc\"; x04[0]")

	// var x05 = [2]string{"abc", "def"}
	// assert('a', x05[0][0], "var x05 =[2]string{\"abc\",\"def\"}; x05[0][0]")
	// assert('b', x05[0][1], "var x05 =[2]string{\"abc\",\"def\"}; x05[0][1]")
	// assert('c', x05[0][2], "var x05 =[2]string{\"abc\",\"def\"}; x05[0][2]")
	// assert('d', x05[1][0], "var x05 =[2]string{\"abc\",\"def\"}; x05[1][0]")
	// assert('e', x05[1][1], "var x05 =[2]string{\"abc\",\"def\"}; x05[1][1]")
	// assert('f', x05[1][2], "var x05 =[2]string{\"abc\",\"def\"}; x05[1][2]")

	// var x06 = struct {
	// 	a int
	// 	b int
	// 	c int
	// }{1, 2, 3}
	// assert(1, x06.a, "var x06 =struct {a int;b int;c int;}{1,2,3,};x06.a")
	// assert(2, x06.b, "var x06 =struct {a int;b int;c int;}{1,2,3,};x06.b")
	// assert(3, x06.c, "var x06 =struct {a int;b int;c int;}{1,2,3,};x06.c")
	// var x07 = [2]struct {
	// 	a int
	// 	b int
	// 	c int
	// }{
	// 	{1, 2, 3},
	// 	{4, 5, 6},
	// }
	// assert(1, x07[0].a, "var x07=[2]struct{a int;b int;c int;}{{1,2,3},{4,5,6},};x07[0].a")
	// assert(2, x07[0].b, "var x07=[2]struct{a int;b int;c int;}{{1,2,3},{4,5,6},};x07[0].b")
	// assert(3, x07[0].c, "var x07=[2]struct{a int;b int;c int;}{{1,2,3},{4,5,6},};x07[0].c")
	// assert(4, x07[1].a, "var x07=[2]struct{a int;b int;c int;}{{1,2,3},{4,5,6},};x07[1].a")
	// assert(5, x07[1].b, "var x07=[2]struct{a int;b int;c int;}{{1,2,3},{4,5,6},};x07[1].b")
	// assert(6, x07[1].c, "var x07=[2]struct{a int;b int;c int;}{{1,2,3},{4,5,6},};x07[1].c")

	// var x08 [3]int = [3]int{1, 2, 3}
	// var x08_2 = x08
	// assert(1, x08_2[0], "var x08 [3]int=[3]int{1,2,3};var x08_2=x08;x08_2[0]")
	// assert(2, x08_2[1], "var x08 [3]int=[3]int{1,2,3};var x08_2=x08;x08_2[1]")
	// assert(3, x08_2[2], "var x08 [3]int=[3]int{1,2,3};var x08_2=x08;x08_2[2]")

	// x09 := 3
	// assert(3, x09, "x09:=3;x09")

	// x010 := [3]int{1, 2, 3}
	// assert(1, x010[0], "x010:=[3]int{1,2,3};x010[0]")
	// assert(2, x010[1], "x010:=[3]int{1,2,3};x010[1]")
	// assert(3, x010[2], "x010:=[3]int{1,2,3};x010[2]")

	// x011 := [2]struct {
	// 	a int
	// 	b string
	// }{{1, "abc"}, {2, "def"}}
	// assert(1, x011[0].a, "x011 := [2]struct{a int;b string}{{1,\"abc\"},{2,\"def\"}};x011[0].a")
	// assert('a', x011[0].b[0], "x011 := [2]struct{a int;b string}{{1,\"abc\"},{2,\"def\"}};x011[0].b[0]")
	// assert('b', x011[0].b[1], "x011 := [2]struct{a int;b string}{{1,\"abc\"},{2,\"def\"}};x011[0].b[1]")
	// assert('c', x011[0].b[2], "x011 := [2]struct{a int;b string}{{1,\"abc\"},{2,\"def\"}};x011[0].b[2]")
	// assert(2, x011[1].a, "x011 := [2]struct{a int;b string}{{1,\"abc\"},{2,\"def\"}};x011[1].a")
	// assert('d', x011[1].b[0], "x011 := [2]struct{a int;b string}{{1,\"abc\"},{2,\"def\"}};x011[0].b[0]")
	// assert('e', x011[1].b[1], "x011 := [2]struct{a int;b string}{{1,\"abc\"},{2,\"def\"}};x011[1].b[1]")
	// assert('f', x011[1].b[2], "x011 := [2]struct{a int;b string}{{1,\"abc\"},{2,\"def\"}};x011[2].b[2]")

	// assert(3, g01, "var g01=3;g01")
	// assert('a', g02, "var g02='a';g02")

	assert(1, g03[0], "var g03=[3]int{1,2,3}; g03[0]")
	assert(2, g03[1], "var g03=[3]int{1,2,3}; g03[1]")
	assert(3, g03[2], "var g03=[3]int{1,2,3}; g03[2]")

	// assert(1, g04[0][0], "var g04=[2][3]int{{1,2,3},{4,5,6}}; g04[0][1]")
	// assert(2, g04[0][1], "var g04=[2][3]int{{1,2,3},{4,5,6}}; g04[0][1]")
	// assert(3, g04[0][2], "var g04=[2][3]int{{1,2,3},{4,5,6}}; g04[0][2]")
	// assert(4, g04[1][0], "var g04=[2][3]int{{1,2,3},{4,5,6}}; g04[1][0]")
	// assert(5, g04[1][1], "var g04=[2][3]int{{1,2,3},{4,5,6}}; g04[1][1]")
	// assert(6, g04[1][2], "var g04=[2][3]int{{1,2,3},{4,5,6}}; g04[1][2]")

	// assert(4, g04_2[0], "var g04_2[0]=g04[1]; g04_2[0]")
	// assert(5, g04_2[1], "var g04_2[0]=g04[1]; g04_2[1]")
	// assert(6, g04_2[2], "var g04_2[0]=g04[1]; g04_2[2]")

	// g04_2[2]=100
	// assert(100, g04_2[2], "var g04_2[0]=g04[1]; g04_2[2]")
	// assert(6, g04[1][2], "var g04=[2][3]int{{1,2,3},{4,5,6}}; g04[1][2]")

	// assert('a', g05[0], "var g05=\"abc\"; g05[0]")
	// assert('b', g05[1], "var g05=\"abc\"; g05[1]")
	// assert('c', g05[2], "var g05=\"abc\"; g05[2]")

	// assert('a', g06[0][0], "var g06 =[2]string{\"abc\",\"def\"}; g06[0][0]")
	// assert('b', g06[0][1], "var g06 =[2]string{\"abc\",\"def\"}; g06[0][1]")
	// assert('c', g06[0][2], "var g06 =[2]string{\"abc\",\"def\"}; g06[0][2]")
	// assert('d', g06[1][0], "var g06 =[2]string{\"abc\",\"def\"}; g06[1][0]")
	// assert('e', g06[1][1], "var g06 =[2]string{\"abc\",\"def\"}; g06[1][1]")
	// assert('f', g06[1][2], "var g06 =[2]string{\"abc\",\"def\"}; g06[1][2]")

	// assert(1, g07[0].a, "var g07=[2]struct{a int;b int;c int;}{{1,2,3},{4,5,6},};g07[0].a")
	// assert(2, g07[0].b, "var g07=[2]struct{a int;b int;c int;}{{1,2,3},{4,5,6},};g07[0].b")
	// assert(3, g07[0].c, "var g07=[2]struct{a int;b int;c int;}{{1,2,3},{4,5,6},};g07[0].c")
	// assert(4, g07[1].a, "var g07=[2]struct{a int;b int;c int;}{{1,2,3},{4,5,6},};g07[1].a")
	// assert(5, g07[1].b, "var g07=[2]struct{a int;b int;c int;}{{1,2,3},{4,5,6},};g07[1].b")
	// assert(6, g07[1].c, "var g07=[2]struct{a int;b int;c int;}{{1,2,3},{4,5,6},};g07[1].c")

	// var x1 [3]int = [3]int{1, 2, 3}
	// assert(1, x1[0], "var x1 [3]int=[3]int{1,2,3}; x1[0]")
	// assert(2, x1[1], "var x1 [3]int=[3]int{1,2,3}; x1[1]")
	// assert(3, x1[2], "var x1 [3]int=[3]int{1,2,3}; x1[2]")

	// var x2 [2][3]int = [2][3]int{{1, 2, 3}, {4, 5, 6}}
	// assert(2, x2[0][1], "var x2 [2][3]int=[2][3]int{{1,2,3},{4,5,6}}; x2[0][1]")
	// assert(4, x2[1][0], "var x2 [2][3]int=[2][3]int{{1,2,3},{4,5,6}}; x2[1][0]")
	// assert(6, x2[1][2], "var x2 [2][3]int=[2][3]int{{1,2,3},{4,5,6}}; x2[1][2]")

	// var x3 [3]int = [3]int{}
	// assert(0, x3[0], "var x3 [3]int=[3]int{}; x3[0]")
	// assert(0, x3[1], "var x3 [3]int=[3]int{}; x3[1]")
	// assert(0, x3[2], "var x3 [3]int=[3]int{}; x3[2]")

	// var x4 [2][3]int = [2][3]int{{1, 2}}
	// assert(2, x4[0][1], "var x4 [2][3]int=[2][3]int{{1,2}}; x4[0][1]")
	// assert(0, x4[1][0], "var x4 [2][3]int=[2][3]int{{1,2}}; x4[1][0]")
	// assert(0, x4[1][2], "var x4 [2][3]int=[2][3]int{{1,2}}; x4[1][2]")
	// var x5 [4]byte = "abc"
	// assert('a', x5[0], "var x5 [4]byte=\"abc\"; x5[0]")
	// assert('c', x5[2], "var x5 [4]byte=\"abc\"; x5[2]")
	// assert(0, x5[3], "var x5 [4]byte=\"abc\"; x5[3]")

	// var x6 string = "def"
	// assert('d', x6[0], "var x6 string=\"def\"; x6[0]")
	// assert('f', x6[2], "var x6 string=\"abc\"; x6[2]")
	// assert(0, x6[3], "var x6 string=\"abc\"; x6[3]")

	// var x7 [2][4]byte = [2][4]byte{"abc", "def"}
	// assert('a', x7[0][0], "var x7 [2][4]byte={\"abc\",\"def\"}; x7[0][0]")
	// assert(0, x7[0][3], "var x7 [2][4]byte={\"abc\",\"def\"}; x7[0][3]")
	// assert('d', x7[1][0], "var x7 [2][4]byte={\"abc\",\"def\"}; x7[1][0]")
	// assert('f', x7[1][2], "var x7 [2][4]byte={\"abc\",\"def\"}; x7[1][2]")

	// var x8 [2]string = [2]string{"abc", "def"}
	// assert('a', x8[0][0], "var x8 [2]string=[2]string{\"abc\",\"def\"}; x8[0][0]")
	// assert(0, x8[0][3], "var x8 [2]string=[2]string{\"abc\",\"def\"}; x8[0][3]")
	// assert('d', x8[1][0], "var x8 [2]string=[2]string{\"abc\",\"def\"}; x8[1][0]")
	// assert('f', x8[1][2], "var x8 [2]string=[2]string{\"abc\",\"def\"}; x8[1][2]")

	// // assert(4, ({ int x[]={1,2,3,4}; x[3]; }));
	// // assert(16, ({ int x[]={1,2,3,4}; sizeof(x); }));
	// // assert(4, ({ char x[]="foo"; sizeof(x); }));

	// type T9 string
	// var x9 T9 = "foo"
	// var y9 T9 = "x"
	// assert(8, Sizeof(x9), "type T9 string; var x9 T9=\"foo\"; var y9 T9=\"x\"; Sizeof(x9)")
	// assert(8, Sizeof(y9), "type T9 string; var x9 T9=\"foo\"; var y9 T9=\"x\"; Sizeof(y9)")
	// var x10 T9 = "foo"
	// var y10 T9 = "x"
	// assert(8, Sizeof(x10), "type T9 string; var x10 T9=\"foo\"; var y10 T9=\"x\"; Sizeof(x10)")
	// assert(8, Sizeof(y10), "type T9 string; var x10 T9=\"foo\"; var y10 T9=\"x\"; Sizeof(y10)")

	// // assert(4, ({ typedef char T[]; T x="foo"; T y="x"; sizeof(x); }));
	// // assert(2, ({ typedef char T[]; T x="foo"; T y="x"; sizeof(y); }));
	// // assert(2, ({ typedef char T[]; T x="x"; T y="foo"; sizeof(x); }));
	// // assert(4, ({ typedef char T[]; T x="x"; T y="foo"; sizeof(y); }));

	// type T11 struct {
	// 	a int
	// 	b int
	// 	c int
	// }
	// var x11 T11 = T11{1, 2, 3}
	// assert(1, x11.a, "type T11 struct {a int; b int; c int;}; var x11 T11=T11{1,2,3}; x11.a")
	// assert(2, x11.b, "type T11 struct {a int; b int; c int;}; var x11 T11=T11{1,2,3}; x11.b")
	// assert(3, x11.c, "type T11 struct {a int; b int; c int;}; var x11 T11=T11{1,2,3}; x11.c")
	// var x12 T11 = T11{1}
	// assert(1, x12.a, "type T11 struct {a int; b int; c int;}; var x12 T11={1}; x12.a")
	// assert(0, x12.b, "type T11 struct {a int; b int; c int;}; var x12 T11={1}; x12.b")
	// assert(0, x12.c, "type T11 struct {a int; b int; c int;}; var x12 T11={1}; x12.c")
	// type T13 struct {
	// 	a int
	// 	b int
	// }
	// var x13 [2]T13 = [2]T13{T13{1, 2}, T13{3, 4}}
	// assert(1, x13[0].a, "type T13 struct {a int; b int;}; var x [2]T13=[2]T13{{1,2},{3,4}}; x13[0].a")
	// assert(2, x13[0].b, "type T13 struct {a int; b int;}; var x [2]T13=[2]T13{{1,2},{3,4}}; x13[0].b")
	// assert(3, x13[1].a, "type T13 struct {a int; b int;}; var x [2]T13=[2]T13{{1,2},{3,4}}; x13[1].a")
	// assert(4, x13[1].b, "type T13 struct {a int; b int;}; var x [2]T13=[2]T13{{1,2},{3,4}}; x13[1].b")
	// type T14 struct {
	// 	a int
	// 	b int
	// }
	// var x14 [2]T14 = [2]T14{T14{1, 2}}
	// assert(0, x14[1].b, "type T14 struct {a int; b int;}; var x14 [2]T14=[2]T14{{1,2}}; x14[1].b")
	// type T15 struct {
	// 	a int
	// 	b int
	// }
	// var x15 T15 = T15{}
	// assert(0, x15.a, "type T15 struct {a int; b int;}; var x15 T15=T15{}; x15.a")
	// assert(0, x15.b, "type T15 struct {a int; b int;}; var x15 T15=T15{}; x15.b")
	// type T16 struct {
	// 	a int
	// 	b int
	// 	c int
	// 	d int
	// 	e int
	// 	f int
	// }
	// var x16 T16 = T16{1, 2, 3, 4, 5, 6}
	// var y16 T16
	// y16 = x16
	// assert(5, y16.e, "type T16 struct {a int;b int;c int;d int;e int;f int;}; var x16 T16=T16{1,2,3,4,5,6};var y16 T16; y16=x16; y16.e")
	// type T17 struct {a int;b int;}; var x17 T17=T17{1,2};var y17 T17;var z17 T17; z17=y17=x17;
	// assert(2, z17.b, "type T17 struct {a int;b int;}; var x17 T17=T17{1,2};var y17 T17,var z17 T17; z=y=x; z.b")
	// type T18 struct {a int;b int;}; var x18 T18=T18{1,2};var y18 T18=x18;
	// assert(1, y18.a, "type T18 struct {a int;b int;}; var x18 T18=T18{1,2};var y T18=x18; y18.a");

	// assert(3, g3, "g3")
	// assert(4, g4, "g4")
	// assert(5, g5, "g5")
	// assert(6, g6, "g6")

	// assert(0, g9[0], "g9[0]")
	// assert(1, g9[1], "g9[1]")
	// assert(2, g9[2], "g9[2]")

	// assert(1, g11[0].a, "g11[0].a")
	// assert(2, g11[0].b, "g11[0].b")
	// assert(3, g11[1].a, "g11[1].a")
	// assert(4, g11[1].b, "g11[1].b")

	// assert(1, g12[0].a[0], "g12[0].a[0]")
	// assert(2, g12[0].a[1], "g12[0].a[1]")
	// assert(0, g12[1].a[0], "g12[1].a[0]")
	// assert(0, g12[1].a[1], "g12[1].a[1]")

	// assert(8, Sizeof(g17), "Sizeof(g17)")

	// assert('f', g17[0], "g17[0]")
	// assert('o', g17[1], "g17[1]")
	// assert('o', g17[2], "g17[2]")
	// assert('b', g17[3], "g17[3]")
	// assert('a', g17[4], "g17[4]")
	// assert('r', g17[5], "g17[5]")

	// assert('f', g17_2[0], "g17_2[0]")
	// assert('o', g17_2[1], "g17_2[1]")
	// assert('o', g17_2[2], "g17_2[2]")
	// assert('b', g17_2[3], "g17_2[3]")
	// assert('a', g17_2[4], "g17_2[4]")
	// assert('r', g17_2[5], "g17_2[5]")

	// g18 = "foo"
	// assert('f', g18[0], "g18[0]")
	// assert('o', g18[1], "g18[1]")
	// assert('o', g18[2], "g18[2]")

	// assert(3, g24, "g24")
	// assert(3, *g25, "*g25")
	// assert(2, *g27, "*g27")
	// assert(3, *g28, "*g28")

	// assert(1, g31[0], "g31[0]")
	// assert(2, g31[1], "g31[1]");
	// assert(3, g31[2], "g31[2]");

	// assert(1, g031[0], "g031[0]")
	// assert(2, g031[1], "g031[1]")
	// assert(3, g031[2], "g031[2]")

	// assert(1, g032[0], "g032[0]")
	// assert(2, g032[1], "g032[1]")
	// assert(3, g032[2], "g032[2]")

	// assert(1, g40[0].a, "g40[0].a")
	// assert(3, g40[1].a, "g40[1].a")

	// assert(1, g41[0].a, "g41[0].a")
	// assert(2, g41[0].b, "g41[0].b")
	// assert(3, g41[1].a, "g41[1].a")
	// assert(4, g41[1].b, "g41[1].b")
	// assert(5, g41[2].a, "g41[2].a")
	// assert(6, g41[2].b, "g41[2].b")
	// var a [3]int=[3]int{1,2,3,};
	// assert(3, a[2], "var a [3]int=[3]int{1,2,3,}; a[2]");
	// var x19 struct {a int;b int;c int;}={1,2,3,};
	// assert(1, x19.a, "var x19 struct {a int;b int;c int;}={1,2,3,}; x19.a");

	println("OK")
}
