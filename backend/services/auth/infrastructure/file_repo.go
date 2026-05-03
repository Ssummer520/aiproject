package infrastructure

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"travel-api/internal/db"
	"travel-api/services/auth/domain"
)

const (
	UsersFile       = "data/users.json"
	TokensFile      = "data/tokens.json"
	ResetTokensFile = "data/reset_tokens.json"
	TokenTTL        = 7 * 24 * time.Hour
	ResetTokenTTL   = 1 * time.Hour
)

type FileUserRepo struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewFileUserRepo() *FileUserRepo {
	database, err := db.Open()
	if err != nil {
		panic(err)
	}
	r := &FileUserRepo{db: database}
	r.migrateFromJSON()
	return r
}

func (r *FileUserRepo) migrateFromJSON() {
	r.mu.Lock()
	defer r.mu.Unlock()

	var count int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count); err != nil || count > 0 {
		return
	}

	b, err := os.ReadFile(UsersFile)
	if err != nil {
		return
	}
	var users []domain.User
	if err := json.Unmarshal(b, &users); err != nil {
		return
	}
	for _, u := range users {
		_, _ = r.db.Exec(
			`INSERT OR IGNORE INTO users(id, email, password_hash, created_at) VALUES(?, ?, ?, ?)`,
			u.ID,
			strings.ToLower(u.Email),
			u.PasswordHash,
			u.CreatedAt.Format(time.RFC3339Nano),
		)
	}
}

func (r *FileUserRepo) FindByEmail(email string) (*domain.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return scanUser(r.db.QueryRow(
		`SELECT id, email, password_hash, created_at FROM users WHERE email = ?`,
		strings.ToLower(email),
	))
}

func (r *FileUserRepo) FindByID(id string) (*domain.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return scanUser(r.db.QueryRow(
		`SELECT id, email, password_hash, created_at FROM users WHERE id = ?`,
		id,
	))
}

func (r *FileUserRepo) Create(u domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec(
		`INSERT INTO users(id, email, password_hash, created_at) VALUES(?, ?, ?, ?)`,
		u.ID,
		strings.ToLower(u.Email),
		u.PasswordHash,
		u.CreatedAt.Format(time.RFC3339Nano),
	)
	return err
}

func (r *FileUserRepo) Update(u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	result, err := r.db.Exec(
		`UPDATE users SET email = ?, password_hash = ?, created_at = ? WHERE id = ?`,
		strings.ToLower(u.Email),
		u.PasswordHash,
		u.CreatedAt.Format(time.RFC3339Nano),
		u.ID,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return os.ErrNotExist
	}
	return nil
}

func scanUser(row *sql.Row) (*domain.User, bool) {
	var u domain.User
	var createdAt string
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &createdAt); err != nil {
		return nil, false
	}
	if parsed, err := time.Parse(time.RFC3339Nano, createdAt); err == nil {
		u.CreatedAt = parsed
	} else if parsed, err := time.Parse(time.RFC3339, createdAt); err == nil {
		u.CreatedAt = parsed
	}
	return &u, true
}

type FileTokenStore struct {
	mu        sync.RWMutex
	db        *sql.DB
	ttl       time.Duration
	tokenType string
	jsonPath  string
}

func NewFileTokenStore(path string, ttl time.Duration) *FileTokenStore {
	database, err := db.Open()
	if err != nil {
		panic(err)
	}
	s := &FileTokenStore{
		db:        database,
		ttl:       ttl,
		tokenType: tokenTypeFromPath(path),
		jsonPath:  path,
	}
	s.migrateFromJSON()
	return s
}

func tokenTypeFromPath(path string) string {
	if filepath.Base(path) == filepath.Base(ResetTokensFile) {
		return "reset"
	}
	return "auth"
}

func (s *FileTokenStore) migrateFromJSON() {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, err := os.ReadFile(s.jsonPath)
	if err != nil {
		return
	}
	var tokens map[string]domain.TokenInfo
	if err := json.Unmarshal(b, &tokens); err != nil {
		return
	}
	for token, info := range tokens {
		if info.UserID == "" || time.Now().After(info.ExpiresAt) {
			continue
		}
		_, _ = s.db.Exec(
			`INSERT OR IGNORE INTO tokens(token, user_id, type, expires_at) VALUES(?, ?, ?, ?)`,
			token,
			info.UserID,
			s.tokenType,
			info.ExpiresAt.Format(time.RFC3339Nano),
		)
	}
}

func (s *FileTokenStore) Set(token, userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, _ = s.db.Exec(
		`INSERT OR REPLACE INTO tokens(token, user_id, type, expires_at) VALUES(?, ?, ?, ?)`,
		token,
		userID,
		s.tokenType,
		time.Now().Add(s.ttl).Format(time.RFC3339Nano),
	)
}

func (s *FileTokenStore) Get(token string) (userID string, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var expiresAt string
	if err := s.db.QueryRow(
		`SELECT user_id, expires_at FROM tokens WHERE token = ? AND type = ?`,
		token,
		s.tokenType,
	).Scan(&userID, &expiresAt); err != nil {
		return "", false
	}
	parsed, err := time.Parse(time.RFC3339Nano, expiresAt)
	if err != nil {
		parsed, err = time.Parse(time.RFC3339, expiresAt)
	}
	if err != nil || time.Now().After(parsed) {
		_, _ = s.db.Exec(`DELETE FROM tokens WHERE token = ? AND type = ?`, token, s.tokenType)
		return "", false
	}
	return userID, true
}

func (s *FileTokenStore) Delete(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, _ = s.db.Exec(`DELETE FROM tokens WHERE token = ? AND type = ?`, token, s.tokenType)
}
