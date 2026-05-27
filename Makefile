all: help

.PHONY: help
help: Makefile
	@sed -n 's/^##//p' $< | awk 'BEGIN {FS = ":"}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## test: Run tests for all modules
.PHONY: test
test:
	@go test $(TESTARGS) ./...
	@cd test && go test $(TESTARGS) ./...
	@cd example && go test $(TESTARGS) ./...

## test-race: Run tests with the Go race detector
.PHONY: test-race
test-race:
	@$(MAKE) test TESTARGS=-race

## lint: Run golangci-lint
.PHONY: lint
lint:
	@./scripts/lint -c .golangci.yml

## fmt: Format Go files
.PHONY: fmt
fmt:
	@gofmt -w $$(find . -name '*.go' -not -path './.tools/*')

## ci: Run lint and tests
.PHONY: ci
ci: lint test test-race
