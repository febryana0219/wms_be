package services

import (
	"errors"
	"wms-be/domain/models"
	"wms-be/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IProductService interface {
	GetProducts(search, warehouseId, category string, page, limit int) ([]models.Product, int, error)
	GetAllProducts() ([]models.Product, error)
	GetProductByID(id string) (models.Product, error)
	CreateProduct(product models.Product) (models.Product, error)
	UpdateProduct(productId string, product models.Product) (models.Product, error)
	DeleteProduct(productId string) error
}

type ProductService struct {
	productRepo repository.ProductRepository
}

// Constructor
func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) GetProducts(search, warehouseId, category string, page, limit int) ([]models.Product, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	products, total, err := s.productRepo.GetProducts(search, warehouseId, category, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	products, _, err := s.productRepo.GetProducts("", "", "", 0, 0)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) CreateProduct(product models.Product) (models.Product, error) {
	existingProduct, err := s.productRepo.GetProductBySKU(product.SKU)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Product{}, err
	}
	if existingProduct.ID != uuid.Nil {
		return models.Product{}, errors.New("product with this SKU already exists")
	}

	createdProduct, err := s.productRepo.CreateProduct(product)
	if err != nil {
		return models.Product{}, err
	}
	return createdProduct, nil
}

func (s *ProductService) UpdateProduct(productId string, product models.Product) (models.Product, error) {
	updatedProduct, err := s.productRepo.UpdateProduct(productId, product)
	if err != nil {
		return models.Product{}, err
	}
	return updatedProduct, nil
}

func (s *ProductService) DeleteProduct(productId string) error {
	return s.productRepo.DeleteProduct(productId)
}

func (s *ProductService) GetProductByID(id string) (models.Product, error) {
	product, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}
