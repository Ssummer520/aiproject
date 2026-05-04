package infrastructure

import (
	"database/sql"
	"sync"
	"time"

	"travel-api/internal/db"
	"travel-api/services/order/domain"
)

type SQLiteOrderRepo struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewSQLiteOrderRepo() *SQLiteOrderRepo {
	database, err := db.Open()
	if err != nil {
		panic(err)
	}
	return &SQLiteOrderRepo{db: database}
}

func (r *SQLiteOrderRepo) ListByUser(userID string) ([]domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rows, err := r.db.Query(`SELECT id, user_id, status, payment_status, total_amount, currency, contact_name, contact_email, created_at, updated_at, COALESCE(cancelled_at, '') FROM orders WHERE user_id = ? ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]domain.Order, 0)
	for rows.Next() {
		order, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		items, err := r.listItems(userID, order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (r *SQLiteOrderRepo) Create(userID string, order domain.Order, item domain.OrderItem) (domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	orderID, err := r.nextOrderID(userID)
	if err != nil {
		return domain.Order{}, err
	}
	itemID := 1
	now := time.Now().Format(time.RFC3339Nano)
	order.ID = orderID
	order.UserID = userID
	order.CreatedAt = now
	order.UpdatedAt = now
	item.ID = itemID
	item.OrderID = orderID
	item.UserID = userID

	tx, err := r.db.Begin()
	if err != nil {
		return domain.Order{}, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO orders(id, user_id, status, payment_status, total_amount, currency, contact_name, contact_email, created_at, updated_at, cancelled_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NULL)`,
		order.ID,
		order.UserID,
		order.Status,
		order.PaymentStatus,
		order.TotalAmount,
		order.Currency,
		order.ContactName,
		order.ContactEmail,
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		return domain.Order{}, err
	}

	_, err = tx.Exec(`INSERT INTO order_items(id, order_id, user_id, product_id, package_id, product_name, package_name, city, cover, travel_date, adults, children, quantity, unit_price, subtotal) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		item.ID,
		item.OrderID,
		item.UserID,
		item.ProductID,
		item.PackageID,
		item.ProductName,
		item.PackageName,
		item.City,
		item.Cover,
		item.TravelDate,
		item.Adults,
		item.Children,
		item.Quantity,
		item.UnitPrice,
		item.Subtotal,
	)
	if err != nil {
		return domain.Order{}, err
	}
	if err := tx.Commit(); err != nil {
		return domain.Order{}, err
	}

	order.Items = []domain.OrderItem{item}
	return order, nil
}

func (r *SQLiteOrderRepo) Cancel(userID string, orderID int) (domain.Order, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().Format(time.RFC3339Nano)
	result, err := r.db.Exec(`UPDATE orders SET status = 'cancelled', updated_at = ?, cancelled_at = ? WHERE user_id = ? AND id = ? AND status != 'cancelled'`, now, now, userID, orderID)
	if err != nil {
		return domain.Order{}, false, err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		order, ok, err := r.get(userID, orderID)
		return order, ok, err
	}
	order, ok, err := r.get(userID, orderID)
	return order, ok, err
}

func (r *SQLiteOrderRepo) get(userID string, orderID int) (domain.Order, bool, error) {
	order, err := scanOrder(r.db.QueryRow(`SELECT id, user_id, status, payment_status, total_amount, currency, contact_name, contact_email, created_at, updated_at, COALESCE(cancelled_at, '') FROM orders WHERE user_id = ? AND id = ?`, userID, orderID))
	if err == sql.ErrNoRows {
		return domain.Order{}, false, nil
	}
	if err != nil {
		return domain.Order{}, false, err
	}
	items, err := r.listItems(userID, orderID)
	if err != nil {
		return domain.Order{}, false, err
	}
	order.Items = items
	return order, true, nil
}

func (r *SQLiteOrderRepo) listItems(userID string, orderID int) ([]domain.OrderItem, error) {
	rows, err := r.db.Query(`SELECT id, order_id, user_id, product_id, package_id, product_name, package_name, city, cover, travel_date, adults, children, quantity, unit_price, subtotal FROM order_items WHERE user_id = ? AND order_id = ? ORDER BY id ASC`, userID, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.OrderItem, 0)
	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.UserID, &item.ProductID, &item.PackageID, &item.ProductName, &item.PackageName, &item.City, &item.Cover, &item.TravelDate, &item.Adults, &item.Children, &item.Quantity, &item.UnitPrice, &item.Subtotal); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *SQLiteOrderRepo) nextOrderID(userID string) (int, error) {
	var next sql.NullInt64
	if err := r.db.QueryRow(`SELECT MAX(id) + 1 FROM orders WHERE user_id = ?`, userID).Scan(&next); err != nil {
		return 0, err
	}
	if !next.Valid || next.Int64 <= 0 {
		return 1, nil
	}
	return int(next.Int64), nil
}

func scanOrder(scanner interface {
	Scan(dest ...interface{}) error
}) (domain.Order, error) {
	var order domain.Order
	err := scanner.Scan(&order.ID, &order.UserID, &order.Status, &order.PaymentStatus, &order.TotalAmount, &order.Currency, &order.ContactName, &order.ContactEmail, &order.CreatedAt, &order.UpdatedAt, &order.CancelledAt)
	return order, err
}
