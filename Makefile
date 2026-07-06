.PHONY: test lint clean

test:
	go test ./... -race -cover

lint:
	golangci-lint run

clean:
	rm -rf bin/