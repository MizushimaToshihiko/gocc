var g1 int
var g2 [4]int
var g3 byte = 3

// var g4 int16 = 4
// var g5 int = 5
// var g6 int64 = 6
// var g7 *int = &g5
// var g8 string = "abc"
// var g8_1 [3]byte = [3]byte{'a', 'b', 'c'}
// var g9 [3]int = [3]int{0, 1, 2}
// var g10 [2]string = [2]string{"foo", "bar"}

// type T_g11 struct {
// 	a byte
// 	b int
// }

// var g11 [2]T_g11 = [2]T_g11{T_g11{1, 2}, T_g11{3, 4}}

// type T_g12 struct {
// 	a [2]int
// }

// var g12 [2]T_g12 = [2]T_g12{T_g12{[2]int{1,2}}, T_g12{[2]int{3,4}}}

// var g12_1 [2]T_g12 = [2]T_g12{T_g12{[2]int{1,2,}}, T_g12{[2]int{3,4,},}}

// var g12_2 []int

// var g12_3 [2]T_g12 = [2]T_g12{{{1,2,}}, {{3,4,},}}

func assert(want int64, ac int64, code *byte) {
	if want == ac {
		printf("\n%s => %ld\n", code, ac)
	} else {
		printf("\n%s => %ld expeted but got %ld\n", code, want, ac)
		exit(1)
	}
}

// func ret3() int {
// 	return 3
// 	return 5
// }

// func add2(x int, y int) int {
// 	return x + y
// }

// func sub2(x int, y int) int {
// 	return x - y
// }

// func add6(a int, b int, c int, d int, e int, f int) int {
// 	return a + b + c + d + e + f
// }

// func addx(x *int, y int) int {
// 	return *x + y
// }

// func subChar(a byte, b byte, c byte) byte {
// 	return a - b - c
// }

// func f31() int {
// 	var a int = 3
// 	return a
// }

// func f32() int {
// 	var a int = 3
// 	var z int = 5
// 	return a + z
// }

// func f33() int {
// 	var foo int = 3
// 	return foo
// }

// func f34() int {
// 	var foo123 int = 3
// 	var bar int = 5
// 	return foo123 + bar
// }

// func f35() int {
// 	if 0 {
// 		return 2
// 	}
// 	return 3
// }

// func f36() int {
// 	if 1 - 1 {
// 		return 2
// 	}
// 	return 3
// }

// func f37() int {
// 	if 2 - 1 {
// 		return 2
// 	}
// 	return 3
// }

// func f40() int {
// 	var i int = 0
// 	for i < 10 {
// 		i = i + 1
// 	}
// 	return i
// }

// func f41() int {
// 	var i int = 0
// 	for {
// 		i = i + 1
// 		if i > 5 {
// 			return i
// 		}
// 	}
// 	return 0
// }

// func f42() int {
// 	var i int = 0
// 	var j int = 0
// 	for i = 0; i <= 10; i = i + 1 {
// 		j = i + j
// 	}
// 	return j
// }

// func f43() int {
// 	for {
// 		return 3
// 		return 5
// 	}
// }

// func f52() int {
// 	var x int = 3
// 	return *&x
// }

// func f53() int {
// 	var x int = 3
// 	var y *int = &x
// 	var z **int = &y
// 	return **z
// }

// func f54() int {
// 	var x int = 3
// 	var y int = 5
// 	return *(&x + 1)
// }

// func f55() int {
// 	var x int = 3
// 	var y int = 5
// 	return *(&y - 1)
// }

// func f56() int {
// 	var x int = 3
// 	var y *int = &x
// 	*y = 5
// 	return x
// }

// func f57() int {
// 	var x int = 3
// 	var y int = 5
// 	*(&x + 1) = 7
// 	return y
// }

// func f58() int {
// 	var x int = 3
// 	var y int = 5
// 	*(&y - 1) = 7
// 	return x
// }

// func f59() int {
// 	var x [2]int
// 	var y *int = &x
// 	*y = 3
// 	return *x
// }

// func f60() int {
// 	var x [3]int
// 	*x = 3
// 	*(x + 1) = 4
// 	*(x + 2) = 5
// 	return *x
// }

// func f61() int {
// 	var x [3]int
// 	*x = 3
// 	*(x + 1) = 4
// 	*(x + 2) = 5
// 	return *(x + 1)
// }

// func f62() int {
// 	var x [3]int
// 	*x = 3
// 	*(x + 1) = 4
// 	*(x + 2) = 5
// 	return *(x + 1)
// }

// func f63() int {
// 	var x [3]int
// 	*x = 3
// 	*(x + 1) = 4
// 	*(x + 2) = 5
// 	return *(x + 2)
// }

// func f64() int {
// 	var x [2][3]int
// 	var y *int = &x
// 	*y = 0
// 	return **x
// }

// func f65() int {
// 	var x [2][3]int
// 	var y *int = &x
// 	*(y + 1) = 1
// 	return *(*x + 1)
// }

// func f66() int {
// 	var x [2][3]int
// 	var y *int = &x
// 	*(y + 2) = 2
// 	return *(*x + 2)
// }

// func f67() int {
// 	var x [2][3]int
// 	var y *int = &x
// 	*(y + 3) = 3
// 	return **(x + 1)
// }

// func f68() int {
// 	var x [2][3]int
// 	var y *int = &x
// 	*(y + 4) = 4
// 	return *(*(x + 1) + 1)
// }

// func f69() int {
// 	var x [2][3]int
// 	var y *int = &x
// 	*(y + 5) = 5
// 	return *(*(x + 1) + 2)
// }

// func f70() int {
// 	var x [2][3]int
// 	var y *int = &x
// 	*(y + 6) = 6
// 	return **(x + 2)
// }

// func f71() int {
// 	var x [3]int
// 	*x = 3
// 	x[1] = 4
// 	x[2] = 5
// 	return *x
// }

// func f72() int {
// 	var x [3]int
// 	*x = 3
// 	x[1] = 4
// 	x[2] = 5
// 	return *(x + 1)
// }

// func f73() int {
// 	var x [3]int
// 	*x = 3
// 	x[1] = 4
// 	x[2] = 5
// 	return *(x + 2)
// }

// func f74() int {
// 	var x [3]int
// 	*x = 3
// 	x[1] = 4
// 	x[2] = 5
// 	return *(x + 2)
// }

// func f75() int {
// 	var x [3]int
// 	*x = 3
// 	x[1] = 4
// 	2[x] = 5
// 	return *(x + 2)
// }

// func f76() int {
// 	var x [2][3]int
// 	var y *int = &x
// 	y[0] = 0
// 	return x[0][0]
// }

// func f77() int {
// 	var x [2][3]int
// 	var y *int = &x
// 	y[1] = 1
// 	return x[0][1]
// }

// func f78() int {
// 	var x [2][3]int
// 	var y *int = &x
// 	y[2] = 2
// 	return x[0][2]
// }

// func f79() int {
// 	var x [2][3]int
// 	var y *int = &x
// 	y[3] = 3
// 	return x[1][0]
// }

// func f80() int {
// 	var x [2][3]int
// 	var y *int = &x
// 	y[4] = 4
// 	return x[1][1]
// }

