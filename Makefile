CFLAGS=-std=c11 -g -fno-common

BINARY_NAME := bin/gocc
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt

SRCS=$(wildcard *.go)

TEST_SRCS=$(wildcard testdata/*.go)
TESTS=$(TEST_SRCS:.go=.exe)

all: build test

build: $(SRCS)
	$(GOBUILD) -o $(BINARY_NAME) -v $^

testdata/%.exe: testdata/%.go
	$(BINARY_NAME) -o testdata/$*.s $^
	$(CC) -static -o $@ testdata/$*.s -xc testdata/common

test: $(TESTS)
	for i in $^; do echo $$i; ./$$i || exit 1; echo; done

clean: 
	$(GOCLEAN)
	rm -f bin/* testdata/*.o testdata/*.s testdata/asm* profile

fmt:
	$(GOFMT) ./...

coverage: $(SRCS)
	$(GOTEST) $^ -coverprofile=profile
	$(GOCMD) tool cover -html=profile

.PHONY: build test clean fmt coverage