#!/bin/bash
gocc=$1

tmp=`mktemp -d /tmp/gocc-test-XXXXXX`
trap 'rm -rf $tmp' INT TERM HUP EXIT
echo > $tmp/empty.go

check() {
  if [ $? -eq 0 ]; then
    echo "testing $1 ... passed"
  else
    echo "testing $1 ... failed"
    exit 1
  fi
}

rm -f $tmp/foo.o $tmp/bar.o
echo 'var x int' > $tmp/foo.go
echo 'var y int' > $tmp/bar.go
(cd $tmp; $OLDPWD/$gocc -c $tmp/foo.go $tmp/bar.go)
[ -f $tmp/foo.o ] && [ -f $tmp/bar.o ]
check 'multiple input files'

rm -f $tmp/foo.s $tmp/bar.s
echo 'var x int' > $tmp/foo.go
echo 'var y int' > $tmp/bar.go
(cd $tmp; $OLDPWD/$gocc -S $tmp/foo.go $tmp/bar.go)
[ -f $tmp/foo.s ] && [ -f $tmp/bar.s ]
check 'multiple input files'

# Run linker
rm -f $tmp/foo $tmp/foo.go
echo 'func main() {}' > $tmp/foo.go
$gocc -o $tmp/foo $tmp/foo.go
$tmp/foo
check linker

echo OK