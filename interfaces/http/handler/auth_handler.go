package handler

import (
	"net/http"
	"os"
	"strconv"
	"wms-be/domain/models"
	"wms-be/domain/repository"
	"wms-be/domain/services"
	"wms-be/infrastructure/jwt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ===================== Helper =====================
func getAccessExpire() int64 {
	valStr := os.Getenv("JWT_ACCESS_EXPIRE")
	if valStr == "" {
		return 3600 // default 1 jam
	}
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return 3600
	}
	return val
}

func buildUserResponse(user *models.User) gin.H {
	return gin.H{
		"id":            user.ID.String(),
		"email":         user.Email,
		"name":          user.Name,
		"role":          user.Role,
		"warehouseId":   user.WarehouseID.String(),
		"warehouseName": user.Warehouse.Name,
		"createdAt":     user.CreatedAt,
	}
}

// ===================== LOGIN =====================
func Login(c *gin.Context, authService services.AuthService, userRepo repository.UserRepository) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request payload"})
		return
	}

	accessToken, refreshToken, err := authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
		return
	}

	user, err := userRepo.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data": gin.H{
			"user": buildUserResponse(user),
			"token": gin.H{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
				"token_type":    "Bearer",
				"expires_in":    getAccessExpire(),
			},
		},
	})
}

// ===================== REFRESH TOKEN =====================
func RefreshToken(c *gin.Context, authService services.AuthService, userRepo repository.UserRepository) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request payload"})
		return
	}

	accessToken, newRefreshToken, err := authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
		return
	}

	// ambil user dari token lama
	userIDStr, err := jwt.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid refresh token"})
		return
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid user ID format"})
		return
	}

	user, err := userRepo.GetUserByID(userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token refreshed successfully",
		"data": gin.H{
			"user": buildUserResponse(user),
			"token": gin.H{
				"access_token":  accessToken,
				"refresh_token": newRefreshToken,
				"token_type":    "Bearer",
				"expires_in":    getAccessExpire(),
			},
		},
	})
}

// ===================== ME =====================
func Me(c *gin.Context, userRepo repository.UserRepository) {
	userIDStr, err := jwt.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid user ID format"})
		return
	}

	user, err := userRepo.GetUserByID(userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User profile fetched successfully",
		"data": gin.H{
			"user": buildUserResponse(user),
		},
	})
}

// ===================== LOGOUT =====================
func Logout(c *gin.Context, authService services.AuthService) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request payload"})
		return
	}

	if err := authService.Logout(req.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Logout successful"})
}
