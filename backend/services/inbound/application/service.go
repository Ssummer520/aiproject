package application

import (
	"errors"
	"strings"

	"travel-api/services/inbound/domain"
	"travel-api/services/inbound/infrastructure"
)

var ErrInvalidInboundRequest = errors.New("invalid_inbound_request")

type InboundService struct {
	repo *infrastructure.SQLiteInboundRepo
}

func NewInboundService() *InboundService {
	return &InboundService{repo: infrastructure.NewSQLiteInboundRepo()}
}
func (s *InboundService) Snapshot() (domain.InboundSnapshot, error) { return s.repo.Snapshot() }
func (s *InboundService) Guide(city string) (domain.CityGuide, bool, error) {
	if strings.TrimSpace(city) == "" {
		return domain.CityGuide{}, false, ErrInvalidInboundRequest
	}
	return s.repo.Guide(city)
}
func (s *InboundService) Concierge(req domain.ConciergeRequest) (domain.ConciergeResponse, error) {
	prompt := strings.ToLower(strings.TrimSpace(req.Prompt))
	city := strings.TrimSpace(req.City)
	if city == "" {
		city = inferCity(prompt)
	}
	if city == "" {
		city = "Hangzhou"
	}
	days := req.Days
	if days <= 0 {
		days = 2
	}
	budget := req.Budget
	if budget <= 0 {
		budget = 1000
	}
	products := []int{109, 111, 113}
	if strings.EqualFold(city, "Shanghai") {
		products = []int{109, 110, 114}
	}
	if strings.Contains(prompt, "driver") || strings.Contains(prompt, "司机") {
		products = append(products, 111)
	}
	return domain.ConciergeResponse{
		City:               city,
		Summary:            "A China inbound plan covering arrival, connectivity, local transport, reservations and bookable city experiences.",
		ChineseMessage:     chineseMessage(city, prompt),
		WeatherAdjustment:  "If rain or heat affects outdoor plans, move museum/tea/food activities to the afternoon and keep cruises for clearer windows.",
		BudgetSuggestion:   budgetSuggestion(budget, days),
		TransportPlan:      transportPlan(city),
		RecommendedProduct: products,
		PracticalChecklist: []string{"Install eSIM before departure", "Save hotel address in Chinese", "Carry passport for tickets/trains", "Keep cash/card fallback", "Book timed attractions early"},
	}, nil
}

func inferCity(prompt string) string {
	for _, city := range []string{"hangzhou", "shanghai", "beijing", "chengdu", "xian"} {
		if strings.Contains(prompt, city) {
			return strings.Title(city)
		}
	}
	if strings.Contains(prompt, "杭州") {
		return "Hangzhou"
	}
	if strings.Contains(prompt, "上海") {
		return "Shanghai"
	}
	if strings.Contains(prompt, "北京") {
		return "Beijing"
	}
	return ""
}
func chineseMessage(city, prompt string) string {
	if strings.Contains(prompt, "driver") || strings.Contains(prompt, "taxi") || strings.Contains(prompt, "司机") {
		return "师傅您好，请带我去这个地址。我不会说中文，如有问题请通过短信联系我。"
	}
	return "您好，我已经预订了行程/门票，请帮我确认预约信息。我会出示护照和电子凭证。"
}
func budgetSuggestion(budget, days int) string {
	perDay := budget / days
	if perDay < 400 {
		return "Use metro, one paid attraction per day, and reserve budget for connectivity and airport transfer."
	}
	if perDay < 900 {
		return "Mix one city pass or guided activity with self-guided meals and metro/taxi transfers."
	}
	return "Add private transfer, city pass, and a guided experience for a smoother inbound trip."
}
func transportPlan(city string) string {
	switch strings.ToLower(city) {
	case "hangzhou":
		return "Book airport transfer for arrival, then use metro/taxi around West Lake; high-speed rail from Shanghai Hongqiao to Hangzhou East is the fastest intercity option."
	case "shanghai":
		return "Use metro inside the city, keep Chinese addresses for taxis, and choose Hongqiao for rail connections."
	case "beijing":
		return "Use airport transfer on arrival, subway for central routes, and pre-book long-distance rail from Beijing South."
	default:
		return "Use eSIM + saved Chinese addresses, combine metro with licensed transfers for arrival and late-night routes."
	}
}
