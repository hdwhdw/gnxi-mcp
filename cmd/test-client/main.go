package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()

	// Create client
	client := mcp.NewClient(&mcp.Implementation{Name: "gnxi-test-client", Version: "v1.0.0"}, nil)

	// Connect to our GNXI MCP server via command transport
	cmd := exec.Command("./gnxi-mcp-server")
	transport := &mcp.CommandTransport{Command: cmd}
	
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer session.Close()

	// List available tools
	fmt.Println("=== Available Tools ===")
	for tool, err := range session.Tools(ctx, nil) {
		if err != nil {
			log.Fatalf("Failed to list tools: %v", err)
		}
		fmt.Printf("Tool: %s - %s\n", tool.Name, tool.Description)
	}

	// Test the get_device_time tool with mock parameters
	fmt.Println("\n=== Testing get_device_time tool ===")
	args := map[string]any{
		"host":     "localhost", // This will fail but shows the tool working
		"port":     9339,
		"insecure": true,
	}

	argsJSON, err := json.Marshal(args)
	if err != nil {
		log.Fatalf("Failed to marshal args: %v", err)
	}

	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "get_device_time",
		Arguments: json.RawMessage(argsJSON),
	})
	if err != nil {
		log.Fatalf("Tool call failed: %v", err)
	}

	// Print the result
	fmt.Printf("Tool call result:\n")
	fmt.Printf("IsError: %t\n", result.IsError)
	
	for i, content := range result.Content {
		switch c := content.(type) {
		case *mcp.TextContent:
			fmt.Printf("Content[%d]: %s\n", i, c.Text)
		default:
			fmt.Printf("Content[%d]: %+v\n", i, content)
		}
	}

	if result.Meta != nil {
		fmt.Printf("Meta: %+v\n", result.Meta)
	}
}