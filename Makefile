GET_TOOL=go get -modfile=tool.mod -tool
RUN_TOOL=go tool -modfile=tool.mod

.PHONY: build
build:
	go build ./...

.PHONY: clean
clean:
	go clean
	
.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	$(GET_TOOL) github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	$(RUN_TOOL) github.com/golangci/golangci-lint/v2/cmd/golangci-lint run -v ./...

.PHONY: tool
tool:
	$(GET_TOOL) github.com/vektra/mockery/v2@latest
	$(GET_TOOL) github.com/hexdigest/gowrap/cmd/gowrap@latest
	$(RUN_TOOL) github.com/vektra/mockery/v2 --name Client && mv mocks/Client.go mocks/client.go
	$(RUN_TOOL) github.com/hexdigest/gowrap/cmd/gowrap gen -p github.com/bloxapp/beacon-kit -i Client -t ./pool/methods.template -o ./pool/methods.go -l ""
