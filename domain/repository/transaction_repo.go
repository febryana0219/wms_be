package repository

import (
	"wms-be/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) (*models.Transaction, error)
	GetTransactionsByWarehouse(warehouseID uuid.UUID, limit, offset int) ([]models.Transaction, error)
	GetTransactionByID(id uuid.UUID) (*models.Transaction, error)
	GetAllTransactions(limit, offset int) ([]models.Transaction, error) // Corrected interface method
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) CreateTransaction(transaction *models.Transaction) (*models.Transaction, error) {
	err := r.db.Create(transaction).Error
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

// Corrected method signature for GetAllTransactions
func (r *transactionRepository) GetAllTransactions(limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Limit(limit).Offset(offset).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) GetTransactionsByWarehouse(warehouseID uuid.UUID, limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("warehouse_id = ?", warehouseID).
		Limit(limit).Offset(offset).
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) GetTransactionByID(id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}
