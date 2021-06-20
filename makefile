.DEFAULT_GOAL=help

help:
	@echo "Usage:"
	@echo "  make [target...]"
	@echo ""
	@echo "Useful commands:"
	@grep -Eh '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-30s %s\n", $$1, $$2}'
	@echo ""


build: ## build the binary(binary name- score) in the current working directory
	@go build -o score ./cmd/cli

move: ## move to /usr/bin so that you can use this binary anywhere.
	@sudo mv score /usr/local/bin

run: build ## run the CLI
	@./score