// func f81() int {
// 	var x [2][3]int
// 	var y *int = &x
// 	y[5] = 5
// 	return x[1][2]
// }

// func f82() int {
// 	var x [2][3]int
// 	var y *int = &x
// 	y[6] = 6
// 	return x[2][0]
// }

// func f90() int {
// 	var x byte = 1
// 	var y byte = 2
// 	return x
// }

// func f91() int {
// 	var x byte = 1
// 	var y byte = 2
// 	return y
// }

// func f97() int {
// 	/*return 1 */
// 	return 2
// }

// func f98() int {
// 	// return 1
// 	return 2
// }

// func f99() int {
// 	var x int = 2
// 	{
// 		var x int = 3
// 	}
// 	return x
// }

// func f100() int {
// 	var x int = 2
// 	{
// 		var x int = 3
// 	}
// 	{
// 		var y int = 4
// 		return x
// 	}
// }

// func f101() int {
// 	var x int = 2
// 	{
// 		x = 3
// 	}
// 	return x
// }

// func f102() int {
// 	type x struct {
// 		a int
// 		b int
// 	}
// 	var x102 x
// 	x102.a = 1
// 	x102.b = 2
// 	return x102.a
// }

// type x103 struct {
// 	a int
// 	b int
// }
// func f103() int {
// 	var x x103
// 	x.a=1
// 	x.b=2
// 	return x.a
// }
// func f104() int {
// 	var x x103
// 	x.a=1
// 	x.b=2
// 	return x.b
// }
// type x105 struct {
// 	a byte
// 	b int
// 	c byte
// }
// func f105() int {
// 	var x x105
// 	x.a=1
// 	x.b=2
// 	x.c=3
// 	return x.a
// }
// func f106() int {
// 	var x x105
// 	x.a=1
// 	x.b=2
// 	x.c=3
// 	return x.b
// }
// func f107() int {
// 	var x x105
// 	x.a=1
// 	x.b=2
// 	x.c=3
// 	return x.c
// }

// type x108 [3]struct{
// 	a int
// 	b int
// }
// func f108() int {
// 	var x x108
// 	var p *int = &x
// 	p[0] = 0
// 	return x[0].a
// }
// func f109() int {
// 	var x x108
// 	var p *int = &x
// 	p[1] = 1
// 	return x[0].b
// }
// func f110() int {
// 	var x x108
// 	var p *int = &x
// 	p[2] = 2
// 	return x[1].a
// }
// func f111() int {
// 	var x x108
// 	var p *int = &x
// 	p[3] = 3
// 	return x[1].b
// }

// type x112 struct {
// 	a [3]int
// 	b [5]int
// }
// func f112() int {
// 	var x x112
// 	var p *int = &x
// 	x.a[0]=6
// 	return p[0]
// }
// func f113() int {
// 	var x x112
// 	var p *int = &x
// 	x.b[0]=7
// 	return p[3]
// }

// type x114 struct {
// 	a struct {
// 		b int
// 	}
// }
// func f114() int {
// 	var x x114
// 	x.a.b = 6
// 	return x.a.b
// }

// // // pointer arithmetic => not supported
// // func f115() int {
// // 	var x int
// // 	var y byte
// // 	var a *int = &x
// // 	var b *int = &y
// // 	return b - a
// // }

// // func f116() int {
// // 	var x byte
// // 	var y int
// // 	var a *int = &x
// // 	var b *int = &y
// // 	return b - a
// // }

// type t117 int
// func f117() int {
// 	var x t117 = 1
// 	return x
// }
// type t118 struct {
// 	a int
// }
// func f118() int {
// 	var x t118
// 	x.a = 1
// 	return x.a
// }
// type t119 int
// func f119() int {
// 	var t119 t119 = 1
// 	return t119
// }
// type t120 struct {
// 	a int
// }
// func f120() int {
// 	{
// 		type t120 int
// 	}
// 	var x t120
// 	x.a = 2
// 	return x.a
// }
// func f121 () int {
// 	var x int
// 	return Sizeof(x)
// }
// func f122() int {
// 	var x int
// 	return Sizeof x
// }
// func f123() int {
// 	var x [4]int
// 	return Sizeof(x)
// }
// func f124() int {
// 	var x [3][4]int
// 	return Sizeof(x)
// }
// func f125() int {
// 	var x [3][4]int
// 	return Sizeof(*x)
// }
// func f126() int {
// 	var x [3][4]int
// 	return Sizeof(**x)
// }
// func f127() int {
// 	var x [3][4]int
// 	return Sizeof(**x) + 1
// }
// func f128() int {
// 	var x [3][4]int
// 	return Sizeof **x + 1
// }
// func f129() int {
// 	var x [3][4]int
// 	return Sizeof(**x + 1)
// }

// func f130() int {
// 	type X struct {
// 		a int
// 	}
// 	var x X
// 	return Sizeof(x)
// }
// func f131() int {
// 	type X struct {
// 		a int
// 		b int
// 	}
// 	var x X
// 	return Sizeof(x)
// }
// func f132() int {
// 	type X struct {
// 		a [3]int
// 	}
// 	var x X
// 	return Sizeof(x)
// }
// func f133() int {
// 	type X struct {
// 		a int
// 	}
// 	var x [4]X
// 	return Sizeof(x)
// }
// func f134() int {
// 	type X struct {
// 		a [3]int
// 	}
// 	var x [2]X
// 	return Sizeof(x)
// }
// func f135() int {
// 	type X struct {
// 		a byte
// 		b int
// 	}
// 	var x X
// 	return Sizeof(x)
// }
// func f136() int {
// 	type X struct {
// 		a int
// 		b byte
// 	}
// 	var x X
// 	return Sizeof(x)
// }

// // func f137() int {
// // 	var x int
// // 	var y byte
// // 	var a *int = &x
// // 	var b *int = &y
// // 	return b - a
// // }

// func subShort(a int16, b int16, c int16) int {
// 	return a - b - c
// }

// func subLong(a int64, b int64, c int64) int {
// 	return a - b - c
// }

// func f138() int {
// 	var x int16
// 	return Sizeof(x)
// }

// func f139() int {
// 	type X struct {
// 		a byte
// 		b int16
// 	}
// 	var x X
// 	return Sizeof(x)
// }

// func f140() int {
// 	var x int64
// 	return Sizeof(x)
// }

// func f141() int {
// 	type X struct {
// 		a byte
// 		b int64
// 	}
// 	var x X
// 	return Sizeof(x)
// }

// // I'm not sure.
// // func g1Ptr() *int {
// // 	return &g1
// // }

// // func f142() *int {
// // 	return *g1Ptr()
// // }

// func f143() bool {
// 	var x bool = 0
// 	return x
// }

// func f144() bool {
// 	var x bool = 1
// 	return x
// }

// func f145() bool {
// 	var x bool = 2
// 	return x
// }

// // this must cause compile error
// // func f146() int {
// // 	var x int = 5
// // 	var y int64 = int64(&x)
// // 	return *int(*y)
// // }

// func charFn() byte {
// 	return 257
// }

