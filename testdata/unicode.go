package test_unicode

func assert(want int, act int, code string)
func println(format string)
func strcmp(s1, s2 string)

#include "test.h"

func main() {
	ASSERT(0, strcmp("αβγ", "\u03B1\u03B2\u03B3"))
  ASSERT(0, strcmp("日本語", "\u65E5\u672C\u8A9E"))
  ASSERT(0, strcmp("日本語", "\U000065E5\U0000672C\U00008A9E"))
  ASSERT(0, strcmp("🌮", "\U0001F32E"))
}