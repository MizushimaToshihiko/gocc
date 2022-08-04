package test_control

func assert(want int, act int, code string)
func println(format ...string)

#include "test.h"
#include <stdbool.h>

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
		if false {
			x1 = 2
		} else {
			x1 = 3
		}
		ASSERT(3, x1)
		if 1 - 1 {
			x1 = 2
		} else {
			x1 = 3
		}
		ASSERT(3, x1)
		if true {
			x1 = 2
		} else {
			x1 = 3
		}
		ASSERT(2, x1)
		if 2 - 1 {
			x1 = 2
		} else {
			x1 = 3
		}
		ASSERT(2, x1)

		var i int = 0
		var j int = 0
		for i = 0; i <= 10; i = i + 1 {
			j = i + j
		}
		ASSERT(55, j)
		var j int = 0
		for i := 0; i <= 10; i = i + 1 {
			j = i + j
		}
		ASSERT(55, j)
		i = 0
		for i < 10 {
			i = i + 1
		}
		ASSERT(10, i)

		i = 1
		{
			i = 2
		}
		i = 3
		ASSERT(3, i)

		i = 0
		for i < 10 {
			i = i + 1
		}
		ASSERT(10, i)
		i = 0
		j = 0
		for i <= 10 {
			j = i + j
			i = i + 1
		}
		ASSERT(55, j)

		// ASSERT(3, (1,2,3))
		// i=2, j=3; (i=5,j)=6;
		// ASSERT(5, i)
		// i=2, j=3; (i=5,j)=6;
		// ASSERT(6, j)

		ASSERT(1, 0 || 1)
		ASSERT(1, 0 || (2-2) || 5)
		ASSERT(0, 0 || 0)
		ASSERT(0, 0 || (2-2))

		ASSERT(0, 0 && 1)
		ASSERT(0, (2-2) && 5)
		ASSERT(1, 1 && 5)

		i = 0
		goto a
	a:
		i++
	b:
		i++
	c:
		i++
		ASSERT(3, i)
		i = 0
		goto e
	d:
		i++
	e:
		i++
	f:
		i++
		ASSERT(2, i)
		i = 0
		goto i
	g:
		i++
	h:
		i++
	i:
		i++
		i
		ASSERT(1, i)

		type foo int
		var x2 foo
		goto foo
		x2 = 2
	foo:
		x2 = 1
		ASSERT(1, x2)

		i = 0
		for ; i < 10; i++ {
			if i == 3 {
				break
			}
		}
		ASSERT(3, i)
		i = 0
		for 1 {
			i++
			if i >= 3 {
				i++
				break
			}
		}
		ASSERT(4, i)
		i = 0
		for ; i < 10; i++ {
			for {
				break
			}
			if i == 3 {
				break
			}
		}
		ASSERT(3, i)
		i = 0
		j = 0
		for ; i < 10; i++ {
			if i > 5 {
				continue
			}
			j++
		}
		ASSERT(10, i)
		i = 0
		j = 0
		for ; i < 10; i++ {
			if i > 5 {
				continue
			}
			j++
		}
		ASSERT(6, j)
		i = 0
		j = 0
		for !i {
			for ; j != 10; j++ {
				continue
			}
			break
		}
		ASSERT(10, j)
		i = 0
		j = 0
		for i < 10 {
			i++
			if i > 5 {
				continue
			}
			j++
		}
		ASSERT(10, i)
		i = 0
		j = 0
		for i < 10 {
			i++
			if i > 5 {
				continue
			}
			j++
		}
		ASSERT(5, j)
		i = 0
		j = 0
		for !i {
			for j != 10 {
				j++
				continue
			}
			break
		}
		ASSERT(10, j)

		ASSERT(5, switchFn(0))
		ASSERT(5, switchFn(3))
		ASSERT(5, switchFn(4))
		ASSERT(6, switchFn(1))
		ASSERT(6, switchFn(2))
		ASSERT(100, switchFn(5))
		ASSERT(100, switchFn(6))
		ASSERT(10, switchFn(8))
		ASSERT(10, switchFn(9))
		ASSERT(10, switchFn(10))
		ASSERT(10, switchFn(11))

		i = 0
		switch i {
		case 0, 3:
			i = 5
		case 1:
			i = 6
		case 2:
			i = 7
		}
		ASSERT(5, i)
		i = 1
		switch i {
		case 0:
			i = 5
		case 1:
			i = 6
		case 2:
			i = 7
		}
		ASSERT(6, i)
		i = 2
		switch i {
		case 0:
			i = 5
		case 1:
			i = 6
		case 2:
			i = 7
		}
		ASSERT(7, i)
		i = 3
		switch i {
		case 0:
			i = 5
		case 1:
			i = 6
		case 2:
			i = 7
		}
		ASSERT(3, i)
		i = 0
		switch i {
		case 0:
			i = 5
		default:
			i = 7
		}
		ASSERT(5, i)
		i = 2
		switch i {
		case 0:
			i = 5
		default:
			i = 7
		}
		ASSERT(7, i)
		i = 0
		switch -1 {
		case 0xffffffff:
			i = 3
		}
		ASSERT(3, i)

		ASSERT(0, 0.0 && 0.0);
		ASSERT(0, 0.0 && 0.1);
		ASSERT(0, 0.3 && 0.0);
		ASSERT(1, 0.3 && 0.5);
		ASSERT(0, 0.0 || 0.0);
		ASSERT(1, 0.0 || 0.1);
		ASSERT(1, 0.3 || 0.0);
		ASSERT(1, 0.3 || 0.5);
		var x2 int; if 0.0 {x2=3;}else{x2=5;};
		ASSERT(5, x2);
		var x3 int; if 0.1 {x3=3;}else{x3=5;};
		ASSERT(3, x3);
		var x4=5; if 0.0{x4=3;};
		ASSERT(5, x4);
		var x5=5; if 0.1{x5=3;};
		ASSERT(3, x5);
		i=10.0; j=0; for ;i!=0;i--,j++{};
		ASSERT(10, j);
		i=10.0; j=0; for i!=0{i--;j++;};
		ASSERT(10, j);

		var x6, y6 = 1, 2
		var z6 int
		switch {
		case x6 < y6: z6=switchFn(x6)
		case x6 > y6: z6=switchFn(y6)
		case x6 == y6: z6=switchFn(x6+y6)
		}
		ASSERT(6, z6)
		x6, y6 = 2, 0
		z6 = 0
		switch {
		case x6 < y6: z6=switchFn(x6)
		case x6 > y6: z6=switchFn(y6)
		case x6 == y6: z6=switchFn(x6+y6)
		}
		ASSERT(5, z6)
		x6, y6 = 3, 3
		z6 = 0
		switch {
		case x6 < y6: z6=switchFn(x6)
		case x6 > y6: z6=switchFn(y6)
		case x6 == y6: z6=switchFn(x6+y6)
		}
		ASSERT(100, z6)

	var x61, y61 float32 = 1.0, 2.0
	println("x61: %f", x61)
	println("y61: %f", y61)
	var z61 int
	switch {
	case x61 < y61:
		println("x61<y61")
		z61 = switchFn(x61)
	case x61 > y61:
		z61 = switchFn(y61)
	case x61 == y61:
		z61 = switchFn(x61 + y61)
	}
	ASSERT(6, z61)
	x61, y61 = 2.0, 0.0
	println("x61: %f", x61)
	println("y61: %f", y61)
	z61 = 0
	switch {
	case x61 < y61:
		z61 = switchFn(x61)
	case x61 > y61:
		z61 = switchFn(y61)
	case x61 == y61:
		z61 = switchFn(x61 + y61)
	}
	ASSERT(5, z61)
	x61, y61 = 3.0, 3.0
	println("x61: %f", x61)
	println("y61: %f", y61)
	z61 = 0
	switch {
	case x61 < y61:
		z61 = switchFn(x61)
	case x61 > y61:
		z61 = switchFn(y61)
	case x61 == y61:
		z61 = switchFn(x61 + y61)
	}
	ASSERT(100, z61)

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
	ASSERT(1, z7)
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
	ASSERT(2, z7)
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
	ASSERT(3, z7)
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
	ASSERT(4, z7)

	var z8 int
	if x8, y8 := switchFn(0), 8; x8 < y8 {
		z8 = 1
	} else if x8 > z8 {
		z8 = 2
	} else {
		z8 = 3
	}
	ASSERT(1, z8)
	z8 = 0
	if x8, y8 := switchFn(200), 8; x8 < y8 {
		z8 = 1
	} else if x8 > z8 {
		z8 = 2
	} else {
		z8 = 3
	}
	ASSERT(2, z8)
	z8 = 10
	if x8, y8 := switchFn(1)+3, 8; x8 < y8 {
		z8 = 1
	} else if x8 > z8 {
		z8 = 2
	} else {
		z8 = 3
	}
	ASSERT(3, z8)

	println("OK")
}
