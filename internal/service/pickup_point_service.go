package service

import (
	"context"
	"fmt"
	"ops_service/internal/model"
	"ops_service/internal/repository"
	"time"
)

type PickupPointService struct {
	repo repository.PickupPointRepository
}

func NewPickupPointService(repo repository.PickupPointRepository) *PickupPointService {
	return &PickupPointService{
		repo: repo,
	}
}

func (s *PickupPointService) CreatePickupPoint(ctx context.Context, city model.City) (*model.PickupPoint, error) {
	if !model.IsCityAllowed(city) {
		return nil, model.ErrCityNotAllowed
	}

	point := &model.PickupPoint{
		RegisteredAt: time.Now(),
		City:         city,
	}

	if err := s.repo.Create(ctx, point); err != nil {
		return nil, fmt.Errorf("failed to create pickup point: %w", err)
	}

	return point, nil
}

func (s *PickupPointService) GetPickupPoint(ctx context.Context, id int) (*model.PickupPoint, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PickupPointService) ListPickupPoints(ctx context.Context, page, pageSize int, startDate, endDate *time.Time) ([]*model.PickupPoint, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return s.repo.List(ctx, page, pageSize, startDate, endDate)
}
