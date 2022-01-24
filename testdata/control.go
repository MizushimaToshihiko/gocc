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
	i = 0
	switch i {
	case 0:
		i = 5
	case 1:
		i = 6
	case 2:
		i = 7
	}
	assert(5, i, "i=0; switch i { case 0:i=5; case 1:i=6; case 2:i=7; } i")
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

	// assert(3, ({ int i=0; switch(-1) { case 0xffffffff: i=3; break; } i; }));

	println("\nOK")
}
