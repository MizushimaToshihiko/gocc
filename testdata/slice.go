package test_slice

func assert(want int, act int, code string)
func println(format ...string)

#include "test.h"

var g01 = 1
var g02 = 3

func strcmp(s1, s2 string) int

func ret3() int {
	return 3
}

func retf3() float64 {
	return 3.5
}

func reti(i int) int {
	return i
}

func printf(format ...string)

var g03 = []string{"abc", "def", "ghi"}
var g04 = [][]string{{"abc", "def", "ghi"}, {"jkl", "mno", "pqr"}}

var a038 = [5]int{0, 0, 0, 0, 0}

func main() {
	var a01 = [2]int{1, 2}
	var s01 = a01[0:2]
	ASSERT(1, s01[0])
	ASSERT(2, s01[1])
	ASSERT(8, Sizeof(s01))
	ASSERT(2, len(s01))
	ASSERT(2, cap(s01))

	var a02 = [6]int{1, 2, 3, 4, 5, 6}
	var s02 = a02[2:5]
	ASSERT(3, s02[0])
	ASSERT(4, s02[1])
	ASSERT(5, s02[2])
	ASSERT(8, Sizeof(s02))
	ASSERT(3, len(s02))
	ASSERT(4, cap(s02))
	s02[0] = 100
	ASSERT(100, a02[2])

	var x021, x0211 = 2.0, 5.0
	s021 := a02[x021:x0211]
	ASSERT(100, s021[0])
	ASSERT(4, s021[1])
	ASSERT(5, s021[2])
	ASSERT(8, Sizeof(s021))
	ASSERT(3, len(s021))
	ASSERT(4, cap(s021))

	var x022, x0221 = 2.0, 5.0
	s022 := a02[x022:x0221]
	ASSERT(100, s022[0])
	ASSERT(4, s022[1])
	ASSERT(5, s022[2])
	ASSERT(8, Sizeof(s022))
	ASSERT(3, len(s022))
	ASSERT(4, cap(s022))

	s023 := a02[g01 : ret3()+3]
	ASSERT(2, s023[0])
	ASSERT(100, s023[1])
	ASSERT(5, len(s023))
	ASSERT(5, cap(s023))

	s024 := a02[g01 : retf3()+3]
	ASSERT(2, s024[0])
	ASSERT(100, s024[1])
	ASSERT(5, len(s024))
	ASSERT(5, cap(s024))

	var s025 []int
	ASSERT(0, len(s025))
	ASSERT(0, cap(s025))
	ASSERT(8, Sizeof(s025))

	var s026 []int = []int{1, 2, 3, 4, 5, 6}
	ASSERT(6, len(s026))
	ASSERT(6, cap(s026))
	ASSERT(8, Sizeof(s026))
	ASSERT(1, s026[0])
	ASSERT(2, s026[1])
	ASSERT(3, s026[2])
	ASSERT(4, s026[3])
	ASSERT(5, s026[4])
	ASSERT(6, s026[5])
	s026[0], s026[1], s026[2], s026[3], s026[4], s026[5] = 100, 101, 102, 103, 104, 105
	ASSERT(100, s026[0])
	ASSERT(101, s026[1])
	ASSERT(102, s026[2])
	ASSERT(103, s026[3])
	ASSERT(104, s026[4])
	ASSERT(105, s026[5])

	var s027 = []int{1, 2, 3, 4, 5, 6}
	ASSERT(6, len(s027))
	ASSERT(6, cap(s027))
	ASSERT(8, Sizeof(s027))
	ASSERT(1, s027[0])
	ASSERT(2, s027[1])
	ASSERT(3, s027[2])
	ASSERT(4, s027[3])
	ASSERT(5, s027[4])
	ASSERT(6, s027[5])

	s028 := []string{"abc", "def", "ghi"}
	ASSERT(8, Sizeof(s028))
	ASSERT(3, len(s028))
	ASSERT(3, cap(s028))
	ASSERT(0, strcmp(s028[0], "abc"))
	ASSERT(0, strcmp(s028[1], "def"))
	ASSERT(0, strcmp(s028[2], "ghi"))

	s029 := [][]string{{"abc", "def", "ghi"}, {"jkl", "mno", "pqr"}}
	ASSERT(8, Sizeof(s029))
	ASSERT(2, len(s029))
	ASSERT(3, len(s029[0]))
	ASSERT(3, len(s029[1]))
	ASSERT(2, cap(s029))
	ASSERT(3, cap(s029[0]))
	ASSERT(3, cap(s029[1]))
	ASSERT(0, strcmp(s029[0][0], "abc"))
	ASSERT(0, strcmp(s029[0][1], "def"))
	ASSERT(0, strcmp(s029[0][2], "ghi"))
	println(s029[0][2])
	ASSERT(0, strcmp(s029[1][0], "jkl"))
	ASSERT(0, strcmp(s029[1][1], "mno"))
	ASSERT(0, strcmp(s029[1][2], "pqr"))

	ASSERT(8, Sizeof(g03))
	ASSERT(3, len(g03))
	ASSERT(3, cap(g03))
	ASSERT(0, strcmp(g03[0], "abc"))
	ASSERT(0, strcmp(g03[1], "def"))
	ASSERT(0, strcmp(g03[2], "ghi"))

	ASSERT(8, Sizeof(g04))
	ASSERT(2, len(g04))
	ASSERT(3, len(g04[0]))
	ASSERT(3, len(g04[1]))
	ASSERT(2, cap(g04))
	ASSERT(3, cap(g04[0]))
	ASSERT(3, cap(g04[1]))
	ASSERT(0, strcmp(g04[0][0], "abc"))
	ASSERT(0, strcmp(g04[0][1], "def"))
	ASSERT(0, strcmp(g04[0][2], "ghi"))
	println(g04[0][2])
	ASSERT(0, strcmp(g04[1][0], "jkl"))
	ASSERT(0, strcmp(g04[1][1], "mno"))
	ASSERT(0, strcmp(g04[1][2], "pqr"))

	var s030 = []struct {
		a []int
		b int
		c []string
	}{
		{
			a: []int{1, 2, 3, 4, 5, 6},
			b: 7,
			c: []string{"abc", "def", "ghi"},
		},
		{
			a: []int{8, 9, 10, 11, 12, 13},
			b: 14,
			c: []string{"jkl", "mno", "pqr"},
		},
	}
	ASSERT(1, s030[0].a[0])
	ASSERT(2, s030[0].a[1])
	ASSERT(3, s030[0].a[2])
	ASSERT(4, s030[0].a[3])
	ASSERT(5, s030[0].a[4])
	ASSERT(6, s030[0].a[5])
	ASSERT(7, s030[0].b)
	ASSERT(0, strcmp(s030[0].c[0], "abc"))
	ASSERT(0, strcmp(s030[0].c[1], "def"))
	ASSERT(0, strcmp(s030[0].c[2], "ghi"))
	ASSERT(8, s030[1].a[0])
	ASSERT(9, s030[1].a[1])
	ASSERT(10, s030[1].a[2])
	ASSERT(11, s030[1].a[3])
	ASSERT(12, s030[1].a[4])
	ASSERT(13, s030[1].a[5])
	ASSERT(14, s030[1].b)
	ASSERT(0, strcmp(s030[1].c[0], "jkl"))
	ASSERT(0, strcmp(s030[1].c[1], "mno"))
	ASSERT(0, strcmp(s030[1].c[2], "pqr"))

	var s031 = []func() int{ret3, ret3}
	ASSERT(2, len(s031))
	ASSERT(2, cap(s031))
	ASSERT(3, s031[0]())
	ASSERT(3, s031[1]())

	var s0311 = []func(int) int{reti, reti}
	ASSERT(2, len(s0311))
	ASSERT(2, cap(s0311))
	ASSERT(3, s0311[0](3))
	ASSERT(4, s0311[1](4))

	var s032 = []float32{1., 2., 3.}
	ASSERT(3, len(s032))
	ASSERT(3, cap(s032))
	ASSERT(1, s032[0] == 1.0)
	ASSERT(1, s032[1] == 2.0)

	var s033 = []float64{1., 2., 3.}
	ASSERT(3, len(s033))
	ASSERT(3, cap(s033))
	ASSERT(1, s033[0] == 1.0)
	ASSERT(1, s033[1] == 2.0)

	s034 := make([]int, 10)
	ASSERT(10, len(s034))
	ASSERT(10, cap(s034))
	ASSERT(0, s034[0])
	ASSERT(0, s034[9])
	s034[1],
		s034[2],
		s034[3],
		s034[4],
		s034[5],
		s034[6],
		s034[7],
		s034[8],
		s034[9] =
		101,
		102,
		103,
		104,
		105,
		106,
		107,
		108,
		109
	ASSERT(101, s034[1])
	ASSERT(102, s034[2])
	ASSERT(103, s034[3])
	ASSERT(104, s034[4])
	ASSERT(105, s034[5])
	ASSERT(106, s034[6])
	ASSERT(107, s034[7])
	ASSERT(108, s034[8])
	ASSERT(109, s034[9])

	s035 := make([]string, 10, 15)
	ASSERT(10, len(s035))
	ASSERT(15, cap(s035))
	ASSERT(0, strcmp(s035[0], ""))
	ASSERT(0, strcmp(s035[9], ""))
	s035[1],
		s035[2],
		s035[3],
		s035[4],
		s035[5],
		s035[6],
		s035[7],
		s035[8],
		s035[9] = "abc",
		"def",
		"ghi",
		"jkl",
		"mno",
		"pqr",
		"stu",
		"vwx",
		"yz"
	ASSERT(0, strcmp(s035[1], "abc"))
	ASSERT(0, strcmp(s035[2], "def"))
	ASSERT(0, strcmp(s035[3], "ghi"))
	ASSERT(0, strcmp(s035[4], "jkl"))
	ASSERT(0, strcmp(s035[5], "mno"))
	ASSERT(0, strcmp(s035[6], "pqr"))
	ASSERT(0, strcmp(s035[7], "stu"))
	ASSERT(0, strcmp(s035[8], "vwx"))
	ASSERT(0, strcmp(s035[9], "yz"))

	s036 := make([]int, 1, 5)
	ASSERT(1, len(s036))
	ASSERT(5, cap(s036))
	ASSERT(0, s036[0])
	s036 = append(s036, 1, 2, 3)
	ASSERT(4, len(s036))
	ASSERT(5, cap(s036))
	ASSERT(0, s036[0])
	ASSERT(1, s036[1])
	ASSERT(2, s036[2])
	ASSERT(3, s036[3])
	ASSERT(1, s036[reti(1)])

	s037 := make([]string, 1, 5)
	ASSERT(1, len(s037))
	ASSERT(5, cap(s037))
	ASSERT(0, strcmp(s037[0], ""))
	s037 = append(s037, "abc", "def", "ghi")
	ASSERT(4, len(s037))
	ASSERT(5, cap(s037))
	ASSERT(0, strcmp(s037[0], ""))
	ASSERT(0, strcmp(s037[1], "abc"))
	ASSERT(0, strcmp(s037[2], "def"))
	ASSERT(0, strcmp(s037[3], "ghi"))

	// case: made a slice from global array, int
	a0381 := [5]int{10, 20, 30, 40, 50}
	s038 := a038[0:2]
	s0381 := a0381[0:2]
	ASSERT(2, len(s038))
	ASSERT(5, cap(s038))
	s038[0] = s0381[0]
	s038[1] = a0381[4]
	ASSERT(10, s038[0])
	ASSERT(50, s038[1])
	s038 = append(s038, 2, 3, 4, 5)
	ASSERT(6, len(s038))
	ASSERT(14, cap(s038))
	ASSERT(10, s038[0])
	ASSERT(50, s038[1])
	ASSERT(2, s038[2])
	ASSERT(3, s038[3])
	ASSERT(4, s038[4])
	ASSERT(5, s038[5])
	ASSERT(10, a038[0])
	ASSERT(50, a038[1])
	ASSERT(0, a038[2])
	ASSERT(0, a038[3])
	ASSERT(0, a038[4])

	// case: made a slice from global array, string
	s0382 := g03[0:2]
	ASSERT(2, len(s0382))
	ASSERT(3, cap(s0382))
	s0382[0] = "aaa"
	s0382[1] = "bbb"
	ASSERT(0, strcmp(s0382[0], "aaa"))
	ASSERT(0, strcmp(s0382[1], "bbb"))
	s0382 = append(s0382, "jkl", "mno", "pqr", "stu")
	ASSERT(6, len(s0382))
	ASSERT(10, cap(s0382))
	ASSERT(0, strcmp(s0382[0], "aaa"))
	ASSERT(0, strcmp(s0382[1], "bbb"))
	ASSERT(0, strcmp(s0382[2], "jkl"))
	ASSERT(0, strcmp(s0382[3], "mno"))
	ASSERT(0, strcmp(s0382[4], "pqr"))
	ASSERT(0, strcmp(s0382[5], "stu"))

	// case: made a slice by initializer, string
	s03821 := []string{"abc", "def", "ghi"}
	ASSERT(3, len(s03821))
	ASSERT(3, cap(s03821))
	s03821[0] = "aaa"
	s03821[1] = "bbb"
	ASSERT(0, strcmp(s03821[0], "aaa"))
	ASSERT(0, strcmp(s03821[1], "bbb"))
	s03821 = append(s03821, "jkl", "mno", "pqr", "stu")
	ASSERT(7, len(s03821))
	ASSERT(10, cap(s03821))
	ASSERT(0, strcmp(s03821[0], "aaa"))
	ASSERT(0, strcmp(s03821[1], "bbb"))
	ASSERT(0, strcmp(s03821[2], "ghi"))
	ASSERT(0, strcmp(s03821[3], "jkl"))
	ASSERT(0, strcmp(s03821[4], "mno"))
	ASSERT(0, strcmp(s03821[5], "pqr"))
	ASSERT(0, strcmp(s03821[6], "stu"))

	// case: made a global slice , string, two dimensions
	ASSERT(2, len(g04))
	ASSERT(2, cap(g04))
	ASSERT(3, len(g04[0]))
	ASSERT(3, cap(g04[0]))
	ASSERT(3, len(g04[1]))
	ASSERT(3, cap(g04[1]))
	g04[0][2] = "ggg"
	s038211 := []string{"aaa", "bbb"}
	g04 = append(g04, s038211)
	g04[1] = append(g04[1], "stu")
	ASSERT(3, len(g04))
	ASSERT(5, cap(g04))
	ASSERT(4, len(g04[0]))
	ASSERT(7, cap(g04[0]))
	ASSERT(4, len(g04[1]))
	ASSERT(7, cap(g04[1]))
	ASSERT(4, len(g04[2]))
	ASSERT(0, strcmp(g04[0][0], "abc"))
	ASSERT(0, strcmp(g04[0][1], "def"))
	ASSERT(0, strcmp(g04[0][2], "ggg"))
	ASSERT(0, strcmp(g04[1][0], "jkl"))
	ASSERT(0, strcmp(g04[1][1], "mno"))
	ASSERT(0, strcmp(g04[1][2], "pqr"))
	ASSERT(0, strcmp(g04[1][3], "stu"))
	ASSERT(0, strcmp(g04[2][0], "aaa"))
	ASSERT(0, strcmp(g04[2][1], "bbb"))

	// case: made a slice by initializer, int
	s03822 := []int{1, 2, 3}
	ASSERT(3, len(s03822))
	ASSERT(3, cap(s03822))
	s03822[0] = 11
	s03822[1] = 22
	ASSERT(11, s03822[0])
	ASSERT(22, s03822[1])
	s03822 = append(s03822, 4, 5, 6, 7)
	ASSERT(7, len(s03822))
	ASSERT(10, cap(s03822))
	ASSERT(11, s03822[0])
	ASSERT(22, s03822[1])
	ASSERT(3, s03822[2])
	ASSERT(4, s03822[3])
	ASSERT(5, s03822[4])
	ASSERT(6, s03822[5])
	ASSERT(7, s03822[6])

	// case: made a slice by 'var', int
	var s03823 []int
	ASSERT(0, len(s03823))
	ASSERT(0, cap(s03823))
	s03823 = append(s03823, 1, 2, 3, 4, 5, 6)
	ASSERT(6, len(s03823))
	ASSERT(6, cap(s03823))
	ASSERT(1, s03823[0])
	ASSERT(2, s03823[1])
	ASSERT(3, s03823[2])
	ASSERT(4, s03823[3])
	ASSERT(5, s03823[4])
	ASSERT(6, s03823[5])
	s03823[2] = 1000
	ASSERT(1000, s03823[2])
	s03823 = append(s03823, 7, 8, 9, 10, 11, 12)
	ASSERT(12, len(s03823))
	ASSERT(18, cap(s03823))
	ASSERT(1, s03823[0])
	ASSERT(2, s03823[1])
	ASSERT(1000, s03823[2])
	ASSERT(4, s03823[3])
	ASSERT(5, s03823[4])
	ASSERT(6, s03823[5])
	ASSERT(7, s03823[6])
	ASSERT(8, s03823[7])
	ASSERT(9, s03823[8])
	ASSERT(10, s03823[9])
	ASSERT(11, s03823[10])
	ASSERT(12, s03823[11])

	// case: made a slice by make function, int
	s0383 := make([]int, 5)
	ASSERT(5, len(s0383))
	ASSERT(5, cap(s0383))
	s0383[0] = 10
	s0383[1] = 50
	ASSERT(10, s0383[0])
	ASSERT(50, s0383[1])
	s0383 = append(s0383, 2, 3, 4, 5)
	ASSERT(9, len(s0383))
	ASSERT(14, cap(s0383))
	ASSERT(10, s0383[0])
	ASSERT(50, s0383[1])
	ASSERT(0, s0383[2])
	ASSERT(0, s0383[3])
	ASSERT(0, s0383[4])
	ASSERT(2, s0383[5])
	ASSERT(3, s0383[6])
	ASSERT(4, s0383[7])
	ASSERT(5, s0383[8])

	// case: made a slice by make function, string
	s039 := []string{"abc", "def"}
	ASSERT(2, len(s039))
	ASSERT(2, cap(s039))
	ASSERT(0, strcmp(s039[0], "abc"))
	ASSERT(0, strcmp(s039[1], "def"))
	s039 = append(s039, "ghi", "jkl", "mno")
	ASSERT(0, strcmp(s039[0], "abc"))
	ASSERT(0, strcmp(s039[1], "def"))
	ASSERT(0, strcmp(s039[2], "ghi"))
	ASSERT(0, strcmp(s039[3], "jkl"))
	ASSERT(0, strcmp(s039[4], "mno"))

	s040 := make([]string, 2, 3)
	ASSERT(2, len(s040))
	ASSERT(3, cap(s040))
	s040[0] = "abc"
	s040[1] = "def"
	ASSERT(0, strcmp(s040[0], "abc"))
	ASSERT(0, strcmp(s040[1], "def"))
	s040 = append(s040, "ghi", "jkl", "mno")
	ASSERT(0, strcmp(s040[0], "abc"))
	ASSERT(0, strcmp(s040[1], "def"))
	ASSERT(0, strcmp(s040[2], "ghi"))
	ASSERT(0, strcmp(s040[3], "jkl"))
	ASSERT(0, strcmp(s040[4], "mno"))

	var a041 = [8]int{0, 1, 2, 3, 4, 5, 6, 7}
	var s041 = make([]int, 6)
	// var s0411 = a041[0:8]
	n0411 := copy(s041, a041[0:]) // n0411 == 6, s041 == []int{0, 1, 2, 3, 4, 5}
	ASSERT(0, s041[0])
	ASSERT(1, s041[1])
	ASSERT(2, s041[2])
	ASSERT(3, s041[3])
	ASSERT(4, s041[4])
	ASSERT(5, s041[5])
	n0412 := copy(s041, s041[2:]) // n0412 == 4, s041 == []int{2, 3, 4, 5, 4, 5}
	ASSERT(6, n0411)
	ASSERT(4, n0412)
	ASSERT(2, s041[0])
	ASSERT(3, s041[1])
	ASSERT(4, s041[2])
	ASSERT(5, s041[3])
	ASSERT(4, s041[4])
	ASSERT(5, s041[5])
	a041[0] = 1000
	ASSERT(2, s041[0])

	var a042 = [6]string{"abc", "def", "ghi", "jkl", "mno", "pqr"}
	var s042 = make([]string, 4)
	n0421 := copy(s042, a042[0:])
	ASSERT(4, n0421)
	ASSERT(0, strcmp(s042[0], "abc"))
	ASSERT(0, strcmp(s042[1], "def"))
	ASSERT(0, strcmp(s042[2], "ghi"))
	ASSERT(0, strcmp(s042[3], "jkl"))
	n0422 := copy(s042, s042[2:])
	ASSERT(2, n0422)
	ASSERT(0, strcmp(s042[0], "ghi"))
	ASSERT(0, strcmp(s042[1], "jkl"))
	ASSERT(0, strcmp(s042[2], "ghi"))
	ASSERT(0, strcmp(s042[3], "jkl"))

	s043 := make([][]int, 6)
	ASSERT(6, len(s043))
	ASSERT(0, len(s043[0]))
	s043[0] = append(s043[0], 1)
	ASSERT(1, len(s043[0]))
	ASSERT(1, s043[0][0])
	ASSERT(1, len(s043[1])) // 他も変わっちゃう

	println("OK")
}
