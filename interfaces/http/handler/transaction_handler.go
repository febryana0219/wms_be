package handler

import (
	"net/http"
	"wms-be/domain/models"
	"wms-be/domain/services"
	"wms-be/interfaces/http/response"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service services.TransactionService
}

func NewTransactionHandler(service services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		response.ErrorMessageResponse(c, err, http.StatusBadRequest)
		return
	}

	createdTransaction, err := h.service.CreateTransaction(&transaction)
	if err != nil {
		response.ErrorMessageResponse(c, err, http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, createdTransaction, "Transaction created successfully")
}
