<div align="center">
  <h1>🤖 V16 Client</h1>

  <h3>Ultra-Lightweight AI Agent for Your Computer</h3>

  <p>
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white" alt="Go">
    <img src="https://img.shields.io/badge/RAM-<50MB-green" alt="RAM">
    <img src="https://img.shields.io/badge/Startup-<2s-blue" alt="Startup">
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="License">
  </p>

  <p>
    <a href="https://v16.ai">Website</a> •
    <a href="#-quick-start">Quick Start</a> •
    <a href="#-features">Features</a> •
    <a href="CREDITS.md">Credits</a>
  </p>
</div>

## 📸 Demo: Telegram Bot in Action

<div align="center">
  <img src="https://pbs.twimg.com/media/HBO3CeRaUAAmYnW?format=jpg&name=large" width="250" alt="File operations">
  <img src="https://pbs.twimg.com/media/HBO3CePbQAAmiPR?format=jpg&name=large" width="250" alt="Web search">
  <img src="https://pbs.twimg.com/media/HBO3DrVasAEw2Cb?format=jpg&name=large" width="250" alt="Desktop analysis">

  <p><i>Real conversation with v16 bot: file operations, web search, desktop analysis</i></p>
</div>

### 🎥 Video Demos

**24/7 Autonomous AWS Log Monitoring**
Watch how V16 autonomously monitors AWS CloudWatch logs every hour and sends intelligent health reports to Telegram:

