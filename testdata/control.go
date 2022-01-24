package test

func main() {
	var x int
	if 0 {
		x = 2
	} else {
		x = 3
	}
	assert(3, x, "if (0) x=2; else x=3; x;")
	if 1 - 1 {
		x = 2
	} else {
		x = 3
	}
	assert(3, x, "if (1-1) {x=2;} else {x=3;}")
	if 1 {
		x = 2
	} else {
		x = 3
	}
	assert(2, x, "if 1 {x=2;} else {x=3;}")
	if 2 - 1 {
		x = 2
	} else {
		x = 3
	}
	assert(2, x, "if 2-1 {x=2;} else {x=3;}")

	var i int = 0
	var j int = 0
	for i = 0; i <= 10; i = i + 1 {
		j = i + j
	}
	assert(55, j, "for i=0; i<=10; i=i+1 {j=i+j;}")
	for i < 10 {
		i = i + 1
	}
	assert(10, i, "for i<10 {i=i+1;}")

	i = 1
	{
		i = 2
	}
	i = 3
	assert(3, i, "i = 1;{i = 2};i = 3")
	// assert(5, ({ ;;; 5; }));

	// assert(10, ({ int i=0; while(i<10) i=i+1; i; }));
	// assert(55, ({ int i=0; int j=0; while(i<=10) {j=i+j; i=i+1;} j; }));

	// assert(3, (1,2,3));
	// assert(5, ({ int i=2, j=3; (i=5,j)=6; i; }));
	// assert(6, ({ int i=2, j=3; (i=5,j)=6; j; }));

	// assert(55, ({ int j=0; for (int i=0; i<=10; i=i+1) j=j+i; j; }));
	// assert(3, ({ int i=3; int j=0; for (int i=0; i<=10; i=i+1) j=j+i; i; }));

	// assert(1, 0||1);
	// assert(1, 0||(2-2)||5);
	// assert(0, 0||0);
	// assert(0, 0||(2-2));

	// assert(0, 0&&1);
	// assert(0, (2-2)&&5);
	// assert(1, 1&&5);

	// assert(3, ({ int i=0; goto a; a: i++; b: i++; c: i++; i; }));
	// assert(2, ({ int i=0; goto e; d: i++; e: i++; f: i++; i; }));
	// assert(1, ({ int i=0; goto i; g: i++; h: i++; i: i++; i; }));

	// assert(1, ({ typedef int foo; goto foo; foo:; 1; }));

	// assert(3, ({ int i=0; for(;i<10;i++) { if (i == 3) break; } i; }));
	// assert(4, ({ int i=0; while (1) { if (i++ == 3) break; } i; }));
	// assert(3, ({ int i=0; for(;i<10;i++) { for (;;) break; if (i == 3) break; } i; }));
	// assert(4, ({ int i=0; while (1) { while(1) break; if (i++ == 3) break; } i; }));

	// assert(10, ({ int i=0; int j=0; for (;i<10;i++) { if (i>5) continue; j++; } i; }));
	// assert(6, ({ int i=0; int j=0; for (;i<10;i++) { if (i>5) continue; j++; } j; }));
	// assert(10, ({ int i=0; int j=0; for(;!i;) { for (;j!=10;j++) continue; break; } j; }));
	// assert(11, ({ int i=0; int j=0; while (i++<10) { if (i>5) continue; j++; } i; }));
	// assert(5, ({ int i=0; int j=0; while (i++<10) { if (i>5) continue; j++; } j; }));
	// assert(11, ({ int i=0; int j=0; while(!i) { while (j++!=10) continue; break; } j; }));

	// assert(5, ({ int i=0; switch(0) { case 0:i=5;break; case 1:i=6;break; case 2:i=7;break; } i; }));
	// assert(6, ({ int i=0; switch(1) { case 0:i=5;break; case 1:i=6;break; case 2:i=7;break; } i; }));
	// assert(7, ({ int i=0; switch(2) { case 0:i=5;break; case 1:i=6;break; case 2:i=7;break; } i; }));
	// assert(0, ({ int i=0; switch(3) { case 0:i=5;break; case 1:i=6;break; case 2:i=7;break; } i; }));
	// assert(5, ({ int i=0; switch(0) { case 0:i=5;break; default:i=7; } i; }));
	// assert(7, ({ int i=0; switch(1) { case 0:i=5;break; default:i=7; } i; }));
	// assert(2, ({ int i=0; switch(1) { case 0: 0; case 1: 0; case 2: 0; i=2; } i; }));
	// assert(0, ({ int i=0; switch(3) { case 0: 0; case 1: 0; case 2: 0; i=2; } i; }));

	// assert(3, ({ int i=0; switch(-1) { case 0xffffffff: i=3; break; } i; }));

}
