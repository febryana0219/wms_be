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

	// Ambil data dari view dashboard_summary
	if err := d.db.Table("public.dashboard_summary").Take(&summary).Error; err != nil {
		return summary, err
	}

	return summary, nil
}
