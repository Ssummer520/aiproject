package infrastructure

import (
	"sync"
)

type InteractionCache struct {
	mu        sync.RWMutex
	favorites map[int]bool
	history   []int
}

func NewInteractionCache() *InteractionCache {
	return &InteractionCache{
		favorites: map[int]bool{
			1: true, // West Lake
			3: true, // Great Wall
		},
		history: []int{2, 4, 5}, // Shanghai, Yellow Mountain, Terracotta Army
	}
}

func (c *InteractionCache) ToggleFavorite(id int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.favorites[id] = !c.favorites[id]
	return c.favorites[id]
}

func (c *InteractionCache) GetFavorites() []int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	favs := make([]int, 0)
	for id, isFav := range c.favorites {
		if isFav {
			favs = append(favs, id)
		}
	}
	return favs
}

func (c *InteractionCache) AddToHistory(id int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Remove if already exists to move to top
	for i, existingID := range c.history {
		if existingID == id {
			c.history = append(c.history[:i], c.history[i+1:]...)
			break
		}
	}

	// Add to top
	c.history = append([]int{id}, c.history...)

	// Limit history size
	if len(c.history) > 10 {
		c.history = c.history[:10]
	}
}

func (c *InteractionCache) GetHistory() []int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.history
}

func (c *InteractionCache) IsFavorite(id int) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.favorites[id]
}
