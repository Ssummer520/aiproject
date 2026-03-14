package main

import (
	"log"
	"net/http"
	"travel-api/services/bff/api"
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
	bffHandler := api.NewBFFHandler()

	mux := http.NewServeMux()
	// BFF Route for home page data
	mux.HandleFunc("/api/home", bffHandler.GetHomePage)

	// Keep old routes for compatibility while migrating frontend
	// ... we'll update frontend to use /api/home soon
	
	log.Println("Server listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", cors(mux)); err != nil {
		log.Fatal("Server failed:", err)
	}
}
