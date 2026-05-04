package infrastructure

import (
	"database/sql"
	"encoding/json"
	"sync"
	"time"

	"travel-api/internal/db"
	"travel-api/services/review/domain"
)

type SQLiteReviewRepo struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewSQLiteReviewRepo() *SQLiteReviewRepo {
	database, err := db.Open()
	if err != nil {
		panic(err)
	}
	repo := &SQLiteReviewRepo{db: database}
	repo.seedDemoData()
	return repo
}

func (r *SQLiteReviewRepo) ListByProduct(productID int, language string) ([]domain.Review, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	query := `SELECT id, user_id, product_id, order_id, rating, quality_score, service_score, value_score, transport_score, family_score, content, images, language, verified, merchant_reply, created_at FROM reviews WHERE product_id = ?`
	args := []interface{}{productID}
	if language != "" {
		query += ` AND language = ?`
		args = append(args, language)
	}
	query += ` ORDER BY verified DESC, created_at DESC, id DESC LIMIT 50`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := make([]domain.Review, 0)
	for rows.Next() {
		review, err := scanReview(rows)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, rows.Err()
}

func (r *SQLiteReviewRepo) Create(review domain.Review) (domain.Review, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if review.CreatedAt == "" {
		review.CreatedAt = time.Now().Format(time.RFC3339Nano)
	}
	var next sql.NullInt64
	if err := r.db.QueryRow(`SELECT MAX(id) + 1 FROM reviews`).Scan(&next); err != nil {
		return domain.Review{}, err
	}
	if next.Valid && next.Int64 > 0 {
		review.ID = int(next.Int64)
	} else {
		review.ID = 1
	}
	images, _ := json.Marshal(review.Images)
	verified := 0
	if review.Verified {
		verified = 1
	}
	_, err := r.db.Exec(`INSERT INTO reviews(id, user_id, product_id, order_id, rating, quality_score, service_score, value_score, transport_score, family_score, content, images, language, verified, merchant_reply, created_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		review.ID, review.UserID, review.ProductID, review.OrderID, review.Rating, review.Scores.Quality, review.Scores.Service, review.Scores.Value, review.Scores.Transport, review.Scores.Family, review.Content, string(images), review.Language, verified, review.MerchantReply, review.CreatedAt)
	if err != nil {
		return domain.Review{}, err
	}
	return review, nil
}

func (r *SQLiteReviewRepo) UserHasProductOrder(userID string, productID int, orderID int) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	query := `SELECT COUNT(1) FROM orders o JOIN order_items i ON i.user_id = o.user_id AND i.order_id = o.id WHERE o.user_id = ? AND i.product_id = ? AND o.status IN ('confirmed', 'paid', 'completed')`
	args := []interface{}{userID, productID}
	if orderID > 0 {
		query += ` AND o.id = ?`
		args = append(args, orderID)
	}
	var count int
	if err := r.db.QueryRow(query, args...).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *SQLiteReviewRepo) seedDemoData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	var count int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM reviews`).Scan(&count); err != nil || count > 0 {
		return
	}
	now := time.Now()
	items := []domain.Review{
		{ID: 1, UserID: "demo", ProductID: 101, Rating: 4.8, Scores: domain.Scores{Quality: 4.8, Service: 4.7, Value: 4.6, Transport: 4.5, Family: 4.6}, Content: "Easy mobile voucher redemption and a beautiful classic West Lake route.", Language: "en", Verified: true, MerchantReply: "Thank you for travelling with us!", CreatedAt: now.AddDate(0, 0, -9).Format(time.RFC3339Nano)},
		{ID: 2, UserID: "demo", ProductID: 103, Rating: 4.7, Scores: domain.Scores{Quality: 4.7, Service: 4.5, Value: 4.3, Transport: 4.2, Family: 4.9}, Content: "Great for families. Entry was smooth with passport and QR voucher.", Language: "en", Verified: true, CreatedAt: now.AddDate(0, 0, -6).Format(time.RFC3339Nano)},
		{ID: 3, UserID: "demo", ProductID: 108, Rating: 4.9, Scores: domain.Scores{Quality: 4.9, Service: 4.9, Value: 4.7, Transport: 4.4, Family: 4.5}, Content: "Local host was warm and the tea tasting felt authentic.", Language: "en", Verified: true, CreatedAt: now.AddDate(0, 0, -4).Format(time.RFC3339Nano)},
		{ID: 4, UserID: "demo", ProductID: 101, Rating: 4.9, Scores: domain.Scores{Quality: 4.9, Service: 4.8, Value: 4.7, Transport: 4.6, Family: 4.7}, Content: "西湖游船很方便，电子凭证兑换很顺畅。", Language: "zh", Verified: true, CreatedAt: now.AddDate(0, 0, -3).Format(time.RFC3339Nano)},
	}
	for _, item := range items {
		images, _ := json.Marshal(item.Images)
		verified := 0
		if item.Verified {
			verified = 1
		}
		_, _ = r.db.Exec(`INSERT OR IGNORE INTO reviews(id, user_id, product_id, order_id, rating, quality_score, service_score, value_score, transport_score, family_score, content, images, language, verified, merchant_reply, created_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, item.ID, item.UserID, item.ProductID, item.OrderID, item.Rating, item.Scores.Quality, item.Scores.Service, item.Scores.Value, item.Scores.Transport, item.Scores.Family, item.Content, string(images), item.Language, verified, item.MerchantReply, item.CreatedAt)
	}
}

func scanReview(scanner interface {
	Scan(dest ...interface{}) error
}) (domain.Review, error) {
	var review domain.Review
	var images string
	var verified int
	err := scanner.Scan(&review.ID, &review.UserID, &review.ProductID, &review.OrderID, &review.Rating, &review.Scores.Quality, &review.Scores.Service, &review.Scores.Value, &review.Scores.Transport, &review.Scores.Family, &review.Content, &images, &review.Language, &verified, &review.MerchantReply, &review.CreatedAt)
	if err != nil {
		return domain.Review{}, err
	}
	_ = json.Unmarshal([]byte(images), &review.Images)
	review.Verified = verified == 1
	return review, nil
}
