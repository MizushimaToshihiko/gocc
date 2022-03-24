package test_control

func assert(want int, act int, code string)
func println(format string)

func switchFn(i int) int {
	switch i {
	case 0, 3, 4:
		return 5
	case 1, 2:
		return 6
	case 5, 6:
		return 100
	default:
		return 10
	}
}

func main() {
	var x1 int
	if 0 {
		x1 = 2
	} else {
		x1 = 3
	}
	assert(3, x1, "if (0) x=2; else x=3; x1;")
	if 1 - 1 {
		x1 = 2
	} else {
		x1 = 3
	}
	assert(3, x1, "if (1-1) {x=2;} else {x=3;}; x1;")
	if 1 {
		x1 = 2
	} else {
		x1 = 3
	}
	assert(2, x1, "if 1 {x=2;} else {x=3;}; x1;")
	if 2 - 1 {
		x1 = 2
	} else {
		x1 = 3
	}
	assert(2, x1, "if 2-1 {x=2;} else {x=3;}; x1;")

	var i int = 0
	var j int = 0
	for i = 0; i <= 10; i = i + 1 {
		j = i + j
	}
	assert(55, j, "var i int=0;var j int=0;for i=0; i<=10; i=i+1 {j=i+j;}")
	var j int = 0
	for i := 0; i <= 10; i = i + 1 {
		j = i + j
	}
	assert(55, j, "for i:=0; i<=10; i=i+1 {j=i+j;}")
	i = 0
	for i < 10 {
		i = i + 1
	}
	assert(10, i, "i=0; for i<10 {i=i+1;}")

	i = 1
	{
		i = 2
	}
	i = 3
	assert(3, i, "i = 1;{i = 2};i = 3")

	i = 0
	for i < 10 {
		i = i + 1
	}
	assert(10, i, "i=0; for i<10 {i=i+1;}")
	i = 0
	j = 0
	for i <= 10 {
		j = i + j
		i = i + 1
	}
	assert(55, j, "i=0; j=0; for i<=10 {j=i+j; i=i+1;} j")

	// assert(3, (1,2,3), "(1,2,3)")
	// i=2, j=3; (i=5,j)=6;
	// assert(5, i, "i=2, j=3; (i=5,j)=6; i")
	// i=2, j=3; (i=5,j)=6;
	// assert(6, j, "i=2, j=3; (i=5,j)=6; j")

	assert(1, 0 || 1, "0||1")
	assert(1, 0 || (2-2) || 5, "0||(2-2)||5")
	assert(0, 0 || 0, "0||0")
	assert(0, 0 || (2-2), "0||(2-2)")

	assert(0, 0 && 1, "0&&1")
	assert(0, (2-2) && 5, "(2-2)&&5")
	assert(1, 1 && 5, "1&&5")

	i = 0
	goto a
a:
	i++
b:
	i++
c:
	i++
	assert(3, i, "i=0; goto a; a: i++; b: i++; c: i++; i")
	i = 0
	goto e
d:
	i++
e:
	i++
f:
	i++
	assert(2, i, "i=0; goto e; d: i++; e: i++; f: i++; i")
	i = 0
	goto i
g:
	i++
h:
	i++
i:
	i++
	i
	assert(1, i, "i=0; goto i; g: i++; h: i++; i: i++; i")

	type foo int
	var x2 foo
	goto foo
	x2 = 2
foo:
	x2 = 1
	assert(1, x2, "type foo int; var x2 foo; goto foo; x2=2; foo:; x2=1;")

	i = 0
	for ; i < 10; i++ {
		if i == 3 {
			break
		}
	}
	assert(3, i, "i=0; for ;i<10;i++ { if i == 3 {break} } i")
	i = 0
	for 1 {
		i++
		if i >= 3 {
			i++
			break
		}
	}
	assert(4, i, "i=0; for 1 { if i == 3 {i++; break;}} i")
	i = 0
	for ; i < 10; i++ {
		for {
			break
		}
		if i == 3 {
			break
		}
	}
	assert(3, i, "i=0; for ;i<10;i++ { for ;; {break;}; if i == 3 {break;} } i")
	i = 0
	j = 0
	for ; i < 10; i++ {
		if i > 5 {
			continue
		}
		j++
	}
	assert(10, i, "i=0; j=0; for ;i<10;i++ { if i>5 {continue;}; j++; }; i")
	i = 0
	j = 0
	for ; i < 10; i++ {
		if i > 5 {
			continue
		}
		j++
	}
	assert(6, j, "i=0; j=0; for ;i<10;i++ { if i>5 {continue;}; j++; } j")
	i = 0
	j = 0
	for !i {
		for ; j != 10; j++ {
			continue
		}
		break
	}
	assert(10, j, "i=0; j=0; for ;!i; { for ;j!=10;j++ {continue;}; break; } j")
	i = 0
	j = 0
	for i < 10 {
		i++
		if i > 5 {
			continue
		}
		j++
	}
	assert(10, i, "i=0; j=0; for i<10 {i++; if i>5 {continue;}; j++; } i")
	i = 0
	j = 0
	for i < 10 {
		i++
		if i > 5 {
			continue
		}
		j++
	}
	assert(5, j, "i=0; j=0; for i<10 {i++; if i>5 {continue;}; j++; } j")
	i = 0
	j = 0
	for !i {
		for j != 10 {
			j++
			continue
		}
		break
	}
	assert(10, j, "i=0; j=0; for !i { for j!=10 {j++; continue;}; break; } j")

	assert(5, switchFn(0), "switchFn(0)")
	assert(5, switchFn(3), "switchFn(3)")
	assert(5, switchFn(4), "switchFn(4)")
	assert(6, switchFn(1), "switchFn(1)")
	assert(6, switchFn(2), "switchFn(2)")
	assert(100, switchFn(5), "switchFn(5)")
	assert(100, switchFn(6), "switchFn(6)")
	assert(10, switchFn(8), "switchFn(8)")
	assert(10, switchFn(9), "switchFn(9)")
	assert(10, switchFn(10), "switchFn(10)")
	assert(10, switchFn(11), "switchFn(11)")

	i = 0
	switch i {
	case 0, 3:
		i = 5
	case 1:
		i = 6
	case 2:
		i = 7
	}
	assert(5, i, "i=0; switch i { case 0,3:i=5; case 1:i=6; case 2:i=7; } i")
	i = 1
	switch i {
	case 0:
		i = 5
	case 1:
		i = 6
	case 2:
		i = 7
	}
	assert(6, i, "i=1; switch i { case 0:i=5; case 1:i=6; case 2:i=7; } i")
	i = 2
	switch i {
	case 0:
		i = 5
	case 1:
		i = 6
	case 2:
		i = 7
	}
	assert(7, i, "i=2; switch i { case 0:i=5; case 1:i=6; case 2:i=7; } i")
	i = 3
	switch i {
	case 0:
		i = 5
	case 1:
		i = 6
	case 2:
		i = 7
	}
	assert(3, i, "i=3; switch i { case 0:i=5; case 1:i=6; case 2:i=7; } i")
	i = 0
	switch i {
	case 0:
		i = 5
	default:
		i = 7
	}
	assert(5, i, "i=0; switch i { case 0:i=5;i; default:i=7; } i")
	i = 2
	switch i {
	case 0:
		i = 5
	default:
		i = 7
	}
	assert(7, i, "i=2; switch i { case 0:i=5;i; default:i=7; } i")
	i = 0
	switch -1 {
	case 0xffffffff:
		i = 3
	}
	assert(3, i, "i=0; switch(-1) { case 0xffffffff: i=3; }; i")

	assert(0, 0.0 && 0.0, "0.0 && 0.0");
	assert(0, 0.0 && 0.1, "0.0 && 0.1");
	assert(0, 0.3 && 0.0, "0.3 && 0.0");
	assert(1, 0.3 && 0.5, "0.3 && 0.5");
	assert(0, 0.0 || 0.0, "0.0 || 0.0");
	assert(1, 0.0 || 0.1, "0.0 || 0.1");
	assert(1, 0.3 || 0.0, "0.3 || 0.0");
	assert(1, 0.3 || 0.5, "0.3 || 0.5");
	var x2 int; if 0.0 {x2=3;}else{x2=5;};
	assert(5, x2, "var x int; if 0.0{x=3;}else{x=5;}; x2");
	var x3 int; if 0.1 {x3=3;}else{x3=5;};
	assert(3, x3, "var x3 int; if 0.1{x3=3;}else{x3=5;};x3");
	var x4=5; if 0.0{x4=3;};
	assert(5, x4, "var x4=5; if 0.0{x4=3;}; x4");
	var x5=5; if 0.1{x5=3;};
	assert(3, x5, "var x5=5; if 0.1{x5=3;};x5");
	i=10.0; j=0; for ;i!=0;i--,j++{};
	assert(10, j, "i=10.0; j=0; for ;i!=0;i--,j++{}; j;");
	i=10.0; j=0; for i!=0{i--;j++;};
	assert(10, j, "i=10.0; j=0; for i!=0{i--;j++}; j;");

	var x6, y6 = 1, 2
	var z6 int
	switch {
	case x6 < y6: z6=switchFn(x6)
	case x6 > y6: z6=switchFn(y6)
	case x6 == y6: z6=switchFn(x6+y6)
	}
	assert(6, z6, "z6")
	x6, y6 = 2, 0
	z6 = 0
	switch {
	case x6 < y6: z6=switchFn(x6)
	case x6 > y6: z6=switchFn(y6)
	case x6 == y6: z6=switchFn(x6+y6)
	}
	assert(5, z6, "z6")
	x6, y6 = 3, 3
	z6 = 0
	switch {
	case x6 < y6: z6=switchFn(x6)
	case x6 > y6: z6=switchFn(y6)
	case x6 == y6: z6=switchFn(x6+y6)
	}
	assert(100, z6, "z6")

	var x6, y6 float32 = 1.0, 2.0
	var z6 int
	switch {
	case x6 < y6: z6=switchFn(x6)
	case x6 > y6: z6=switchFn(y6)
	case x6 == y6: z6=switchFn(x6+y6)
	}
	assert(6, z6, "z6")
	x6, y6 = 2.0, 0.0
	z6 = 0
	switch {
	case x6 < y6: z6=switchFn(x6)
	case x6 > y6: z6=switchFn(y6)
	case x6 == y6: z6=switchFn(x6+y6)
	}
	assert(5, z6, "z6")
	x6, y6 = 3.0, 3.0
	z6 = 0
	switch {
	case x6 < y6: z6=switchFn(x6)
	case x6 > y6: z6=switchFn(y6)
	case x6 == y6: z6=switchFn(x6+y6)
	}
	assert(100, z6, "z6")

	var z7 int
	switch x7 := switchFn(0); {
	case x7 == 5:
		z7 = 1
	case x7 == 6:
		z7 = 2
	case x7 == 100:
		z7 = 3
	case x7 == 10:
		z7 = 4
	}
	assert(1, z7, "z7")
	z7 = 0
	switch x7 := switchFn(1); {
	case x7 == 5:
		z7 = 1
	case x7 == 6:
		z7 = 2
	case x7 == 100:
		z7 = 3
	case x7 == 10:
		z7 = 4
	}
	assert(2, z7, "z7")
	z7 = 0
	switch x7 := switchFn(5); {
	case x7 == 5:
		z7 = 1
	case x7 == 6:
		z7 = 2
	case x7 == 100:
		z7 = 3
	case x7 == 10:
		z7 = 4
	}
	assert(3, z7, "z7")
	z7 = 0
	switch x7 := switchFn(7); {
	case x7 == 5:
		z7 = 1
	case x7 == 6:
		z7 = 2
	case x7 == 100:
		z7 = 3
	case x7 == 10:
		z7 = 4
	}
	assert(4, z7, "z7")

	var z8 int
	if x8, y8 := switchFn(0), 8; x8 < y8 {
		z8 = 1
	} else if x8 > z8 {
		z8 = 2
	} else {
		z8 = 3
	}
	assert(1, z8, "z8")
	z8 = 0
	if x8, y8 := switchFn(200), 8; x8 < y8 {
		z8 = 1
	} else if x8 > z8 {
		z8 = 2
	} else {
		z8 = 3
	}
	assert(2, z8, "z8")
	z8 = 10
	if x8, y8 := switchFn(1)+3, 8; x8 < y8 {
		z8 = 1
	} else if x8 > z8 {
		z8 = 2
	} else {
		z8 = 3
	}
	assert(3, z8, "z8")

	println("OK")
}
