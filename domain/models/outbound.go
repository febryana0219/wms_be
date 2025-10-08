package models

import (
	"time"

	"github.com/google/uuid"
)

type Outbound struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProductID          uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	Product            Product   `gorm:"foreignKey:ProductID" json:"product"`
	WarehouseID        uuid.UUID `gorm:"type:uuid;not null;index" json:"warehouse_id"`
	Warehouse          Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse"`
	Quantity           int       `json:"quantity"`
	DestinationType    string    `gorm:"type:outbound_destination_type;default:'customer'" json:"destination_type"`
	DestinationName    string    `gorm:"type:varchar(100);not null" json:"destination_name"`
	DestinationContact string    `gorm:"type:varchar(100);" json:"destination_contact,omitempty"`
	ReferenceNumber    string    `gorm:"type:varchar(100);" json:"reference_number,omitempty"`
	UnitPrice          float64   `json:"unit_price,omitempty"`
	TotalPrice         float64   `json:"total_price,omitempty"`
	Notes              string    `json:"notes,omitempty"`
	ShippedDate        time.Time `json:"shipped_date"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedBy          uuid.UUID `gorm:"type:uuid;not null;index" json:"created_by"`
	User               User      `gorm:"foreignKey:CreatedBy" json:"user"`
}

func (Outbound) TableName() string {
	return "outbounds"
}
