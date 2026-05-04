package infrastructure

import (
	"database/sql"
	"encoding/json"
	"strings"

	"travel-api/internal/db"
	"travel-api/services/inbound/domain"
)

type SQLiteInboundRepo struct{ db *sql.DB }

func NewSQLiteInboundRepo() *SQLiteInboundRepo {
	database, err := db.Open()
	if err != nil {
		panic(err)
	}
	repo := &SQLiteInboundRepo{db: database}
	repo.seed()
	return repo
}

func (r *SQLiteInboundRepo) seed() {
	toolkit := []domain.ToolkitItem{
		{ID: 1, Key: "esim-wifi", Title: "eSIM / WiFi setup", TitleZh: "eSIM / WiFi 设置", Category: "connectivity", Description: "Stay online after landing with eSIM or pocket WiFi backup.", Steps: []string{"Install eSIM before departure", "Keep passport name consistent", "Save hotel address offline"}, CTA: "Book eSIM", ProductID: 109},
		{ID: 2, Key: "payment", Title: "Mobile payment guide", TitleZh: "移动支付指南", Category: "payment", Description: "Prepare card, cash and mobile payment fallback for China.", Steps: []string{"Bind an international card where supported", "Carry small cash", "Ask merchants for card/cash fallback"}, CTA: "Open guide", ProductID: 0},
		{ID: 3, Key: "reservation", Title: "Attraction reservation", TitleZh: "景点预约规则", Category: "reservation", Description: "Many museums and temples require real-name reservations.", Steps: []string{"Use passport name", "Book popular slots early", "Bring physical passport"}, CTA: "See rules", ProductID: 0},
		{ID: 4, Key: "phrases", Title: "Useful Chinese phrases", TitleZh: "常用中文短语", Category: "language", Description: "One-tap phrases for drivers, hotels and suppliers.", Steps: []string{"请带我去这个地址", "我已经预订了，请帮我确认", "可以使用现金或银行卡吗？"}, CTA: "Ask concierge", ProductID: 0},
	}
	for _, item := range toolkit {
		_, _ = r.db.Exec(`INSERT OR IGNORE INTO inbound_toolkit(id, key, title, title_zh, category, description, steps, cta, product_id) VALUES(?,?,?,?,?,?,?,?,?)`, item.ID, item.Key, item.Title, item.TitleZh, item.Category, item.Description, mustJSON(item.Steps), item.CTA, item.ProductID)
	}
	rails := []domain.RailRoute{{ID: 1, From: "Shanghai Hongqiao", To: "Hangzhou East", Duration: "45-65 min", Frequency: "Every 10-20 min", PriceFrom: 73, Tip: "Arrive 45 minutes early for passport check.", ProductID: 110}, {ID: 2, From: "Beijing South", To: "Xi'an North", Duration: "4.5-6 h", Frequency: "Hourly daytime departures", PriceFrom: 515, Tip: "Choose aisle seats for luggage access.", ProductID: 0}}
	for _, item := range rails {
		_, _ = r.db.Exec(`INSERT OR IGNORE INTO inbound_rails(id, from_station, to_station, duration, frequency, price_from, tip, product_id) VALUES(?,?,?,?,?,?,?,?)`, item.ID, item.From, item.To, item.Duration, item.Frequency, item.PriceFrom, item.Tip, item.ProductID)
	}
	transfers := []domain.TransferOption{{ID: 1, City: "Hangzhou", From: "Hangzhou Xiaoshan Airport", To: "West Lake hotels", Vehicle: "5-seat car", PriceFrom: 188, DriverTip: "Show driver your hotel name in Chinese.", ProductID: 111}, {ID: 2, City: "Beijing", From: "Capital Airport", To: "Downtown Beijing", Vehicle: "Business van", PriceFrom: 268, DriverTip: "Confirm terminal and arrival hall before landing.", ProductID: 112}}
	for _, item := range transfers {
		_, _ = r.db.Exec(`INSERT OR IGNORE INTO inbound_transfers(id, city, from_place, to_place, vehicle, price_from, driver_tip, product_id) VALUES(?,?,?,?,?,?,?,?)`, item.ID, item.City, item.From, item.To, item.Vehicle, item.PriceFrom, item.DriverTip, item.ProductID)
	}
	passes := []domain.CityPass{{ID: 1, City: "Hangzhou", Name: "Hangzhou 2-Day Culture Pass", Duration: "2 days", Includes: []string{"West Lake cruise", "Lingyin Temple route", "Airport transfer coupon", "AI route"}, PriceFrom: 399, ProductID: 113, AISuggested: true}, {ID: 2, City: "Shanghai", Name: "Shanghai Night Pass", Duration: "1 night", Includes: []string{"Huangpu river cruise", "Night transfer", "Food street guide"}, PriceFrom: 299, ProductID: 114, AISuggested: true}, {ID: 3, City: "Beijing", Name: "Beijing Culture Pass", Duration: "2 days", Includes: []string{"Forbidden City guidance", "Temple route", "Chinese phrases"}, PriceFrom: 459, ProductID: 0, AISuggested: false}}
	for _, item := range passes {
		_, _ = r.db.Exec(`INSERT OR IGNORE INTO city_passes(id, city, name, duration, includes, price_from, product_id, ai_suggested) VALUES(?,?,?,?,?,?,?,?)`, item.ID, item.City, item.Name, item.Duration, mustJSON(item.Includes), item.PriceFrom, item.ProductID, boolToInt(item.AISuggested))
	}
	guides := []domain.CityGuide{{City: "Hangzhou", BestSeason: "March-May and September-November", Weather: "Lake areas can be humid; pack light rain gear.", Transport: "Metro + taxi works well; airport transfer is easiest for late arrivals.", Payment: "Keep card/cash fallback near scenic areas.", Connectivity: "Install eSIM before landing; save hotel address offline.", Reservation: "Lingyin Temple and museums can need passport-based booking.", LanguageTips: []string{"请带我去西湖附近的酒店", "我用护照预约了门票"}, SafetyTips: []string{"Avoid unlicensed lake taxis", "Confirm meeting point names in Chinese"}, DietaryTips: []string{"Vegetarian options: say 我吃素", "For allergies: 我对...过敏"}, FamilyTips: []string{"West Lake cruise is stroller-friendly", "Plan indoor tea museum backup on rainy days"}}, {City: "Shanghai", BestSeason: "April-May and October-November", Weather: "Windy riverside nights; bring a light jacket.", Transport: "Metro is efficient; taxis need Chinese destination names.", Payment: "Cards are easier in hotels/malls than small vendors.", Connectivity: "eSIM works well; keep VPN needs checked before travel.", Reservation: "Museums and observation decks may require timed slots.", LanguageTips: []string{"请带我去外滩", "我想去虹桥火车站"}, SafetyTips: []string{"Watch for tea-house scams near tourist streets"}, DietaryTips: []string{"Halal restaurants are available near People's Square"}, FamilyTips: []string{"Avoid rush-hour metro with luggage"}}}
	for _, item := range guides {
		_, _ = r.db.Exec(`INSERT OR IGNORE INTO inbound_city_guides(city, best_season, weather, transport, payment, connectivity, reservation, language_tips, safety_tips, dietary_tips, family_tips) VALUES(?,?,?,?,?,?,?,?,?,?,?)`, item.City, item.BestSeason, item.Weather, item.Transport, item.Payment, item.Connectivity, item.Reservation, mustJSON(item.LanguageTips), mustJSON(item.SafetyTips), mustJSON(item.DietaryTips), mustJSON(item.FamilyTips))
	}
}

