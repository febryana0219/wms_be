package models

import (
	"github.com/google/uuid"
)

type LowStockProduct struct {
	ID               string    `json:"id"`
	SKU              string    `json:"sku"`
	Name             string    `json:"name"`
	Category         string    `json:"category"`
	Description      string    `json:"description"`
	Price            float64   `json:"price"`
	MinStock         int       `json:"min_stock"`
	Stock            int       `json:"stock"`
	ReservedStock    int       `json:"reserved_stock"`
	WarehouseID      uuid.UUID `json:"warehouse_id"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        *string   `json:"created_at"`
	UpdatedAt        *string   `json:"updated_at"`
	AvailableStock   int       `json:"available_stock"`
	WarehouseName    string    `json:"warehouse_name"`
	ShortageQuantity int       `json:"shortage_quantity"`
	DashboardID      string    `json:"dashboard_id"`
}

func (LowStockProduct) TableName() string {
	return "low_stock_products"
}
