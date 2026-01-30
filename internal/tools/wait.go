package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/afalcongonzalez/surfmate.io/internal/browser"
	"github.com/mark3labs/mcp-go/mcp"
)

// WaitForUserHandler handles the wait_for_user tool
func WaitForUserHandler(mgr *browser.Manager) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		reason := "User intervention required"
		if r, ok := req.Params.Arguments["reason"].(string); ok && r != "" {
			reason = r
		}

		timeout := 5 * time.Minute
		if t, ok := req.Params.Arguments["timeout"].(float64); ok && t > 0 {
			timeout = time.Duration(t) * time.Second
		}

		fmt.Printf("Waiting for user: %s (timeout: %v)\n", reason, timeout)

		if err := mgr.WaitForUser(timeout); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("wait failed: %v", err)), nil
		}

		hasCaptcha := mgr.HasCaptcha()
		if hasCaptcha {
			return mcp.NewToolResultText("Wait completed but captcha still detected. User may need more time."), nil
		}

		return mcp.NewToolResultText("User intervention completed successfully."), nil
	}
}

// WaitForUserTool returns the tool definition for wait_for_user
func WaitForUserTool() mcp.Tool {
	return mcp.NewTool(
		"wait_for_user",
		mcp.WithDescription("Pause execution and wait for user to complete an action (like solving a captcha or logging in). Polls until captcha elements disappear or timeout."),
		mcp.WithString("reason",
			mcp.Description("Reason for waiting (e.g., 'Solve captcha', 'Complete login')"),
		),
		mcp.WithNumber("timeout",
			mcp.Description("Maximum time to wait in seconds (default: 300)"),
		),
	)
}
