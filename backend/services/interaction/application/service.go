package application

import (
	interactionInfra "travel-api/services/interaction/infrastructure"
)

type InteractionService struct {
	cache *interactionInfra.InteractionCache
}

func NewInteractionService() *InteractionService {
	return &InteractionService{
		cache: interactionInfra.NewInteractionCache(),
	}
}

func (s *InteractionService) ToggleFavorite(id int) bool {
	return s.cache.ToggleFavorite(id)
}

func (s *InteractionService) GetFavorites() []int {
	return s.cache.GetFavorites()
}

func (s *InteractionService) AddToHistory(id int) {
	s.cache.AddToHistory(id)
}

func (s *InteractionService) GetHistory() []int {
	return s.cache.GetHistory()
}

func (s *InteractionService) IsFavorite(id int) bool {
	return s.cache.IsFavorite(id)
}
