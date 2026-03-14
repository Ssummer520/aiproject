package infrastructure

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

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
	mu    sync.RWMutex
	users []domain.User
	path  string
}

func NewFileUserRepo() *FileUserRepo {
	r := &FileUserRepo{users: []domain.User{}, path: UsersFile}
	r.load()
	return r
}

func (r *FileUserRepo) load() {
	r.mu.Lock()
	defer r.mu.Unlock()
	dir := filepath.Dir(r.path)
	_ = os.MkdirAll(dir, 0755)
	b, err := os.ReadFile(r.path)
	if err != nil {
		return
	}
	_ = json.Unmarshal(b, &r.users)
	if r.users == nil {
		r.users = []domain.User{}
	}
}

func (r *FileUserRepo) save() {
	b, _ := json.MarshalIndent(r.users, "", "  ")
	_ = os.MkdirAll(filepath.Dir(r.path), 0755)
	_ = os.WriteFile(r.path, b, 0644)
}

func (r *FileUserRepo) FindByEmail(email string) (*domain.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.users {
		if r.users[i].Email == email {
			u := r.users[i]
			return &u, true
		}
	}
	return nil, false
}

func (r *FileUserRepo) FindByID(id string) (*domain.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := range r.users {
		if r.users[i].ID == id {
			u := r.users[i]
			return &u, true
		}
	}
	return nil, false
}

func (r *FileUserRepo) Create(u domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users = append(r.users, u)
	r.save()
	return nil
}

func (r *FileUserRepo) Update(u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.users {
		if r.users[i].ID == u.ID {
			r.users[i] = *u
			r.save()
			return nil
		}
	}
	return os.ErrNotExist
}

type FileTokenStore struct {
	mu     sync.RWMutex
	tokens map[string]domain.TokenInfo
	path   string
	ttl    time.Duration
}

func NewFileTokenStore(path string, ttl time.Duration) *FileTokenStore {
	s := &FileTokenStore{tokens: make(map[string]domain.TokenInfo), path: path, ttl: ttl}
	s.load()
	return s
}

func (s *FileTokenStore) load() {
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = os.MkdirAll(filepath.Dir(s.path), 0755)
	b, err := os.ReadFile(s.path)
	if err != nil {
		return
	}
	_ = json.Unmarshal(b, &s.tokens)
	if s.tokens == nil {
		s.tokens = make(map[string]domain.TokenInfo)
	}
}

func (s *FileTokenStore) save() {
	b, _ := json.MarshalIndent(s.tokens, "", "  ")
	_ = os.WriteFile(s.path, b, 0644)
}

func (s *FileTokenStore) Set(token, userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token] = domain.TokenInfo{UserID: userID, ExpiresAt: time.Now().Add(s.ttl)}
	s.save()
}

func (s *FileTokenStore) Get(token string) (userID string, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	info, ok := s.tokens[token]
	if !ok || time.Now().After(info.ExpiresAt) {
		return "", false
	}
	return info.UserID, true
}

func (s *FileTokenStore) Delete(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokens, token)
	s.save()
}
