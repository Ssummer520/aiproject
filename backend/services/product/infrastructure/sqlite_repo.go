package infrastructure

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"travel-api/internal/db"
	"travel-api/services/product/domain"
)

type SQLiteProductRepo struct {
	mu sync.RWMutex
	db *sql.DB
}

func NewSQLiteProductRepo() *SQLiteProductRepo {
	database, err := db.Open()
	if err != nil {
		panic(err)
	}
	repo := &SQLiteProductRepo{db: database}
	repo.seedDemoData()
	return repo
}

func (r *SQLiteProductRepo) Search(filters domain.SearchFilters) ([]domain.Product, error) {
	r.mu.RLock()
	rows, err := r.db.Query(`SELECT id, destination_id, city, category, type, name, subtitle, description, cover, images, tags, rating, review_count, booked_count, base_price, currency, instant_confirm, free_cancel, duration, meeting_point, included, excluded, usage, policy, status FROM products WHERE status = 'active'`)
	if err != nil {
		r.mu.RUnlock()
		return nil, err
	}

	candidates := make([]domain.Product, 0)
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			_ = rows.Close()
			r.mu.RUnlock()
			return nil, err
		}
		candidates = append(candidates, product)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		r.mu.RUnlock()
		return nil, err
	}
	if err := rows.Close(); err != nil {
		r.mu.RUnlock()
		return nil, err
	}
	r.mu.RUnlock()

	products := make([]domain.Product, 0, len(candidates))
	for _, product := range candidates {
		packages, err := r.ListPackages(product.ID)
		if err != nil {
			return nil, err
		}
		availability, err := r.ListAvailability(product.ID, filters.Date)
		if err != nil {
			return nil, err
		}
		product.Packages = packages
		product.Availability = availability
		if !matchesFilters(product, filters) {
			continue
		}
		products = append(products, product)
	}

	sortProducts(products, filters.Sort)
	return products, nil
}

func (r *SQLiteProductRepo) Get(id int) (domain.Product, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, err := scanProduct(r.db.QueryRow(`SELECT id, destination_id, city, category, type, name, subtitle, description, cover, images, tags, rating, review_count, booked_count, base_price, currency, instant_confirm, free_cancel, duration, meeting_point, included, excluded, usage, policy, status FROM products WHERE id = ? AND status = 'active'`, id))
	if err == sql.ErrNoRows {
		return domain.Product{}, false, nil
	}
	if err != nil {
		return domain.Product{}, false, err
	}

	packages, err := r.ListPackages(id)
	if err != nil {
		return domain.Product{}, false, err
	}
	availability, err := r.ListAvailability(id, "")
	if err != nil {
		return domain.Product{}, false, err
	}
	product.Packages = packages
	product.Availability = availability
	return product, true, nil
}

func (r *SQLiteProductRepo) GetByDestinationID(destinationID int) (domain.Product, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, err := scanProduct(r.db.QueryRow(`SELECT id, destination_id, city, category, type, name, subtitle, description, cover, images, tags, rating, review_count, booked_count, base_price, currency, instant_confirm, free_cancel, duration, meeting_point, included, excluded, usage, policy, status FROM products WHERE destination_id = ? AND status = 'active' ORDER BY CASE WHEN id < 109 THEN 0 ELSE 1 END ASC, CASE WHEN type = 'bundle' THEN 1 ELSE 0 END ASC, booked_count DESC, id ASC LIMIT 1`, destinationID))
	if err == sql.ErrNoRows {
		return domain.Product{}, false, nil
	}
	if err != nil {
		return domain.Product{}, false, err
	}

	packages, err := r.ListPackages(product.ID)
	if err != nil {
		return domain.Product{}, false, err
	}
	availability, err := r.ListAvailability(product.ID, "")
	if err != nil {
		return domain.Product{}, false, err
	}
	product.Packages = packages
	product.Availability = availability
	return product, true, nil
}

