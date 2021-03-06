package test_slice

func assert(want int, act int, code string)
func println(format ...string)

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
	assert(1, s01[0], "s01[0]")
	assert(2, s01[1], "s01[1]")
	assert(8, Sizeof(s01), "Sizeof(x01)")
	assert(2, len(s01), "len(s01)")
	assert(2, cap(s01), "cap(s01)")

	var a02 = [6]int{1, 2, 3, 4, 5, 6}
	var s02 = a02[2:5]
	assert(3, s02[0], "s02[0]")
	assert(4, s02[1], "s02[1]")
	assert(5, s02[2], "s02[2]")
	assert(8, Sizeof(s02), "Sizeof(x02)")
	assert(3, len(s02), "len(s02)")
	assert(4, cap(s02), "cap(s02)")
	s02[0] = 100
	assert(100, a02[2], "a02[2]")

	var x021, x0211 = 2.0, 5.0
	s021 := a02[x021:x0211]
	assert(100, s021[0], "s021[0]")
	assert(4, s021[1], "s021[1]")
	assert(5, s021[2], "s021[2]")
	assert(8, Sizeof(s021), "Sizeof(x021)")
	assert(3, len(s021), "len(s021)")
	assert(4, cap(s021), "cap(s021)")

	var x022, x0221 = 2.0, 5.0
	s022 := a02[x022:x0221]
	assert(100, s022[0], "s022[0]")
	assert(4, s022[1], "s022[1]")
	assert(5, s022[2], "s022[2]")
	assert(8, Sizeof(s022), "Sizeof(x022)")
	assert(3, len(s022), "len(s022)")
	assert(4, cap(s022), "cap(s022)")

	s023 := a02[g01 : ret3()+3]
	assert(2, s023[0], "s023[0]")
	assert(100, s023[1], "s023[1]")
	assert(5, len(s023), "len(s023)")
	assert(5, cap(s023), "cap(s023)")

	s024 := a02[g01 : retf3()+3]
	assert(2, s024[0], "s024[0]")
	assert(100, s024[1], "s024[1]")
	assert(5, len(s024), "len(s024)")
	assert(5, cap(s024), "cap(s024)")

	var s025 []int
	assert(0, len(s025), "len(s025)")
	assert(0, cap(s025), "cap(s025)")
	assert(8, Sizeof(s025), "Sizeof(s025)")

	var s026 []int = []int{1, 2, 3, 4, 5, 6}
	assert(6, len(s026), "len(s026)")
	assert(6, cap(s026), "cap(s026)")
	assert(8, Sizeof(s026), "Sizeof(s026)")
	assert(1, s026[0], "s026[0]")
	assert(2, s026[1], "s026[1]")
	assert(3, s026[2], "s026[2]")
	assert(4, s026[3], "s026[3]")
	assert(5, s026[4], "s026[4]")
	assert(6, s026[5], "s026[5]")
	s026[0], s026[1], s026[2], s026[3], s026[4], s026[5] = 100, 101, 102, 103, 104, 105
	assert(100, s026[0], "s026[0]")
	assert(101, s026[1], "s026[1]")
	assert(102, s026[2], "s026[2]")
	assert(103, s026[3], "s026[3]")
	assert(104, s026[4], "s026[4]")
	assert(105, s026[5], "s026[5]")

	var s027 = []int{1, 2, 3, 4, 5, 6}
	assert(6, len(s027), "len(s027)")
	assert(6, cap(s027), "cap(s027)")
	assert(8, Sizeof(s027), "Sizeof(s027)")
	assert(1, s027[0], "s027[0]")
	assert(2, s027[1], "s027[1]")
	assert(3, s027[2], "s027[2]")
	assert(4, s027[3], "s027[3]")
	assert(5, s027[4], "s027[4]")
	assert(6, s027[5], "s027[5]")

	s028 := []string{"abc", "def", "ghi"}
	assert(8, Sizeof(s028), "Sizeof(s028)")
	assert(3, len(s028), "len(s028)")
	assert(3, cap(s028), "cap(s028)")
	assert(0, strcmp(s028[0], "abc"), "strcmp(s028[0], \"abc\")")
	assert(0, strcmp(s028[1], "def"), "strcmp(s028[1], \"def\")")
	assert(0, strcmp(s028[2], "ghi"), "strcmp(s028[2], \"ghi\")")

	s029 := [][]string{{"abc", "def", "ghi"}, {"jkl", "mno", "pqr"}}
	assert(8, Sizeof(s029), "Sizeof(s029)")
	assert(2, len(s029), "len(s029)")
	assert(3, len(s029[0]), "len(s029[0])")
	assert(3, len(s029[1]), "len(s029[1])")
	assert(2, cap(s029), "cap(s029)")
	assert(3, cap(s029[0]), "cap(s029[0])")
	assert(3, cap(s029[1]), "cap(s029[1])")
	assert(0, strcmp(s029[0][0], "abc"), "strcmp(s029[0][0], \"abc\")")
	assert(0, strcmp(s029[0][1], "def"), "strcmp(s029[0][1], \"def\")")
	assert(0, strcmp(s029[0][2], "ghi"), "strcmp(s029[0][2], \"ghi\")")
	println(s029[0][2])
	assert(0, strcmp(s029[1][0], "jkl"), "strcmp(s029[1][0], \"jkl\")")
	assert(0, strcmp(s029[1][1], "mno"), "strcmp(s029[1][1], \"mno\")")
	assert(0, strcmp(s029[1][2], "pqr"), "strcmp(s029[1][2], \"pqr\")")

	assert(8, Sizeof(g03), "Sizeof(g03)")
	assert(3, len(g03), "len(g03)")
	assert(3, cap(g03), "cap(g03)")
	assert(0, strcmp(g03[0], "abc"), "strcmp(g03[0], \"abc\")")
	assert(0, strcmp(g03[1], "def"), "strcmp(g03[1], \"def\")")
	assert(0, strcmp(g03[2], "ghi"), "strcmp(g03[2], \"ghi\")")

	assert(8, Sizeof(g04), "Sizeof(g04)")
	assert(2, len(g04), "len(g04)")
	assert(3, len(g04[0]), "len(g04[0])")
	assert(3, len(g04[1]), "len(g04[1])")
	assert(2, cap(g04), "cap(g04)")
	assert(3, cap(g04[0]), "cap(g04[0])")
	assert(3, cap(g04[1]), "cap(g04[1])")
	assert(0, strcmp(g04[0][0], "abc"), "strcmp(g04[0][0], \"abc\")")
	assert(0, strcmp(g04[0][1], "def"), "strcmp(g04[0][1], \"def\")")
	assert(0, strcmp(g04[0][2], "ghi"), "strcmp(g04[0][2], \"ghi\")")
	println(g04[0][2])
	assert(0, strcmp(g04[1][0], "jkl"), "strcmp(g04[1][0], \"jkl\")")
	assert(0, strcmp(g04[1][1], "mno"), "strcmp(g04[1][1], \"mno\")")
	assert(0, strcmp(g04[1][2], "pqr"), "strcmp(g04[1][2], \"pqr\")")

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
	assert(1, s030[0].a[0], "s030[0].a[0]")
	assert(2, s030[0].a[1], "s030[0].a[1]")
	assert(3, s030[0].a[2], "s030[0].a[2]")
	assert(4, s030[0].a[3], "s030[0].a[3]")
	assert(5, s030[0].a[4], "s030[0].a[4]")
	assert(6, s030[0].a[5], "s030[0].a[5]")
	assert(7, s030[0].b, "s030[0].b")
	assert(0, strcmp(s030[0].c[0], "abc"), "strcmp(s030[0].c[0], \"abc\")")
	assert(0, strcmp(s030[0].c[1], "def"), "strcmp(s030[0].c[1], \"def\")")
	assert(0, strcmp(s030[0].c[2], "ghi"), "strcmp(s030[0].c[2], \"ghi\")")
	assert(8, s030[1].a[0], "s030[1].a[0]")
	assert(9, s030[1].a[1], "s030[1].a[1]")
	assert(10, s030[1].a[2], "s030[1].a[2]")
	assert(11, s030[1].a[3], "s030[1].a[3]")
	assert(12, s030[1].a[4], "s030[1].a[4]")
	assert(13, s030[1].a[5], "s030[1].a[5]")
	assert(14, s030[1].b, "s030[1].b")
	assert(0, strcmp(s030[1].c[0], "jkl"), "strcmp(s030[1].c[0], \"abc\")")
	assert(0, strcmp(s030[1].c[1], "mno"), "strcmp(s030[1].c[1], \"def\")")
	assert(0, strcmp(s030[1].c[2], "pqr"), "strcmp(s030[1].c[2], \"ghi\")")

	var s031 = []func() int{ret3, ret3}
	assert(2, len(s031), "len(s031)")
	assert(2, cap(s031), "cap(s031)")
	assert(3, s031[0](), "s031[0]()")
	assert(3, s031[1](), "s031[1]()")

	var s0311 = []func(int) int{reti, reti}
	assert(2, len(s0311), "len(s0311)")
	assert(2, cap(s0311), "cap(s0311)")
	assert(3, s0311[0](3), "s0311[0](3)")
	assert(4, s0311[1](4), "s0311[1](4)")

	var s032 = []float32{1., 2., 3.}
	assert(3, len(s032), "len(s032)")
	assert(3, cap(s032), "cap(s032)")
	assert(1, s032[0] == 1.0, "s032[0]==1.0")
	assert(1, s032[1] == 2.0, "s032[1]==2.0")

	var s033 = []float64{1., 2., 3.}
	assert(3, len(s033), "len(s033)")
	assert(3, cap(s033), "cap(s033)")
	assert(1, s033[0] == 1.0, "s033[0]==1.0")
	assert(1, s033[1] == 2.0, "s033[1]==2.0")

	s034 := make([]int, 10)
	assert(10, len(s034), "len(s034)")
	assert(10, cap(s034), "cap(s034)")
	assert(0, s034[0], "s034[0]")
	assert(0, s034[9], "s034[9]")
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
	assert(101, s034[1], "s034[1]")
	assert(102, s034[2], "s034[2]")
	assert(103, s034[3], "s034[3]")
	assert(104, s034[4], "s034[4]")
	assert(105, s034[5], "s034[5]")
	assert(106, s034[6], "s034[6]")
	assert(107, s034[7], "s034[7]")
	assert(108, s034[8], "s034[8]")
	assert(109, s034[9], "s034[9]")

	s035 := make([]string, 10, 15)
	assert(10, len(s035), "len(s035)")
	assert(15, cap(s035), "cap(s035)")
	assert(0, strcmp(s035[0], ""), "strcmp(s035[0], \"\")")
	assert(0, strcmp(s035[9], ""), "strcmp(s034[9], \"\")")
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
	assert(0, strcmp(s035[1], "abc"), "strcmp(s035[1], \"abc\")")
	assert(0, strcmp(s035[2], "def"), "strcmp(s035[2], \"def\")")
	assert(0, strcmp(s035[3], "ghi"), "strcmp(s035[3], \"ghi\")")
	assert(0, strcmp(s035[4], "jkl"), "strcmp(s035[4], \"jkl\")")
	assert(0, strcmp(s035[5], "mno"), "strcmp(s035[5], \"mno\")")
	assert(0, strcmp(s035[6], "pqr"), "strcmp(s035[6], \"pqr\")")
	assert(0, strcmp(s035[7], "stu"), "strcmp(s035[7], \"stu\")")
	assert(0, strcmp(s035[8], "vwx"), "strcmp(s035[8], \"vwx\")")
	assert(0, strcmp(s035[9], "yz"), "strcmp(s035[9], \"yz\")")

	s036 := make([]int, 1, 5)
	assert(1, len(s036), "len(s036)")
	assert(5, cap(s036), "cap(s036)")
	assert(0, s036[0], "s036[0]")
	s036 = append(s036, 1, 2, 3)
	assert(4, len(s036), "len(s036)")
	assert(5, cap(s036), "cap(s036)")
	assert(0, s036[0], "s036[0]")
	assert(1, s036[1], "s036[1]")
	assert(2, s036[2], "s036[2]")
	assert(3, s036[3], "s036[3]")
	assert(1, s036[reti(1)], "s036[reti(1)]")

	s037 := make([]string, 1, 5)
	assert(1, len(s037), "len(s037)")
	assert(5, cap(s037), "cap(s037)")
	assert(0, strcmp(s037[0], ""), "strcmp(s037[0], \"\")")
	s037 = append(s037, "abc", "def", "ghi")
	assert(4, len(s037), "len(s037)")
	assert(5, cap(s037), "cap(s037)")
	assert(0, strcmp(s037[0], ""), "strcmp(s037[0], \"\")")
	assert(0, strcmp(s037[1], "abc"), "strcmp(s037[1], \"abc\")")
	assert(0, strcmp(s037[2], "def"), "strcmp(s037[2], \"def\")")
	assert(0, strcmp(s037[3], "ghi"), "strcmp(s037[3], \"ghi\")")

	// case: made a slice from global array, int
	a0381 := [5]int{10, 20, 30, 40, 50}
	s038 := a038[0:2]
	s0381 := a0381[0:2]
	assert(2, len(s038), "len(s038)")
	assert(5, cap(s038), "cap(s038)")
	s038[0] = s0381[0]
	s038[1] = a0381[4]
	assert(10, s038[0], "s038[0]")
	assert(50, s038[1], "s038[1]")
	s038 = append(s038, 2, 3, 4, 5)
	assert(6, len(s038), "len(s038)")
	assert(14, cap(s038), "cap(s038)")
	assert(10, s038[0], "s038[0]")
	assert(50, s038[1], "s038[1]")
	assert(2, s038[2], "s038[2]")
	assert(3, s038[3], "s038[3]")
	assert(4, s038[4], "s038[4]")
	assert(5, s038[5], "s038[5]")
	assert(10, a038[0], "a038[0]")
	assert(50, a038[1], "a038[1]")
	assert(0, a038[2], "a038[2]")
	assert(0, a038[3], "a038[3]")
	assert(0, a038[4], "a038[4]")

	// case: made a slice from global array, string
	s0382 := g03[0:2]
	assert(2, len(s0382), "len(s0382)")
	assert(3, cap(s0382), "cap(s0382)")
	s0382[0] = "aaa"
	s0382[1] = "bbb"
	assert(0, strcmp(s0382[0], "aaa"), "strcmp(s0382[0], \"aaa\")")
	assert(0, strcmp(s0382[1], "bbb"), "strcmp(s0382[1], \"bbb\")")
	s0382 = append(s0382, "jkl", "mno", "pqr", "stu")
	assert(6, len(s0382), "len(s0382)")
	assert(10, cap(s0382), "cap(s0382)")
	assert(0, strcmp(s0382[0], "aaa"), "strcmp(s0382[0], \"aaa\")")
	assert(0, strcmp(s0382[1], "bbb"), "strcmp(s0382[1], \"bbb\")")
	assert(0, strcmp(s0382[2], "jkl"), "strcmp(s0382[2], \"jkl\")")
	assert(0, strcmp(s0382[3], "mno"), "strcmp(s0382[3], \"mno\")")
	assert(0, strcmp(s0382[4], "pqr"), "strcmp(s0382[4], \"pqr\")")
	assert(0, strcmp(s0382[5], "stu"), "strcmp(s0382[5], \"stu\")")

	// case: made a slice by initializer, string
	s03821 := []string{"abc", "def", "ghi"}
	assert(3, len(s03821), "len(s03821)")
	assert(3, cap(s03821), "cap(s03821)")
	s03821[0] = "aaa"
	s03821[1] = "bbb"
	assert(0, strcmp(s03821[0], "aaa"), "strcmp(s03821[0], \"aaa\")")
	assert(0, strcmp(s03821[1], "bbb"), "strcmp(s03821[1], \"bbb\")")
	s03821 = append(s03821, "jkl", "mno", "pqr", "stu")
	assert(7, len(s03821), "len(s03821)")
	assert(10, cap(s03821), "cap(s03821)")
	assert(0, strcmp(s03821[0], "aaa"), "strcmp(s03821[0], \"aaa\")")
	assert(0, strcmp(s03821[1], "bbb"), "strcmp(s03821[1], \"bbb\")")
	assert(0, strcmp(s03821[2], "ghi"), "strcmp(s03821[2], \"ghi\")")
	assert(0, strcmp(s03821[3], "jkl"), "strcmp(s03821[3], \"jkl\")")
	assert(0, strcmp(s03821[4], "mno"), "strcmp(s03821[4], \"mno\")")
	assert(0, strcmp(s03821[5], "pqr"), "strcmp(s03821[5], \"pqr\")")
	assert(0, strcmp(s03821[6], "stu"), "strcmp(s03821[6], \"stu\")")

	// case: made a global slice , string, two dimensions
	assert(2, len(g04), "len(g04)")
	assert(2, cap(g04), "cap(g04)")
	assert(3, len(g04[0]), "len(g04[0])")
	assert(3, cap(g04[0]), "cap(g04[0])")
	assert(3, len(g04[1]), "len(g04[1])")
	assert(3, cap(g04[1]), "cap(g04[1])")
	g04[0][2] = "ggg"
	s038211 := []string{"aaa", "bbb"}
	g04 = append(g04, s038211)
	g04[1] = append(g04[1], "stu")
	assert(3, len(g04), "len(g04)")
	assert(5, cap(g04), "cap(g04)")
	assert(4, len(g04[0]), "len(g04[0])")
	assert(7, cap(g04[0]), "cap(g04[0])")
	assert(4, len(g04[1]), "len(g04[1])")
	assert(7, cap(g04[1]), "cap(g04[1])")
	assert(4, len(g04[2]), "len(g04[1])")
	assert(0, strcmp(g04[0][0], "abc"), "strcmp(g04[0][0], \"abc\")")
	assert(0, strcmp(g04[0][1], "def"), "strcmp(g04[0][1], \"def\")")
	assert(0, strcmp(g04[0][2], "ggg"), "strcmp(g04[0][2], \"ggg\")")
	assert(0, strcmp(g04[1][0], "jkl"), "strcmp(g04[1][0], \"jkl\")")
	assert(0, strcmp(g04[1][1], "mno"), "strcmp(g04[1][1], \"mno\")")
	assert(0, strcmp(g04[1][2], "pqr"), "strcmp(g04[1][2], \"pqr\")")
	assert(0, strcmp(g04[1][3], "stu"), "strcmp(g04[1][3], \"stu\")")
	assert(0, strcmp(g04[2][0], "aaa"), "strcmp(g04[2][0], \"aaa\")")
	assert(0, strcmp(g04[2][1], "bbb"), "strcmp(g04[2][1], \"bbb\")")

	// case: made a slice by initializer, int
	s03822 := []int{1, 2, 3}
	assert(3, len(s03822), "len(s03822)")
	assert(3, cap(s03822), "cap(s03822)")
	s03822[0] = 11
	s03822[1] = 22
	assert(11, s03822[0], "s03822[0])")
	assert(22, s03822[1], "s03822[1])")
	s03822 = append(s03822, 4, 5, 6, 7)
	assert(7, len(s03822), "len(s03832)")
	assert(10, cap(s03822), "cap(s03832)")
	assert(11, s03822[0], "s03822[0]")
	assert(22, s03822[1], "s03822[1]")
	assert(3, s03822[2], "s03822[2]")
	assert(4, s03822[3], "s03822[3]")
	assert(5, s03822[4], "s03822[4]")
	assert(6, s03822[5], "s03822[5]")
	assert(7, s03822[6], "s03822[6]")

	// case: made a slice by 'var', int
	var s03823 []int
	assert(0, len(s03823), "len(s03823)")
	assert(0, cap(s03823), "cap(s03823)")
	s03823 = append(s03823, 1, 2, 3, 4, 5, 6)
	assert(6, len(s03823), "len(s03823)")
	assert(6, cap(s03823), "cap(s03823)")
	assert(1, s03823[0], "s03823[0]")
	assert(2, s03823[1], "s03823[1]")
	assert(3, s03823[2], "s03823[2]")
	assert(4, s03823[3], "s03823[3]")
	assert(5, s03823[4], "s03823[4]")
	assert(6, s03823[5], "s03823[5]")
	s03823[2] = 1000
	assert(1000, s03823[2], "s03823[2]")
	s03823 = append(s03823, 7, 8, 9, 10, 11, 12)
	assert(12, len(s03823), "len(s03823)")
	assert(18, cap(s03823), "cap(s03823)")
	assert(1, s03823[0], "s03823[0]")
	assert(2, s03823[1], "s03823[1]")
	assert(1000, s03823[2], "s03823[2]")
	assert(4, s03823[3], "s03823[3]")
	assert(5, s03823[4], "s03823[4]")
	assert(6, s03823[5], "s03823[5]")
	assert(7, s03823[6], "s03823[6]")
	assert(8, s03823[7], "s03823[7]")
	assert(9, s03823[8], "s03823[8]")
	assert(10, s03823[9], "s03823[9]")
	assert(11, s03823[10], "s03823[10]")
	assert(12, s03823[11], "s03823[11]")

	// case: made a slice by make function, int
	s0383 := make([]int, 5)
	assert(5, len(s0383), "len(s0383)")
	assert(5, cap(s0383), "cap(s0383)")
	s0383[0] = 10
	s0383[1] = 50
	assert(10, s0383[0], "s0383[0]")
	assert(50, s0383[1], "s0383[1]")
	s0383 = append(s0383, 2, 3, 4, 5)
	assert(9, len(s0383), "len(s0383)")
	assert(14, cap(s0383), "cap(s0383)")
	assert(10, s0383[0], "s0383[0]")
	assert(50, s0383[1], "s0383[1]")
	assert(0, s0383[2], "s0383[2]")
	assert(0, s0383[3], "s0383[3]")
	assert(0, s0383[4], "s0383[4]")
	assert(2, s0383[5], "s0383[5]")
	assert(3, s0383[6], "s0383[6]")
	assert(4, s0383[7], "s0383[7]")
	assert(5, s0383[8], "s0383[8]")

	// case: made a slice by make function, string
	s039 := []string{"abc", "def"}
	assert(2, len(s039), "len(s039)")
	assert(2, cap(s039), "cap(s039)")
	assert(0, strcmp(s039[0], "abc"), "strcmp(s039[0], \"abc\")")
	assert(0, strcmp(s039[1], "def"), "strcmp(s039[1], \"def\")")
	s039 = append(s039, "ghi", "jkl", "mno")
	assert(0, strcmp(s039[0], "abc"), "strcmp(s039[0], \"abc\")")
	assert(0, strcmp(s039[1], "def"), "strcmp(s039[1], \"def\")")
	assert(0, strcmp(s039[2], "ghi"), "strcmp(s039[2], \"ghi\")")
	assert(0, strcmp(s039[3], "jkl"), "strcmp(s039[3], \"jkl\")")
	assert(0, strcmp(s039[4], "mno"), "strcmp(s039[4], \"mno\")")

	s040 := make([]string, 2, 3)
	assert(2, len(s040), "len(s040)")
	assert(3, cap(s040), "cap(s040)")
	s040[0] = "abc"
	s040[1] = "def"
	assert(0, strcmp(s040[0], "abc"), "strcmp(s040[0], \"abc\")")
	assert(0, strcmp(s040[1], "def"), "strcmp(s040[1], \"def\")")
	s040 = append(s040, "ghi", "jkl", "mno")
	assert(0, strcmp(s040[0], "abc"), "strcmp(s040[0], \"abc\")")
	assert(0, strcmp(s040[1], "def"), "strcmp(s040[1], \"def\")")
	assert(0, strcmp(s040[2], "ghi"), "strcmp(s040[2], \"ghi\")")
	assert(0, strcmp(s040[3], "jkl"), "strcmp(s040[3], \"jkl\")")
	assert(0, strcmp(s040[4], "mno"), "strcmp(s040[4], \"mno\")")

	var a041 = [8]int{0, 1, 2, 3, 4, 5, 6, 7}
	var s041 = make([]int, 6)
	// var s0411 = a041[0:8]
	n0411 := copy(s041, a041[0:]) // n0411 == 6, s041 == []int{0, 1, 2, 3, 4, 5}
	assert(0, s041[0], "s041[0]")
	assert(1, s041[1], "s041[1]")
	assert(2, s041[2], "s041[2]")
	assert(3, s041[3], "s041[3]")
	assert(4, s041[4], "s041[4]")
	assert(5, s041[5], "s041[5]")
	n0412 := copy(s041, s041[2:]) // n0412 == 4, s041 == []int{2, 3, 4, 5, 4, 5}
	assert(6, n0411, "n0411")
	assert(4, n0412, "n0412")
	assert(2, s041[0], "s041[0]")
	assert(3, s041[1], "s041[1]")
	assert(4, s041[2], "s041[2]")
	assert(5, s041[3], "s041[3]")
	assert(4, s041[4], "s041[4]")
	assert(5, s041[5], "s041[5]")
	a041[0] = 1000
	assert(2, s041[0], "s041[0]")

	var a042 = [6]string{"abc", "def", "ghi", "jkl", "mno", "pqr"}
	var s042 = make([]string, 4)
	n0421 := copy(s042, a042[0:])
	assert(4, n0421, "n0421")
	assert(0, strcmp(s042[0], "abc"), "strcmp(s042[0], \"abc\")")
	assert(0, strcmp(s042[1], "def"), "strcmp(s042[1], \"def\")")
	assert(0, strcmp(s042[2], "ghi"), "strcmp(s042[2], \"ghi\")")
	assert(0, strcmp(s042[3], "jkl"), "strcmp(s042[3], \"jkl\")")
	n0422 := copy(s042, s042[2:])
	assert(2, n0422, "n0422")
	assert(0, strcmp(s042[0], "ghi"), "strcmp(s042[0], \"ghi\")")
	assert(0, strcmp(s042[1], "jkl"), "strcmp(s042[1], \"jkl\")")
	assert(0, strcmp(s042[2], "ghi"), "strcmp(s042[2], \"ghi\")")
	assert(0, strcmp(s042[3], "jkl"), "strcmp(s042[3], \"jkl\")")

	s043 := make([][]int, 6)
	assert(6, len(s043), "len(s043)")
	assert(0, len(s043[0]), "len(s043[0])")
	s043[0] = append(s043[0], 1)
	assert(1, len(s043[0]), "len(s043[0])")
	assert(1, s043[0][0], "s043[0][0]")
	assert(1, len(s043[1]), "len(s043[1])") // ????????????????????????

	println("OK")
}
