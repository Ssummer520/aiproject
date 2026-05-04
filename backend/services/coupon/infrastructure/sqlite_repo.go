package infrastructure

import (
	"database/sql"
	"strings"
	"sync"
	"time"

	"travel-api/internal/db"
	"travel-api/services/coupon/domain"
)

type SQLiteCouponRepo struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewSQLiteCouponRepo() *SQLiteCouponRepo {
	database, err := db.Open()
	if err != nil {
		panic(err)
	}
	repo := &SQLiteCouponRepo{db: database}
	repo.seedDemoData()
	return repo
}

func (r *SQLiteCouponRepo) ListActive() ([]domain.Coupon, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rows, err := r.db.Query(`SELECT id, code, name, discount_type, discount_value, min_spend, valid_from, valid_to, usage_limit, status FROM coupons WHERE status = 'active' ORDER BY min_spend ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	coupons := make([]domain.Coupon, 0)
	for rows.Next() {
		coupon, err := scanCoupon(rows)
		if err != nil {
			return nil, err
		}
		coupons = append(coupons, coupon)
	}
	return coupons, rows.Err()
}

func (r *SQLiteCouponRepo) GetByCode(code string) (domain.Coupon, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	coupon, err := scanCoupon(r.db.QueryRow(`SELECT id, code, name, discount_type, discount_value, min_spend, valid_from, valid_to, usage_limit, status FROM coupons WHERE UPPER(code) = UPPER(?) LIMIT 1`, strings.TrimSpace(code)))
	if err == sql.ErrNoRows {
		return domain.Coupon{}, false, nil
	}
	if err != nil {
		return domain.Coupon{}, false, err
	}
	return coupon, true, nil
}

func (r *SQLiteCouponRepo) seedDemoData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	validFrom := now.AddDate(0, -1, 0).Format("2006-01-02")
	validTo := now.AddDate(1, 0, 0).Format("2006-01-02")
	items := []domain.Coupon{
		{ID: 1, Code: "WELCOME80", Name: "New traveller ¥80 off", DiscountType: "amount", DiscountValue: 80, MinSpend: 300, ValidFrom: validFrom, ValidTo: validTo, UsageLimit: 0, Status: "active"},
		{ID: 2, Code: "CHINA10", Name: "10% off China experiences", DiscountType: "percent", DiscountValue: 10, MinSpend: 200, ValidFrom: validFrom, ValidTo: validTo, UsageLimit: 0, Status: "active"},
		{ID: 3, Code: "FAMILY120", Name: "Family trip ¥120 off", DiscountType: "amount", DiscountValue: 120, MinSpend: 800, ValidFrom: validFrom, ValidTo: validTo, UsageLimit: 0, Status: "active"},
	}
	for _, item := range items {
		_, _ = r.db.Exec(`INSERT OR IGNORE INTO coupons(id, code, name, discount_type, discount_value, min_spend, valid_from, valid_to, usage_limit, status) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, item.ID, item.Code, item.Name, item.DiscountType, item.DiscountValue, item.MinSpend, item.ValidFrom, item.ValidTo, item.UsageLimit, item.Status)
	}
}

func scanCoupon(scanner interface {
	Scan(dest ...interface{}) error
}) (domain.Coupon, error) {
	var coupon domain.Coupon
	err := scanner.Scan(&coupon.ID, &coupon.Code, &coupon.Name, &coupon.DiscountType, &coupon.DiscountValue, &coupon.MinSpend, &coupon.ValidFrom, &coupon.ValidTo, &coupon.UsageLimit, &coupon.Status)
	return coupon, err
}
