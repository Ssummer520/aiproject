package application

import (
	"errors"
	"math"
	"strings"
	"time"

	"travel-api/services/coupon/domain"
	"travel-api/services/coupon/infrastructure"
)

var (
	ErrCouponNotFound = errors.New("coupon_not_found")
	ErrCouponInvalid  = errors.New("coupon_invalid")
)

type CouponService struct {
	repo *infrastructure.SQLiteCouponRepo
}

func NewCouponService() *CouponService {
	return &CouponService{repo: infrastructure.NewSQLiteCouponRepo()}
}

func (s *CouponService) ListActive() ([]domain.Coupon, error) {
	return s.repo.ListActive()
}

func (s *CouponService) Validate(code string, amount float64) (domain.ValidationResult, error) {
	code = strings.TrimSpace(code)
	result := domain.ValidationResult{OriginalAmount: roundMoney(amount), FinalAmount: roundMoney(amount)}
	if code == "" {
		return result, nil
	}
	coupon, ok, err := s.repo.GetByCode(code)
	if err != nil {
		return domain.ValidationResult{}, err
	}
	if !ok {
		result.Error = ErrCouponNotFound.Error()
		return result, ErrCouponNotFound
	}
	if coupon.Status != "active" || !isWithinWindow(coupon) {
		result.Coupon = coupon
		result.Error = ErrCouponInvalid.Error()
		return result, ErrCouponInvalid
	}
	if amount < coupon.MinSpend {
		result.Coupon = coupon
		result.Error = "coupon_min_spend_not_met"
		return result, ErrCouponInvalid
	}

	discount := 0.0
	if coupon.DiscountType == "percent" {
		discount = amount * coupon.DiscountValue / 100
	} else {
		discount = coupon.DiscountValue
	}
	discount = math.Min(discount, amount)
	result.Valid = true
	result.Coupon = coupon
	result.DiscountAmount = roundMoney(discount)
	result.FinalAmount = roundMoney(amount - discount)
	return result, nil
}

func isWithinWindow(coupon domain.Coupon) bool {
	today := time.Now().Format("2006-01-02")
	return (coupon.ValidFrom == "" || coupon.ValidFrom <= today) && (coupon.ValidTo == "" || coupon.ValidTo >= today)
}

func roundMoney(value float64) float64 {
	return math.Round(value*100) / 100
}
