MODULE := github.com/narvanalabs/flkr
VERSION ?= dev
LDFLAGS := -ldflags "-X $(MODULE)/cmd.Version=$(VERSION)"

.PHONY: build test vet lint clean

build:
	go build $(LDFLAGS) -o flkr .

test:
	go test ./... -v

vet:
	go vet ./...

lint: vet
	@echo "lint OK"

clean:
	rm -f flkr

all: vet test build