// // I will postpone it
// // func f147() int {
// // 	var x, y, z, a, b, c int = 1, 2, 3, 4, 5, 6
// // 	printf("x: %d\n", x)
// // 	printf("y: %d\n", y)
// // 	printf("z: %d\n", z)
// // 	printf("a: %d\n", a)
// // 	printf("b: %d\n", b)
// // 	printf("c: %d\n", c)
// // 	return x + y + z + a + b + c
// // }
// // func f148() int {
// // 	var x int = 1
// // 	var y int = 2
// // 	var z int = 3
// // 	var a int = 4
// // 	var b int = 5
// // 	var c int = 6
// // 	printf("x: %d\n", x)
// // 	printf("y: %d\n", y)
// // 	printf("z: %d\n", z)
// // 	printf("a: %d\n", a)
// // 	printf("b: %d\n", b)
// // 	printf("c: %d\n", c)
// // 	return x + y + z + a + b + c
// // }

// func f149() int {
// 	var i int = 2
// 	i++
// 	return i
// }

// func f150() int {
// 	var i int = 2
// 	i--
// 	return i
// }

// func f151() int {
// 	var a [3]int
// 	a[0]=0
// 	a[1]=1
// 	a[2]=2
// 	a[0]++
// 	a[1]--
// 	a[2]--
// 	return a[0]+a[1]+a[2]
// }

// func f152() int {
// 	type x struct {
// 		y x114
// 	}
// 	var z [3]x
// 	z[0].y.a.b = 6
// 	z[0].y.a.b++
// 	return z[0].y.a.b
// }

// func f153() int {
// 	var i int = 2
// 	i += 5
// 	return i
// }

// func f154() int {
// 	var i int = 2
// 	return i+=5
// }

// func f155() int {
// 	var i int=5
// 	i-=2
// 	return i
// }

// func f156() int {
// 	var i int=5
// 	return i-=2
// }

// func f157() int {
// 	var i int=3
// 	i*=2
// 	return i
// }

// func f158() int {
// 	var i int=3
// 	return i*=2
// }

// func f159() int {
// 	var i int=6
// 	i/=2
// 	return i
// }

// func f160() int {
// 	var i int=6
// 	return i/=2
// }

// func f161() int {
// 	var i int=0
// 	for ;i<10;i++{
// 		if i == 3 {
// 			break
// 		}
// 	}
// 	return i
// }

// func f162() int {
// 	var i int=0
// 	for {
// 		if i++ == 3 {
// 			break
// 		}
// 	}
// 	return i
// }

// func f163() int {
// 	var i int=0
// 	for ;i<10;i++ {
// 		for {
// 			break
// 		}
// 		if i == 3 {
// 			break
// 		}
// 	}
// 	return i
// }

// func f164() int {
// 	var i int=0
// 	for {
// 		for {
// 			break
// 		}
// 		if i++ == 3 {
// 			break
// 		}
// 	}
// 	return i
// }

// func f165() int {
// 	var i int=0
// 	var j int=0
// 	for ;i<10;i++{
// 		if i>5 {
// 			continue
// 		}
// 		j++
// 	}
// 	return i
// }

// func f166() int {
// 	var i int=0
// 	var j int=0
// 	for ;i<10;i++{
// 		if i>5 {
// 			continue
// 		}
// 		j++
// 	}
// 	return j
// }

// func f167() int {
// 	var i int=0
// 	var j int=0
// 	for ;!i;{
// 		for ;j!=10;j++ {
// 			continue
// 		}
// 		break
// 	}
// 	return j
// }

// func f168() int {
// 	var i int=0
// 	var j int=0
// 	for i++<10 {
// 		if i>5 {
// 			continue
// 		}
// 		j++
// 	}
// 	return i
// }

// func f169() int {
// 	var i int=0
// 	var j int=0
// 	for i++<10 {
// 		if i>5 {
// 			continue
// 		}
// 		j++
// 	}
// 	return j
// }

// func f170() int {
// 	var i int=0
// 	var j int=0
// 	for !i {
// 		for j++!=10 {
// 			continue
// 		}
// 		break
// 	}
// 	return j
// }

// func f171() int {
// 	var i int=0
// 	goto a
// 	a:
// 		i++
// 	b:
// 		i++
// 	c:
// 		i++
// 	return i
// }

// func f172() int {
// 	var i int=0
// 	goto e
// 	d:
// 		i++
// 	e:
// 		i++
// 	f:
// 		i++
// 	return i
// }

// func f173() int {
// 	var i int=0
// 	goto i
// 	g:
// 		i++
// 	h:
// 		i++
// 	i:
// 		i++
// 	return i
// }

// func f174() int {
// 	var i int=0
// 	switch 0 {
// 	case 0:
// 		i=5
// 	case 1:
// 		i=6
// 	case 2:
// 		i=7
// 	}
// 	return i
// }
// func f175() int {
// 	var i int=0
// 	switch 1 {
// 	case 0:
// 		i=5
// 	case 1:
// 		i=6
// 	case 2:
// 		i=7
// 	}
// 	return i
// }

// func f176() int {
// 	var i int=0
// 	switch 2 {
// 	case 0:
// 		i=5
// 	case 1:
// 		i=6
// 	case 2:
// 		i=7
// 	}
// 	return i
// }

// func f177() int {
// 	var i int=0
// 	switch 3 {
// 	case 0:
// 		i=5
// 	case 1:
// 		i=6
// 	case 2:
// 		i=7
// 	}
// 	return i
// }

// func f178() int {
// 	var i int=0
// 	switch 0 {
// 	case 0:
// 		i=5
// 	default:
// 		i=7
// 	}
// 	return i
// }

// func f179() int {
// 	var i int=0
// 	switch 1 {
// 	case 0:
// 		i=5
// 	default:
// 		i=7
// 	}
// 	return i
// }

// func voidFn() {}

// func f180() int {
// 	var i int=1
// 	i<<=0
// 	return i
// }

// func f181() int {
// 	var i int=1
// 	i<<=3
// 	return i
// }

// func f182() int {
// 	var i int=5
// 	i<<=1
// 	return i
// }

// func f183() int {
// 	var i int=5
// 	i>>=1
// 	return i
// }

// func f184() int {
// 	var i int=-1
// 	i>>=1
// 	return i
// }

// func f185() int {
// 	var i int=0
// 	switch 3 {
// 	case 5-2+0*3:
// 		i++
// 	}
// 	return i
// }

// func f186() int {
// 	var x [1+1]int
// 	return Sizeof(x)
// }

// // func f187() int {
// // 	var y [2]int = [2]int{1, 2}
// // 	var x [2]int = y
// // 	return x[0]
// // }

// // lvar-initializers
// func f188() int {
// 	var x [2][2]int = [2][2]int{[2]int{1, 2}, [2]int{4, 5}}
// 	return x[1][1]
// }

// // // error occurs
// // func f189() int {
// // 	var x [2][3]int = [2][3]int{{1, 2, 3}, {4, 5, 6}}
// // 	return x[1][1]
// // }

// func f190() int {
// 	var x [3]int = [3]int{1,2,3}
// 	return x[0]
// }

// func f191() int {
// 	var x [3]int = [3]int{1,2,3}
// 	return x[1]
// }

// func f192() int {
// 	var x [3]int = [3]int{1,2,3}
// 	return x[2]
// }

// func f193() int {
// 	var x [2][3]int = [2][3]int{[3]int{1,2,3},[3]int{4,5,6}}
// 	return x[0][1]
// }