func (r *SQLiteProductRepo) ListPackages(productID int) ([]domain.Package, error) {
	rows, err := r.db.Query(`SELECT id, product_id, name, description, price, original_price, unit, min_quantity, max_quantity, age_rule, refund_policy, confirm_type, voucher_type FROM product_packages WHERE product_id = ? ORDER BY price ASC`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	packages := make([]domain.Package, 0)
	for rows.Next() {
		var pkg domain.Package
		if err := rows.Scan(&pkg.ID, &pkg.ProductID, &pkg.Name, &pkg.Description, &pkg.Price, &pkg.OriginalPrice, &pkg.Unit, &pkg.MinQuantity, &pkg.MaxQuantity, &pkg.AgeRule, &pkg.RefundPolicy, &pkg.ConfirmType, &pkg.VoucherType); err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
	}
	return packages, rows.Err()
}

func (r *SQLiteProductRepo) GetPackage(packageID int) (domain.Package, bool, error) {
	var pkg domain.Package
	err := r.db.QueryRow(`SELECT id, product_id, name, description, price, original_price, unit, min_quantity, max_quantity, age_rule, refund_policy, confirm_type, voucher_type FROM product_packages WHERE id = ?`, packageID).Scan(&pkg.ID, &pkg.ProductID, &pkg.Name, &pkg.Description, &pkg.Price, &pkg.OriginalPrice, &pkg.Unit, &pkg.MinQuantity, &pkg.MaxQuantity, &pkg.AgeRule, &pkg.RefundPolicy, &pkg.ConfirmType, &pkg.VoucherType)
	if err == sql.ErrNoRows {
		return domain.Package{}, false, nil
	}
	if err != nil {
		return domain.Package{}, false, err
	}
	return pkg, true, nil
}

func (r *SQLiteProductRepo) ListAvailability(productID int, date string) ([]domain.Availability, error) {
	query := `SELECT id, product_id, package_id, date, price, stock, status FROM product_availability WHERE product_id = ?`
	args := []interface{}{productID}
	if date != "" {
		query += ` AND date = ?`
		args = append(args, date)
	} else {
		query += ` AND date >= ?`
		args = append(args, time.Now().Format("2006-01-02"))
	}
	query += ` ORDER BY date ASC, package_id ASC LIMIT 240`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	availability := make([]domain.Availability, 0)
	for rows.Next() {
		var item domain.Availability
		if err := rows.Scan(&item.ID, &item.ProductID, &item.PackageID, &item.Date, &item.Price, &item.Stock, &item.Status); err != nil {
			return nil, err
		}
		availability = append(availability, item)
	}
	return availability, rows.Err()
}

func (r *SQLiteProductRepo) GetAvailability(packageID int, date string) (domain.Availability, bool, error) {
	var item domain.Availability
	err := r.db.QueryRow(`SELECT id, product_id, package_id, date, price, stock, status FROM product_availability WHERE package_id = ? AND date = ?`, packageID, date).Scan(&item.ID, &item.ProductID, &item.PackageID, &item.Date, &item.Price, &item.Stock, &item.Status)
	if err == sql.ErrNoRows {
		return domain.Availability{}, false, nil
	}
	if err != nil {
		return domain.Availability{}, false, err
	}
	return item, true, nil
}

func (r *SQLiteProductRepo) IncrementBookedCount(productID int, quantity int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec(`UPDATE products SET booked_count = booked_count + ?, updated_at = ? WHERE id = ?`, quantity, time.Now().Format(time.RFC3339Nano), productID)
	return err
}

func (r *SQLiteProductRepo) seedDemoData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	var count int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM products`).Scan(&count); err != nil {
		return
	}
	if count == 0 {
		for _, product := range demoProducts() {
			r.insertProduct(product)
		}
		for _, pkg := range demoPackages() {
			r.insertPackage(pkg)
		}
	}
	for _, product := range inboundDemoProducts() {
		r.insertProduct(product)
	}
	for _, pkg := range inboundDemoPackages() {
		r.insertPackage(pkg)
	}
	r.seedAvailability()
}

func (r *SQLiteProductRepo) insertProduct(product domain.Product) {
	now := time.Now().Format(time.RFC3339Nano)
	_, _ = r.db.Exec(`INSERT OR IGNORE INTO products(id, destination_id, city, category, type, name, subtitle, description, cover, images, tags, rating, review_count, booked_count, base_price, currency, instant_confirm, free_cancel, duration, meeting_point, included, excluded, usage, policy, status, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		product.ID,
		product.DestinationID,
		product.City,
		product.Category,
		product.Type,
		product.Name,
		product.Subtitle,
		product.Description,
		product.Cover,
		mustJSON(product.Images),
		mustJSON(product.Tags),
		product.Rating,
		product.ReviewCount,
		product.BookedCount,
		product.BasePrice,
		product.Currency,
		boolToInt(product.InstantConfirm),
		boolToInt(product.FreeCancel),
		product.Duration,
		product.MeetingPoint,
		mustJSON(product.Included),
		mustJSON(product.Excluded),
		product.Usage,
		product.Policy,
		product.Status,
		now,
		now,
	)
}

