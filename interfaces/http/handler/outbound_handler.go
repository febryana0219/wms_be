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

type OutboundHandler struct {
	outboundService *services.OutboundService
}

func NewOutboundHandler(outboundService *services.OutboundService) *OutboundHandler {
	return &OutboundHandler{outboundService: outboundService}
}

type OutboundResponse struct {
	ID                 string  `json:"id"`
	ProductID          string  `json:"product_id"`
	ProductName        string  `json:"product_name"`
	ProductSKU         string  `json:"product_sku"`
	WarehouseID        string  `json:"warehouse_id"`
	WarehouseName      string  `json:"warehouse_name"`
	Quantity           int     `json:"quantity"`
	DestinationType    string  `json:"destination_type"`
	DestinationName    string  `json:"destination_name"`
	DestinationContact string  `json:"destination_contact,omitempty"`
	ReferenceNumber    string  `json:"reference_number,omitempty"`
	UnitPrice          float64 `json:"unit_price,omitempty"`
	TotalPrice         float64 `json:"total_price,omitempty"`
	Notes              string  `json:"notes,omitempty"`
	ShippedDate        string  `json:"shipped_date"`
	CreatedAt          string  `json:"created_at"`
	CreatedBy          string  `json:"created_by"`
	CreatedByName      string  `json:"created_by_name"`
}

func mapOutboundToResponse(outbound models.Outbound) OutboundResponse {
	return OutboundResponse{
		ID:                 outbound.ID.String(),
		ProductID:          outbound.ProductID.String(),
		ProductName:        outbound.Product.Name,
		ProductSKU:         outbound.Product.SKU,
		WarehouseID:        outbound.WarehouseID.String(),
		WarehouseName:      outbound.Warehouse.Name,
		Quantity:           outbound.Quantity,
		DestinationName:    outbound.DestinationName,
		DestinationType:    outbound.DestinationType,
		DestinationContact: outbound.DestinationContact,
		ReferenceNumber:    outbound.ReferenceNumber,
		UnitPrice:          outbound.UnitPrice,
		TotalPrice:         outbound.TotalPrice,
		Notes:              outbound.Notes,
		ShippedDate:        outbound.ShippedDate.Format(time.RFC3339),
		CreatedAt:          outbound.CreatedAt.Format(time.RFC3339),
		CreatedBy:          outbound.CreatedBy.String(),
		CreatedByName:      outbound.User.Name,
	}
}

// GET /outbounds
func (h *OutboundHandler) GetOutbounds(c *gin.Context) {
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

	outbounds, total, err := h.outboundService.GetOutbounds(search, warehouseId, page, limit)
	if err != nil {
		response.ErrorMessageResponse(c, err, 500)
		return
	}

	resp := make([]OutboundResponse, len(outbounds))
	for i, p := range outbounds {
		resp[i] = mapOutboundToResponse(p)
	}

	if isPaginated {
		response.PaginatedResponse(c, "outbounds", resp, total, page, limit)
	} else {
		response.SuccessResponse(c, resp, "Outbounds fetched successfully")
	}
}

// POST /outbounds
func (h *OutboundHandler) CreateOutbound(c *gin.Context) {
	var req struct {
		ProductID          string  `json:"product_id"`
		WarehouseID        string  `json:"warehouse_id"`
		Quantity           int     `json:"quantity"`
		DestinationType    string  `json:"destination_type"`
		DestinationName    string  `json:"destination_name"`
		DestinationContact string  `json:"destination_contact,omitempty"`
		ReferenceNumber    string  `json:"reference_number,omitempty"`
		UnitPrice          float64 `json:"unit_price,omitempty"`
		Notes              string  `json:"notes,omitempty"`
		ShippedDate        string  `json:"shipped_date"`
		CreatedBy          string  `json:"created_by"`
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
	receivedDate, err := time.Parse(time.RFC3339, req.ShippedDate)
	if err != nil {
		response.ErrorMessageResponse(c, err, 400)
		return
	}

	// Create the Outbound model
	outbound := models.Outbound{
		ProductID:          productID,
		WarehouseID:        warehouseID,
		Quantity:           req.Quantity,
		DestinationName:    req.DestinationName,
		DestinationContact: req.DestinationContact,
		ReferenceNumber:    req.ReferenceNumber,
		UnitPrice:          req.UnitPrice,
		TotalPrice:         req.UnitPrice * float64(req.Quantity),
		Notes:              req.Notes,
		ShippedDate:        receivedDate,
		CreatedAt:          time.Now(),
		CreatedBy:          createdBy,
	}

	// Create the outbound record
	createdOutbound, err := h.outboundService.CreateOutbound(outbound)
	if err != nil {
		response.ErrorMessageResponse(c, err, 500)
		return
	}

	// Respond with the created outbound record
	response.SuccessResponse(c, mapOutboundToResponse(createdOutbound), "Outbound created successfully")
}
