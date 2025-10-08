package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Standard success response
func SuccessResponse(c *gin.Context, data any, message string) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// Standard paginated response
func PaginatedResponse(c *gin.Context, key string, data any, total int, page, limit int) {
	totalPages := (total + limit - 1) / limit
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			key:           data,
			"total":       total,
			"pages":       totalPages,
			"currentPage": page,
		},
	})
}

// Standard success message response
func SuccessMessageResponse(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
	})
}

// Standard error message response
func ErrorMessageResponse(c *gin.Context, err error, code int) {
	c.JSON(code, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}
