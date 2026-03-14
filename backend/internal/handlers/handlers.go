package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"travel-api/internal/models"
	"travel-api/internal/store"
)

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}

// GetRecommendations 首页推荐（全部或精选）
func GetRecommendations(w http.ResponseWriter, r *http.Request) {
	s := store.GetStore()
	list := s.ListAll()
	writeJSON(w, map[string]interface{}{"list": list})
}

// GetRecentAndFavorites 最近浏览 + 收藏
func GetRecentAndFavorites(w http.ResponseWriter, r *http.Request) {
	s := store.GetStore()
	recentIDs := s.GetRecentIDs()
	favMap := s.GetFavoriteIDs()

	var recent, favorites []models.Destination
	for _, id := range recentIDs {
		d := s.GetDestination(id)
		if d != nil {
			d.IsFavorite = favMap[d.ID]
			recent = append(recent, *d)
		}
	}
	for id := range favMap {
		d := s.GetDestination(id)
		if d != nil {
			d.IsFavorite = true
			favorites = append(favorites, *d)
		}
	}
	writeJSON(w, map[string]interface{}{
		"recent":    recent,
		"favorites": favorites,
	})
}

// GetNearby 周边目的地（按经纬度距离排序，可选 lat, lng, limit）
func GetNearby(w http.ResponseWriter, r *http.Request) {
	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lng, _ := strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10
	}
	if lat == 0 && lng == 0 {
		lat, lng = 30.25, 120.15 // 默认杭州附近
	}

	s := store.GetStore()
	list := s.ListAll()
	type item struct {
		models.Destination
		DistanceKm float64 `json:"distance_km"`
	}
	var items []item
	for _, d := range list {
		km := haversineKm(lat, lng, d.Lat, d.Lng)
		items = append(items, item{Destination: d, DistanceKm: math.Round(km*10) / 10})
	}
	// 按距离排序
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].DistanceKm < items[i].DistanceKm {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	if len(items) > limit {
		items = items[:limit]
	}
	writeJSON(w, map[string]interface{}{"list": items})
}

func haversineKm(lat1, lng1, lat2, lng2 float64) float64 {
	const r = 6371 // 地球半径 km
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return r * c
}

// RecordView 记录浏览
func RecordView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if id <= 0 {
		writeJSON(w, map[string]interface{}{"ok": false})
		return
	}
	store.GetStore().RecordView(id)
	writeJSON(w, map[string]interface{}{"ok": true})
}

// ToggleFavorite 切换收藏
func ToggleFavorite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if id <= 0 {
		writeJSON(w, map[string]interface{}{"ok": false, "is_favorite": false})
		return
	}
	isFav := store.GetStore().ToggleFavorite(id)
	writeJSON(w, map[string]interface{}{"ok": true, "is_favorite": isFav})
}
