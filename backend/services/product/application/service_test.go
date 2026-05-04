package application

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"travel-api/services/product/domain"
)

func init() {
	os.Setenv("TRAVEL_DB_PATH", filepath.Join(os.TempDir(), fmt.Sprintf("chinatravel-product-test-%d.db", os.Getpid())))
}

func TestProductServiceSearchAndDetailSeededData(t *testing.T) {
	service := NewProductService()

	result, err := service.Search(domain.SearchFilters{City: "Hangzhou", Category: "Tickets"})
	if err != nil {
		t.Fatalf("search products: %v", err)
	}
	if result.Total == 0 {
		t.Fatal("expected seeded Hangzhou ticket products")
	}

	product, err := service.Get(101)
	if err != nil {
		t.Fatalf("get product detail: %v", err)
	}
	if product.ID != 101 || product.Name == "" {
		t.Fatalf("unexpected product detail: %#v", product)
	}
	if len(product.Packages) != 2 {
		t.Fatalf("expected 2 seeded packages, got %d", len(product.Packages))
	}
	if len(product.Availability) == 0 {
		t.Fatal("expected seeded availability")
	}
}

func TestProductServiceFilters(t *testing.T) {
	service := NewProductService()

	result, err := service.Search(domain.SearchFilters{})
	if err != nil {
		t.Fatalf("search all products: %v", err)
	}
	if result.Total != 8 {
		t.Fatalf("expected 8 seeded products, got %d", result.Total)
	}

	transport, err := service.Search(domain.SearchFilters{Type: "transport"})
	if err != nil {
		t.Fatalf("search transport: %v", err)
	}
	if transport.Total != 1 || transport.Results[0].Type != "transport" {
		t.Fatalf("expected only transport product, got %#v", transport.Results)
	}
}

func TestProductServiceGetByDestinationID(t *testing.T) {
	service := NewProductService()

	product, err := service.GetByDestinationID(1)
	if err != nil {
		t.Fatalf("get product by destination: %v", err)
	}
	if product.ID != 101 || product.DestinationID != 1 {
		t.Fatalf("unexpected destination mapping: %#v", product)
	}
	if _, err := service.GetByDestinationID(999999); err != ErrProductNotFound {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}
