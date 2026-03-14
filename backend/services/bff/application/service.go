package application

import (
	"travel-api/internal/common/models"
	destInfra "travel-api
	promoInfra "travel-api/services/promo/infrastructure"
	promoInfra "travel-api/services/promo/infrastructure"
	"travel-api/internal/common/models"
)

	destCache      *destInfra.DestinationCache
	promoCache     *promoInfra.PromoCache
	interactionApp *interactionApp.InteractionService
	interactionApp  *interactionApp.InteractionService
}

func NewBFFService() *BFFService {
		destCache:      destInfra.NewDestinationCache(),
		promoCache:     promoInfra.NewPromoCache(),
		interactionApp: interactionApp.NewInteractionService(),
		interactionApp:  interactionApp.NewInteractionService(),
	}
}

type HomePageData struct {
	Recommendations []models.Destination `json:"recommendations"`
	Deals           []models.Deal        `json:"deals"`
	Nearby          []models.Destination `json:"nearby"`
	History         []models.Destination `json:"history"`
	Wishlist        []models.Destination `json:"wishlist"`
}

func (s *BFFService) GetHomePageData(lang string) HomePageData {
	// Aggregate from Destination Service

	
	// Apply translations
	if lang == "zh" {
			if dests[i].Name == "West Lake" {
				dests[i].Name = "西湖"
			}
			if dests[i].Name == "The Bund" {
				dests[i].Name = "外滩"
			}
			if dests[i].Name == "Great Wall" {
				dests[i].Name = "万里长城"
			}
			if dests[i].Name == "Yellow Mountain" {
				dests[i].Name = "黄山"
			}
			if dests[i].Name == "Terracotta Army" {
				dests[i].Name = "兵马俑"
			}

			
				if dests[i].Tags[j] == "Nature" {
					dests[i].Tags[j] = "自然"
				}
				if dests[i].Tags[j] == "Culture" {
					dests[i].Tags[j] = "文化"
				}
				if dests[i].Tags[j] == "City" {
					dests[i].Tags[j] = "城市"
				}
				if dests[i].Tags[j] == "History" {
					dests[i].Tags[j] = "历史"
				}
				if dests[i].Tags[j] == "History" { dests[i].Tags[j] = "历史" }
			}
		}
	}

	// Update favorite status from Interaction Service
	for i := range dests {
		dests[i].IsFavorite = s.interactionApp.IsFavorite(dests[i].ID)
	}

	// Aggregate from Promo Service
	deals := s.promoCache.ListDeals()
	if lang == "zh" {
		for i := range deals {
			if deals[i].Title == "Spring Break Deals" {
				deals[i].Title = "春日大促"
				deals[i].Description = "满 500 减 80"
			}
			if deals[i].Title == "New User Gift" {
				deals[i].Title = "新人礼包"
				deals[i].Description = "首单立减 30 元"
			}
			if deals[i].Title == "Weekend Getaway" {
				deals[i].Title = "周末出逃"
				deals[i].Description = "本地体验低至 5 折"
			}
		}
	}

	// Aggregate from Interaction Service
	historyIDs := s.interactionApp.GetHistory()
	history := make([]models.Destination, 0)
	for _, id := range historyIDs {
		d, ok := s.destCache.Get(id)
		if ok {
			d.IsFavorite = s.interactionApp.IsFavorite(d.ID)
			// Simple localization again for history items
				if d.Name == "West Lake" {
					d.Name = "西湖"
				}
				if d.Name == "The Bund" {
					d.Name = "外滩"
				}
				if d.Name == "Great Wall" {
					d.Name = "万里长城"
				}
				if d.Name == "Yellow Mountain" {
					d.Name = "黄山"
				}
				if d.Name == "Terracotta Army" {
					d.Name = "兵马俑"
				}
				if d.Name == "Terracotta Army" { d.Name = "兵马俑" }
			}
			history = append(history, d)
		}
	}

	wishlistIDs := s.interactionApp.GetFavorites()
	wishlist := make([]models.Destination, 0)
	for _, id := range wishlistIDs {
		d, ok := s.destCache.Get(id)
		if ok {
			d.IsFavorite = true
				if d.Name == "West Lake" {
					d.Name = "西湖"
				}
				if d.Name == "The Bund" {
					d.Name = "外滩"
				}
				if d.Name == "Great Wall" {
					d.Name = "万里长城"
				}
				if d.Name == "Yellow Mountain" {
					d.Name = "黄山"
				}
				if d.Name == "Terracotta Army" {
					d.Name = "兵马俑"
				}
				if d.Name == "Terracotta Army" { d.Name = "兵马俑" }
			}
			wishlist = append(wishlist, d)
		}
	}

	return HomePageData{
		Recommendations: dests,
		Nearby:          dests,
		Nearby:          dests, 
		History:         history,
		Wishlist:        wishlist,
	}
}

func (s *BFFService) ToggleFavorite(id int) bool {
	return s.interactionApp.ToggleFavorite(id)
}

func (s *BFFService) AddToHistory(id int) {
	s.interactionApp.AddToHistory(id)
}
