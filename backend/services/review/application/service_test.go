package application

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	orderApp "travel-api/services/order/application"
	orderDomain "travel-api/services/order/domain"
	productApp "travel-api/services/product/application"
	"travel-api/services/review/domain"
)

func init() {
	os.Setenv("TRAVEL_DB_PATH", filepath.Join(os.TempDir(), fmt.Sprintf("chinatravel-review-test-%d.db", os.Getpid())))
}

func TestReviewServiceListsSeededSummary(t *testing.T) {
	_ = productApp.NewProductService()
	service := NewReviewService()

	result, err := service.List(101, "en")
	if err != nil {
		t.Fatalf("list reviews: %v", err)
	}
	if result.Summary.Total == 0 || result.Summary.AverageRating <= 0 || len(result.Reviews) == 0 {
		t.Fatalf("expected seeded review summary, got %#v", result)
	}
	for _, review := range result.Reviews {
		if review.Language != "en" {
			t.Fatalf("expected english language filter, got %#v", review)
		}
	}
}

func TestReviewServiceCreatesVerifiedReviewForMatchingOrder(t *testing.T) {
	productService := productApp.NewProductService()
	orderService := orderApp.NewOrderService(productService)
	order, err := orderService.Create("review-user", orderDomain.CreateOrderRequest{ProductID: 101, PackageID: 1011, TravelDate: orderAppTestDate(), Adults: 1})
	if err != nil {
		t.Fatalf("create order for review: %v", err)
	}
	service := NewReviewService()

	review, err := service.Create("review-user", 101, domain.CreateReviewRequest{
		OrderID:  order.ID,
		Rating:   4.7,
		Content:  "Smooth voucher redemption and clear meeting point.",
		Language: "en",
		Scores: domain.Scores{
			Quality:   4.8,
			Service:   4.6,
			Value:     4.5,
			Transport: 4.2,
			Family:    4.4,
		},
	})
	if err != nil {
		t.Fatalf("create review: %v", err)
	}
	if !review.Verified || review.ProductID != 101 || review.OrderID != order.ID || review.Rating != 4.7 {
		t.Fatalf("unexpected created review: %#v", review)
	}
}

func TestReviewServiceRejectsAnonymousOrUnmatchedOrder(t *testing.T) {
	service := NewReviewService()

	_, err := service.Create("", 101, domain.CreateReviewRequest{Rating: 5, Content: "Great"})
	if err != ErrInvalidReviewRequest {
		t.Fatalf("expected invalid anonymous review, got %v", err)
	}
	_, err = service.Create("no-order-user", 101, domain.CreateReviewRequest{Rating: 5, Content: "Great"})
	if err != ErrReviewNotAllowed {
		t.Fatalf("expected review not allowed, got %v", err)
	}
}

func orderAppTestDate() string {
	return time.Now().AddDate(0, 0, 1).Format("2006-01-02")
}
