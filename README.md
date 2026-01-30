# Surfmate.io

Browser automation MCP server for AI. Opens a visible browser so you can handle captchas and logins while AI controls navigation.

## Quick Install

### macOS / Linux

```bash
git clone https://github.com/afalcongonzalez/surfmate.io.git
cd surfmate.io
go build -o surfmate.io .
sudo mv surfmate.io /usr/local/bin/
```

### Windows (PowerShell)

```powershell
git clone https://github.com/afalcongonzalez/surfmate.io.git
cd surfmate.io
go build -o surfmate.io.exe .
# Move surfmate.io.exe to a folder in your PATH
```

### One-liner

**macOS / Linux:**
```bash
git clone https://github.com/afalcongonzalez/surfmate.io.git && cd surfmate.io && go build -o surfmate.io . && sudo mv surfmate.io /usr/local/bin/
```

**Windows (PowerShell):**
```powershell
git clone https://github.com/afalcongonzalez/surfmate.io.git; cd surfmate.io; go build -o surfmate.io.exe .
```

## Requirements

- [Go 1.23+](https://go.dev/dl/)
- Chromium-based browser (Chrome, Edge, or Brave)

## Claude Desktop Setup

Add to your config file:

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

## Usage

Ask Claude things like:

- "Navigate to https://example.com and take a screenshot"
- "Search for 'MCP protocol' on Google and extract the results"
- "Go to my email and wait for me to log in"

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

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `BROWSER_PATH` | auto-detect | Path to browser binary |
| `BROWSER_TIMEOUT` | `30s` | Operation timeout |
| `VIEWPORT_WIDTH` | `1280` | Browser viewport width |
| `VIEWPORT_HEIGHT` | `800` | Browser viewport height |
| `HEADLESS` | `false` | Run in headless mode |

## Development

```bash
# Run locally
go run ./main.go

# HTTP mode (for testing)
go run ./main.go -http -port 8080

# Build
make build

# Build all platforms
make build-all

# Docker (headless HTTP mode)
make docker-dev
```

## License

MIT
