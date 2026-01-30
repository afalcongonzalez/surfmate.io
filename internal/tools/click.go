package tools

import (
	"context"
	"fmt"

	"github.com/afalcongonzalez/surfmate.io/internal/browser"
	"github.com/mark3labs/mcp-go/mcp"
)

// ClickHandler handles the click tool
func ClickHandler(mgr *browser.Manager) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		selector, ok := req.Params.Arguments["selector"].(string)
		if !ok || selector == "" {
			return mcp.NewToolResultError("selector parameter is required"), nil
		}

		if err := mgr.Click(selector); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("click failed: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Clicked element: %s", selector)), nil
	}
}

// ClickTool returns the tool definition for click
func ClickTool() mcp.Tool {
	return mcp.NewTool(
		"click",
		mcp.WithDescription("Click an element on the page by CSS selector or XPath."),
		mcp.WithString("selector",
			mcp.Required(),
			mcp.Description("CSS selector or XPath to the element to click"),
		),
	)
}