// func f194() int {
// 	var x [2][3]int = [2][3]int{[3]int{1,2,3},[3]int{4,5,6}}
// 	return x[1][0]
// }

// func f195() int {
// 	var x [2][3]int = [2][3]int{[3]int{1,2,3},[3]int{4,5,6}}
// 	return x[1][2]
// }

// func f196() int {
// 	var x [2][3]int = [2][3]int{[3]int{1, 2}}
// 	return x[0][1]
// }

// func f197() int {
// 	var x [2][3]int=[2][3]int{[3]int{1, 2}}
// 	return x[1][0]
// }

// func f198() int {
// 	var x [2][3]int=[2][3]int{[3]int{1, 2}}
// 	return x[1][2]
// }

// func f199() int {
// 	var x string="abc"
// 	return x[0]
// }

// func f200() int {
// 	var x string = "abc"
// 	return x[2]
// }

// func f201() int {
// 	var x string
// 	x = "abc"
// 	return x[2]
// }

// func f202() int {
// 	type T struct {
// 		a int
// 		b int
// 	}
// 	var x T = T{1, 2}
// 	return x.a
// }

// func f203() int {
// 	type T struct {
// 		a int
// 		b int
// 		c int
// 	}
// 	var x T = T{1, 2, 3}
// 	return x.a
// }

// func f204() int {
// 	type T struct {
// 		a int
// 		b int
// 		c int
// 	}
// 	var x T = T{1, 2, 3}
// 	return x.b
// }

// func f205() int {
// 	type T struct {
// 		a int
// 		b int
// 		c int
// 	}
// 	var x T = T{1, 2, 3}
// 	return x.c
// }

// func f206() int {
// 	type T struct {
// 		a int
// 		b int
// 		c int
// 	}
// 	var x T = T{1}
// 	return x.a
// }

// func f207() int {
// 	type T struct {
// 		a int
// 		b int
// 		c int
// 	}
// 	var x T = T{1}
// 	return x.b
// }

// func f208() int {
// 	type T struct {
// 		a int
// 		b int
// 		c int
// 	}
// 	var x T = T{1}
// 	return x.c
// }

// func f209() int {
// 	type T struct {
// 		a int
// 		b int
// 	}
// 	var x [2]T = [2]T{T{1,2},T{3,4}}
// 	return x[0].a
// }

// func f210() int {
// 	type T struct {
// 		a int
// 		b int
// 	}
// 	var x [2]T = [2]T{T{1,2},T{3,4}}
// 	return x[0].b
// }

// func f211() int {
// 	type T struct {
// 		a int
// 		b int
// 	}
// 	var x [2]T = [2]T{T{1,2},T{3,4}}
// 	return x[1].a
// }

// func f212() int {
// 	type T struct {
// 		a int
// 		b int
// 	}
// 	var x [2]T = [2]T{T{1,2},T{3,4}}
// 	return x[1].b
// }

// func f213() string {
// 	var x string
// 	x = "abc"
// 	return x
// }

// func f214() string {
// 	return g8
// }

// func f215() []byte {
// 	var x []byte
// 	return x
// }

// func f216() int {
// 	var x int
// 	var y int
// 	x=1,y=2
// 	return x+y
// }

// func f217() int {
// 	var x [2][3]int = [2][3]int{{1,2,3},{4,5,6}}
// 	return x[1][2]
// }

// // type cast
// func f218() int {
// 	var x byte = 1
// 	return int(x)
// }

// func f219() int {
// 	var x int64 = 8590066177
// 	return int(x)
// }

// func f220() int16 {
// 	var x int64 = 8590066177
// 	return int16(x)
// }

// func f221() byte {
// 	var x int64 = 8590066177
// 	return byte(x)
// }

// func main() {
// 	assert(0, 0, "0")
// 	assert(42, 42, "42")
// 	assert(5, 5, "5")
// 	assert(41, 12+34-5, "12 + 34 - 5")
// 	assert(15, 5*(9-6), "5*(9-6)")
// 	assert(4, (3+5)/2, "(3+5)/2")
// 	assert(-10, -10, "-10")
// 	assert(10, - -10, "- -10")
// 	assert(10, - -+10, "- - +10")

// 	assert(0, 0==1, "0==1")
// 	assert(1, 42==42, "42==42")
// 	assert(1, 0!=1, "0!=1")
// 	assert(0, 42!=42, "42!=42")

// 	assert(1, 0<1, "0<1")
// 	assert(0, 1<1, "1<1")
// 	assert(0, 2<1, "2<1")
// 	assert(1, 0<=1, "0<=1")
// 	assert(1, 1<=1, "1<=1")
// 	assert(0, 2<=1, "2<=1")

// 	assert(1, 1>0, "1>0")
// 	assert(0, 1>1, "1>1")
// 	assert(0, 1>2, "1>2")
// 	assert(1, 1>=0, "1>=0")
// 	assert(1, 1>=1, "1>=1")
// 	assert(0, 1>=2, "1>=2")

// 	assert(3, ret3(), "ret3()")

// 	assert(8, add2(3, 5), "add(3, 5)")
// 	assert(2, sub2(5, 3), "sub(5, 3)")
// 	assert(21, add6(1, 2, 3, 4, 5, 6), "add6(1, 2, 3, 4, 5, 6)")
// 	assert(55, fib(9), "fib(9)")

// 	assert(0, g1, "g1")
// 	g1 = 3
// 	assert(3, g1, "g1")

// 	g2[0] = 0
// 	g2[1] = 1
// 	g2[2] = 2
// 	g2[3] = 3
// 	assert(0, g2[0], "g2[0]")
// 	assert(1, g2[1], "g2[1]")
// 	assert(2, g2[2], "g2[2]")
// 	assert(3, g2[3], "g2[3]")

// 	assert(1, subChar(7, 3, 3), "subChar(7, 3, 3)")

// 	assert(97, "abc"[0], "\"abc\"[0]")
// 	assert(98, "abc"[1], "\"abc\"[1]")
// 	assert(99, "abc"[2], "\"abc\"[2]")
// 	assert(0, "abc"[3], "\"abc\"[3]")

// 	assert(7, "\a"[0], "\"\\a\"[0]")
// 	assert(8, "\b"[0], "\"\\b\"[0]")
// 	assert(9, "\t"[0], "\"\\t\"[0]")
// 	assert(10, "\n"[0], "\"\\n\"[0]")
// 	assert(11, "\v"[0], "\"\\v\"[0]")
// 	assert(12, "\f"[0], "\"\\f\"[0]")
// 	assert(13, "\r"[0], "\"\\r\"[0]")
// 	assert(27, "\e"[0], "\"\\e\"[0]")
// 	assert(0, "\0"[0], "\"\\0\"[0]")

// 	assert(106, "\j"[0], "\"\\j\"[0]")
// 	assert(107, "\k"[0], "\"\\k\"[0]")
// 	assert(108, "\l"[0], "\"\\l\"[0]")

// 	assert(3, f31(), "func f31() int {\n\tvar a int=3\n\treturn a\n}")
// 	assert(8, f32(), "func f32() int {\n\tvar a int=3\n\tvar z int=5\n\treturn a+z\n}")
// 	assert(3, f33(), "func f33() int {\n\tvar foo int=3\n\treturn foo\n}")
// 	assert(8, f34(), "func f34() int {\n\tvar foo123 int=3\n\tvar bar int=5\n\treturn foo123+bar\n}")

