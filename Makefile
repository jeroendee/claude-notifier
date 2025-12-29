# Makefile for claude-notifier
# macOS menu bar notification app for Claude Code

# Variables
BINARY_NAME := claude-notifier
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR := bin
INSTALL_DIR := /usr/local/bin
LAUNCHAGENT_DIR := $(HOME)/Library/LaunchAgents
PLIST_NAME := com.dee.claude-notifier.plist
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

# Default target
.DEFAULT_GOAL := build

# Phony declarations
.PHONY: all build test lint clean install uninstall install-launchagent uninstall-launchagent install-hook help

# Build the binary
build:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/notifier ./cmd/notifier

# Run tests
test:
	go test -v ./...

# Run linter (optional, if golangci-lint is available)
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, skipping lint"; \
	fi

# Remove build artifacts
clean:
	rm -rf $(BUILD_DIR)/

# Install binary to /usr/local/bin
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)"
	install -m 755 $(BUILD_DIR)/notifier $(INSTALL_DIR)/$(BINARY_NAME)

# Remove installed binary
uninstall:
	@echo "Removing $(BINARY_NAME) from $(INSTALL_DIR)"
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)

# Install LaunchAgent for auto-start on login
install-launchagent: install
	@echo "Installing LaunchAgent to $(LAUNCHAGENT_DIR)"
	@mkdir -p $(LAUNCHAGENT_DIR)
	cp scripts/$(PLIST_NAME) $(LAUNCHAGENT_DIR)/$(PLIST_NAME)
	launchctl load $(LAUNCHAGENT_DIR)/$(PLIST_NAME)
	@echo "LaunchAgent installed and loaded"

# Uninstall LaunchAgent
uninstall-launchagent:
	@echo "Unloading and removing LaunchAgent"
	-launchctl unload $(LAUNCHAGENT_DIR)/$(PLIST_NAME) 2>/dev/null
	rm -f $(LAUNCHAGENT_DIR)/$(PLIST_NAME)
	@echo "LaunchAgent uninstalled"

# Install Claude Code hook script
install-hook:
	@echo "Installing Claude Code hook script"
	@mkdir -p $(HOME)/.claude/hooks
	cp scripts/claude-hook.sh $(HOME)/.claude/hooks/claude-hook.sh
	chmod +x $(HOME)/.claude/hooks/claude-hook.sh
	@echo ""
	@echo "Hook script installed to ~/.claude/hooks/claude-hook.sh"
	@echo ""
	@echo "To configure Claude Code to use this hook, add to ~/.claude/settings.json:"
	@echo '  {'
	@echo '    "hooks": {'
	@echo '      "Stop": ['
	@echo '        {'
	@echo '          "matcher": "",'
	@echo '          "hooks": ['
	@echo '            {'
	@echo '              "type": "command",'
	@echo '              "command": "~/.claude/hooks/claude-hook.sh"'
	@echo '            }'
	@echo '          ]'
	@echo '        }'
	@echo '      ]'
	@echo '    }'
	@echo '  }'

# Show help
help:
	@echo "Targets:"
	@echo "  build      Build the binary to $(BUILD_DIR)/notifier"
	@echo "  test       Run tests"
	@echo "  lint       Run linter (if golangci-lint is available)"
	@echo "  clean      Remove $(BUILD_DIR)/ directory"
	@echo "  install              Install binary to $(INSTALL_DIR)/$(BINARY_NAME)"
	@echo "  uninstall            Remove binary from $(INSTALL_DIR)"
	@echo "  install-launchagent  Install LaunchAgent for auto-start on login"
	@echo "  uninstall-launchagent Unload and remove LaunchAgent"
	@echo "  install-hook         Install Claude Code hook script"
	@echo "  help                 Show this help"
	@echo ""
	@echo "Variables:"
	@echo "  BINARY_NAME = $(BINARY_NAME)"
	@echo "  VERSION     = $(VERSION)"
	@echo "  BUILD_DIR   = $(BUILD_DIR)"
	@echo "  INSTALL_DIR = $(INSTALL_DIR)"
