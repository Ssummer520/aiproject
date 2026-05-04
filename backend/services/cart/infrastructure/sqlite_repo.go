package infrastructure

import (
	"database/sql"
	"sync"
	"time"

	"travel-api/internal/db"
	"travel-api/services/cart/domain"
)

type SQLiteCartRepo struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewSQLiteCartRepo() *SQLiteCartRepo {
	database, err := db.Open()
	if err != nil {
		panic(err)
	}
	return &SQLiteCartRepo{db: database}
}

func (r *SQLiteCartRepo) ListRaw(userID string) ([]domain.CartItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rows, err := r.db.Query(`SELECT id, user_id, product_id, package_id, travel_date, adults, children, quantity, selected_options, created_at FROM cart_items WHERE user_id = ? ORDER BY created_at DESC, id DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.CartItem, 0)
	for rows.Next() {
		var item domain.CartItem
		if err := rows.Scan(&item.ID, &item.UserID, &item.ProductID, &item.PackageID, &item.TravelDate, &item.Adults, &item.Children, &item.Quantity, &item.SelectedOptions, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *SQLiteCartRepo) Add(userID string, item domain.CartItem) (domain.CartItem, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id, err := r.nextID(userID)
	if err != nil {
		return domain.CartItem{}, err
	}
	item.ID = id
	item.UserID = userID
	item.CreatedAt = time.Now().Format(time.RFC3339Nano)
	_, err = r.db.Exec(`INSERT INTO cart_items(id, user_id, product_id, package_id, travel_date, adults, children, quantity, selected_options, created_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, item.ID, userID, item.ProductID, item.PackageID, item.TravelDate, item.Adults, item.Children, item.Quantity, item.SelectedOptions, item.CreatedAt)
	return item, err
}

func (r *SQLiteCartRepo) Clear(userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec(`DELETE FROM cart_items WHERE user_id = ?`, userID)
	return err
}

func (r *SQLiteCartRepo) nextID(userID string) (int, error) {
	var next sql.NullInt64
	if err := r.db.QueryRow(`SELECT MAX(id) + 1 FROM cart_items WHERE user_id = ?`, userID).Scan(&next); err != nil {
		return 0, err
	}
	if !next.Valid || next.Int64 <= 0 {
		return 1, nil
	}
	return int(next.Int64), nil
}
