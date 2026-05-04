package domain

type Merchant struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	ContactEmail string  `json:"contact_email"`
	Phone        string  `json:"phone"`
	City         string  `json:"city"`
	Status       string  `json:"status"`
	Rating       float64 `json:"rating"`
	CreatedAt    string  `json:"created_at"`
}

type InventoryItem struct {
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	MerchantID  int     `json:"merchant_id"`
	Merchant    string  `json:"merchant"`
	PackageID   int     `json:"package_id"`
	PackageName string  `json:"package_name"`
	Date        string  `json:"date"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Status      string  `json:"status"`
}

type InventoryUpdateRequest struct {
	PackageID int     `json:"package_id"`
	Date      string  `json:"date"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
	Status    string  `json:"status"`
}

type PlatformOrder struct {
	ID            int     `json:"id"`
	UserID        string  `json:"user_id"`
	Status        string  `json:"status"`
	PaymentStatus string  `json:"payment_status"`
	TotalAmount   float64 `json:"total_amount"`
	Currency      string  `json:"currency"`
	ProductName   string  `json:"product_name"`
	PackageName   string  `json:"package_name"`
	City          string  `json:"city"`
	TravelDate    string  `json:"travel_date"`
	CreatedAt     string  `json:"created_at"`
}

type RefundRequest struct {
	ID           int     `json:"id"`
	UserID       string  `json:"user_id"`
	OrderID      int     `json:"order_id"`
	Reason       string  `json:"reason"`
	RefundAmount float64 `json:"refund_amount"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type CreateRefundRequest struct {
	UserID  string `json:"user_id"`
	OrderID int    `json:"order_id"`
	Reason  string `json:"reason"`
}

type UserProfile struct {
	UserID              string `json:"user_id"`
	DisplayName         string `json:"display_name"`
	Avatar              string `json:"avatar"`
	Phone               string `json:"phone"`
	Nationality         string `json:"nationality"`
	PassportName        string `json:"passport_name"`
	Language            string `json:"language"`
	Currency            string `json:"currency"`
	TravelPreferences   string `json:"travel_preferences"`
	DietaryRestrictions string `json:"dietary_restrictions"`
	MembershipLevel     string `json:"membership_level"`
	PointsBalance       int    `json:"points_balance"`
	UpdatedAt           string `json:"updated_at"`
}

type CMSArticle struct {
	ID        int    `json:"id"`
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	Category  string `json:"category"`
	City      string `json:"city"`
	Language  string `json:"language"`
	Summary   string `json:"summary"`
	Content   string `json:"content"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
}

type DashboardMetrics struct {
	GMV                 float64      `json:"gmv"`
	OrderCount          int          `json:"order_count"`
	ProductCount        int          `json:"product_count"`
	MerchantCount       int          `json:"merchant_count"`
	RefundRate          float64      `json:"refund_rate"`
	AIItineraryCount    int          `json:"ai_itinerary_count"`
	CartItemCount       int          `json:"cart_item_count"`
	ReviewCount         int          `json:"review_count"`
	MemberCount         int          `json:"member_count"`
	PublishedCMSCount   int          `json:"published_cms_count"`
	CityHeat            []MetricPair `json:"city_heat"`
	OperationalWarnings []string     `json:"operational_warnings"`
}

type MetricPair struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

type Snapshot struct {
	Metrics   DashboardMetrics `json:"metrics"`
	Merchants []Merchant       `json:"merchants"`
	Inventory []InventoryItem  `json:"inventory"`
	Orders    []PlatformOrder  `json:"orders"`
	Refunds   []RefundRequest  `json:"refunds"`
	CMS       []CMSArticle     `json:"cms"`
	Profile   UserProfile      `json:"profile"`
}
