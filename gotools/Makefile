BIN_DIR=$(abspath ./bin)
BIN_PATH="$(BIN_DIR):$(PATH)"

.PHONY: tools
tools:
	GOBIN=$(BIN_DIR) go -C tools install tool

.PHONY: gen
gen: tools
	@PATH=$(BIN_PATH) go generate ./...
	@PATH=$(BIN_PATH) buf generate