// 	assert(3, f35(), "func f35() int {\n\tif 0 {\n\t\treturn 2\n\t}\n\treturn 3\n}")
// 	assert(3, f36(), "func f36() int {\n\tif 1-1{\n\t\treturn 2\n\t}\n\treturn 3\n}")
// 	assert(2, f37(), "func f37() int {\n\tif 2-1{\n\t\treturn 2\n\t}\n\treturn 3\n}")

// 	assert(10, f40(), "func f40() int {\n\tvar i int=0\n\tfor i<10 {\n\t\ti=i+1\n\t}\n\treturn i\n}")
// 	assert(6, f41(), "func f41() int {\n\tvar i int=0\n\tfor {\n\t\ti=i+1\n\t\tif i>5 {\n\t\t\treturn i\n\t\t}\n\t}\n\treturn 0\n}")
// 	assert(55, f42(), "func f42() int {\n\tvar i int=0\n\tvar j int=0\n\tfor i=0; i<=10; i=i+1 {\n\t\tj=i+j\n\t}\n\treturn j\n}")
// 	assert(3, f43(), "func f43() int {\nfor ;; {\n\treturn 3\n\treturn 5\n}\n}")

// 	assert(3, f52(), "func f52() int {\n\tvar x int=3\n\treturn *&x\n}")
// 	assert(3, f53(), "func f53() int {\n\tvar x int=3\n\tvar y *int=&x\n\tvar z **int=&y\n\treturn **z\n}")
// 	assert(5, f54(), "func f54() int {\n\tvar x int=3\n\tvar y int=5\n\treturn *(&x+1)\n}")
// 	assert(3, f55(), "func f55() int {\n\tvar x int=3\n\tvar y int=5\n\treturn *(&y-1)\n}")
// 	assert(5, f56(), "func f56() int {\n\tvar x int=3\n\tvar y *int=&x\n\t*y=5\n\treturn x\n}")
// 	assert(7, f57(), "func f57() int {\n\tvar x int=3\n\tvar y int=5\n\t*(&x+1)=7\n\treturn y\n}")
// 	assert(7, f58(), "func f58() int {\n\tvar x int=3\n\tvar y int=5\n\t*(&y-1)=7\n\treturn x\n}")
// 	assert(3, f59(), "func f59() int {\n\tvar x [2]int\n\tvar y *int=&x\n\t*y=3\n\treturn *x\n}")

// 	assert(3, f60(), "func f60() int {\n\tvar x [3]int\n\t*x=3\n\t*(x+1)=4\n\t*(x+2)=5\n\treturn *x\n}")
// 	assert(4, f61(), "func f61() int {\n\tvar x [3]int\n\t*x=3\n\t*(x+1)=4\n\t*(x+2)=5\n\treturn *(x+1)\n}")
// 	assert(5, f63(), "func f63() int {\n\tvar x [3]int\n\t*x=3\n\t*(x+1)=4\n\t*(x+2)=5\n\treturn *(x+2)\n}")

// 	assert(0, f64(), "func f64() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\t*y=0\n\treturn **x\n}")
// 	assert(1, f65(), "func f65() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\t*(y+1)=1\n\treturn *(*x+1)\n}")
// 	assert(2, f66(), "func f66() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\t*(y+2)=2\n\treturn *(*x+2)\n}")
// 	assert(3, f67(), "func f67() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\t*(y+3)=3\n\treturn **(x+1)\n}")
// 	assert(4, f68(), "func f68() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\t*(y+4)=4\n\treturn *(*(x+1)+1)\n}")
// 	assert(5, f69(), "func f69() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\t*(y+5)=5\n\treturn *(*(x+1)+2)\n}")
// 	assert(6, f70(), "func f70() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\t*(y+6)=6\n\treturn **(x+2)\n}")

// 	assert(3, f71(), "func f71() int {\n\tvar x [3]int\n\t*x=3\n\tx[1]=4\n\tx[2]=5\n\treturn *x\n}")
// 	assert(4, f72(), "func f72() int {\n\tvar x [3]int\n\t*x=3\n\tx[1]=4\n\tx[2]=5\n\treturn *(x+1)\n}")
// 	assert(5, f73(), "func f73() int {\n\tvar x [3]int\n\t*x=3\n\tx[1]=4\n\tx[2]=5\n\treturn *(x+2)\n}")
// 	assert(5, f74(), "func f74() int {\n\tvar x [3]int\n\t*x=3\n\tx[1]=4\n\tx[2]=5\n\treturn *(x+2)\n}")
// 	assert(5, f75(), "func f75() int {\n\tvar x [3]int\n\t*x=3\n\tx[1]=4\n\t2[x]=5\n\treturn *(x+2)\n}")

// 	assert(0, f76(), "func f76() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\ty[0]=0\n\treturn x[0][0]\n}")
// 	assert(1, f77(), "func f77() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\ty[1]=1\n\treturn x[0][1]\n}")
// 	assert(2, f78(), "func f78() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\ty[2]=2\n\treturn x[0][2]\n}")
// 	assert(3, f79(), "func f79() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\ty[3]=3\n\treturn x[1][0]\n}")
// 	assert(4, f80(), "func f80() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\ty[4]=4\n\treturn x[1][1]\n}")
// 	assert(5, f81(), "func f81() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\ty[5]=5\n\treturn x[1][2]\n}")
// 	assert(6, f82(), "func f82() int {\n\tvar x [2][3]int\n\tvar y *int=x\n\ty[6]=6\n\treturn x[2][0]\n}")

// 	assert(1, f90(), "func f90() int {\n\tvar x byte=1\n\tvar y byte=2\n\treturn x\n}")
// 	assert(2, f91(), "func f91() int {\n\tvar x byte=1\n\tvar y byte=2\n\treturn y\n}")

// 	assert(2, f97(), "func f97() int {\n\t/*return 1 */\n\treturn 2\n}")
// 	assert(2, f98(), "func f98() int {\n\t// return 1\nreturn 2\n}")

// 	assert(2, f99(), "func f99() int {\n\tvar x int=2\n\t{\n\t\tvar x int=3\n\t}\n\treturn x\n}")
// 	assert(2, f100(), "func f100() int {\n\tvar x int=2\n\t{\n\t\tvar x int=3\n\t}\n\t{\n\t\tvar y int=4\n\t\treturn x\n\t}\n}")
// 	assert(3, f101(), "func f101() int {\n\tvar x int=2\n\t{\n\t\tx=3\n\t}\n\treturn x\n}")

// 	assert(1, f102(), "func f102() int {\n\ttype x102 struct {\n\t\tint a\n\t\tint b\n\t}\n\tvar x102 x\n\tx102.a = 1\n\tx102.b = 2\n\treturn x102.a\n}")

// 	assert(1, f103(), "type x103 struct {\n\ta int\n\tb int\n}\nfunc f103() int {\n\tvar x x103\n\tx.a=1\n\tx.b=2\n\treturn x.a\n}")
// 	assert(2, f104(), "type x103 struct {\n\ta int\n\tb int\n}\nfunc f103() int {\n\tvar x x103\n\tx.a=1\n\tx.b=2\n\treturn x.b\n}")

