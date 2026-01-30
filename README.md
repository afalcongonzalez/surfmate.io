# Surfmate.io

Browser automation MCP server for AI. Opens a visible browser so you can handle captchas and logins while AI controls navigation.

## Quick Install

### macOS (Apple Silicon)

```bash
curl -L -o surfmate.io https://github.com/afalcongonzalez/surfmate.io/releases/latest/download/surfmate.io-darwin-arm64
xattr -d com.apple.quarantine surfmate.io
chmod +x surfmate.io
sudo mv surfmate.io /usr/local/bin/
```

### macOS (Intel)

```bash
curl -L -o surfmate.io https://github.com/afalcongonzalez/surfmate.io/releases/latest/download/surfmate.io-darwin-amd64
xattr -d com.apple.quarantine surfmate.io
chmod +x surfmate.io
sudo mv surfmate.io /usr/local/bin/
```

### Linux (x64)

```bash
curl -L -o surfmate.io https://github.com/afalcongonzalez/surfmate.io/releases/latest/download/surfmate.io-linux-amd64
chmod +x surfmate.io
sudo mv surfmate.io /usr/local/bin/
```

### Linux (ARM64)

```bash
curl -L -o surfmate.io https://github.com/afalcongonzalez/surfmate.io/releases/latest/download/surfmate.io-linux-arm64
chmod +x surfmate.io
sudo mv surfmate.io /usr/local/bin/
```

### Windows (PowerShell as Admin)

```powershell
Invoke-WebRequest -Uri "https://github.com/afalcongonzalez/surfmate.io/releases/latest/download/surfmate.io-windows-amd64.exe" -OutFile "surfmate.io.exe"
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\bin"
Move-Item surfmate.io.exe "$env:USERPROFILE\bin\"
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$env:USERPROFILE\bin", "User")
```

### Troubleshooting

**macOS "cannot be opened" error:**
```bash
xattr -d com.apple.quarantine surfmate.io
```

**Windows SmartScreen warning:** Click "More info" → "Run anyway"

## Requirements

- Chromium-based browser (Chrome, Edge, or Brave)

---

## Setup for Claude Desktop

Add to your Claude Desktop config:

| OS | Config Path |
|----|-------------|
| macOS | `~/Library/Application Support/Claude/claude_desktop_config.json` |
| Windows | `%APPDATA%\Claude\claude_desktop_config.json` |
| Linux | `~/.config/Claude/claude_desktop_config.json` |

```json
{
  "mcpServers": {
    "surfmate": {
      "command": "surfmate.io"
    }
  }
}
```

Restart Claude Desktop after saving.

---

## Setup for Claude Code

Add to your project's `.claude/settings.json`:

```json
{
  "mcpServers": {
    "surfmate": {
      "command": "surfmate.io"
    }
  }
}
```

---

## Setup for ChatGPT (via ngrok)

### 1. Start surfmate.io in HTTP mode

```bash
surfmate.io -http -port 8080
```

### 2. Expose with ngrok

Install [ngrok](https://ngrok.com/download) if you haven't, then in a new terminal:

```bash
ngrok http 8080
```

Copy the HTTPS URL (e.g., `https://abc123.ngrok.io`)

### 3. Create a ChatGPT Action

1. Go to [chat.openai.com](https://chat.openai.com) and open a GPT you can edit (or create one)
2. Click **Configure** → **Create new action**
3. For the schema, click "Import from URL" and enter: `YOUR_NGROK_URL/openapi.yaml`
4. Set Authentication to **None**
5. Save and test!

### 4. Use it

Ask ChatGPT things like:
- "Navigate to https://example.com and tell me what you see"
- "Go to Wikipedia and search for 'artificial intelligence'"

> **Tip:** Keep both terminal windows open. The browser appears on your machine so you can solve captchas.

---

## Usage Examples

Once configured, ask your AI:

- "Navigate to https://example.com and take a screenshot"
- "Search for 'MCP protocol' on Google and extract the results"
- "Go to my email and wait for me to log in"

---

## MCP Tools

| Tool | Description |
|------|-------------|
| `navigate(url)` | Navigate to URL, returns title and captcha status |
| `click(selector)` | Click element by CSS selector or XPath |
| `type(selector, text, submit?)` | Type into input, optionally press Enter |
| `scroll(direction?, amount?, selector?)` | Scroll page or element into view |
| `get_page_content(include_html?)` | Get page text or HTML |
| `screenshot(full_page?, selector?, quality?)` | Capture as base64 PNG |
| `extract_text(selector, multiple?)` | Extract text from elements |
| `wait_for_user(reason?, timeout?)` | Pause for captcha/login resolution |

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `BROWSER_PATH` | auto-detect | Path to browser binary |
| `BROWSER_TIMEOUT` | `30s` | Operation timeout |
| `VIEWPORT_WIDTH` | `1280` | Browser viewport width |
| `VIEWPORT_HEIGHT` | `800` | Browser viewport height |
| `HEADLESS` | `false` | Run in headless mode |

---

## Build from Source

Requires [Go 1.23+](https://go.dev/dl/)

```bash
git clone https://github.com/afalcongonzalez/surfmate.io.git
cd surfmate.io
go build -o surfmate.io .
sudo mv surfmate.io /usr/local/bin/
```

---

## Development

```bash
# Run locally
go run ./main.go

# HTTP mode
go run ./main.go -http -port 8080

# Build all platforms
make build-all

# Docker (headless)
make docker-dev
```

## License

MIT
