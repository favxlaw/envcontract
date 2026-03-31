.PHONY: test lint build clean

test:
	go test ./... -race -cover

lint:
	golangci-lint run

build:
	go build ./cmd/envcontract/...

clean:
	rm -rf bin/