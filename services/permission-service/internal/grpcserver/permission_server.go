package grpcserver

import (
	"context"
	"fmt"
	"os"
	"strings"

	"permission-service/internal/services"
	permissionpb "permission-service/proto"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type PermissionServer struct {
	permissionpb.UnimplementedPermissionServiceServer
	service *services.PermissionService
}

func NewPermissionServer(service *services.PermissionService) *PermissionServer {
	return &PermissionServer{
		service: service,
	}
}

func (s *PermissionServer) CreatePermission(
	ctx context.Context,
	req *permissionpb.CreatePermissionRequest,
) (*permissionpb.CreatePermissionResponse, error) {
	if req.TaskId == "" {
		return nil, status.Error(codes.InvalidArgument, "task_id required")
	}

	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.service.Create(userID, req.TaskId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &permissionpb.CreatePermissionResponse{Created: true}, nil
}

func (s *PermissionServer) CheckPermission(
	ctx context.Context,
	req *permissionpb.CheckPermissionRequest,
) (*permissionpb.CheckPermissionResponse, error) {
	if req.TaskId == "" {
		return nil, status.Error(codes.InvalidArgument, "task_id required")
	}

	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	allowed, err := s.service.Check(userID, req.TaskId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &permissionpb.CheckPermissionResponse{Allowed: allowed}, nil
}

func userIDFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization metadata")
	}

	authHeader := values[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", status.Error(codes.Unauthenticated, "invalid authorization format")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", status.Error(codes.Internal, "JWT_SECRET not set")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return "", status.Error(codes.Unauthenticated, "invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "invalid claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return "", status.Error(codes.Unauthenticated, "invalid user_id")
	}

	return userID, nil
}
