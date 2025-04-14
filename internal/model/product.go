package model

import (
	"time"
)

// ProductType represents the type of product
type ProductType string

const (
	// ProductTypeElectronics represents electronics products
	ProductTypeElectronics ProductType = "электроника"
	// ProductTypeClothing represents clothing products
	ProductTypeClothing ProductType = "одежда"
	// ProductTypeFootwear represents footwear products
	ProductTypeFootwear ProductType = "обувь"
)

// Product represents a product in the system
type Product struct {
	ID           int         `json:"id" db:"id"`
	ReceivedAt   time.Time   `json:"received_at" db:"received_at"`
	Type         ProductType `json:"type" db:"type"`
	ReceiptID    int         `json:"receipt_id" db:"receipt_id"`
	AdditionOrder int         `json:"addition_order" db:"addition_order"`
} 