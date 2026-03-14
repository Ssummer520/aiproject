package application

import (
	"travel-api/internal/common/models"
	destInfra "travel-api/services/destination/infrastructure"
	interactionApp "travel-api/services/interaction/application"
	promoInfra "travel-api/services/promo/infrastructure"
)

type BFFService struct {
	destCache      *destInfra.DestinationCache
	promoCache     *promoInfra.PromoCache
	interactionApp *interactionApp.InteractionService
}

func NewBFFService() *BFFService {
	return &BFFService{
		destCache:      destInfra.NewDestinationCache(),
		promoCache:     promoInfra.NewPromoCache(),
		interactionApp: interactionApp.NewInteractionService(),
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
	dests := s.destCache.ListAll()

	// Update favorite status from Interaction Service
	for i := range dests {
		dests[i].IsFavorite = s.interactionApp.IsFavorite(dests[i].ID)
	}

	// Apply translations
	if lang == "zh" {
		for i := range dests {
			switch dests[i].Name {
			case "West Lake":
				dests[i].Name = "西湖"
				dests[i].Description = "杭州著名的淡水湖，以其自然风光和文化底蕴闻名。"
				dests[i].Policy = "入住前 48 小时可免费取消。"
			case "The Bund":
				dests[i].Name = "外滩"
				dests[i].Description = "上海标志性的滨江地带，展示了殖民时期建筑和未来感十足的天际线。"
				dests[i].Policy = "不可退款预订。"
			case "Great Wall":
				dests[i].Name = "万里长城"
				dests[i].Description = "世界七大奇迹之一，横跨中国北方数千英里。"
				dests[i].Policy = "入住前 24 小时可免费取消。"
			case "Yellow Mountain":
				dests[i].Name = "黄山"
				dests[i].Description = "以奇松、怪石、云海、温泉“四绝”著称。"
				dests[i].Policy = "入住前 72 小时可申请改期。"
			case "Terracotta Army":
				dests[i].Name = "兵马俑"
				dests[i].Description = "秦始皇陵的随葬品，被誉为“世界第八大奇迹”。"
				dests[i].Policy = "不支持取消预订。"
			}
			for j := range dests[i].Tags {
				switch dests[i].Tags[j] {
				case "Nature":
					dests[i].Tags[j] = "自然"
				case "Culture":
					dests[i].Tags[j] = "文化"
				case "City":
					dests[i].Tags[j] = "城市"
				case "History":
					dests[i].Tags[j] = "历史"
				}
			}
		}
	}

	// Aggregate from Promo Service
	deals := s.promoCache.ListDeals()
	if lang == "zh" {
		for i := range deals {
			switch deals[i].Title {
			case "Spring Break Deals":
				deals[i].Title = "春日大促"
				deals[i].Description = "满 500 减 80"
			case "New User Gift":
				deals[i].Title = "新人礼包"
				deals[i].Description = "首单立减 30 元"
			case "Weekend Getaway":
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
			if lang == "zh" {
				switch d.Name {
				case "West Lake":
					d.Name = "西湖"
				case "The Bund":
					d.Name = "外滩"
				case "Great Wall":
					d.Name = "万里长城"
				case "Yellow Mountain":
					d.Name = "黄山"
				case "Terracotta Army":
					d.Name = "兵马俑"
				}
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
			if lang == "zh" {
				switch d.Name {
				case "West Lake":
					d.Name = "西湖"
				case "The Bund":
					d.Name = "外滩"
				case "Great Wall":
					d.Name = "万里长城"
				case "Yellow Mountain":
					d.Name = "黄山"
				case "Terracotta Army":
					d.Name = "兵马俑"
				}
			}
			wishlist = append(wishlist, d)
		}
	}

	return HomePageData{
		Recommendations: dests,
		Deals:           deals,
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
