package tools

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/afalcongonzalez/surfmate.io/internal/browser"
	"github.com/mark3labs/mcp-go/mcp"
)

// ScreenshotHandler handles the screenshot tool
func ScreenshotHandler(mgr *browser.Manager) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fullPage := false
		if f, ok := req.Params.Arguments["full_page"].(bool); ok {
			fullPage = f
		}

		selector := ""
		if s, ok := req.Params.Arguments["selector"].(string); ok {
			selector = s
		}

		quality := 80
		if q, ok := req.Params.Arguments["quality"].(float64); ok {
			quality = int(q)
		}

		data, err := mgr.Screenshot(fullPage, selector, quality)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("screenshot failed: %v", err)), nil
		}

		encoded := base64.StdEncoding.EncodeToString(data)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewImageContent(encoded, "image/png"),
			},
		}, nil
	}
}

// ScreenshotTool returns the tool definition for screenshot
func ScreenshotTool() mcp.Tool {
	return mcp.NewTool(
		"screenshot",
		mcp.WithDescription("Capture a screenshot of the current page. Returns a base64 encoded PNG image."),
		mcp.WithBoolean("full_page",
			mcp.Description("Capture the full scrollable page (default: false, captures viewport only)"),
		),
		mcp.WithString("selector",
			mcp.Description("If provided, capture only this element"),
		),
		mcp.WithNumber("quality",
			mcp.Description("Image quality 1-100 (default: 80)"),
		),
	)
}
