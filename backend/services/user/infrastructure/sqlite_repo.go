package infrastructure

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"travel-api/internal/db"
	"travel-api/services/user/domain"
)

var ErrDuplicateDocument = errors.New("duplicate_document")

type SQLiteUserRepo struct {
	mu     sync.RWMutex
	db     *sql.DB
	crypto *PIICrypto
}

func NewSQLiteUserRepo() *SQLiteUserRepo {
	database, err := db.Open()
	if err != nil {
		panic(err)
	}
	crypto, err := NewPIICrypto()
	if err != nil {
		panic(err)
	}
	return &SQLiteUserRepo{db: database, crypto: crypto}
}

func (r *SQLiteUserRepo) EnsureUserDefaults(userID, email string) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil
	}
	now := time.Now().Format(time.RFC3339Nano)
	name := strings.TrimSpace(strings.Split(email, "@")[0])
	if name == "" {
		name = "Traveler"
	}
	_, err := r.db.Exec(`INSERT OR IGNORE INTO user_profiles(user_id, display_name, avatar, phone, nationality, passport_name, language, currency, travel_preferences, dietary_restrictions, membership_level, points_balance, travel_style, budget_level, family_type, accessibility_needs, updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		userID, name, "", "", "", "", "en", "CNY", "culture,food,family", "[]", "Silver", 300, "culture", "mid_range", "", "[]", now)
	if err != nil {
		return err
	}
	benefits, _ := json.Marshal(defaultBenefits("Silver"))
	_, err = r.db.Exec(`INSERT OR IGNORE INTO memberships(user_id, level, points_balance, valid_until, benefits, updated_at) VALUES(?,?,?,?,?,?)`,
		userID, "Silver", 300, time.Now().AddDate(1, 0, 0).Format("2006-01-02"), string(benefits), now)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`INSERT OR IGNORE INTO user_roles(user_id, role_code, created_at) VALUES(?,?,?)`, userID, "customer", now)
	return err
}

func (r *SQLiteUserRepo) GetProfile(userID, email string) (domain.UserProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if err := r.EnsureUserDefaults(userID, email); err != nil {
		return domain.UserProfile{}, err
	}
	var p domain.UserProfile
	var dietaryRaw, accessibilityRaw string
	err := r.db.QueryRow(`SELECT user_id, display_name, avatar, phone, nationality, language, currency, COALESCE(travel_style, ''), COALESCE(budget_level, ''), COALESCE(family_type, ''), COALESCE(dietary_restrictions, ''), COALESCE(accessibility_needs, ''), updated_at FROM user_profiles WHERE user_id=?`, userID).
		Scan(&p.UserID, &p.DisplayName, &p.Avatar, &p.Phone, &p.Nationality, &p.Language, &p.Currency, &p.TravelStyle, &p.BudgetLevel, &p.FamilyType, &dietaryRaw, &accessibilityRaw, &p.UpdatedAt)
	if err != nil {
		return domain.UserProfile{}, err
	}
	p.DietaryRestrictions = decodeList(dietaryRaw)
	p.AccessibilityNeeds = decodeList(accessibilityRaw)
	return p, nil
}

func (r *SQLiteUserRepo) UpsertProfile(userID string, p domain.UserProfile) (domain.UserProfile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().Format(time.RFC3339Nano)
	dietary := encodeList(p.DietaryRestrictions)
	accessibility := encodeList(p.AccessibilityNeeds)
	travelPreferences := strings.Join(cleanList([]string{p.TravelStyle, p.BudgetLevel, p.FamilyType}), ",")
	if travelPreferences == "" {
		travelPreferences = "culture,food,family"
	}
	_, err := r.db.Exec(`INSERT INTO user_profiles(user_id, display_name, avatar, phone, nationality, passport_name, language, currency, travel_preferences, dietary_restrictions, membership_level, points_balance, travel_style, budget_level, family_type, accessibility_needs, updated_at)
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
		ON CONFLICT(user_id) DO UPDATE SET display_name=excluded.display_name, avatar=excluded.avatar, phone=excluded.phone, nationality=excluded.nationality, language=excluded.language, currency=excluded.currency, travel_preferences=excluded.travel_preferences, dietary_restrictions=excluded.dietary_restrictions, travel_style=excluded.travel_style, budget_level=excluded.budget_level, family_type=excluded.family_type, accessibility_needs=excluded.accessibility_needs, updated_at=excluded.updated_at`,
		userID, p.DisplayName, p.Avatar, p.Phone, p.Nationality, "", p.Language, p.Currency, travelPreferences, dietary, "Silver", 300, p.TravelStyle, p.BudgetLevel, p.FamilyType, accessibility, now)
	if err != nil {
		return domain.UserProfile{}, err
	}
	return r.getProfileNoLock(userID)
}

func (r *SQLiteUserRepo) ListTravelers(userID string) ([]domain.TravelerProfile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rows, err := r.db.Query(`SELECT id, user_id, name, gender, birth_date, document_type, document_no_masked, nationality, phone, email, is_default, created_at, updated_at FROM traveler_profiles WHERE user_id=? AND deleted_at IS NULL ORDER BY is_default DESC, updated_at DESC, id DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTravelers(rows)
}

