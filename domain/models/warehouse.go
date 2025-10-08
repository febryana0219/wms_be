package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Warehouse struct {
	ID                 uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name               string    `gorm:"type:varchar(100);not null" json:"name"`
	Code               string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Address            string    `gorm:"type:text" json:"address"`
	Phone              string    `gorm:"type:varchar(50)" json:"phone"`
	Email              string    `gorm:"type:varchar(100)" json:"email"`
	Manager            string    `gorm:"type:varchar(100)" json:"manager"`
	IsActive           bool      `gorm:"default:true" json:"isActive"`
	Capacity           int       `gorm:"type:integer" json:"capacity"`
	CurrentUtilization int       `gorm:"type:integer" json:"currentUtilization"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

func (Warehouse) TableName() string {
	return "warehouses"
}

func (warehouse *Warehouse) BeforeCreate(tx *gorm.DB) error {
	if warehouse.ID == uuid.Nil {
		warehouse.ID = uuid.New()
	}
	return nil
}
