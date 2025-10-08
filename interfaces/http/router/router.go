package router

import (
	"time"
	"wms-be/domain/repository"
	"wms-be/domain/services"
	"wms-be/infrastructure/database"
	"wms-be/infrastructure/middleware"
	"wms-be/interfaces/http/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	warehouseRepo repository.WarehouseRepository,
	productRepo repository.ProductRepository,
	transactionRepo repository.TransactionRepository,
	inboundRepo repository.InboundRepository,
	outboundRepo repository.OutboundRepository,
	orderRepo repository.OrderRepository,
) *gin.Engine {
	r := gin.Default()

	// CORS Middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize Services
	authService := services.NewAuthService(userRepo, refreshTokenRepo)
	warehouseService := services.NewWarehouseService(warehouseRepo)
	productService := services.NewProductService(productRepo)
	transactionService := services.NewTransactionService(transactionRepo)
	inboundService := services.NewInboundService(inboundRepo)
	outboundService := services.NewOutboundService(outboundRepo)
	orderService := services.NewOrderService(orderRepo)

	// Ambil koneksi DB dari package database
	db := database.GetDB()
	transactionHistoryService := services.NewTransactionHistoryService(db)

	// Initialize Handlers
	warehouseHandler := handler.NewWarehouseHandler(warehouseService)
	productHandler := handler.NewProductHandler(productService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	transactionHistoryHandler := handler.NewTransactionHistoryHandler(transactionHistoryService)
	inboundHandler := handler.NewInboundHandler(inboundService)
	outboundHandler := handler.NewOutboundHandler(outboundService)
	orderHandler := handler.NewOrderHandler(orderService)

	// API Prefix
	api := r.Group("/api")

	// Public Auth Routes
	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/login", func(c *gin.Context) {
			handler.Login(c, authService, userRepo)
		})
		authRoutes.POST("/refresh_token", func(c *gin.Context) {
			handler.RefreshToken(c, authService, userRepo)
		})
	}

	// Protected Auth Routes
	authProtected := api.Group("/auth").Use(middleware.AuthMiddleware())
	{
		authProtected.GET("/me", func(c *gin.Context) {
			handler.Me(c, userRepo)
		})
		authProtected.GET("/profile", func(c *gin.Context) {
			handler.Me(c, userRepo)
		})
		authProtected.POST("/logout", func(c *gin.Context) {
			handler.Logout(c, authService)
		})
	}

	// Warehouse Routes
	warehouseRoutes := api.Group("/warehouses").Use(middleware.AuthMiddleware())
	{
		warehouseRoutes.GET("", warehouseHandler.GetWarehouses)
		warehouseRoutes.POST("", warehouseHandler.CreateWarehouse)
		warehouseRoutes.GET("/:id", warehouseHandler.GetWarehouseByID)
		warehouseRoutes.PUT("/:id", warehouseHandler.UpdateWarehouse)
		warehouseRoutes.DELETE("/:id", warehouseHandler.DeleteWarehouse)
	}

	// Product Routes
	productRoutes := api.Group("/products").Use(middleware.AuthMiddleware())
	{
		productRoutes.GET("", productHandler.GetProducts)
		productRoutes.POST("", productHandler.CreateProduct)
		productRoutes.GET("/:id", productHandler.GetProductByID)
		productRoutes.PUT("/:id", productHandler.UpdateProduct)
		productRoutes.DELETE("/:id", productHandler.DeleteProduct)
	}

	// Transaction Routes
	transactionRoutes := api.Group("/transactions").Use(middleware.AuthMiddleware())
	{
		transactionRoutes.POST("", transactionHandler.CreateTransaction)
		transactionRoutes.GET("", transactionHistoryHandler.GetTransactions)
	}

	// Inbound Routes
	inboundRoutes := api.Group("/inbounds").Use(middleware.AuthMiddleware())
	{
		inboundRoutes.GET("", inboundHandler.GetInbounds)
		inboundRoutes.POST("", inboundHandler.CreateInbound)
	}

	// Outbound Routes
	outboundRoutes := api.Group("/outbounds").Use(middleware.AuthMiddleware())
	{
		outboundRoutes.GET("", outboundHandler.GetOutbounds)
		outboundRoutes.POST("", outboundHandler.CreateOutbound)
	}

	// Order Routes
	orderRoutes := api.Group("/orders").Use(middleware.AuthMiddleware())
	{
		orderRoutes.POST("", orderHandler.CreateOrder)
		orderRoutes.GET("", orderHandler.GetOrders)
		orderRoutes.PUT("/:id/status", orderHandler.UpdateOrderStatus)
		orderRoutes.GET("/:id", orderHandler.GetOrderByID)
	}

	// Dashboard Routes
	dashboardRoutes := api.Group("/dashboard").Use(middleware.AuthMiddleware())
	{
		dashboardRoutes.GET("/stats", handler.GetDashboard)
	}

	return r
}
