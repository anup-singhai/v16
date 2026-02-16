# Changelog

All notable changes to V16 Client will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-02-15

### Added - V16 Capabilities

#### Desktop Automation Tool
- Screenshot capture with custom save paths
- Mouse movement and clicking (left/right/middle buttons)
- Keyboard text input automation
- Screen size and mouse position detection
- Window list and active window detection
- Uses `robotgo` library for cross-platform support

#### Browser Automation Tool
- Navigate to URLs
- Take full-page screenshots
- Click elements via CSS selectors
- Type text into form fields
- Extract HTML content or element text
- Execute custom JavaScript
- Wait for elements to be visible
- Uses `chromedp` (Chrome DevTools Protocol)
- Headless Chrome by default

#### Terminal/PTY Tool
- Open interactive shell sessions (bash/zsh)
- Multi-session support with UUID identifiers
- Send commands to sessions
- Retrieve output from sessions
- Resize terminal windows
- Session management (list/close)
- Thread-safe with mutex locks
- Uses `creack/pty` library

#### Enhanced Code Tools
- **Grep Tool**: Fast recursive code search
  - Uses ripgrep with fallback to grep
  - Regex pattern support
  - File pattern filtering (*.go, *.js, etc.)
  - Case-sensitive/insensitive options
  - Context lines around matches
  - Configurable max results
  - Output truncation for large results

- **Glob Tool**: File pattern matching
  - Supports ** patterns (e.g., **/*.ts)
  - Sort by name or modification time
  - Include/exclude directories
  - Relative path display
  - Configurable max results

- **Git Tool**: Repository operations
  - 9 operations: status, diff, log, commit, push, pull, branch, add, show
  - Workspace-aware execution
  - Output truncation for large diffs
  - Comprehensive error handling

#### Platform Integration
- V16 connector framework for v16.ai platform
- Command bridge with 30+ command mappings
- Connection lifecycle management
- CLI commands: `v16 connect --token`, `v16 disconnect`
- Agent helper methods: GetTool(), GetToolNames()
- WebSocket client implementation pending (v0.2.0)

### Changed - Branding & Architecture

#### Rebranding from PicoClaw
- Binary renamed: `picoclaw` → `v16`
- Configuration path: `~/.picoclaw/` → `~/.v16/`
- Commands updated:
  - `picoclaw onboard` → `v16 init`
  - `picoclaw agent` → `v16 chat`
  - `picoclaw cron` → `v16 schedule`
- Package path: `github.com/sipeed/picoclaw` → `github.com/anup-singhai/v16`

#### Documentation
- New README with V16 branding
- Added CREDITS.md with full attribution chain
- MIT License with proper attribution
- Updated all help text and CLI output

### Technical Details

#### Dependencies Added
- `github.com/go-vgo/robotgo` v1.0.0 - Desktop automation
- `github.com/chromedp/chromedp` v0.14.2 - Browser automation
- `github.com/creack/pty` v1.1.24 - Terminal/PTY support
- `github.com/googollee/go-socket.io` v1.7.0 - WebSocket (placeholder)

#### Tool Count
- **21 Total Tools** (up from 15)
  - 15 original PicoClaw tools preserved
  - 3 V16 desktop/browser/terminal tools
  - 3 enhanced code tools (grep, glob, git)

#### Performance
- Binary size: ~27MB (darwin/arm64)
- RAM usage: <50MB (maintained)
- Startup time: <2s (maintained)
- Build: Go 1.21+

### Preserved from PicoClaw

All original PicoClaw features maintained:
- 15 core tools (files, shell, web, hardware, cron, etc.)
- 14+ LLM providers (OpenAI, Anthropic, Gemini, DeepSeek, etc.)
- 11 communication channels (Telegram, Discord, Slack, etc.)
- Skills system
- Heartbeat system
- Memory management with auto-summarization
- Subagent spawning
- Security sandbox
- Hardware integration (I2C/SPI on Linux)

## Attribution

V16 Client is based on [PicoClaw](https://github.com/sipeed/picoclaw) by Sipeed,
which was inspired by [nanobot](https://github.com/HKUDS/nanobot) by HKUDS.

See [CREDITS.md](CREDITS.md) for detailed acknowledgments.

[0.1.0]: https://github.com/anup-singhai/v16/releases/tag/v0.1.0
