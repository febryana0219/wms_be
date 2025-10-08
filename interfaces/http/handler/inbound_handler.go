package handler

import (
	"strconv"
	"time"
	"wms-be/domain/models"
	"wms-be/domain/services"
	"wms-be/interfaces/http/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InboundHandler struct {
	inboundService *services.InboundService
}

func NewInboundHandler(inboundService *services.InboundService) *InboundHandler {
	return &InboundHandler{inboundService: inboundService}
}

type InboundResponse struct {
	ID              string  `json:"id"`
	ProductID       string  `json:"product_id"`
	ProductName     string  `json:"product_name"`
	ProductSKU      string  `json:"product_sku"`
	WarehouseID     string  `json:"warehouse_id"`
	WarehouseName   string  `json:"warehouse_name"`
	Quantity        int     `json:"quantity"`
	SupplierName    string  `json:"supplier_name"`
	SupplierContact string  `json:"supplier_contact,omitempty"`
	ReferenceNumber string  `json:"reference_number,omitempty"`
	UnitCost        float64 `json:"unit_cost,omitempty"`
	TotalCost       float64 `json:"total_cost,omitempty"`
	Notes           string  `json:"notes,omitempty"`
	ReceivedDate    string  `json:"received_date"`
	CreatedAt       string  `json:"created_at"`
	CreatedBy       string  `json:"created_by"`
	CreatedByName   string  `json:"created_by_name"`
}

func mapInboundToResponse(inbound models.Inbound) InboundResponse {
	productName := ""
	productSKU := ""
	if inbound.Product != (models.Product{}) {
		productName = inbound.Product.Name
		productSKU = inbound.Product.SKU
	}

	warehouseName := ""
	if inbound.Warehouse != (models.Warehouse{}) {
		warehouseName = inbound.Warehouse.Name
	}

	createdByName := ""
	if inbound.User != (models.User{}) {
		createdByName = inbound.User.Name
	}

	// Return the populated InboundResponse
	return InboundResponse{
		ID:              inbound.ID.String(),
		ProductID:       inbound.ProductID.String(),
		ProductName:     productName,
		ProductSKU:      productSKU,
		WarehouseID:     inbound.WarehouseID.String(),
		WarehouseName:   warehouseName,
		Quantity:        inbound.Quantity,
		SupplierName:    inbound.SupplierName,
		SupplierContact: inbound.SupplierContact,
		ReferenceNumber: inbound.ReferenceNumber,
		UnitCost:        inbound.UnitCost,
		TotalCost:       inbound.TotalCost,
		Notes:           inbound.Notes,
		ReceivedDate:    inbound.ReceivedDate.Format(time.RFC3339),
		CreatedAt:       inbound.CreatedAt.Format(time.RFC3339),
		CreatedBy:       inbound.CreatedBy.String(),
		CreatedByName:   createdByName,
	}
}

// GET /inbounds
func (h *InboundHandler) GetInbounds(c *gin.Context) {
	search := c.Query("search")
	warehouseId := c.Query("warehouseId")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	// === cek apakah ada parameter pagination/filter ===
	isPaginated := pageStr != "" || limitStr != "" || search != "" || warehouseId != ""

	var (
		page, limit int
		err         error
	)

	if isPaginated {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 10
		}
	} else {
		// jika tidak ada pagination, ambil semua produk
		page = 1
		limit = 1000000 // jumlah sangat besar agar ambil semua
	}

	inbounds, total, err := h.inboundService.GetInbounds(search, warehouseId, page, limit)
	if err != nil {
		response.ErrorMessageResponse(c, err, 500)
		return
	}

	resp := make([]InboundResponse, len(inbounds))
	for i, p := range inbounds {
		resp[i] = mapInboundToResponse(p)
	}

	if isPaginated {
		response.PaginatedResponse(c, "inbounds", resp, total, page, limit)
	} else {
		response.SuccessResponse(c, resp, "Inbounds fetched successfully")
	}
}

// POST /inbounds
func (h *InboundHandler) CreateInbound(c *gin.Context) {
	var req struct {
		ProductID       string  `json:"product_id"`
		WarehouseID     string  `json:"warehouse_id"`
		Quantity        int     `json:"quantity"`
		SupplierName    string  `json:"supplier_name"`
		SupplierContact string  `json:"supplier_contact,omitempty"`
		ReferenceNumber string  `json:"reference_number,omitempty"`
		UnitCost        float64 `json:"unit_cost,omitempty"`
		Notes           string  `json:"notes,omitempty"`
		ReceivedDate    string  `json:"received_date"`
		CreatedBy       string  `json:"created_by"`
	}

	// Bind incoming JSON request to the struct
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorMessageResponse(c, err, 400)
		return
	}

	// Convert strings to UUIDs for ProductID, WarehouseID, and CreatedBy
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		response.ErrorMessageResponse(c, err, 400)
		return
	}

	warehouseID, err := uuid.Parse(req.WarehouseID)
	if err != nil {
		response.ErrorMessageResponse(c, err, 400)
		return
	}

	createdBy, err := uuid.Parse(req.CreatedBy)
	if err != nil {
		response.ErrorMessageResponse(c, err, 400)
		return
	}

	// Parse the received date string into time.Time
	receivedDate, err := time.Parse(time.RFC3339, req.ReceivedDate)
	if err != nil {
		response.ErrorMessageResponse(c, err, 400)
		return
	}

	// Create the Inbound model
	inbound := models.Inbound{
		ProductID:       productID,
		WarehouseID:     warehouseID,
		Quantity:        req.Quantity,
		SupplierName:    req.SupplierName,
		SupplierContact: req.SupplierContact,
		ReferenceNumber: req.ReferenceNumber,
		UnitCost:        req.UnitCost,
		TotalCost:       req.UnitCost * float64(req.Quantity),
		Notes:           req.Notes,
		ReceivedDate:    receivedDate,
		CreatedAt:       time.Now(),
		CreatedBy:       createdBy,
	}

	// Create the inbound record
	createdInbound, err := h.inboundService.CreateInbound(inbound)
	if err != nil {
		response.ErrorMessageResponse(c, err, 500)
		return
	}

	// Respond with the created inbound record
	response.SuccessResponse(c, mapInboundToResponse(createdInbound), "Inbound created successfully")
}
