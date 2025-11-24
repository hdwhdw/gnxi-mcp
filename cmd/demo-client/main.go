package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: demo-client <device_ip> <port>")
		fmt.Println("Example: demo-client 192.168.1.100 9339")
		fmt.Println("Example: demo-client 10.250.0.101 8080")
		os.Exit(1)
	}

	deviceIP := os.Args[1]
	portStr := os.Args[2]
	
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}
	ctx := context.Background()

	// Create client
	client := mcp.NewClient(&mcp.Implementation{Name: "gnxi-demo-client", Version: "v1.0.0"}, nil)

	// Connect to our GNXI MCP server
	cmd := exec.Command("./gnxi-mcp-server")
	transport := &mcp.CommandTransport{Command: cmd}
	
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer session.Close()

	fmt.Printf("=== Getting time from device: %s:%d ===\n", deviceIP, port)
	
	// Call get_device_time with provided device IP and port
	args := map[string]any{
		"host":     deviceIP,
		"port":     port,
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
}