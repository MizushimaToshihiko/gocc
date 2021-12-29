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

$(BINARY_NAME): $(SRCS)
	$(GOBUILD) -o $(BINARY_NAME) -v $^

test: $(SRCS)
	$(GOTEST) $^ -cover -count 1 -timeout 600s

clean: 
	$(GOCLEAN)
	rm -f bin/* testdata/*.o testdata/*.s testdata/asm* profile

fmt:
	$(GOFMT) ./...

coverage: $(SRCS)
	$(GOTEST) $^ -coverprofile=profile
	$(GOCMD) tool cover -html=profile

.PHONY: test clean fmt coverage