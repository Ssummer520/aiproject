package domain

type UserProfile struct {
	UserID              string   `json:"user_id"`
	DisplayName         string   `json:"display_name"`
	Avatar              string   `json:"avatar"`
	Phone               string   `json:"phone"`
	Nationality         string   `json:"nationality"`
	Language            string   `json:"language"`
	Currency            string   `json:"currency"`
	TravelStyle         string   `json:"travel_style"`
	BudgetLevel         string   `json:"budget_level"`
	FamilyType          string   `json:"family_type"`
	DietaryRestrictions []string `json:"dietary_restrictions"`
	AccessibilityNeeds  []string `json:"accessibility_needs"`
	UpdatedAt           string   `json:"updated_at"`
}

type TravelerProfile struct {
	ID               int    `json:"id"`
	UserID           string `json:"user_id"`
	Name             string `json:"name"`
	Gender           string `json:"gender"`
	BirthDate        string `json:"birth_date"`
	DocumentType     string `json:"document_type"`
	DocumentNoMasked string `json:"document_no_masked"`
	Nationality      string `json:"nationality"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	IsDefault        bool   `json:"is_default"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

type TravelerProfileInput struct {
	Name         string `json:"name"`
	Gender       string `json:"gender"`
	BirthDate    string `json:"birth_date"`
	DocumentType string `json:"document_type"`
	DocumentNo   string `json:"document_no"`
	Nationality  string `json:"nationality"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	IsDefault    bool   `json:"is_default"`
}

type Membership struct {
	UserID        string   `json:"user_id"`
	Level         string   `json:"level"`
	PointsBalance int      `json:"points_balance"`
	ValidUntil    string   `json:"valid_until"`
	Benefits      []string `json:"benefits"`
	UpdatedAt     string   `json:"updated_at"`
}

type Role struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type TravelerSnapshotRequest struct {
	TravelerIDs []int                  `json:"traveler_ids"`
	Travelers   []TravelerProfileInput `json:"travelers"`
}

type OrderTravelerSnapshot struct {
	ID               int    `json:"id"`
	OrderID          int    `json:"order_id"`
	UserID           string `json:"user_id"`
	SourceTravelerID int    `json:"source_traveler_id,omitempty"`
	Name             string `json:"name"`
	Gender           string `json:"gender"`
	BirthDate        string `json:"birth_date"`
	DocumentType     string `json:"document_type"`
	DocumentNoMasked string `json:"document_no_masked"`
	Nationality      string `json:"nationality"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
}
