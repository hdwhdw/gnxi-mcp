package main

import (
	"context"
	"log"

	"github.com/example/gnxi-mcp/internal/server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Create MCP server
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "gnxi-mcp-server",
		Version: "1.0.0",
	}, nil)

	// Add the get_device_time tool
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_device_time",
		Description: "Get system time from a network device using gNOI system service",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args server.GetDeviceTimeArgs) (*mcp.CallToolResult, any, error) {
		return server.GetDeviceTime(ctx, args)
	})

	// Run server on stdio transport
	if err := mcpServer.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}