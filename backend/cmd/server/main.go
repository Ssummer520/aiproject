package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"travel-api/internal/contextkeys"
	authApi "travel-api/services/auth/api"
	authApp "travel-api/services/auth/application"
	"travel-api/services/bff/api"
	couponApi "travel-api/services/coupon/api"
	couponApp "travel-api/services/coupon/application"
	orderApi "travel-api/services/order/api"
	orderApp "travel-api/services/order/application"
	productApi "travel-api/services/product/api"
	productApp "travel-api/services/product/application"
	reviewApi "travel-api/services/review/api"
	reviewApp "travel-api/services/review/application"
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
	productService := productApp.NewProductService()
	couponService := couponApp.NewCouponService()
	reviewService := reviewApp.NewReviewService()
	orderService := orderApp.NewOrderServiceWithCoupon(productService, couponService)
	authHandler := authApi.NewAuthHandlerWithService(authService)
	bffHandler := api.NewBFFHandler()
	productHandler := productApi.NewProductHandlerWithService(productService)
	couponHandler := couponApi.NewCouponHandler(couponService)
	reviewHandler := reviewApi.NewReviewHandler(reviewService)
	orderHandler := orderApi.NewOrderHandler(orderService)

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
	mux.Handle("/api/v1/search", authMiddleware(authService, http.HandlerFunc(bffHandler.Search)))
	mux.Handle("/api/v1/category/", authMiddleware(authService, http.HandlerFunc(bffHandler.GetCategory)))
	mux.Handle("/api/v1/city/", authMiddleware(authService, http.HandlerFunc(bffHandler.GetCity)))
	mux.Handle("/api/v1/destinations/", authMiddleware(authService, http.HandlerFunc(bffHandler.HandleDestinations)))
	mux.Handle("/api/v1/bookings", authMiddleware(authService, http.HandlerFunc(bffHandler.HandleBookings)))
	mux.Handle("/api/v1/bookings/", authMiddleware(authService, http.HandlerFunc(bffHandler.HandleBookingActions)))
	mux.Handle("/api/v1/notifications", authMiddleware(authService, http.HandlerFunc(bffHandler.HandleNotifications)))
	mux.HandleFunc("/api/v1/products", productHandler.HandleProducts)
	mux.Handle("/api/v1/products/", authMiddleware(authService, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(strings.Trim(r.URL.Path, "/"), "/reviews") {
			reviewHandler.HandleProductReviews(w, r)
			return
		}
		productHandler.HandleProductDetail(w, r)
	})))
	mux.HandleFunc("/api/v1/coupons", couponHandler.HandleCoupons)
	mux.HandleFunc("/api/v1/coupons/validate", couponHandler.HandleValidate)
	mux.Handle("/api/v1/orders", authMiddleware(authService, http.HandlerFunc(orderHandler.HandleOrders)))
	mux.Handle("/api/v1/orders/", authMiddleware(authService, http.HandlerFunc(orderHandler.HandleOrderActions)))

	const addr = ":8888"
	log.Println("Server listening on http://localhost:8888")
	if err := http.ListenAndServe(addr, cors(mux)); err != nil {
		log.Fatal("Server failed:", err)
	}
}
