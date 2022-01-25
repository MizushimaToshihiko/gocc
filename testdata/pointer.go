package test

func main() {
	var x1 int = 3
	assert(3, *&x1, "var x int=3; *&x1")
	var x2 int = 3
	var y2 *int = &x2
	var z2 **int = &y2
	assert(3, **z2, "var x2 int=3; var y2 *int=&x2; var z2 **int=&y2; **z2")
	var x3 int = 3
	var y3 int = 5
	// Not supported in Go.
	assert(5, *(&x3 + 1), "var x3 int=3; var y3 int=5; *(&x3+1)")
	assert(3, *(&y3 - 1), "var x3 int=3; var y3 int=5; *(&y3-1)")
	assert(5, *(&x3 - (-1)), "var x3 int=3; var y3 int=5; *(&x3-(-1))")
	var x4 int = 3
	var y4 *int = &x4
	*y4 = 5
	assert(5, x4, "var x4 int=3; var y4 *int=&x4; *y4=5; x4")
	var x5 int = 3
	var y5 int = 5
	*(&x5 + 1) = 7
	assert(7, y5, "var x5 int=3; var y5 int=5; *(&x5+1)=7; y5")
	var x6 int = 3
	var y6 int = 5
	*(&y6 - 2 + 1) = 7
	assert(7, x6, "var x6 int=3; var y6 int=5; *(&y6-2+1)=7; x6")
	var x7 int = 3
	assert(5, (&x7+2)-&x7+3, "var x7 int=3; (&x+2)-&x+3")
	var x8 [2]int
	var y8 *int = &x8
	*y8 = 3
	assert(3, *x8, "var x8 [2]int; var y8 *int=&x8; *y8=3; *x8")
	var x9 [3]int
	*x9 = 3
	*(x9 + 1) = 4
	*(x9 + 2) = 5
	assert(3, *x9, "var x9 [3]int; *x9=3; *(x9+1)=4; *(x9+2)=5; *x9")
	assert(4, *(x9 + 1), "var x9 [3]int; *x9=3; *(x9+1)=4; *(x9+2)=5; *(x9+1)")
	assert(5, *(x9 + 2), "var x9 [3]int; *x9=3; *(x9+1)=4; *(x9+2)=5; *(x9+2)")
	var x10 [2][3]int
	var y10 *int = x10
	*y10 = 0
	assert(0, **x10, "var x10 [2][3]int; var y10 *int=x10; *y10=0; **x10")
	*(y10 + 1) = 1
	assert(1, *(*x10 + 1), "var x10 [2][3]int; var y10 *int=x10; *(y10+1)=1; *(*x10+1)")
	*(y10 + 2) = 2
	assert(2, *(*x10 + 2), "var x10 [2][3]int; var y10 *int=x10; *(y10+2)=2; *(*x10+2)")
	*(y10 + 3) = 3
	assert(3, **(x10 + 1), "var x10 [2][3]int; var y10 *int=x10; *(y10+3)=3; **(x10+1)")
	*(y10 + 4) = 4
	assert(4, *(*(x10 + 1) + 1), "var x10 [2][3]int; var y10 *int=x10; *(y10+4)=4; *(*(x10+1)+1)")
	*(y10 + 5) = 5
	assert(5, *(*(x10 + 1) + 2), "var x10 [2][3]int; var y10 *int=x10; *(y10+5)=5; *(*(x10+1)+2)")

	// assert(3, ({ int x[3]; *x=3; x[1]=4; x[2]=5; *x; }));
	// assert(4, ({ int x[3]; *x=3; x[1]=4; x[2]=5; *(x+1); }));
	// assert(5, ({ int x[3]; *x=3; x[1]=4; x[2]=5; *(x+2); }));
	// assert(5, ({ int x[3]; *x=3; x[1]=4; x[2]=5; *(x+2); }));
	// assert(5, ({ int x[3]; *x=3; x[1]=4; 2[x]=5; *(x+2); }));

	// assert(0, ({ int x[2][3]; int *y=x; y[0]=0; x[0][0]; }));
	// assert(1, ({ int x[2][3]; int *y=x; y[1]=1; x[0][1]; }));
	// assert(2, ({ int x[2][3]; int *y=x; y[2]=2; x[0][2]; }));
	// assert(3, ({ int x[2][3]; int *y=x; y[3]=3; x[1][0]; }));
	// assert(4, ({ int x[2][3]; int *y=x; y[4]=4; x[1][1]; }));
	// assert(5, ({ int x[2][3]; int *y=x; y[5]=5; x[1][2]; }));

	printf("OK\n")
}
