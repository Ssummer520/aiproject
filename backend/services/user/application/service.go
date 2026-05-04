package application

import (
	"database/sql"
	"errors"
	"strings"

	orderDomain "travel-api/services/order/domain"
	"travel-api/services/user/domain"
	"travel-api/services/user/infrastructure"
)

var (
	ErrInvalidUserRequest = errors.New("invalid_user_request")
	ErrTravelerNotFound   = errors.New("traveler_not_found")
	ErrDocumentDuplicate  = errors.New("document_duplicate")
)

type UserRepository interface {
	EnsureUserDefaults(userID, email string) error
	GetProfile(userID, email string) (domain.UserProfile, error)
	UpsertProfile(userID string, profile domain.UserProfile) (domain.UserProfile, error)
	ListTravelers(userID string) ([]domain.TravelerProfile, error)
	GetTraveler(userID string, travelerID int) (domain.TravelerProfile, bool, error)
	CreateTraveler(userID string, input domain.TravelerProfileInput) (domain.TravelerProfile, error)
	UpdateTraveler(userID string, travelerID int, input domain.TravelerProfileInput, hasDocument bool) (domain.TravelerProfile, error)
	DeleteTraveler(userID string, travelerID int) error
	SetDefaultTraveler(userID string, travelerID int) (domain.TravelerProfile, error)
	GetMembership(userID string) (domain.Membership, error)
	GetRoles(userID string) ([]domain.Role, error)
	BuildOrderSnapshots(userID string, req domain.TravelerSnapshotRequest) ([]domain.OrderTravelerSnapshot, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService() *UserService {
	return &UserService{repo: infrastructure.NewSQLiteUserRepo()}
}

func NewUserServiceWithRepo(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) EnsureUserDefaults(userID, email string) error {
	if strings.TrimSpace(userID) == "" {
		return ErrInvalidUserRequest
	}
	return s.repo.EnsureUserDefaults(userID, email)
}

func (s *UserService) GetProfile(userID, email string) (domain.UserProfile, error) {
	if strings.TrimSpace(userID) == "" {
		return domain.UserProfile{}, ErrInvalidUserRequest
	}
	return s.repo.GetProfile(userID, email)
}

func (s *UserService) UpsertProfile(userID string, profile domain.UserProfile) (domain.UserProfile, error) {
	if strings.TrimSpace(userID) == "" {
		return domain.UserProfile{}, ErrInvalidUserRequest
	}
	profile.DisplayName = strings.TrimSpace(profile.DisplayName)
	profile.Avatar = strings.TrimSpace(profile.Avatar)
	profile.Phone = strings.TrimSpace(profile.Phone)
	profile.Nationality = strings.TrimSpace(profile.Nationality)
	profile.Language = normalizeChoice(profile.Language, "en")
	profile.Currency = normalizeChoice(strings.ToUpper(profile.Currency), "CNY")
	profile.TravelStyle = strings.TrimSpace(profile.TravelStyle)
	profile.BudgetLevel = strings.TrimSpace(profile.BudgetLevel)
	profile.FamilyType = strings.TrimSpace(profile.FamilyType)
	profile.DietaryRestrictions = normalizeList(profile.DietaryRestrictions)
	profile.AccessibilityNeeds = normalizeList(profile.AccessibilityNeeds)
	if profile.DisplayName == "" {
		profile.DisplayName = "Traveler"
	}
	return s.repo.UpsertProfile(userID, profile)
}

func (s *UserService) ListTravelers(userID string) ([]domain.TravelerProfile, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrInvalidUserRequest
	}
	return s.repo.ListTravelers(userID)
}

func (s *UserService) GetTraveler(userID string, travelerID int) (domain.TravelerProfile, error) {
	if strings.TrimSpace(userID) == "" || travelerID <= 0 {
		return domain.TravelerProfile{}, ErrInvalidUserRequest
	}
	traveler, ok, err := s.repo.GetTraveler(userID, travelerID)
	if err != nil {
		return domain.TravelerProfile{}, err
	}
	if !ok {
		return domain.TravelerProfile{}, ErrTravelerNotFound
	}
	return traveler, nil
}

func (s *UserService) CreateTraveler(userID string, input domain.TravelerProfileInput) (domain.TravelerProfile, error) {
	if strings.TrimSpace(userID) == "" {
		return domain.TravelerProfile{}, ErrInvalidUserRequest
	}
	normalized, err := normalizeTraveler(input, true)
	if err != nil {
		return domain.TravelerProfile{}, err
	}
	traveler, err := s.repo.CreateTraveler(userID, normalized)
	return traveler, mapRepoErr(err)
}

func (s *UserService) UpdateTraveler(userID string, travelerID int, input domain.TravelerProfileInput) (domain.TravelerProfile, error) {
	if strings.TrimSpace(userID) == "" || travelerID <= 0 {
		return domain.TravelerProfile{}, ErrInvalidUserRequest
	}
	current, err := s.GetTraveler(userID, travelerID)
	if err != nil {
		return domain.TravelerProfile{}, err
	}
	hasDocument := strings.TrimSpace(input.DocumentNo) != ""
	if !hasDocument {
		input.DocumentType = current.DocumentType
	}
	normalized, err := normalizeTraveler(input, hasDocument)
	if err != nil {
		return domain.TravelerProfile{}, err
	}
	traveler, err := s.repo.UpdateTraveler(userID, travelerID, normalized, hasDocument)
	return traveler, mapRepoErr(err)
}

func (s *UserService) DeleteTraveler(userID string, travelerID int) error {
	if strings.TrimSpace(userID) == "" || travelerID <= 0 {
		return ErrInvalidUserRequest
	}
	return mapRepoErr(s.repo.DeleteTraveler(userID, travelerID))
}

func (s *UserService) SetDefaultTraveler(userID string, travelerID int) (domain.TravelerProfile, error) {
	if strings.TrimSpace(userID) == "" || travelerID <= 0 {
		return domain.TravelerProfile{}, ErrInvalidUserRequest
	}
	traveler, err := s.repo.SetDefaultTraveler(userID, travelerID)
	return traveler, mapRepoErr(err)
}

func (s *UserService) GetMembership(userID string) (domain.Membership, error) {
	if strings.TrimSpace(userID) == "" {
		return domain.Membership{}, ErrInvalidUserRequest
	}
	return s.repo.GetMembership(userID)
}

func (s *UserService) GetRoles(userID string) ([]domain.Role, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrInvalidUserRequest
	}
	return s.repo.GetRoles(userID)
}

