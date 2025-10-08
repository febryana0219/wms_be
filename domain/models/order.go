package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID           uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderNumber  string      `gorm:"type:varchar(50);unique;not null" json:"order_number"`
	CustomerID   string      `gorm:"type:varchar(100);not null" json:"customer_id"`
	CustomerName string      `gorm:"type:varchar(100);not null" json:"customer_name"`
	Status       string      `gorm:"type:order_status;default:pending_payment" json:"status"`
	TotalAmount  float64     `gorm:"type:numeric(15,2);default:0.00" json:"total_amount"`
	WarehouseID  uuid.UUID   `gorm:"type:uuid;not null" json:"warehouse_id"`
	Warehouse    Warehouse   `gorm:"foreignKey:WarehouseID"`
	Notes        string      `gorm:"type:text" json:"notes"`
	ExpiresAt    time.Time   `gorm:"type:timestamptz" json:"expires_at"`
	CreatedAt    time.Time   `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"type:timestamptz;default:now()" json:"updated_at"`
	OrderItems   []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items"`
}

type OrderItem struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID    uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`
	ProductID  uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	Product    Product   `gorm:"foreignKey:ProductID"`
	Quantity   int       `gorm:"not null" json:"quantity"`
	UnitPrice  float64   `gorm:"type:numeric(15,2);not null" json:"unit_price"`
	TotalPrice float64   `gorm:"type:numeric(15,2);not null" json:"total_price"`
	CreatedAt  time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`
}

func (Order) TableName() string {
	return "orders"
}

func (OrderItem) TableName() string {
	return "order_items"
}
