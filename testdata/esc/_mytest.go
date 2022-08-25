package test_unicode

func assert(want int, act int, code string)
func println(format string)
func strcmp(s1, s2 string)

#include "test.h"

func main() {
	ASSERT(0, strcmp("αβγ", "αβγ"))
	ASSERT(0, strcmp("日本語", "日本語"))
	ASSERT(0, strcmp("日本語", "日本語"))
	ASSERT(0, strcmp("🌮", "🌮"))

	ASSERT(0, strcmp("ÿ", "ÿ"))
	ASSERT(0, strcmp("ÿ", "ÿ"))
	ASSERT(0, strcmp("あ", "あ"))

	// ASSERT(946, 'β')
	// ASSERT(12354, 'あ')
	// ASSERT(127843, '🍣')
	println("OK")
}