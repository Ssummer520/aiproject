package models

type Destination struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	City        string   `json:"city"`
	Rating      float64  `json:"rating"`
	Cover       string   `json:"cover"`
	Tags        []string `json:"tags"`
	Lat         float64  `json:"lat"`
	Lng         float64  `json:"lng"`
	IsFavorite  bool     `json:"is_favorite"`
	Price       float64  `json:"price"`
	ReviewCount int      `json:"review_count"`
	BookedCount int      `json:"booked_count"`
	Description string   `json:"description"`
	Amenities   []string `json:"amenities"`
	HostName    string   `json:"host_name"`
	Policy      string   `json:"policy"`
}

type Deal struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"` // primary, secondary, accent
	Badge       string `json:"badge"`
	Expiry      string `json:"expiry"`
}

type Booking struct {
	ID            int     `json:"id"`
	UserID        string  `json:"user_id"`
	DestinationID int     `json:"destination_id"`
	Name          string  `json:"name"`
	City          string  `json:"city"`
	Cover         string  `json:"cover"`
	CheckIn       string  `json:"check_in"`
	CheckOut      string  `json:"check_out"`
	Guests        int     `json:"guests"`
	TotalPrice    float64 `json:"total_price"`
	Status        string  `json:"status"`
	CreatedAt     string  `json:"created_at"`
}
