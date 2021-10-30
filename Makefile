BINARY_NAME := bin/compiler
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt

SRCS=$(wildcard *.go)

all: build test

build: $(SRCS)
	$(GOBUILD) -o $(BINARY_NAME) -v $^

test: $(SRCS)
	$(GOTEST) $^ -cover -v -count 1 -timeout 60s # -v

clean: 
	$(GOCLEAN)
	rm -f testdata/temp* testdata/in*.c testdata/asm.s testdata/funcs_file profile

fmt:
	$(GOFMT) ./...

coverage: $(SRCS)
	$(GOTEST) $^ -coverprofile=profile
	$(GOCMD) tool cover -html=profile

.PHONY: test clean fmt coverage