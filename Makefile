all: help

.PHONY : help
help : Makefile
	@sed -n 's/^##//p' $< | awk 'BEGIN {FS = ":"}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

TOOLS_MOD_DIR := ./tools
TOOLS_DIR := $(abspath ./.tools)
$(TOOLS_DIR)/golangci-lint: $(TOOLS_MOD_DIR)/go.mod $(TOOLS_MOD_DIR)/go.sum $(TOOLS_MOD_DIR)/tools.go
	@echo BUILD golangci-lint
	@cd $(TOOLS_MOD_DIR) && \
	go build -o $(TOOLS_DIR)/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

$(TOOLS_DIR)/godoc: $(TOOLS_MOD_DIR)/go.mod $(TOOLS_MOD_DIR)/go.sum $(TOOLS_MOD_DIR)/tools.go
	@echo BUILD godoc
	@cd $(TOOLS_MOD_DIR) && \
	go build -o $(TOOLS_DIR)/godoc golang.org/x/tools/cmd/godoc

## tools: Build all tools
tools: $(TOOLS_DIR)/golangci-lint

## lint: Run golangci-lint
.PHONY: lint
lint: $(TOOLS_DIR)/golangci-lint
	@echo LINT
	@$(TOOLS_DIR)/golangci-lint run -c .github/linters/.golangci.yaml --out-format colored-line-number
	@printf "LINT... \033[0;32m [OK] \033[0m\n"

## test: Run test
.PHONY: test
test:
	@echo TEST
	@go test  ./...
	@printf "TEST... \033[0;32m [OK] \033[0m\n"

## test/coverage: Run test and generate coverage report
.PHONY: test/coverage
test/coverage:
	@go test ./... -coverprofile=coverage.txt -covermode=atomic
	@go tool cover -html=coverage.txt -o coverage.html

.PHONY: $(BIN_DIR)
$(BIN_DIR):
	@mkdir -p $@

## godoc: View godoc
PKG_NAME:=$(shell cat go.mod | grep module | cut -d' ' -f2)
.PHONY: godoc
godoc: $(TOOLS_DIR)/godoc
	@echo "Open http://localhost:6060/pkg/$(PKG_NAME) on browser."
	$(TOOLS_DIR)/godoc -http localhost:6060

.PHONY: lint-ci
lint-ci: $(TOOLS_DIR)/golangci-lint
	@echo LINT
	@$(TOOLS_DIR)/golangci-lint run -c .github/linters/.golangci.yaml
	@printf "LINT... \033[0;32m [OK] \033[0m\n"
