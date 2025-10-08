package services

import (
	"wms-be/domain/models"
	"wms-be/domain/repository"
)

type IOutboundService interface {
	GetOutbounds(search, warehouseId string, page, limit int) ([]models.Outbound, int, error)
	GetAllOutbounds() ([]models.Outbound, error)
	CreateOutbound(outbound models.Outbound) (models.Outbound, error)
}

type OutboundService struct {
	outboundRepo repository.OutboundRepository
}

// Constructor
func NewOutboundService(outboundRepo repository.OutboundRepository) *OutboundService {
	return &OutboundService{outboundRepo: outboundRepo}
}

func (s *OutboundService) GetOutbounds(search, warehouseId string, page, limit int) ([]models.Outbound, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	outbounds, total, err := s.outboundRepo.GetOutbounds(search, warehouseId, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return outbounds, total, nil
}

func (s *OutboundService) GetAllOutbounds() ([]models.Outbound, error) {
	outbounds, _, err := s.outboundRepo.GetOutbounds("", "", 0, 0)
	if err != nil {
		return nil, err
	}
	return outbounds, nil
}

func (s *OutboundService) CreateOutbound(outbound models.Outbound) (models.Outbound, error) {
	createdOutbound, err := s.outboundRepo.CreateOutbound(outbound)
	if err != nil {
		return models.Outbound{}, err
	}
	return createdOutbound, nil
}
