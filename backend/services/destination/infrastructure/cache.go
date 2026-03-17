package infrastructure

import (
	"sync"
	"travel-api/internal/common/models"
)

type DestinationCache struct {
	mu   sync.RWMutex
	data map[int]models.Destination
}

func NewDestinationCache() *DestinationCache {
	return &DestinationCache{
		data: make(map[int]models.Destination),
	}
}

func (c *DestinationCache) Get(id int) (models.Destination, bool) {
	// First check fallback data (without lock)
	fallbacks := []models.Destination{
		{ID: 1, Name: "West Lake", City: "Hangzhou", Rating: 4.9, Lat: 30.25, Lng: 120.15,
			Tags: []string{"Nature", "Culture"}, Price: 168, ReviewCount: 1200, BookedCount: 14,
			Cover:       "https://images.unsplash.com/photo-1547981609-4b6bfe67ca0b?w=800",
			Description: "A legendary freshwater lake in Hangzhou, known for its scenic beauty and cultural significance.",
			Amenities:   []string{"Wi-Fi", "Tea Garden", "Guided Tour"},
			HostName:    "Local Guide Association",
			Policy:      "Free cancellation up to 48h before arrival."},
		{ID: 2, Name: "The Bund", City: "Shanghai", Rating: 4.8, Lat: 31.24, Lng: 121.49,
			Tags: []string{"City", "Night View"}, Price: 268, ReviewCount: 3500, BookedCount: 45,
			Cover:       "https://images.unsplash.com/photo-1548115184-bc65ee498ad0?w=800",
			Description: "The iconic waterfront area in Shanghai, showcasing colonial architecture and futuristic skyline.",
			Amenities:   []string{"Elevator", "Rooftop Bar", "Parking"},
			HostName:    "Skyline Properties",
			Policy:      "Non-refundable booking."},
		{ID: 3, Name: "Great Wall", City: "Beijing", Rating: 5.0, Lat: 40.43, Lng: 116.57,
			Tags: []string{"History", "Hiking"}, Price: 198, ReviewCount: 8900, BookedCount: 78,
			Cover:       "https://images.unsplash.com/photo-1508804185872-d7badad00f7d?w=800",
			Description: "One of the greatest wonders of the world, stretching thousands of miles across northern China.",
			Amenities:   []string{"Hiking Map", "First Aid Kit", "Bottle of Water"},
			HostName:    "Heritage Tours",
			Policy:      "Cancel up to 24h before for a full refund."},
		{ID: 4, Name: "Yellow Mountain", City: "Huangshan", Rating: 4.9, Lat: 30.13, Lng: 118.16,
			Tags: []string{"Nature", "Hiking"}, Price: 228, ReviewCount: 2100, BookedCount: 23,
			Cover:       "https://images.unsplash.com/photo-1525113190471-9969be29263a?w=800",
			Description: "Famous for its peculiar pines, odd rocks, sea of clouds, and hot springs.",
			Amenities:   []string{"Mountain Gear", "Sunrise Viewpoint", "Oxygen Tank"},
			HostName:    "Peak Adventures",
			Policy:      "Reschedule possible up to 72h before."},
		{ID: 5, Name: "Terracotta Army", City: "Xi'an", Rating: 4.8, Lat: 34.38, Lng: 109.28,
			Tags: []string{"History", "Archaeology"}, Price: 158, ReviewCount: 5400, BookedCount: 56,
			Cover:       "https://images.unsplash.com/photo-1523482580672-f109ba8cb9be?w=800",
			Description: "A collection of terracotta sculptures depicting the armies of Qin Shi Huang, the first Emperor of China.",
			Amenities:   []string{"Museum Pass", "Audio Guide", "Souvenir Shop"},
			HostName:    "Qin Dynasty Experts",
			Policy:      "No cancellations allowed."},
	}
	for _, f := range fallbacks {
		if f.ID == id {
			return f, true
		}
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	d, ok := c.data[id]
	return d, ok
}

func (c *DestinationCache) ListAll() []models.Destination {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.data) == 0 {
		// Fallback initial data with enriched sub-features
		return []models.Destination{
			{
				ID: 1, Name: "West Lake", City: "Hangzhou", Rating: 4.9, Lat: 30.25, Lng: 120.15,
				Tags: []string{"Nature", "Culture"}, Price: 168, ReviewCount: 1200, BookedCount: 14,
				Cover:       "https://images.unsplash.com/photo-1547981609-4b6bfe67ca0b?w=800",
				Description: "A legendary freshwater lake in Hangzhou, known for its scenic beauty and cultural significance.",
				Amenities:   []string{"Wi-Fi", "Tea Garden", "Guided Tour"},
				HostName:    "Local Guide Association",
				Policy:      "Free cancellation up to 48h before arrival.",
			},
			{
				ID: 2, Name: "The Bund", City: "Shanghai", Rating: 4.8, Lat: 31.24, Lng: 121.49,
				Tags: []string{"City", "Night View"}, Price: 268, ReviewCount: 3500, BookedCount: 45,
				Cover:       "https://images.unsplash.com/photo-1548115184-bc65ee498ad0?w=800",
				Description: "The iconic waterfront area in Shanghai, showcasing colonial architecture and futuristic skyline.",
				Amenities:   []string{"Elevator", "Rooftop Bar", "Parking"},
				HostName:    "Skyline Properties",
				Policy:      "Non-refundable booking.",
			},
			{
				ID: 3, Name: "Great Wall", City: "Beijing", Rating: 5.0, Lat: 40.43, Lng: 116.57,
				Tags: []string{"History", "Hiking"}, Price: 198, ReviewCount: 8900, BookedCount: 78,
				Cover:       "https://images.unsplash.com/photo-1508804185872-d7badad00f7d?w=800",
				Description: "One of the greatest wonders of the world, stretching thousands of miles across northern China.",
				Amenities:   []string{"Hiking Map", "First Aid Kit", "Bottle of Water"},
				HostName:    "Heritage Tours",
				Policy:      "Cancel up to 24h before for a full refund.",
			},
			{
				ID: 4, Name: "Yellow Mountain", City: "Huangshan", Rating: 4.9, Lat: 30.13, Lng: 118.16,
				Tags: []string{"Nature", "Hiking"}, Price: 228, ReviewCount: 2100, BookedCount: 23,
				Cover:       "https://images.unsplash.com/photo-1525113190471-9969be29263a?w=800",
				Description: "Famous for its peculiar pines, odd rocks, sea of clouds, and hot springs.",
				Amenities:   []string{"Mountain Gear", "Sunrise Viewpoint", "Oxygen Tank"},
				HostName:    "Peak Adventures",
				Policy:      "Reschedule possible up to 72h before.",
			},
			{
				ID: 5, Name: "Terracotta Army", City: "Xi'an", Rating: 4.8, Lat: 34.38, Lng: 109.28,
				Tags: []string{"History", "Archaeology"}, Price: 158, ReviewCount: 5400, BookedCount: 56,
				Cover:       "https://images.unsplash.com/photo-1523482580672-f109ba8cb9be?w=800",
				Description: "A collection of terracotta sculptures depicting the armies of Qin Shi Huang, the first Emperor of China.",
				Amenities:   []string{"Museum Pass", "Audio Guide", "Souvenir Shop"},
				HostName:    "Qin Dynasty Experts",
				Policy:      "No cancellations allowed.",
			},
		}
	}

	list := make([]models.Destination, 0, len(c.data))
	for _, d := range c.data {
		list = append(list, d)
	}
	return list
}

func (c *DestinationCache) Set(d models.Destination) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[d.ID] = d
}
