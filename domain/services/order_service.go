package services

import (
	"errors"
	"wms-be/domain/models"
	"wms-be/domain/repository"
)

type OrderService interface {
	CreateOrder(order *models.Order) (*models.Order, error)
	UpdateOrderStatus(id string, status string) (*models.Order, error)
	GetOrders(page, limit int, filters map[string]interface{}) ([]models.Order, int64, error)
	GetOrderByID(id string) (*models.Order, error)
}

type orderService struct {
	orderRepo repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) OrderService {
	return &orderService{orderRepo: orderRepo}
}

func (s *orderService) CreateOrder(order *models.Order) (*models.Order, error) {
	if order == nil {
		return nil, errors.New("order cannot be nil")
	}

	// Create order beserta items
	createdOrder, err := s.orderRepo.CreateOrder(order)
	if err != nil {
		return nil, err
	}

	// Preload items setelah create agar response lengkap
	createdOrderWithItems, err := s.orderRepo.GetOrderByID(createdOrder.ID.String())
	if err != nil {
		return nil, err
	}

	return createdOrderWithItems, nil
}

func (s *orderService) UpdateOrderStatus(id string, status string) (*models.Order, error) {
	if id == "" {
		return nil, errors.New("order ID cannot be empty")
	}
	if status == "" {
		return nil, errors.New("status cannot be empty")
	}
	return s.orderRepo.UpdateOrderStatus(id, status)
}

func (s *orderService) GetOrders(page, limit int, filters map[string]interface{}) ([]models.Order, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return s.orderRepo.GetOrders(page, limit, filters)
}

func (s *orderService) GetOrderByID(id string) (*models.Order, error) {
	if id == "" {
		return nil, errors.New("order ID cannot be empty")
	}
	return s.orderRepo.GetOrderByID(id)
}
