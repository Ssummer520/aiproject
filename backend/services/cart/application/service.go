package application

import (
	"errors"
	"strings"
	"time"

	"travel-api/services/cart/domain"
	"travel-api/services/cart/infrastructure"
	orderApp "travel-api/services/order/application"
	orderDomain "travel-api/services/order/domain"
	productApp "travel-api/services/product/application"
)

var (
	ErrInvalidCartRequest = errors.New("invalid_cart_request")
	ErrCartEmpty          = errors.New("cart_empty")
)

type CartService struct {
	repo           *infrastructure.SQLiteCartRepo
	productService *productApp.ProductService
	orderService   *orderApp.OrderService
}

func NewCartService(productService *productApp.ProductService, orderService *orderApp.OrderService) *CartService {
	return &CartService{repo: infrastructure.NewSQLiteCartRepo(), productService: productService, orderService: orderService}
}

func (s *CartService) List(userID string) (domain.CartSummary, error) {
	if strings.TrimSpace(userID) == "" {
		return domain.CartSummary{}, ErrInvalidCartRequest
	}
	raw, err := s.repo.ListRaw(userID)
	if err != nil {
		return domain.CartSummary{}, err
	}
	items := make([]domain.CartItem, 0, len(raw))
	total := 0.0
	currency := "CNY"
	for _, item := range raw {
		enriched, err := s.enrich(item)
		if err != nil {
			return domain.CartSummary{}, err
		}
		items = append(items, enriched)
		total += enriched.Subtotal
	}
	return domain.CartSummary{Items: items, TotalAmount: roundMoney(total), Currency: currency}, nil
}

func (s *CartService) Add(userID string, req domain.AddCartItemRequest) (domain.CartSummary, error) {
	if strings.TrimSpace(userID) == "" || req.ProductID <= 0 || req.PackageID <= 0 || strings.TrimSpace(req.TravelDate) == "" {
		return domain.CartSummary{}, ErrInvalidCartRequest
	}
	travelDate, err := time.Parse("2006-01-02", req.TravelDate)
	if err != nil || travelDate.Before(startOfToday()) {
		return domain.CartSummary{}, ErrInvalidCartRequest
	}
	if req.Adults < 0 || req.Children < 0 || req.Adults+req.Children <= 0 {
		return domain.CartSummary{}, ErrInvalidCartRequest
	}
	pkg, err := s.productService.Package(req.PackageID)
	if err != nil || pkg.ProductID != req.ProductID {
		return domain.CartSummary{}, ErrInvalidCartRequest
	}
	quantity := req.Adults + req.Children
	if quantity < pkg.MinQuantity || quantity > pkg.MaxQuantity {
		return domain.CartSummary{}, ErrInvalidCartRequest
	}
	availability, ok, err := s.productService.AvailabilityForPackage(req.PackageID, req.TravelDate)
	if err != nil || !ok || availability.Status != "available" || availability.Stock < quantity {
		return domain.CartSummary{}, ErrInvalidCartRequest
	}
	_, err = s.repo.Add(userID, domain.CartItem{ProductID: req.ProductID, PackageID: req.PackageID, TravelDate: req.TravelDate, Adults: req.Adults, Children: req.Children, Quantity: quantity, SelectedOptions: req.SelectedOptions})
	if err != nil {
		return domain.CartSummary{}, err
	}
	return s.List(userID)
}

func (s *CartService) Clear(userID string) error {
	if strings.TrimSpace(userID) == "" {
		return ErrInvalidCartRequest
	}
	return s.repo.Clear(userID)
}

func (s *CartService) Checkout(userID string, req domain.CheckoutRequest) ([]orderDomain.Order, error) {
	summary, err := s.List(userID)
	if err != nil {
		return nil, err
	}
	if len(summary.Items) == 0 {
		return nil, ErrCartEmpty
	}
	orders := make([]orderDomain.Order, 0, len(summary.Items))
	for index, item := range summary.Items {
		couponCode := ""
		if index == 0 {
			couponCode = req.CouponCode
		}
		order, err := s.orderService.Create(userID, orderDomain.CreateOrderRequest{ProductID: item.ProductID, PackageID: item.PackageID, TravelDate: item.TravelDate, Adults: item.Adults, Children: item.Children, CouponCode: couponCode, ContactName: "Cart Guest", ContactEmail: "cart@example.com"})
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	_ = s.repo.Clear(userID)
	return orders, nil
}

func (s *CartService) enrich(item domain.CartItem) (domain.CartItem, error) {
	product, err := s.productService.Get(item.ProductID)
	if err != nil {
		return domain.CartItem{}, err
	}
	pkg, err := s.productService.Package(item.PackageID)
	if err != nil {
		return domain.CartItem{}, err
	}
	availability, ok, err := s.productService.AvailabilityForPackage(item.PackageID, item.TravelDate)
	if err != nil {
		return domain.CartItem{}, err
	}
	unitPrice := pkg.Price
	if ok {
		unitPrice = availability.Price
	}
	item.ProductName = product.Name
	item.PackageName = pkg.Name
	item.City = product.City
	item.Cover = product.Cover
	item.Quantity = item.Adults + item.Children
	item.UnitPrice = unitPrice
	item.Subtotal = roundMoney(float64(item.Adults)*unitPrice + float64(item.Children)*unitPrice*0.7)
	return item, nil
}

func startOfToday() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func roundMoney(value float64) float64 {
	return float64(int(value*100+0.5)) / 100
}
