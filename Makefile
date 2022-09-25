.DEFAULT_GOAL := help

.PHONY: test coverage help html html-gen html-wsl

test: ## Run unit tests
	@go test -race ./...

coverage: ## Run unit tests with coverage
	@go test -v -race -cover -coverpkg=./... -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out

html:
	@go tool cover -html=./coverage.out

html-gen: ## explorer.exe ./coverage.html
	@go tool cover -html=./coverage.out -o ./coverage.html

html-wsl: html-gen ## wslview ./coverage.html
	@explorer.exe `wslpath -w ./coverage.html` || true

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'