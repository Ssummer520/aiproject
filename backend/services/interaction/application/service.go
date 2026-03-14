package application

import (
	interactionInfra "travel-api/services/interaction/infrastructure"
)

type InteractionService struct {
	repo *interactionInfra.FileInteractionRepo
}

func NewInteractionService() *InteractionService {
	return &InteractionService{
		repo: interactionInfra.NewFileInteractionRepo(),
	}
}

func (s *InteractionService) ToggleFavorite(userID string, id int) bool {
	if userID == "" {
		return false
	}
	u := s.repo.Get(userID)
	u.Favorites[id] = !u.Favorites[id]
	s.repo.Set(userID, u)
	return u.Favorites[id]
}

func (s *InteractionService) GetFavorites(userID string) []int {
	if userID == "" {
		return nil
	}
	u := s.repo.Get(userID)
	favs := make([]int, 0)
	for id, isFav := range u.Favorites {
		if isFav {
			favs = append(favs, id)
		}
	}
	return favs
}

func (s *InteractionService) AddToHistory(userID string, id int) {
	if userID == "" {
		return
	}
	u := s.repo.Get(userID)
	for i, existingID := range u.History {
		if existingID == id {
			u.History = append(u.History[:i], u.History[i+1:]...)
			break
		}
	}
	u.History = append([]int{id}, u.History...)
	if len(u.History) > 20 {
		u.History = u.History[:20]
	}
	s.repo.Set(userID, u)
}

func (s *InteractionService) GetHistory(userID string) []int {
	if userID == "" {
		return nil
	}
	return s.repo.Get(userID).History
}

func (s *InteractionService) IsFavorite(userID string, id int) bool {
	if userID == "" {
		return false
	}
	return s.repo.Get(userID).Favorites[id]
}
