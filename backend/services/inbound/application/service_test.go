package application

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"travel-api/services/inbound/domain"
	productApp "travel-api/services/product/application"
)

func init() {
	os.Setenv("TRAVEL_DB_PATH", filepath.Join(os.TempDir(), fmt.Sprintf("chinatravel-inbound-test-%d.db", os.Getpid())))
}

func TestInboundSnapshotSeedsChinaTravelDifferentiators(t *testing.T) {
	productApp.NewProductService()
	service := NewInboundService()
	snapshot, err := service.Snapshot()
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if len(snapshot.Toolkit) < 4 || len(snapshot.Rails) == 0 || len(snapshot.Transfers) == 0 || len(snapshot.Passes) == 0 || len(snapshot.Guides) == 0 {
		t.Fatalf("expected inbound snapshot data, got %#v", snapshot)
	}
	if snapshot.Toolkit[0].ProductID == 0 {
		t.Fatalf("expected toolkit item linked to bookable product: %#v", snapshot.Toolkit[0])
	}
}

func TestInboundCityGuideAndConcierge(t *testing.T) {
	productApp.NewProductService()
	service := NewInboundService()
	guide, ok, err := service.Guide("Hangzhou")
	if err != nil || !ok {
		t.Fatalf("guide: ok=%v err=%v", ok, err)
	}
	if !strings.Contains(guide.Transport, "airport") && !strings.Contains(guide.Transport, "Airport") {
		t.Fatalf("expected airport/transport guidance, got %#v", guide)
	}

	answer, err := service.Concierge(domain.ConciergeRequest{Prompt: "Need driver and eSIM in Shanghai", City: "Shanghai", Budget: 1200, Days: 2})
	if err != nil {
		t.Fatalf("concierge: %v", err)
	}
	if answer.City != "Shanghai" || answer.ChineseMessage == "" || len(answer.PracticalChecklist) == 0 {
		t.Fatalf("unexpected concierge answer: %#v", answer)
	}
}
