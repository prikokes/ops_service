package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ops_service/internal/model"
	"time"
)

// ReceiptRepo is a PostgreSQL implementation of ReceiptRepository
type ReceiptRepo struct {
	db *sql.DB
}

// NewReceiptRepo creates a new ReceiptRepo
func NewReceiptRepo(db *sql.DB) *ReceiptRepo {
	return &ReceiptRepo{db: db}
}

// Create creates a new receipt in the database
func (r *ReceiptRepo) Create(ctx context.Context, receipt *model.Receipt) error {
	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if there's already an open receipt for the pickup point
	var count int
	checkQuery := `
		SELECT COUNT(*)
		FROM receipts
		WHERE pickup_point_id = $1 AND status = $2
	`
	err = tx.QueryRowContext(ctx, checkQuery, receipt.PickupPointID, model.ReceiptStatusInProgress).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing receipts: %w", err)
	}

	if count > 0 {
		return model.ErrReceiptAlreadyClosed
	}

	// Create the receipt
	query := `
		INSERT INTO receipts (created_at, pickup_point_id, status)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err = tx.QueryRowContext(
		ctx,
		query,
		receipt.CreatedAt,
		receipt.PickupPointID,
		receipt.Status,
	).Scan(&receipt.ID)
	if err != nil {
		return fmt.Errorf("failed to create receipt: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetByID returns a receipt by ID
func (r *ReceiptRepo) GetByID(ctx context.Context, id int) (*model.Receipt, error) {
	query := `
		SELECT id, created_at, pickup_point_id, status
		FROM receipts
		WHERE id = $1
	`

	receipt := &model.Receipt{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&receipt.ID,
		&receipt.CreatedAt,
		&receipt.PickupPointID,
		&receipt.Status,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("receipt not found")
		}
		return nil, err
	}

	return receipt, nil
}

// GetLastOpenByPickupPointID returns the last open receipt for a pickup point
func (r *ReceiptRepo) GetLastOpenByPickupPointID(ctx context.Context, pickupPointID int) (*model.Receipt, error) {
	query := `
		SELECT id, created_at, pickup_point_id, status
		FROM receipts
		WHERE pickup_point_id = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	receipt := &model.Receipt{}
	err := r.db.QueryRowContext(ctx, query, pickupPointID, model.ReceiptStatusInProgress).Scan(
		&receipt.ID,
		&receipt.CreatedAt,
		&receipt.PickupPointID,
		&receipt.Status,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNoOpenReceipt
		}
		return nil, err
	}

	return receipt, nil
}

// Close marks a receipt as closed
func (r *ReceiptRepo) Close(ctx context.Context, id int) error {
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
	err = tx.QueryRowContext(ctx, checkQuery, id).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("receipt not found")
		}
		return fmt.Errorf("failed to check receipt status: %w", err)
	}

	if status == model.ReceiptStatusClosed {
		return model.ErrReceiptAlreadyClosed
	}

	// Close the receipt
	query := `
		UPDATE receipts
		SET status = $1
		WHERE id = $2
	`

	result, err := tx.ExecContext(ctx, query, model.ReceiptStatusClosed, id)
	if err != nil {
		return fmt.Errorf("failed to close receipt: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if affected == 0 {
		return errors.New("receipt not found")
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetByDateRange returns a list of receipts for a pickup point within a date range
func (r *ReceiptRepo) GetByDateRange(ctx context.Context, pickupPointID int, startDate, endDate time.Time) ([]*model.Receipt, error) {
	query := `
		SELECT id, created_at, pickup_point_id, status
		FROM receipts
		WHERE pickup_point_id = $1 AND created_at BETWEEN $2 AND $3
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, pickupPointID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query receipts: %w", err)
	}
	defer rows.Close()

	var receipts []*model.Receipt
	for rows.Next() {
		receipt := &model.Receipt{}
		if err := rows.Scan(&receipt.ID, &receipt.CreatedAt, &receipt.PickupPointID, &receipt.Status); err != nil {
			return nil, fmt.Errorf("failed to scan receipt: %w", err)
		}
		receipts = append(receipts, receipt)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over receipts: %w", err)
	}

	return receipts, nil
} 