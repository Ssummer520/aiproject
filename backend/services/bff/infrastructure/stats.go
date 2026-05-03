package infrastructure

import (
	"database/sql"
	"encoding/json"
	"os"
	"sort"
	"sync"
	"time"

	"travel-api/internal/common/models"
	"travel-api/internal/db"
)

const (
	statsFile         = "data/stats.json"
	bookingsFile      = "data/bookings.json"
	notificationsFile = "data/notifications.json"
)

type StatsStore struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewStatsStore() *StatsStore {
	database, err := db.Open()
	if err != nil {
		panic(err)
	}
	store := &StatsStore{db: database}
	store.migrateFromJSON()
	return store
}

func (s *StatsStore) IncrementView(id int) {
	s.increment(id, "views", 1)
}

func (s *StatsStore) IncrementFavorite(id int, delta int) {
	s.increment(id, "favorites", delta)
}

func (s *StatsStore) increment(id int, column string, delta int) {
	if id <= 0 || (column != "views" && column != "favorites") {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	_, _ = s.db.Exec(`INSERT OR IGNORE INTO stats(destination_id, views, favorites) VALUES(?, 0, 0)`, id)
	if delta >= 0 {
		_, _ = s.db.Exec(`UPDATE stats SET `+column+` = `+column+` + ? WHERE destination_id = ?`, delta, id)
		return
	}
	_, _ = s.db.Exec(`UPDATE stats SET `+column+` = MAX(0, `+column+` + ?) WHERE destination_id = ?`, delta, id)
}

func (s *StatsStore) TopByViews(limit int) []int {
	return s.topBy("views", limit)
}

func (s *StatsStore) TopByFavorites(limit int) []int {
	return s.topBy("favorites", limit)
}

func (s *StatsStore) topBy(column string, limit int) []int {
	if limit <= 0 || (column != "views" && column != "favorites") {
		return []int{}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	rows, err := s.db.Query(`SELECT destination_id FROM stats WHERE `+column+` > 0 ORDER BY `+column+` DESC, destination_id ASC LIMIT ?`, limit)
	if err != nil {
		return []int{}
	}
	defer rows.Close()

	ids := make([]int, 0, limit)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

func (s *StatsStore) migrateFromJSON() {
	s.mu.Lock()
	defer s.mu.Unlock()

	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM stats`).Scan(&count); err != nil || count > 0 {
		return
	}
	b, err := os.ReadFile(statsFile)
	if err != nil {
		return
	}
	var payload struct {
		Views     map[int]int `json:"views"`
		Favorites map[int]int `json:"favorites"`
	}
	if err := json.Unmarshal(b, &payload); err != nil {
		return
	}
	ids := make(map[int]bool)
	for id := range payload.Views {
		ids[id] = true
	}
	for id := range payload.Favorites {
		ids[id] = true
	}
	for id := range ids {
		_, _ = s.db.Exec(
			`INSERT OR REPLACE INTO stats(destination_id, views, favorites) VALUES(?, ?, ?)`,
			id,
			payload.Views[id],
			payload.Favorites[id],
		)
	}
}

type BookingStore struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewBookingStore() *BookingStore {
	database, err := db.Open()
	if err != nil {
		panic(err)
	}
	store := &BookingStore{db: database}
	store.migrateFromJSON()
	return store
}

func (s *BookingStore) CreateBooking(userID string, dest models.Destination, checkIn, checkOut string, guests int) models.Booking {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID(userID)
	booking := models.Booking{
		ID:            id,
		UserID:        userID,
		DestinationID: dest.ID,
		Name:          dest.Name,
		City:          dest.City,
		Cover:         dest.Cover,
		CheckIn:       checkIn,
		CheckOut:      checkOut,
		Guests:        guests,
		TotalPrice:    dest.Price * float64(calculateNights(checkIn, checkOut)),
		Status:        "confirmed",
		CreatedAt:     time.Now().Format("2006-01-02"),
	}

	_, _ = s.db.Exec(
		`INSERT INTO bookings(id, user_id, destination_id, name, city, cover, check_in, check_out, guests, total_price, status, created_at, cancelled_at)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		booking.ID,
		booking.UserID,
		booking.DestinationID,
		booking.Name,
		booking.City,
		booking.Cover,
		booking.CheckIn,
		booking.CheckOut,
		booking.Guests,
		booking.TotalPrice,
		booking.Status,
		booking.CreatedAt,
		nullString(booking.CancelledAt),
	)

	return booking
}

func (s *BookingStore) GetUserBookings(userID string) []models.Booking {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query(
		`SELECT id, user_id, destination_id, name, city, cover, check_in, check_out, guests, total_price, status, created_at, cancelled_at
		 FROM bookings WHERE user_id = ? ORDER BY id DESC`,
		userID,
	)
	if err != nil {
		return []models.Booking{}
	}
	defer rows.Close()

	bookings := make([]models.Booking, 0)
	for rows.Next() {
		booking, ok := scanBooking(rows)
		if ok {
			bookings = append(bookings, booking)
		}
	}
	return bookings
}

func (s *BookingStore) CancelBooking(userID string, bookingID int) (models.Booking, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	booking, ok := s.getBookingLocked(userID, bookingID)
	if !ok {
		return models.Booking{}, false
	}
	if booking.Status == "cancelled" {
		return booking, true
	}
	booking.Status = "cancelled"
	booking.CancelledAt = time.Now().Format("2006-01-02 15:04:05")
	_, _ = s.db.Exec(
		`UPDATE bookings SET status = ?, cancelled_at = ? WHERE user_id = ? AND id = ?`,
		booking.Status,
		booking.CancelledAt,
		userID,
		bookingID,
	)
	return booking, true
}

func (s *BookingStore) nextID(userID string) int {
	var id int
	_ = s.db.QueryRow(`SELECT COALESCE(MAX(id), 0) + 1 FROM bookings WHERE user_id = ?`, userID).Scan(&id)
	if id <= 0 {
		return 1
	}
	return id
}

func (s *BookingStore) getBookingLocked(userID string, bookingID int) (models.Booking, bool) {
	row := s.db.QueryRow(
		`SELECT id, user_id, destination_id, name, city, cover, check_in, check_out, guests, total_price, status, created_at, cancelled_at
		 FROM bookings WHERE user_id = ? AND id = ?`,
		userID,
		bookingID,
	)
	return scanBooking(row)
}

func (s *BookingStore) migrateFromJSON() {
	s.mu.Lock()
	defer s.mu.Unlock()

	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM bookings`).Scan(&count); err != nil || count > 0 {
		return
	}
	b, err := os.ReadFile(bookingsFile)
	if err != nil {
		return
	}
	var data map[string][]models.Booking
	if err := json.Unmarshal(b, &data); err != nil {
		return
	}
	for userID, bookings := range data {
		sort.SliceStable(bookings, func(i, j int) bool { return bookings[i].ID < bookings[j].ID })
		for _, booking := range bookings {
			if booking.UserID == "" {
				booking.UserID = userID
			}
			_, _ = s.db.Exec(
				`INSERT OR REPLACE INTO bookings(id, user_id, destination_id, name, city, cover, check_in, check_out, guests, total_price, status, created_at, cancelled_at)
				 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				booking.ID,
				booking.UserID,
				booking.DestinationID,
				booking.Name,
				booking.City,
				booking.Cover,
				booking.CheckIn,
				booking.CheckOut,
				booking.Guests,
				booking.TotalPrice,
				booking.Status,
				booking.CreatedAt,
				nullString(booking.CancelledAt),
			)
		}
	}
}

