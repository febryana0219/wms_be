package handler

import (
	"net/http"
	"strconv"
	"wms-be/domain/services"
	"wms-be/interfaces/http/response"

	"github.com/gin-gonic/gin"
)

type TransactionHistoryHandler struct {
	service *services.TransactionHistoryService
}

func NewTransactionHistoryHandler(service *services.TransactionHistoryService) *TransactionHistoryHandler {
	return &TransactionHistoryHandler{service: service}
}

func (h *TransactionHistoryHandler) GetTransactions(c *gin.Context) {
	txType := c.Query("type")
	warehouseId := c.Query("warehouseId")
	dateFrom := c.Query("dateFrom")
	dateTo := c.Query("dateTo")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	transactions, total, err := h.service.GetTransactions(txType, warehouseId, dateFrom, dateTo, page, limit)
	if err != nil {
		response.ErrorMessageResponse(c, err, http.StatusInternalServerError)
		return
	}

	response.PaginatedResponse(c, "transactions", transactions, int(total), page, limit)
}
