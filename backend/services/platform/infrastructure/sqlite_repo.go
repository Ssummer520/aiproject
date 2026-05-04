package infrastructure

import (
	"database/sql"
	"strings"
	"time"

	"travel-api/internal/db"
	"travel-api/services/platform/domain"
)

type SQLitePlatformRepo struct{ db *sql.DB }

func NewSQLitePlatformRepo() *SQLitePlatformRepo {
	database, err := db.Open()
	if err != nil {
		panic(err)
	}
	repo := &SQLitePlatformRepo{db: database}
	repo.seed()
	return repo
}

func (r *SQLitePlatformRepo) seed() {
	now := time.Now().Format(time.RFC3339Nano)
	merchants := []domain.Merchant{
		{ID: 1, Name: "Hangzhou Culture Travel Co.", ContactEmail: "ops-hz@chinatravel.demo", Phone: "+86 571 8888 0001", City: "Hangzhou", Status: "active", Rating: 4.8, CreatedAt: now},
		{ID: 2, Name: "Shanghai Night Experience Ltd.", ContactEmail: "ops-sh@chinatravel.demo", Phone: "+86 21 8888 0002", City: "Shanghai", Status: "active", Rating: 4.7, CreatedAt: now},
		{ID: 3, Name: "China Arrival Essentials", ContactEmail: "arrival@chinatravel.demo", Phone: "+86 10 8888 0003", City: "Beijing", Status: "active", Rating: 4.6, CreatedAt: now},
	}
	for _, m := range merchants {
		_, _ = r.db.Exec(`INSERT OR IGNORE INTO merchants(id, name, contact_email, phone, city, status, rating, created_at) VALUES(?,?,?,?,?,?,?,?)`, m.ID, m.Name, m.ContactEmail, m.Phone, m.City, m.Status, m.Rating, m.CreatedAt)
	}
	maps := map[int]int{101: 1, 102: 1, 103: 2, 104: 2, 105: 3, 106: 3, 107: 3, 108: 1}
	for productID, merchantID := range maps {
		_, _ = r.db.Exec(`INSERT OR IGNORE INTO merchant_products(product_id, merchant_id) VALUES(?,?)`, productID, merchantID)
	}
	articles := []domain.CMSArticle{
		{ID: 1, Slug: "china-payment-guide", Title: "China payment guide for overseas travellers", Category: "payment", City: "China", Language: "en", Summary: "Cards, mobile pay, cash and common payment tips.", Content: "Prepare an international card, small cash, and check mobile payment setup before arrival.", Status: "published", UpdatedAt: now},
		{ID: 2, Slug: "hangzhou-attraction-reservation", Title: "杭州景点预约指南", Category: "reservation", City: "Hangzhou", Language: "zh", Summary: "西湖、灵隐寺、博物馆等预约注意事项。", Content: "热门景区建议提前预约，并随身携带护照。", Status: "published", UpdatedAt: now},
	}
	for _, a := range articles {
		_, _ = r.db.Exec(`INSERT OR IGNORE INTO cms_articles(id, slug, title, category, city, language, summary, content, status, updated_at) VALUES(?,?,?,?,?,?,?,?,?,?)`, a.ID, a.Slug, a.Title, a.Category, a.City, a.Language, a.Summary, a.Content, a.Status, a.UpdatedAt)
	}
}