func (s *UserService) BuildOrderSnapshots(userID string, req domain.TravelerSnapshotRequest) ([]orderDomain.Traveler, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrInvalidUserRequest
	}
	if len(req.TravelerIDs)+len(req.Travelers) == 0 {
		return []orderDomain.Traveler{}, nil
	}
	for _, id := range req.TravelerIDs {
		if id <= 0 {
			return nil, ErrInvalidUserRequest
		}
	}
	for _, traveler := range req.Travelers {
		if _, err := normalizeTraveler(traveler, true); err != nil {
			return nil, err
		}
	}
	snapshots, err := s.repo.BuildOrderSnapshots(userID, req)
	if err != nil {
		return nil, mapRepoErr(err)
	}
	travelers := make([]orderDomain.Traveler, 0, len(snapshots))
	for _, snapshot := range snapshots {
		travelers = append(travelers, orderDomain.Traveler{
			SourceTravelerID: snapshot.SourceTravelerID,
			Name:             snapshot.Name,
			Gender:           snapshot.Gender,
			BirthDate:        snapshot.BirthDate,
			DocumentType:     snapshot.DocumentType,
			DocumentNoMasked: snapshot.DocumentNoMasked,
			Nationality:      snapshot.Nationality,
			Phone:            snapshot.Phone,
			Email:            snapshot.Email,
		})
	}
	return travelers, nil
}

func normalizeTraveler(input domain.TravelerProfileInput, requireDocument bool) (domain.TravelerProfileInput, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Gender = strings.TrimSpace(input.Gender)
	input.BirthDate = strings.TrimSpace(input.BirthDate)
	input.DocumentType = strings.ToUpper(strings.TrimSpace(input.DocumentType))
	input.DocumentNo = strings.TrimSpace(input.DocumentNo)
	input.Nationality = strings.TrimSpace(input.Nationality)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	if input.Name == "" {
		return domain.TravelerProfileInput{}, ErrInvalidUserRequest
	}
	if requireDocument && (input.DocumentType == "" || input.DocumentNo == "") {
		return domain.TravelerProfileInput{}, ErrInvalidUserRequest
	}
	if !requireDocument && input.DocumentType == "" {
		return domain.TravelerProfileInput{}, ErrInvalidUserRequest
	}
	return input, nil
}

func normalizeList(values []string) []string {
	cleaned := []string{}
	seen := map[string]bool{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		cleaned = append(cleaned, value)
	}
	return cleaned
}

func normalizeChoice(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func mapRepoErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, infrastructure.ErrDuplicateDocument) {
		return ErrDocumentDuplicate
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrTravelerNotFound
	}
	return err
}
