package services

import (
	"time"
	"wms-be/domain/models"

	"gorm.io/gorm"
)

type TransactionHistoryService struct {
	db *gorm.DB
}

func NewTransactionHistoryService(db *gorm.DB) *TransactionHistoryService {
	return &TransactionHistoryService{db: db}
}

// GetTransactions with filters and pagination
func (s *TransactionHistoryService) GetTransactions(
	txType string,
	warehouseId string,
	dateFrom string,
	dateTo string,
	page int,
	limit int,
) ([]models.TransactionHistory, int64, error) {

	var transactions []models.TransactionHistory
	var total int64

	query := s.db.Model(&models.TransactionHistory{})

	// Apply filters
	if txType != "" {
		query = query.Where("type = ?", txType)
	}
	if warehouseId != "" {
		query = query.Where("warehouse_id = ?", warehouseId)
	}

	layout := "2006-01-02" // format frontend

	if dateFrom != "" && dateTo != "" {
		from, err := time.Parse(layout, dateFrom)
		if err != nil {
			return nil, 0, err
		}
		to, err := time.Parse(layout, dateTo)
		if err != nil {
			return nil, 0, err
		}
		// include seluruh hari dateTo
		to = to.AddDate(0, 0, 1)
		query = query.Where("created_at >= ? AND created_at < ?", from, to)
	} else if dateFrom != "" {
		from, err := time.Parse(layout, dateFrom)
		if err != nil {
			return nil, 0, err
		}
		query = query.Where("created_at >= ?", from)
	} else if dateTo != "" {
		to, err := time.Parse(layout, dateTo)
		if err != nil {
			return nil, 0, err
		}
		to = to.AddDate(0, 0, 1)
		query = query.Where("created_at < ?", to)
	}

	// Hitung total data
	query.Count(&total)

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Ambil data
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}
