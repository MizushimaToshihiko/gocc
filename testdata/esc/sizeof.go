package test

func main() {
	assert(1, Sizeof(byte), "Sizeof(byte)")
	// assert(2, sizeof(short));
	// assert(2, sizeof(short int));
	// assert(2, sizeof(int short));
	// assert(4, sizeof(int));
	// assert(8, sizeof(long));
	// assert(8, sizeof(long int));
	// assert(8, sizeof(long int));
	// assert(8, sizeof(char *));
	// assert(8, sizeof(int *));
	// assert(8, sizeof(long *));
	// assert(8, sizeof(int **));
	// assert(8, sizeof(int(*)[4]));
	// assert(32, sizeof(int*[4]));
	// assert(16, sizeof(int[4]));
	// assert(48, sizeof(int[3][4]));
	// assert(8, sizeof(struct {int a; int b;}));

	// assert(8, sizeof(-10 + (long)5));
	// assert(8, sizeof(-10 - (long)5));
	// assert(8, sizeof(-10 * (long)5));
	// assert(8, sizeof(-10 / (long)5));
	// assert(8, sizeof((long)-10 + 5));
	// assert(8, sizeof((long)-10 - 5));
	// assert(8, sizeof((long)-10 * 5));
	// assert(8, sizeof((long)-10 / 5));

	// assert(1, ({ char i; sizeof(++i); }));
	// assert(1, ({ char i; sizeof(i++); }));

	// assert(8, sizeof(int(*)[10]));
	// assert(8, sizeof(int(*)[][10]));

	printf("OK\n")
}
