package server

import (
	"context"
	"fmt"
	"time"

	"github.com/example/gnxi-mcp/internal/gnoi"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetDeviceTimeArgs defines the input parameters for the get_device_time tool
type GetDeviceTimeArgs struct {
	Host     string `json:"host" jsonschema:"Device hostname or IP address"`
	Port     int    `json:"port" jsonschema:"gRPC port number"`
	Insecure bool   `json:"insecure" jsonschema:"Use insecure connection (default: true)"`
}

// GetDeviceTime implements the MCP tool to retrieve device time via gNOI
func GetDeviceTime(ctx context.Context, args GetDeviceTimeArgs) (
	*mcp.CallToolResult,
	any,
	error,
) {
	// Default to insecure if not specified
	if args.Port == 0 {
		args.Port = 9339 // Standard gNOI port
	}

	address := fmt.Sprintf("%s:%d", args.Host, args.Port)
	
	// Create gNOI client
	client, err := gnoi.NewClient(address, args.Insecure)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to connect to device: %v", err)},
			},
		}, nil, err
	}
	defer client.Close()

	// Get device time
	timeResp, err := client.GetTime(ctx)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Failed to get device time: %v", err)},
			},
		}, nil, err
	}

	// Convert timestamp to human readable format
	t := time.Unix(0, int64(timeResp.Time))
	formatted := t.Format(time.RFC3339)

	output := map[string]any{
		"timestamp":      timeResp.Time,
		"formatted_time": formatted,
		"device":        address,
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Device %s time: %s (timestamp: %d)", 
					address, formatted, timeResp.Time),
			},
		},
	}, output, nil
}