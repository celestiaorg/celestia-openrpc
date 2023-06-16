DOCKER := $(shell which docker)

## help: Show this help message
help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
.PHONY: help

## clean: clean testcache
clean:
	@echo "--> Clearing testcache"
	@go clean --testcache
.PHONY: clean

## cover: generate to code coverage report.
cover:
	@echo "--> Generating Code Coverage"
	@go install github.com/ory/go-acc@latest
	@go-acc -o coverage.txt `go list ./...`
.PHONY: cover

## lint: Run linters golangci-lint and markdownlint.
lint:
	@echo "--> Running golangci-lint"
	@golangci-lint run
	@echo "--> Running markdownlint"
	@markdownlint --config .markdownlint.yaml '**/*.md'
	@echo "--> Running yamllint"
	@yamllint --no-warnings . -c .yamllint.yml

.PHONY: lint

## test-unit: Running unit tests
test-unit:
	@echo "--> Running unit tests"
	@go test `go list ./...`
.PHONY: test-unit

## test-unit-race: Running unit tests with data race detector
test-unit-race:
	@echo "--> Running unit tests with data race detector"
	@go test -race -count=1 `go list ./...`
.PHONY: test-unit-race

### test-all: Run tests with and without data race
test-all:
	@$(MAKE) test-unit
	@$(MAKE) test-unit-race
.PHONY: test-all

proto-gen:
	@echo "Generating Protobuf files"
	@go run github.com/bufbuild/buf/cmd/buf generate
.PHONY: proto-gen

# These targets are provided for convenience and are intended for local
# execution only.
proto-lint:
	@echo "Linting Protobuf files"
	@go run github.com/bufbuild/buf/cmd/buf lint
.PHONY: proto-lint

proto-format:
	@echo "Formatting Protobuf files"
	@find . -name '*.proto' -path "./proto/*" -exec clang-format -i {} \;
.PHONY: proto-format
