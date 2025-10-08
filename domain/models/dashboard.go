package models

type DashboardSummary struct {
	TotalProducts     int `json:"totalProducts"`
	TotalWarehouses   int `json:"totalWarehouses"`
	TotalOrders       int `json:"totalOrders"`
	TotalTransactions int `json:"totalTransactions"`
	ActiveWarehouses  int `json:"activeWarehouses"`
	PendingOrders     int `json:"pendingOrders"`
}
