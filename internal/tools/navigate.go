package tools

import (
	"context"
	"fmt"

	"github.com/afalcongonzalez/surfmate.io/internal/browser"
	"github.com/mark3labs/mcp-go/mcp"
)

// NavigateHandler handles the navigate tool
func NavigateHandler(mgr *browser.Manager) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url, ok := req.Params.Arguments["url"].(string)
		if !ok || url == "" {
			return mcp.NewToolResultError("url parameter is required"), nil
		}

		if err := mgr.Navigate(url); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("navigation failed: %v", err)), nil
		}

		if err := mgr.WaitLoad(); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("page load failed: %v", err)), nil
		}

		title, _ := mgr.GetTitle()
		currentURL, _ := mgr.GetURL()
		hasCaptcha := mgr.HasCaptcha()

		result := fmt.Sprintf("Navigated to: %s\nTitle: %s\nCaptcha detected: %v", currentURL, title, hasCaptcha)
		return mcp.NewToolResultText(result), nil
	}
}

// NavigateTool returns the tool definition for navigate
func NavigateTool() mcp.Tool {
	return mcp.NewTool(
		"navigate",
		mcp.WithDescription("Navigate to a URL. Returns the page title and whether a captcha was detected."),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("The URL to navigate to"),
		),
	)
}
