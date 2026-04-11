BINARY      := portwatch
CMD_PKG     := ./cmd/portwatch
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS     := -ldflags "-X main.version=$(VERSION)"
OUT_DIR     := bin

.PHONY: all build test lint clean run

all: build

build:
	@mkdir -p $(OUT_DIR)
	go build $(LDFLAGS) -o $(OUT_DIR)/$(BINARY) $(CMD_PKG)

test:
	go test ./...

test-integration:
	PORTWATCH_RUN_INTEGRATION=1 go test ./...

lint:
	@which golangci-lint > /dev/null 2>&1 || (echo "golangci-lint not found" && exit 1)
	golangci-lint run ./...

run: build
	$(OUT_DIR)/$(BINARY)

clean:
	rm -rf $(OUT_DIR)
