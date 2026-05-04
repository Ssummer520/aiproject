package application

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"travel-api/services/platform/domain"
	productApp "travel-api/services/product/application"
)

func init() {
	os.Setenv("TRAVEL_DB_PATH", filepath.Join(os.TempDir(), fmt.Sprintf("chinatravel-platform-test-%d.db", os.Getpid())))
}

func futureDate() string { return time.Now().AddDate(0, 0, 1).Format("2006-01-02") }

func TestPlatformSnapshotSeedsOperatingData(t *testing.T) {
	productApp.NewProductService()
	service := NewPlatformService()
	snapshot, err := service.Snapshot("platform-user", "operator@example.com")
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if snapshot.Metrics.ProductCount == 0 || len(snapshot.Merchants) == 0 || len(snapshot.Inventory) == 0 || len(snapshot.CMS) == 0 {
		t.Fatalf("expected seeded platform data, got %#v", snapshot)
	}
	if snapshot.Profile.UserID != "platform-user" || snapshot.Profile.PointsBalance == 0 {
		t.Fatalf("expected seeded profile, got %#v", snapshot.Profile)
	}
}

func TestPlatformInventoryProfileAndCMS(t *testing.T) {
	productApp.NewProductService()
	service := NewPlatformService()
	updated, err := service.UpdateInventory(domain.InventoryUpdateRequest{PackageID: 1011, Date: futureDate(), Price: 99, Stock: 7, Status: "available"})
	if err != nil {
		t.Fatalf("update inventory: %v", err)
	}
	if updated.Stock != 7 || updated.Price != 99 || updated.ProductID == 0 {
		t.Fatalf("unexpected inventory update: %#v", updated)
	}
	if _, err := service.UpdateInventory(domain.InventoryUpdateRequest{PackageID: 0, Date: futureDate(), Stock: -1}); err != ErrInvalidPlatformRequest {
		t.Fatalf("expected invalid inventory request, got %v", err)
	}

	profile, err := service.UpdateProfile("member-user", domain.UserProfile{DisplayName: "Alan", Language: "zh", Currency: "CNY", MembershipLevel: "Gold", PointsBalance: 880})
	if err != nil {
		t.Fatalf("update profile: %v", err)
	}
	if profile.MembershipLevel != "Gold" || profile.PointsBalance != 880 {
		t.Fatalf("unexpected profile: %#v", profile)
	}

	article, err := service.CreateCMS(domain.CMSArticle{Slug: "arrival-checklist", Title: "Arrival checklist", Category: "guide", City: "China", Language: "en", Summary: "Before arrival", Content: "Passport, payment, eSIM", Status: "published"})
	if err != nil {
		t.Fatalf("create cms: %v", err)
	}
	if article.ID == 0 || article.Status != "published" {
		t.Fatalf("unexpected cms article: %#v", article)
	}
}

func TestPlatformRefundValidation(t *testing.T) {
	productApp.NewProductService()
	service := NewPlatformService()
	if _, err := service.CreateRefund(domain.CreateRefundRequest{UserID: "", OrderID: 1}); err != ErrInvalidPlatformRequest {
		t.Fatalf("expected invalid refund request, got %v", err)
	}
}
