CWD    := $(shell pwd)
NAME    := pts
VERSION := 0.1.0

LDFLAGS := -s -w \
           -X 'main.BuildVersion=$(VERSION)' \
           -X 'main.BuildGitRev=$(shell git rev-list HEAD --count)' \
           -X 'main.BuildGitCommit=$(shell git describe --abbrev=0 --always)' \
           -X 'main.BuildDate=$(shell date -u -R)'

PACKAGES := $(shell go list ./... | grep -v /vendor/)

default:
	@ echo "no default target for Makefile"

clean:
	@ rm -rf $(NAME) ./releases ./build

glide-vc:
	@ glide-vc --only-code --no-tests --no-legal-files

fmt:
	@ go fmt $(PACKAGES)
	@ goimports -w .

build: \
    build-linux \
    build-darwin \
    build-windows

build-linux: clean fmt
	@ GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o releases/$(NAME)-$(VERSION)-linux-amd64

build-darwin: clean fmt
	@ GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o releases/$(NAME)-$(VERSION)-darwin-amd64

build-windows: clean fmt
	@ GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o releases/$(NAME)-$(VERSION)-windows-amd64.exe

sha256sum: build
	@ for f in $(shell ls ./releases); do \
		cd $(CWD)/releases; sha256sum "$$f" >> $$f.sha256; \
	done

release: sha256sum
