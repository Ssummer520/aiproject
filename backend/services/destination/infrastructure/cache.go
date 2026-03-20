package infrastructure

import (
	"sort"
	"sync"
	"travel-api/internal/common/models"
)

type DestinationCache struct {
	mu   sync.RWMutex
	data map[int]models.Destination
}

var fallbackDestinations = []models.Destination{
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
	{
		ID: 6, Name: "Lingyin Temple", City: "Hangzhou", Rating: 4.7, Lat: 30.24, Lng: 120.10,
		Tags: []string{"Culture", "History"}, Price: 128, ReviewCount: 2200, BookedCount: 36,
		Cover:       "https://images.unsplash.com/photo-1605649487212-47bdab064df7?w=800",
		Description: "A thousand-year-old Buddhist temple tucked into forested hills near West Lake.",
		Amenities:   []string{"Temple Audio Tour", "Incense Kit", "Shuttle"},
		HostName:    "Lingyin Cultural Tours",
		Policy:      "Free cancellation up to 24h before arrival.",
	},
	{
		ID: 7, Name: "Shanghai Disney Resort", City: "Shanghai", Rating: 4.8, Lat: 31.14, Lng: 121.66,
		Tags: []string{"Theme Park", "Family"}, Price: 399, ReviewCount: 7600, BookedCount: 92,
		Cover:       "https://images.unsplash.com/photo-1531259683007-016a7b628fc3?w=800",
		Description: "A full-day magical theme park experience with live shows and signature attractions.",
		Amenities:   []string{"Fast Pass", "Locker", "Family Rest Area"},
		HostName:    "Shanghai Resort Operations",
		Policy:      "Reschedule once up to 48h before entry.",
	},
	{
		ID: 8, Name: "Wuzhen Water Town", City: "Hangzhou", Rating: 4.7, Lat: 30.74, Lng: 120.49,
		Tags: []string{"Culture", "City"}, Price: 188, ReviewCount: 3100, BookedCount: 47,
		Cover:       "https://images.unsplash.com/photo-1517554558302-64d2fca8f9e2?w=800",
		Description: "Ancient canals, stone bridges, and classic Jiangnan architecture in a heritage town.",
		Amenities:   []string{"Boat Ride", "Night Pass", "Local Guide"},
		HostName:    "Jiangnan Heritage Co.",
		Policy:      "Free cancellation up to 24h before arrival.",
	},
	{
		ID: 9, Name: "Summer Palace", City: "Beijing", Rating: 4.8, Lat: 39.99, Lng: 116.27,
		Tags: []string{"History", "Nature"}, Price: 138, ReviewCount: 6400, BookedCount: 61,
		Cover:       "https://images.unsplash.com/photo-1508807527081-8f81e0f6f8b8?w=800",
		Description: "Imperial gardens, lake scenery, and palace halls make this a classic Beijing day trip.",
		Amenities:   []string{"Garden Map", "Boat Ticket", "Audio Guide"},
		HostName:    "Capital Heritage Group",
		Policy:      "Free cancellation up to 24h before arrival.",
	},
	{
		ID: 10, Name: "Forbidden City", City: "Beijing", Rating: 4.9, Lat: 39.92, Lng: 116.40,
		Tags: []string{"History", "Museum"}, Price: 179, ReviewCount: 9200, BookedCount: 89,
		Cover:       "https://images.unsplash.com/photo-1508807527081-8f81e0f6f8b8?w=800",
		Description: "Explore the grand palace complex that served as the imperial home for nearly 500 years.",
		Amenities:   []string{"Priority Entry", "Audio Guide", "Museum Pass"},
		HostName:    "Forbidden City Experiences",
		Policy:      "No cancellation within 24h of visit date.",
	},
	{
		ID: 11, Name: "Chengdu Panda Base", City: "Chengdu", Rating: 4.9, Lat: 30.73, Lng: 104.14,
		Tags: []string{"Nature", "Family"}, Price: 149, ReviewCount: 5100, BookedCount: 73,
		Cover:       "https://images.unsplash.com/photo-1535930749574-1399327ce78f?w=800",
		Description: "Meet giant pandas in a conservation-first habitat with educational exhibits.",
		Amenities:   []string{"Park Shuttle", "Photo Spot", "Guide Service"},
		HostName:    "Chengdu Wildlife Center",
		Policy:      "Free cancellation up to 48h before arrival.",
	},
	{
		ID: 12, Name: "Jinli Ancient Street", City: "Chengdu", Rating: 4.6, Lat: 30.65, Lng: 104.05,
		Tags: []string{"Food", "Culture"}, Price: 118, ReviewCount: 2800, BookedCount: 31,
		Cover:       "https://images.unsplash.com/photo-1553856622-d1b352e1f6dc?w=800",
		Description: "A lively old street known for Sichuan snacks, local crafts, and evening performances.",
		Amenities:   []string{"Snack Coupon", "Local Show", "Walking Route"},
		HostName:    "Chengdu Local Walks",
		Policy:      "Free cancellation up to 24h before arrival.",
	},
	{
		ID: 13, Name: "Xi'an City Wall", City: "Xi'an", Rating: 4.7, Lat: 34.26, Lng: 108.95,
		Tags: []string{"History", "Hiking"}, Price: 126, ReviewCount: 3300, BookedCount: 42,
		Cover:       "https://images.unsplash.com/photo-1523482580672-f109ba8cb9be?w=800",
		Description: "Walk or bike along one of the best-preserved ancient city walls in China.",
		Amenities:   []string{"Bike Rental", "Night View Route", "Guidebook"},
		HostName:    "Xi'an Culture Travel",
		Policy:      "Free cancellation up to 24h before arrival.",
	},
	{
		ID: 14, Name: "Oriental Pearl Tower", City: "Shanghai", Rating: 4.7, Lat: 31.24, Lng: 121.50,
		Tags: []string{"City", "Night View"}, Price: 238, ReviewCount: 4700, BookedCount: 58,
		Cover:       "https://images.unsplash.com/photo-1548115184-bc65ee498ad0?w=800",
		Description: "Skyline observatory with panoramic views over the Huangpu River and downtown Shanghai.",
		Amenities:   []string{"Fast Track", "Observation Deck", "Souvenir Store"},
		HostName:    "Pearl Tower Group",
		Policy:      "Reschedule up to 24h before arrival.",
	},
	{
		ID: 15, Name: "Longjing Tea Plantation", City: "Hangzhou", Rating: 4.8, Lat: 30.22, Lng: 120.10,
		Tags: []string{"Nature", "Food"}, Price: 158, ReviewCount: 1900, BookedCount: 28,
		Cover:       "https://images.unsplash.com/photo-1515488042361-ee00e0ddd4e4?w=800",
		Description: "Tea field trails, hand-picking sessions, and tasting workshops in Longjing village.",
		Amenities:   []string{"Tea Tasting", "Farm Tour", "Shuttle"},
		HostName:    "Longjing Estate",
		Policy:      "Free cancellation up to 24h before arrival.",
	},
	{
		ID: 16, Name: "Universal Beijing Resort", City: "Beijing", Rating: 4.8, Lat: 39.90, Lng: 116.66,
		Tags: []string{"Theme Park", "Family"}, Price: 429, ReviewCount: 6100, BookedCount: 84,
		Cover:       "https://images.unsplash.com/photo-1543968996-ee822b8176ba?w=800",
		Description: "Cinematic rides and immersive worlds designed for all-age family adventures.",
		Amenities:   []string{"Express Pass", "Locker", "Family Lounge"},
		HostName:    "Universal Beijing",
		Policy:      "Reschedule up to 72h before entry.",
	},
	{
		ID: 17, Name: "Zhujiajiao Water Town", City: "Shanghai", Rating: 4.6, Lat: 31.11, Lng: 121.05,
		Tags: []string{"Culture", "Nature"}, Price: 139, ReviewCount: 2500, BookedCount: 37,
		Cover:       "https://images.unsplash.com/photo-1496939376851-89342e90adcd?w=800",
		Description: "An easy day-trip old town with canals, arched bridges, and riverside tea houses.",
		Amenities:   []string{"Boat Ticket", "Local Snacks", "Walking Map"},
		HostName:    "Shanghai Heritage Trails",
		Policy:      "Free cancellation up to 24h before arrival.",
	},
	{
		ID: 18, Name: "Mount Qingcheng", City: "Chengdu", Rating: 4.7, Lat: 30.90, Lng: 103.57,
		Tags: []string{"Nature", "Hiking"}, Price: 169, ReviewCount: 2100, BookedCount: 33,
		Cover:       "https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?w=800",
		Description: "A lush mountain route with Taoist heritage sites, streams, and forest air.",
		Amenities:   []string{"Trail Guide", "Cable Car", "Tea Rest Stop"},
		HostName:    "Sichuan Mountain Club",
		Policy:      "Free cancellation up to 24h before arrival.",
	},
}

func NewDestinationCache() *DestinationCache {
	return &DestinationCache{
		data: make(map[int]models.Destination),
	}
}

func (c *DestinationCache) Get(id int) (models.Destination, bool) {
	c.mu.RLock()
	d, ok := c.data[id]
	c.mu.RUnlock()
	if ok {
		return d, true
	}

	for _, f := range fallbackDestinations {
		if f.ID == id {
			return f, true
		}
	}
	return models.Destination{}, false
}

func (c *DestinationCache) ListAll() []models.Destination {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.data) == 0 {
		out := make([]models.Destination, len(fallbackDestinations))
		copy(out, fallbackDestinations)
		return out
	}

	list := make([]models.Destination, 0, len(c.data))
	for _, d := range c.data {
		list = append(list, d)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	return list
}

func (c *DestinationCache) Set(d models.Destination) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[d.ID] = d
}