[![24/7 AWS Log Monitoring Agent](https://img.youtube.com/vi/7Q7MuoR7iaE/maxresdefault.jpg)](https://www.youtube.com/watch?v=7Q7MuoR7iaE)

**[Watch: Autonomous Agent Demo](https://www.youtube.com/watch?v=7Q7MuoR7iaE)** - See the agent analyze logs, detect errors, and provide actionable insights automatically.

---

## 🌟 What is V16 Client?

V16 Client is an **ultra-lightweight AI agent** that turns your computer into an intelligent assistant. With less than 50MB RAM usage and 2-second startup, it's 20x more efficient than traditional AI agents while providing more capabilities.

Connect to [v16.ai](https://v16.ai) to personalize your agent from anywhere, or run standalone with your own LLM provider.

### Key Highlights

- **🪶 Ultra-Lightweight**: <50MB RAM (vs 1GB+ for Node.js agents)
- **⚡ Lightning Fast**: 2-second startup
- **📦 Single Binary**: No Node.js, Python, or complex dependencies
- **🌍 Runs Anywhere**: From $10 boards to enterprise servers
- **🔌 Dual Mode**: Standalone OR connected to v16.ai platform

## 🚀 Quick Start

### Option 1: Telegram Bot (Easiest - 2 minutes!)

The fastest way to get started:

```bash
# 1. Build v16
git clone https://github.com/anup-singhai/v16-client.git
cd v16-client
make build
./build/v16 init

# 2. Create Telegram bot
# - Open Telegram, search for @BotFather
# - Send: /newbot
# - Follow prompts, copy your bot token

# 3. Configure (edit ~/.v16/config.json)
{
  "agents": {
    "defaults": {
      "provider": "moonshot",  // or anthropic, openai, etc.
      "model": "moonshot-v1-128k"
    }
  },
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "YOUR_BOT_TOKEN"
    }
  },
  "providers": {
    "moonshot": {
      "api_key": "sk-...",
      "api_base": "https://api.moonshot.ai/v1"
    }
  }
}

# 4. Start gateway
./build/v16 gateway

# 5. Chat with your bot on Telegram!
# Search for your bot and send: "List files in this directory"
```

### Option 2: CLI Mode

```bash
# 1. Build & Install
git clone https://github.com/anup-singhai/v16-client.git
cd v16-client
make build
make install

# 2. Initialize
v16 init

# 3. Configure provider (edit ~/.v16/config.json)
{
  "agents": {
    "defaults": {
      "provider": "anthropic",
      "model": "claude-3-5-sonnet-20241022"
    }
  },
  "providers": {
    "anthropic": {
      "api_key": "sk-ant-..."
    }
  }
}

# 4. Start using
v16 chat -m "Write a Python script to analyze CSV files"
v16 chat  # Interactive mode
```

### Option 3: Web UI (New! 🎨)

Configure your agent through a browser instead of editing JSON files:

```bash
# 1. Build & Install
git clone https://github.com/anup-singhai/v16-client.git
cd v16-client
make build
make install

# 2. Initialize
v16 init

# 3. Start Web UI
v16 web

# 4. Configure in browser
# Open http://localhost:8080
# - Set your LLM provider & API keys
# - Configure Telegram/Discord/Slack
# - Enable autonomous behaviors
# - Click Save

# 5. Start gateway with new config
v16 gateway
```

**Web UI Features:**
- 🎛️ Form-based configuration (no JSON editing)
- 📊 Live status dashboard
- 🔑 Provider API key management
- 📱 Channel configuration (Telegram, Discord, Slack)
- 💾 Auto-validation before saving

See [WEB_UI.md](WEB_UI.md) for details.

### Supported LLM Providers (14+)

| Provider | Models | API Base |
|----------|--------|----------|
| **Anthropic** | Claude 3.5 Sonnet, Opus, Haiku | https://api.anthropic.com |
| **Moonshot AI** | Kimi K2.5, moonshot-v1-128k | https://api.moonshot.ai/v1 |
| **OpenAI** | GPT-4, GPT-3.5 | https://api.openai.com/v1 |
| **Google** | Gemini Pro, Flash | https://generativelanguage.googleapis.com |
| **Groq** | Llama 3, Mixtral (fast) | https://api.groq.com/openai/v1 |
| **DeepSeek** | DeepSeek-V2 | https://api.deepseek.com |
| **Zhipu** | GLM-4 | https://open.bigmodel.cn/api/paas/v4 |
| **Local** | Ollama, vLLM | http://localhost |

And more: OpenRouter, Nvidia, GitHub Copilot

## ✨ Features

### Core Capabilities

| Category | Features |
|----------|----------|
| **Desktop** | ✅ Screen capture, mouse/keyboard control, window management |
| **Browser** | ✅ Navigate, fill forms, extract data, screenshots |
| **Terminal** | ✅ Interactive shell sessions, command execution, multi-session |
| **Files** | ✅ Read, write, edit, search (grep), pattern match (glob) |
| **Code** | ✅ Git operations (status, diff, log, commit, push, pull, branch) |
| **Web** | Search (Brave/DuckDuckGo), fetch content |
| **System** | Cron scheduling, task tracking |
| **Hardware** | I2C/SPI devices (Linux only) |

### AI Features

- **21 Tools**: Desktop, browser, terminal, files, code, web, hardware, and more
- **14+ LLM Providers**: Mix and match providers (Anthropic, Moonshot/Kimi, OpenAI, Gemini, etc.)
- **Multi-Channel**: **Telegram** (easiest!), Discord, Slack, WhatsApp, QQ, DingTalk, Feishu, LINE
- **Skills System**: Install and create custom skills
- **Memory**: Automatic context summarization
- **Scheduling**: Cron jobs and periodic tasks
- **Subagents**: Spawn async tasks in background

## 🚀 New V16 Capabilities

Phase 3 complete! The agent now has these powerful tools:

**Desktop Control** (8 actions)
- Screenshot capture and save
- Mouse movement and clicking (left/right/middle)
- Keyboard text input
- Screen size and mouse position info
- Window list and active window detection

**Browser Automation** (8 actions)
- Navigate to URLs
- Take full-page screenshots
- Click elements via CSS selectors
- Type into form fields
- Extract HTML content or element text
- Execute JavaScript
- Wait for elements to be visible

**Terminal/PTY** (6 actions)
- Open interactive shell sessions
- Send commands to sessions
- Get output from sessions
- Resize terminal windows
- Multi-session support
- Session management (list/close)

**Enhanced Code Tools** (3 tools)
- **Grep**: Fast code search with ripgrep, regex support, file filtering
- **Glob**: Pattern matching (**.go, **/*.ts), sorted by name/mtime
- **Git**: 9 operations (status, diff, log, commit, push, pull, branch, add, show)

All 21 tools are automatically available to the LLM during conversations!

## 🎯 Use Cases

**Personal Productivity**
- Schedule and manage tasks
- Automate repetitive workflows
- Research and summarize information

**Software Development**
- Write and review code
- Run tests and create PRs
- Search codebases

**IoT & Embedded**
- Control sensors via I2C/SPI
- Run on $10 boards (LicheeRV-Nano, Raspberry Pi)
- Edge AI applications

**Server Management**
- Monitor health and logs
- Automated maintenance
- Alert handling

## 🤖 Autonomous Agents (24/7 Monitoring & Automation)

V16 Client supports **autonomous agents** - AI agents that run on scheduled intervals (cron) to monitor systems, analyze data, and send reports automatically. No human intervention needed!

**📺 [Watch Demo: 24/7 AWS Log Monitoring Agent](https://youtu.be/7Q7MuoR7iaE)**

### Real-World Example: AWS Log Monitoring Agent

This autonomous agent monitors your AWS ECS logs every hour, analyzes errors, and sends detailed health reports to Telegram:

```json
{
  "agents": {
    "defaults": {
      "workspace": "~/.v16/workspace",
      "restrict_to_workspace": false,
      "provider": "moonshot",
      "model": "moonshot-v1-128k",
      "max_tokens": 8192,
      "temperature": 0.3,
      "max_tool_iterations": 20
    },
    "list": [
      {
        "id": "agent-logs-1",
        "name": "AWSLogMonitor",
        "personality": "You are a senior DevOps engineer analyzing AWS ECS logs. Parse ALL log events to extract: API endpoints, HTTP status code distribution, response times, error patterns. Provide INSIGHTS: identify trends (errors increasing?), anomalies (sudden spikes?), performance degradation, which endpoints failing. Show actual error messages with timestamps. Give actionable recommendations. Use exec tool with --profile production. Be specific, not generic.",
        "provider": "moonshot",
        "model": "moonshot-v1-128k",
        "max_tokens": 8192,
        "temperature": 0.3,
        "workspace": "~/.v16/workspace",
        "restrict_to_workspace": false,
        "enabled": true,
        "telegram_chat_id": "YOUR_TELEGRAM_CHAT_ID",
        "cron_jobs": [
          {
            "schedule": "0 0 * * * *",
            "task": "1) Execute: aws logs filter-log-events --log-group-name /ecs/your-api-service --start-time $(($(date +%s) - 3600))000 --profile production 2) Analyze the JSON output 3) Send message to Telegram with: total requests, success/error rates, top errors with snippets, performance insights, recommendations"
          }
        ]
      }
    ]
  },
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "YOUR_BOT_TOKEN"
    }
  },
  "providers": {
    "moonshot": {
      "api_key": "sk-...",
      "api_base": "https://api.moonshot.ai/v1"
    }
  }
}
```

### What This Agent Does

**Every Hour (automatically):**
1. ✅ Fetches last hour of AWS CloudWatch logs using AWS CLI
2. ✅ Parses JSON output to extract metrics
3. ✅ Analyzes HTTP status codes (200, 404, 500, etc.)
4. ✅ Identifies error patterns and trends
5. ✅ Calculates success/error rates and performance metrics
6. ✅ Sends comprehensive health report to Telegram

**Example Report Sent to Telegram:**
```
After analyzing the AWS ECS logs, here are the key insights:

1. Total Requests: 1,247
2. Success/Error Rates:
   - Success Rate: 94.3% (1,176 out of 1,247 requests)
   - Error Rate: 5.7% (71 out of 1,247 requests)

3. Top Errors with Snippets:
   - Error 404: 45 occurrences
     - Timestamp: 2026-02-15T14:23:12Z
     - Message: "GET /api/items/123 404 45ms"

   - Error 500: 18 occurrences
     - Timestamp: 2026-02-15T14:45:33Z
     - Message: "POST /api/users 500 234ms - Database connection timeout"

   - Error 503: 8 occurrences
     - Timestamp: 2026-02-15T15:12:08Z
     - Message: "GET /api/orders 503 1205ms - Service unavailable"

4. Performance Insights:
   - Average Response Time: 187ms
   - 95th Percentile: 421ms
   - Slowest Endpoint: /api/analytics (avg 1.2s)

5. Recommendations:
   - 🔴 CRITICAL: Investigate database connection timeouts on /api/users endpoint
   - 🟡 WARNING: /api/items/123 shows high 404 rate - check if resource exists
   - 🟢 OPTIMIZE: /api/analytics response time degrading - consider caching
```

### Cron Schedule Format

The schedule uses 6 fields (includes seconds):
```
┌─────────── second (0-59)
│ ┌─────────── minute (0-59)
│ │ ┌─────────── hour (0-23)
│ │ │ ┌─────────── day of month (1-31)
│ │ │ │ ┌─────────── month (1-12)
│ │ │ │ │ ┌─────────── day of week (0-6, Sunday=0)
│ │ │ │ │ │
│ │ │ │ │ │
0 0 * * * *  ← Run every hour (at :00)
0 */30 * * * *  ← Run every 30 minutes
0 0 */6 * * *  ← Run every 6 hours
0 0 9 * * *  ← Run daily at 9:00 AM
0 0 0 * * 1  ← Run every Monday at midnight
```

### More Autonomous Agent Ideas

**Security Monitoring**
```json
{
  "name": "SecurityMonitor",
  "personality": "You are a security analyst. Monitor system logs for suspicious activity, failed login attempts, unusual network traffic. Alert on potential security threats.",
  "cron_jobs": [{
    "schedule": "0 */15 * * * *",
    "task": "Check /var/log/auth.log for failed SSH attempts in last 15 mins. Send alert if >5 failures from same IP."
  }]
}
```

**Database Health Check**
```json
{
  "name": "DatabaseMonitor",
  "personality": "You are a DBA. Monitor database performance metrics, slow queries, connection pool usage, and disk space.",
  "cron_jobs": [{
    "schedule": "0 */5 * * * *",
    "task": "Execute: SELECT COUNT(*) FROM pg_stat_activity; Check slow queries. Report if connections >80% or queries >1s."
  }]
}
```

**API Uptime Monitor**
```json
{
  "name": "UptimeMonitor",
  "personality": "You are a reliability engineer. Test API endpoints for availability and response times.",
  "cron_jobs": [{
    "schedule": "0 */10 * * * *",
    "task": "Curl production endpoints. Check status codes and response times. Alert if any endpoint fails or >500ms response."
  }]
}
```

**Disk Space Monitor**
```json
{
  "name": "DiskSpaceMonitor",
  "personality": "You are a system administrator. Monitor disk usage across all mounted filesystems.",
  "cron_jobs": [{
    "schedule": "0 0 */4 * * *",
    "task": "Run: df -h. Check all partitions. Alert if any partition >85% full with cleanup recommendations."
  }]
}
```

### Interactive Intelligence

When you message the agent on Telegram, it intelligently responds to your questions:

**You:** "what are those 500 errors?"
**Agent:** *[Immediately queries AWS logs with filter for 500 errors]*
```
Found 18 instances of 500 errors in the last hour:

1. Timestamp: 2026-02-15T14:45:33Z
   Endpoint: POST /api/users
   Error: "Database connection timeout after 5000ms"

2. Timestamp: 2026-02-15T14:47:12Z
   Endpoint: POST /api/users
   Error: "Database connection timeout after 5000ms"

Root Cause: Database connection pool exhausted
Recommendation: Increase pool size or optimize slow queries
```

**You:** "look into logs"
**Agent:** *[Analyzes recent logs and provides detailed breakdown]*

The agent proactively investigates issues using real AWS CLI commands, not generic explanations!

### Key Features

- ✅ **24/7 Monitoring**: Runs autonomously on schedule
- ✅ **Smart Analysis**: Uses AI to extract insights from logs
- ✅ **Multi-Channel**: Send reports to Telegram, Slack, Discord
- ✅ **AWS Integration**: Full AWS CLI access with profile support
- ✅ **Interactive**: Ask follow-up questions, get real-time analysis
- ✅ **Actionable**: Provides specific recommendations, not just alerts
- ✅ **Persistent**: Auto-corrects errors, retries failed commands

### Getting Started

1. **Configure Telegram bot** (see Quick Start)
2. **Add autonomous agent** to `~/.v16/config.json`
3. **Set up AWS credentials** (if monitoring AWS)
4. **Start gateway**: `v16 gateway`
5. **Check Telegram** for hourly reports!

The agent runs in the background, continuously monitoring your systems and sending insights when it detects issues.

## 🏗️ Architecture

### Standalone Mode
Run locally with direct LLM access:
```
User → v16 chat → LLM Provider → Tools → Response
```

### Connected Mode (Coming Soon)
Connect to v16.ai for remote control:
```
v16.ai Dashboard → Socket.IO → v16 client → Tools → Response
```

**The Workflow:**
1. Visit [v16.ai](https://v16.ai) and sign up
2. Download v16 client (single binary)
3. Run `v16 connect --token <from dashboard>`
4. Personalize your agent via web interface
5. Agent executes locally with full system access

## 📊 Comparison

| Feature | Traditional Agents | V16 Client |
|---------|-------------------|------------|
| RAM Usage | 1GB+ | <50MB |
| Startup | 500s+ | <2s |
| Binary Size | 200MB+ | ~15MB |
| Dependencies | Node.js/Python | None |
| LLM Providers | 2-3 | 14+ |
| Channels | 1-2 | 11 |
| Hardware | Servers | $10+ boards |

## 🛠️ Building from Source

```bash
git clone https://github.com/anup-singhai/v16.git
cd v16

# Build
go build -o v16 ./cmd/v16

# Or use Make
make build

# Install
make install
```

## 🤝 Contributing

We welcome contributions!

- **Report bugs**: [GitHub Issues](https://github.com/anup-singhai/v16/issues)
- **Feature requests**: [Discussions](https://github.com/anup-singhai/v16/discussions)
- **Pull requests**: Always welcome!

## 📜 License

MIT License - see [LICENSE](LICENSE)

### Attribution

V16 Client is based on [PicoClaw](https://github.com/sipeed/picoclaw) by Sipeed, which was inspired by [nanobot](https://github.com/HKUDS/nanobot) by HKUDS.

See [CREDITS.md](CREDITS.md) for full acknowledgments.

## 🔗 Links

- **Website**: [v16.ai](https://v16.ai)
- **GitHub**: [github.com/anup-singhai/v16](https://github.com/anup-singhai/v16)
- **Original PicoClaw**: [github.com/sipeed/picoclaw](https://github.com/sipeed/picoclaw)

---

<div align="center">
  <p>Built with ❤️ by the V16 community</p>
  <p>Based on PicoClaw | Inspired by nanobot</p>
</div>
