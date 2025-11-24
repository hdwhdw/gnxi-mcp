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
	client := mcp.NewClient(&mcp.Implementation{Name: "sonic-test-client", Version: "v1.0.0"}, nil)

	// Connect to our GNXI MCP server
	cmd := exec.Command("./gnxi-mcp-server")
	transport := &mcp.CommandTransport{Command: cmd}
	
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer session.Close()

	fmt.Println("=== Testing with SONiC device 10.250.0.101:8080 ===")
	
	// Call get_device_time with the correct SONiC device endpoint
	args := map[string]any{
		"host":     "10.250.0.101",
		"port":     8080,  // Use the correct port
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
	if result.IsError {
		fmt.Printf("❌ Error calling device:\n")
	} else {
		fmt.Printf("✅ Successfully got device time:\n")
	}
	
	for _, content := range result.Content {
		if textContent, ok := content.(*mcp.TextContent); ok {
			fmt.Printf("%s\n", textContent.Text)
		}
	}

	// Also print the structured output if available
	fmt.Printf("\n=== Raw Result ===\n")
	fmt.Printf("IsError: %t\n", result.IsError)
	if len(result.Content) > 0 {
		fmt.Printf("Content Count: %d\n", len(result.Content))
	}
}