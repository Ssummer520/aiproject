package infrastructure

import (
	"database/sql"
	"sync"
	"time"

	"travel-api/internal/db"
	"travel-api/services/itinerary/domain"
)

type SQLiteItineraryRepo struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewSQLiteItineraryRepo() *SQLiteItineraryRepo {
	database, err := db.Open()
	if err != nil {
		panic(err)
	}
	return &SQLiteItineraryRepo{db: database}
}

func (r *SQLiteItineraryRepo) List(userID string) ([]domain.Itinerary, error) {
	r.mu.RLock()
	rows, err := r.db.Query(`SELECT id, user_id, title, city, start_date, end_date, guests, budget, status, created_at, updated_at FROM itineraries WHERE user_id = ? ORDER BY start_date ASC, id DESC`, userID)
	if err != nil {
		r.mu.RUnlock()
		return nil, err
	}
	items := make([]domain.Itinerary, 0)
	for rows.Next() {
		item, err := scanItinerary(rows)
		if err != nil {
			_ = rows.Close()
			r.mu.RUnlock()
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		r.mu.RUnlock()
		return nil, err
	}
	_ = rows.Close()
	r.mu.RUnlock()

	for i := range items {
		children, err := r.ListItems(userID, items[i].ID)
		if err != nil {
			return nil, err
		}
		items[i].Items = children
	}
	return items, nil
}

func (r *SQLiteItineraryRepo) Get(userID string, itineraryID int) (domain.Itinerary, bool, error) {
	r.mu.RLock()
	itinerary, err := scanItinerary(r.db.QueryRow(`SELECT id, user_id, title, city, start_date, end_date, guests, budget, status, created_at, updated_at FROM itineraries WHERE user_id = ? AND id = ?`, userID, itineraryID))
	r.mu.RUnlock()
	if err == sql.ErrNoRows {
		return domain.Itinerary{}, false, nil
	}
	if err != nil {
		return domain.Itinerary{}, false, err
	}
	items, err := r.ListItems(userID, itinerary.ID)
	if err != nil {
		return domain.Itinerary{}, false, err
	}
	itinerary.Items = items
	return itinerary, true, nil
}

func (r *SQLiteItineraryRepo) Create(userID string, itinerary domain.Itinerary) (domain.Itinerary, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id, err := r.nextItineraryID(userID)
	if err != nil {
		return domain.Itinerary{}, err
	}
	now := time.Now().Format(time.RFC3339Nano)
	itinerary.ID = id
	itinerary.UserID = userID
	itinerary.CreatedAt = now
	itinerary.UpdatedAt = now
	if itinerary.Status == "" {
		itinerary.Status = "draft"
	}
	_, err = r.db.Exec(`INSERT INTO itineraries(id, user_id, title, city, start_date, end_date, guests, budget, status, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, itinerary.ID, userID, itinerary.Title, itinerary.City, itinerary.StartDate, itinerary.EndDate, itinerary.Guests, itinerary.Budget, itinerary.Status, itinerary.CreatedAt, itinerary.UpdatedAt)
	if err != nil {
		return domain.Itinerary{}, err
	}
	for _, item := range itinerary.Items {
		if _, err := r.addItemLocked(userID, itinerary.ID, item); err != nil {
			return domain.Itinerary{}, err
		}
	}
	created, _, err := r.getLocked(userID, itinerary.ID)
	return created, err
}

func (r *SQLiteItineraryRepo) AddItem(userID string, itineraryID int, item domain.ItineraryItem) (domain.Itinerary, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok, err := r.getLocked(userID, itineraryID); err != nil || !ok {
		return domain.Itinerary{}, ok, err
	}
	if _, err := r.addItemLocked(userID, itineraryID, item); err != nil {
		return domain.Itinerary{}, false, err
	}
	itinerary, ok, err := r.getLocked(userID, itineraryID)
	return itinerary, ok, err
}

func (r *SQLiteItineraryRepo) MoveItem(userID string, itineraryID int, itemID int, direction string) (domain.Itinerary, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	items, err := r.listItemsLocked(userID, itineraryID)
	if err != nil || len(items) == 0 {
		return domain.Itinerary{}, false, err
	}
	index := -1
	for i, item := range items {
		if item.ID == itemID {
			index = i
			break
		}
	}
	if index < 0 {
		return domain.Itinerary{}, false, nil
	}
	target := index
	if direction == "up" && index > 0 {
		target = index - 1
	}
	if direction == "down" && index < len(items)-1 {
		target = index + 1
	}
	if target != index {
		items[index].SortOrder, items[target].SortOrder = items[target].SortOrder, items[index].SortOrder
		_, _ = r.db.Exec(`UPDATE itinerary_items SET sort_order = ? WHERE user_id = ? AND itinerary_id = ? AND id = ?`, items[index].SortOrder, userID, itineraryID, items[index].ID)
		_, _ = r.db.Exec(`UPDATE itinerary_items SET sort_order = ? WHERE user_id = ? AND itinerary_id = ? AND id = ?`, items[target].SortOrder, userID, itineraryID, items[target].ID)
	}
	itinerary, ok, err := r.getLocked(userID, itineraryID)
	return itinerary, ok, err
}

func (r *SQLiteItineraryRepo) ListItems(userID string, itineraryID int) ([]domain.ItineraryItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.listItemsLocked(userID, itineraryID)
}

func (r *SQLiteItineraryRepo) getLocked(userID string, itineraryID int) (domain.Itinerary, bool, error) {
	itinerary, err := scanItinerary(r.db.QueryRow(`SELECT id, user_id, title, city, start_date, end_date, guests, budget, status, created_at, updated_at FROM itineraries WHERE user_id = ? AND id = ?`, userID, itineraryID))
	if err == sql.ErrNoRows {
		return domain.Itinerary{}, false, nil
	}
	if err != nil {
		return domain.Itinerary{}, false, err
	}
	items, err := r.listItemsLocked(userID, itineraryID)
	if err != nil {
		return domain.Itinerary{}, false, err
	}
	itinerary.Items = items
	return itinerary, true, nil
}

func (r *SQLiteItineraryRepo) addItemLocked(userID string, itineraryID int, item domain.ItineraryItem) (domain.ItineraryItem, error) {
	id, err := r.nextItemID(userID, itineraryID)
	if err != nil {
		return domain.ItineraryItem{}, err
	}
	if item.DayIndex <= 0 {
		item.DayIndex = 1
	}
	if item.ItemType == "" {
		item.ItemType = "note"
	}
	item.ID = id
	item.UserID = userID
	item.ItineraryID = itineraryID
	item.SortOrder = id
	_, err = r.db.Exec(`INSERT INTO itinerary_items(id, itinerary_id, user_id, day_index, start_time, end_time, item_type, product_id, destination_id, title, note, estimated_cost, sort_order) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, item.ID, itineraryID, userID, item.DayIndex, item.StartTime, item.EndTime, item.ItemType, item.ProductID, item.DestinationID, item.Title, item.Note, item.EstimatedCost, item.SortOrder)
	return item, err
}

func (r *SQLiteItineraryRepo) listItemsLocked(userID string, itineraryID int) ([]domain.ItineraryItem, error) {
	rows, err := r.db.Query(`SELECT id, itinerary_id, user_id, day_index, start_time, end_time, item_type, product_id, destination_id, title, note, estimated_cost, sort_order FROM itinerary_items WHERE user_id = ? AND itinerary_id = ? ORDER BY day_index ASC, sort_order ASC, id ASC`, userID, itineraryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.ItineraryItem, 0)
	for rows.Next() {
		var item domain.ItineraryItem
		if err := rows.Scan(&item.ID, &item.ItineraryID, &item.UserID, &item.DayIndex, &item.StartTime, &item.EndTime, &item.ItemType, &item.ProductID, &item.DestinationID, &item.Title, &item.Note, &item.EstimatedCost, &item.SortOrder); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *SQLiteItineraryRepo) nextItineraryID(userID string) (int, error) {
	var next sql.NullInt64
	if err := r.db.QueryRow(`SELECT MAX(id) + 1 FROM itineraries WHERE user_id = ?`, userID).Scan(&next); err != nil {
		return 0, err
	}
	if !next.Valid || next.Int64 <= 0 {
		return 1, nil
	}
	return int(next.Int64), nil
}

func (r *SQLiteItineraryRepo) nextItemID(userID string, itineraryID int) (int, error) {
	var next sql.NullInt64
	if err := r.db.QueryRow(`SELECT MAX(id) + 1 FROM itinerary_items WHERE user_id = ? AND itinerary_id = ?`, userID, itineraryID).Scan(&next); err != nil {
		return 0, err
	}
	if !next.Valid || next.Int64 <= 0 {
		return 1, nil
	}
	return int(next.Int64), nil
}

func scanItinerary(scanner interface {
	Scan(dest ...interface{}) error
}) (domain.Itinerary, error) {
	var item domain.Itinerary
	err := scanner.Scan(&item.ID, &item.UserID, &item.Title, &item.City, &item.StartDate, &item.EndDate, &item.Guests, &item.Budget, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}