func (r *SQLitePlatformRepo) ListMerchants() ([]domain.Merchant, error) {
	rows, err := r.db.Query(`SELECT id, name, contact_email, phone, city, status, rating, created_at FROM merchants ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []domain.Merchant{}
	for rows.Next() {
		var m domain.Merchant
		if err := rows.Scan(&m.ID, &m.Name, &m.ContactEmail, &m.Phone, &m.City, &m.Status, &m.Rating, &m.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

func (r *SQLitePlatformRepo) ListInventory() ([]domain.InventoryItem, error) {
	rows, err := r.db.Query(`SELECT p.id, p.name, COALESCE(m.id,0), COALESCE(m.name,''), pp.id, pp.name, a.date, a.price, a.stock, a.status FROM product_availability a JOIN products p ON p.id=a.product_id JOIN product_packages pp ON pp.id=a.package_id LEFT JOIN merchant_products mp ON mp.product_id=p.id LEFT JOIN merchants m ON m.id=mp.merchant_id WHERE a.date >= ? ORDER BY a.date ASC, p.id ASC LIMIT 80`, time.Now().Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []domain.InventoryItem{}
	for rows.Next() {
		var item domain.InventoryItem
		if err := rows.Scan(&item.ProductID, &item.ProductName, &item.MerchantID, &item.Merchant, &item.PackageID, &item.PackageName, &item.Date, &item.Price, &item.Stock, &item.Status); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *SQLitePlatformRepo) UpdateInventory(req domain.InventoryUpdateRequest) (domain.InventoryItem, error) {
	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = "available"
	}
	_, err := r.db.Exec(`UPDATE product_availability SET price = CASE WHEN ? > 0 THEN ? ELSE price END, stock = ?, status = ? WHERE package_id = ? AND date = ?`, req.Price, req.Price, req.Stock, status, req.PackageID, req.Date)
	if err != nil {
		return domain.InventoryItem{}, err
	}
	var item domain.InventoryItem
	err = r.db.QueryRow(`SELECT p.id, p.name, COALESCE(m.id,0), COALESCE(m.name,''), pp.id, pp.name, a.date, a.price, a.stock, a.status FROM product_availability a JOIN products p ON p.id=a.product_id JOIN product_packages pp ON pp.id=a.package_id LEFT JOIN merchant_products mp ON mp.product_id=p.id LEFT JOIN merchants m ON m.id=mp.merchant_id WHERE a.package_id=? AND a.date=?`, req.PackageID, req.Date).Scan(&item.ProductID, &item.ProductName, &item.MerchantID, &item.Merchant, &item.PackageID, &item.PackageName, &item.Date, &item.Price, &item.Stock, &item.Status)
	return item, err
}

func (r *SQLitePlatformRepo) ListOrders() ([]domain.PlatformOrder, error) {
	rows, err := r.db.Query(`SELECT o.id, o.user_id, o.status, o.payment_status, o.total_amount, o.currency, oi.product_name, oi.package_name, oi.city, oi.travel_date, o.created_at FROM orders o LEFT JOIN order_items oi ON oi.user_id=o.user_id AND oi.order_id=o.id ORDER BY o.created_at DESC LIMIT 80`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []domain.PlatformOrder{}
	for rows.Next() {
		var item domain.PlatformOrder
		if err := rows.Scan(&item.ID, &item.UserID, &item.Status, &item.PaymentStatus, &item.TotalAmount, &item.Currency, &item.ProductName, &item.PackageName, &item.City, &item.TravelDate, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *SQLitePlatformRepo) CreateRefund(req domain.CreateRefundRequest) (domain.RefundRequest, error) {
	now := time.Now().Format(time.RFC3339Nano)
	var amount float64
	if err := r.db.QueryRow(`SELECT total_amount FROM orders WHERE user_id=? AND id=?`, req.UserID, req.OrderID).Scan(&amount); err != nil {
		return domain.RefundRequest{}, err
	}
	var id int
	_ = r.db.QueryRow(`SELECT COALESCE(MAX(id),0)+1 FROM refund_requests`).Scan(&id)
	_, err := r.db.Exec(`INSERT INTO refund_requests(id, user_id, order_id, reason, refund_amount, status, created_at, updated_at) VALUES(?,?,?,?,?,?,?,?)`, id, req.UserID, req.OrderID, strings.TrimSpace(req.Reason), amount, "requested", now, now)
	if err != nil {
		return domain.RefundRequest{}, err
	}
	_, _ = r.db.Exec(`UPDATE orders SET status='refunding', payment_status='refunding', updated_at=? WHERE user_id=? AND id=?`, now, req.UserID, req.OrderID)
	return domain.RefundRequest{ID: id, UserID: req.UserID, OrderID: req.OrderID, Reason: req.Reason, RefundAmount: amount, Status: "requested", CreatedAt: now, UpdatedAt: now}, nil
}

func (r *SQLitePlatformRepo) ListRefunds() ([]domain.RefundRequest, error) {
	rows, err := r.db.Query(`SELECT id, user_id, order_id, reason, refund_amount, status, created_at, updated_at FROM refund_requests ORDER BY created_at DESC LIMIT 80`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []domain.RefundRequest{}
	for rows.Next() {
		var item domain.RefundRequest
		if err := rows.Scan(&item.ID, &item.UserID, &item.OrderID, &item.Reason, &item.RefundAmount, &item.Status, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *SQLitePlatformRepo) GetProfile(userID, email string) (domain.UserProfile, error) {
	now := time.Now().Format(time.RFC3339Nano)
	_, _ = r.db.Exec(`INSERT OR IGNORE INTO user_profiles(user_id, display_name, avatar, phone, nationality, passport_name, language, currency, travel_preferences, dietary_restrictions, membership_level, points_balance, updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)`, userID, strings.Split(email, "@")[0], "", "", "", "", "en", "CNY", "culture,food,family", "", "Silver", 300, now)
	var p domain.UserProfile
	err := r.db.QueryRow(`SELECT user_id, display_name, avatar, phone, nationality, passport_name, language, currency, travel_preferences, dietary_restrictions, membership_level, points_balance, updated_at FROM user_profiles WHERE user_id=?`, userID).Scan(&p.UserID, &p.DisplayName, &p.Avatar, &p.Phone, &p.Nationality, &p.PassportName, &p.Language, &p.Currency, &p.TravelPreferences, &p.DietaryRestrictions, &p.MembershipLevel, &p.PointsBalance, &p.UpdatedAt)
	return p, err
}

func (r *SQLitePlatformRepo) UpsertProfile(userID string, p domain.UserProfile) (domain.UserProfile, error) {
	now := time.Now().Format(time.RFC3339Nano)
	level := strings.TrimSpace(p.MembershipLevel)
	if level == "" {
		level = "Silver"
	}
	_, err := r.db.Exec(`INSERT INTO user_profiles(user_id, display_name, avatar, phone, nationality, passport_name, language, currency, travel_preferences, dietary_restrictions, membership_level, points_balance, updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?) ON CONFLICT(user_id) DO UPDATE SET display_name=excluded.display_name, phone=excluded.phone, nationality=excluded.nationality, passport_name=excluded.passport_name, language=excluded.language, currency=excluded.currency, travel_preferences=excluded.travel_preferences, dietary_restrictions=excluded.dietary_restrictions, membership_level=excluded.membership_level, points_balance=excluded.points_balance, updated_at=excluded.updated_at`, userID, p.DisplayName, p.Avatar, p.Phone, p.Nationality, p.PassportName, p.Language, p.Currency, p.TravelPreferences, p.DietaryRestrictions, level, p.PointsBalance, now)
	if err != nil {
		return domain.UserProfile{}, err
	}
	return r.GetProfile(userID, "")
}

func (r *SQLitePlatformRepo) ListCMS() ([]domain.CMSArticle, error) {
	rows, err := r.db.Query(`SELECT id, slug, title, category, city, language, summary, content, status, updated_at FROM cms_articles ORDER BY updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []domain.CMSArticle{}
	for rows.Next() {
		var a domain.CMSArticle
		if err := rows.Scan(&a.ID, &a.Slug, &a.Title, &a.Category, &a.City, &a.Language, &a.Summary, &a.Content, &a.Status, &a.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, rows.Err()
}

func (r *SQLitePlatformRepo) CreateCMS(a domain.CMSArticle) (domain.CMSArticle, error) {
	now := time.Now().Format(time.RFC3339Nano)
	if a.Status == "" {
		a.Status = "draft"
	}
	_ = r.db.QueryRow(`SELECT COALESCE(MAX(id),0)+1 FROM cms_articles`).Scan(&a.ID)
	a.UpdatedAt = now
	_, err := r.db.Exec(`INSERT INTO cms_articles(id, slug, title, category, city, language, summary, content, status, updated_at) VALUES(?,?,?,?,?,?,?,?,?,?)`, a.ID, a.Slug, a.Title, a.Category, a.City, a.Language, a.Summary, a.Content, a.Status, a.UpdatedAt)
	return a, err
}

func (r *SQLitePlatformRepo) Metrics() (domain.DashboardMetrics, error) {
	var m domain.DashboardMetrics
	_ = r.db.QueryRow(`SELECT COALESCE(SUM(total_amount),0), COUNT(*) FROM orders WHERE status IN ('paid','completed')`).Scan(&m.GMV, &m.OrderCount)
	_ = r.db.QueryRow(`SELECT COUNT(*) FROM products WHERE status='active'`).Scan(&m.ProductCount)
	_ = r.db.QueryRow(`SELECT COUNT(*) FROM merchants WHERE status='active'`).Scan(&m.MerchantCount)
	_ = r.db.QueryRow(`SELECT COUNT(*) FROM refund_requests`).Scan(&m.PublishedCMSCount)
	_ = r.db.QueryRow(`SELECT COUNT(*) FROM itineraries WHERE title LIKE '%AI%' OR title LIKE '%Plan%'`).Scan(&m.AIItineraryCount)
	_ = r.db.QueryRow(`SELECT COUNT(*) FROM cart_items`).Scan(&m.CartItemCount)
	_ = r.db.QueryRow(`SELECT COUNT(*) FROM reviews`).Scan(&m.ReviewCount)
	_ = r.db.QueryRow(`SELECT COUNT(*) FROM user_profiles`).Scan(&m.MemberCount)
	_ = r.db.QueryRow(`SELECT COUNT(*) FROM cms_articles WHERE status='published'`).Scan(&m.PublishedCMSCount)
	var refunding int
	_ = r.db.QueryRow(`SELECT COUNT(*) FROM orders WHERE status IN ('refunding','refunded')`).Scan(&refunding)
	if m.OrderCount > 0 {
		m.RefundRate = float64(refunding) / float64(m.OrderCount)
	}
	rows, err := r.db.Query(`SELECT oi.city, COUNT(*) FROM order_items oi GROUP BY oi.city ORDER BY COUNT(*) DESC LIMIT 6`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var pair domain.MetricPair
			_ = rows.Scan(&pair.Label, &pair.Value)
			m.CityHeat = append(m.CityHeat, pair)
		}
	}
	if m.CartItemCount > 0 {
		m.OperationalWarnings = append(m.OperationalWarnings, "cart_items_waiting_for_checkout")
	}
	return m, nil
}
