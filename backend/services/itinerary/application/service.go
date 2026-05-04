package application

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"travel-api/services/itinerary/domain"
	"travel-api/services/itinerary/infrastructure"
	productApp "travel-api/services/product/application"
	productDomain "travel-api/services/product/domain"
)

var (
	ErrInvalidItineraryRequest = errors.New("invalid_itinerary_request")
	ErrItineraryNotFound       = errors.New("itinerary_not_found")
)

type ItineraryService struct {
	repo           *infrastructure.SQLiteItineraryRepo
	productService *productApp.ProductService
}

func NewItineraryService(productService *productApp.ProductService) *ItineraryService {
	return &ItineraryService{repo: infrastructure.NewSQLiteItineraryRepo(), productService: productService}
}

func (s *ItineraryService) List(userID string) ([]domain.Itinerary, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrInvalidItineraryRequest
	}
	return s.repo.List(userID)
}

func (s *ItineraryService) Get(userID string, itineraryID int) (domain.Itinerary, error) {
	if strings.TrimSpace(userID) == "" || itineraryID <= 0 {
		return domain.Itinerary{}, ErrInvalidItineraryRequest
	}
	itinerary, ok, err := s.repo.Get(userID, itineraryID)
	if err != nil {
		return domain.Itinerary{}, err
	}
	if !ok {
		return domain.Itinerary{}, ErrItineraryNotFound
	}
	return itinerary, nil
}

func (s *ItineraryService) Create(userID string, req domain.CreateItineraryRequest) (domain.Itinerary, error) {
	if strings.TrimSpace(userID) == "" || strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.City) == "" {
		return domain.Itinerary{}, ErrInvalidItineraryRequest
	}
	if req.Guests <= 0 {
		req.Guests = 1
	}
	if req.StartDate == "" {
		req.StartDate = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	}
	if req.EndDate == "" {
		req.EndDate = req.StartDate
	}
	itinerary := domain.Itinerary{Title: strings.TrimSpace(req.Title), City: strings.TrimSpace(req.City), StartDate: req.StartDate, EndDate: req.EndDate, Guests: req.Guests, Budget: req.Budget, Status: req.Status, Items: req.Items}
	return s.repo.Create(userID, itinerary)
}

func (s *ItineraryService) AddItem(userID string, itineraryID int, req domain.AddItemRequest) (domain.Itinerary, error) {
	if strings.TrimSpace(userID) == "" || itineraryID <= 0 {
		return domain.Itinerary{}, ErrInvalidItineraryRequest
	}
	item := domain.ItineraryItem{DayIndex: req.DayIndex, StartTime: req.StartTime, EndTime: req.EndTime, ItemType: req.ItemType, ProductID: req.ProductID, DestinationID: req.DestinationID, Title: strings.TrimSpace(req.Title), Note: strings.TrimSpace(req.Note), EstimatedCost: req.EstimatedCost}
	if item.Title == "" && item.ProductID > 0 {
		product, err := s.productService.Get(item.ProductID)
		if err == nil {
			item.Title = product.Name
			item.EstimatedCost = product.BasePrice
			if item.ItemType == "" {
				item.ItemType = "product"
			}
		}
	}
	if item.Title == "" {
		return domain.Itinerary{}, ErrInvalidItineraryRequest
	}
	itinerary, ok, err := s.repo.AddItem(userID, itineraryID, item)
	if err != nil {
		return domain.Itinerary{}, err
	}
	if !ok {
		return domain.Itinerary{}, ErrItineraryNotFound
	}
	return itinerary, nil
}

func (s *ItineraryService) MoveItem(userID string, itineraryID int, itemID int, direction string) (domain.Itinerary, error) {
	if strings.TrimSpace(userID) == "" || itineraryID <= 0 || itemID <= 0 || (direction != "up" && direction != "down") {
		return domain.Itinerary{}, ErrInvalidItineraryRequest
	}
	itinerary, ok, err := s.repo.MoveItem(userID, itineraryID, itemID, direction)
	if err != nil {
		return domain.Itinerary{}, err
	}
	if !ok {
		return domain.Itinerary{}, ErrItineraryNotFound
	}
	return itinerary, nil
}

func (s *ItineraryService) Generate(userID string, req domain.GenerateRequest) (domain.Itinerary, error) {
	if strings.TrimSpace(userID) == "" {
		return domain.Itinerary{}, ErrInvalidItineraryRequest
	}
	city := strings.TrimSpace(req.City)
	if city == "" {
		city = inferCity(req.Prompt)
	}
	if city == "" {
		city = "Hangzhou"
	}
	days := req.Days
	if days <= 0 {
		days = inferDays(req.Prompt)
	}
	if days <= 0 {
		days = 2
	}
	if days > 5 {
		days = 5
	}
	guests := req.Guests
	if guests <= 0 {
		guests = 2
	}
	startDate := req.StartDate
	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	}
	start, _ := time.Parse("2006-01-02", startDate)
	endDate := start.AddDate(0, 0, days-1).Format("2006-01-02")

	products, _ := s.productService.Search(productDomain.SearchFilters{City: city, Sort: "booked"})
	items := make([]domain.ItineraryItem, 0, days*3)
	slots := []struct{ start, end, label string }{{"09:00", "11:30", "Morning"}, {"14:00", "17:00", "Afternoon"}, {"19:00", "21:00", "Evening"}}
	for day := 1; day <= days; day++ {
		for slotIndex, slot := range slots {
			product := pickProduct(products.Results, day, slotIndex)
			item := domain.ItineraryItem{DayIndex: day, StartTime: slot.start, EndTime: slot.end, ItemType: "note", Title: fmt.Sprintf("%s %s inspiration", city, slot.label), Note: "AI generated planning block", EstimatedCost: 0}
			if product.ID > 0 {
				item.ItemType = "product"
				item.ProductID = product.ID
				item.DestinationID = product.DestinationID
				item.Title = product.Name
				item.Note = product.Subtitle
				item.EstimatedCost = product.BasePrice
			}
			items = append(items, item)
		}
	}
	itinerary := domain.Itinerary{Title: fmt.Sprintf("%s %d-day AI itinerary", city, days), City: city, StartDate: startDate, EndDate: endDate, Guests: guests, Budget: req.Budget, Status: "draft", Items: items}
	if req.Save {
		return s.repo.Create(userID, itinerary)
	}
	itinerary.UserID = userID
	return itinerary, nil
}

func inferCity(prompt string) string {
	lower := strings.ToLower(prompt)
	for _, city := range []string{"Hangzhou", "Shanghai", "Beijing", "Chengdu", "Xi'an"} {
		if strings.Contains(lower, strings.ToLower(city)) || strings.Contains(prompt, map[string]string{"Hangzhou": "杭州", "Shanghai": "上海", "Beijing": "北京", "Chengdu": "成都", "Xi'an": "西安"}[city]) {
			return city
		}
	}
	return ""
}

func inferDays(prompt string) int {
	for day := 1; day <= 5; day++ {
		if strings.Contains(prompt, fmt.Sprintf("%d day", day)) || strings.Contains(prompt, fmt.Sprintf("%d-day", day)) || strings.Contains(prompt, fmt.Sprintf("%d日", day)) || strings.Contains(prompt, fmt.Sprintf("%d天", day)) {
			return day
		}
	}
	return 0
}

func pickProduct(products []productDomain.Product, day int, slotIndex int) productDomain.Product {
	if len(products) == 0 {
		return productDomain.Product{}
	}
	return products[((day-1)*3+slotIndex)%len(products)]
}
