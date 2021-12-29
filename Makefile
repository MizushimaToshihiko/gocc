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

all: $(BINARY_NAME) test

$(BINARY_NAME): $(SRCS)
	$(GOBUILD) -o $(BINARY_NAME) -v $^

testdata/%.exe: testdata/%.go
	$(BINARY_NAME) -o testdata/$*.s $^
	gcc -static -g -o $@ testdata/$*.s

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

.PHONY: test clean fmt coverage