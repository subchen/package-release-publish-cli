GOPATH := $(shell pwd)/.build/go

export GOPATH

build:
	cd $(GOPATH)/src/github.com/subchen/storm && go build -o storm
