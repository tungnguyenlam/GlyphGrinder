BINARY := glyphgrinder

.PHONY: verify
verify: ## Run the canonical verification suite (fmt, vet, build, smoke, tests, tidy)
	./scripts/verify.sh

.PHONY: run
run: ## Play the game (needs a real terminal / TTY)
	go run .

.PHONY: build
build: ## Build the binary into ./$(BINARY)
	go build -o $(BINARY) .

.PHONY: test
test: ## Run tests only (verify is the command docs reference)
	go test -count=1 ./...

.PHONY: fmt
fmt: ## Format all Go files
	gofmt -w .

.PHONY: clean
clean:
	rm -f $(BINARY)

.PHONY: help
help:
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'
