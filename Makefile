SHELL := bash

default: build

build:
	/bin/bash -c "export GOARCH=$(GOARCH) && export GO_BUILD_FLAGS="$(GO_BUILD_FLAGS)" && hack/build.sh"

test:
	go test -v ./...

lint:
	golangci-lint run -v

coverprofile:
	hack/coverprofile.sh

dependencies:
	hack/update-deps.sh

listen: build
	socat UNIX-LISTEN:ssh.sock,fork EXEC:bin/wsl2-ssh-pageant.exe

.PHONY: \
	default \
	build \
	test \
	lint \
	coverprofile \
	dependencies \
	listen \
	$(NULL)
