package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionType string

const (
	Transfer       TransactionType = "transfer"
	PendingPayment TransactionType = "pending_payment"
	Confirmed      TransactionType = "confirmed"
	Processing     TransactionType = "processing"
	Shipped        TransactionType = "shipped"
	Delivered      TransactionType = "delivered"
	Cancelled      TransactionType = "cancelled"
	Expired        TransactionType = "expired"
)

type Transaction struct {
	ID              uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Type            TransactionType `gorm:"type:transaction_type;not null" json:"type"`
	ProductID       uuid.UUID       `gorm:"type:uuid;not null;index" json:"product_id"`
	Quantity        int             `gorm:"not null" json:"quantity"`
	WarehouseID     uuid.UUID       `gorm:"type:uuid;not null;index" json:"warehouse_id"`
	ToWarehouseID   *uuid.UUID      `gorm:"type:uuid;index" json:"to_warehouse_id,omitempty"`
	ReferenceNumber string          `gorm:"type:varchar(100);index" json:"reference_number,omitempty"`
	Notes           string          `gorm:"type:text" json:"notes,omitempty"`
	CreatedBy       uuid.UUID       `gorm:"type:uuid;not null;index" json:"created_by"`
	CreatedAt       time.Time       `gorm:"type:timestamptz;default:now();index" json:"created_at"`
}

func (Transaction) TableName() string {
	return "transactions"
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	return
}
