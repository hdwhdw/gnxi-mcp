package gnoi

import (
	"context"
	"crypto/tls"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	systempb "github.com/openconfig/gnoi/system"
)

// Client wraps gNOI system service calls
type Client struct {
	conn *grpc.ClientConn
	sysClient systempb.SystemClient
}

// NewClient creates a new gNOI client
func NewClient(address string, useInsecure bool) (*Client, error) {
	var opts []grpc.DialOption
	
	if useInsecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})))
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", address, err)
	}

	return &Client{
		conn:      conn,
		sysClient: systempb.NewSystemClient(conn),
	}, nil
}

// GetTime retrieves the system time from the device
func (c *Client) GetTime(ctx context.Context) (*systempb.TimeResponse, error) {
	req := &systempb.TimeRequest{}
	return c.sysClient.Time(ctx, req)
}

// Close closes the gRPC connection
func (c *Client) Close() error {
	return c.conn.Close()
}