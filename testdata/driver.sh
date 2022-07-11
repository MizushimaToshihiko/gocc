#!/bin/bash
gocc=$1

tmp=`mktemp -d /tmp/gocc-test-XXXXXX`
trap 'rm -rf $tmp' INT TERM HUP EXIT
echo > $tmp/empty.c

check() {
  if [ $? -eq 0 ]; then
    echo "testing $1 ... passed"
  else
    echo "testing $1 ... failed"
    exit 1
  fi
}

rm -f $tmp/foo.go $tmp/bar.go
echo 'var x int' > $tmp/foo.go
echo 'var y int' > $tmp/bar.go
(cd $tmp; $OLDPWD/$gocc $tmp/foo.go $tmp/bar.go)
[ -f $tmp/foo.s ] && [ -f $tmp/bar.s ]
check 'multiple input files'

echo OK