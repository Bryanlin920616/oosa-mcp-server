package oosa

import (
	"github.com/mark3labs/mcp-go/server"
)

func NewServer(version string) *server.MCPServer {
	// Create a new MCP server
	s := server.NewMCPServer(
		"oosa-mcp-server",
		version,
		server.WithResourceCapabilities(true, true),
		server.WithLogging())

	// Add resources
	// TODO: Add resources

	// Add tools
	// TODO: Add tools

	// Add prompts
	// TODO: Add prompts

	return s
}
