package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

const DefaultPath = "data/travel.db"

var (
	shared  *sql.DB
	once    sync.Once
	initErr error
)

func Open() (*sql.DB, error) {
	once.Do(func() {
		path := os.Getenv("TRAVEL_DB_PATH")
		if path == "" {
			path = DefaultPath
		}
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			initErr = err
			return
		}
		db, err := sql.Open("sqlite3", path+"?_foreign_keys=on&_busy_timeout=5000")
		if err != nil {
			initErr = err
			return
		}
		db.SetMaxOpenConns(1)
		if err := migrate(db); err != nil {
			_ = db.Close()
			initErr = err
			return
		}
		shared = db
	})
	return shared, initErr
}

func migrate(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS tokens (
			token TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			type TEXT NOT NULL,
			expires_at TEXT NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_tokens_user_type ON tokens(user_id, type);`,
		`CREATE TABLE IF NOT EXISTS favorites (
			user_id TEXT NOT NULL,
			destination_id INTEGER NOT NULL,
			is_favorite INTEGER NOT NULL DEFAULT 1,
			updated_at TEXT NOT NULL,
			PRIMARY KEY(user_id, destination_id)
		);`,
		`CREATE TABLE IF NOT EXISTS history (
			user_id TEXT NOT NULL,
			destination_id INTEGER NOT NULL,
			viewed_at TEXT NOT NULL,
			PRIMARY KEY(user_id, destination_id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_history_user_viewed ON history(user_id, viewed_at DESC);`,
		`CREATE TABLE IF NOT EXISTS stats (
			destination_id INTEGER PRIMARY KEY,
			views INTEGER NOT NULL DEFAULT 0,
			favorites INTEGER NOT NULL DEFAULT 0
		);`,
		`CREATE TABLE IF NOT EXISTS bookings (
			id INTEGER NOT NULL,
			user_id TEXT NOT NULL,
			destination_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			city TEXT NOT NULL,
			cover TEXT NOT NULL,
			check_in TEXT NOT NULL,
			check_out TEXT NOT NULL,
			guests INTEGER NOT NULL,
			total_price REAL NOT NULL,
			status TEXT NOT NULL,
			created_at TEXT NOT NULL,
			cancelled_at TEXT,
			PRIMARY KEY(user_id, id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_bookings_user_created ON bookings(user_id, created_at DESC);`,
		`CREATE TABLE IF NOT EXISTS notifications (
			id INTEGER NOT NULL,
			user_id TEXT NOT NULL,
			type TEXT NOT NULL,
			title TEXT NOT NULL,
			message TEXT NOT NULL,
			link TEXT NOT NULL,
			read INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			PRIMARY KEY(user_id, id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_user_created ON notifications(user_id, created_at DESC);`,
		`CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY,
			destination_id INTEGER NOT NULL,
			city TEXT NOT NULL,
			category TEXT NOT NULL,
			type TEXT NOT NULL,
			name TEXT NOT NULL,
			subtitle TEXT NOT NULL,
			description TEXT NOT NULL,
			cover TEXT NOT NULL,
			images TEXT NOT NULL,
			tags TEXT NOT NULL,
			rating REAL NOT NULL,
			review_count INTEGER NOT NULL,
			booked_count INTEGER NOT NULL,
			base_price REAL NOT NULL,
			currency TEXT NOT NULL,
			instant_confirm INTEGER NOT NULL DEFAULT 1,
			free_cancel INTEGER NOT NULL DEFAULT 1,
			duration TEXT NOT NULL,
			meeting_point TEXT NOT NULL,
			included TEXT NOT NULL,
			excluded TEXT NOT NULL,
			usage TEXT NOT NULL,
			policy TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_products_city_category ON products(city, category);`,
		`CREATE INDEX IF NOT EXISTS idx_products_type_status ON products(type, status);`,
		`CREATE TABLE IF NOT EXISTS product_packages (
			id INTEGER PRIMARY KEY,
			product_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			price REAL NOT NULL,
			original_price REAL NOT NULL,
			unit TEXT NOT NULL,
			min_quantity INTEGER NOT NULL,
			max_quantity INTEGER NOT NULL,
			age_rule TEXT NOT NULL,
			refund_policy TEXT NOT NULL,
			confirm_type TEXT NOT NULL,
			voucher_type TEXT NOT NULL,
			FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_product_packages_product ON product_packages(product_id);`,
		`CREATE TABLE IF NOT EXISTS product_availability (
			id INTEGER PRIMARY KEY,
			product_id INTEGER NOT NULL,
			package_id INTEGER NOT NULL,
			date TEXT NOT NULL,
			price REAL NOT NULL,
			stock INTEGER NOT NULL,
			status TEXT NOT NULL,
			UNIQUE(package_id, date),
			FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE,
			FOREIGN KEY(package_id) REFERENCES product_packages(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_product_availability_product_date ON product_availability(product_id, date);`,
		`CREATE TABLE IF NOT EXISTS orders (
			id INTEGER NOT NULL,
			user_id TEXT NOT NULL,
			status TEXT NOT NULL,
			payment_status TEXT NOT NULL,
			original_amount REAL NOT NULL DEFAULT 0,
			discount_amount REAL NOT NULL DEFAULT 0,
			total_amount REAL NOT NULL,
			coupon_code TEXT NOT NULL DEFAULT '',
			currency TEXT NOT NULL,
			contact_name TEXT NOT NULL,
			contact_email TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			cancelled_at TEXT,
			PRIMARY KEY(user_id, id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_orders_user_created ON orders(user_id, created_at DESC);`,
		`ALTER TABLE orders ADD COLUMN original_amount REAL NOT NULL DEFAULT 0;`,
		`ALTER TABLE orders ADD COLUMN discount_amount REAL NOT NULL DEFAULT 0;`,
		`ALTER TABLE orders ADD COLUMN coupon_code TEXT NOT NULL DEFAULT '';`,
		`CREATE TABLE IF NOT EXISTS order_items (
			id INTEGER NOT NULL,
			order_id INTEGER NOT NULL,
			user_id TEXT NOT NULL,
			product_id INTEGER NOT NULL,
			package_id INTEGER NOT NULL,
			product_name TEXT NOT NULL,
			package_name TEXT NOT NULL,
			city TEXT NOT NULL,
			cover TEXT NOT NULL,
			usage TEXT NOT NULL DEFAULT '',
			travel_date TEXT NOT NULL,
			adults INTEGER NOT NULL,
			children INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			unit_price REAL NOT NULL,
			subtotal REAL NOT NULL,
			PRIMARY KEY(user_id, order_id, id),
			FOREIGN KEY(user_id, order_id) REFERENCES orders(user_id, id) ON DELETE CASCADE
		);`,
		`ALTER TABLE order_items ADD COLUMN usage TEXT NOT NULL DEFAULT '';`,
		`CREATE INDEX IF NOT EXISTS idx_order_items_user_order ON order_items(user_id, order_id);`,
		`CREATE TABLE IF NOT EXISTS coupons (
			id INTEGER PRIMARY KEY,
			code TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			discount_type TEXT NOT NULL,
			discount_value REAL NOT NULL,
			min_spend REAL NOT NULL,
			valid_from TEXT NOT NULL,
			valid_to TEXT NOT NULL,
			usage_limit INTEGER NOT NULL DEFAULT 0,
			status TEXT NOT NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_coupons_code_status ON coupons(code, status);`,
		`CREATE TABLE IF NOT EXISTS reviews (
			id INTEGER PRIMARY KEY,
			user_id TEXT NOT NULL,
			product_id INTEGER NOT NULL,
			order_id INTEGER NOT NULL DEFAULT 0,
			rating REAL NOT NULL,
			quality_score REAL NOT NULL DEFAULT 0,
			service_score REAL NOT NULL DEFAULT 0,
			value_score REAL NOT NULL DEFAULT 0,
			transport_score REAL NOT NULL DEFAULT 0,
			family_score REAL NOT NULL DEFAULT 0,
			content TEXT NOT NULL,
			images TEXT NOT NULL DEFAULT '[]',
			language TEXT NOT NULL DEFAULT 'en',
			verified INTEGER NOT NULL DEFAULT 0,
			merchant_reply TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_reviews_product_created ON reviews(product_id, created_at DESC);`,

		`CREATE TABLE IF NOT EXISTS itineraries (
			id INTEGER NOT NULL,
			user_id TEXT NOT NULL,
			title TEXT NOT NULL,
			city TEXT NOT NULL,
			start_date TEXT NOT NULL,
			end_date TEXT NOT NULL,
			guests INTEGER NOT NULL,
			budget REAL NOT NULL DEFAULT 0,
			status TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			PRIMARY KEY(user_id, id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_itineraries_user_status ON itineraries(user_id, status, start_date);`,
		`CREATE TABLE IF NOT EXISTS itinerary_items (
			id INTEGER NOT NULL,
			itinerary_id INTEGER NOT NULL,
			user_id TEXT NOT NULL,
			day_index INTEGER NOT NULL,
			start_time TEXT NOT NULL,
			end_time TEXT NOT NULL,
			item_type TEXT NOT NULL,
			product_id INTEGER NOT NULL DEFAULT 0,
			destination_id INTEGER NOT NULL DEFAULT 0,
			title TEXT NOT NULL,
			note TEXT NOT NULL,
			estimated_cost REAL NOT NULL DEFAULT 0,
			sort_order INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY(user_id, itinerary_id, id),
			FOREIGN KEY(user_id, itinerary_id) REFERENCES itineraries(user_id, id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_itinerary_items_order ON itinerary_items(user_id, itinerary_id, day_index, sort_order);`,
		`CREATE TABLE IF NOT EXISTS cart_items (
			id INTEGER NOT NULL,
			user_id TEXT NOT NULL,
			product_id INTEGER NOT NULL,
			package_id INTEGER NOT NULL,
			travel_date TEXT NOT NULL,
			adults INTEGER NOT NULL,
			children INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			selected_options TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			PRIMARY KEY(user_id, id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_cart_items_user_created ON cart_items(user_id, created_at DESC);`,
		`CREATE TABLE IF NOT EXISTS merchants (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			contact_email TEXT NOT NULL,
			phone TEXT NOT NULL,
			city TEXT NOT NULL,
			status TEXT NOT NULL,
			rating REAL NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS merchant_products (
			product_id INTEGER PRIMARY KEY,
			merchant_id INTEGER NOT NULL,
			FOREIGN KEY(product_id) REFERENCES products(id) ON DELETE CASCADE,
			FOREIGN KEY(merchant_id) REFERENCES merchants(id) ON DELETE CASCADE
		);`,
		`CREATE INDEX IF NOT EXISTS idx_merchant_products_merchant ON merchant_products(merchant_id);`,
		`CREATE TABLE IF NOT EXISTS refund_requests (
			id INTEGER PRIMARY KEY,
			user_id TEXT NOT NULL,
			order_id INTEGER NOT NULL,
			reason TEXT NOT NULL,
			refund_amount REAL NOT NULL,
			status TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_refund_requests_user_order ON refund_requests(user_id, order_id);`,
		`CREATE TABLE IF NOT EXISTS user_profiles (
			user_id TEXT PRIMARY KEY,
			display_name TEXT NOT NULL,
			avatar TEXT NOT NULL,
			phone TEXT NOT NULL,
			nationality TEXT NOT NULL,
			passport_name TEXT NOT NULL,
			language TEXT NOT NULL,
			currency TEXT NOT NULL,
			travel_preferences TEXT NOT NULL,
			dietary_restrictions TEXT NOT NULL,
			membership_level TEXT NOT NULL,
			points_balance INTEGER NOT NULL DEFAULT 0,
			updated_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS cms_articles (
			id INTEGER PRIMARY KEY,
			slug TEXT NOT NULL UNIQUE,
			title TEXT NOT NULL,
			category TEXT NOT NULL,
			city TEXT NOT NULL,
			language TEXT NOT NULL,
			summary TEXT NOT NULL,
			content TEXT NOT NULL,
			status TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_cms_articles_category_city ON cms_articles(category, city, status);`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			if strings.HasPrefix(statement, `ALTER TABLE`) && strings.Contains(err.Error(), "duplicate column name") {
				continue
			}
			return err
		}
	}
	return nil
}
