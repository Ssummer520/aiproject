package domain

type Product struct {
	ID             int            `json:"id"`
	DestinationID  int            `json:"destination_id"`
	City           string         `json:"city"`
	Category       string         `json:"category"`
	Type           string         `json:"type"`
	Name           string         `json:"name"`
	Subtitle       string         `json:"subtitle"`
	Description    string         `json:"description"`
	Cover          string         `json:"cover"`
	Images         []string       `json:"images"`
	Tags           []string       `json:"tags"`
	Rating         float64        `json:"rating"`
	ReviewCount    int            `json:"review_count"`
	BookedCount    int            `json:"booked_count"`
	BasePrice      float64        `json:"base_price"`
	Currency       string         `json:"currency"`
	InstantConfirm bool           `json:"instant_confirm"`
	FreeCancel     bool           `json:"free_cancel"`
	Duration       string         `json:"duration"`
	MeetingPoint   string         `json:"meeting_point"`
	Included       []string       `json:"included"`
	Excluded       []string       `json:"excluded"`
	Usage          string         `json:"usage"`
	Policy         string         `json:"policy"`
	Status         string         `json:"status"`
	Packages       []Package      `json:"packages,omitempty"`
	Availability   []Availability `json:"availability,omitempty"`
}

type Package struct {
	ID            int     `json:"id"`
	ProductID     int     `json:"product_id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	OriginalPrice float64 `json:"original_price"`
	Unit          string  `json:"unit"`
	MinQuantity   int     `json:"min_quantity"`
	MaxQuantity   int     `json:"max_quantity"`
	AgeRule       string  `json:"age_rule"`
	RefundPolicy  string  `json:"refund_policy"`
	ConfirmType   string  `json:"confirm_type"`
	VoucherType   string  `json:"voucher_type"`
}

type Availability struct {
	ID        int     `json:"id"`
	ProductID int     `json:"product_id"`
	PackageID int     `json:"package_id"`
	Date      string  `json:"date"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
	Status    string  `json:"status"`
}

type SearchFilters struct {
	Query          string
	City           string
	Category       string
	Type           string
	MinPrice       float64
	MaxPrice       float64
	RatingMin      float64
	InstantConfirm *bool
	FreeCancel     *bool
	Sort           string
}

type SearchResult struct {
	Results []Product `json:"results"`
	Total   int       `json:"total"`
}
