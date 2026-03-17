package infrastructure

import (
	"sync"
	"time"

	"travel-api/internal/common/models"
)

type StatsStore struct {
	mu        sync.RWMutex
	views     map[int]int
	favorites map[int]int
}

func NewStatsStore() *StatsStore {
	return &StatsStore{
		views:     make(map[int]int),
		favorites: make(map[int]int),
	}
}

func (s *StatsStore) IncrementView(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.views[id]++
}

func (s *StatsStore) IncrementFavorite(id int, delta int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.favorites[id] += delta
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
	return &BookingStore{
		bookings: make(map[string][]models.Booking),
	}
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

	return booking
}

func (s *BookingStore) GetUserBookings(userID string) []models.Booking {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.bookings[userID]
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
