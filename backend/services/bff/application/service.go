package application

import (
	"strings"

	"travel-api/internal/common/models"
	bffInfra "travel-api/services/bff/infrastructure"
	destInfra "travel-api/services/destination/infrastructure"

	interactionApp "travel-api/services/interaction/application"
	promoInfra "travel-api/services/promo/infrastructure"
)

type BFFService struct {
	destCache         *destInfra.DestinationCache
	promoCache        *promoInfra.PromoCache
	interactionApp    *interactionApp.InteractionService
	statsStore        *bffInfra.StatsStore
	bookingStore      *bffInfra.BookingStore
	notificationStore *bffInfra.NotificationStore
}

func NewBFFService() *BFFService {
	return &BFFService{
		destCache:         destInfra.NewDestinationCache(),
		promoCache:        promoInfra.NewPromoCache(),
		interactionApp:    interactionApp.NewInteractionService(),
		statsStore:        bffInfra.NewStatsStore(),
		bookingStore:      bffInfra.NewBookingStore(),
		notificationStore: bffInfra.NewNotificationStore(),
	}
}

type HomePageData struct {
	Recommendations  []models.Destination `json:"recommendations"`
	Deals            []models.Deal        `json:"deals"`
	Nearby           []models.Destination `json:"nearby"`
	History          []models.Destination `json:"history"`
	Wishlist         []models.Destination `json:"wishlist"`
	TrendingThisWeek []models.Destination `json:"trending_this_week"`
	MostViewedNearby []models.Destination `json:"most_viewed_nearby"`
}

func (s *BFFService) GetHomePageData(lang, userID string) HomePageData {
	dests := s.destCache.ListAll()

	for i := range dests {
		dests[i].IsFavorite = s.interactionApp.IsFavorite(userID, dests[i].ID)
	}

	for i := range dests {
		applyZhDestination(&dests[i], lang)
	}
	if lang == "zh" {
		for i := range dests {
			switch dests[i].Name {
			case "West Lake":
				dests[i].Description = "杭州著名的淡水湖，以其自然风光和文化底蕴闻名。"
				dests[i].Policy = "入住前 48 小时可免费取消。"
			case "The Bund":
				dests[i].Description = "上海标志性的滨江地带，展示了殖民时期建筑和未来感十足的天际线。"
				dests[i].Policy = "不可退款预订。"
			case "Great Wall":
				dests[i].Description = "世界七大奇迹之一，横跨中国北方数千英里。"
				dests[i].Policy = "入住前 24 小时可免费取消。"
			case "Yellow Mountain":
				dests[i].Description = "以奇松、怪石、云海、温泉四绝著称。"
				dests[i].Policy = "入住前 72 小时可申请改期。"
			case "Terracotta Army":
				dests[i].Description = "秦始皇陵的随葬品，被誉为世界第八大奇迹。"
				dests[i].Policy = "不支持取消预订。"
			}
		}
	}

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
			case "Family Fun Pass":
				deals[i].Title = "亲子畅玩套票"
				deals[i].Description = "2 大 1 小组合立减 20%"
			case "Museum Night Special":
				deals[i].Title = "博物馆夜场特惠"
				deals[i].Description = "夜场联票 99 元起"
			case "Foodie Trail Package":
				deals[i].Title = "城市觅食路线包"
				deals[i].Description = "美食 + 地陪组合 129 元起"
			}
		}
	}

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

	trendingIDs := s.statsStore.TopByFavorites(10)
	trendingThisWeek := make([]models.Destination, 0)
	if len(trendingIDs) > 0 {
		for _, id := range trendingIDs {
			d, ok := s.destCache.Get(id)
			if ok {
				d.IsFavorite = s.interactionApp.IsFavorite(userID, d.ID)
				applyZhDestination(&d, lang)
				trendingThisWeek = append(trendingThisWeek, d)
			}
		}
	}
	trendingThisWeek = s.fillWithDefaults(
		trendingThisWeek,
		[]int{1, 3, 2, 16, 11, 10, 14, 5, 4},
		lang,
		userID,
		8,
	)

	mostViewedIDs := s.statsStore.TopByViews(10)
	mostViewedNearby := make([]models.Destination, 0)
	if len(mostViewedIDs) > 0 {
		for _, id := range mostViewedIDs {
			d, ok := s.destCache.Get(id)
			if ok {
				d.IsFavorite = s.interactionApp.IsFavorite(userID, d.ID)
				applyZhDestination(&d, lang)
				mostViewedNearby = append(mostViewedNearby, d)
			}
		}
	}
	mostViewedNearby = s.fillWithDefaults(
		mostViewedNearby,
		[]int{1, 2, 3, 6, 14, 9, 11, 13},
		lang,
		userID,
		8,
	)

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
	case "Lingyin Temple":
		d.Name = "灵隐寺"
	case "Shanghai Disney Resort":
		d.Name = "上海迪士尼度假区"
	case "Wuzhen Water Town":
		d.Name = "乌镇水乡"
	case "Summer Palace":
		d.Name = "颐和园"
	case "Forbidden City":
		d.Name = "故宫"
	case "Chengdu Panda Base":
		d.Name = "成都大熊猫繁育研究基地"
	case "Jinli Ancient Street":
		d.Name = "锦里古街"
	case "Xi'an City Wall":
		d.Name = "西安城墙"
	case "Oriental Pearl Tower":
		d.Name = "东方明珠"
	case "Longjing Tea Plantation":
		d.Name = "龙井茶园"
	case "Universal Beijing Resort":
		d.Name = "北京环球度假区"
	case "Zhujiajiao Water Town":
		d.Name = "朱家角古镇"
	case "Mount Qingcheng":
		d.Name = "青城山"
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
		case "Night View":
			d.Tags[j] = "夜景"
		case "Hiking":
			d.Tags[j] = "徒步"
		case "Archaeology":
			d.Tags[j] = "考古"
		case "Theme Park":
			d.Tags[j] = "主题乐园"
		case "Family":
			d.Tags[j] = "亲子"
		case "Museum":
			d.Tags[j] = "博物馆"
		case "Food":
			d.Tags[j] = "美食"
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