func (r *SQLiteUserRepo) GetTraveler(userID string, travelerID int) (domain.TravelerProfile, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.getTravelerNoLock(userID, travelerID)
}

func (r *SQLiteUserRepo) CreateTraveler(userID string, input domain.TravelerProfileInput) (domain.TravelerProfile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().Format(time.RFC3339Nano)
	encrypted, err := r.crypto.Encrypt(input.DocumentNo)
	if err != nil {
		return domain.TravelerProfile{}, err
	}
	if input.IsDefault {
		if _, err := r.db.Exec(`UPDATE traveler_profiles SET is_default=0, updated_at=? WHERE user_id=? AND deleted_at IS NULL`, now, userID); err != nil {
			return domain.TravelerProfile{}, err
		}
	}
	result, err := r.db.Exec(`INSERT INTO traveler_profiles(user_id, name, gender, birth_date, document_type, document_no_encrypted, document_no_hash, document_no_masked, nationality, phone, email, is_default, created_at, updated_at, deleted_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,NULL)`,
		userID, input.Name, input.Gender, input.BirthDate, input.DocumentType, encrypted, r.crypto.Hash(input.DocumentType, input.DocumentNo), MaskDocument(input.DocumentNo), input.Nationality, input.Phone, input.Email, boolInt(input.IsDefault), now, now)
	if isUniqueErr(err) {
		return domain.TravelerProfile{}, ErrDuplicateDocument
	}
	if err != nil {
		return domain.TravelerProfile{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return domain.TravelerProfile{}, err
	}
	traveler, ok, err := r.getTravelerNoLock(userID, int(id))
	if err != nil || !ok {
		return domain.TravelerProfile{}, err
	}
	return traveler, nil
}

func (r *SQLiteUserRepo) UpdateTraveler(userID string, travelerID int, input domain.TravelerProfileInput, hasDocument bool) (domain.TravelerProfile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().Format(time.RFC3339Nano)
	if input.IsDefault {
		if _, err := r.db.Exec(`UPDATE traveler_profiles SET is_default=0, updated_at=? WHERE user_id=? AND deleted_at IS NULL AND id<>?`, now, userID, travelerID); err != nil {
			return domain.TravelerProfile{}, err
		}
	}
	if hasDocument {
		encrypted, err := r.crypto.Encrypt(input.DocumentNo)
		if err != nil {
			return domain.TravelerProfile{}, err
		}
		_, err = r.db.Exec(`UPDATE traveler_profiles SET name=?, gender=?, birth_date=?, document_type=?, document_no_encrypted=?, document_no_hash=?, document_no_masked=?, nationality=?, phone=?, email=?, is_default=?, updated_at=? WHERE user_id=? AND id=? AND deleted_at IS NULL`,
			input.Name, input.Gender, input.BirthDate, input.DocumentType, encrypted, r.crypto.Hash(input.DocumentType, input.DocumentNo), MaskDocument(input.DocumentNo), input.Nationality, input.Phone, input.Email, boolInt(input.IsDefault), now, userID, travelerID)
		if isUniqueErr(err) {
			return domain.TravelerProfile{}, ErrDuplicateDocument
		}
		if err != nil {
			return domain.TravelerProfile{}, err
		}
	} else {
		_, err := r.db.Exec(`UPDATE traveler_profiles SET name=?, gender=?, birth_date=?, nationality=?, phone=?, email=?, is_default=?, updated_at=? WHERE user_id=? AND id=? AND deleted_at IS NULL`,
			input.Name, input.Gender, input.BirthDate, input.Nationality, input.Phone, input.Email, boolInt(input.IsDefault), now, userID, travelerID)
		if err != nil {
			return domain.TravelerProfile{}, err
		}
	}
	traveler, ok, err := r.getTravelerNoLock(userID, travelerID)
	if err != nil || !ok {
		return domain.TravelerProfile{}, sql.ErrNoRows
	}
	return traveler, nil
}

func (r *SQLiteUserRepo) DeleteTraveler(userID string, travelerID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	result, err := r.db.Exec(`UPDATE traveler_profiles SET deleted_at=?, is_default=0 WHERE user_id=? AND id=? AND deleted_at IS NULL`, time.Now().Format(time.RFC3339Nano), userID, travelerID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *SQLiteUserRepo) SetDefaultTraveler(userID string, travelerID int) (domain.TravelerProfile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok, err := r.getTravelerNoLock(userID, travelerID); err != nil || !ok {
		if err != nil {
			return domain.TravelerProfile{}, err
		}
		return domain.TravelerProfile{}, sql.ErrNoRows
	}
	now := time.Now().Format(time.RFC3339Nano)
	tx, err := r.db.Begin()
	if err != nil {
		return domain.TravelerProfile{}, err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`UPDATE traveler_profiles SET is_default=0, updated_at=? WHERE user_id=? AND deleted_at IS NULL`, now, userID); err != nil {
		return domain.TravelerProfile{}, err
	}
	if _, err := tx.Exec(`UPDATE traveler_profiles SET is_default=1, updated_at=? WHERE user_id=? AND id=? AND deleted_at IS NULL`, now, userID, travelerID); err != nil {
		return domain.TravelerProfile{}, err
	}
	if err := tx.Commit(); err != nil {
		return domain.TravelerProfile{}, err
	}
	traveler, _, err := r.getTravelerNoLock(userID, travelerID)
	return traveler, err
}

func (r *SQLiteUserRepo) GetMembership(userID string) (domain.Membership, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if err := r.EnsureUserDefaults(userID, ""); err != nil {
		return domain.Membership{}, err
	}
	var m domain.Membership
	var benefitsRaw string
	err := r.db.QueryRow(`SELECT user_id, level, points_balance, valid_until, benefits, updated_at FROM memberships WHERE user_id=?`, userID).
		Scan(&m.UserID, &m.Level, &m.PointsBalance, &m.ValidUntil, &benefitsRaw, &m.UpdatedAt)
	if err != nil {
		return domain.Membership{}, err
	}
	m.Benefits = decodeList(benefitsRaw)
	return m, nil
}

func (r *SQLiteUserRepo) GetRoles(userID string) ([]domain.Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if err := r.EnsureUserDefaults(userID, ""); err != nil {
		return nil, err
	}
	rows, err := r.db.Query(`SELECT r.code, r.name FROM user_roles ur JOIN roles r ON r.code=ur.role_code WHERE ur.user_id=? ORDER BY r.code`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	roles := []domain.Role{}
	for rows.Next() {
		var role domain.Role
		if err := rows.Scan(&role.Code, &role.Name); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (r *SQLiteUserRepo) BuildOrderSnapshots(userID string, req domain.TravelerSnapshotRequest) ([]domain.OrderTravelerSnapshot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	snapshots := make([]domain.OrderTravelerSnapshot, 0, len(req.TravelerIDs)+len(req.Travelers))
	for _, id := range req.TravelerIDs {
		traveler, ok, err := r.getTravelerNoLock(userID, id)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, sql.ErrNoRows
		}
		snapshots = append(snapshots, domain.OrderTravelerSnapshot{
			SourceTravelerID: traveler.ID,
			Name:             traveler.Name,
			Gender:           traveler.Gender,
			BirthDate:        traveler.BirthDate,
			DocumentType:     traveler.DocumentType,
			DocumentNoMasked: traveler.DocumentNoMasked,
			Nationality:      traveler.Nationality,
			Phone:            traveler.Phone,
			Email:            traveler.Email,
		})
	}
	for _, traveler := range req.Travelers {
		snapshots = append(snapshots, domain.OrderTravelerSnapshot{
			Name:             traveler.Name,
			Gender:           traveler.Gender,
			BirthDate:        traveler.BirthDate,
			DocumentType:     traveler.DocumentType,
			DocumentNoMasked: MaskDocument(traveler.DocumentNo),
			Nationality:      traveler.Nationality,
			Phone:            traveler.Phone,
			Email:            traveler.Email,
		})
	}
	return snapshots, nil
}

func (r *SQLiteUserRepo) getProfileNoLock(userID string) (domain.UserProfile, error) {
	var p domain.UserProfile
	var dietaryRaw, accessibilityRaw string
	err := r.db.QueryRow(`SELECT user_id, display_name, avatar, phone, nationality, language, currency, COALESCE(travel_style, ''), COALESCE(budget_level, ''), COALESCE(family_type, ''), COALESCE(dietary_restrictions, ''), COALESCE(accessibility_needs, ''), updated_at FROM user_profiles WHERE user_id=?`, userID).
		Scan(&p.UserID, &p.DisplayName, &p.Avatar, &p.Phone, &p.Nationality, &p.Language, &p.Currency, &p.TravelStyle, &p.BudgetLevel, &p.FamilyType, &dietaryRaw, &accessibilityRaw, &p.UpdatedAt)
	if err != nil {
		return domain.UserProfile{}, err
	}
	p.DietaryRestrictions = decodeList(dietaryRaw)
	p.AccessibilityNeeds = decodeList(accessibilityRaw)
	return p, nil
}

func (r *SQLiteUserRepo) getTravelerNoLock(userID string, travelerID int) (domain.TravelerProfile, bool, error) {
	var traveler domain.TravelerProfile
	var isDefault int
	err := r.db.QueryRow(`SELECT id, user_id, name, gender, birth_date, document_type, document_no_masked, nationality, phone, email, is_default, created_at, updated_at FROM traveler_profiles WHERE user_id=? AND id=? AND deleted_at IS NULL`, userID, travelerID).
		Scan(&traveler.ID, &traveler.UserID, &traveler.Name, &traveler.Gender, &traveler.BirthDate, &traveler.DocumentType, &traveler.DocumentNoMasked, &traveler.Nationality, &traveler.Phone, &traveler.Email, &isDefault, &traveler.CreatedAt, &traveler.UpdatedAt)
	if err == sql.ErrNoRows {
		return domain.TravelerProfile{}, false, nil
	}
	if err != nil {
		return domain.TravelerProfile{}, false, err
	}
	traveler.IsDefault = isDefault == 1
	return traveler, true, nil
}

func scanTravelers(rows *sql.Rows) ([]domain.TravelerProfile, error) {
	travelers := []domain.TravelerProfile{}
	for rows.Next() {
		var traveler domain.TravelerProfile
		var isDefault int
		if err := rows.Scan(&traveler.ID, &traveler.UserID, &traveler.Name, &traveler.Gender, &traveler.BirthDate, &traveler.DocumentType, &traveler.DocumentNoMasked, &traveler.Nationality, &traveler.Phone, &traveler.Email, &isDefault, &traveler.CreatedAt, &traveler.UpdatedAt); err != nil {
			return nil, err
		}
		traveler.IsDefault = isDefault == 1
		travelers = append(travelers, traveler)
	}
	return travelers, rows.Err()
}

func encodeList(values []string) string {
	cleaned := cleanList(values)
	if len(cleaned) == 0 {
		return "[]"
	}
	b, _ := json.Marshal(cleaned)
	return string(b)
}

func decodeList(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{}
	}
	var values []string
	if strings.HasPrefix(raw, "[") && json.Unmarshal([]byte(raw), &values) == nil {
		return cleanList(values)
	}
	return cleanList(strings.Split(raw, ","))
}

func cleanList(values []string) []string {
	seen := map[string]bool{}
	cleaned := []string{}
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

func defaultBenefits(level string) []string {
	switch strings.ToLower(level) {
	case "platinum":
		return []string{"priority_support", "refund_priority", "exclusive_deals"}
	case "gold":
		return []string{"priority_support", "member_deals"}
	default:
		return []string{"member_deals", "points_rewards"}
	}
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func isUniqueErr(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "unique")
}