func scanBooking(scanner interface {
	Scan(dest ...interface{}) error
}) (models.Booking, bool) {
	var booking models.Booking
	var cancelledAt sql.NullString
	if err := scanner.Scan(
		&booking.ID,
		&booking.UserID,
		&booking.DestinationID,
		&booking.Name,
		&booking.City,
		&booking.Cover,
		&booking.CheckIn,
		&booking.CheckOut,
		&booking.Guests,
		&booking.TotalPrice,
		&booking.Status,
		&booking.CreatedAt,
		&cancelledAt,
	); err != nil {
		return models.Booking{}, false
	}
	if cancelledAt.Valid {
		booking.CancelledAt = cancelledAt.String
	}
	return booking, true
}

func calculateNights(checkIn, checkOut string) int {
	inTime, err := time.Parse("2006-01-02", checkIn)
	if err != nil {
		return 1
	}
	outTime, err := time.Parse("2006-01-02", checkOut)
	if err != nil {
		return 1
	}
	nights := int(outTime.Sub(inTime).Hours() / 24)
	if nights <= 0 {
		return 1
	}
	return nights
}

type NotificationStore struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewNotificationStore() *NotificationStore {
	database, err := db.Open()
	if err != nil {
		panic(err)
	}
	store := &NotificationStore{db: database}
	store.migrateFromJSON()
	return store
}