func (r *SQLiteProductRepo) insertPackage(pkg domain.Package) {
	_, _ = r.db.Exec(`INSERT OR IGNORE INTO product_packages(id, product_id, name, description, price, original_price, unit, min_quantity, max_quantity, age_rule, refund_policy, confirm_type, voucher_type) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		pkg.ID,
		pkg.ProductID,
		pkg.Name,
		pkg.Description,
		pkg.Price,
		pkg.OriginalPrice,
		pkg.Unit,
		pkg.MinQuantity,
		pkg.MaxQuantity,
		pkg.AgeRule,
		pkg.RefundPolicy,
		pkg.ConfirmType,
		pkg.VoucherType,
	)
}

func (r *SQLiteProductRepo) seedAvailability() {
	packages, err := r.db.Query(`SELECT id, product_id, price FROM product_packages`)
	if err != nil {
		return
	}
	type seedPackage struct {
		ID        int
		ProductID int
		Price     float64
	}
	items := make([]seedPackage, 0)
	for packages.Next() {
		var item seedPackage
		if err := packages.Scan(&item.ID, &item.ProductID, &item.Price); err == nil {
			items = append(items, item)
		}
	}
	_ = packages.Close()

	today := time.Now()
	for _, item := range items {
		for offset := 0; offset < 45; offset++ {
			date := today.AddDate(0, 0, offset).Format("2006-01-02")
			weekendBoost := 0.0
			weekday := today.AddDate(0, 0, offset).Weekday()
			if weekday == time.Saturday || weekday == time.Sunday {
				weekendBoost = 18
			}
			stock := 18 + (item.ID+offset)%13
			status := "available"
			if offset%17 == 0 && offset > 0 {
				stock = 0
				status = "sold_out"
			}
			_, _ = r.db.Exec(`INSERT OR IGNORE INTO product_availability(product_id, package_id, date, price, stock, status) VALUES(?, ?, ?, ?, ?, ?)`, item.ProductID, item.ID, date, item.Price+weekendBoost, stock, status)
		}
	}
}

func scanProduct(scanner interface {
	Scan(dest ...interface{}) error
}) (domain.Product, error) {
	var product domain.Product
	var images string
	var tags string
	var included string
	var excluded string
	var instantConfirm int
	var freeCancel int
	err := scanner.Scan(&product.ID, &product.DestinationID, &product.City, &product.Category, &product.Type, &product.Name, &product.Subtitle, &product.Description, &product.Cover, &images, &tags, &product.Rating, &product.ReviewCount, &product.BookedCount, &product.BasePrice, &product.Currency, &instantConfirm, &freeCancel, &product.Duration, &product.MeetingPoint, &included, &excluded, &product.Usage, &product.Policy, &product.Status)
	if err != nil {
		return domain.Product{}, err
	}
	product.Images = parseStringSlice(images)
	product.Tags = parseStringSlice(tags)
	product.Included = parseStringSlice(included)
	product.Excluded = parseStringSlice(excluded)
	product.InstantConfirm = instantConfirm == 1
	product.FreeCancel = freeCancel == 1
	return product, nil
}

func matchesFilters(product domain.Product, filters domain.SearchFilters) bool {
	query := strings.ToLower(strings.TrimSpace(filters.Query))
	if query != "" {
		haystack := strings.ToLower(strings.Join([]string{product.Name, product.Subtitle, product.Description, product.City, product.Category, product.Type, product.Duration, product.MeetingPoint, product.Usage, strings.Join(product.Tags, " "), strings.Join(product.Included, " ")}, " "))
		if !strings.Contains(haystack, query) {
			return false
		}
	}
	if filters.City != "" && !strings.EqualFold(product.City, filters.City) {
		return false
	}
	if filters.Category != "" && !strings.EqualFold(product.Category, filters.Category) && !containsFold(product.Tags, filters.Category) {
		return false
	}
	if filters.Type != "" && !strings.EqualFold(product.Type, filters.Type) {
		return false
	}
	if filters.Duration != "" && !strings.Contains(strings.ToLower(product.Duration), strings.ToLower(filters.Duration)) {
		return false
	}
	if filters.Language != "" && !containsAnyFold([]string{product.Description, product.Subtitle, strings.Join(product.Tags, " "), strings.Join(product.Included, " ")}, filters.Language) {
		return false
	}
	for _, feature := range filters.Features {
		if feature != "" && !containsAnyFold([]string{product.Description, product.Subtitle, product.Policy, product.Usage, strings.Join(product.Tags, " "), strings.Join(product.Included, " ")}, feature) {
			return false
		}
	}
	if filters.MinPrice > 0 && product.BasePrice < filters.MinPrice {
		return false
	}
	if filters.MaxPrice > 0 && product.BasePrice > filters.MaxPrice {
		return false
	}
	if filters.RatingMin > 0 && product.Rating < filters.RatingMin {
		return false
	}
	if filters.InstantConfirm != nil && product.InstantConfirm != *filters.InstantConfirm {
		return false
	}
	if filters.FreeCancel != nil && product.FreeCancel != *filters.FreeCancel {
		return false
	}
	if filters.Date != "" || filters.AvailableToday != nil || filters.AvailableTomorrow != nil || filters.Adults+filters.Children > 0 || filters.VoucherType != "" {
		return hasMatchingPackageAvailability(product, filters)
	}
	return true
}

func sortProducts(products []domain.Product, sortBy string) {
	switch sortBy {
	case "price_asc":
		sort.Slice(products, func(i, j int) bool { return products[i].BasePrice < products[j].BasePrice })
	case "price_desc":
		sort.Slice(products, func(i, j int) bool { return products[i].BasePrice > products[j].BasePrice })
	case "rating":
		sort.Slice(products, func(i, j int) bool { return products[i].Rating > products[j].Rating })
	case "booked":
		sort.Slice(products, func(i, j int) bool { return products[i].BookedCount > products[j].BookedCount })
	case "discount":
		sort.Slice(products, func(i, j int) bool { return maxPackageDiscount(products[i]) > maxPackageDiscount(products[j]) })
	case "distance":
		sort.Slice(products, func(i, j int) bool { return products[i].ID < products[j].ID })
	default:
		sort.Slice(products, func(i, j int) bool {
			left := products[i].Rating*100 + float64(products[i].BookedCount)
			right := products[j].Rating*100 + float64(products[j].BookedCount)
			return left > right
		})
	}
}

func mustJSON(values []string) string {
	b, err := json.Marshal(values)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func parseStringSlice(raw string) []string {
	var values []string
	if err := json.Unmarshal([]byte(raw), &values); err != nil {
		return []string{}
	}
	return values
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func containsFold(values []string, target string) bool {
	for _, value := range values {
		if strings.EqualFold(value, target) || strings.Contains(strings.ToLower(value), strings.ToLower(target)) {
			return true
		}
	}
	return false
}

func demoProducts() []domain.Product {
	return []domain.Product{
		{
			ID: 101, DestinationID: 1, City: "Hangzhou", Category: "Tickets", Type: "ticket", Name: "West Lake Scenic Boat Ticket", Subtitle: "Classic lake cruise with mobile voucher", Description: "Cruise across West Lake, pass Su Causeway and Three Pools Mirroring the Moon, and enjoy a flexible entry window designed for first-time visitors.", Cover: "https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=900", Images: []string{"https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=1200", "https://images.unsplash.com/photo-1547981609-4b6bfe67ca0b?w=1200"}, Tags: []string{"Boat", "Nature", "Mobile voucher"}, Rating: 4.8, ReviewCount: 3260, BookedCount: 128, BasePrice: 88, Currency: "CNY", InstantConfirm: true, FreeCancel: true, Duration: "45 min", MeetingPoint: "West Lake Hubin Pier", Included: []string{"Boat ticket", "Mobile voucher", "Basic route map"}, Excluded: []string{"Hotel pickup", "Meals", "Personal expenses"}, Usage: "Show your mobile voucher and passport at the pier counter.", Policy: "Free cancellation up to 24 hours before travel date.", Status: "active",
		},
		{
			ID: 102, DestinationID: 6, City: "Hangzhou", Category: "Tours", Type: "tour", Name: "Lingyin Temple Half-Day Guided Tour", Subtitle: "English guide, tea village stop, small group", Description: "Explore Lingyin Temple with an English-speaking guide, learn Buddhist culture, and finish with a relaxed Longjing tea village walk.", Cover: "https://images.unsplash.com/photo-1605649487212-47bdab064df7?w=900", Images: []string{"https://images.unsplash.com/photo-1605649487212-47bdab064df7?w=1200", "https://images.unsplash.com/photo-1515488042361-ee00e0ddd4e4?w=1200"}, Tags: []string{"Culture", "English guide", "Tea"}, Rating: 4.9, ReviewCount: 1180, BookedCount: 74, BasePrice: 268, Currency: "CNY", InstantConfirm: true, FreeCancel: true, Duration: "4 hours", MeetingPoint: "Lingyin Temple Main Gate", Included: []string{"English-speaking guide", "Temple entry", "Tea tasting"}, Excluded: []string{"Lunch", "Hotel pickup", "Gratuities"}, Usage: "Meet your guide 10 minutes before departure at the main gate.", Policy: "Free cancellation up to 48 hours before start time.", Status: "active",
		},
		{
			ID: 103, DestinationID: 7, City: "Shanghai", Category: "Tickets", Type: "ticket", Name: "Shanghai Disney Resort 1-Day Ticket", Subtitle: "Instant confirmation for family trips", Description: "Spend a magical day at Shanghai Disney Resort with mobile entry and optional priority add-ons for families.", Cover: "https://images.unsplash.com/photo-1531259683007-906a3792e424?w=900", Images: []string{"https://images.unsplash.com/photo-1531259683007-906a3792e424?w=1200", "https://images.unsplash.com/photo-1543968996-ee822b8176ba?w=1200"}, Tags: []string{"Theme Park", "Family", "Bestseller"}, Rating: 4.7, ReviewCount: 8650, BookedCount: 310, BasePrice: 475, Currency: "CNY", InstantConfirm: true, FreeCancel: false, Duration: "1 day", MeetingPoint: "Shanghai Disney Resort entrance", Included: []string{"1-day park admission", "Mobile voucher"}, Excluded: []string{"Food", "Transport", "Premier access"}, Usage: "Scan the QR voucher with your passport at the park entrance.", Policy: "Non-refundable after booking confirmation.", Status: "active",
		},
		{
			ID: 104, DestinationID: 2, City: "Shanghai", Category: "Experiences", Type: "experience", Name: "Huangpu River Night Cruise", Subtitle: "Bund skyline, mobile voucher, evening sail", Description: "See Shanghai from the water with a night cruise along the Huangpu River, passing the Bund and Pudong skyline.", Cover: "https://images.unsplash.com/photo-1548115184-bc65ee498ad0?w=900", Images: []string{"https://images.unsplash.com/photo-1548115184-bc65ee498ad0?w=1200", "https://images.unsplash.com/photo-1547981609-4b6bfe67ca0b?w=1200"}, Tags: []string{"Night View", "Cruise", "Couples"}, Rating: 4.8, ReviewCount: 5420, BookedCount: 196, BasePrice: 128, Currency: "CNY", InstantConfirm: true, FreeCancel: true, Duration: "50 min", MeetingPoint: "Shiliupu Wharf", Included: []string{"Cruise ticket", "Mobile voucher"}, Excluded: []string{"Hotel pickup", "Food and drinks"}, Usage: "Redeem your voucher at the wharf ticket office before boarding.", Policy: "Free cancellation up to 24 hours before departure.", Status: "active",
		},
		{
			ID: 105, DestinationID: 10, City: "Beijing", Category: "Tours", Type: "tour", Name: "Forbidden City English Guided Walk", Subtitle: "Verified guide with palace highlights", Description: "Walk through the Forbidden City with a licensed English-speaking guide and understand imperial stories, architecture, and daily court life.", Cover: "https://images.unsplash.com/photo-1508807527081-8f81e0f6f8b8?w=900", Images: []string{"https://images.unsplash.com/photo-1508807527081-8f81e0f6f8b8?w=1200", "https://images.unsplash.com/photo-1508804185872-d7badad00f7d?w=1200"}, Tags: []string{"Museum", "History", "English guide"}, Rating: 4.9, ReviewCount: 6320, BookedCount: 152, BasePrice: 198, Currency: "CNY", InstantConfirm: true, FreeCancel: true, Duration: "3 hours", MeetingPoint: "Meridian Gate security entrance", Included: []string{"Admission reservation", "English guide", "Headset"}, Excluded: []string{"Hotel pickup", "Meals"}, Usage: "Bring your passport. The guide will contact you before the tour.", Policy: "Free cancellation up to 48 hours before start time.", Status: "active",
		},
		{
			ID: 106, DestinationID: 3, City: "Beijing", Category: "Transport", Type: "transport", Name: "Mutianyu Great Wall Private Transfer", Subtitle: "Hotel pickup, flexible return, optional ticket", Description: "Book a private transfer from downtown Beijing to Mutianyu Great Wall with a bilingual driver and flexible waiting time.", Cover: "https://images.unsplash.com/photo-1508804185872-d7badad00f7d?w=900", Images: []string{"https://images.unsplash.com/photo-1508804185872-d7badad00f7d?w=1200", "https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?w=1200"}, Tags: []string{"Private transfer", "Great Wall", "Family"}, Rating: 4.8, ReviewCount: 980, BookedCount: 61, BasePrice: 588, Currency: "CNY", InstantConfirm: true, FreeCancel: true, Duration: "8 hours", MeetingPoint: "Your hotel lobby inside Beijing 5th Ring Road", Included: []string{"Private car", "Driver", "Fuel and tolls"}, Excluded: []string{"Great Wall ticket", "Cable car", "Meals"}, Usage: "Driver details will be sent after confirmation.", Policy: "Free cancellation up to 24 hours before pickup.", Status: "active",
		},
		{
			ID: 107, DestinationID: 11, City: "Chengdu", Category: "Tickets", Type: "ticket", Name: "Chengdu Panda Base Morning Ticket", Subtitle: "Best for early panda activity", Description: "Visit Chengdu Panda Base in the morning when pandas are most active, with optional shuttle and family-friendly guidance.", Cover: "https://images.unsplash.com/photo-1535930749574-1399327ce78f?w=900", Images: []string{"https://images.unsplash.com/photo-1535930749574-1399327ce78f?w=1200", "https://images.unsplash.com/photo-1523482580672-f109ba8cb9be?w=1200"}, Tags: []string{"Nature", "Family", "Animals"}, Rating: 4.9, ReviewCount: 4880, BookedCount: 207, BasePrice: 72, Currency: "CNY", InstantConfirm: true, FreeCancel: true, Duration: "3 hours", MeetingPoint: "Chengdu Panda Base ticket gate", Included: []string{"Admission ticket", "Mobile voucher"}, Excluded: []string{"Transport", "Guide"}, Usage: "Use your voucher and passport at the ticket gate.", Policy: "Free cancellation up to 24 hours before visit date.", Status: "active",
		},
		{
			ID: 108, DestinationID: 15, City: "Hangzhou", Category: "Experiences", Type: "experience", Name: "Longjing Tea Picking Workshop", Subtitle: "Tea farm walk, tasting, local host", Description: "Join a Longjing tea host for a field walk, seasonal picking workshop, roasting demonstration, and tasting session.", Cover: "https://images.unsplash.com/photo-1515488042361-ee00e0ddd4e4?w=900", Images: []string{"https://images.unsplash.com/photo-1515488042361-ee00e0ddd4e4?w=1200", "https://images.unsplash.com/photo-1511920170033-f8396924c348?w=1200"}, Tags: []string{"Tea", "Culture", "Local host"}, Rating: 4.8, ReviewCount: 1460, BookedCount: 88, BasePrice: 168, Currency: "CNY", InstantConfirm: true, FreeCancel: true, Duration: "2.5 hours", MeetingPoint: "Longjing Village visitor center", Included: []string{"Tea workshop", "Tasting", "Local host"}, Excluded: []string{"Hotel pickup", "Lunch"}, Usage: "Meet your host at the visitor center and show the mobile voucher.", Policy: "Free cancellation up to 24 hours before start time.", Status: "active",
		},
	}
}

func demoPackages() []domain.Package {
	packages := make([]domain.Package, 0)
	for _, product := range demoProducts() {
		packages = append(packages,
			domain.Package{ID: product.ID*10 + 1, ProductID: product.ID, Name: "Standard", Description: fmt.Sprintf("Standard %s package", strings.ToLower(product.Type)), Price: product.BasePrice, OriginalPrice: product.BasePrice + 30, Unit: "person", MinQuantity: 1, MaxQuantity: 9, AgeRule: "Adult 12+", RefundPolicy: product.Policy, ConfirmType: "instant", VoucherType: "mobile"},
			domain.Package{ID: product.ID*10 + 2, ProductID: product.ID, Name: "Family / Plus", Description: "Better value option for families or travelers who want a smoother experience", Price: product.BasePrice + 58, OriginalPrice: product.BasePrice + 98, Unit: "person", MinQuantity: 1, MaxQuantity: 9, AgeRule: "Adult and child supported", RefundPolicy: product.Policy, ConfirmType: "instant", VoucherType: "mobile"},
		)
	}
	return packages
}

func containsAnyFold(values []string, target string) bool {
	for _, value := range values {
		if strings.Contains(strings.ToLower(value), strings.ToLower(target)) {
			return true
		}
	}
	return false
}

func maxPackageDiscount(product domain.Product) float64 {
	best := 0.0
	for _, pkg := range product.Packages {
		if pkg.OriginalPrice > pkg.Price && pkg.OriginalPrice-pkg.Price > best {
			best = pkg.OriginalPrice - pkg.Price
		}
	}
	return best
}

func hasMatchingPackageAvailability(product domain.Product, filters domain.SearchFilters) bool {
	date := strings.TrimSpace(filters.Date)
	if filters.AvailableToday != nil && *filters.AvailableToday {
		date = time.Now().Format("2006-01-02")
	}
	if filters.AvailableTomorrow != nil && *filters.AvailableTomorrow {
		date = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	}
	quantity := filters.Adults + filters.Children
	if quantity <= 0 {
		quantity = 1
	}
	for _, pkg := range product.Packages {
		if filters.VoucherType != "" && !strings.EqualFold(pkg.VoucherType, filters.VoucherType) {
			continue
		}
		if quantity < pkg.MinQuantity || quantity > pkg.MaxQuantity {
			continue
		}
		for _, availability := range product.Availability {
			if availability.PackageID != pkg.ID || availability.Status != "available" || availability.Stock < quantity {
				continue
			}
			if date != "" && availability.Date != date {
				continue
			}
			return true
		}
	}
	return false
}

func inboundDemoProducts() []domain.Product {
	return []domain.Product{
		{ID: 109, DestinationID: 0, City: "China", Category: "Essentials", Type: "essential", Name: "China eSIM 7-Day Data Pack", Subtitle: "Landing-ready connectivity for overseas travellers", Description: "Install before departure and stay online after arrival in China. Includes English setup guide and hotel-address offline checklist.", Cover: "https://images.unsplash.com/photo-1516321318423-f06f85e504b3?w=1200", Images: []string{"https://images.unsplash.com/photo-1516321318423-f06f85e504b3?w=1200"}, Tags: []string{"Inbound", "eSIM", "WiFi", "English support"}, Rating: 4.7, ReviewCount: 86, BookedCount: 620, BasePrice: 68, Currency: "CNY", InstantConfirm: true, FreeCancel: true, Duration: "7 days", MeetingPoint: "Delivered by email", Included: []string{"7-day data pack", "English setup guide", "China arrival checklist"}, Excluded: []string{"Phone device", "Voice calls"}, Usage: "Install eSIM before departure and activate after landing.", Policy: "Free cancellation before activation.", Status: "active"},
		{ID: 110, DestinationID: 0, City: "Hangzhou", Category: "Transport", Type: "transport", Name: "Shanghai Hongqiao to Hangzhou East Rail Helper", Subtitle: "High-speed rail guidance with station transfer tips", Description: "English rail guidance for Shanghai-Hangzhou trips, including passport checks, station navigation and timing reminders.", Cover: "https://images.unsplash.com/photo-1474487548417-781cb71495f3?w=1200", Images: []string{"https://images.unsplash.com/photo-1474487548417-781cb71495f3?w=1200"}, Tags: []string{"High-speed rail", "Inbound", "Shanghai", "Hangzhou"}, Rating: 4.6, ReviewCount: 42, BookedCount: 310, BasePrice: 73, Currency: "CNY", InstantConfirm: true, FreeCancel: true, Duration: "45-65 min", MeetingPoint: "Shanghai Hongqiao Railway Station", Included: []string{"Rail planning guide", "Chinese station phrases", "Passport timing reminder"}, Excluded: []string{"Physical rail ticket", "Station escort"}, Usage: "Use the guide before departure and keep passport ready for station checks.", Policy: "Planning service refundable before travel date.", Status: "active"},
		{ID: 111, DestinationID: 1, City: "Hangzhou", Category: "Transport", Type: "transport", Name: "Hangzhou Airport to West Lake Transfer", Subtitle: "Private arrival transfer with Chinese hotel confirmation", Description: "Licensed car transfer for overseas travellers arriving at Hangzhou Xiaoshan Airport, with English pickup instructions and Chinese address confirmation.", Cover: "https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?w=1200", Images: []string{"https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?w=1200"}, Tags: []string{"Airport transfer", "Inbound", "Family", "English support"}, Rating: 4.8, ReviewCount: 95, BookedCount: 520, BasePrice: 188, Currency: "CNY", InstantConfirm: true, FreeCancel: true, Duration: "1 hour", MeetingPoint: "Hangzhou Xiaoshan Airport arrivals", Included: []string{"5-seat car", "Driver contact", "Chinese hotel address check"}, Excluded: []string{"Extra waiting time", "Child seat unless selected"}, Usage: "Show mobile voucher and pickup name at arrivals.", Policy: "Free cancellation 24 hours before pickup.", Status: "active"},
		{ID: 112, DestinationID: 3, City: "Beijing", Category: "Transport", Type: "transport", Name: "Beijing Airport Downtown Transfer", Subtitle: "Business van arrival service for long-haul flights", Description: "Airport-to-hotel transfer with terminal confirmation, driver contact and bilingual arrival message.", Cover: "https://images.unsplash.com/photo-1508804185872-d7badad00f7d?w=1200", Images: []string{"https://images.unsplash.com/photo-1508804185872-d7badad00f7d?w=1200"}, Tags: []string{"Airport transfer", "Beijing", "Inbound"}, Rating: 4.7, ReviewCount: 61, BookedCount: 280, BasePrice: 268, Currency: "CNY", InstantConfirm: true, FreeCancel: true, Duration: "1.5 hours", MeetingPoint: "Beijing airport arrivals", Included: []string{"Business van", "Driver contact", "Bilingual pickup note"}, Excluded: []string{"Tolls for route changes"}, Usage: "Confirm terminal and show voucher to driver.", Policy: "Free cancellation 24 hours before pickup.", Status: "active"},
		{ID: 113, DestinationID: 1, City: "Hangzhou", Category: "City Pass", Type: "bundle", Name: "Hangzhou 2-Day Culture Pass", Subtitle: "West Lake, Lingyin route, transfer coupon and AI plan", Description: "A China inbound city pass combining scenic experiences, arrival essentials and AI route planning for Hangzhou.", Cover: "https://images.unsplash.com/photo-1547981609-4b6bfe67ca0b?w=1200", Images: []string{"https://images.unsplash.com/photo-1547981609-4b6bfe67ca0b?w=1200"}, Tags: []string{"City Pass", "Hangzhou", "Inbound", "AI itinerary"}, Rating: 4.9, ReviewCount: 132, BookedCount: 740, BasePrice: 399, Currency: "CNY", InstantConfirm: true, FreeCancel: true, Duration: "2 days", MeetingPoint: "Digital pass", Included: []string{"West Lake cruise", "Lingyin route guide", "Transfer coupon", "AI itinerary"}, Excluded: []string{"Meals unless stated"}, Usage: "Use digital pass QR and follow AI route in My Trips.", Policy: "Free cancellation before first redemption.", Status: "active"},
		{ID: 114, DestinationID: 2, City: "Shanghai", Category: "City Pass", Type: "bundle", Name: "Shanghai Night Pass", Subtitle: "River cruise, night transfer and food street guide", Description: "A night-focused Shanghai pass for overseas travellers with transfer support and bilingual food guidance.", Cover: "https://images.unsplash.com/photo-1548115184-bc65ee498ad0?w=1200", Images: []string{"https://images.unsplash.com/photo-1548115184-bc65ee498ad0?w=1200"}, Tags: []string{"City Pass", "Shanghai", "Night", "Inbound"}, Rating: 4.8, ReviewCount: 77, BookedCount: 410, BasePrice: 299, Currency: "CNY", InstantConfirm: true, FreeCancel: true, Duration: "1 night", MeetingPoint: "The Bund", Included: []string{"Huangpu River cruise", "Night transfer", "Food guide"}, Excluded: []string{"Dinner bill"}, Usage: "Redeem digital pass at each listed stop.", Policy: "Free cancellation before first redemption.", Status: "active"},
	}
}

func inboundDemoPackages() []domain.Package {
	return []domain.Package{
		{ID: 1091, ProductID: 109, Name: "7-day eSIM", Description: "Data pack with English setup guide", Price: 68, OriginalPrice: 88, Unit: "device", MinQuantity: 1, MaxQuantity: 5, AgeRule: "One device per QR", RefundPolicy: "Before activation", ConfirmType: "instant", VoucherType: "mobile"},
		{ID: 1101, ProductID: 110, Name: "Rail helper", Description: "Shanghai-Hangzhou high-speed rail planning", Price: 73, OriginalPrice: 88, Unit: "person", MinQuantity: 1, MaxQuantity: 9, AgeRule: "Passport required", RefundPolicy: "Before travel date", ConfirmType: "instant", VoucherType: "mobile"},
		{ID: 1111, ProductID: 111, Name: "5-seat car", Description: "Airport to West Lake hotels", Price: 188, OriginalPrice: 238, Unit: "car", MinQuantity: 1, MaxQuantity: 4, AgeRule: "Child seat on request", RefundPolicy: "24h free cancellation", ConfirmType: "instant", VoucherType: "mobile"},
		{ID: 1121, ProductID: 112, Name: "Business van", Description: "Airport to downtown Beijing", Price: 268, OriginalPrice: 328, Unit: "van", MinQuantity: 1, MaxQuantity: 6, AgeRule: "Child seat on request", RefundPolicy: "24h free cancellation", ConfirmType: "instant", VoucherType: "mobile"},
		{ID: 1131, ProductID: 113, Name: "2-day pass", Description: "Hangzhou inbound culture pass", Price: 399, OriginalPrice: 488, Unit: "person", MinQuantity: 1, MaxQuantity: 6, AgeRule: "Passport redemption", RefundPolicy: "Before first use", ConfirmType: "instant", VoucherType: "mobile"},
		{ID: 1141, ProductID: 114, Name: "Night pass", Description: "Shanghai night route pass", Price: 299, OriginalPrice: 368, Unit: "person", MinQuantity: 1, MaxQuantity: 6, AgeRule: "Passport redemption", RefundPolicy: "Before first use", ConfirmType: "instant", VoucherType: "mobile"},
	}
}
