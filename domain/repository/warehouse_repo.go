package repository

import (
	"wms-be/domain/models"
	"wms-be/infrastructure/database"

	"github.com/google/uuid"
)

type WarehouseRepository interface {
	CreateWarehouse(warehouse *models.Warehouse) error
	GetWarehouseByID(id uuid.UUID) (*models.Warehouse, error)
	GetWarehouses() ([]models.Warehouse, error)
	UpdateWarehouse(warehouse *models.Warehouse) error
	DeleteWarehouse(id uuid.UUID) error
}

type warehouseRepository struct{}

func NewWarehouseRepository() WarehouseRepository {
	return &warehouseRepository{}
}

func (r *warehouseRepository) CreateWarehouse(warehouse *models.Warehouse) error {
	return database.DB.Create(warehouse).Error
}

func (r *warehouseRepository) GetWarehouseByID(id uuid.UUID) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	if err := database.DB.Where("id = ?", id).First(&warehouse).Error; err != nil {
		return nil, err
	}
	return &warehouse, nil
}

func (r *warehouseRepository) GetWarehouses() ([]models.Warehouse, error) {
	var warehouses []models.Warehouse
	if err := database.DB.Find(&warehouses).Error; err != nil {
		return nil, err
	}
	return warehouses, nil
}

func (r *warehouseRepository) UpdateWarehouse(warehouse *models.Warehouse) error {
	return database.DB.Save(warehouse).Error
}

func (r *warehouseRepository) DeleteWarehouse(id uuid.UUID) error {
	return database.DB.Where("id = ?", id).Delete(&models.Warehouse{}).Error
}
