package repository

import (
	"wms-be/domain/models"
	"wms-be/infrastructure/database"

	"gorm.io/gorm"
)

type OutboundRepository interface {
	GetOutbounds(search, warehouseId string, page, limit int) ([]models.Outbound, int, error)
	CreateOutbound(outbound models.Outbound) (models.Outbound, error)
}

type outboundRepo struct {
	db *gorm.DB
}

func NewOutboundRepository() OutboundRepository {
	return &outboundRepo{db: database.GetDB()}
}

// GetOutbounds retrieves outbound records with optional filters.
func (r *outboundRepo) GetOutbounds(search, warehouseId string, page, limit int) ([]models.Outbound, int, error) {
	var outbounds []models.Outbound
	var total int64

	// Start query
	query := r.db.Model(&models.Outbound{}).
		Joins("JOIN products ON products.id = outbounds.product_id").
		Preload("Warehouse").
		Preload("Product").
		Preload("User")

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"products.name ILIKE ? OR products.sku ILIKE ? OR outbounds.destination_name ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	if warehouseId != "" {
		query = query.Where("warehouse_id = ?", warehouseId)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page - 1) * limit).Limit(limit).Find(&outbounds).Error
	if err != nil {
		return nil, 0, err
	}

	return outbounds, int(total), nil
}

func (r *outboundRepo) CreateOutbound(outbound models.Outbound) (models.Outbound, error) {
	err := r.db.Create(&outbound).Error
	return outbound, err
}
