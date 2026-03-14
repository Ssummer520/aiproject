package application

import (
	destInfra "travel-api/services/destination/infrastructure"
	promoInfra "travel-api/services/promo/infrastructure"
	"travel-api/internal/common/models"
)

type BFFService struct {
	destCache  *destInfra.DestinationCache
	promoCache *promoInfra.PromoCache
}

func NewBFFService() *BFFService {
	return &BFFService{
		destCache:  destInfra.NewDestinationCache(),
		promoCache: promoInfra.NewPromoCache(),
	}
}

type HomePageData struct {
	Recommendations []models.Destination `json:"recommendations"`
	Deals           []models.Deal        `json:"deals"`
	Nearby          []models.Destination `json:"nearby"`
}

func (s *BFFService) GetHomePageData(lang string) HomePageData {
	// Aggregate from Destination Service
	dests := s.destCache.ListAll()
	
	// Simple translation logic for demo
	if lang == "zh" {
		for i := range dests {
			if dests[i].Name == "West Lake" { dests[i].Name = "西湖" }
			if dests[i].Name == "The Bund" { dests[i].Name = "外滩" }
			if dests[i].Name == "Great Wall" { dests[i].Name = "万里长城" }
			if dests[i].Name == "Yellow Mountain" { dests[i].Name = "黄山" }
			if dests[i].Name == "Terracotta Army" { dests[i].Name = "兵马俑" }
			
			for j := range dests[i].Tags {
				if dests[i].Tags[j] == "Nature" { dests[i].Tags[j] = "自然" }
				if dests[i].Tags[j] == "Culture" { dests[i].Tags[j] = "文化" }
				if dests[i].Tags[j] == "City" { dests[i].Tags[j] = "城市" }
				if dests[i].Tags[j] == "History" { dests[i].Tags[j] = "历史" }
			}
		}
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

	return HomePageData{
		Recommendations: dests,
		Deals:           deals,
		Nearby:          dests, 
	}
}
