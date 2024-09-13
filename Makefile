GOARCH=arm64
GOOS=linux
PREFIX=linux-arm64
BINARY_NAME=webnm
SOURCE_DIR=./cmd/webnm
BUILD_DIR=build

OUTPUT_NAME=$(BUILD_DIR)/$(PREFIX)-$(BINARY_NAME)

.PHONY: all
all: build check

.PHONY: build
build:
	@rm -rf $(BUILD_DIR)
	@mkdir -p $(BUILD_DIR)
	@echo "Building $(OUTPUT_NAME)..."
	GOARCH=$(GOARCH) GOOS=$(GOOS) go build -o $(OUTPUT_NAME) $(SOURCE_DIR)
	@echo "Build complete."

.PHONY: check
check:
	@if [ -f $(OUTPUT_NAME) ]; then \
		echo "Binary $(OUTPUT_NAME) exists."; \
	else \
		echo "Binary $(OUTPUT_NAME) does not exist. Build might have failed."; \
		exit 1; \
	fi

.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f $(BUILD_DIR)/$(PREFIX)-$(BINARY_NAME)
	@rmdir $(BUILD_DIR) 2>/dev/null || true
	@echo "Clean complete."