func (r *SQLiteInboundRepo) Snapshot() (domain.InboundSnapshot, error) {
	toolkit, err := r.Toolkit()
	if err != nil {
		return domain.InboundSnapshot{}, err
	}
	rails, _ := r.Rails()
	transfers, _ := r.Transfers("")
	passes, _ := r.Passes()
	guides, _ := r.Guides()
	return domain.InboundSnapshot{Toolkit: toolkit, Rails: rails, Transfers: transfers, Passes: passes, Guides: guides}, nil
}
func (r *SQLiteInboundRepo) Toolkit() ([]domain.ToolkitItem, error) {
	rows, err := r.db.Query(`SELECT id, key, title, title_zh, category, description, steps, cta, product_id FROM inbound_toolkit ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.ToolkitItem{}
	for rows.Next() {
		var item domain.ToolkitItem
		var steps string
		if err := rows.Scan(&item.ID, &item.Key, &item.Title, &item.TitleZh, &item.Category, &item.Description, &steps, &item.CTA, &item.ProductID); err != nil {
			return nil, err
		}
		item.Steps = parseSlice(steps)
		out = append(out, item)
	}
	return out, rows.Err()
}
func (r *SQLiteInboundRepo) Rails() ([]domain.RailRoute, error) {
	rows, err := r.db.Query(`SELECT id, from_station, to_station, duration, frequency, price_from, tip, product_id FROM inbound_rails ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.RailRoute{}
	for rows.Next() {
		var item domain.RailRoute
		if err := rows.Scan(&item.ID, &item.From, &item.To, &item.Duration, &item.Frequency, &item.PriceFrom, &item.Tip, &item.ProductID); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}
func (r *SQLiteInboundRepo) Transfers(city string) ([]domain.TransferOption, error) {
	q := `SELECT id, city, from_place, to_place, vehicle, price_from, driver_tip, product_id FROM inbound_transfers`
	args := []interface{}{}
	if strings.TrimSpace(city) != "" {
		q += ` WHERE lower(city)=lower(?)`
		args = append(args, city)
	}
	q += ` ORDER BY id`
	rows, err := r.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.TransferOption{}
	for rows.Next() {
		var item domain.TransferOption
		if err := rows.Scan(&item.ID, &item.City, &item.From, &item.To, &item.Vehicle, &item.PriceFrom, &item.DriverTip, &item.ProductID); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}
func (r *SQLiteInboundRepo) Passes() ([]domain.CityPass, error) {
	rows, err := r.db.Query(`SELECT id, city, name, duration, includes, price_from, product_id, ai_suggested FROM city_passes ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.CityPass{}
	for rows.Next() {
		var item domain.CityPass
		var includes string
		var ai int
		if err := rows.Scan(&item.ID, &item.City, &item.Name, &item.Duration, &includes, &item.PriceFrom, &item.ProductID, &ai); err != nil {
			return nil, err
		}
		item.Includes = parseSlice(includes)
		item.AISuggested = ai == 1
		out = append(out, item)
	}
	return out, rows.Err()
}
func (r *SQLiteInboundRepo) Guides() ([]domain.CityGuide, error) {
	rows, err := r.db.Query(`SELECT city, best_season, weather, transport, payment, connectivity, reservation, language_tips, safety_tips, dietary_tips, family_tips FROM inbound_city_guides ORDER BY city`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.CityGuide{}
	for rows.Next() {
		var item domain.CityGuide
		var language, safety, dietary, family string
		if err := rows.Scan(&item.City, &item.BestSeason, &item.Weather, &item.Transport, &item.Payment, &item.Connectivity, &item.Reservation, &language, &safety, &dietary, &family); err != nil {
			return nil, err
		}
		item.LanguageTips = parseSlice(language)
		item.SafetyTips = parseSlice(safety)
		item.DietaryTips = parseSlice(dietary)
		item.FamilyTips = parseSlice(family)
		out = append(out, item)
	}
	return out, rows.Err()
}
func (r *SQLiteInboundRepo) Guide(city string) (domain.CityGuide, bool, error) {
	guides, err := r.Guides()
	if err != nil {
		return domain.CityGuide{}, false, err
	}
	for _, g := range guides {
		if strings.EqualFold(g.City, city) {
			return g, true, nil
		}
	}
	return domain.CityGuide{}, false, nil
}

func mustJSON(v []string) string { b, _ := json.Marshal(v); return string(b) }
func parseSlice(raw string) []string {
	var out []string
	_ = json.Unmarshal([]byte(raw), &out)
	if out == nil {
		return []string{}
	}
	return out
}
func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
