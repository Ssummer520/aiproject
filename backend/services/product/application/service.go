package application

import (
	"errors"
	"strconv"
	"strings"

	"travel-api/services/product/domain"
	"travel-api/services/product/infrastructure"
)

var ErrProductNotFound = errors.New("product_not_found")

type ProductService struct {
	repo *infrastructure.SQLiteProductRepo
}

func NewProductService() *ProductService {
	return &ProductService{repo: infrastructure.NewSQLiteProductRepo()}
}

func (s *ProductService) Search(filters domain.SearchFilters) (domain.SearchResult, error) {
	products, err := s.repo.Search(filters)
	if err != nil {
		return domain.SearchResult{}, err
	}
	return domain.SearchResult{Results: products, Total: len(products)}, nil
}

func (s *ProductService) Get(id int) (domain.Product, error) {
	product, ok, err := s.repo.Get(id)
	if err != nil {
		return domain.Product{}, err
	}
	if !ok {
		return domain.Product{}, ErrProductNotFound
	}
	return product, nil
}

func (s *ProductService) GetByDestinationID(destinationID int) (domain.Product, error) {
	product, ok, err := s.repo.GetByDestinationID(destinationID)
	if err != nil {
		return domain.Product{}, err
	}
	if !ok {
		return domain.Product{}, ErrProductNotFound
	}
	return product, nil
}

func (s *ProductService) Availability(productID int, date string) ([]domain.Availability, error) {
	return s.repo.ListAvailability(productID, date)
}

func (s *ProductService) Package(packageID int) (domain.Package, error) {
	pkg, ok, err := s.repo.GetPackage(packageID)
	if err != nil {
		return domain.Package{}, err
	}
	if !ok {
		return domain.Package{}, ErrProductNotFound
	}
	return pkg, nil
}

func (s *ProductService) AvailabilityForPackage(packageID int, date string) (domain.Availability, bool, error) {
	return s.repo.GetAvailability(packageID, date)
}

func (s *ProductService) IncrementBookedCount(productID int, quantity int) error {
	return s.repo.IncrementBookedCount(productID, quantity)
}

func ParseBoolPointer(value string) *bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return nil
	}
	return &parsed
}
