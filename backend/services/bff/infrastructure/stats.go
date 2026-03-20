package infrastructure

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"travel-api/internal/common/models"
)

const (
	statsFile         = "data/stats.json"
	bookingsFile      = "data/bookings.json"
	notificationsFile = "data/notifications.json"
)

type StatsStore struct {
	mu        sync.RWMutex
	views     map[int]int
	favorites map[int]int
}

func NewStatsStore() *StatsStore {
	store := &StatsStore{
		views:     make(map[int]int),
		favorites: make(map[int]int),
	}
	store.load()
	return store
}

func (s *StatsStore) IncrementView(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.views[id]++
	s.saveLocked()
}

func (s *StatsStore) IncrementFavorite(id int, delta int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.favorites[id] += delta
	if s.favorites[id] < 0 {
		s.favorites[id] = 0
	}
	s.saveLocked()
}

func (s *StatsStore) TopByViews(limit int) []int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	type pair struct {
		id    int
		count int
	}

	var sorted []pair
	for id, count := range s.views {
		sorted = append(sorted, pair{id, count})
	}

	// Simple bubble sort
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].count > sorted[i].count {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	result := make([]int, 0, limit)
	for i := 0; i < len(sorted) && i < limit; i++ {
		result = append(result, sorted[i].id)
	}

	return result
}

func (s *StatsStore) load() {
	_ = os.MkdirAll(filepath.Dir(statsFile), 0755)
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
	if payload.Views != nil {
		s.views = payload.Views
	}
	if payload.Favorites != nil {
		s.favorites = payload.Favorites
	}
}

func (s *StatsStore) saveLocked() {
	payload := struct {
		Views     map[int]int `json:"views"`
		Favorites map[int]int `json:"favorites"`
	}{
		Views:     s.views,
		Favorites: s.favorites,
	}
	b, _ := json.MarshalIndent(payload, "", "  ")
	_ = os.WriteFile(statsFile, b, 0644)
}

func (s *StatsStore) TopByFavorites(limit int) []int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	type pair struct {
		id    int
		count int
	}

	var sorted []pair
	for id, count := range s.favorites {
		sorted = append(sorted, pair{id, count})
	}

	// Simple bubble sort
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].count > sorted[i].count {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	result := make([]int, 0, limit)
	for i := 0; i < len(sorted) && i < limit; i++ {
		result = append(result, sorted[i].id)
	}

	return result
}

// BookingStore stores bookings
type BookingStore struct {
	mu       sync.RWMutex
	bookings map[string][]models.Booking
}

func NewBookingStore() *BookingStore {
	store := &BookingStore{
		bookings: make(map[string][]models.Booking),
	}
	store.load()
	return store
}

func (s *BookingStore) CreateBooking(userID string, dest models.Destination, checkIn, checkOut string, guests int) models.Booking {
	s.mu.Lock()
	defer s.mu.Unlock()

	booking := models.Booking{
		ID:            len(s.bookings[userID]) + 1,
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

	s.bookings[userID] = append(s.bookings[userID], booking)
	s.saveLocked()

	return booking
}

func (s *BookingStore) GetUserBookings(userID string) []models.Booking {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if bookings, ok := s.bookings[userID]; ok {
		result := make([]models.Booking, len(bookings))
		copy(result, bookings)
		return result
	}
	return []models.Booking{}
}

func (s *BookingStore) CancelBooking(userID string, bookingID int) (models.Booking, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	userBookings, ok := s.bookings[userID]
	if !ok {
		return models.Booking{}, false
	}

	for i := range userBookings {
		if userBookings[i].ID != bookingID {
			continue
		}
		if userBookings[i].Status == "cancelled" {
			return userBookings[i], true
		}
		userBookings[i].Status = "cancelled"
		userBookings[i].CancelledAt = time.Now().Format("2006-01-02 15:04:05")
		s.bookings[userID] = userBookings
		s.saveLocked()
		return userBookings[i], true
	}

	return models.Booking{}, false
}

func (s *BookingStore) load() {
	_ = os.MkdirAll(filepath.Dir(bookingsFile), 0755)
	b, err := os.ReadFile(bookingsFile)
	if err != nil {
		return
	}
	_ = json.Unmarshal(b, &s.bookings)
	if s.bookings == nil {
		s.bookings = make(map[string][]models.Booking)
	}
}

func (s *BookingStore) saveLocked() {
	b, _ := json.MarshalIndent(s.bookings, "", "  ")
	_ = os.WriteFile(bookingsFile, b, 0644)
}

func calculateNights(checkIn, checkOut string) int {
	// Simple calculation - in production would parse dates properly
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

// NotificationStore stores user notifications
type NotificationStore struct {
	mu            sync.RWMutex
	notifications map[string][]models.Notification
}

func NewNotificationStore() *NotificationStore {
	store := &NotificationStore{
		notifications: make(map[string][]models.Notification),
	}
	store.load()
	return store
}

func (s *NotificationStore) AddNotification(userID string, notification models.Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()
	notification.ID = len(s.notifications[userID]) + 1
	notification.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	s.notifications[userID] = append([]models.Notification{notification}, s.notifications[userID]...)
	s.saveLocked()
}

func (s *NotificationStore) GetNotifications(userID string) []models.Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if notifications, ok := s.notifications[userID]; ok {
		result := make([]models.Notification, len(notifications))
		copy(result, notifications)
		return result
	}
	return []models.Notification{}
}

func (s *NotificationStore) MarkAsRead(userID string, notificationID int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.notifications[userID] {
		if s.notifications[userID][i].ID == notificationID {
			s.notifications[userID][i].Read = true
			s.saveLocked()
			break
		}
	}
}

func (s *NotificationStore) GetUnreadCount(userID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	count := 0
	for _, n := range s.notifications[userID] {
		if !n.Read {
			count++
		}
	}
	return count
}

func (s *NotificationStore) load() {
	_ = os.MkdirAll(filepath.Dir(notificationsFile), 0755)
	b, err := os.ReadFile(notificationsFile)
	if err != nil {
		return
	}
	_ = json.Unmarshal(b, &s.notifications)
	if s.notifications == nil {
		s.notifications = make(map[string][]models.Notification)
	}
}

func (s *NotificationStore) saveLocked() {
	b, _ := json.MarshalIndent(s.notifications, "", "  ")
	_ = os.WriteFile(notificationsFile, b, 0644)
}
