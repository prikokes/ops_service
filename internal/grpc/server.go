package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"ops_service/internal/config"
	"ops_service/internal/repository"
	"time"

	"google.golang.org/grpc"
)

// Server represents gRPC server
type Server struct {
	server *grpc.Server
	repo   repository.PickupPointRepository
}

// NewServer creates a new gRPC server
func NewServer(repo repository.PickupPointRepository) *Server {
	return &Server{
		server: grpc.NewServer(),
		repo:   repo,
	}
}

// Run starts the gRPC server
func (s *Server) Run(cfg config.GRPCConfig) error {
	addr := fmt.Sprintf(":%d", cfg.Port)
	
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	
	// Register services
	service := NewPickupPointService(s.repo)
	RegisterPickupPointServiceServer(s.server, service)
	
	log.Printf("Starting gRPC server on port %d", cfg.Port)
	
	return s.server.Serve(lis)
}

// Shutdown gracefully stops the gRPC server
func (s *Server) Shutdown() {
	s.server.GracefulStop()
}

// PickupPointGRPCService implements gRPC service for pickup points
type PickupPointGRPCService struct {
	UnimplementedPickupPointServiceServer
	repo repository.PickupPointRepository
}

// NewPickupPointService creates a new PickupPointGRPCService
func NewPickupPointService(repo repository.PickupPointRepository) *PickupPointGRPCService {
	return &PickupPointGRPCService{
		repo: repo,
	}
}

// GetAllPickupPoints returns all pickup points
func (s *PickupPointGRPCService) GetAllPickupPoints(ctx context.Context, req *GetAllPickupPointsRequest) (*GetAllPickupPointsResponse, error) {
	// Use repository to get all pickup points
	points, _, err := s.repo.List(ctx, 1, 1000, nil, nil) // Get all pickup points
	if err != nil {
		return nil, fmt.Errorf("failed to get pickup points: %w", err)
	}
	
	// Convert to gRPC model
	var protoPoints []*PickupPoint
	for _, point := range points {
		protoPoints = append(protoPoints, &PickupPoint{
			Id:           int32(point.ID),
			RegisteredAt: point.RegisteredAt.Format(time.RFC3339),
			City:         string(point.City),
		})
	}
	
	return &GetAllPickupPointsResponse{
		PickupPoints: protoPoints,
	}, nil
} 