package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/afalcongonzalez/surfmate.io/internal/browser"
	"github.com/mark3labs/mcp-go/mcp"
)

// ExtractTextHandler handles the extract_text tool
func ExtractTextHandler(mgr *browser.Manager) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		selector, ok := req.Params.Arguments["selector"].(string)
		if !ok || selector == "" {
			return mcp.NewToolResultError("selector parameter is required"), nil
		}

		multiple := false
		if m, ok := req.Params.Arguments["multiple"].(bool); ok {
			multiple = m
		}

		texts, err := mgr.ExtractText(selector, multiple)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("extract failed: %v", err)), nil
		}

		if len(texts) == 0 {
			return mcp.NewToolResultText("No text found for selector: " + selector), nil
		}

		if multiple {
			result := fmt.Sprintf("Found %d elements:\n%s", len(texts), strings.Join(texts, "\n---\n"))
			return mcp.NewToolResultText(result), nil
		}

		return mcp.NewToolResultText(texts[0]), nil
	}
}

// ExtractTextTool returns the tool definition for extract_text
func ExtractTextTool() mcp.Tool {
	return mcp.NewTool(
		"extract_text",
		mcp.WithDescription("Extract text content from elements matching a selector."),
		mcp.WithString("selector",
			mcp.Required(),
			mcp.Description("CSS selector or XPath to the element(s)"),
		),
		mcp.WithBoolean("multiple",
			mcp.Description("Extract text from all matching elements (default: false, returns first match only)"),
		),
	)
}
