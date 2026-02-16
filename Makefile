.PHONY: all build install uninstall clean help test

# Build variables
BINARY_NAME=v16
BUILD_DIR=build
CMD_DIR=cmd/$(BINARY_NAME)
MAIN_GO=$(CMD_DIR)/main.go

# Version
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT=$(shell git rev-parse --short=8 HEAD 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date +%FT%T%z)
GO_VERSION=$(shell $(GO) version | awk '{print $$3}')
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.gitCommit=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME) -X main.goVersion=$(GO_VERSION)"

# Go variables
GO?=go
GOFLAGS?=-v

# Installation
INSTALL_PREFIX?=$(HOME)/.local
INSTALL_BIN_DIR=$(INSTALL_PREFIX)/bin
INSTALL_MAN_DIR=$(INSTALL_PREFIX)/share/man/man1

# Workspace and Skills
V16_HOME?=$(HOME)/.v16
WORKSPACE_DIR?=$(V16_HOME)/workspace
WORKSPACE_SKILLS_DIR=$(WORKSPACE_DIR)/skills
BUILTIN_SKILLS_DIR=$(CURDIR)/skills

# OS detection
UNAME_S:=$(shell uname -s)
UNAME_M:=$(shell uname -m)

# Platform-specific settings
ifeq ($(UNAME_S),Linux)
	PLATFORM=linux
	ifeq ($(UNAME_M),x86_64)
		ARCH=amd64
	else ifeq ($(UNAME_M),aarch64)
		ARCH=arm64
	else ifeq ($(UNAME_M),riscv64)
		ARCH=riscv64
	else
		ARCH=$(UNAME_M)
	endif
else ifeq ($(UNAME_S),Darwin)
	PLATFORM=darwin
	ifeq ($(UNAME_M),x86_64)
		ARCH=amd64
	else ifeq ($(UNAME_M),arm64)
		ARCH=arm64
	else
		ARCH=$(UNAME_M)
	endif
else
	PLATFORM=$(UNAME_S)
	ARCH=$(UNAME_M)
endif

BINARY_PATH=$(BUILD_DIR)/$(BINARY_NAME)-$(PLATFORM)-$(ARCH)

# Default target
all: build

## generate: Run generate
generate:
	@echo "Run generate..."
	@rm -r ./$(CMD_DIR)/workspace 2>/dev/null || true
	@$(GO) generate ./...
	@echo "Run generate complete"

## build: Build the v16 binary for current platform
build: generate
	@echo "Building $(BINARY_NAME) for $(PLATFORM)/$(ARCH)..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_PATH) ./$(CMD_DIR)
	@echo "Build complete: $(BINARY_PATH)"
	@ln -sf $(BINARY_NAME)-$(PLATFORM)-$(ARCH) $(BUILD_DIR)/$(BINARY_NAME)

## build-all: Build v16 for all platforms
build-all: generate
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)
	GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./$(CMD_DIR)
	GOOS=linux GOARCH=riscv64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-riscv64 ./$(CMD_DIR)
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(CMD_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)
	@echo "All builds complete"

## install: Install v16 to system and copy builtin skills
install: build
	@echo "Installing $(BINARY_NAME)..."
	@mkdir -p $(INSTALL_BIN_DIR)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_BIN_DIR)/$(BINARY_NAME)
	@chmod +x $(INSTALL_BIN_DIR)/$(BINARY_NAME)
	@echo "Installed binary to $(INSTALL_BIN_DIR)/$(BINARY_NAME)"
	@echo "Installation complete!"

## uninstall: Remove v16 from system
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(INSTALL_BIN_DIR)/$(BINARY_NAME)
	@echo "Removed binary from $(INSTALL_BIN_DIR)/$(BINARY_NAME)"
	@echo "Note: Only the executable file has been deleted."
	@echo "If you need to delete all configurations (config.json, workspace, etc.), run 'make uninstall-all'"

## uninstall-all: Remove v16 and all data
uninstall-all:
	@echo "Removing workspace and skills..."
	@rm -rf $(V16_HOME)
	@echo "Removed workspace: $(V16_HOME)"
	@echo "Complete uninstallation done!"

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

## fmt: Format Go code
vet:
	@$(GO) vet ./...

## fmt: Format Go code
test:
	@$(GO) test ./...

## fmt: Format Go code
fmt:
	@$(GO) fmt ./...

## deps: Update dependencies
deps:
	@$(GO) get -u ./...
	@$(GO) mod tidy

## run: Build and run v16
run: build
	@$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

