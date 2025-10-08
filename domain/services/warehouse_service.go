package services

import (
	"errors"
	"wms-be/domain/models"
	"wms-be/domain/repository"

	"github.com/google/uuid"
)

type WarehouseService interface {
	CreateWarehouse(warehouse *models.Warehouse) error
	GetWarehouseByID(id uuid.UUID) (*models.Warehouse, error)
	GetWarehouses() ([]models.Warehouse, error)
	UpdateWarehouse(warehouse *models.Warehouse) error
	DeleteWarehouse(id uuid.UUID) error
}

type warehouseService struct {
	warehouseRepo repository.WarehouseRepository
}

func NewWarehouseService(warehouseRepo repository.WarehouseRepository) WarehouseService {
	return &warehouseService{
		warehouseRepo: warehouseRepo,
	}
}

func (s *warehouseService) CreateWarehouse(warehouse *models.Warehouse) error {
	return s.warehouseRepo.CreateWarehouse(warehouse)
}

func (s *warehouseService) GetWarehouseByID(id uuid.UUID) (*models.Warehouse, error) {
	warehouse, err := s.warehouseRepo.GetWarehouseByID(id)
	if err != nil {
		return nil, errors.New("warehouse not found")
	}
	return warehouse, nil
}

func (s *warehouseService) GetWarehouses() ([]models.Warehouse, error) {
	return s.warehouseRepo.GetWarehouses()
}

func (s *warehouseService) UpdateWarehouse(warehouse *models.Warehouse) error {
	return s.warehouseRepo.UpdateWarehouse(warehouse)
}

func (s *warehouseService) DeleteWarehouse(id uuid.UUID) error {
	return s.warehouseRepo.DeleteWarehouse(id)
}
