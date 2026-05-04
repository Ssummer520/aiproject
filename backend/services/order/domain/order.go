package domain

type Order struct {
	ID            int         `json:"id"`
	UserID        string      `json:"user_id"`
	Status        string      `json:"status"`
	PaymentStatus string      `json:"payment_status"`
	TotalAmount   float64     `json:"total_amount"`
	Currency      string      `json:"currency"`
	ContactName   string      `json:"contact_name"`
	ContactEmail  string      `json:"contact_email"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
	CancelledAt   string      `json:"cancelled_at,omitempty"`
	Items         []OrderItem `json:"items"`
}

type OrderItem struct {
	ID          int     `json:"id"`
	OrderID     int     `json:"order_id"`
	UserID      string  `json:"user_id"`
	ProductID   int     `json:"product_id"`
	PackageID   int     `json:"package_id"`
	ProductName string  `json:"product_name"`
	PackageName string  `json:"package_name"`
	City        string  `json:"city"`
	Cover       string  `json:"cover"`
	TravelDate  string  `json:"travel_date"`
	Adults      int     `json:"adults"`
	Children    int     `json:"children"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Subtotal    float64 `json:"subtotal"`
}

type CreateOrderRequest struct {
	ProductID    int    `json:"product_id"`
	PackageID    int    `json:"package_id"`
	TravelDate   string `json:"travel_date"`
	Adults       int    `json:"adults"`
	Children     int    `json:"children"`
	ContactName  string `json:"contact_name"`
	ContactEmail string `json:"contact_email"`
}
