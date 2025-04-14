package grpc

import (
	"context"
	"google.golang.org/grpc"
)

// Temporary stubs for gRPC classes until we generate them from proto files

// UnimplementedPickupPointServiceServer is a temporary stub
type UnimplementedPickupPointServiceServer struct{}

// PickupPointServiceServer is a temporary stub interface
type PickupPointServiceServer interface {
	GetAllPickupPoints(context.Context, *GetAllPickupPointsRequest) (*GetAllPickupPointsResponse, error)
}

// RegisterPickupPointServiceServer is a temporary stub function
func RegisterPickupPointServiceServer(s *grpc.Server, srv PickupPointServiceServer) {}

// GetAllPickupPointsRequest is a temporary stub
type GetAllPickupPointsRequest struct{}

// GetAllPickupPointsResponse is a temporary stub
type GetAllPickupPointsResponse struct {
	PickupPoints []*PickupPoint `json:"pickup_points"`
}

// PickupPoint is a temporary stub
type PickupPoint struct {
	Id           int32  `json:"id"`
	RegisteredAt string `json:"registered_at"`
	City         string `json:"city"`
} 