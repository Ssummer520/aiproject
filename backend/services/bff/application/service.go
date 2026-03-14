package application

import (
	"travel-api/internal/common/models"
	bffInfra "travel-api/services/bff/infrastructure"
	destInfra "travel-api/services/destination/infrastructure"
	
	interactionApp "travel-api/services/interaction/application"
	promoInfra "travel-api/services/promo/infrastructure"
)

type BFFService struct {
	destCache      *destInfra.DestinationCache
	promoCache     *promoInfra.PromoCache
	interactionApp *interactionApp.InteractionService
	statsStore     *bffInfra.StatsStore
}

func NewBFFService() *BFFService {
	return &BFFService{
		destCache:      destInfra.NewDestinationCache(),
		promoCache:     promoInfra.NewPromoCache(),
		interactionApp: interactionApp.NewInteractionService(),
		statsStore:     bffInfra.NewStatsStore(),
	}
}

	Recommendations  []models.Destination `json:"recommendations"`
	Deals            []models.Deal        `json:"deals"`
	Nearby           []models.Destination `json:"nearby"`
	History          []models.Destination `json:"history"`
	Wishlist         []models.Destination `json:"wishlist"`
	TrendingThisWeek []models.Destination `json:"trending_this_week"`
	MostViewedNearby []models.Destination `json:"most_viewed_nearby"`
	MostViewedNearby  []models.Destination `json:"most_viewed_nearby"`
}

func (s *BFFService) GetHomePageData(lang, userID string) HomePageData {
	// Aggregate from Destination Service
	dests := s.destCache.ListAll()

	// Update favorite status from Interaction Service (per user)
	for i := range dests {
		dests[i].IsFavorite = s.interactionApp.IsFavorite(userID, dests[i].ID)
	}

	for i := range dests {
		applyZhDestination(&dests[i], lang)
	}
	if lang == "zh" {
		for i := range dests {
			switch dests[i].Name {
			case "西湖":
				dests[i].Description = "杭州著名的淡水湖，以其自然风光和文化底蕴闻名。"
				dests[i].Policy = "入住前 48 小时可免费取消。"
			case "外滩":
				dests[i].Description = "上海标志性的滨江地带，展示了殖民时期建筑和未来感十足的天际线。"
				dests[i].Policy = "不可退款预订。"
			case "万里长城":
				dests[i].Description = "世界七大奇迹之一，横跨中国北方数千英里。"
				dests[i].Policy = "入住前 24 小时可免费取消。"
			case "黄山":
				dests[i].Description = "以奇松、怪石、云海、温泉“四绝”著称。"
				dests[i].Policy = "入住前 72 小时可申请改期。"
			case "兵马俑":
				dests[i].Description = "秦始皇陵的随葬品，被誉为“世界第八大奇迹”。"
				dests[i].Policy = "不支持取消预订。"
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

	// Aggregate from Interaction Service (per user; empty if not logged in)
	historyIDs := s.interactionApp.GetHistory(userID)
	history := make([]models.Destination, 0)
	for _, id := range historyIDs {
		d, ok := s.destCache.Get(id)
		if ok {
			d.IsFavorite = s.interactionApp.IsFavorite(userID, d.ID)
			applyZhDestination(&d, lang)
			history = append(history, d)
		}
	}

	wishlistIDs := s.interactionApp.GetFavorites(userID)
	wishlist := make([]models.Destination, 0)
	for _, id := range wishlistIDs {
		d, ok := s.destCache.Get(id)
		if ok {
			d.IsFavorite = true
			applyZhDestination(&d, lang)
			wishlist = append(wishlist, d)
		}
	}

	// 排行榜：最近一周喜欢最多 = 收藏数排序；周边点击榜 = 浏览量排序
	trendingIDs := s.statsStore.TopByFavorites(10)
	trendingThisWeek := make([]models.Destination, 0)
	for _, id := range trendingIDs {
		d, ok := s.destCache.Get(id)
		if ok {
			d.IsFavorite = s.interactionApp.IsFavorite(userID, d.ID)
			applyZhDestination(&d, lang)
			trendingThisWeek = append(trendingThisWeek, d)
		}
	}
	mostViewedIDs := s.statsStore.TopByViews(10)
	mostViewedNearby := make([]models.Destination, 0)
	for _, id := range mostViewedIDs {
		d, ok := s.destCache.Get(id)
		if ok {
			d.IsFavorite = s.interactionApp.IsFavorite(userID, d.ID)
			applyZhDestination(&d, lang)
			mostViewedNearby = append(mostViewedNearby, d)
		}
	}

	return HomePageData{
		Recommendations:  dests,
		Deals:            deals,
		Nearby:           dests,
		History:          history,
		Wishlist:         wishlist,
		TrendingThisWeek: trendingThisWeek,
		MostViewedNearby: mostViewedNearby,
	}
}

func applyZhDestination(d *models.Destination, lang string) {
	if lang != "zh" {
		return
	}
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
	for j := range d.Tags {
		switch d.Tags[j] {
		case "Nature":
			d.Tags[j] = "自然"
		case "Culture":
			d.Tags[j] = "文化"
		case "City":
			d.Tags[j] = "城市"
		case "History":
			d.Tags[j] = "历史"
		}
	}
}

func (s *BFFService) ToggleFavorite(userID string, id int) bool {
	result := s.interactionApp.ToggleFavorite(userID, id)
	if result {
		s.statsStore.IncrementFavorite(id, 1)
	} else {
		s.statsStore.IncrementFavorite(id, -1)
	}
	return result
}

func (s *BFFService) AddToHistory(userID string, id int) {
	s.interactionApp.AddToHistory(userID, id)
	s.statsStore.IncrementView(id)
}
