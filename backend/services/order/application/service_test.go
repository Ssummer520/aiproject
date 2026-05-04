package application

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	couponApp "travel-api/services/coupon/application"
	"travel-api/services/order/domain"
	productApp "travel-api/services/product/application"
)

func init() {
	os.Setenv("TRAVEL_DB_PATH", filepath.Join(os.TempDir(), fmt.Sprintf("chinatravel-order-test-%d.db", os.Getpid())))
}

func newTestOrderService() *OrderService {
	return NewOrderService(productApp.NewProductService())
}

func nextTravelDate() string {
	return time.Now().AddDate(0, 0, 1).Format("2006-01-02")
}

func TestOrderServiceCreatePersistsPricedProductOrder(t *testing.T) {
	service := newTestOrderService()

	order, err := service.Create("user-1", domain.CreateOrderRequest{
		ProductID:    101,
		PackageID:    1011,
		TravelDate:   nextTravelDate(),
		Adults:       2,
		Children:     1,
		ContactName:  "Alan",
		ContactEmail: "alan@example.com",
	})
	if err != nil {
		t.Fatalf("create order: %v", err)
	}
	if order.ID == 0 || order.Status != "paid" || order.PaymentStatus != "paid_mock" {
		t.Fatalf("unexpected order status: %#v", order)
	}
	if len(order.Items) != 1 {
		t.Fatalf("expected one order item, got %d", len(order.Items))
	}
	if order.Items[0].ProductID != 101 || order.Items[0].PackageID != 1011 {
		t.Fatalf("unexpected order item: %#v", order.Items[0])
	}
	if order.Items[0].Usage == "" {
		t.Fatalf("expected order item usage instructions")
	}
	if order.TotalAmount <= 0 || order.TotalAmount != order.Items[0].Subtotal {
		t.Fatalf("unexpected totals: order=%v item=%v", order.TotalAmount, order.Items[0].Subtotal)
	}

	orders, err := service.List("user-1")
	if err != nil {
		t.Fatalf("list orders: %v", err)
	}
	if len(orders) != 1 || len(orders[0].Items) != 1 {
		t.Fatalf("expected persisted order with item, got %#v", orders)
	}
	if orders[0].Items[0].Usage == "" {
		t.Fatalf("expected persisted usage instructions")
	}
}

func TestOrderServiceCreateRejectsInvalidRequests(t *testing.T) {
	service := newTestOrderService()

	cases := []struct {
		name string
		req  domain.CreateOrderRequest
		err  error
	}{
		{
			name: "empty user still rejected by service boundary",
			req:  domain.CreateOrderRequest{ProductID: 101, PackageID: 1011, TravelDate: nextTravelDate(), Adults: 1},
			err:  ErrInvalidOrderRequest,
		},
		{
			name: "zero travelers",
			req:  domain.CreateOrderRequest{ProductID: 101, PackageID: 1011, TravelDate: nextTravelDate()},
			err:  ErrInvalidOrderRequest,
		},
		{
			name: "past date",
			req:  domain.CreateOrderRequest{ProductID: 101, PackageID: 1011, TravelDate: time.Now().AddDate(0, 0, -1).Format("2006-01-02"), Adults: 1},
			err:  ErrInvalidOrderRequest,
		},
		{
			name: "package does not belong to product",
			req:  domain.CreateOrderRequest{ProductID: 101, PackageID: 1021, TravelDate: nextTravelDate(), Adults: 1},
			err:  ErrPackageNotFound,
		},
		{
			name: "too many travelers",
			req:  domain.CreateOrderRequest{ProductID: 101, PackageID: 1011, TravelDate: nextTravelDate(), Adults: 10},
			err:  ErrInvalidOrderRequest,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			userID := "user-invalid"
			if tc.name == "empty user still rejected by service boundary" {
				userID = ""
			}
			_, err := service.Create(userID, tc.req)
			if err != tc.err {
				t.Fatalf("expected %v, got %v", tc.err, err)
			}
		})
	}
}

