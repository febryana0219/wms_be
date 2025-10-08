package repository

import (
	"wms-be/domain/models"
	"wms-be/infrastructure/database" // Importing the database package

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetProducts(search, warehouseId, category string, page, limit int) ([]models.Product, int, error)
	GetProductBySKU(sku string) (models.Product, error)
	GetProductByID(id string) (models.Product, error)
	CreateProduct(product models.Product) (models.Product, error)
	UpdateProduct(productId string, product models.Product) (models.Product, error)
	DeleteProduct(productId string) error
}

type productRepo struct {
	db *gorm.DB
}

func NewProductRepository() ProductRepository {
	return &productRepo{db: database.GetDB()}
}

func (r *productRepo) GetProducts(search, warehouseId, category string, page, limit int) ([]models.Product, int, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{}).Where("is_active = true").Preload("Warehouse")

	if search != "" {
		query = query.Where("name ILIKE ? OR sku ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if warehouseId != "" {
		query = query.Where("warehouse_id = ?", warehouseId)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page - 1) * limit).Limit(limit).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, int(total), nil
}

func (r *productRepo) GetProductBySKU(sku string) (models.Product, error) {
	var product models.Product
	err := r.db.Where("sku = ?", sku).Preload("Warehouse").First(&product).Error
	return product, err
}

func (r *productRepo) GetProductByID(id string) (models.Product, error) {
	var product models.Product
	err := r.db.Where("id = ?", id).Preload("Warehouse").First(&product).Error
	return product, err
}

func (r *productRepo) CreateProduct(product models.Product) (models.Product, error) {
	err := r.db.Create(&product).Error
	return product, err
}

func (r *productRepo) UpdateProduct(productId string, product models.Product) (models.Product, error) {
	var existingProduct models.Product
	err := r.db.Where("id = ?", productId).First(&existingProduct).Error
	if err != nil {
		return models.Product{}, err
	}

	// Update semua field yang dikirim
	if product.Name != "" {
		existingProduct.Name = product.Name
	}
	if product.SKU != "" {
		existingProduct.SKU = product.SKU
	}
	if product.Description != "" {
		existingProduct.Description = product.Description
	}
	if product.Category != "" {
		existingProduct.Category = product.Category
	}
	if product.Price != 0 {
		existingProduct.Price = product.Price
	}
	// stock dan reservedStock
	existingProduct.Stock = product.Stock
	existingProduct.ReservedStock = product.ReservedStock
	existingProduct.MinStock = product.MinStock

	// warehouse
	if product.WarehouseID != uuid.Nil {
		existingProduct.WarehouseID = product.WarehouseID
	}

	// simpan
	err = r.db.Save(&existingProduct).Error
	if err != nil {
		return models.Product{}, err
	}

	// preload warehouse
	r.db.Preload("Warehouse").First(&existingProduct, "id = ?", productId)
	return existingProduct, nil
}

func (r *productRepo) DeleteProduct(productId string) error {
	err := r.db.Where("id = ?", productId).Delete(&models.Product{}).Error
	return err
}
