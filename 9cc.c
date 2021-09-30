#include <ctype.h>
#include <stdarg.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// types of token
typedef enum {
  TK_RESERVED,  // symbol
  TK_NUM,       // integer
  TK_EOF,       // end of the input
} TokeKind;

typedef struct Token Token;

struct Token {
  TokenKind kind;
  Token *next;
  int val;
  char *str;
};

// the global variableof the current token
Token *token;

void error (char *fmt, ...) {
  va_list ap;
  va_start(ap, fmt);
  vfprintf(stderr, "\n");
  exit(1);
}

int main(int argc, char **argv) {
  if (argc != 2) {
    fprintf(stderr, "invalid the number of arguments\n");
    return 1;
  }

  char *p = argv[1];

  printf(".intel_syntax noprefix\n");
  printf(".global main\n");
  printf("main:\n");
  printf("  mov rax, %ld\n", strtol(p, &p, 10));

  while (*p) {
    if (*p == '+') {
      p++;
      printf("  add rax, %ld\n", strtol(p, &p, 10));
      continue;
    }

    if (*p == '-') {
      p++;
      printf("  sub rax, %ld\n", strtol(p, &p, 10));
      continue;
    }

    fprintf(stderr, "unexpected character: '%c'\n", *p);
    return 1;
  }
  printf("  ret\n");
  return 0;
}
