# Makefile for claude-notifier
# macOS menu bar notification app for Claude Code

# Variables
BINARY_NAME := claude-notifier
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR := bin
LAUNCHAGENT_DIR := $(HOME)/Library/LaunchAgents
PLIST_NAME := com.dee.claude-notifier.plist
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -ldflags "-X github.com/jeroendee/claude-notifier/internal/version.Version=$(VERSION) -X github.com/jeroendee/claude-notifier/internal/version.Commit=$(COMMIT) -X github.com/jeroendee/claude-notifier/internal/version.BuildDate=$(BUILD_DATE)"

# Default target
.DEFAULT_GOAL := build

# Phony declarations
.PHONY: all build test coverage lint clean install uninstall install-launchagent uninstall-launchagent install-hook help

# Build the binary
build:
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/claude-notifier

# Run tests
test:
	go test -v ./...

# Run tests with coverage report
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run linter (optional, if golangci-lint is available)
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, skipping lint"; \
	fi

# Remove build artifacts
clean:
	rm -rf $(BUILD_DIR)/ coverage.out coverage.html

# Install binary to GOPATH/bin (or GOBIN if set)
install:
	go install $(LDFLAGS) ./cmd/claude-notifier

# Remove installed binary from GOPATH/bin
uninstall:
	@GOBIN=$$(go env GOBIN); \
	if [ -z "$$GOBIN" ]; then GOBIN=$$(go env GOPATH)/bin; fi; \
	echo "Removing $(BINARY_NAME) from $$GOBIN"; \
	rm -f "$$GOBIN/$(BINARY_NAME)"

# Install LaunchAgent for auto-start on login
install-launchagent: install
	@echo "Installing LaunchAgent to $(LAUNCHAGENT_DIR)..."
	@go run ./cmd/install-launchagent
	launchctl load $(LAUNCHAGENT_DIR)/$(PLIST_NAME)
	@echo "LaunchAgent installed and loaded"

# Uninstall LaunchAgent
uninstall-launchagent:
	@echo "Unloading and removing LaunchAgent"
	-launchctl unload $(LAUNCHAGENT_DIR)/$(PLIST_NAME) 2>/dev/null
	rm -f $(LAUNCHAGENT_DIR)/$(PLIST_NAME)
	@echo "LaunchAgent uninstalled"

# Install Claude Code hook scripts
install-hook:
	@echo "Installing Claude Code hook scripts..."
	@go run ./cmd/install-hook
	@echo ""
	@echo "To configure Claude Code to use these hooks, add to ~/.claude/settings.json:"
	@echo '  {'
	@echo '    "hooks": {'
	@echo '      "Stop": [{'
	@echo '        "hooks": [{ "type": "command", "command": "~/.claude/hooks/claude-hook.sh" }]'
	@echo '      }],'
	@echo '      "Notification": [{'
	@echo '        "matcher": "idle_prompt",'
	@echo '        "hooks": [{ "type": "command", "command": "~/.claude/hooks/notification-hook.sh" }]'
	@echo '      }],'
	@echo '      "PermissionRequest": [{'
	@echo '        "matcher": "",'
	@echo '        "hooks": [{ "type": "command", "command": "~/.claude/hooks/notification-hook.sh" }]'
	@echo '      }]'
	@echo '    }'
	@echo '  }'
	@echo ""
	@echo "Hook types:"
	@echo "  - Stop: Notifies when Claude finishes responding (with SESSION summary)"
	@echo "  - Notification (idle_prompt): Notifies after 60s waiting for user input"
	@echo "  - PermissionRequest: Notifies when Claude needs permission approval"

# Show help
help:
	@echo "Targets:"
	@echo "  build      Build the binary to $(BUILD_DIR)/$(BINARY_NAME)"
	@echo "  test       Run tests"
	@echo "  coverage   Run tests with HTML coverage report"
	@echo "  lint       Run linter (if golangci-lint is available)"
	@echo "  clean      Remove $(BUILD_DIR)/ directory"
	@echo "  install              Install binary to GOPATH/bin"
	@echo "  uninstall            Remove binary from GOPATH/bin"
	@echo "  install-launchagent  Install LaunchAgent for auto-start on login"
	@echo "  uninstall-launchagent Unload and remove LaunchAgent"
	@echo "  install-hook         Install Claude Code hook script"
	@echo "  help                 Show this help"
	@echo ""
	@echo "Variables:"
	@echo "  BINARY_NAME = $(BINARY_NAME)"
	@echo "  VERSION     = $(VERSION)"
	@echo "  BUILD_DIR   = $(BUILD_DIR)"
