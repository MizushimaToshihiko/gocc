package test

func main() {
	var i int = 0
	switch 3 {
	case 5 - 2 + 0*3:
		i++
	}
	assert(1, i, "var i int =0; switch(3) { case 5-2+0*3: i++; }")
	// var x [1 + 1]int
	// assert(8, Sizeof(x), "var x [1+1]int; Sizeof(x)")
	// assert(6, ({ char x[8-2]; Sizeof(x); }));
	// assert(6, ({ char x[2*3]; Sizeof(x); }));
	// assert(3, ({ char x[12/4]; Sizeof(x); }));
	// assert(2, ({ char x[12%10]; Sizeof(x); }));
	// assert(0b100, ({ char x[0b110&0b101]; Sizeof(x); }));
	// assert(0b111, ({ char x[0b110|0b101]; Sizeof(x); }));
	// assert(0b110, ({ char x[0b111^0b001]; Sizeof(x); }));
	// assert(4, ({ char x[1<<2]; Sizeof(x); }));
	// assert(2, ({ char x[4>>1]; Sizeof(x); }));
	// assert(2, ({ char x[(1==1)+1]; Sizeof(x); }));
	// assert(1, ({ char x[(1!=1)+1]; Sizeof(x); }));
	// assert(1, ({ char x[(1<1)+1]; Sizeof(x); }));
	// assert(2, ({ char x[(1<=1)+1]; Sizeof(x); }));
	// assert(2, ({ char x[1?2:3]; Sizeof(x); }));
	// assert(3, ({ char x[0?2:3]; Sizeof(x); }));
	// assert(3, ({ char x[(1,3)]; Sizeof(x); }));
	// assert(2, ({ char x[!0+1]; Sizeof(x); }));
	// assert(1, ({ char x[!1+1]; Sizeof(x); }));
	// assert(2, ({ char x[~-3]; Sizeof(x); }));
	// assert(2, ({ char x[(5||6)+1]; Sizeof(x); }));
	// assert(1, ({ char x[(0||0)+1]; Sizeof(x); }));
	// assert(2, ({ char x[(1&&1)+1]; Sizeof(x); }));
	// assert(1, ({ char x[(1&&0)+1]; Sizeof(x); }));
	// assert(3, ({ char x[(int)3]; Sizeof(x); }));
	// assert(15, ({ char x[(char)0xffffff0f]; Sizeof(x); }));
	// assert(0x10f, ({ char x[(short)0xffff010f]; Sizeof(x); }));
	// assert(4, ({ char x[(int)0xfffffffffff+5]; Sizeof(x); }));
	// assert(8, ({ char x[(int*)0+2]; Sizeof(x); }));
	// assert(12, ({ char x[(int*)16-1]; Sizeof(x); }));
	// assert(3, ({ char x[(int*)16-(int*)4]; Sizeof(x); }));

	printf("OK\n")
}
