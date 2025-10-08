package models

import (
	"time"

	"github.com/google/uuid"
)

type Inbound struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProductID       uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	Product         Product   `gorm:"foreignKey:ProductID"`
	WarehouseID     uuid.UUID `gorm:"type:uuid;not null" json:"warehouse_id"`
	Warehouse       Warehouse `gorm:"foreignKey:WarehouseID"`
	Quantity        int       `json:"quantity"`
	SupplierName    string    `gorm:"type:varchar(100);not null" json:"supplier_name"`
	SupplierContact string    `gorm:"type:varchar(100);" json:"supplier_contact,omitempty"`
	ReferenceNumber string    `gorm:"type:varchar(100);" json:"reference_number,omitempty"`
	UnitCost        float64   `json:"unit_cost,omitempty"`
	TotalCost       float64   `json:"total_cost,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	ReceivedDate    time.Time `json:"received_date"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	User            User      `gorm:"foreignKey:CreatedBy"`
}

func (Inbound) TableName() string {
	return "inbounds"
}
