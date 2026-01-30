package tools

import (
	"github.com/afalcongonzalez/surfmate.io/internal/browser"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAll registers all browser tools with the MCP server
func RegisterAll(s *server.MCPServer, mgr *browser.Manager) {
	// Navigation
	s.AddTool(NavigateTool(), NavigateHandler(mgr))

	// Interaction
	s.AddTool(ClickTool(), ClickHandler(mgr))
	s.AddTool(TypeTool(), TypeHandler(mgr))
	s.AddTool(ScrollTool(), ScrollHandler(mgr))

	// Content
	s.AddTool(GetPageContentTool(), GetPageContentHandler(mgr))
	s.AddTool(ExtractTextTool(), ExtractTextHandler(mgr))
	s.AddTool(ScreenshotTool(), ScreenshotHandler(mgr))

	// User intervention
	s.AddTool(WaitForUserTool(), WaitForUserHandler(mgr))
}
