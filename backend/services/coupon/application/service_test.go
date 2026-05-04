package application

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func init() {
	os.Setenv("TRAVEL_DB_PATH", filepath.Join(os.TempDir(), fmt.Sprintf("chinatravel-coupon-test-%d.db", os.Getpid())))
}

func TestCouponServiceValidateAmountCoupon(t *testing.T) {
	service := NewCouponService()

	result, err := service.Validate("WELCOME80", 320)
	if err != nil {
		t.Fatalf("validate amount coupon: %v", err)
	}
	if !result.Valid || result.DiscountAmount != 80 || result.FinalAmount != 240 {
		t.Fatalf("unexpected amount coupon result: %#v", result)
	}
}

func TestCouponServiceValidatePercentCoupon(t *testing.T) {
	service := NewCouponService()

	result, err := service.Validate("CHINA10", 500)
	if err != nil {
		t.Fatalf("validate percent coupon: %v", err)
	}
	if !result.Valid || result.DiscountAmount != 50 || result.FinalAmount != 450 {
		t.Fatalf("unexpected percent coupon result: %#v", result)
	}
}

func TestCouponServiceRejectsMinSpendAndMissingCode(t *testing.T) {
	service := NewCouponService()

	if result, err := service.Validate("FAMILY120", 200); err != ErrCouponInvalid || result.Valid {
		t.Fatalf("expected min spend failure, got result=%#v err=%v", result, err)
	}
	if result, err := service.Validate("NOPE", 500); err != ErrCouponNotFound || result.Valid {
		t.Fatalf("expected missing coupon failure, got result=%#v err=%v", result, err)
	}
}
