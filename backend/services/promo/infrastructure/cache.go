package infrastructure

import (
	"sync"
	"travel-api/internal/common/models"
)

type PromoCache struct {
	mu   sync.RWMutex
	data []models.Deal
}

func NewPromoCache() *PromoCache {
	return &PromoCache{}
}

func (c *PromoCache) ListDeals() []models.Deal {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.data) == 0 {
		// Fallback initial data
		return []models.Deal{
			{ID: 1, Title: "Spring Break Deals", Description: "Save 80¥ on bookings over 500¥", Type: "primary", Badge: "Limited Time", Expiry: "12:45:03"},
			{ID: 2, Title: "New User Gift", Description: "30¥ OFF your first trip in China", Type: "secondary", Badge: "", Expiry: ""},
			{ID: 3, Title: "Weekend Getaway", Description: "Up to 50% off for local experiences", Type: "accent", Badge: "", Expiry: ""},
			{ID: 4, Title: "Family Fun Pass", Description: "Bundle 2 adult + 1 kid tickets and save 20%", Type: "primary", Badge: "Popular", Expiry: "36:00:00"},
			{ID: 5, Title: "Museum Night Special", Description: "Late-entry museum packages from 99¥", Type: "secondary", Badge: "", Expiry: ""},
			{ID: 6, Title: "Foodie Trail Package", Description: "Street food + local guide combo from 129¥", Type: "accent", Badge: "", Expiry: ""},
		}
	}
	return c.data
}
