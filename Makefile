BINARY_NAME := bin/compiler
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt

SRCS=$(wildcard *.go)

all: build test

build: main.go
	$(GOBUILD) -o $(BINARY_NAME) -v $^

test: $(SRCS)
	$(GOTEST) $^ -v -cover

clean: 
	$(GOCLEAN)
	rm -f temporary* profile

fmt:
	$(GOFMT) ./...

coverage: $(SRCS)
	$(GOTEST) $^ -coverprofile=profile
	$(GOCMD) tool cover -html=profile

.PHONY: test clean fmt coverage