package tools

import (
	"context"
	"fmt"

	"github.com/afalcongonzalez/surfmate.io/internal/browser"
	"github.com/mark3labs/mcp-go/mcp"
)

// ScrollHandler handles the scroll tool
func ScrollHandler(mgr *browser.Manager) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		direction := "down"
		if d, ok := req.Params.Arguments["direction"].(string); ok && d != "" {
			direction = d
		}

		amount := 300
		if a, ok := req.Params.Arguments["amount"].(float64); ok {
			amount = int(a)
		}

		selector := ""
		if s, ok := req.Params.Arguments["selector"].(string); ok {
			selector = s
		}

		if err := mgr.Scroll(direction, amount, selector); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("scroll failed: %v", err)), nil
		}

		if selector != "" {
			return mcp.NewToolResultText(fmt.Sprintf("Scrolled element into view: %s", selector)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Scrolled %s by %d pixels", direction, amount)), nil
	}
}

// ScrollTool returns the tool definition for scroll
func ScrollTool() mcp.Tool {
	return mcp.NewTool(
		"scroll",
		mcp.WithDescription("Scroll the page or scroll an element into view."),
		mcp.WithString("direction",
			mcp.Description("Direction to scroll: up, down, left, right (default: down)"),
		),
		mcp.WithNumber("amount",
			mcp.Description("Amount in pixels to scroll (default: 300)"),
		),
		mcp.WithString("selector",
			mcp.Description("If provided, scroll this element into view instead of scrolling the page"),
		),
	)
}
