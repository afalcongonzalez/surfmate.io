package tools

import (
	"context"
	"fmt"

	"github.com/afalcongonzalez/surfmate.io/internal/browser"
	"github.com/mark3labs/mcp-go/mcp"
)

// GetPageContentHandler handles the get_page_content tool
func GetPageContentHandler(mgr *browser.Manager) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		includeHTML := false
		if h, ok := req.Params.Arguments["include_html"].(bool); ok {
			includeHTML = h
		}

		content, err := mgr.GetContent(includeHTML)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get content: %v", err)), nil
		}

		return mcp.NewToolResultText(content), nil
	}
}

// GetPageContentTool returns the tool definition for get_page_content
func GetPageContentTool() mcp.Tool {
	return mcp.NewTool(
		"get_page_content",
		mcp.WithDescription("Get the text content or HTML of the current page."),
		mcp.WithBoolean("include_html",
			mcp.Description("Return full HTML instead of text content (default: false)"),
		),
	)
}
