package service

import (
	"context"
	"fmt"
	"ops_service/internal/model"
	"ops_service/internal/repository"
	"time"
)

type ReceiptService struct {
	repo        repository.ReceiptRepository
	pickupRepo  repository.PickupPointRepository
	productRepo repository.ProductRepository
}

func NewReceiptService(
	repo repository.ReceiptRepository,
	pickupRepo repository.PickupPointRepository,
	productRepo repository.ProductRepository,
) *ReceiptService {
	return &ReceiptService{
		repo:        repo,
		pickupRepo:  pickupRepo,
		productRepo: productRepo,
	}
}

func (s *ReceiptService) CreateReceipt(ctx context.Context, pickupPointID int) (*model.Receipt, error) {
	_, err := s.pickupRepo.GetByID(ctx, pickupPointID)
	if err != nil {
		return nil, fmt.Errorf("failed to find pickup point: %w", err)
	}

	receipt := &model.Receipt{
		CreatedAt:     time.Now(),
		PickupPointID: pickupPointID,
		Status:        model.ReceiptStatusInProgress,
	}

	if err := s.repo.Create(ctx, receipt); err != nil {
		return nil, fmt.Errorf("failed to create receipt: %w", err)
	}

	return receipt, nil
}

func (s *ReceiptService) GetReceipt(ctx context.Context, id int) (*model.Receipt, error) {
	receipt, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	products, err := s.productRepo.GetByReceiptID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get receipt products: %w", err)
	}

	receipt.Products = products

	return receipt, nil
}

func (s *ReceiptService) GetLastOpenReceipt(ctx context.Context, pickupPointID int) (*model.Receipt, error) {
	return s.repo.GetLastOpenByPickupPointID(ctx, pickupPointID)
}

func (s *ReceiptService) CloseReceipt(ctx context.Context, id int) error {
	return s.repo.Close(ctx, id)
}

func (s *ReceiptService) AddProduct(ctx context.Context, pickupPointID int, productType model.ProductType) (*model.Product, error) {
	receipt, err := s.repo.GetLastOpenByPickupPointID(ctx, pickupPointID)
	if err != nil {
		return nil, err
	}

	product := &model.Product{
		ReceivedAt: time.Now(),
		Type:       productType,
		ReceiptID:  receipt.ID,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to add product: %w", err)
	}

	return product, nil
}

func (s *ReceiptService) DeleteLastProduct(ctx context.Context, pickupPointID int) error {
	receipt, err := s.repo.GetLastOpenByPickupPointID(ctx, pickupPointID)
	if err != nil {
		return err
	}

	return s.productRepo.DeleteLast(ctx, receipt.ID)
}

func (s *ReceiptService) GetReceiptsByDateRange(ctx context.Context, pickupPointID int, startDate, endDate time.Time) ([]*model.Receipt, error) {
	receipts, err := s.repo.GetByDateRange(ctx, pickupPointID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	for _, receipt := range receipts {
		products, err := s.productRepo.GetByReceiptID(ctx, receipt.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get receipt products: %w", err)
		}
		receipt.Products = products
	}

	return receipts, nil
}
