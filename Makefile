BINARY := linden
BUILD_DIR := bin

.PHONY: build install test tidy

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/linden

install:
	go install ./cmd/linden

test:
	go test ./...

test-skills:
	go test ./skills -count=1

tidy:
	go mod tidy
