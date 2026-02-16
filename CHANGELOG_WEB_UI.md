# Web UI Feature - Changelog

## Summary

Added a web-based configuration UI to v16-client, making it easier for users to configure agents without editing JSON files manually.

## New Files

### Backend
- **`pkg/web/server.go`** - HTTP server with REST API
  - `GET /api/config` - Retrieve current configuration
  - `POST /api/config/save` - Save configuration
  - `GET /api/status` - Get system status

### Frontend
- **`pkg/web/static/index.html`** - Web UI HTML
  - 4 tabs: Status, Agents, Channels, Providers
  - Dark theme matching V16 brand colors (#8B7355)
  - Responsive design

- **`pkg/web/static/app.js`** - Frontend JavaScript
  - Real-time config loading
  - Form validation
  - API integration
  - Success/error messaging

### Documentation
- **`WEB_UI.md`** - Complete web UI documentation
  - Quick start guide
  - Features overview
  - Security notes
  - Troubleshooting
  - Advanced usage (systemd, nginx)

## Modified Files

### `cmd/v16/main.go`
- Added `web` package import
- Added `webCmd()` function
- Added `web` and `ui` commands to switch statement
- Updated help text with web UI instructions

### `README.md`
- Added "Option 3: Web UI" section
- Highlighted key web UI features

## Usage

```bash
# Basic usage
v16 web

# Custom port
v16 web --port 3000

# Custom address (bind to all interfaces)
v16 web --addr 0.0.0.0:8080
```

## Features

### Status Dashboard
- Config file path
- Workspace location
- Active provider & model
- Channel status indicators (enabled/disabled)

### Agent Configuration
- Workspace path
- Provider selection (OpenRouter, Anthropic, OpenAI, Gemini, Groq)
- Model name
- Max tokens
- Temperature
- Max tool iterations
- Restrict to workspace toggle

### Channel Configuration
- **Telegram**: Enable/disable, bot token
- **Discord**: Enable/disable, bot token
- **Slack**: Enable/disable, bot token

### Provider API Keys
- OpenRouter
- Anthropic
- OpenAI
- Gemini
- Groq

## Technical Details

### Architecture
```
Browser (http://localhost:8080)
    ↓
HTTP Server (pkg/web/server.go)
    ↓
Config File (~/.v16/config.json)
    ↓
V16 Gateway (reads config on startup)
```

### API Endpoints

**GET /api/config**
```json
{
  "agents": {...},
  "channels": {...},
  "providers": {...}
}
```

**POST /api/config/save**
```json
{
  "success": true,
  "message": "Configuration saved successfully"
}
```

**GET /api/status**
```json
{
  "config_path": "~/.v16/config.json",
  "workspace": "~/workspace",
  "provider": "anthropic",
  "model": "claude-3-5-sonnet",
  "channels": {
    "telegram": true,
    "discord": false,
    "slack": false
  }
}
```

### UI Design

**Color Palette** (matching V16 brand):
- Primary: #8B7355 (brown)
- Secondary: #A67C52 (lighter brown)
- Background: #0a0a0a (dark)
- Cards: #1a1a1a
- Borders: #333
- Text: #e0e0e0

**Responsive Breakpoints**:
- Mobile: < 768px (single column)
- Desktop: >= 768px (grid layout)

## Security Considerations

1. **Default Binding**: localhost only (not exposed to network)
2. **No Authentication**: Add nginx reverse proxy for auth
3. **Plain Text Storage**: API keys saved to config.json unencrypted
4. **Network Exposure**: Use `--addr 0.0.0.0:8080` with caution

## Future Enhancements

- [ ] Live agent chat interface
- [ ] Cron job visual editor
- [ ] Skill marketplace integration
- [ ] Real-time log viewer
- [ ] System metrics (CPU, RAM, disk)
- [ ] Multi-agent management
- [ ] Built-in authentication
- [ ] Theme switcher (dark/light)
- [ ] WebSocket for live updates
- [ ] Config export/import

## Testing

```bash
# Build
cd /Users/anupsingh/projects/v16/v16-client
go build -o build/v16 cmd/v16/main.go

# Test web command
./build/v16 web

# In browser
open http://localhost:8080

# Test API
curl http://localhost:8080/api/status
curl http://localhost:8080/api/config
```

## Benefits

**For Beginners:**
- No JSON syntax errors
- Visual validation
- Inline help text
- Easier onboarding

**For Advanced Users:**
- Quick config changes
- Visual status monitoring
- Remote configuration (with nginx)
- API for automation

## Comparison: CLI vs Web UI

| Task | CLI | Web UI |
|------|-----|--------|
| Initial setup | `v16 init` + edit JSON | `v16 web` + forms |
| Add Telegram bot | Edit JSON, find syntax | Enable toggle, paste token |
| Change model | Edit JSON string | Select dropdown |
| Verify config | Read JSON | Visual status tab |
| API key typo | Silent failure | Immediate validation |
| Learning curve | Steeper | Gentle |
| Power user speed | Faster | Slower |

## Deployment Example

### Linux systemd
```bash
# /etc/systemd/system/v16-web.service
[Unit]
Description=V16 Web UI
After=network.target

[Service]
Type=simple
User=v16
ExecStart=/usr/local/bin/v16 web --addr localhost:8080
Restart=always

[Install]
WantedBy=multi-user.target
```

### macOS launchd
```xml
<!-- ~/Library/LaunchAgents/ai.v16.web.plist -->
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>ai.v16.web</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/v16</string>
        <string>web</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>
```

### Docker
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o v16 cmd/v16/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/v16 /usr/local/bin/
EXPOSE 8080
CMD ["v16", "web", "--addr", "0.0.0.0:8080"]
```

## Version

- **Added in**: v1.1.0
- **Status**: Beta
- **Tested on**: macOS, Linux
- **Browser Support**: Chrome, Firefox, Safari, Edge

## Contributors

- Initial implementation: Claude Code session
- Testing: V16 team
- Documentation: Integrated

## License

MIT (same as v16-client)
