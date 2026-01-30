package tools

import (
	"context"
	"fmt"

	"github.com/afalcongonzalez/surfmate.io/internal/browser"
	"github.com/mark3labs/mcp-go/mcp"
)

// TypeHandler handles the type tool
func TypeHandler(mgr *browser.Manager) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		selector, ok := req.Params.Arguments["selector"].(string)
		if !ok || selector == "" {
			return mcp.NewToolResultError("selector parameter is required"), nil
		}

		text, ok := req.Params.Arguments["text"].(string)
		if !ok {
			return mcp.NewToolResultError("text parameter is required"), nil
		}

		submit := false
		if s, ok := req.Params.Arguments["submit"].(bool); ok {
			submit = s
		}

		if err := mgr.Type(selector, text, submit); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("type failed: %v", err)), nil
		}

		result := fmt.Sprintf("Typed into element: %s", selector)
		if submit {
			result += " (submitted)"
		}
		return mcp.NewToolResultText(result), nil
	}
}

// TypeTool returns the tool definition for type
func TypeTool() mcp.Tool {
	return mcp.NewTool(
		"type",
		mcp.WithDescription("Type text into an input element. Optionally press Enter to submit."),
		mcp.WithString("selector",
			mcp.Required(),
			mcp.Description("CSS selector or XPath to the input element"),
		),
		mcp.WithString("text",
			mcp.Required(),
			mcp.Description("Text to type into the element"),
		),
		mcp.WithBoolean("submit",
			mcp.Description("Press Enter after typing to submit (default: false)"),
		),
	)
}
