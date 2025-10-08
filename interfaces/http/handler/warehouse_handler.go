package handler

import (
	"net/http"
	"time"
	"wms-be/domain/models"
	"wms-be/domain/services"
	"wms-be/interfaces/http/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// WarehouseHandler handles warehouse-related requests
type WarehouseHandler struct {
	warehouseService services.WarehouseService
}

// Constructor
func NewWarehouseHandler(warehouseService services.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{
		warehouseService: warehouseService,
	}
}

// WarehouseResponse represents API response for a warehouse
type WarehouseResponse struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Code               string    `json:"code"`
	Address            string    `json:"address"`
	Phone              string    `json:"phone"`
	Email              string    `json:"email"`
	Manager            string    `json:"manager"`
	IsActive           bool      `json:"isActive"`
	Capacity           int       `json:"capacity"`
	CurrentUtilization int       `json:"currentUtilization"`
	CreatedAt          time.Time `json:"createdAt"`
}

// mapWarehouseToResponse maps models.Warehouse to WarehouseResponse
func mapWarehouseToResponse(w models.Warehouse) WarehouseResponse {
	return WarehouseResponse{
		ID:                 w.ID.String(),
		Name:               w.Name,
		Code:               w.Code,
		Address:            w.Address,
		Phone:              w.Phone,
		Email:              w.Email,
		Manager:            w.Manager,
		IsActive:           w.IsActive,
		Capacity:           w.Capacity,
		CurrentUtilization: w.CurrentUtilization,
		CreatedAt:          w.CreatedAt,
	}
}

// GetWarehouses handles GET /warehouses
func (h *WarehouseHandler) GetWarehouses(c *gin.Context) {
	warehouses, err := h.warehouseService.GetWarehouses()
	if err != nil {
		response.ErrorMessageResponse(c, err, http.StatusInternalServerError)
		return
	}

	var resp []WarehouseResponse
	for _, w := range warehouses {
		resp = append(resp, mapWarehouseToResponse(w))
	}

	response.SuccessResponse(c, resp, "Warehouses fetched successfully")
}

// GetWarehouseByID handles GET /warehouses/:id
func (h *WarehouseHandler) GetWarehouseByID(c *gin.Context) {
	id := c.Param("id")
	warehouseID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorMessageResponse(c, err, http.StatusBadRequest)
		return
	}

	warehouse, err := h.warehouseService.GetWarehouseByID(warehouseID)
	if err != nil {
		response.ErrorMessageResponse(c, err, http.StatusNotFound)
		return
	}

	response.SuccessResponse(c, mapWarehouseToResponse(*warehouse), "Warehouse fetched successfully")
}

// CreateWarehouse handles POST /warehouses
func (h *WarehouseHandler) CreateWarehouse(c *gin.Context) {
	var warehouse models.Warehouse
	if err := c.ShouldBindJSON(&warehouse); err != nil {
		response.ErrorMessageResponse(c, err, http.StatusBadRequest)
		return
	}

	if warehouse.ID == uuid.Nil {
		warehouse.ID = uuid.New()
	}

	if err := h.warehouseService.CreateWarehouse(&warehouse); err != nil {
		response.ErrorMessageResponse(c, err, http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, mapWarehouseToResponse(warehouse), "Warehouse created successfully")
}

// UpdateWarehouse handles PUT /warehouses/:id
func (h *WarehouseHandler) UpdateWarehouse(c *gin.Context) {
	id := c.Param("id")
	warehouseID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorMessageResponse(c, err, http.StatusBadRequest)
		return
	}

	var warehouse models.Warehouse
	if err := c.ShouldBindJSON(&warehouse); err != nil {
		response.ErrorMessageResponse(c, err, http.StatusBadRequest)
		return
	}

	warehouse.ID = warehouseID
	if err := h.warehouseService.UpdateWarehouse(&warehouse); err != nil {
		response.ErrorMessageResponse(c, err, http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, mapWarehouseToResponse(warehouse), "Warehouse updated successfully")
}

// DeleteWarehouse handles DELETE /warehouses/:id
func (h *WarehouseHandler) DeleteWarehouse(c *gin.Context) {
	id := c.Param("id")
	warehouseID, err := uuid.Parse(id)
	if err != nil {
		response.ErrorMessageResponse(c, err, http.StatusBadRequest)
		return
	}

	if err := h.warehouseService.DeleteWarehouse(warehouseID); err != nil {
		response.ErrorMessageResponse(c, err, http.StatusInternalServerError)
		return
	}

	response.SuccessMessageResponse(c, "Warehouse deleted successfully")
}
