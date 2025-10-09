package models

import (
	"wms-be/infrastructure/database"
)

type TransactionHistory struct {
	ID              string  `json:"id"`
	Type            string  `json:"type"`
	ProductID       string  `json:"product_id"`
	ProductName     string  `json:"product_name"`
	SKU             string  `json:"sku"`
	Quantity        float64 `json:"quantity"`
	WarehouseID     string  `json:"warehouse_id"`
	WarehouseName   string  `json:"warehouse_name"`
	ToWarehouseID   *string `json:"to_warehouse_id"`
	ToWarehouseName *string `json:"to_warehouse_name"`
	ReferenceNumber string  `json:"reference_number"`
	Notes           string  `json:"notes"`
	CreatedBy       string  `json:"created_by"`
	CreatedByName   string  `json:"created_by_name"`
	CreatedAt       string  `json:"created_at"`
	DashboardID     string  `json:"dashboard_id"`
}

func GetAllTransactionHistory() ([]TransactionHistory, error) {
	var histories []TransactionHistory
	db := database.GetDB()

	err := db.Table("transaction_histories").Find(&histories).Error
	return histories, err
}
