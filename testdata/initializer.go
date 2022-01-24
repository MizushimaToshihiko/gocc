package test

var g3 byte = 3
var g4 int16 = 4
var g5 int = 5
var g6 int64 = 6

func main() {
	var x [3]int={1,2,3}
  assert(1, x[0], "var x [3]int={1,2,3}; x[0]");
  // assert(2, ({ int x[3]={1,2,3}; x[1]; }));
  // assert(3, ({ int x[3]={1,2,3}; x[2]; }));
  // assert(3, ({ int x[3]={1,2,3}; x[2]; }));

  // assert(2, ({ int x[2][3]={{1,2,3},{4,5,6}}; x[0][1]; }));
  // assert(4, ({ int x[2][3]={{1,2,3},{4,5,6}}; x[1][0]; }));
  // assert(6, ({ int x[2][3]={{1,2,3},{4,5,6}}; x[1][2]; }));

  // assert(0, ({ int x[3]={}; x[0]; }));
  // assert(0, ({ int x[3]={}; x[1]; }));
  // assert(0, ({ int x[3]={}; x[2]; }));

  // assert(2, ({ int x[2][3]={{1,2}}; x[0][1]; }));
  // assert(0, ({ int x[2][3]={{1,2}}; x[1][0]; }));
  // assert(0, ({ int x[2][3]={{1,2}}; x[1][2]; }));

  // assert('a', ({ char x[4]="abc"; x[0]; }));
  // assert('c', ({ char x[4]="abc"; x[2]; }));
  // assert(0, ({ char x[4]="abc"; x[3]; }));
  // assert('a', ({ char x[2][4]={"abc","def"}; x[0][0]; }));
  // assert(0, ({ char x[2][4]={"abc","def"}; x[0][3]; }));
  // assert('d', ({ char x[2][4]={"abc","def"}; x[1][0]; }));
  // assert('f', ({ char x[2][4]={"abc","def"}; x[1][2]; }));

  // assert(4, ({ int x[]={1,2,3,4}; x[3]; }));
  // assert(16, ({ int x[]={1,2,3,4}; sizeof(x); }));
  // assert(4, ({ char x[]="foo"; sizeof(x); }));

  // assert(4, ({ typedef char T[]; T x="foo"; T y="x"; sizeof(x); }));
  // assert(2, ({ typedef char T[]; T x="foo"; T y="x"; sizeof(y); }));
  // assert(2, ({ typedef char T[]; T x="x"; T y="foo"; sizeof(x); }));
  // assert(4, ({ typedef char T[]; T x="x"; T y="foo"; sizeof(y); }));

  // assert(1, ({ struct {int a; int b; int c;} x={1,2,3}; x.a; }));
  // assert(2, ({ struct {int a; int b; int c;} x={1,2,3}; x.b; }));
  // assert(3, ({ struct {int a; int b; int c;} x={1,2,3}; x.c; }));
  // assert(1, ({ struct {int a; int b; int c;} x={1}; x.a; }));
  // assert(0, ({ struct {int a; int b; int c;} x={1}; x.b; }));
  // assert(0, ({ struct {int a; int b; int c;} x={1}; x.c; }));

  // assert(1, ({ struct {int a; int b;} x[2]={{1,2},{3,4}}; x[0].a; }));
  // assert(2, ({ struct {int a; int b;} x[2]={{1,2},{3,4}}; x[0].b; }));
  // assert(3, ({ struct {int a; int b;} x[2]={{1,2},{3,4}}; x[1].a; }));
  // assert(4, ({ struct {int a; int b;} x[2]={{1,2},{3,4}}; x[1].b; }));

  // assert(0, ({ struct {int a; int b;} x[2]={{1,2}}; x[1].b; }));

  // assert(0, ({ struct {int a; int b;} x={}; x.a; }));
  // assert(0, ({ struct {int a; int b;} x={}; x.b; }));

  // assert(5, ({ typedef struct {int a,b,c,d,e,f;} T; T x={1,2,3,4,5,6}; T y; y=x; y.e; }));
  // assert(2, ({ typedef struct {int a,b;} T; T x={1,2}; T y, z; z=y=x; z.b; }));

  // assert(1, ({ typedef struct {int a,b;} T; T x={1,2}; T y=x; y.a; }));

  // assert(4, ({ union { int a; char b[4]; } x={0x01020304}; x.b[0]; }));
  // assert(3, ({ union { int a; char b[4]; } x={0x01020304}; x.b[1]; }));

  // assert(0x01020304, ({ union { struct { char a,b,c,d; } e; int f; } x={{4,3,2,1}}; x.f; }));

  // assert(3, g3);
  // assert(4, g4);
  // assert(5, g5);
  // assert(6, g6);

  println("\nOK");
}
