package tools

import (
	"context"
	"fmt"

	"github.com/afalcongonzalez/surfmate.io/internal/browser"
	"github.com/mark3labs/mcp-go/mcp"
)

// OpenBrowserHandler handles the open_browser tool
func OpenBrowserHandler(mgr *browser.Manager) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if mgr.IsLaunched() {
			return mcp.NewToolResultText("Browser is already open."), nil
		}

		if err := mgr.Launch(); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to launch browser: %v", err)), nil
		}

		return mcp.NewToolResultText("Browser launched successfully."), nil
	}
}

// OpenBrowserTool returns the tool definition for open_browser
func OpenBrowserTool() mcp.Tool {
	return mcp.NewTool(
		"open_browser",
		mcp.WithDescription("Open the browser. Must be called before using any other browser tools."),
	)
}
