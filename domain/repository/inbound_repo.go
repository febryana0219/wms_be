package repository

import (
	"wms-be/domain/models"
	"wms-be/infrastructure/database"

	"gorm.io/gorm"
)

type InboundRepository interface {
	GetInbounds(search, warehouseId string, page, limit int) ([]models.Inbound, int, error)
	CreateInbound(inbound models.Inbound) (models.Inbound, error)
}

type inboundRepo struct {
	db *gorm.DB
}

func NewInboundRepository() InboundRepository {
	return &inboundRepo{db: database.GetDB()}
}

func (r *inboundRepo) GetInbounds(search, warehouseId string, page, limit int) ([]models.Inbound, int, error) {
	var inbounds []models.Inbound
	var total int64

	// Start query
	query := r.db.Model(&models.Inbound{}).
		Joins("JOIN products ON products.id = inbounds.product_id").
		Preload("Warehouse").
		Preload("Product").
		Preload("User")

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"products.name ILIKE ? OR products.sku ILIKE ? OR inbounds.supplier_name ILIKE ?",
			searchPattern, searchPattern, searchPattern,
		)
	}

	if warehouseId != "" {
		query = query.Where("inbounds.warehouse_id = ?", warehouseId)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page - 1) * limit).Limit(limit).Find(&inbounds).Error
	if err != nil {
		return nil, 0, err
	}

	return inbounds, int(total), nil
}

func (r *inboundRepo) CreateInbound(inbound models.Inbound) (models.Inbound, error) {
	err := r.db.Create(&inbound).Error
	return inbound, err
}
