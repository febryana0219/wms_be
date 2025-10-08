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

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

type ProductResponse struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	SKU            string  `json:"sku"`
	Description    string  `json:"description"`
	Price          float64 `json:"price"`
	Stock          int     `json:"stock"`
	ReservedStock  int     `json:"reservedStock"`
	AvailableStock int     `json:"availableStock"`
	WarehouseID    string  `json:"warehouseId"`
	WarehouseName  string  `json:"warehouseName"`
	Category       string  `json:"category"`
	MinStock       int     `json:"minStock"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}

func mapProductToResponse(p models.Product) ProductResponse {
	return ProductResponse{
		ID:             p.ID.String(),
		Name:           p.Name,
		SKU:            p.SKU,
		Description:    p.Description,
		Price:          p.Price,
		Stock:          p.Stock,
		ReservedStock:  p.ReservedStock,
		AvailableStock: p.AvailableStock,
		WarehouseID:    p.WarehouseID.String(),
		WarehouseName:  p.Warehouse.Name,
		Category:       p.Category,
		MinStock:       p.MinStock,
		CreatedAt:      p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      p.UpdatedAt.Format(time.RFC3339),
	}
}

// GET /products
func (h *ProductHandler) GetProducts(c *gin.Context) {
	search := c.Query("search")
	warehouseId := c.Query("warehouseId")
	category := c.Query("category")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	// === cek apakah ada parameter pagination/filter ===
	isPaginated := pageStr != "" || limitStr != "" || search != "" || warehouseId != "" || category != ""

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

	products, total, err := h.productService.GetProducts(search, warehouseId, category, page, limit)
	if err != nil {
		response.ErrorMessageResponse(c, err, 500)
		return
	}

	resp := make([]ProductResponse, len(products))
	for i, p := range products {
		resp[i] = mapProductToResponse(p)
	}

	if isPaginated {
		response.PaginatedResponse(c, "products", resp, total, page, limit)
	} else {
		response.SuccessResponse(c, resp, "Products fetched successfully")
	}
}

// POST /products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	// request struct tanpa AvailableStock
	var req struct {
		SKU           string  `json:"sku" binding:"required"`
		Name          string  `json:"name" binding:"required"`
		Category      string  `json:"category"`
		Description   string  `json:"description"`
		Price         float64 `json:"price"`
		Stock         int     `json:"stock"`
		ReservedStock int     `json:"reservedStock"`
		MinStock      int     `json:"minStock"`
		WarehouseID   string  `json:"warehouseId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorMessageResponse(c, err, 400)
		return
	}

	product := models.Product{
		SKU:           req.SKU,
		Name:          req.Name,
		Category:      req.Category,
		Description:   req.Description,
		Price:         req.Price,
		Stock:         req.Stock,
		ReservedStock: req.ReservedStock,
		MinStock:      req.MinStock,
		WarehouseID:   uuid.MustParse(req.WarehouseID),
	}

	createdProduct, err := h.productService.CreateProduct(product)
	if err != nil {
		response.ErrorMessageResponse(c, err, 500)
		return
	}

	response.SuccessResponse(c, mapProductToResponse(createdProduct), "Product created successfully")
}

// GET /products/:id
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	product, err := h.productService.GetProductByID(id)
	if err != nil {
		response.ErrorMessageResponse(c, err, 500)
		return
	}

	response.SuccessResponse(c, mapProductToResponse(product), "")
}

// PUT /products/:id
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	// request struct tanpa AvailableStock
	var req struct {
		SKU           string  `json:"sku"`
		Name          string  `json:"name"`
		Category      string  `json:"category"`
		Description   string  `json:"description"`
		Price         float64 `json:"price"`
		Stock         int     `json:"stock"`
		ReservedStock int     `json:"reservedStock"`
		MinStock      int     `json:"minStock"`
		WarehouseID   string  `json:"warehouseId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorMessageResponse(c, err, 400)
		return
	}

	product := models.Product{
		SKU:           req.SKU,
		Name:          req.Name,
		Category:      req.Category,
		Description:   req.Description,
		Price:         req.Price,
		Stock:         req.Stock,
		ReservedStock: req.ReservedStock,
		MinStock:      req.MinStock,
	}

	if req.WarehouseID != "" {
		product.WarehouseID = uuid.MustParse(req.WarehouseID)
	}

	updatedProduct, err := h.productService.UpdateProduct(id, product)
	if err != nil {
		response.ErrorMessageResponse(c, err, 500)
		return
	}

	response.SuccessResponse(c, mapProductToResponse(updatedProduct), "Product updated successfully")
}

// DELETE /products/:id
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := h.productService.DeleteProduct(id); err != nil {
		response.ErrorMessageResponse(c, err, 500)
		return
	}

	response.SuccessResponse(c, nil, "Product deleted successfully")
}
