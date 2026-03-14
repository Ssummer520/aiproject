package main

import (
	"log"
	"net/http"
	"travel-api/internal/handlers"
)

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/recommendations", handlers.GetRecommendations)
	mux.HandleFunc("/api/recent-favorites", handlers.GetRecentAndFavorites)
	mux.HandleFunc("/api/nearby", handlers.GetNearby)
	mux.HandleFunc("/api/view", handlers.RecordView)
	mux.HandleFunc("/api/favorite", handlers.ToggleFavorite)

	log.Println("Server listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", cors(mux)); err != nil {
		log.Fatal("Server failed:", err)
	}
}
