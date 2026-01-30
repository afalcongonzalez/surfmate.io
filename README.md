# Surfmate.io

Browser automation MCP server that enables AI models to browse the web. Opens a visible browser for human intervention (captchas, logins) while the AI controls navigation.

## Installation

Download the latest release for your platform from the [releases page](https://github.com/afalcongonzalez/surfmate.io/releases) or build from source:

```bash
go build -o surfmate.io ./main.go
```

## Configuration

Add to your Claude Desktop config:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "surfmate.io": {
      "command": "surfmate.io"
    }
  }
}
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `BROWSER_PATH` | auto-detect | Path to browser binary |
| `BROWSER_TIMEOUT` | `30s` | Operation timeout |
| `VIEWPORT_WIDTH` | `1280` | Browser viewport width |
| `VIEWPORT_HEIGHT` | `800` | Browser viewport height |

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

## How It Works

1. Surfmate.io launches a visible browser window using your system's Chrome, Edge, or Brave
2. AI models send commands via MCP to navigate, click, type, and extract content
3. When captchas or logins are detected, the AI can call `wait_for_user()` to pause
4. You solve the captcha or complete the login in the visible browser
5. AI resumes control after you're done

## Development

```bash
# Run locally
go run ./main.go

# Build
make build

# Build all platforms
make build-all

# Test
make test
```

## License

MIT
