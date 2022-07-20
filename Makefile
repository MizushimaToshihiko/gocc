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

all: clean build test

build: $(SRCS)
	$(GOBUILD) -o $(BINARY_NAME) -v $^

testdata/%.exe: testdata/%.go
	$(BINARY_NAME) -c -o $(^D)/$*.o $^
	$(CC) -static -g -o $@ $(^D)/$*.o -xc $(^D)/common

test: $(TESTS)
	for i in $^; do echo $$i; ./$$i || exit 1; echo; done
	./$(<D)/driver.sh $(BINARY_NAME)


clean: 
	$(GOCLEAN)
	rm -f bin/* testdata/*.s testdata/*.exe testdata/*.o profile

fmt:
	$(GOFMT) ./...

coverage: $(SRCS)
	$(GOTEST) $^ -coverprofile=profile
	$(GOCMD) tool cover -html=profile

.PHONY: build test clean fmt coverage