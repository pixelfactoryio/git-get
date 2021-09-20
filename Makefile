test:
	@go test -v -race -coverprofile coverage.txt -covermode atomic ./...
.PHONY: test

lint:
	@golangci-lint run ./...
.PHONY: lint

vet:
	@go vet ./...
.PHONY: vet
