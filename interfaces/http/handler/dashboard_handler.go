package handler

import (
	"net/http"
	"wms-be/domain/services"
	"wms-be/infrastructure/database"
	"wms-be/interfaces/http/response"

	"github.com/gin-gonic/gin"
)

func GetDashboard(c *gin.Context) {
	db := database.GetDB()

	dashboardService := services.NewDashboardService(db)

	data, err := dashboardService.GetDashboardSummary()
	if err != nil {
		response.ErrorMessageResponse(c, err, http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(c, data, "Dashboard summary fetched successfully")
}
