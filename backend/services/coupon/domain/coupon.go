package domain

type Coupon struct {
	ID            int     `json:"id"`
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	DiscountType  string  `json:"discount_type"`
	DiscountValue float64 `json:"discount_value"`
	MinSpend      float64 `json:"min_spend"`
	ValidFrom     string  `json:"valid_from"`
	ValidTo       string  `json:"valid_to"`
	UsageLimit    int     `json:"usage_limit"`
	Status        string  `json:"status"`
}

type ValidationResult struct {
	Valid          bool    `json:"valid"`
	Coupon         Coupon  `json:"coupon,omitempty"`
	OriginalAmount float64 `json:"original_amount"`
	DiscountAmount float64 `json:"discount_amount"`
	FinalAmount    float64 `json:"final_amount"`
	Error          string  `json:"error,omitempty"`
}

type ValidateRequest struct {
	Code   string  `json:"code"`
	Amount float64 `json:"amount"`
}
