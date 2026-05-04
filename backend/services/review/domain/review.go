package domain

type Scores struct {
	Quality   float64 `json:"quality"`
	Service   float64 `json:"service"`
	Value     float64 `json:"value"`
	Transport float64 `json:"transport"`
	Family    float64 `json:"family"`
}

type Review struct {
	ID            int      `json:"id"`
	UserID        string   `json:"user_id"`
	ProductID     int      `json:"product_id"`
	OrderID       int      `json:"order_id"`
	Rating        float64  `json:"rating"`
	Scores        Scores   `json:"scores"`
	Content       string   `json:"content"`
	Images        []string `json:"images"`
	Language      string   `json:"language"`
	Verified      bool     `json:"verified"`
	MerchantReply string   `json:"merchant_reply"`
	CreatedAt     string   `json:"created_at"`
}

type Summary struct {
	AverageRating float64 `json:"average_rating"`
	Total         int     `json:"total"`
	Quality       float64 `json:"quality"`
	Service       float64 `json:"service"`
	Value         float64 `json:"value"`
	Transport     float64 `json:"transport"`
	Family        float64 `json:"family"`
}

type ListResult struct {
	Summary Summary  `json:"summary"`
	Reviews []Review `json:"reviews"`
}

type CreateReviewRequest struct {
	OrderID  int      `json:"order_id"`
	Rating   float64  `json:"rating"`
	Scores   Scores   `json:"scores"`
	Content  string   `json:"content"`
	Images   []string `json:"images"`
	Language string   `json:"language"`
}
