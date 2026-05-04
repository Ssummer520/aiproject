package domain

type ToolkitItem struct {
	ID          int      `json:"id"`
	Key         string   `json:"key"`
	Title       string   `json:"title"`
	TitleZh     string   `json:"title_zh"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
	CTA         string   `json:"cta"`
	ProductID   int      `json:"product_id"`
}

type RailRoute struct {
	ID        int     `json:"id"`
	From      string  `json:"from"`
	To        string  `json:"to"`
	Duration  string  `json:"duration"`
	Frequency string  `json:"frequency"`
	PriceFrom float64 `json:"price_from"`
	Tip       string  `json:"tip"`
	ProductID int     `json:"product_id"`
}

type TransferOption struct {
	ID        int     `json:"id"`
	City      string  `json:"city"`
	From      string  `json:"from"`
	To        string  `json:"to"`
	Vehicle   string  `json:"vehicle"`
	PriceFrom float64 `json:"price_from"`
	DriverTip string  `json:"driver_tip"`
	ProductID int     `json:"product_id"`
}

type CityPass struct {
	ID          int      `json:"id"`
	City        string   `json:"city"`
	Name        string   `json:"name"`
	Duration    string   `json:"duration"`
	Includes    []string `json:"includes"`
	PriceFrom   float64  `json:"price_from"`
	ProductID   int      `json:"product_id"`
	AISuggested bool     `json:"ai_suggested"`
}

type CityGuide struct {
	City         string   `json:"city"`
	BestSeason   string   `json:"best_season"`
	Weather      string   `json:"weather"`
	Transport    string   `json:"transport"`
	Payment      string   `json:"payment"`
	Connectivity string   `json:"connectivity"`
	Reservation  string   `json:"reservation"`
	LanguageTips []string `json:"language_tips"`
	SafetyTips   []string `json:"safety_tips"`
	DietaryTips  []string `json:"dietary_tips"`
	FamilyTips   []string `json:"family_tips"`
}

type ConciergeRequest struct {
	Prompt string `json:"prompt"`
	City   string `json:"city"`
	Budget int    `json:"budget"`
	Days   int    `json:"days"`
}

type ConciergeResponse struct {
	City               string   `json:"city"`
	Summary            string   `json:"summary"`
	ChineseMessage     string   `json:"chinese_message"`
	WeatherAdjustment  string   `json:"weather_adjustment"`
	BudgetSuggestion   string   `json:"budget_suggestion"`
	TransportPlan      string   `json:"transport_plan"`
	RecommendedProduct []int    `json:"recommended_product_ids"`
	PracticalChecklist []string `json:"practical_checklist"`
}

type InboundSnapshot struct {
	Toolkit   []ToolkitItem    `json:"toolkit"`
	Rails     []RailRoute      `json:"rails"`
	Transfers []TransferOption `json:"transfers"`
	Passes    []CityPass       `json:"passes"`
	Guides    []CityGuide      `json:"guides"`
}
