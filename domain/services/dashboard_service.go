package services

import (
	"wms-be/domain/models"

	"gorm.io/gorm"
)

type DashboardService struct {
	db *gorm.DB
}

func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{db: db}
}

func (d *DashboardService) GetDashboardSummary() (models.DashboardSummary, error) {
	var summary models.DashboardSummary

	// Ambil 1 baris dari view dashboard_summary
	if err := d.db.Table("public.dashboard_summary").Take(&summary).Error; err != nil {
		return summary, err
	}

	// Ambil 10 data terakhir dari transaction_histories
	var histories []models.TransactionHistory
	if err := d.db.Table("public.transaction_histories").
		Where("dashboard_id = ?", summary.ID).
		Order("created_at DESC").
		Limit(10).
		Find(&histories).Error; err != nil {
		return summary, err
	}

	// Ambil 10 data low stock dari view low_stock_products
	var lowStocks []models.LowStockProduct
	if err := d.db.Table("public.low_stock_products").
		Where("dashboard_id = ?", summary.ID).
		Order("available_stock ASC").
		Limit(10).
		Find(&lowStocks).Error; err != nil {
		return summary, err
	}

	// Simpan hasil query ke struct utama
	summary.TransactionHistories = histories
	summary.LowStockProducts = lowStocks

	return summary, nil
}
