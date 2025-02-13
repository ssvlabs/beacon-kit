.PHONY: build
build:
	go build ./...

.PHONY: clean
clean:
	go clean
	
.PHONY: test
test:
	go test ./...

# Generate test mocks with mockery and other things with gowrap.
generate:
	go install github.com/vektra/mockery/v2@latest
	go install github.com/hexdigest/gowrap/cmd/gowrap@latest
	(mockery --name Client && mv mocks/Client.go mocks/client.go)
	go generate ./...