package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/afalcongonzalez/surfmate.io/internal/browser"
	"github.com/afalcongonzalez/surfmate.io/internal/config"
	httpserver "github.com/afalcongonzalez/surfmate.io/internal/server"
	"github.com/afalcongonzalez/surfmate.io/internal/tools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Parse flags
	httpMode := flag.Bool("http", false, "Run as HTTP server instead of MCP")
	port := flag.Int("port", 8080, "HTTP server port (only used with -http)")
	flag.Parse()

	cfg := config.Load()

	// Create browser manager
	mgr := browser.GetManager(cfg)

	// Handle shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		mgr.Close()
		os.Exit(0)
	}()

	if *httpMode {
		// Run HTTP server for ChatGPT Actions
		srv := httpserver.NewHTTPServer(mgr, *port)
		if err := srv.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "HTTP server error: %v\n", err)
			mgr.Close()
			os.Exit(1)
		}
	} else {
		// Run MCP server (default)
		s := server.NewMCPServer(
			"surfmate.io",
			"1.0.0",
			server.WithToolCapabilities(true),
		)

		// Register all tools
		tools.RegisterAll(s, mgr)

		// Start stdio transport
		if err := server.ServeStdio(s); err != nil {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			mgr.Close()
			os.Exit(1)
		}
	}
}
