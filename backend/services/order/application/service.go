package application

import (
	"errors"
	"strings"
	"time"

	couponApp "travel-api/services/coupon/application"
	"travel-api/services/order/domain"
	"travel-api/services/order/infrastructure"
	productApp "travel-api/services/product/application"
)

var (
	ErrInvalidOrderRequest = errors.New("invalid_order_request")
	ErrProductNotFound     = errors.New("product_not_found")
	ErrPackageNotFound     = errors.New("package_not_found")
	ErrAvailabilityClosed  = errors.New("availability_closed")
)

type OrderService struct {
	repo           *infrastructure.SQLiteOrderRepo
	productService *productApp.ProductService
	couponService  *couponApp.CouponService
}

func NewOrderService(productService *productApp.ProductService) *OrderService {
	return &OrderService{repo: infrastructure.NewSQLiteOrderRepo(), productService: productService, couponService: couponApp.NewCouponService()}
}

func NewOrderServiceWithCoupon(productService *productApp.ProductService, couponService *couponApp.CouponService) *OrderService {
	return &OrderService{repo: infrastructure.NewSQLiteOrderRepo(), productService: productService, couponService: couponService}
}

func (s *OrderService) List(userID string) ([]domain.Order, error) {
	return s.repo.ListByUser(userID)
}

func (s *OrderService) Create(userID string, req domain.CreateOrderRequest) (domain.Order, error) {
	if strings.TrimSpace(userID) == "" {
		return domain.Order{}, ErrInvalidOrderRequest
	}
	if req.ProductID <= 0 || req.PackageID <= 0 || strings.TrimSpace(req.TravelDate) == "" {
		return domain.Order{}, ErrInvalidOrderRequest
	}
	if req.Adults < 0 || req.Children < 0 || req.Adults+req.Children <= 0 {
		return domain.Order{}, ErrInvalidOrderRequest
	}
	travelDate, err := time.Parse("2006-01-02", req.TravelDate)
	if err != nil || travelDate.Before(startOfToday()) {
		return domain.Order{}, ErrInvalidOrderRequest
	}

	product, err := s.productService.Get(req.ProductID)
	if err == productApp.ErrProductNotFound {
		return domain.Order{}, ErrProductNotFound
	}
	if err != nil {
		return domain.Order{}, err
	}
	pkg, err := s.productService.Package(req.PackageID)
	if err == productApp.ErrProductNotFound {
		return domain.Order{}, ErrPackageNotFound
	}
	if err != nil {
		return domain.Order{}, err
	}
	if pkg.ProductID != product.ID {
		return domain.Order{}, ErrPackageNotFound
	}

	quantity := req.Adults + req.Children
	if quantity < pkg.MinQuantity || quantity > pkg.MaxQuantity {
		return domain.Order{}, ErrInvalidOrderRequest
	}

	availability, ok, err := s.productService.AvailabilityForPackage(pkg.ID, req.TravelDate)
	if err != nil {
		return domain.Order{}, err
	}
	if !ok || availability.Status != "available" || availability.Stock < quantity {
		return domain.Order{}, ErrAvailabilityClosed
	}

	unitPrice := availability.Price
	childrenDiscount := unitPrice * 0.7
	subtotal := roundMoney(float64(req.Adults)*unitPrice + float64(req.Children)*childrenDiscount)
	discountAmount := 0.0
	couponCode := strings.TrimSpace(req.CouponCode)
	if couponCode != "" {
		validation, err := s.couponService.Validate(couponCode, subtotal)
		if err != nil || !validation.Valid {
			return domain.Order{}, ErrInvalidOrderRequest
		}
		discountAmount = validation.DiscountAmount
		couponCode = validation.Coupon.Code
	}
	order := domain.Order{
		Status:         "paid",
		PaymentStatus:  "paid_mock",
		OriginalAmount: subtotal,
		DiscountAmount: discountAmount,
		TotalAmount:    roundMoney(subtotal - discountAmount),
		CouponCode:     couponCode,
		Currency:       product.Currency,
		ContactName:    strings.TrimSpace(req.ContactName),
		ContactEmail:   strings.TrimSpace(req.ContactEmail),
	}
	if order.ContactName == "" {
		order.ContactName = "Guest"
	}
	if order.ContactEmail == "" {
		order.ContactEmail = "guest@example.com"
	}

	item := domain.OrderItem{
		ProductID:   product.ID,
		PackageID:   pkg.ID,
		ProductName: product.Name,
		PackageName: pkg.Name,
		City:        product.City,
		Cover:       product.Cover,
		Usage:       product.Usage,
		TravelDate:  req.TravelDate,
		Adults:      req.Adults,
		Children:    req.Children,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		Subtotal:    subtotal,
	}

	created, err := s.repo.Create(userID, order, item)
	if err != nil {
		return domain.Order{}, err
	}
	_ = s.productService.IncrementBookedCount(product.ID, quantity)
	return created, nil
}

func (s *OrderService) Cancel(userID string, orderID int) (domain.Order, bool, error) {
	return s.repo.Cancel(userID, orderID)
}

func startOfToday() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func roundMoney(value float64) float64 {
	return float64(int(value*100+0.5)) / 100
}

func (s *OrderService) Complete(userID string, orderID int) (domain.Order, bool, error) {
	return s.repo.UpdateStatus(userID, orderID, "completed", "paid_mock")
}

func (s *OrderService) Refund(userID string, orderID int) (domain.Order, bool, error) {
	return s.repo.UpdateStatus(userID, orderID, "refunded", "refunded")
}
