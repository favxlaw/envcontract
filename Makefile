.PHONY: test lint clean build

test:
	go test ./... -race -cover

lint:
	golangci-lint run

build:
	go build -o bin/envcontract ./cmd/envcontract/

clean:
	rm -rf bin/
