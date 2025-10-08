package main

import (
	"wms-be/config"
	"wms-be/domain/repository"
	"wms-be/infrastructure/database"
	"wms-be/infrastructure/jwt"
	"wms-be/interfaces/http/router"
)

func main() {
	// Initialize configurations
	config.InitConfig()

	// Initialize JWT
	jwt.Init()

	// Initialize database
	database.InitDB()

	// Initialize repositories
	userRepo := repository.NewUserRepository()
	refreshTokenRepo := repository.NewRefreshTokenRepository()
	warehouseRepo := repository.NewWarehouseRepository()
	productRepo := repository.NewProductRepository()
	transactionRepo := repository.NewTransactionRepository(database.GetDB())
	inboundRepo := repository.NewInboundRepository()
	outboundRepo := repository.NewOutboundRepository()
	orderRepo := repository.NewOrderRepository(database.GetDB())

	// Setup router with all repositories (including Inbound Repository)
	r := router.SetupRouter(
		userRepo,
		refreshTokenRepo,
		warehouseRepo,
		productRepo,
		transactionRepo,
		inboundRepo,
		outboundRepo,
		orderRepo,
	)

	// Run the server on port 8000
	r.Run(":8000")
}
