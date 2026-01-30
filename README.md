# Surfmate.io

Let AI browse the web for you. A browser window opens on your screen so you can handle captchas and logins while the AI does everything else.

Works with **Claude Desktop** and **ChatGPT**.

## Install

### macOS (Apple Silicon)

```bash
curl -fsSL https://github.com/afalcongonzalez/surfmate.io/releases/latest/download/surfmate.io-darwin-arm64 -o ~/surfmate.io && chmod +x ~/surfmate.io && sudo mv ~/surfmate.io /usr/local/bin/surfmate.io
```

### macOS (Intel)

```bash
curl -fsSL https://github.com/afalcongonzalez/surfmate.io/releases/latest/download/surfmate.io-darwin-amd64 -o ~/surfmate.io && chmod +x ~/surfmate.io && sudo mv ~/surfmate.io /usr/local/bin/surfmate.io
```

### Linux

```bash
curl -fsSL https://github.com/afalcongonzalez/surfmate.io/releases/latest/download/surfmate.io-linux-amd64 -o ~/surfmate.io && chmod +x ~/surfmate.io && sudo mv ~/surfmate.io /usr/local/bin/surfmate.io
```

### Windows (PowerShell as Admin)

```powershell
irm https://github.com/afalcongonzalez/surfmate.io/releases/latest/download/surfmate.io-windows-amd64.exe -OutFile surfmate.io.exe; mkdir -Force "$env:USERPROFILE\bin"; mv surfmate.io.exe "$env:USERPROFILE\bin\"
```

> **Note:** You'll need your password for the install. On macOS, if you see a security warning, run `xattr -d com.apple.quarantine /usr/local/bin/surfmate.io`. On Windows, click "More info" → "Run anyway" if SmartScreen appears.

---

## Setup for Claude Desktop

1. Open Claude Desktop → **Settings** → **Developer** → **Edit Config**

2. Paste this:

```json
{
  "mcpServers": {
    "surfmate": {
      "command": "surfmate.io"
    }
  }
}
```

3. Restart Claude Desktop

4. Try asking: *"Open the browser and go to news.ycombinator.com"*

---

## Setup for ChatGPT

1. **Start the server** (keep this terminal open):
   ```bash
   surfmate.io -http -port 8080
   ```

2. **Make it accessible** with [ngrok](https://ngrok.com/download) (new terminal):
   ```bash
   ngrok http 8080
   ```
   Copy the `https://` URL it shows.

3. **Add to ChatGPT:**
   - Go to ChatGPT → Create or edit a GPT
   - Click **Configure** → **Add actions** → **Create new action**
   - Click **Import from URL** and paste: `YOUR_NGROK_URL/openapi.yaml`
   - Save!

4. Try asking: *"Open the browser and navigate to example.com"*

> **Keep both terminals open** while using ChatGPT. The browser appears on your computer so you can handle captchas.

---

## How it works

1. You ask your AI to open the browser
2. A browser window opens on your screen
3. The AI controls the browser (clicking, typing, reading)
4. When there's a login or captcha, the AI waits for you to handle it
5. You solve it in the browser window, then the AI continues

---

## Example chat sessions

### Basic browsing

```
You: Open the browser and go to Hacker News

AI: [calls open_browser]
    Browser launched successfully.

    [calls navigate with url="https://news.ycombinator.com"]
    Navigated to: https://news.ycombinator.com/
    Title: Hacker News
    Captcha detected: false

    I've opened Hacker News. I can see the front page with today's top stories.
```

### Searching and extracting content

```
You: Open a browser, search Wikipedia for "artificial intelligence" and
     give me the first paragraph

AI: [calls open_browser]
    Browser launched successfully.

    [calls navigate with url="https://wikipedia.org"]
    Navigated to: https://www.wikipedia.org/

    [calls type with selector="input[name='search']", text="artificial intelligence", submit=true]
    Typed text and submitted.

    [calls extract_text with selector=".mw-parser-output p"]
    "Artificial intelligence (AI), in its broadest sense, is intelligence
    exhibited by machines..."

    Here's the first paragraph from the Wikipedia article on AI: ...
```

### Handling logins

```
You: Go to my email and wait for me to log in

AI: [calls open_browser]
    Browser launched successfully.

    [calls navigate with url="https://gmail.com"]
    Navigated to: https://accounts.google.com/...
    Captcha detected: false

    I've opened Gmail. Please log in to your account in the browser window.
    Let me know when you're done.

You: Ok I'm logged in

AI: [calls get_page_content]
    "Primary inbox - 3 unread messages..."

    I can see your inbox now. You have 3 unread messages. Would you like
    me to summarize them?
```

### Taking screenshots

```
You: Open a browser, go to example.com and take a screenshot

AI: [calls open_browser]
    Browser launched successfully.

    [calls navigate with url="https://example.com"]
    Navigated to: https://example.com/
    Title: Example Domain

    [calls screenshot]
    [returns base64 image]

    Here's a screenshot of example.com. It shows the default example domain
    page with a heading and a link to more information.
```

### Solving captchas

```
You: Search Google for "best pizza near me"

AI: [calls open_browser]
    Browser launched successfully.

    [calls navigate with url="https://google.com"]
    Navigated to: https://google.com/
    Captcha detected: true

    I see there's a captcha on the page. Please solve it in the browser
    window, and I'll continue once you're done.

    [calls wait_for_user with reason="Solve captcha"]
    User completed action.

    [calls type with selector="input[name='q']", text="best pizza near me", submit=true]
    Typed text and submitted.

    I've searched for "best pizza near me". Here are the results...
```

---

## Troubleshooting

**macOS says the app can't be opened:**
The install command handles this, but if needed:
```bash
xattr -d com.apple.quarantine /usr/local/bin/surfmate.io
```

**Windows SmartScreen warning:**
Click "More info" → "Run anyway"

**Browser doesn't open:**
Make sure you have Chrome, Edge, or Brave installed.

---

## License

MIT