// 	assert(1, f105(), "type x105 struct {\n\ta type\n\tb int\n\tc byte\n\t}\nfunc f105() int {\n\tvar x x105\n\tx.a=1\n\tx.b=2\n\tx.c=3\n\treturn x.a\n}")
// 	assert(2, f106(), "type x106 struct {\n\ta type\n\tb int\n\tc byte\n\t}\nfunc f106() int {\n\tvar x x105\n\tx.a=1\n\tx.b=2\n\tx.c=3\n\treturn x.b\n}")
// 	assert(3, f107(), "type x107 struct {\n\ta type\n\tb int\n\tc byte\n\t}\nfunc f107() int {\n\tvar x x105\n\tx.a=1\n\tx.b=2\n\tx.c=3\n\treturn x.c\n}")

// 	assert(6, f112(), "type x112 struct {\n\ta [3]int\n\tb [5]int\n}\nfunc f112 int {\n\tvar x x112\n\tp *int=&x\n\tx.a[0]=6\n\treturn p[0]\n}")
// 	assert(7, f113(), "type x112 struct {\n\ta [3]int\n\tb [5]int\n}\nfunc f113 int {\n\tvar x x112\n\tp *int=&x\n\tx.b[7]=7\n\treturn p[3]\n}")

// 	assert(6, f114(), "type x114 struct {\n\t type a struct{\n\t\tb int\n\t}\n}\nfunc f114 int {\n\tvar x x114\n\tx.a.b=6\n\treturn x.a.b\n}")

// 	assert(0, f108(), "type x108 [3]struct {\n\ta int\n\tb int\n}\nfunc f108() int {\n\tvar x x108\n\tp *int=x\n\tp[0]=0\n\treturn x[0].a\n}")
// 	assert(1, f109(), "type x108 [3]struct {\n\ta int\n\tb int\n}\nfunc f109() int {\n\tvar x x108\n\tp *int=x\n\tp[1]=1\n\treturn x[0].b\n}")
// 	assert(2, f110(), "type x108 [3]struct {\n\ta int\n\tb int\n}\nfunc f110() int {\n\tvar x x108\n\tp *int=x\n\tp[2]=2\n\treturn x[1].a\n}")
// 	assert(3, f111(), "type x108 [3]struct {\n\ta int\n\tb int\n}\nfunc f111() int {\n\tvar x x108\n\tvar p*int=x\n\tp[3]=3\n\t return x[1].b\n}")

// 	// assert(7, f115(), "func f115() int {\n\tvar x int\n\tvar y byte\n\tvar a int = &x\n\tvar b int = &y\n\treturn b - a\n}")

// 	// assert(1, f116(), "func f116() int {\n\tvar x byte\n\tvar y int\n\tvar a int = &x\n\tvar b int = &y\n\treturn b - a\n}")

// 	assert(1, f117(), "type t117 int\nfunc f117() int{\n\tvar x t = 1\n\treturn x\n}")
// 	assert(1, f118(), "type t118 struct {\n\ta int\n}\nfunc f118() int{\n\tvar x t118\n\tx.a = 1\n\treturn x.a\n}")
// 	assert(1, f119(), "type t119 int\nfunc f119() int{\n\tvar t110 t119 = 1\n\treturn t\n}")
// 	assert(2, f120(), "type t120 struct {\n\ta int\n}\nfunc f120(){\n\t{\n\t\ttype t120 int\n\t}\n\tvar x t120\n\tx.a = 2\n\treturn x.a\n}")

// 	assert(4, f121(), "var x int\nreturn Sizeof(x)")
// 	assert(4, f122(), "var x int\nreturn Sizeof x")

// 	assert(16, f123(), "var x [4]int\nreturn Sizeof(x)")

// 	assert(48, f124(), "var x [3][4]int\nreturn Sizeof(x)")
// 	assert(16, f125(), "var x [3][4]int\nreturn Sizeof(*x)")
// 	assert(4, f126(), "var x [3][4]int\nreturn Sizeof(**x)")
// 	assert(5, f127(), "var x [3][4]int\nreturn Sizeof(**x) + 1")
// 	assert(5, f128(), "var x [3][4]int\nreturn Sizeof **x + 1")
// 	assert(4, f129(), "var x [3][4]int\nreturn Sizeof(**x + 1)")

// 	assert(4, Sizeof(g1), "Sizeof(g1)")
// 	assert(16, Sizeof(g2), "Sizeof(g2)")

// 	assert(4, f130(), "type X struct {\n\ta int\n}\nvar x X\nreturn Sizeof(x)")
// 	assert(8, f131(), "type X struct {\n\ta int\n\tb int\n}\nvar x X\nreturn Sizeof(x)")
// 	assert(12, f132(), "type X struct {\n\ta [3]int\n}\nvar x X\nreturn Sizeof(x)")
// 	assert(16, f133(), "type X struct {\n\ta int\n}\nvar x [4]X\nreturn Sizeof(x)")
// 	assert(24, f134(), "type X struct {\n\ta [3]int\n}\n var x [2]X\nreturn Sizeof(x)")
// 	assert(8, f135(), "type X struct {\n\ta byte\n\tb int\n}\nvar x X\nreturn Sizeof(x)")
// 	assert(8, f136(), "type X struct {\n\ta int\n\tb byte\n}\n var x X\nreturn Sizeof(x)")
// 	// assert(7, f137(), "var x int\nvar y byte\nvar a int=&x\nvar b int=&y\nreturn b - a")

// 	assert(2, f138(), "var x int64\nreturn Sizeof(x)")
// 	assert(4, f139(), "type X struct {\n\ta byte\n\tb int16\n}\nvar x X\nreturn Sizeof(x)")

// 	assert(8, f140(), "var x int64\nreturn Sizeof(x)")
// 	assert(16, f141(), "type X struct {\n\ta byte\n\tb int64\n}\nvar x X\nreturn Sizeof(x)")

// 	assert(1, subShort(7, 3, 3), "subShort(7, 3, 3)")
// 	assert(1, subLong(7, 3, 3), "subLong(7, 3, 3)")

// 			// I'm not sure.
// 			// assert(3, f142(), "return *g1Ptr()")

// 	assert(0, f143(), "var x bool = 0\nreturn x")
// 	assert(1, f144(), "var x bool = 1\nreturn x")
// 	assert(1, f145(), "var x bool = 2\nreturn x")

// 	assert(4, Sizeof(0), "Sizeof(0)")
// 	assert(4294967297, 4294967297, "4294967297")
// 	assert(8, Sizeof(4294967297), "Sizeof(4294967297)")

// 	assert(131585, int(8590066177), "int(8590066177)")
// 	assert(513, int16(8590066177), "int16(8590066177)")
// 	assert(1, byte(8590066177), "byte(8590066177)")
// 	assert(1, bool(1), "bool(1)")
// 	assert(1, bool(2), "bool(2)")
// 	assert(0, bool(byte(256)), "bool(byte(256))")
// 	assert(1, int64(1), "int64(1)")

// 	assert(97, 'a', "'a'")
// 	assert(10, '\n', "\'\\n\'")

