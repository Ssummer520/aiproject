package infrastructure

import (
	"database/sql"
	"encoding/json"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"travel-api/internal/db"
)

const interactionsFile = "data/interactions.json"

type UserInteraction struct {
	History   []int        `json:"history"`
	Favorites map[int]bool `json:"favorites"`
}

type FileInteractionRepo struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewFileInteractionRepo() *FileInteractionRepo {
	database, err := db.Open()
	if err != nil {
		panic(err)
	}
	r := &FileInteractionRepo{db: database}
	r.migrateFromJSON()
	return r
}

func (r *FileInteractionRepo) migrateFromJSON() {
	r.mu.Lock()
	defer r.mu.Unlock()

	var favCount int
	var historyCount int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM favorites`).Scan(&favCount); err != nil {
		return
	}
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM history`).Scan(&historyCount); err != nil {
		return
	}
	if favCount > 0 || historyCount > 0 {
		return
	}

	b, err := os.ReadFile(interactionsFile)
	if err != nil {
		return
	}
	var data map[string]UserInteraction
	if err := json.Unmarshal(b, &data); err != nil {
		return
	}
	now := time.Now()
	for userID, interaction := range data {
		for id, isFav := range interaction.Favorites {
			_, _ = r.db.Exec(
				`INSERT OR REPLACE INTO favorites(user_id, destination_id, is_favorite, updated_at) VALUES(?, ?, ?, ?)`,
				userID,
				id,
				boolToInt(isFav),
				now.Format(time.RFC3339Nano),
			)
		}
		for index, id := range interaction.History {
			viewedAt := now.Add(-time.Duration(index) * time.Minute).Format(time.RFC3339Nano)
			_, _ = r.db.Exec(
				`INSERT OR REPLACE INTO history(user_id, destination_id, viewed_at) VALUES(?, ?, ?)`,
				userID,
				id,
				viewedAt,
			)
		}
	}
}

func (r *FileInteractionRepo) Get(userID string) UserInteraction {
	r.mu.RLock()
	defer r.mu.RUnlock()

	interaction := UserInteraction{History: []int{}, Favorites: make(map[int]bool)}

	favRows, err := r.db.Query(
		`SELECT destination_id, is_favorite FROM favorites WHERE user_id = ?`,
		userID,
	)
	if err == nil {
		defer favRows.Close()
		for favRows.Next() {
			var id int
			var isFavorite int
			if err := favRows.Scan(&id, &isFavorite); err == nil {
				interaction.Favorites[id] = isFavorite == 1
			}
		}
	}

	historyRows, err := r.db.Query(
		`SELECT destination_id FROM history WHERE user_id = ? ORDER BY viewed_at DESC LIMIT 20`,
		userID,
	)
	if err == nil {
		defer historyRows.Close()
		for historyRows.Next() {
			var id int
			if err := historyRows.Scan(&id); err == nil {
				interaction.History = append(interaction.History, id)
			}
		}
	}

	return interaction
}

func (r *FileInteractionRepo) Set(userID string, u UserInteraction) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if u.Favorites == nil {
		u.Favorites = make(map[int]bool)
	}
	now := time.Now().Format(time.RFC3339Nano)

	favoriteIDs := make([]int, 0, len(u.Favorites))
	for id := range u.Favorites {
		favoriteIDs = append(favoriteIDs, id)
	}
	sort.Ints(favoriteIDs)
	for _, id := range favoriteIDs {
		_, _ = r.db.Exec(
			`INSERT OR REPLACE INTO favorites(user_id, destination_id, is_favorite, updated_at) VALUES(?, ?, ?, ?)`,
			userID,
			id,
			boolToInt(u.Favorites[id]),
			now,
		)
	}

	_, _ = r.db.Exec(`DELETE FROM history WHERE user_id = ?`, userID)
	for index, id := range u.History {
		viewedAt := time.Now().Add(-time.Duration(index) * time.Minute).Format(time.RFC3339Nano)
		_, _ = r.db.Exec(
			`INSERT OR REPLACE INTO history(user_id, destination_id, viewed_at) VALUES(?, ?, ?)`,
			userID,
			id,
			viewedAt,
		)
	}
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func parseJSONFavoriteKey(key string) int {
	id, _ := strconv.Atoi(key)
	return id
}
