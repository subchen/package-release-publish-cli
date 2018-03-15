CWD     := $(shell pwd)
NAME    := sha256sum-files
VERSION := 0.1.0

PACKAGES := $(shell go list ./... | grep -v /vendor/)

default:
	@ echo "no default target for Makefile"

clean:
	@ rm -rf $(NAME) ./_releases

glide-vc:
	@ glide-vc --only-code --no-tests --no-legal-files

fmt:
	@ go fmt $(PACKAGES)
	@ goimports -w *.go

build: clean
	@ go run ../go-build/main.go -n $(NAME) -v $(VERSION)

release: build
	@ go run main.go ./_releases/