func (s *BFFService) SearchDestinations(lang, userID, query, city, category string, minPrice, maxPrice int) map[string]interface{} {
	dests := s.destCache.ListAll()
	results := make([]models.Destination, 0)

	for i := range dests {
		d := &dests[i]

		if query != "" {
			match := false
			lowerQuery := strings.ToLower(query)
			if strings.Contains(strings.ToLower(d.Name), lowerQuery) {
				match = true
			}
			if strings.Contains(strings.ToLower(d.City), lowerQuery) {
				match = true
			}
			if strings.Contains(strings.ToLower(d.Description), lowerQuery) {
				match = true
			}
			if !match {
				continue
			}
		}

		if city != "" && d.City != city {
			continue
		}

		if category != "" {
			found := false
			for _, tag := range d.Tags {
				if strings.Contains(strings.ToLower(tag), strings.ToLower(category)) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if minPrice > 0 && int(d.Price) < minPrice {
			continue
		}
		if maxPrice > 0 && int(d.Price) > maxPrice {
			continue
		}

		d.IsFavorite = s.interactionApp.IsFavorite(userID, d.ID)
		applyZhDestination(d, lang)
		results = append(results, *d)
	}

	return map[string]interface{}{
		"results": results,
		"total":   len(results),
	}
}

func (s *BFFService) GetCategoryData(lang, userID, category string) map[string]interface{} {
	dests := s.destCache.ListAll()
	results := make([]models.Destination, 0)

	for i := range dests {
		d := &dests[i]

		found := false
		lowerCat := strings.ToLower(category)
		if category == "all" {
			found = true
		}
		for _, tag := range d.Tags {
			if strings.Contains(strings.ToLower(tag), lowerCat) {
				found = true
				break
			}
		}
		if strings.Contains(strings.ToLower(d.Name), lowerCat) {
			found = true
		}

		if !found {
			continue
		}

		d.IsFavorite = s.interactionApp.IsFavorite(userID, d.ID)
		applyZhDestination(d, lang)
		results = append(results, *d)
	}

	return map[string]interface{}{
		"results": results,
		"total":   len(results),
	}
}

func (s *BFFService) GetCityData(lang, userID, city string) map[string]interface{} {
	dests := s.destCache.ListAll()
	results := make([]models.Destination, 0)

	for i := range dests {
		d := &dests[i]

		if !strings.Contains(strings.ToLower(d.City), strings.ToLower(city)) {
			continue
		}

		d.IsFavorite = s.interactionApp.IsFavorite(userID, d.ID)
		applyZhDestination(d, lang)
		results = append(results, *d)
	}

	return map[string]interface{}{
		"results": results,
		"total":   len(results),
	}
}

func (s *BFFService) GetDestinationDetail(lang, userID string, id int) models.Destination {
	d, ok := s.destCache.Get(id)
	if !ok {
		return models.Destination{}
	}

	d.IsFavorite = s.interactionApp.IsFavorite(userID, d.ID)
	applyZhDestination(&d, lang)
	if userID != "" {
		s.AddToHistory(userID, id)
	} else {
		s.statsStore.IncrementView(id)
	}

	return d
}

func (s *BFFService) GetUserBookings(userID string) []models.Booking {
	return s.bookingStore.GetUserBookings(userID)
}

func (s *BFFService) CreateBooking(userID string, destID int, checkIn, checkOut string, guests int) models.Booking {
	d, ok := s.destCache.Get(destID)
	if !ok {
		return models.Booking{}
	}

	booking := s.bookingStore.CreateBooking(userID, d, checkIn, checkOut, guests)

	// Add notification
	s.notificationStore.AddNotification(userID, models.Notification{
		Type:    "booking_confirmed",
		Title:   "Booking Confirmed",
		Message: "Your booking for " + d.Name + " has been confirmed!",
		Link:    "/trips",
	})

	return booking
}

func (s *BFFService) CancelBooking(userID string, bookingID int) (models.Booking, bool) {
	booking, ok := s.bookingStore.CancelBooking(userID, bookingID)
	if !ok {
		return models.Booking{}, false
	}

	s.notificationStore.AddNotification(userID, models.Notification{
		Type:    "booking_cancelled",
		Title:   "Booking Cancelled",
		Message: "Your booking for " + booking.Name + " has been cancelled.",
		Link:    "/trips",
	})

	return booking, true
}

func (s *BFFService) GetNotifications(userID string) []models.Notification {
	return s.notificationStore.GetNotifications(userID)
}

func (s *BFFService) GetUnreadNotificationCount(userID string) int {
	return s.notificationStore.GetUnreadCount(userID)
}

func (s *BFFService) MarkNotificationAsRead(userID string, notificationID int) {
	s.notificationStore.MarkAsRead(userID, notificationID)
}

func (s *BFFService) getDestByID(id int) models.Destination {
	d, _ := s.destCache.Get(id)
	// Always return with basic fields populated
	if d.ID == 0 {
		switch id {
		case 1:
			d = models.Destination{ID: 1, Name: "West Lake", City: "Hangzhou", Rating: 4.9, Lat: 30.25, Lng: 120.15, Tags: []string{"Nature", "Culture"}, Price: 168, ReviewCount: 1200, BookedCount: 14, Cover: "https://images.unsplash.com/photo-1547981609-4b6bfe67ca0b?w=800", Description: "A legendary freshwater lake in Hangzhou, known for its scenic beauty and cultural significance.", Amenities: []string{"Wi-Fi", "Tea Garden", "Guided Tour"}, HostName: "Local Guide Association", Policy: "Free cancellation up to 48h before arrival."}
		case 2:
			d = models.Destination{ID: 2, Name: "The Bund", City: "Shanghai", Rating: 4.8, Lat: 31.24, Lng: 121.49, Tags: []string{"City", "Night View"}, Price: 268, ReviewCount: 3500, BookedCount: 45, Cover: "https://images.unsplash.com/photo-1548115184-bc65ee498ad0?w=800", Description: "The iconic waterfront area in Shanghai, showcasing colonial architecture and futuristic skyline.", Amenities: []string{"Elevator", "Rooftop Bar", "Parking"}, HostName: "Skyline Properties", Policy: "Non-refundable booking."}
		case 3:
			d = models.Destination{ID: 3, Name: "Great Wall", City: "Beijing", Rating: 5.0, Lat: 40.43, Lng: 116.57, Tags: []string{"History", "Hiking"}, Price: 198, ReviewCount: 8900, BookedCount: 78, Cover: "https://images.unsplash.com/photo-1508804185872-d7badad00f7d?w=800", Description: "One of the greatest wonders of the world, stretching thousands of miles across northern China.", Amenities: []string{"Hiking Map", "First Aid Kit", "Bottle of Water"}, HostName: "Heritage Tours", Policy: "Cancel up to 24h before for a full refund."}
		case 4:
			d = models.Destination{ID: 4, Name: "Yellow Mountain", City: "Huangshan", Rating: 4.9, Lat: 30.13, Lng: 118.16, Tags: []string{"Nature", "Hiking"}, Price: 228, ReviewCount: 2100, BookedCount: 23, Cover: "https://images.unsplash.com/photo-1525113190471-9969be29263a?w=800", Description: "Famous for its peculiar pines, odd rocks, sea of clouds, and hot springs.", Amenities: []string{"Mountain Gear", "Sunrise Viewpoint", "Oxygen Tank"}, HostName: "Peak Adventures", Policy: "Reschedule possible up to 72h before."}
		case 5:
			d = models.Destination{ID: 5, Name: "Terracotta Army", City: "Xi'an", Rating: 4.8, Lat: 34.38, Lng: 109.28, Tags: []string{"History", "Archaeology"}, Price: 158, ReviewCount: 5400, BookedCount: 56, Cover: "https://images.unsplash.com/photo-1523482580672-f109ba8cb9be?w=800", Description: "A collection of terracotta sculptures depicting the armies of Qin Shi Huang, the first Emperor of China.", Amenities: []string{"Museum Pass", "Audio Guide", "Souvenir Shop"}, HostName: "Qin Dynasty Experts", Policy: "No cancellations allowed."}
		}
	}
	return d
}

func (s *BFFService) fillWithDefaults(
	list []models.Destination,
	defaultIDs []int,
	lang, userID string,
	target int,
) []models.Destination {
	out := make([]models.Destination, 0, len(list)+len(defaultIDs))
	seen := make(map[int]bool, len(list)+len(defaultIDs))

	for _, d := range list {
		if d.ID == 0 || seen[d.ID] {
			continue
		}
		seen[d.ID] = true
		out = append(out, d)
	}

	for _, id := range defaultIDs {
		if target > 0 && len(out) >= target {
			break
		}
		if seen[id] {
			continue
		}
		d := s.getDestByID(id)
		if d.ID == 0 {
			continue
		}
		d.IsFavorite = s.interactionApp.IsFavorite(userID, d.ID)
		applyZhDestination(&d, lang)
		seen[d.ID] = true
		out = append(out, d)
	}

	if target > 0 && len(out) > target {
		return out[:target]
	}
	return out
}
