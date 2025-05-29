package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jerop/cloud-build-mcp-server/internal/tools"
	"github.com/mark3labs/mcp-go/server"
)

// StartServer initializes and starts the MCP server.
func StartServer() {
	s := server.NewMCPServer(
		"Cloud Build MCP Server",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)
	slog.Info("Adding tools and resources to the server.")
	ctx := context.Background()
	tools.Add(ctx, s)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
