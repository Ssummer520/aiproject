package domain

type Itinerary struct {
	ID        int             `json:"id"`
	UserID    string          `json:"user_id"`
	Title     string          `json:"title"`
	City      string          `json:"city"`
	StartDate string          `json:"start_date"`
	EndDate   string          `json:"end_date"`
	Guests    int             `json:"guests"`
	Budget    float64         `json:"budget"`
	Status    string          `json:"status"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
	Items     []ItineraryItem `json:"items"`
}

type ItineraryItem struct {
	ID            int     `json:"id"`
	ItineraryID   int     `json:"itinerary_id"`
	UserID        string  `json:"user_id"`
	DayIndex      int     `json:"day_index"`
	StartTime     string  `json:"start_time"`
	EndTime       string  `json:"end_time"`
	ItemType      string  `json:"item_type"`
	ProductID     int     `json:"product_id"`
	DestinationID int     `json:"destination_id"`
	Title         string  `json:"title"`
	Note          string  `json:"note"`
	EstimatedCost float64 `json:"estimated_cost"`
	SortOrder     int     `json:"sort_order"`
}

type CreateItineraryRequest struct {
	Title     string          `json:"title"`
	City      string          `json:"city"`
	StartDate string          `json:"start_date"`
	EndDate   string          `json:"end_date"`
	Guests    int             `json:"guests"`
	Budget    float64         `json:"budget"`
	Status    string          `json:"status"`
	Items     []ItineraryItem `json:"items"`
}

type AddItemRequest struct {
	DayIndex      int     `json:"day_index"`
	StartTime     string  `json:"start_time"`
	EndTime       string  `json:"end_time"`
	ItemType      string  `json:"item_type"`
	ProductID     int     `json:"product_id"`
	DestinationID int     `json:"destination_id"`
	Title         string  `json:"title"`
	Note          string  `json:"note"`
	EstimatedCost float64 `json:"estimated_cost"`
}

type GenerateRequest struct {
	Prompt    string  `json:"prompt"`
	City      string  `json:"city"`
	Days      int     `json:"days"`
	Guests    int     `json:"guests"`
	Budget    float64 `json:"budget"`
	StartDate string  `json:"start_date"`
	Save      bool    `json:"save"`
}
