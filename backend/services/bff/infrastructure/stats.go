package infrastructure

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
)

const statsFile = "data/stats.json"

type StatsStore struct {
	mu        sync.RWMutex
	Views     map[int]int `json:"views"`
	Favorites map[int]int `json:"favorites"`
	path      string
}

func NewStatsStore() *StatsStore {
	s := &StatsStore{
		Views:     make(map[int]int),
		Favorites: make(map[int]int),
		path:      statsFile,
	}
	s.load()
	return s
}

func (s *StatsStore) load() {
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = os.MkdirAll(filepath.Dir(s.path), 0755)
	b, err := os.ReadFile(s.path)
	if err != nil {
		return
	}
	var raw struct {
		Views     map[string]int `json:"views"`
		Favorites map[string]int `json:"favorites"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return
	}
	if raw.Views != nil {
		s.Views = make(map[int]int)
		for k, v := range raw.Views {
			id, _ := strconv.Atoi(k)
			if id > 0 {
				s.Views[id] = v
			}
		}
	}
	if raw.Favorites != nil {
		s.Favorites = make(map[int]int)
		for k, v := range raw.Favorites {
			id, _ := strconv.Atoi(k)
			if id > 0 {
				s.Favorites[id] = v
			}
		}
	}
	if s.Views == nil {
		s.Views = make(map[int]int)
	}
	if s.Favorites == nil {
		s.Favorites = make(map[int]int)
	}
}

func (s *StatsStore) save() {
	raw := struct {
		Views     map[string]int `json:"views"`
		Favorites map[string]int `json:"favorites"`
	}{
		Views:     make(map[string]int),
		Favorites: make(map[string]int),
	}
	for k, v := range s.Views {
		raw.Views[strconv.Itoa(k)] = v
	}
	for k, v := range s.Favorites {
		raw.Favorites[strconv.Itoa(k)] = v
	}
	b, _ := json.MarshalIndent(raw, "", "  ")
	_ = os.WriteFile(s.path, b, 0644)
}

func (s *StatsStore) IncrementView(id int) {
	if id <= 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Views[id]++
	s.save()
}

func (s *StatsStore) IncrementFavorite(id int, delta int) {
	if id <= 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Favorites[id] += delta
	if s.Favorites[id] < 0 {
		s.Favorites[id] = 0
	}
	s.save()
}

// TopByViews returns destination IDs sorted by view count descending.
func (s *StatsStore) TopByViews(limit int) []int {
	s.mu.RLock()
	type pair struct{ id, count int }
	var pairs []pair
	for id, count := range s.Views {
		pairs = append(pairs, pair{id, count})
	}
	s.mu.RUnlock()
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].count > pairs[j].count })
	if limit <= 0 || limit > len(pairs) {
		limit = len(pairs)
	}
	out := make([]int, 0, limit)
	for i := 0; i < limit && i < len(pairs); i++ {
		out = append(out, pairs[i].id)
	}
	return out
}

// TopByFavorites returns destination IDs sorted by favorite count descending.
func (s *StatsStore) TopByFavorites(limit int) []int {
	s.mu.RLock()
	type pair struct{ id, count int }
	var pairs []pair
	for id, count := range s.Favorites {
		pairs = append(pairs, pair{id, count})
	}
	s.mu.RUnlock()
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].count > pairs[j].count })
	if limit <= 0 || limit > len(pairs) {
		limit = len(pairs)
	}
	out := make([]int, 0, limit)
	for i := 0; i < limit && i < len(pairs); i++ {
		out = append(out, pairs[i].id)
	}
	return out
}
