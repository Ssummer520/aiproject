package application

import (
	"errors"
	"testing"

	"travel-api/services/user/domain"
	"travel-api/services/user/infrastructure"
)

type fakeUserRepo struct {
	profile   domain.UserProfile
	travelers map[int]domain.TravelerProfile
}

func (r *fakeUserRepo) EnsureUserDefaults(userID, email string) error { return nil }
func (r *fakeUserRepo) GetProfile(userID, email string) (domain.UserProfile, error) {
	return r.profile, nil
}
func (r *fakeUserRepo) UpsertProfile(userID string, profile domain.UserProfile) (domain.UserProfile, error) {
	profile.UserID = userID
	r.profile = profile
	return profile, nil
}
func (r *fakeUserRepo) ListTravelers(userID string) ([]domain.TravelerProfile, error) {
	items := []domain.TravelerProfile{}
	for _, traveler := range r.travelers {
		items = append(items, traveler)
	}
	return items, nil
}
func (r *fakeUserRepo) GetTraveler(userID string, travelerID int) (domain.TravelerProfile, bool, error) {
	traveler, ok := r.travelers[travelerID]
	return traveler, ok, nil
}
func (r *fakeUserRepo) CreateTraveler(userID string, input domain.TravelerProfileInput) (domain.TravelerProfile, error) {
	for _, traveler := range r.travelers {
		if traveler.DocumentNoMasked == infrastructure.MaskDocument(input.DocumentNo) {
			return domain.TravelerProfile{}, infrastructure.ErrDuplicateDocument
		}
	}
	id := len(r.travelers) + 1
	traveler := domain.TravelerProfile{ID: id, UserID: userID, Name: input.Name, DocumentType: input.DocumentType, DocumentNoMasked: infrastructure.MaskDocument(input.DocumentNo), IsDefault: input.IsDefault}
	r.travelers[id] = traveler
	return traveler, nil
}
func (r *fakeUserRepo) UpdateTraveler(userID string, travelerID int, input domain.TravelerProfileInput, hasDocument bool) (domain.TravelerProfile, error) {
	traveler, ok := r.travelers[travelerID]
	if !ok {
		return domain.TravelerProfile{}, errors.New("missing")
	}
	traveler.Name = input.Name
	traveler.IsDefault = input.IsDefault
	if hasDocument {
		traveler.DocumentType = input.DocumentType
		traveler.DocumentNoMasked = infrastructure.MaskDocument(input.DocumentNo)
	}
	r.travelers[travelerID] = traveler
	return traveler, nil
}
func (r *fakeUserRepo) DeleteTraveler(userID string, travelerID int) error { return nil }
func (r *fakeUserRepo) SetDefaultTraveler(userID string, travelerID int) (domain.TravelerProfile, error) {
	traveler := r.travelers[travelerID]
	traveler.IsDefault = true
	r.travelers[travelerID] = traveler
	return traveler, nil
}
func (r *fakeUserRepo) GetMembership(userID string) (domain.Membership, error) {
	return domain.Membership{UserID: userID, Level: "Silver", PointsBalance: 300}, nil
}
func (r *fakeUserRepo) GetRoles(userID string) ([]domain.Role, error) {
	return []domain.Role{{Code: "customer", Name: "Customer"}}, nil
}
func (r *fakeUserRepo) BuildOrderSnapshots(userID string, req domain.TravelerSnapshotRequest) ([]domain.OrderTravelerSnapshot, error) {
	snapshots := []domain.OrderTravelerSnapshot{}
	for _, id := range req.TravelerIDs {
		traveler, ok := r.travelers[id]
		if !ok {
			return nil, errors.New("missing")
		}
		snapshots = append(snapshots, domain.OrderTravelerSnapshot{SourceTravelerID: traveler.ID, Name: traveler.Name, DocumentType: traveler.DocumentType, DocumentNoMasked: traveler.DocumentNoMasked})
	}
	for _, traveler := range req.Travelers {
		snapshots = append(snapshots, domain.OrderTravelerSnapshot{Name: traveler.Name, DocumentType: traveler.DocumentType, DocumentNoMasked: infrastructure.MaskDocument(traveler.DocumentNo)})
	}
	return snapshots, nil
}

func newFakeService() *UserService {
	return NewUserServiceWithRepo(&fakeUserRepo{travelers: map[int]domain.TravelerProfile{}})
}

func TestUserServiceNormalizesProfileDefaults(t *testing.T) {
	service := newFakeService()
	profile, err := service.UpsertProfile("user-1", domain.UserProfile{
		DisplayName:         " Alan ",
		Currency:            "usd",
		DietaryRestrictions: []string{"halal", "halal", ""},
	})
	if err != nil {
		t.Fatalf("upsert profile: %v", err)
	}
	if profile.DisplayName != "Alan" || profile.Language != "en" || profile.Currency != "USD" {
		t.Fatalf("unexpected normalized profile: %#v", profile)
	}
	if len(profile.DietaryRestrictions) != 1 || profile.DietaryRestrictions[0] != "halal" {
		t.Fatalf("expected de-duplicated dietary list, got %#v", profile.DietaryRestrictions)
	}
}

func TestUserServiceTravelerValidationAndMasking(t *testing.T) {
	service := newFakeService()
	_, err := service.CreateTraveler("user-1", domain.TravelerProfileInput{Name: "Alan"})
	if err != ErrInvalidUserRequest {
		t.Fatalf("expected invalid request, got %v", err)
	}
	traveler, err := service.CreateTraveler("user-1", domain.TravelerProfileInput{Name: "Alan", DocumentType: "passport", DocumentNo: "E12345678", IsDefault: true})
	if err != nil {
		t.Fatalf("create traveler: %v", err)
	}
	if traveler.DocumentType != "PASSPORT" || traveler.DocumentNoMasked != "E1***5678" || !traveler.IsDefault {
		t.Fatalf("unexpected traveler: %#v", traveler)
	}
	_, err = service.CreateTraveler("user-1", domain.TravelerProfileInput{Name: "Copy", DocumentType: "PASSPORT", DocumentNo: "E12345678"})
	if err != ErrDocumentDuplicate {
		t.Fatalf("expected duplicate document, got %v", err)
	}
}

func TestUserServiceBuildsOrderSnapshotsWithoutFullDocumentNumber(t *testing.T) {
	service := newFakeService()
	created, err := service.CreateTraveler("user-1", domain.TravelerProfileInput{Name: "Alan", DocumentType: "ID", DocumentNo: "330100199001011234"})
	if err != nil {
		t.Fatalf("create traveler: %v", err)
	}
	snapshots, err := service.BuildOrderSnapshots("user-1", domain.TravelerSnapshotRequest{
		TravelerIDs: []int{created.ID},
		Travelers:   []domain.TravelerProfileInput{{Name: "Manual", DocumentType: "PASSPORT", DocumentNo: "P99887766"}},
	})
	if err != nil {
		t.Fatalf("build snapshots: %v", err)
	}
	if len(snapshots) != 2 {
		t.Fatalf("expected two snapshots, got %#v", snapshots)
	}
	if snapshots[0].DocumentNoMasked == "330100199001011234" || snapshots[1].DocumentNoMasked == "P99887766" {
		t.Fatalf("snapshot leaked full document numbers: %#v", snapshots)
	}
}
