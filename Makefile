all: lint test
.PHONY: all

lint:
	golangci-lint run ./...
.PHONY: lint

test:
	go test -race ./...
.PHONY: test
