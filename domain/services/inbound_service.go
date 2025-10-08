package services

import (
	"wms-be/domain/models"
	"wms-be/domain/repository"
)

type IInboundService interface {
	GetInbounds(search, warehouseId string, page, limit int) ([]models.Inbound, int, error)
	GetAllInbounds() ([]models.Inbound, error)
	CreateInbound(inbound models.Inbound) (models.Inbound, error)
}

type InboundService struct {
	inboundRepo repository.InboundRepository
}

// Constructor
func NewInboundService(inboundRepo repository.InboundRepository) *InboundService {
	return &InboundService{inboundRepo: inboundRepo}
}

func (s *InboundService) GetInbounds(search, warehouseId string, page, limit int) ([]models.Inbound, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	inbounds, total, err := s.inboundRepo.GetInbounds(search, warehouseId, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return inbounds, total, nil
}

func (s *InboundService) GetAllInbounds() ([]models.Inbound, error) {
	inbounds, _, err := s.inboundRepo.GetInbounds("", "", 0, 0)
	if err != nil {
		return nil, err
	}
	return inbounds, nil
}

func (s *InboundService) CreateInbound(inbound models.Inbound) (models.Inbound, error) {
	createdInbound, err := s.inboundRepo.CreateInbound(inbound)
	if err != nil {
		return models.Inbound{}, err
	}
	return createdInbound, nil
}
