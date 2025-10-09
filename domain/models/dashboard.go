package models

type DashboardSummary struct {
	ID                   string               `json:"id"`
	TotalProducts        int                  `json:"totalProducts"`
	TotalWarehouses      int                  `json:"totalWarehouses"`
	TotalOrders          int                  `json:"totalOrders"`
	TotalTransactions    int                  `json:"totalTransactions"`
	ActiveWarehouses     int                  `json:"activeWarehouses"`
	PendingOrders        int                  `json:"pendingOrders"`
	TransactionHistories []TransactionHistory `json:"transactionHistories" gorm:"foreignKey:DashboardID"`
	LowStockProducts     []LowStockProduct    `json:"lowStockProducts" gorm:"foreignKey:DashboardID"`
}