func (s *NotificationStore) AddNotification(userID string, notification models.Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()
	notification.ID = s.nextID(userID)
	notification.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	_, _ = s.db.Exec(
		`INSERT INTO notifications(id, user_id, type, title, message, link, read, created_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`,
		notification.ID,
		userID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.Link,
		boolToInt(notification.Read),
		notification.CreatedAt,
	)
}

func (s *NotificationStore) GetNotifications(userID string) []models.Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rows, err := s.db.Query(
		`SELECT id, type, title, message, link, read, created_at FROM notifications WHERE user_id = ? ORDER BY id DESC`,
		userID,
	)
	if err != nil {
		return []models.Notification{}
	}
	defer rows.Close()

	notifications := make([]models.Notification, 0)
	for rows.Next() {
		var notification models.Notification
		var read int
		if err := rows.Scan(
			&notification.ID,
			&notification.Type,
			&notification.Title,
			&notification.Message,
			&notification.Link,
			&read,
			&notification.CreatedAt,
		); err == nil {
			notification.Read = read == 1
			notifications = append(notifications, notification)
		}
	}
	return notifications
}

func (s *NotificationStore) MarkAsRead(userID string, notificationID int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, _ = s.db.Exec(`UPDATE notifications SET read = 1 WHERE user_id = ? AND id = ?`, userID, notificationID)
}

func (s *NotificationStore) GetUnreadCount(userID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var count int
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE user_id = ? AND read = 0`, userID).Scan(&count)
	return count
}

func (s *NotificationStore) nextID(userID string) int {
	var id int
	_ = s.db.QueryRow(`SELECT COALESCE(MAX(id), 0) + 1 FROM notifications WHERE user_id = ?`, userID).Scan(&id)
	if id <= 0 {
		return 1
	}
	return id
}

func (s *NotificationStore) migrateFromJSON() {
	s.mu.Lock()
	defer s.mu.Unlock()

	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM notifications`).Scan(&count); err != nil || count > 0 {
		return
	}
	b, err := os.ReadFile(notificationsFile)
	if err != nil {
		return
	}
	var data map[string][]models.Notification
	if err := json.Unmarshal(b, &data); err != nil {
		return
	}
	for userID, notifications := range data {
		sort.SliceStable(notifications, func(i, j int) bool { return notifications[i].ID < notifications[j].ID })
		for _, notification := range notifications {
			_, _ = s.db.Exec(
				`INSERT OR REPLACE INTO notifications(id, user_id, type, title, message, link, read, created_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`,
				notification.ID,
				userID,
				notification.Type,
				notification.Title,
				notification.Message,
				notification.Link,
				boolToInt(notification.Read),
				notification.CreatedAt,
			)
		}
	}
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func nullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}
