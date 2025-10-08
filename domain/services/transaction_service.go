package services

import (
	"fmt"
	"wms-be/domain/models"
	"wms-be/domain/repository"
	"wms-be/infrastructure/database"
)

type TransactionService interface {
	CreateTransaction(transaction *models.Transaction) (*models.Transaction, error)
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{
		repo: repo,
	}
}

func (s *transactionService) CreateTransaction(transaction *models.Transaction) (*models.Transaction, error) {
	createdTransaction, err := s.repo.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}
	if transaction.Type == "transfer" {
		err := s.transferStock(transaction)
		if err != nil {
			return nil, fmt.Errorf("failed to execute transfer stock procedure: %v", err)
		}
	}

	return createdTransaction, nil
}

func (s *transactionService) transferStock(transaction *models.Transaction) error {
	db := database.GetDB()

	err := db.Exec(`
    SELECT public.transfer_stock(?, ?, ?, ?, ?, ?, ?);
`, transaction.ProductID, transaction.WarehouseID, transaction.ToWarehouseID, transaction.Quantity,
		transaction.ReferenceNumber, transaction.Notes, transaction.CreatedBy).Error

	if err != nil {
		return fmt.Errorf("error executing transfer stock procedure: %v", err)
	}

	return nil
}
