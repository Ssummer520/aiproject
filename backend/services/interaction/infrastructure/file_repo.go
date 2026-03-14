package infrastructure

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

const interactionsFile = "data/interactions.json"

type UserInteraction struct {
	History   []int        `json:"history"`
	Favorites map[int]bool `json:"favorites"`
}

type FileInteractionRepo struct {
	mu   sync.RWMutex
	data map[string]UserInteraction
	path string
}

func NewFileInteractionRepo() *FileInteractionRepo {
	r := &FileInteractionRepo{data: make(map[string]UserInteraction), path: interactionsFile}
	r.load()
	return r
}

func (r *FileInteractionRepo) load() {
	r.mu.Lock()
	defer r.mu.Unlock()
	_ = os.MkdirAll(filepath.Dir(r.path), 0755)
	b, err := os.ReadFile(r.path)
	if err != nil {
		return
	}
	_ = json.Unmarshal(b, &r.data)
	if r.data == nil {
		r.data = make(map[string]UserInteraction)
	}
}

func (r *FileInteractionRepo) save() {
	b, _ := json.MarshalIndent(r.data, "", "  ")
	_ = os.WriteFile(r.path, b, 0644)
}

func (r *FileInteractionRepo) Get(userID string) UserInteraction {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.data[userID]
	if !ok {
		return UserInteraction{History: nil, Favorites: make(map[int]bool)}
	}
	if u.Favorites == nil {
		u.Favorites = make(map[int]bool)
	}
	return u
}

func (r *FileInteractionRepo) Set(userID string, u UserInteraction) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u.Favorites == nil {
		u.Favorites = make(map[int]bool)
	}
	r.data[userID] = u
	r.save()
}
