package test

func assert(want int, act int, code string)
func println(format ...string)

func main() {

	type x struct {
		a int
		b int
	}
	var y1 = &x{1, 2}
	assert(1, y1.a, "y1.a")
	var y2 = &[2]x{{1, 2}, {3, 4}}
	assert(1, y2[0].a, "y2[0].a")
	// assert(3, *&x1, "var x int=3; *&x1")
	// var x2 int = 3
	// var y2 *int = &x2
	// var z2 **int = &y2
	// assert(3, **z2, "var x2 int=3; var y2 *int=&x2; var z2 **int=&y2; **z2")
	// var x3 int = 3
	// var y3 int = 5
	// assert(5, *(&x3 + 1), "var x3 int=3; var y3 int=5; *(&x3+1)")       // Not supported in Go.
	// assert(3, *(&y3 - 1), "var x3 int=3; var y3 int=5; *(&y3-1)")       // Not supported in Go.
	// assert(5, *(&x3 - (-1)), "var x3 int=3; var y3 int=5; *(&x3-(-1))") // Not supported in Go.
	// var x4 int = 3
	// var y4 *int = &x4
	// *y4 = 5
	// assert(5, x4, "var x4 int=3; var y4 *int=&x4; *y4=5; x4")
	// var x5 int = 3
	// var y5 int = 5
	// *(&x5 + 1) = 7
	// assert(7, y5, "var x5 int=3; var y5 int=5; *(&x5+1)=7; y5")
	// var x6 int = 3
	// var y6 int = 5
	// *(&y6 - 2 + 1) = 7
	// assert(7, x6, "var x6 int=3; var y6 int=5; *(&y6-2+1)=7; x6")
	// var x7 int = 3
	// assert(5, (&x7+2)-&x7+3, "var x7 int=3; (&x+2)-&x+3")
	// var x8 [2]int
	// var y8 *int = &x8
	// *y8 = 3
	// assert(3, *x8, "var x8 [2]int; var y8 *int=&x8; *y8=3; *x8")
	// var x9 [3]int
	// *x9 = 3
	// *(x9 + 1) = 4
	// *(x9 + 2) = 5
	// assert(3, *x9, "var x9 [3]int; *x9=3; *(x9+1)=4; *(x9+2)=5; *x9")
	// assert(4, *(x9 + 1), "var x9 [3]int; *x9=3; *(x9+1)=4; *(x9+2)=5; *(x9+1)")
	// assert(5, *(x9 + 2), "var x9 [3]int; *x9=3; *(x9+1)=4; *(x9+2)=5; *(x9+2)")
	// var x10 [2][3]int
	// var y10 *int = x10
	// *y10 = 0
	// assert(0, **x10, "var x10 [2][3]int; var y10 *int=x10; *y10=0; **x10")
	// *(y10 + 1) = 1
	// assert(1, *(*x10 + 1), "var x10 [2][3]int; var y10 *int=x10; *(y10+1)=1; *(*x10+1)")
	// *(y10 + 2) = 2
	// assert(2, *(*x10 + 2), "var x10 [2][3]int; var y10 *int=x10; *(y10+2)=2; *(*x10+2)")
	// *(y10 + 3) = 3
	// assert(3, **(x10 + 1), "var x10 [2][3]int; var y10 *int=x10; *(y10+3)=3; **(x10+1)")
	// *(y10 + 4) = 4
	// assert(4, *(*(x10 + 1) + 1), "var x10 [2][3]int; var y10 *int=x10; *(y10+4)=4; *(*(x10+1)+1)")
	// *(y10 + 5) = 5
	// assert(5, *(*(x10 + 1) + 2), "var x10 [2][3]int; var y10 *int=x10; *(y10+5)=5; *(*(x10+1)+2)")
	// var x11 [3]int
	// *x11 = 3
	// x11[1] = 4
	// x11[2] = 5
	// assert(3, *x11, "var x11 [3]int; *x11=3; x11[1]=4; x11[2]=5; *x11")
	// assert(4, *(x11 + 1), "var x11 [3]int; *x11=3; x11[1]=4; x11[2]=5; *(x11+1)")
	// assert(5, *(x11 + 2), "var x11 [3]int; *x11=3; x11[1]=4; x11[2]=5; *(x11+2)")
	// 2[x11] = 5
	// assert(5, *(x11 + 2), "var x11 [3]int; *x11=3; x11[1]=4; 2[x11]=5; *(x11+2)")
	// var x12 [2][3]int
	// var y12 *int = x12
	// y12[0] = 0
	// assert(0, x12[0][0], "var x12 [2][3]int; var y12 *int=x12; y12[0]=0; x12[0][0]")
	// y12[1] = 1
	// assert(1, x12[0][1], "var x12 [2][3]int; var y12 *int=x12; y12[1]=1; x12[0][1]")
	// y12[2] = 2
	// assert(2, x12[0][2], "var x12 [2][3]int; var y12 *int=x12; y12[2]=2; x12[0][2]")
	// y12[3] = 3
	// assert(3, x12[1][0], "var x12 [2][3]int; var y12 *int=x12; y12[3]=3; x12[1][0]")
	// y12[4] = 4
	// assert(4, x12[1][1], "var x12 [2][3]int; var y12 *int=x12; y12[4]=4; x12[1][1]")
	// y12[5] = 5
	// assert(5, x12[1][2], "var x12 [2][3]int; var y12 *int=x12; y12[5]=5; x12[1][2]")

	println("OK")
}
