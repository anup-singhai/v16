# V16 Web UI

Configure your V16 agents through a browser interface instead of editing JSON files manually.

## Quick Start

```bash
# Start the web UI
v16 web

# Custom port
v16 web --port 3000

# Custom address
v16 web --addr 0.0.0.0:8080
```

Then open your browser to: **http://localhost:8080**

## Features

### 1. **Status Dashboard**
- View current configuration path
- Check workspace location
- See active provider and model
- Monitor channel status (Telegram, Discord, Slack)

### 2. **Agent Configuration**
Configure your agent's behavior:
- **Workspace Path**: Where agent files are stored
- **Provider**: Choose LLM provider (OpenRouter, Anthropic, OpenAI, Gemini, Groq)
- **Model**: Specify model name
- **Max Tokens**: Response length limit
- **Temperature**: Response creativity (0-2)
- **Max Tool Iterations**: How many tool calls agent can make
- **Restrict to Workspace**: Prevent file access outside workspace

### 3. **Channel Configuration**
Enable and configure communication channels:

**Telegram**
- Enable/disable
- Bot token (get from @BotFather)

**Discord**
- Enable/disable
- Bot token

**Slack**
- Enable/disable
- Bot token

### 4. **Provider API Keys**
Add API keys for LLM providers:
- OpenRouter
- Anthropic
- OpenAI
- Gemini
- Groq

## Usage Flow

1. **Start Web UI**
   ```bash
   v16 web
   ```

2. **Configure in Browser**
   - Open http://localhost:8080
   - Navigate through tabs: Status, Agents, Channels, Providers
   - Fill in your configuration
   - Click "Save Configuration"

3. **Restart Gateway**
   ```bash
   # Stop existing gateway (Ctrl+C)
   # Start with new config
   v16 gateway
   ```

## Security Notes

- The web UI runs on localhost by default (not accessible from network)
- To allow network access, use: `v16 web --addr 0.0.0.0:8080`
- **Warning**: Don't expose the web UI to the internet without authentication
- API keys are saved to `~/.v16/config.json` in plain text
- Use firewall rules if running on a server

## CLI vs Web UI

| Feature | CLI | Web UI |
|---------|-----|--------|
| Edit config | Manual JSON editing | Form-based |
| Validation | None | Real-time |
| Documentation | Separate | Inline help text |
| Learning curve | Steeper | Gentle |
| Power users | ✅ Faster | ❌ Slower |
| Beginners | ❌ Harder | ✅ Easier |

## Configuration File Location

Both CLI and Web UI modify the same file:
```
~/.v16/config.json
```

You can edit this file directly or use the web UI - changes are synced.

## Troubleshooting

### Web UI won't start
```bash
# Check if port is in use
lsof -i :8080

# Use different port
v16 web --port 8081
```

### Can't access from network
```bash
# Bind to all interfaces (⚠️ security risk)
v16 web --addr 0.0.0.0:8080
```

### Changes not applied
After saving configuration via web UI, restart the gateway:
```bash
pkill v16  # If running in background
v16 gateway
```

### API keys not working
- Ensure no trailing spaces
- Check provider-specific key format:
  - OpenRouter: `sk-or-v1-...`
  - Anthropic: `sk-ant-...`
  - OpenAI: `sk-...`
  - Groq: `gsk_...`
  - Gemini: `AIza...`

## Advanced Usage

### Running as Background Service

**systemd (Linux)**
```ini
[Unit]
Description=V16 Web UI
After=network.target

[Service]
Type=simple
User=your-username
WorkingDirectory=/home/your-username
ExecStart=/usr/local/bin/v16 web --addr localhost:8080
Restart=always

[Install]
WantedBy=multi-user.target
```

**launchd (macOS)**
Create `~/Library/LaunchAgents/ai.v16.web.plist`:
```xml
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
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
```

Then:
```bash
launchctl load ~/Library/LaunchAgents/ai.v16.web.plist
launchctl start ai.v16.web
```

### Nginx Reverse Proxy

To add authentication and HTTPS:

```nginx
server {
    listen 443 ssl;
    server_name v16.yourdomain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        auth_basic "V16 Admin";
        auth_basic_user_file /etc/nginx/.htpasswd;

        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## API Endpoints

The web UI exposes these endpoints:

- `GET /api/config` - Get current configuration
- `POST /api/config/save` - Save configuration
- `GET /api/status` - Get system status

You can use these programmatically:

```bash
# Get config
curl http://localhost:8080/api/config

# Get status
curl http://localhost:8080/api/status

# Save config
curl -X POST http://localhost:8080/api/config/save \
  -H "Content-Type: application/json" \
  -d @config.json
```

## Future Features

- [ ] Live agent chat interface
- [ ] Cron job management
- [ ] Skill installation
- [ ] Log viewer
- [ ] System metrics
- [ ] Multi-agent management
- [ ] Authentication/authorization
- [ ] Dark/light theme toggle

## Contributing

Want to improve the web UI? Contributions welcome!

```bash
# UI files location
pkg/web/static/index.html  # HTML structure
pkg/web/static/app.js      # JavaScript logic
pkg/web/server.go          # Go backend

# After changes, rebuild
go build -o v16 cmd/v16/main.go
```
