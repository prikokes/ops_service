package repository

import (
	"context"
	"database/sql"
	"ops_service/internal/model"
	"time"
)

// UserRepository defines methods to work with users
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
}

// PickupPointRepository defines methods to work with pickup points
type PickupPointRepository interface {
	Create(ctx context.Context, point *model.PickupPoint) error
	GetByID(ctx context.Context, id int) (*model.PickupPoint, error)
	List(ctx context.Context, page, pageSize int, startDate, endDate *time.Time) ([]*model.PickupPoint, int, error)
}

// ReceiptRepository defines methods to work with receipts
type ReceiptRepository interface {
	Create(ctx context.Context, receipt *model.Receipt) error
	GetByID(ctx context.Context, id int) (*model.Receipt, error)
	GetLastOpenByPickupPointID(ctx context.Context, pickupPointID int) (*model.Receipt, error)
	Close(ctx context.Context, id int) error
	GetByDateRange(ctx context.Context, pickupPointID int, startDate, endDate time.Time) ([]*model.Receipt, error)
}

// ProductRepository defines methods to work with products
type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	GetByReceiptID(ctx context.Context, receiptID int) ([]*model.Product, error)
	DeleteLast(ctx context.Context, receiptID int) error
	GetLastByReceiptID(ctx context.Context, receiptID int) (*model.Product, error)
}

// Repository defines the interface for all repositories
type Repository struct {
	User        UserRepository
	PickupPoint PickupPointRepository
	Receipt     ReceiptRepository
	Product     ProductRepository
}

// NewRepository creates a new repository instance
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User:        NewUserRepo(db),
		PickupPoint: NewPickupPointRepo(db),
		Receipt:     NewReceiptRepo(db),
		Product:     NewProductRepo(db),
	}
} 