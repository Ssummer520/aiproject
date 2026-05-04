package application

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"travel-api/services/itinerary/domain"
	productApp "travel-api/services/product/application"
)

func init() {
	os.Setenv("TRAVEL_DB_PATH", filepath.Join(os.TempDir(), fmt.Sprintf("chinatravel-itinerary-test-%d.db", os.Getpid())))
}

func testDate(offset int) string { return time.Now().AddDate(0, 0, offset).Format("2006-01-02") }

func TestItineraryServiceCreateAddMoveAndList(t *testing.T) {
	service := NewItineraryService(productApp.NewProductService())
	itinerary, err := service.Create("user-itinerary", domain.CreateItineraryRequest{Title: "Hangzhou family trip", City: "Hangzhou", StartDate: testDate(3), EndDate: testDate(4), Guests: 3, Budget: 1200})
	if err != nil {
		t.Fatalf("create itinerary: %v", err)
	}
	if itinerary.ID == 0 || itinerary.Status != "draft" {
		t.Fatalf("unexpected itinerary: %#v", itinerary)
	}

	itinerary, err = service.AddItem("user-itinerary", itinerary.ID, domain.AddItemRequest{DayIndex: 1, StartTime: "09:00", EndTime: "11:00", ProductID: 101})
	if err != nil {
		t.Fatalf("add product item: %v", err)
	}
	if len(itinerary.Items) != 1 || itinerary.Items[0].Title == "" || itinerary.Items[0].ItemType != "product" {
		t.Fatalf("unexpected item: %#v", itinerary.Items)
	}

	itinerary, err = service.AddItem("user-itinerary", itinerary.ID, domain.AddItemRequest{DayIndex: 1, StartTime: "14:00", EndTime: "15:00", ItemType: "note", Title: "Cafe break"})
	if err != nil {
		t.Fatalf("add note item: %v", err)
	}
	moved, err := service.MoveItem("user-itinerary", itinerary.ID, itinerary.Items[1].ID, "up")
	if err != nil {
		t.Fatalf("move item: %v", err)
	}
	if len(moved.Items) != 2 || moved.Items[0].Title != "Cafe break" {
		t.Fatalf("expected note moved up, got %#v", moved.Items)
	}

	items, err := service.List("user-itinerary")
	if err != nil || len(items) != 1 || len(items[0].Items) != 2 {
		t.Fatalf("list itineraries: items=%#v err=%v", items, err)
	}
}

func TestItineraryServiceGenerateAISave(t *testing.T) {
	service := NewItineraryService(productApp.NewProductService())
	itinerary, err := service.Generate("ai-user", domain.GenerateRequest{Prompt: "杭州 2天 亲子 低预算", Budget: 1000, Save: true})
	if err != nil {
		t.Fatalf("generate itinerary: %v", err)
	}
	if itinerary.City != "Hangzhou" || len(itinerary.Items) != 6 {
		t.Fatalf("expected Hangzhou 2-day itinerary, got %#v", itinerary)
	}
	if itinerary.ID == 0 {
		t.Fatalf("saved AI itinerary should have id")
	}
}

func TestItineraryServiceRejectsInvalidRequests(t *testing.T) {
	service := NewItineraryService(productApp.NewProductService())
	if _, err := service.Create("", domain.CreateItineraryRequest{Title: "A", City: "Hangzhou"}); err != ErrInvalidItineraryRequest {
		t.Fatalf("expected invalid user, got %v", err)
	}
	if _, err := service.AddItem("user", 999, domain.AddItemRequest{Title: "Missing"}); err != ErrItineraryNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}
