package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ops_service/internal/model"
)

// ProductRepo is a PostgreSQL implementation of ProductRepository
type ProductRepo struct {
	db *sql.DB
}

// NewProductRepo creates a new ProductRepo
func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

// Create creates a new product in the database
func (r *ProductRepo) Create(ctx context.Context, product *model.Product) error {
	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if receipt exists and is open
	var status model.ReceiptStatus
	checkQuery := `
		SELECT status
		FROM receipts
		WHERE id = $1
	`
	err = tx.QueryRowContext(ctx, checkQuery, product.ReceiptID).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("receipt not found")
		}
		return fmt.Errorf("failed to check receipt status: %w", err)
	}

	if status != model.ReceiptStatusInProgress {
		return model.ErrReceiptAlreadyClosed
	}

	// Get the max order for this receipt
	var maxOrder sql.NullInt32
	orderQuery := `
		SELECT MAX(addition_order)
		FROM products
		WHERE receipt_id = $1
	`
	err = tx.QueryRowContext(ctx, orderQuery, product.ReceiptID).Scan(&maxOrder)
	if err != nil {
		return fmt.Errorf("failed to get max order: %w", err)
	}

	// Set the addition order
	if maxOrder.Valid {
		product.AdditionOrder = int(maxOrder.Int32) + 1
	} else {
		product.AdditionOrder = 1
	}

	// Create the product
	query := `
		INSERT INTO products (received_at, type, receipt_id, addition_order)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err = tx.QueryRowContext(
		ctx,
		query,
		product.ReceivedAt,
		product.Type,
		product.ReceiptID,
		product.AdditionOrder,
	).Scan(&product.ID)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetByReceiptID returns all products for a receipt, ordered by addition order
func (r *ProductRepo) GetByReceiptID(ctx context.Context, receiptID int) ([]*model.Product, error) {
	query := `
		SELECT id, received_at, type, receipt_id, addition_order
		FROM products
		WHERE receipt_id = $1
		ORDER BY addition_order
	`

	rows, err := r.db.QueryContext(ctx, query, receiptID)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		product := &model.Product{}
		if err := rows.Scan(&product.ID, &product.ReceivedAt, &product.Type, &product.ReceiptID, &product.AdditionOrder); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over products: %w", err)
	}

	return products, nil
}

// DeleteLast deletes the last added product in a receipt
func (r *ProductRepo) DeleteLast(ctx context.Context, receiptID int) error {
	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if receipt exists and is open
	var status model.ReceiptStatus
	checkQuery := `
		SELECT status
		FROM receipts
		WHERE id = $1
	`
	err = tx.QueryRowContext(ctx, checkQuery, receiptID).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("receipt not found")
		}
		return fmt.Errorf("failed to check receipt status: %w", err)
	}

	if status != model.ReceiptStatusInProgress {
		return model.ErrReceiptAlreadyClosed
	}

	// Get the max order product for this receipt
	maxOrderQuery := `
		SELECT MAX(addition_order)
		FROM products
		WHERE receipt_id = $1
	`
	var maxOrder sql.NullInt32
	err = tx.QueryRowContext(ctx, maxOrderQuery, receiptID).Scan(&maxOrder)
	if err != nil {
		return fmt.Errorf("failed to get max order: %w", err)
	}

	if !maxOrder.Valid {
		return errors.New("no products to delete")
	}

	// Delete the product
	deleteQuery := `
		DELETE FROM products
		WHERE receipt_id = $1 AND addition_order = $2
	`
	result, err := tx.ExecContext(ctx, deleteQuery, receiptID, maxOrder.Int32)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if affected == 0 {
		return errors.New("product not found")
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetLastByReceiptID returns the last added product in a receipt
func (r *ProductRepo) GetLastByReceiptID(ctx context.Context, receiptID int) (*model.Product, error) {
	query := `
		SELECT id, received_at, type, receipt_id, addition_order
		FROM products
		WHERE receipt_id = $1
		ORDER BY addition_order DESC
		LIMIT 1
	`

	product := &model.Product{}
	err := r.db.QueryRowContext(ctx, query, receiptID).Scan(
		&product.ID,
		&product.ReceivedAt,
		&product.Type,
		&product.ReceiptID,
		&product.AdditionOrder,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return product, nil
} 