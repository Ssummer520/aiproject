package application

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"travel-api/services/cart/domain"
	couponApp "travel-api/services/coupon/application"
	orderApp "travel-api/services/order/application"
	productApp "travel-api/services/product/application"
)

func init() {
	os.Setenv("TRAVEL_DB_PATH", filepath.Join(os.TempDir(), fmt.Sprintf("chinatravel-cart-test-%d.db", os.Getpid())))
}

func cartDate() string { return time.Now().AddDate(0, 0, 1).Format("2006-01-02") }

func newCartService() *CartService {
	products := productApp.NewProductService()
	orders := orderApp.NewOrderServiceWithCoupon(products, couponApp.NewCouponService())
	return NewCartService(products, orders)
}

func TestCartServiceAddListAndClear(t *testing.T) {
	service := newCartService()
	summary, err := service.Add("cart-user", domain.AddCartItemRequest{ProductID: 101, PackageID: 1011, TravelDate: cartDate(), Adults: 2})
	if err != nil {
		t.Fatalf("add cart: %v", err)
	}
	if len(summary.Items) != 1 || summary.TotalAmount <= 0 || summary.Items[0].ProductName == "" {
		t.Fatalf("unexpected summary: %#v", summary)
	}
	if err := service.Clear("cart-user"); err != nil {
		t.Fatalf("clear cart: %v", err)
	}
	summary, err = service.List("cart-user")
	if err != nil || len(summary.Items) != 0 {
		t.Fatalf("expected empty cart, got %#v err=%v", summary, err)
	}
}

func TestCartServiceCheckoutCreatesOrdersAndClearsCart(t *testing.T) {
	service := newCartService()
	_, err := service.Add("checkout-user", domain.AddCartItemRequest{ProductID: 101, PackageID: 1011, TravelDate: cartDate(), Adults: 4})
	if err != nil {
		t.Fatalf("add cart: %v", err)
	}
	orders, err := service.Checkout("checkout-user", domain.CheckoutRequest{CouponCode: "WELCOME80"})
	if err != nil {
		t.Fatalf("checkout: %v", err)
	}
	if len(orders) != 1 || orders[0].DiscountAmount != 80 || orders[0].Status != "paid" {
		t.Fatalf("unexpected orders: %#v", orders)
	}
	summary, _ := service.List("checkout-user")
	if len(summary.Items) != 0 {
		t.Fatalf("cart should be cleared after checkout: %#v", summary)
	}
}

func TestCartServiceRejectsInvalidAndEmptyCheckout(t *testing.T) {
	service := newCartService()
	if _, err := service.Add("cart-invalid", domain.AddCartItemRequest{ProductID: 101}); err != ErrInvalidCartRequest {
		t.Fatalf("expected invalid cart request, got %v", err)
	}
	if _, err := service.Checkout("empty-cart", domain.CheckoutRequest{}); err != ErrCartEmpty {
		t.Fatalf("expected empty cart, got %v", err)
	}
}
