SHELL := bash

default: build

# Build the binary
build:
	/bin/bash -c "export GOARCH=$(GOARCH) && export GO_BUILD_FLAGS="$(GO_BUILD_FLAGS)" && hack/build.sh"

# Run unit tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run -v

# Format code
fmt:
	gofmt -s -w ./cmd ./pkg

# Validate that all generated files are up to date
validate:
	hack/validate.sh

# Generate cover profile
coverprofile:
	hack/coverprofile.sh

# Update dependencies
dependencies:
	hack/update-deps.sh

# Build and start the listener for testing
listen: build
	socat UNIX-LISTEN:ssh.sock,fork EXEC:bin/wsl2-ssh-pageant.exe

# Show this help message
help:
	@echo "Available targets:"
	@echo ""
	@awk '/^#/{c=substr($$0,3);next}c&&/^[[:alpha:]][[:alnum:]_-]+:/{print substr($$1,1,index($$1,":")),c}1{c=0}' $(MAKEFILE_LIST) | column -s: -t
	@echo ""
	@echo "Run 'make <target>' to execute a specific target."

.PHONY: \
	default \
	build \
	test \
	lint \
	fmt \
	validate \
	coverprofile \
	dependencies \
	listen \
	help \
	$(NULL)
