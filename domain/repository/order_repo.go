package repository

import (
	"wms-be/domain/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *models.Order) (*models.Order, error)
	UpdateOrderStatus(id string, status string) (*models.Order, error)
	GetOrders(page, limit int, filters map[string]interface{}) ([]models.Order, int64, error)
	GetOrderByID(id string) (*models.Order, error)
}

type orderRepo struct {
	DB *gorm.DB
}

func NewOrderRepository(DB *gorm.DB) OrderRepository {
	return &orderRepo{DB: DB}
}

// CreateOrder membuat order beserta items
func (r *orderRepo) CreateOrder(order *models.Order) (*models.Order, error) {
	err := r.DB.Create(order).Error
	if err != nil {
		return nil, err
	}
	// Preload nested relation setelah create
	err = r.DB.Preload("Warehouse").Preload("OrderItems.Product").First(order, "id = ?", order.ID).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

// UpdateOrderStatus update status dan kembalikan order lengkap
func (r *orderRepo) UpdateOrderStatus(id string, status string) (*models.Order, error) {
	var order models.Order
	// Update status
	err := r.DB.Model(&order).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		return nil, err
	}
	// Ambil kembali order beserta items & product
	err = r.DB.Preload("Warehouse").Preload("OrderItems.Product").First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetOrders ambil list order dengan filter dan pagination
func (r *orderRepo) GetOrders(page, limit int, filters map[string]interface{}) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.DB.Model(&models.Order{})

	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("order_number LIKE ?", "%"+search+"%")
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if warehouseId, ok := filters["warehouse_id"].(string); ok && warehouseId != "" {
		query = query.Where("warehouse_id = ?", warehouseId)
	}

	// Hitung total
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Ambil data dengan preload nested relation
	err = query.Preload("Warehouse").Preload("OrderItems.Product").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetOrderByID ambil order lengkap berdasarkan ID
func (r *orderRepo) GetOrderByID(id string) (*models.Order, error) {
	var order models.Order
	// Preload nested relation
	err := r.DB.Preload("Warehouse").Preload("OrderItems.Product").First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
