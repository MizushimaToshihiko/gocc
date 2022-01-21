package test

func main() {
	assert(0, 0, "0")
	assert(42, 42, "42")
	assert(5, 5, "5")
	assert(41, 12+34-5, "12 + 34 - 5")
	assert(15, 5*(9-6), "5*(9-6)")
	assert(4, (3+5)/2, "(3+5)/2")
	assert(10, -10+20, "-10+20")
	assert(10, - -10, "- -10")
	assert(10, - -+10, "- - +10")

	assert(0, 0 == 1, "0==1")
	assert(1, 42 == 42, "42==42")
	assert(1, 0 != 1, "0!=1")
	assert(0, 42 != 42, "42!=42")

	assert(1, 0 < 1, "0<1")
	assert(0, 1 < 1, "1<1")
	assert(0, 2 < 1, "2<1")
	assert(1, 0 <= 1, "0<=1")
	assert(1, 1 <= 1, "1<=1")
	assert(0, 2 <= 1, "2<=1")

	assert(1, 1 > 0, "1>0")
	assert(0, 1 > 1, "1>1")
	assert(0, 1 > 2, "1>2")
	assert(1, 1 >= 0, "1>=0")
	assert(1, 1 >= 1, "1>=1")
	assert(0, 1 >= 2, "1>=2")

	assert(4294967297, 4294967297, "4294967297")
	assert(0, 1073741824*100/100, "1073741824 * 100 / 100")

	var i int
	i = 2, i += 5
	assert(7, i, "i=2, i+=5, i")
	i = 5, i -= 2
	assert(3, i, "i=5, i-=2, i")
	i=3, i*=2
	assert(6, i, "i=3, i*=2, i")
	i=6; i/=2
	assert(3,i,"i=6; i/=2, i")
	i=2; i++
	assert(3, i,"i=2; i++, i")
	i=2; i--
	assert(1, i, "i=2; i--, i")

	assert(0, !1, "!1")
	assert(0, !2, "!2")
	assert(1, !0, "!0")
	assert(1, !byte(0), "!byte(0)")
	assert(3, int64(3), "int64(3)")
	// ASSERT(4, sizeof(!(char)0));
	// ASSERT(4, sizeof(!(long)0));

	assert(-1, ^0, "^0")
	assert(0, ^-1, "^-1")

	assert(5, 17%6, "17%6")
	// assert(5, (int64(17))%6, "(int64(17))%6")
	i = 10
	i %= 4
	assert(2, i, "i=10; i%=4")
	var i int64
	i = 10
	i %= 4
	assert(2, i, "i=10; i%=4")

	assert(0, 0&1, "0&1")
	assert(1, 3&1, "3&1")
	assert(3, 7&3, "7&3")
	assert(10, -1&10, "-1&10")

	assert(1, 0|1, "0|1")
	// assert(0b10011, 0b10000|0b00011, "0b10000|0b00011")

	assert(0, 0^0, "0^0")
	// assert(0, 0b1111^0b1111, "0b1111^0b1111")
	// assert(0b110100, 0b111000^0b001100, "0b111000^0b001100")

	printf("\nOK\n")
}
