.PHONY: build test lint clean

build:
	go build -o bin/envcontract ./cmd/envcontract

test:
	go test ./... -race -cover

lint:
	golangci-lint run

clean:
	rm -rf bin/