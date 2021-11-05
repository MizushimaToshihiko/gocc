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
	$(GOTEST) $^ -cover -v -count 1 -timeout 120s # -v

clean: 
	$(GOCLEAN)
	rm -f testdata/tmp testdata/asm.s profile

fmt:
	$(GOFMT) ./...

coverage: $(SRCS)
	$(GOTEST) $^ -coverprofile=profile
	$(GOCMD) tool cover -html=profile

.PHONY: test clean fmt coverage