// 	assert(1, charFn(), "charFn()")

// 			// I will postpone it.
// 			// assert(21, f148(), "var x int\nvar y int\nvar z int\nvar a int\nvar b int\nvar c int\nreturn  x+y+z+a+b+c")
// 			// assert(21, f147(), "var x,y,z,a,b,c int = 1,2,3,4,5,6\nreturn x+y+z+a+b+c")

// 	assert(3, f149(), "f149:\nvar i int = 2\ni++\nreturn i")
// 	assert(1, f150(), "f150:\nvar i int = 2\ni--\nreturn i")
// 	assert(2, f151(), "f151:\nvar a [3]int\na[0]=0\na[1]=1\na[2]=2\na[0]++\na[1]--\na[2]--\nreturn a[0]+a[1]+a[2]")
// 	assert(7, f152(), "f152:\ntype x struct {\n\ty x114\n}\nvar z [3]x\nz[0].y.a.b = 6\nz[0].y.a.b++\nreturn z[0].y.a.b")

// 	assert(7, f153(), "f153:\nvar i int=2\ni+=5\nreturn i")
// 	assert(7, f154(), "f154:\nvar i int=2\nreturn i+=5")
// 	assert(3, f155(), "f155:\nvar i int=5\ni-=2\nreturn i")
// 	assert(3, f156(), "f156:\nvar i int=5\nreturn i-=2")
// 	assert(6, f157(), "f157:\nvar i int=3\ni*=2\nreturn i")
// 	assert(6, f158(), "f158:\nvar i int=3\nreturn i*=2")
// 	assert(3, f159(), "f159:\nvar i int=6\ni/=2\nreturn i")
// 	assert(3, f160(), "f160:\nvar i int=6\nreturn i/=2")

// 	assert(0, !1, "!1")
// 	assert(0, !2, "!2")
// 	assert(1, !0, "!0")

// 	assert(-1, ^0, "^0")
// 	assert(0, ^-1, "^-1")

// 	assert(0, 0&1, "0&1")
// 	assert(1, 3&1, "3&1")
// 	assert(3, 7&3, "7&3")
// 	assert(10, -1&10, " -1&10")

// 	assert(1, 0|1, "0|1")
// 	assert(3, 2|1, "2|1")
// 	assert(3, 1|3, "1|3")

// 	assert(0, 0^0, "0^0")
// 	assert(0, 8^8, "8^8")
// 	assert(4, 7^3, "7^3")
// 	assert(2, 7^5, "7^5")
// 	assert(4, (3+4)&^(2+1), "(3+4)&^(2+1)")

// 	assert(1, 0||1, "0||1")
// 	assert(1, 0||(2-2)||5, "0||(2-2)||5")
// 	assert(0, 0||0, "0||0")
// 	assert(0, 0||(2-2), "0||(2-2)")

// 	assert(0, 0&&1, "0&&1")
// 	assert(0, (2-2)&&5, "(2-2)&&5")
// 	assert(1, 1&&5, "1&&5")

// 	assert(3, f161(), "f161:\nvar i int=0\nfor ;i<10;i++{\n\tif i == 3 {\n\t\tbreak\n\t}\n}\nreturn i")
// 	assert(4, f162(), "f162:\nvar i int=0\nfor {\n\tif i++ == 3 {\n\t\tbreak\n\t}\n}\nreturn i")
// 	assert(3, f163(), "f163:\nvar i int=0\nfor ;i<10;i++ {\n\tfor {\n\t\tbrea\n\t}\n\tif i == 3 {\n\t\tbreakt}\n}\nreturn i")
// 	assert(4, f164(), "f164:\nvar i int=0\nfor {\n\tfor {\n\t\tbreak\n\t}\n\tif i++ == 3 {\n\t\tbreak\n\t}\n}\nreturn i")

// 	assert(10, f165(), "f165: var i int=0\nvar j int=0\nfor ;i<10;i++{\n\tif i>5 {\n\t\tcontinue\n\t}\n\tj++\n}\nreturn i")
// 	assert(6, f166(), "f166: var i int=0\nvar j int=0\nfor ;i<10;i++{\n\tif i>5 {\n\t\tcontinue\n\t}\n\tj++\n}\nreturn j")
// 	assert(10, f167(), "f167: var i int=0\nvar j int=0\nfor ;!i{\n\tfor ;j!=10;j++ {\n\t\tcontinue\n\t}\n\tbreak\n}\nreturn j")
// 	assert(11, f168(), "f168: var i int=0\nvar j int=0\nfor i++<10 {\n\tif i>5 {\n\t\tcontinue\n\t}\n\tj++\n}\nreturn i")
// 	assert(5, f169(), "f169: var i int=0\nvar j int=0\nfor i++<10 {\n\tif i>5 {\n\t\tcontinue\n\t}\n\tj++\n}\nreturn j")
// 	assert(11, f170(), "f170: var i int=0\nvar j int=0\nfor !i {\n\tfor j++!=10 {\n\t\tcontinue\n}\n\tbreak\n}\nreturn j")

// 	assert(3, f171(), "f171:\nvar i int=0\ngoto a\na:\n\ti++\nb:\n\ti++\nc:\n\ti++\nreturn i")
// 	assert(2, f172(), "f172:\nvar i int=0\ngoto e\nd:\n\ti++\ne:\n\ti++\nf:\n\ti++\nreturn i")
// 	assert(1, f173(), "f173:\nvar i int=0\ngoto i\ng:\n\ti++\nh:\n\ti++\ni:\n\ti++\nreturn i")

// 	assert(5, f174(), "f174:\nvar i int=0\nswitch 0 {\ncase 0:\n\ti=5\ncase 1:\n\n\ti=6\ncase 2:\n\ti=7\n}\nreturn i")
// 	assert(6, f175(), "f175:\nvar i int=0\nswitch 1 {\ncase 0:\n\ti=5\ncase 1:\n\n\ti=6\ncase 2:\n\ti=7\n}\nreturn i")
// 	assert(7, f176(), "f176:\nvar i int=0\nswitch 2 {\ncase 0:\n\ti=5\ncase 1:\n\n\ti=6\ncase 2:\n\ti=7\n}\nreturn i")
// 	assert(0, f177(), "f177:\nvar i int=0\nswitch 3 {\ncase 0:\n\ti=5\ncase 1:\n\n\ti=6\ncase 2:\n\ti=7\n}\nreturn i")
// 	assert(5, f178(), "f178:\nvar i int=0\nswitch 0 {\ncase 0:\n\ti=5\ndefault:\n\ti=7\n}\nreturn i")
// 	assert(7, f179(), "f179:\nvar i int=0\nswitch 1 {\ncase 0:\n\ti=5\ndefault:\n\ti=7\n}\nreturn i")

// 	voidFn()

// 	assert(1, 1<<0, "1<<0")
// 	assert(8, 1<<3, "1<<3")
// 	assert(10, 5<<1, "5<<1")
// 	assert(2, 5>>1, "5>>1")
// 	assert(-1, -1>>1, "-1>>1")
// 	assert(1, f180(), "f180:\nvar i int=1\ni<<=0\nreturn i")
// 	assert(8, f181(), "f181:\nvar i int=1\ni<<=3\nreturn i")
// 	assert(10, f182(), "f182:\nvar i int=5\ni<<=1\nreturn i")
// 	assert(2, f183(), "f183:\nvar i int=5\ni>>=1\nreturn i")
// 	assert(-1, -1, "-1")
// 	assert(-1, f184(), "f184:\nvar i int=-1\ni>>=1\nreturn i")

