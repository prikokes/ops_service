package model

import (
	"errors"
	"time"
)

// ReceiptStatus represents the status of a receipt
type ReceiptStatus string

const (
	// ReceiptStatusInProgress represents a receipt that is still in progress
	ReceiptStatusInProgress ReceiptStatus = "in_progress"
	// ReceiptStatusClosed represents a receipt that has been closed
	ReceiptStatusClosed ReceiptStatus = "closed"
)

// Common errors for receipt operations
var (
	ErrReceiptAlreadyClosed = errors.New("receipt already closed")
	ErrNoOpenReceipt        = errors.New("no open receipt available")
)

// Receipt represents a product receipt at a pickup point
type Receipt struct {
	ID            int           `json:"id" db:"id"`
	CreatedAt     time.Time     `json:"created_at" db:"created_at"`
	PickupPointID int           `json:"pickup_point_id" db:"pickup_point_id"`
	Status        ReceiptStatus `json:"status" db:"status"`
	Products      []*Product    `json:"products,omitempty" db:"-"` // Products is populated separately
} 