func TestOrderServiceCancelOnlyUserOrder(t *testing.T) {
	service := newTestOrderService()
	order, err := service.Create("user-cancel", domain.CreateOrderRequest{ProductID: 104, PackageID: 1041, TravelDate: nextTravelDate(), Adults: 1})
	if err != nil {
		t.Fatalf("create order: %v", err)
	}

	if _, ok, err := service.Cancel("other-user", order.ID); err != nil || ok {
		t.Fatalf("other user should not cancel order: ok=%v err=%v", ok, err)
	}

	cancelled, ok, err := service.Cancel("user-cancel", order.ID)
	if err != nil || !ok {
		t.Fatalf("cancel own order: ok=%v err=%v", ok, err)
	}
	if cancelled.Status != "cancelled" || cancelled.CancelledAt == "" {
		t.Fatalf("expected cancelled order with timestamp, got %#v", cancelled)
	}
}

func TestOrderServiceCreateAppliesCouponDiscount(t *testing.T) {
	service := NewOrderServiceWithCoupon(productApp.NewProductService(), couponApp.NewCouponService())

	order, err := service.Create("user-coupon", domain.CreateOrderRequest{
		ProductID:   101,
		PackageID:   1011,
		TravelDate:  nextTravelDate(),
		Adults:      4,
		Children:    0,
		CouponCode:  "WELCOME80",
		ContactName: "Coupon User",
	})
	if err != nil {
		t.Fatalf("create order with coupon: %v", err)
	}
	if order.OriginalAmount <= 0 || order.DiscountAmount != 80 || order.TotalAmount != order.OriginalAmount-order.DiscountAmount {
		t.Fatalf("unexpected coupon totals: %#v", order)
	}
	if order.CouponCode != "WELCOME80" {
		t.Fatalf("expected normalized coupon code, got %q", order.CouponCode)
	}

	orders, err := service.List("user-coupon")
	if err != nil {
		t.Fatalf("list coupon orders: %v", err)
	}
	if len(orders) != 1 || orders[0].DiscountAmount != 80 || orders[0].CouponCode != "WELCOME80" {
		t.Fatalf("expected persisted coupon totals, got %#v", orders)
	}
}

func TestOrderServiceRejectsInvalidCoupon(t *testing.T) {
	service := NewOrderServiceWithCoupon(productApp.NewProductService(), couponApp.NewCouponService())

	_, err := service.Create("user-bad-coupon", domain.CreateOrderRequest{
		ProductID:   101,
		PackageID:   1011,
		TravelDate:  nextTravelDate(),
		Adults:      1,
		CouponCode:  "MISSING",
		ContactName: "Bad Coupon",
	})
	if err != ErrInvalidOrderRequest {
		t.Fatalf("expected invalid order request for bad coupon, got %v", err)
	}
}

func TestOrderServiceStatusTransitions(t *testing.T) {
	service := newTestOrderService()
	order, err := service.Create("user-status", domain.CreateOrderRequest{ProductID: 108, PackageID: 1081, TravelDate: nextTravelDate(), Adults: 1})
	if err != nil {
		t.Fatalf("create status order: %v", err)
	}

	completed, ok, err := service.Complete("user-status", order.ID)
	if err != nil || !ok {
		t.Fatalf("complete order: ok=%v err=%v", ok, err)
	}
	if completed.Status != "completed" || completed.PaymentStatus != "paid_mock" {
		t.Fatalf("expected completed paid_mock order, got %#v", completed)
	}

	refunded, ok, err := service.Refund("user-status", order.ID)
	if err != nil || !ok {
		t.Fatalf("refund order: ok=%v err=%v", ok, err)
	}
	if refunded.Status != "refunded" || refunded.PaymentStatus != "refunded" {
		t.Fatalf("expected refunded order, got %#v", refunded)
	}

	cancelled, ok, err := service.Cancel("user-status", order.ID)
	if err != nil || !ok {
		t.Fatalf("cancel refunded order should return existing order: ok=%v err=%v", ok, err)
	}
	if cancelled.Status != "refunded" {
		t.Fatalf("refunded order should not become cancelled, got %#v", cancelled)
	}
}