// 	assert(1, f185(), "f185:\nvar i int=0\nswitch 3 {\ncase 5-2+0*3:\n\ti++\n}\nreturn i")
// 	assert(8, f186(), "f186:\nvar x [1+1]int\nreturn Sizeof(x)")

// 	assert(5, f188(), "f188:\nvar x [2][2]int = [2][2]int{{1, 2}, {4, 5}}\n	return x[1][1]")

// 	assert(1, f190(), "f190:\nvar x [3]int = [3]int{1,2,3}\nreturn x[0]")
// 	assert(2, f191(), "f191:\nvar x [3]int = [3]int{1,2,3}\nreturn x[1]")
// 	assert(3, f192(), "f192:\nvar x [3]int = [3]int{1,2,3}\nreturn x[2]")

// 	assert(2, f193(), "f193:\nvar x [2][3]int = [2][3]int{{1,2,3},{4,5,6}}\nreturn x[0][1]")
// 	assert(4, f194(), "f194:\nvar x [2][3]int = [2][3]int{{1,2,3},{4,5,6}}\nreturn x[1][0]")
// 	assert(6, f195(), "f195:\nvar x [2][3]int = [2][3]int{{1,2,3},{4,5,6}}\nreturn x[1][2]")

// 	assert(2, f196(), "f196:\nvar x [2][3]int={{1,2}}\nreturn x[0][1]")
// 	assert(0, f197(), "f197:\nvar x [2][3]int={{1,2}}\nreturn x[1][0]")
// 	assert(0, f198(), "f198:\nvar x [2][3]int={{1,2}}\nreturn x[1][2]")

// 	assert('a', f199(), "var x string=\"abc\"\nreturn x[0]")
// 	assert('c', f200(), "var x string=\"abc\"\nreturn x[2]")
// 	assert('c', f201(), "var x string\nx=\"abc\"\nreturn x")

// 	assert(1, f202(), "f202:\ntype T struct {\n\ta int\n\tb int\n}\nvar x T = T{1, 2}\nreturn x.a")
// 	assert(1, f203(), "f203:\ntype T struct {\n\ta int\n\tb int\n\tc int\n}\nvar x T = T{1, 2, 3}\nreturn x.a")
// 	assert(2, f204(), "f204:\ntype T struct {\n\ta int\n\tb int\n\tc int\n}\nvar x T = T{1, 2, 3}\nreturn x.b")
// 	assert(3, f205(), "f205:\ntype T struct {\n\ta int\n\tb int\n\tc int\n}\nvar x T = T{1, 2, 3}\nreturn x.c")
// 	assert(1, f206(), "f206:\ntype T struct {\n\ta int\n\tb int\n\tc int\n}\nvar x T = T{1}\nreturn x.a")
// 	assert(0, f207(), "f207:\ntype T struct {\n\ta int\n\tb int\n\tc int\n}\nvar x T = T{1}\nreturn x.b")
// 	assert(0, f208(), "f208:\ntype T struct {\n\ta int\n\tb int\n\tc int\n}\nvar x T = T{1}\nreturn x.c")

// 	assert(1, f209(), "f209:\ntype T struct {\n\ta int\n\tb int\n}\nvar x [2]T = [2]T{{1,2},{3,4}}\nreturn x[0].a")
// 	assert(2, f210(), "f210:\ntype T struct {\n\ta int\n\tb int\n}\nvar x [2]T = [2]T{{1,2},{3,4}}\nreturn x[1].b")
// 	assert(3, f211(), "f211:\ntype T struct {\n\ta int\n\tb int\n}\nvar x [2]T = [2]T{{1,2},{3,4}}\nreturn x[0].a")
// 	assert(4, f212(), "f212:\ntype T struct {\n\ta int\n\tb int\n}\nvar x [2]T = [2]T{{1,2},{3,4}}\nreturn x[1].b")

// 	assert(3, g3, "g3")
// 	assert(4, g4, "g4")
// 	assert(5, g5, "g5")
// 	assert(6, g6, "g6")
// 	assert(5, *g7, "*g7")
// 	assert(0, strcmp(f213(), "abc"), "strcmp(f213(),\"abc\")")
// 	assert(0, strcmp(f214(), "abc"), "strcmp(f214(), \"abc\")")
// 	assert(0, strcmp(g8, "abc"), "strcmp(g8, \"abc\")")

// 	assert(0, g9[0], "g9[0]")
// 	assert(1, g9[1], "g9[1]")
// 	assert(2, g9[2], "g9[2]")

// 	assert(0, strcmp(g10[0], "foo"), "strcmp(g10[0], \"foo\")")
// 	assert(0, strcmp(g10[1], "bar"), "strcmp(g10[1], \"bar\")")
// 	assert(0, g10[1][3], "g10[1][3]")
// 	assert(2, Sizeof(g10) / Sizeof(*g10), "Sizeof(g10) / Sizeof(*g10)")

// 	assert(1, g11[0].a, "g11[0].a")
// 	assert(2, g11[0].b, "g11[0].b")
// 	assert(3, g11[1].a, "g11[1].a")
// 	assert(4, g11[1].b, "g11[1].b")

// 	assert(1, g12[0].a[0], "g12[0].a[0]")
// 	assert(2, g12[0].a[1], "g12[0].a[1]")
// 	assert(3, g12[1].a[0], "g12[1].a[0]")
// 	assert(4, g12[1].a[1], "g12[1].a[1]")

// 	assert('a', g8_1[0], "g8_1[0]")

// 	assert(1, g12_1[0].a[0], "g12_1[0].a[0]")
// 	assert(2, g12_1[0].a[1], "g12_1[0].a[1]")
// 	assert(3, g12_1[1].a[0], "g12_1[1].a[0]")
// 	assert(4, g12_1[1].a[1], "g12_1[1].a[1]")

// 	assert(0, Sizeof(g12_2), "Sizeof(g12_2)")
// 	assert(0, Sizeof(f215()), "f215:\nvar x []byte\nreturn x")

// 	assert(3, f216(), "f216\nvar x int\nvar y int\nx=1,y=2\nreturn x+y")

// 	assert(6, f217(), "f217:\nvar x [2][3]int = [2][3]int{{1,2,3},{4,5,6}}\nreturn x[1][2]")

// 	assert(4, g12_3[1].a[1], "g12_3[1].a[1]")

// 	// type cast
// 	assert(1, f218(), "f218:\nvar x byte\nreturn int(x)")
// 	assert(131585, f219(), "f219:\nvar x int64 = 8590066177\nreturn int(x)")
// 	assert(513, f220(), "f220:\nvar x int64 = 8590066177\nreturn int16(x)")
// 	assert(1, f221(), "f221:\nvar x int64 = 8590066177\nreturn byte(x)")

// 	printf("\nOK\n")
// }

// func fib(x int) int {
// 	if x <= 1 {
// 		return 1
// 	}
// 	return fib(x-1) + fib(x-2)
// }
