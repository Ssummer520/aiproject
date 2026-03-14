package models

import "time"

// Destination 目的地/景点
type Destination struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Cover      string    `json:"cover"`
	Summary    string    `json:"summary"`
	Lat        float64   `json:"lat"`
	Lng        float64   `json:"lng"`
	City       string    `json:"city"`
	Tags       []string  `json:"tags"`
	Rating     float64   `json:"rating"`
	ViewedAt   time.Time `json:"viewed_at,omitempty"`
	IsFavorite bool      `json:"is_favorite,omitempty"`
	DistanceKm float64   `json:"distance_km,omitempty"`
}
