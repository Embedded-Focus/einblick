.PHONY: build test race vet fmt lint check

build:
	CGO_ENABLED=0 go build -trimpath ./cmd/einblick

test:
	go test ./...

race:
	go test -race ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

lint:
	golangci-lint run

check: fmt vet test race lint build