## whatsapp-setup: Setup WhatsApp bridge (one-time)
whatsapp-setup:
	@echo "Setting up WhatsApp bridge..."
	@mkdir -p ../whatsapp-bridge
	@cd ../whatsapp-bridge && npm init -y && npm install whatsapp-web.js ws qrcode-terminal
	@echo 'const { Client, LocalAuth } = require("whatsapp-web.js");' > ../whatsapp-bridge/server.js
	@echo 'const WebSocket = require("ws");' >> ../whatsapp-bridge/server.js
	@echo 'const qrcode = require("qrcode-terminal");' >> ../whatsapp-bridge/server.js
	@echo 'const client = new Client({ authStrategy: new LocalAuth(), puppeteer: { headless: true, args: ["--no-sandbox"] } });' >> ../whatsapp-bridge/server.js
	@echo 'const wss = new WebSocket.Server({ port: 3001 });' >> ../whatsapp-bridge/server.js
	@echo 'const clients = new Set();' >> ../whatsapp-bridge/server.js
	@echo 'client.on("qr", (qr) => { console.log("\\n📱 Scan QR:\\n"); qrcode.generate(qr, { small: true }); });' >> ../whatsapp-bridge/server.js
	@echo 'client.on("ready", () => console.log("✅ WhatsApp ready"));' >> ../whatsapp-bridge/server.js
	@echo 'client.on("message", async (msg) => { clients.forEach(ws => { if (ws.readyState === WebSocket.OPEN) ws.send(JSON.stringify({ type: "message", from: msg.from, content: msg.body, timestamp: Date.now() })); }); });' >> ../whatsapp-bridge/server.js
	@echo 'wss.on("connection", (ws) => { clients.add(ws); console.log("✅ v16 connected"); ws.on("message", async (data) => { const msg = JSON.parse(data); if (msg.type === "send") await client.sendMessage(msg.to, msg.content); }); ws.on("close", () => clients.delete(ws)); });' >> ../whatsapp-bridge/server.js
	@echo 'client.initialize(); console.log("🚀 Bridge: ws://localhost:3001");' >> ../whatsapp-bridge/server.js
	@echo "✅ WhatsApp bridge ready"

## whatsapp-start: Start WhatsApp + v16 gateway
whatsapp-start:
	@if [ ! -f ~/.v16/config.json ]; then echo "⚠️  Run './v16 init' first"; exit 1; fi
	@mkdir -p ~/.v16/logs
	@echo "Starting WhatsApp bridge..."
	@cd ../whatsapp-bridge && nohup node server.js > ~/.v16/logs/whatsapp.log 2>&1 & echo $$! > /tmp/whatsapp-bridge.pid
	@sleep 3
	@tail -20 ~/.v16/logs/whatsapp.log
	@echo ""
	@echo "Starting v16 gateway..."
	@$(BUILD_DIR)/$(BINARY_NAME) gateway

## whatsapp-stop: Stop WhatsApp services
whatsapp-stop:
	@if [ -f /tmp/whatsapp-bridge.pid ]; then kill $$(cat /tmp/whatsapp-bridge.pid) 2>/dev/null || true; rm /tmp/whatsapp-bridge.pid; fi
	@pkill -f "v16 gateway" || true
	@echo "✅ Stopped"

## help: Show this help message
help:
	@echo "v16 Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'
	@echo ""
	@echo "Examples:"
	@echo "  make build              # Build for current platform"
	@echo "  make install            # Install to ~/.local/bin"
	@echo "  make whatsapp-setup     # Setup WhatsApp (one-time)"
	@echo "  make whatsapp-start     # Start WhatsApp + gateway"
	@echo ""
	@echo "WhatsApp Setup:"
	@echo "  1. make whatsapp-setup"
	@echo "  2. ./build/v16 init && add API key to ~/.v16/config.json"
	@echo "  3. make whatsapp-start"
	@echo ""
	@echo "Environment Variables:"
	@echo "  INSTALL_PREFIX          # Installation prefix (default: ~/.local)"
	@echo "  WORKSPACE_DIR           # Workspace directory (default: ~/.v16/workspace)"
	@echo "  VERSION                 # Version string (default: git describe)"
	@echo ""
	@echo "Current Configuration:"
	@echo "  Platform: $(PLATFORM)/$(ARCH)"
	@echo "  Binary: $(BINARY_PATH)"
	@echo "  Install Prefix: $(INSTALL_PREFIX)"
	@echo "  Workspace: $(WORKSPACE_DIR)"
