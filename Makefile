test:
	go test -v -cover ./...

golangci-lint:
	golangci-lint run ./...

.PHONY: test golangci-lint