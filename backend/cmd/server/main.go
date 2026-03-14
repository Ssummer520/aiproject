package main

import (
	"context"
	"log"
	"net/http"

	"travel-api/internal/contextkeys"
	authApi "travel-api/services/auth/api"
	authApp "travel-api/services/auth/application"
	"travel-api/services/bff/api"
)

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept-Language, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(authService *authApp.AuthService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := authApi.ExtractBearer(r)
		userID, _ := authService.ValidateToken(token)
		ctx := context.WithValue(r.Context(), contextkeys.UserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	authService := authApp.NewAuthService()
	authHandler := authApi.NewAuthHandler()
	bffHandler := api.NewBFFHandler()

	mux := http.NewServeMux()

	// Auth (no userID required)
	mux.HandleFunc("/api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("/api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("/api/v1/auth/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("/api/v1/auth/reset-password", authHandler.ResetPassword)
	mux.HandleFunc("/api/v1/auth/me", authHandler.Me)
	mux.HandleFunc("/api/v1/auth/logout", authHandler.Logout)

	// BFF (auth middleware sets userID in context; empty if not logged in)
	mux.Handle("/api/v1/home", authMiddleware(authService, http.HandlerFunc(bffHandler.GetHomePage)))
	mux.Handle("/api/v1/destinations/", authMiddleware(authService, http.HandlerFunc(bffHandler.HandleDestinations)))

	log.Println("Server listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", cors(mux)); err != nil {
		log.Fatal("Server failed:", err)
	}
}
