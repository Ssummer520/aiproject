package domain

type CartItem struct {
	ID              int     `json:"id"`
	UserID          string  `json:"user_id"`
	ProductID       int     `json:"product_id"`
	PackageID       int     `json:"package_id"`
	ProductName     string  `json:"product_name"`
	PackageName     string  `json:"package_name"`
	City            string  `json:"city"`
	Cover           string  `json:"cover"`
	TravelDate      string  `json:"travel_date"`
	Adults          int     `json:"adults"`
	Children        int     `json:"children"`
	Quantity        int     `json:"quantity"`
	UnitPrice       float64 `json:"unit_price"`
	Subtotal        float64 `json:"subtotal"`
	SelectedOptions string  `json:"selected_options"`
	CreatedAt       string  `json:"created_at"`
}

type AddCartItemRequest struct {
	ProductID       int    `json:"product_id"`
	PackageID       int    `json:"package_id"`
	TravelDate      string `json:"travel_date"`
	Adults          int    `json:"adults"`
	Children        int    `json:"children"`
	SelectedOptions string `json:"selected_options"`
}

type CartSummary struct {
	Items       []CartItem `json:"items"`
	TotalAmount float64    `json:"total_amount"`
	Currency    string     `json:"currency"`
}

type CheckoutRequest struct {
	CouponCode string `json:"coupon_code"`
}
