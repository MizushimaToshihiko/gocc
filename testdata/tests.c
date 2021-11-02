// -*- c -*-

// line comment

/*
* block comment
*/

int assert(int expected, int actual, char *code) {
  if (expected == actual) {
    printf("%s => %d\n", code, actual);
  } else {
    printf("%s => %d expected but got %d\n", code, expected, actual);
    exit(1);
  }
}

int main() {
  assert(8, ({ int a=3; int z=5; a+z; }), "int a=3; int z=5; a+z");

  assert(0, 0, "0");
  assert(42, 42, "42");
  assert(5, 5, "0");
  assert(41,  12 + 34 - 5 , " 12 + 34 - 5 ");
  assert(5, 5, "0");
  assert(15, 5*(9-6), "5*(9-6)");
  assert(4, (3+5)/2, "(3+5)/2");
  assert(-10, -10, "0");
  assert(10, - -10, "- -10");
  assert(10, - - +10, "- - +10");

  // assert(2, ({ int x=2; { int x=3; } x; }), "int x=2; { int x=3; } x;");
  // assert(2, ({ int x=2; { int x=3; } int y=4; x; }), "int x=2; { int x=3; } int y=4; x;");
  // assert(3, ({ int x=2; { x=3; } x; }), "int x=2; { x=3; } x;");

  printf("OK\n");
  return 0;
}