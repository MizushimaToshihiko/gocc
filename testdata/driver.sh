#!/bin/bash
gocc=$1

tmp=`mktemp -d /tmp/gocc-test-XXXXXX`
trap 'rm -rf $tmp' INT TERM HUP EXIT
echo "package empty" > $tmp/empty.go

check() {
  if [ $? -eq 0 ]; then
    echo "testing $1 ... passed"
  else
    echo "testing $1 ... failed"
    exit 1
  fi
}

# -o
rm -f $tmp/out
$gocc -c -o $tmp/out $tmp/empty.go
[ -f $tmp/out ]
check -o

# -h
$gocc -h 2>&1 | grep -q gocc
check -h

# -S
echo 'func main() {}' | $gocc -S -o - - | grep -q 'main:'
check -S

# Default output file
rm -f $tmp/out.o $tmp/out.s
echo 'func main() {}' > $tmp/out.go
(cd $tmp; $OLDPWD/$gocc -c out.go)
[ -f $tmp/out.o ]
check 'default output file'

(cd $tmp; $OLDPWD/$gocc -c -S out.go)
[ -f $tmp/out.s ]
check 'default output file'

# Multiple input files
rm -f $tmp/foo.o $tmp/bar.o
echo 'var x int' > $tmp/foo.go
echo 'var y int' > $tmp/bar.go
(cd $tmp; $OLDPWD/$gocc -c $tmp/foo.go $tmp/bar.go)
[ -f $tmp/foo.o ] && [ -f $tmp/bar.o ]
check 'multiple input files'

rm -f $tmp/foo.s $tmp/bar.s
echo 'var x int' > $tmp/foo.go
echo 'var y int' > $tmp/bar.go
(cd $tmp; $OLDPWD/$gocc -c -S $tmp/foo.go $tmp/bar.go)
[ -f $tmp/foo.s ] && [ -f $tmp/bar.s ]
check 'multiple input files'

# Run linker
rm -f $tmp/foo $tmp/foo.go
echo 'func main() {}' > $tmp/foo.go
$gocc -o $tmp/foo $tmp/foo.go
$tmp/foo
check linker

rm -f $tmp/foo
echo 'func Bar() int; func main() { return Bar(); }' > $tmp/foo.go
echo 'func Bar() int { return 42; }' > $tmp/bar.go
$gocc -o $tmp/foo $tmp/foo.go $tmp/bar.go
$tmp/foo
[ "$?" = 42 ]
check linker

echo OK