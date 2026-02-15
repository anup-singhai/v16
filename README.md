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
