package clients

import (
	"context"
	"fmt"
	"time"

	permissionpb "task-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type PermissionClient struct {
	conn   *grpc.ClientConn
	client permissionpb.PermissionServiceClient
}

func NewPermissionClient(addr string) (*PermissionClient, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &PermissionClient{
		conn:   conn,
		client: permissionpb.NewPermissionServiceClient(conn),
	}, nil
}

func (c *PermissionClient) Close() error {
	return c.conn.Close()
}

func (c *PermissionClient) Create(ctx context.Context, taskID string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	ctx, err := contextWithAuth(ctx)
	if err != nil {
		return err
	}

	_, err = c.client.CreatePermission(ctx, &permissionpb.CreatePermissionRequest{
		TaskId: taskID,
	})

	return err
}

func (c *PermissionClient) Check(ctx context.Context, taskID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	ctx, err := contextWithAuth(ctx)
	if err != nil {
		return false, err
	}

	resp, err := c.client.CheckPermission(ctx, &permissionpb.CheckPermissionRequest{
		TaskId: taskID,
	})
	if err != nil {
		return false, err
	}

	return resp.Allowed, nil
}

func contextWithAuth(ctx context.Context) (context.Context, error) {
	authHeader, ok := ctx.Value("auth_header").(string)
	if !ok || authHeader == "" {
		return nil, fmt.Errorf("missing auth header in context")
	}

	return metadata.AppendToOutgoingContext(ctx, "authorization", authHeader), nil
}
