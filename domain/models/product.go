package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SKU            string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"sku"`
	Name           string    `gorm:"type:varchar(100);not null" json:"name"`
	Category       string    `gorm:"type:varchar(100)" json:"category"`
	Description    string    `gorm:"type:text" json:"description"`
	Price          float64   `gorm:"type:decimal(15,2);not null;default:0.00" json:"price"`
	MinStock       int       `gorm:"type:integer;not null;default:0" json:"min_stock"`
	Stock          int       `gorm:"type:integer;not null;default:0" json:"stock"`
	ReservedStock  int       `gorm:"type:integer;not null;default:0" json:"reserved_stock"`
	AvailableStock int       `gorm:"->;type:integer" json:"available_stock"` // read-only
	WarehouseID    uuid.UUID `gorm:"type:uuid;not null" json:"warehouse_id"`
	Warehouse      Warehouse `gorm:"foreignKey:WarehouseID"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}
