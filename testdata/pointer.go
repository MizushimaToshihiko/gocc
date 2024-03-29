package test_pointer

func assert(want int, act int, code string)
func println(format ...string)

#include "test.h"

func main() {

	/*
		#include <stdio.h>
		int main(void){
		    int x[2] = {1, 2};
		    int (*y)[2] = &x;

		    printf("%p\n", y[1]);
		    printf("%d\n", (*y)[1]);

		    typedef struct t01 t01;
		    struct t01 {
		        int a;
		        int b;
		    };

		    t01 x01[2] = {{1,2},{3,4}};
		    t01 (*y01)[2] = &x01;
		    printf("%p\n", y01[0]);
		    printf("%p\n", &x01[0]);
		    printf("%p\n", y01[1]);
		    printf("%d\n", (*y01)[0]);
		    printf("%d\n", (*y01)[1]);
		}
	*/
	type x01 struct {
		a int
		b int
	}
	var y01 = &x01{1, 2}
	ASSERT(1, y01.a)
	var y02 = &[2]x01{{1, 2}, {3, 4}}
	// ASSERT(1, y02[0].a)
	// ASSERT(2, y02[0].b)
	// ASSERT(3, y02[1].a) //
	// ASSERT(4, y02[1].b) //
	ASSERT(1, (*y02)[0].a)
	ASSERT(2, (*y02)[0].b)
	ASSERT(3, (*y02)[1].a)
	ASSERT(4, (*y02)[1].b)

	var x03 = [2]int{1, 2}
	var y03 *[2]int = &x03
	ASSERT(1, (*y03)[0])
	ASSERT(2, (*y03)[1])
	// ASSERT(2, y03[1]) // y03[0] and y03[1] is pointer address?

	var y031 = &x03
	// ASSERT(1, y031[0])       // is pointer address?
	// ASSERT(2, y031[1])       // is pointer address?
	ASSERT(1, (*y031)[0]) //
	ASSERT(2, (*y031)[1]) //

	var x1 int = 3
	ASSERT(3, *&x1)
	var x2 int = 3
	var y2 *int = &x2
	var z2 **int = &y2
	ASSERT(3, **z2)
	var x3 int = 3
	var y3 int = 5
	ASSERT(5, *(&x3 + 1))       // Not supported in Go.
	ASSERT(3, *(&y3 - 1))       // Not supported in Go.
	ASSERT(5, *(&x3 - (-1))) // Not supported in Go.
	var x4 int = 3
	var y4 *int = &x4
	*y4 = 5
	ASSERT(5, x4)
	var x5 int = 3
	var y5 int = 5
	*(&x5 + 1) = 7
	ASSERT(7, y5)
	var x6 int = 3
	var y6 int = 5
	*(&y6 - 2 + 1) = 7
	ASSERT(7, x6)
	var x7 int = 3
	ASSERT(5, (&x7+2)-&x7+3)
	var x8 [2]int
	var y8 *int = &x8
	*y8 = 3
	ASSERT(3, *x8)
	var x9 [3]int
	*x9 = 3
	*(x9 + 1) = 4
	*(x9 + 2) = 5
	ASSERT(3, *x9)
	ASSERT(4, *(x9 + 1))
	ASSERT(5, *(x9 + 2))
	var x10 [2][3]int
	var y10 *int = x10
	*y10 = 0
	ASSERT(0, **x10)
	*(y10 + 1) = 1
	ASSERT(1, *(*x10 + 1))
	*(y10 + 2) = 2
	ASSERT(2, *(*x10 + 2))
	*(y10 + 3) = 3
	ASSERT(3, **(x10 + 1))
	*(y10 + 4) = 4
	ASSERT(4, *(*(x10 + 1) + 1))
	*(y10 + 5) = 5
	ASSERT(5, *(*(x10 + 1) + 2))
	var x11 [3]int
	*x11 = 3
	x11[1] = 4
	x11[2] = 5
	ASSERT(3, *x11)
	ASSERT(4, *(x11 + 1))
	ASSERT(5, *(x11 + 2))
	2[x11] = 5
	ASSERT(5, *(x11 + 2))
	var x12 [2][3]int
	var y12 *int = x12
	y12[0] = 0
	ASSERT(0, x12[0][0])
	y12[1] = 1
	ASSERT(1, x12[0][1])
	y12[2] = 2
	ASSERT(2, x12[0][2])
	y12[3] = 3
	ASSERT(3, x12[1][0])
	y12[4] = 4
	ASSERT(4, x12[1][1])
	y12[5] = 5
	ASSERT(5, x12[1][2])

	println("OK")
}
