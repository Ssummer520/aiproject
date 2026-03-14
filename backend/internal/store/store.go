package store

import (
	"sync"
	"travel-api/internal/models"
)

// Store 内存存储（演示用，可替换为数据库）
type Store struct {
	mu           sync.RWMutex
	Destinations []models.Destination
	Favorites    map[int]bool // destinationId -> true
	RecentViews  []int        // destination IDs, recent first
}

var defaultStore *Store

func init() {
	defaultStore = &Store{
		Favorites:   make(map[int]bool),
		RecentViews: []int{},
		Destinations: []models.Destination{
			{ID: 1, Name: "西湖", Cover: "https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=400", Summary: "杭州西湖，世界文化遗产", Lat: 30.2444, Lng: 120.1492, City: "杭州", Tags: []string{"湖泊", "文化遗产"}, Rating: 4.8},
			{ID: 2, Name: "外滩", Cover: "https://images.unsplash.com/photo-1547981609-4b6bfe67ca0b?w=400", Summary: "上海外滩万国建筑博览", Lat: 31.2397, Lng: 121.4912, City: "上海", Tags: []string{"夜景", "建筑"}, Rating: 4.7},
			{ID: 3, Name: "故宫", Cover: "https://images.unsplash.com/photo-1508807527081-8f81e0f6f8b8?w=400", Summary: "北京故宫博物院", Lat: 39.9163, Lng: 116.3972, City: "北京", Tags: []string{"古建筑", "博物馆"}, Rating: 4.9},
			{ID: 4, Name: "黄山", Cover: "https://images.unsplash.com/photo-1545569341-9eb8b30979d9?w=400", Summary: "黄山奇松怪石云海", Lat: 30.1333, Lng: 118.3333, City: "黄山", Tags: []string{"山岳", "自然"}, Rating: 4.8},
			{ID: 5, Name: "乌镇", Cover: "https://images.unsplash.com/photo-1517554558302-64d2fca8f9e2?w=400", Summary: "江南水乡古镇", Lat: 30.7411, Lng: 120.4856, City: "嘉兴", Tags: []string{"古镇", "水乡"}, Rating: 4.6},
			{ID: 6, Name: "灵隐寺", Cover: "https://images.unsplash.com/photo-1605649487212-47bdab064df7?w=400", Summary: "杭州灵隐禅寺", Lat: 30.2417, Lng: 120.0967, City: "杭州", Tags: []string{"寺庙", "人文"}, Rating: 4.5},
			{ID: 7, Name: "迪士尼乐园", Cover: "https://images.unsplash.com/photo-1531259683007-906a3792e424?w=400", Summary: "上海迪士尼度假区", Lat: 31.1447, Lng: 121.6572, City: "上海", Tags: []string{"主题乐园", "亲子"}, Rating: 4.7},
			{ID: 8, Name: "千岛湖", Cover: "https://images.unsplash.com/photo-1476514525535-07fb3b4ae5f1?w=400", Summary: "千岛湖风景区", Lat: 29.6050, Lng: 119.0389, City: "杭州", Tags: []string{"湖泊", "度假"}, Rating: 4.6},
		},
	}
	// 初始一些收藏和浏览记录
	defaultStore.Favorites[1] = true
	defaultStore.Favorites[3] = true
	defaultStore.RecentViews = []int{2, 5, 1, 4}
}

func GetStore() *Store { return defaultStore }

func (s *Store) GetDestination(id int) *models.Destination {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.Destinations {
		if s.Destinations[i].ID == id {
			d := s.Destinations[i]
			d.IsFavorite = s.Favorites[id]
			return &d
		}
	}
	return nil
}

func (s *Store) ListAll() []models.Destination {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]models.Destination, len(s.Destinations))
	copy(out, s.Destinations)
	for i := range out {
		out[i].IsFavorite = s.Favorites[out[i].ID]
	}
	return out
}

// GetRecentIDs 返回最近浏览的 ID 列表（顺序）
func (s *Store) GetRecentIDs() []int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ids := make([]int, len(s.RecentViews))
	copy(ids, s.RecentViews)
	return ids
}

// GetFavoriteIDs 返回收藏的 ID 集合
func (s *Store) GetFavoriteIDs() map[int]bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m := make(map[int]bool)
	for k, v := range s.Favorites {
		if v {
			m[k] = true
		}
	}
	return m
}

// RecordView 记录浏览（将 id 移到最近浏览最前）
func (s *Store) RecordView(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	newList := []int{id}
	for _, i := range s.RecentViews {
		if i != id {
			newList = append(newList, i)
		}
	}
	if len(newList) > 20 {
		newList = newList[:20]
	}
	s.RecentViews = newList
}

// ToggleFavorite 切换收藏
func (s *Store) ToggleFavorite(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Favorites[id] = !s.Favorites[id]
	return s.Favorites[id]
}
