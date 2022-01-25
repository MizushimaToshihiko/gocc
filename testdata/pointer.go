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
	// assert(5, ({ int x=3; int *y=&x; *y=5; x; }));
	// assert(7, ({ int x=3; int y=5; *(&x+1)=7; y; }));
	// assert(7, ({ int x=3; int y=5; *(&y-2+1)=7; x; }));
	// assert(5, ({ int x=3; (&x+2)-&x+3; }));
	// assert(8, ({ int x, y; x=3; y=5; x+y; }));
	// assert(8, ({ int x=3, y=5; x+y; }));

	// assert(3, ({ int x[2]; int *y=&x; *y=3; *x; }));

	// assert(3, ({ int x[3]; *x=3; *(x+1)=4; *(x+2)=5; *x; }));
	// assert(4, ({ int x[3]; *x=3; *(x+1)=4; *(x+2)=5; *(x+1); }));
	// assert(5, ({ int x[3]; *x=3; *(x+1)=4; *(x+2)=5; *(x+2); }));

	// assert(0, ({ int x[2][3]; int *y=x; *y=0; **x; }));
	// assert(1, ({ int x[2][3]; int *y=x; *(y+1)=1; *(*x+1); }));
	// assert(2, ({ int x[2][3]; int *y=x; *(y+2)=2; *(*x+2); }));
	// assert(3, ({ int x[2][3]; int *y=x; *(y+3)=3; **(x+1); }));
	// assert(4, ({ int x[2][3]; int *y=x; *(y+4)=4; *(*(x+1)+1); }));
	// assert(5, ({ int x[2][3]; int *y=x; *(y+5)=5; *(*(x+1)+2); }));

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
