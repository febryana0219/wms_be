package handler

import (
	"net/http"
	"strconv"
	"time"
	"wms-be/domain/models"
	"wms-be/domain/services"
	"wms-be/interfaces/http/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	OrderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{OrderService: orderService}
}

// Response structs
type OrderItemResponse struct {
	ID          string  `json:"id"`
	OrderID     string  `json:"order_id"`
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
	CreatedAt   string  `json:"created_at"`
}

type OrderResponse struct {
	ID            string              `json:"id"`
	OrderNumber   string              `json:"order_number"`
	CustomerID    string              `json:"customer_id"`
	CustomerName  string              `json:"customer_name"`
	Status        string              `json:"status"`
	TotalAmount   float64             `json:"total_amount"`
	WarehouseID   string              `json:"warehouse_id"`
	WarehouseName string              `json:"warehouse_name"`
	Notes         string              `json:"notes"`
	ExpiresAt     string              `json:"expires_at"`
	CreatedAt     string              `json:"created_at"`
	UpdatedAt     string              `json:"updated_at"`
	Items         []OrderItemResponse `json:"items"`
}

// mapOrderToResponse konversi Order ke OrderResponse
func mapOrderToResponse(order models.Order) OrderResponse {
	// Mapping items
	items := make([]OrderItemResponse, 0) // jangan nil
	for _, i := range order.OrderItems {
		productName := ""
		if i.Product != (models.Product{}) {
			productName = i.Product.Name
		}
		items = append(items, OrderItemResponse{
			ID:          i.ID.String(),
			OrderID:     i.OrderID.String(),
			ProductID:   i.ProductID.String(),
			ProductName: productName,
			Quantity:    i.Quantity,
			UnitPrice:   i.UnitPrice,
			TotalPrice:  i.TotalPrice,
			CreatedAt:   i.CreatedAt.Format(time.RFC3339),
		})
	}

	warehouseName := ""
	if order.Warehouse != (models.Warehouse{}) {
		warehouseName = order.Warehouse.Name
	}

	return OrderResponse{
		ID:            order.ID.String(),
		OrderNumber:   order.OrderNumber,
		CustomerID:    order.CustomerID,
		CustomerName:  order.CustomerName,
		Status:        order.Status,
		TotalAmount:   order.TotalAmount,
		WarehouseID:   order.WarehouseID.String(),
		WarehouseName: warehouseName,
		Notes:         order.Notes,
		ExpiresAt:     order.ExpiresAt.Format(time.RFC3339),
		CreatedAt:     order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     order.UpdatedAt.Format(time.RFC3339),
		Items:         items,
	}
}

// GET /orders
func (h *OrderHandler) GetOrders(c *gin.Context) {
	search := c.Query("search")
	status := c.Query("status")
	warehouseId := c.Query("warehouse_id")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	filters := make(map[string]interface{})
	if search != "" {
		filters["search"] = search
	}
	if status != "" {
		filters["status"] = status
	}
	if warehouseId != "" {
		filters["warehouse_id"] = warehouseId
	}

	orders, total, err := h.OrderService.GetOrders(page, limit, filters)
	if err != nil {
		response.ErrorMessageResponse(c, err, http.StatusInternalServerError)
		return
	}

	var orderResponses []OrderResponse
	orderResponses = make([]OrderResponse, 0) // pastikan tidak nil
	for _, o := range orders {
		orderResponses = append(orderResponses, mapOrderToResponse(o))
	}

	response.PaginatedResponse(c, "orders", orderResponses, int(total), page, limit)
}

// GET /orders/:id
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	order, err := h.OrderService.GetOrderByID(id)
	if err != nil {
		response.ErrorMessageResponse(c, err, http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, mapOrderToResponse(*order), "Order retrieved successfully")
}

// POST /orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var orderRequest struct {
		OrderNumber  string `json:"order_number"`
		CustomerID   string `json:"customer_id"`
		CustomerName string `json:"customer_name"`
		WarehouseID  string `json:"warehouse_id"`
		Items        []struct {
			ProductID string `json:"product_id"`
			Quantity  int    `json:"quantity"`
		} `json:"items"`
		Notes     string    `json:"notes"`
		ExpiresAt time.Time `json:"expires_at"`
	}

	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		response.ErrorMessageResponse(c, err, http.StatusBadRequest)
		return
	}

	warehouseUUID, err := uuid.Parse(orderRequest.WarehouseID)
	if err != nil {
		response.ErrorMessageResponse(c, err, http.StatusBadRequest)
		return
	}

	order := &models.Order{
		OrderNumber:  orderRequest.OrderNumber,
		CustomerID:   orderRequest.CustomerID,
		CustomerName: orderRequest.CustomerName,
		WarehouseID:  warehouseUUID,
		Notes:        orderRequest.Notes,
		ExpiresAt:    orderRequest.ExpiresAt,
	}

	for _, item := range orderRequest.Items {
		productUUID, err := uuid.Parse(item.ProductID)
		if err != nil {
			response.ErrorMessageResponse(c, err, http.StatusBadRequest)
			return
		}

		order.OrderItems = append(order.OrderItems, models.OrderItem{
			ProductID: productUUID,
			Quantity:  item.Quantity,
		})
	}

	createdOrder, err := h.OrderService.CreateOrder(order)
	if err != nil {
		response.ErrorMessageResponse(c, err, http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, mapOrderToResponse(*createdOrder), "Order created successfully")
}

// PUT /orders/:id/status
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var statusUpdate struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		response.ErrorMessageResponse(c, err, http.StatusBadRequest)
		return
	}

	order, err := h.OrderService.UpdateOrderStatus(id, statusUpdate.Status)
	if err != nil {
		response.ErrorMessageResponse(c, err, http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, mapOrderToResponse(*order), "Order status updated successfully")
}
