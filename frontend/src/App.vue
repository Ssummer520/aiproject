<template>
  <div class="travel-home" @mousemove="handleMouseMove">
    <header class="site-header" :class="{ 'header--visible': headerVisible }">
      <a href="/" class="header-logo">
        <span class="logo-icon">✈️</span>
        <span>ChinaTravel</span>
      </a>
      
      <nav class="header-nav">
        <a href="#" class="header-nav-link">{{ $t('nav.guides') }}</a>
        <a href="#" class="header-nav-link">{{ $t('nav.myTrips') }}</a>
        <a href="#" class="header-nav-link" @click.prevent="scrollToHistory">{{ $t('nav.history') }}</a>
        <a href="#" class="header-nav-link" @click.prevent="scrollToWishlist">{{ $t('nav.wishlist') }}</a>
      </nav>

      <div class="header-actions">
        <button class="map-toggle-header">
          <span class="map-icon">🗺️</span>
          <span>Map</span>
        </button>
        <button class="action-btn" @click="toggleLang" title="Switch Language/Currency">🌐 {{ locale.toUpperCase() }}</button>
        
        <div class="user-profile" v-if="isLoggedIn">
          <div class="user-name">{{ user?.email }}</div>
          <div class="user-avatar">{{ (user?.email || '?')[0].toUpperCase() }}</div>
          <button class="logout-btn" @click="logout">Log out</button>
        </div>
        <button v-else class="signin-btn" @click="showAuthModal = 'login'">Sign in</button>
      </div>
    </header>

    <!-- Airbnb 风格：全屏沉浸式 Hero -->
    <div class="hero" :class="{ 'hero--collapsed': heroCollapsed }">
      <div
        v-for="(img, i) in heroImages"
        :key="img"
        class="hero-bg"
        :class="{ active: i === heroIndex }"
      >
        <img :src="img" :alt="''" class="hero-bg-img" @error="onImgError" />
      </div>
      <div class="hero-overlay" />
      <div class="hero-content">
        <h1 class="page-title">{{ $t('hero.title') }}</h1>
        <p class="page-subtitle">{{ $t('hero.subtitle') }}</p>
        
        <div class="hero-search-container">
          <div class="hero-search-bar">
            <div class="search-input-group">
              <span class="search-icon-static">🔍</span>
              <input v-model="keyword" type="text" class="omnibox-input" :placeholder="searchPlaceholder" />
            </div>
            <button class="search-submit" @click="onSearch">
              <span>{{ $t('nav.search') }}</span>
            </button>
          </div>
          <div class="hero-popular-tags">
            <span>{{ $t('nav.popular') }}:</span>
            <a href="#">Hangzhou</a>
            <a href="#">Shanghai</a>
            <a href="#">Beijing</a>
            <a href="#">Xi'an</a>
            <a href="#">Chengdu</a>
          </div>
        </div>
      </div>
    </div>

    <div class="page-layout">
      <!-- 左侧：过滤器 (Booking 风格) -->
      <aside class="sidebar sidebar-left">
        <!-- Categories Widget (Replaces simple filter) -->
        <div class="sidebar-widget categories-widget">
          <h3 class="filter-title">{{ $t('common.categories') }}</h3>
          <div class="sidebar-category-list">
            <div v-for="cat in visibleCategories" :key="cat.id" class="cat-group">
              <a href="#" 
                 class="sidebar-cat-item" 
                 :class="{ active: expandedCats.has(cat.id) }"
                 @click.prevent="toggleCat(cat.id)"
              >
                <div class="cat-main">
                  <span class="cat-icon">{{ cat.icon }}</span>
                  <span>{{ cat.label }}</span>
                </div>
                <span v-if="cat.children && cat.children.length" class="cat-chevron" :class="{ rotated: expandedCats.has(cat.id) }">›</span>
              </a>
              
              <div v-if="expandedCats.has(cat.id) && cat.children && cat.children.length" class="cat-children">
                <a v-for="child in cat.children" :key="child.id" href="#" class="sidebar-cat-child">
                  {{ child.label }}
                </a>
              </div>
            </div>
            
            <button v-if="categoryTree.length > 6" class="cat-show-more" :class="{ expanded: showAllCats }" @click="showAllCats = !showAllCats">
              {{ showAllCats ? 'Show Less' : 'Show More' }}
            </button>
          </div>
        </div>

        <div class="sidebar-widget nearby-widget">
          <h3 class="filter-title">{{ $t('nearby.title') }}</h3>
          <div v-if="nearbyLoading" class="loading">Loading...</div>
          <div v-else-if="nearbyError" class="error">{{ nearbyError }}</div>
          <div v-else class="nearby-list-unified">
            <a v-for="d in nearby.slice(0, 5)" :key="d.id" class="nearby-item-unified" href="#" @click.prevent="openDetail(d)">
              <div class="name">
                <div class="icon-box">📍</div>
                <span>{{ d.name }}</span>
              </div>
              <span class="dist">{{ d.distance_km }}km</span>
            </a>
          </div>
        </div>
      </aside>

      <!-- 主内容区 -->
      <main class="page-main">
        <div class="content-wrap">
          <!-- Experience Categories (Klook/Trip Style) - REMOVED, MOVED TO HEADER -->
          
          <!-- 首页推荐：一整块卡片区，单卡更友好 -->
          <section class="section recommendations-section">
            <div class="recommendations-block">
              <div class="section-header">
                <div class="header-left">
                  <h2 class="section-title">{{ $t('recommendations.title') }}</h2>
                  <p v-if="isLoggedIn" class="section-subtitle">{{ $t('recommendations.locationHint') }}</p>
                  <p v-else class="section-subtitle section-subtitle--muted">{{ locale === 'zh' ? '登录后可同步收藏与浏览记录' : 'Sign in to sync wishlist & history' }}</p>
                </div>
                <a href="#" class="view-all-link">View all</a>
              </div>
              <div v-if="recLoading" class="loading">Loading...</div>
              <div v-else-if="recError" class="error">{{ recError }}</div>
              <div
                v-else
                class="card-carousel card-carousel--horizontal"
              >
                <div class="carousel-track">
                  <a
                    v-for="(d, idx) in displayRecommendations"
                    :key="'rec-' + idx"
                    class="dest-card carousel-item"
                    href="#"
                    @click.prevent="openDetail(d)"
                  >
                    <div class="cover-wrap">
                      <img :src="d.cover" :alt="d.name" class="cover" loading="lazy" @error="onRecCoverError(d.id, $event)" />
                      <button type="button" class="fav-btn" :class="{ favorited: d.is_favorite && isLoggedIn }" @click.prevent.stop="toggleFav(d)">{{ (d.is_favorite && isLoggedIn) ? '♥' : '♡' }}</button>
                      <div class="card-badge" v-if="idx % 5 === 0">{{ $t('common.rareFind') }}</div>
                    </div>
                    <div class="body">
                      <div class="card-header">
                        <div class="name">{{ d.name }}</div>
                        <div class="rating">★ {{ d.rating }}</div>
                      </div>
                      <div class="meta">{{ d.city }}</div>
                      <div class="tags">
                        <span v-for="t in (d.tags || []).slice(0, 2)" :key="t" class="tag">{{ t }}</span>
                      </div>
                      <div class="price">
                        <span class="amount">¥{{ 168 + idx * 10 }}</span>
                        <span class="unit">{{ $t('common.night') }}</span>
                      </div>
                      <div class="trust-signal">
                        <span class="reviews">{{ $t('common.reviews', { count: 100 + idx * 50 }) }}</span>
                        <span class="booked">{{ $t('common.booked', { count: 14 }) }}</span>
                      </div>
                    </div>
                  </a>
                </div>
              </div>
            </div>
          </section>

          <section class="section hot-activities-section">
            <div class="section-header">
              <div class="header-left">
                <h2 class="section-title">{{ $t('hotActivities.title') }}</h2>
              </div>
            </div>
            <div class="carousel-wrap">
              <button type="button" class="carousel-nav-float carousel-nav-float--left" @click="scrollNearbyActivities(-1)" aria-label="Previous">‹</button>
              <button type="button" class="carousel-nav-float carousel-nav-float--right" @click="scrollNearbyActivities(1)" aria-label="Next">›</button>
              <div ref="nearbyActivitiesRef" class="card-carousel card-carousel--horizontal card-carousel--silky">
                <a v-for="(d, idx) in nearby" :key="'nearby-' + idx" class="dest-card carousel-item" href="#" @click.prevent="openDetail(d)">
                  <div class="cover-wrap">
                    <img :src="d.cover" :alt="d.name" class="cover" loading="lazy" @error="onImgError" />
                    <div class="card-badge">HOT</div>
                  </div>
                  <div class="body">
                    <div class="name">{{ d.name }}</div>
                    <div class="meta">{{ d.city }} · {{ d.distance_km }}km</div>
                  </div>
                </a>
              </div>
            </div>
          </section>

          <!-- 排行榜：最近一周喜欢最多 / 周边点击榜 -->
          <section v-if="trendingThisWeek.length" class="section leaderboard-section">
            <h2 class="section-title">{{ locale === 'zh' ? '本周最爱 · 收藏榜' : 'Trending this week · Most liked' }}</h2>
            <div class="leaderboard-list">
              <a v-for="(d, idx) in trendingThisWeek" :key="'trend-' + d.id" class="leaderboard-row" href="#" @click.prevent="openDetail(d)">
                <span class="leaderboard-rank" :class="{ 'leaderboard-rank--top': idx < 3 }">{{ idx + 1 }}</span>
                <img :src="d.cover" :alt="d.name" class="leaderboard-thumb" @error="onImgError" />
                <div class="leaderboard-info">
                  <span class="leaderboard-name">{{ d.name }}</span>
                  <span class="leaderboard-meta">{{ d.city }}</span>
                </div>
                <button type="button" class="fav-btn fav-btn--small" :class="{ favorited: d.is_favorite && isLoggedIn }" @click.prevent.stop="toggleFav(d)">{{ (d.is_favorite && isLoggedIn) ? '♥' : '♡' }}</button>
              </a>
            </div>
          </section>
          <section v-if="mostViewedNearby.length" class="section leaderboard-section">
            <h2 class="section-title">{{ locale === 'zh' ? '周边人气 · 点击榜' : 'Most viewed nearby' }}</h2>
            <div class="leaderboard-list">
              <a v-for="(d, idx) in mostViewedNearby" :key="'view-' + d.id" class="leaderboard-row" href="#" @click.prevent="openDetail(d)">
                <span class="leaderboard-rank" :class="{ 'leaderboard-rank--top': idx < 3 }">{{ idx + 1 }}</span>
                <img :src="d.cover" :alt="d.name" class="leaderboard-thumb" @error="onImgError" />
                <div class="leaderboard-info">
                  <span class="leaderboard-name">{{ d.name }}</span>
                  <span class="leaderboard-meta">{{ d.city }}</span>
                </div>
                <button type="button" class="fav-btn fav-btn--small" :class="{ favorited: d.is_favorite && isLoggedIn }" @click.prevent.stop="toggleFav(d)">{{ (d.is_favorite && isLoggedIn) ? '♥' : '♡' }}</button>
              </a>
            </div>
          </section>

          <section class="section travel-guide">
            <h2 class="section-title">{{ $t('common.travelGuide') }}</h2>
            <div class="guide-grid">
              <div class="guide-card">
                <img src="https://images.unsplash.com/photo-1540959733332-eab4deabeeaf?w=400" alt="Guide 1" class="guide-img" @error="onImgError" />
                <div class="guide-info">
                  <h3>Top 10 Street Foods in Chengdu</h3>
                  <span>{{ $t('common.readMore') }} →</span>
                </div>
              </div>
              <div class="guide-card">
                <img src="https://images.unsplash.com/photo-1518548419970-58e3b4079ab2?w=400" alt="Guide 2" class="guide-img" @error="onImgError" />
                <div class="guide-info">
                  <h3>Hidden Gems of Beijing</h3>
                  <span>{{ $t('common.readMore') }} →</span>
                </div>
              </div>
            </div>
          </section>

          <!-- Travel Blog/Stories Section -->
          <section class="section travel-blog">
            <h2 class="section-title">{{ $t('common.travelBlog') }}</h2>
            <div class="blog-grid">
              <div class="blog-card blog-card-lg">
                <img src="https://images.unsplash.com/photo-1476514525535-07fb3b4ae5f1?w=800" class="blog-img" @error="onImgError" />
                <div class="blog-content">
                  <span class="blog-tag">Road Trip</span>
                  <h3>Switzerland of the East: A Week in Tibet</h3>
                  <p>Discovering the hidden valleys and sacred mountains of the Himalayas.</p>
                </div>
              </div>
              <div class="blog-col">
                <div class="blog-card blog-card-sm">
                  <img src="https://images.unsplash.com/photo-1528127269322-539801943592?w=600" class="blog-img" @error="onImgError" />
                  <div class="blog-content">
                    <span class="blog-tag">Culture</span>
                    <h3>The Ancient Art of Tea Making in Hangzhou</h3>
                  </div>
                </div>
                <div class="blog-card blog-card-sm">
                  <img src="https://images.unsplash.com/photo-1506377247377-2a5b3b417ebb?w=600" class="blog-img" @error="onImgError" />
                  <div class="blog-content">
                    <span class="blog-tag">Food</span>
                    <h3>Why Spicy Food is Life in Sichuan</h3>
                  </div>
                </div>
              </div>
            </div>
          </section>
    </div>
      </main>

      <!-- 右侧：最近浏览 / 收藏（需登录）+ Deals & 灵感 -->
      <aside class="sidebar sidebar-right">
          <div class="sidebar-widget history-wishlist-widget">
            <h3 class="widget-title">{{ $t('nav.history') }} / {{ $t('nav.wishlist') }}</h3>
            <template v-if="!isLoggedIn">
              <p class="auth-required-hint">Sign in to view your history and wishlist.</p>
              <button class="auth-required-btn" @click="showAuthModal = 'login'">Sign in</button>
            </template>
            <template v-else>
              <div class="tabs tabs-compact">
                <button :class="{ active: activeSidebarTab === 'history' }" @click="activeSidebarTab = 'history'">Recent</button>
                <button :class="{ active: activeSidebarTab === 'wishlist' }" @click="activeSidebarTab = 'wishlist'">Wishlist</button>
              </div>
              <div v-if="activeSidebarTab === 'history'" class="sidebar-dest-list">
                <a v-for="d in history" :key="d.id" class="sidebar-dest-row" href="#" @click.prevent="openDetail(d)">
                  <img :src="d.cover" :alt="d.name" class="sidebar-dest-thumb" @error="onImgError" />
                  <span class="sidebar-dest-name">{{ d.name }}</span>
                </a>
                <p v-if="history.length === 0" class="empty-hint-mini">No recent views</p>
              </div>
              <div v-else class="sidebar-dest-list">
                <a v-for="d in wishlist" :key="d.id" class="sidebar-dest-row" href="#" @click.prevent="openDetail(d)">
                  <img :src="d.cover" :alt="d.name" class="sidebar-dest-thumb" @error="onImgError" />
                  <span class="sidebar-dest-name">{{ d.name }}</span>
                </a>
                <p v-if="wishlist.length === 0" class="empty-hint-mini">No wishlist items</p>
              </div>
            </template>
          </div>
          <!-- Deals Widget -->
          <div class="sidebar-widget deals-widget">
            <h3 class="widget-title">🔥 {{ $t('deals.title') }}</h3>
            <div class="sidebar-deals-list">
              <div v-for="deal in deals" :key="deal.id" class="sidebar-deal-card" :class="'deal-' + deal.type">
                <div class="deal-content-mini">
                  <h4>{{ deal.title }}</h4>
                  <p>{{ deal.description }}</p>
                  <button class="deal-btn-mini">{{ $t('deals.explore') }}</button>
                </div>
              </div>
            </div>
          </div>

        <div class="sidebar-widget inspiration-widget">
          <h3 class="widget-title">{{ $t('common.inspiration') }}</h3>
          <div class="inspiration-list">
            <article class="mini-inspiration">
              <span class="ins-badge">Topic</span>
              <h4>48h in Hangzhou</h4>
              <p>Cycling West Lake & Lingyin Temple</p>
            </article>
            <article class="mini-inspiration">
              <span class="ins-badge hot">Hot</span>
              <h4>Family Fun List</h4>
              <p>Top theme parks in Shanghai & beyond</p>
            </article>
          </div>
        </div>

        <div class="sidebar-widget trust-widget">
          <div class="trust-item">
            <span class="trust-icon">🔒</span>
            <div class="trust-text">
              <strong>{{ $t('trust.securePayment') }}</strong>
              <p>{{ $t('trust.securePaymentDesc') }}</p>
            </div>
          </div>
          <div class="trust-item">
            <span class="trust-icon">🎧</span>
            <div class="trust-text">
              <strong>{{ $t('trust.support') }}</strong>
              <p>{{ $t('trust.supportDesc') }}</p>
            </div>
          </div>
        </div>
      </aside>
    </div>

    <!-- 右侧浮动 AI 小助手：自动轮播提示 + 轻微动效 -->
    <div
      class="ai-float-wrap"
      @mouseenter="pauseAiHint = true"
      @mouseleave="pauseAiHint = false"
    >
      <button
        type="button"
        class="ai-float-btn"
        :class="{ 'ai-float-btn--open': showAiHint, 'ai-float-btn--pulse': aiPulse }"
        @click="onAiFloatClick"
        aria-label="AI travel assistant"
      >
        <span class="ai-float-icon">✨</span>
      </button>
      <Transition name="ai-hint">
        <div v-if="showAiHint" class="ai-float-hint">
          <p class="ai-float-hint-text">{{ locale === 'zh' ? '不知道去哪玩？问我呀' : 'Where to go? Ask me!' }}</p>
          <span class="ai-float-hint-arrow"></span>
        </div>
      </Transition>
    </div>

    <!-- 浮动地图按钮 (欧美用户狂爱) -->
    <!-- Removed as per request, moved to header -->

    <!-- 信任信号页脚 -->
    <footer class="site-footer">
      <div class="footer-trust-bar">
        <span>✅ {{ $t('trust.verifiedReviews') }}</span>
        <span>🛡️ {{ $t('trust.secureBooking') }}</span>
        <span>🌍 {{ $t('trust.globalSupport') }}</span>
      </div>
      <div class="footer-links">
        <p>© 2026 ChinaTravel, Inc. · Created by Alan Wang · <a href="#">Privacy</a> · <a href="#">Terms</a></p>
      </div>
    </footer>

    <!-- Auth Modal: Login / Register / Forgot / Reset -->
    <div v-if="showAuthModal" class="modal-overlay auth-modal-overlay" @click.self="showAuthModal = null">
      <div class="auth-modal-card">
        <button class="modal-close" @click="showAuthModal = null">×</button>
        <template v-if="showAuthModal === 'login'">
          <h2 class="auth-modal-title">Sign in</h2>
          <form @submit.prevent="doLogin" class="auth-form">
            <input v-model="authEmail" type="email" placeholder="Email" required class="auth-input" />
            <input v-model="authPassword" type="password" placeholder="Password" required class="auth-input" />
            <p v-if="authError" class="auth-error">{{ authError }}</p>
            <button type="submit" class="auth-submit">Sign in</button>
            <button type="button" class="auth-link" @click="showAuthModal = 'forgot'">Forgot password?</button>
            <button type="button" class="auth-link" @click="showAuthModal = 'register'">Create account</button>
          </form>
        </template>
        <template v-else-if="showAuthModal === 'register'">
          <h2 class="auth-modal-title">Create account</h2>
          <form @submit.prevent="doRegister" class="auth-form">
            <input v-model="authEmail" type="email" placeholder="Email" required class="auth-input" />
            <input v-model="authPassword" type="password" placeholder="Password (min 6)" required minlength="6" class="auth-input" />
            <input v-model="authConfirmPassword" type="password" placeholder="Confirm password" class="auth-input" />
            <p v-if="authError" class="auth-error">{{ authError }}</p>
            <button type="submit" class="auth-submit">Register</button>
            <button type="button" class="auth-link" @click="showAuthModal = 'login'">Already have an account? Sign in</button>
          </form>
        </template>
        <template v-else-if="showAuthModal === 'forgot'">
          <h2 class="auth-modal-title">Forgot password</h2>
          <form @submit.prevent="doForgotPassword" class="auth-form">
            <input v-model="authEmail" type="email" placeholder="Email" required class="auth-input" />
            <p v-if="authError" class="auth-error">{{ authError }}</p>
            <p v-if="authSuccess" class="auth-success">{{ authSuccess }}</p>
            <button type="submit" class="auth-submit">Send reset link</button>
            <button type="button" class="auth-link" @click="showAuthModal = 'login'">Back to Sign in</button>
          </form>
        </template>
        <template v-else-if="showAuthModal === 'reset'">
          <h2 class="auth-modal-title">Reset password</h2>
          <form @submit.prevent="doResetPassword" class="auth-form">
            <input v-model="authResetToken" type="text" placeholder="Reset token (from email)" class="auth-input" />
            <input v-model="authPassword" type="password" placeholder="New password (min 6)" required minlength="6" class="auth-input" />
            <input v-model="authConfirmPassword" type="password" placeholder="Confirm new password" class="auth-input" />
            <p v-if="authError" class="auth-error">{{ authError }}</p>
            <button type="submit" class="auth-submit">Reset password</button>
            <button type="button" class="auth-link" @click="showAuthModal = 'login'">Back to Sign in</button>
          </form>
        </template>
      </div>
    </div>

    <!-- Destination Detail Modal (Sub-features demo) -->
    <div v-if="showDetailModal" class="modal-overlay" @click="showDetailModal = false">
      <div class="modal-content" @click.stop>
        <button class="modal-close" @click="showDetailModal = false">×</button>
        <div class="modal-body" v-if="currentDest">
          <div class="modal-header-section">
            <img :src="currentDest.cover" :alt="currentDest.name" class="modal-hero-img" @error="onImgError" />
            <div class="modal-title-box">
              <h2>{{ currentDest.name }}</h2>
              <p>{{ currentDest.city }} · ★ {{ currentDest.rating }} ({{ currentDest.review_count }} reviews)</p>
            </div>
          </div>
          
          <div class="modal-main-grid">
            <div class="modal-left-col">
              <div class="modal-section">
                <h3>About this place</h3>
                <p>{{ currentDest.description }}</p>
                <div class="modal-tags">
                  <span v-for="t in currentDest.tags" :key="t" class="modal-tag">{{ t }}</span>
                </div>
              </div>

              <div class="modal-section">
                <h3>Amenities</h3>
                <div class="amenities-grid">
                  <div v-for="a in currentDest.amenities" :key="a" class="amenity-item">
                    <span>{{ a }}</span>
                  </div>
                </div>
              </div>

              <div class="modal-section">
                <h3>Policy</h3>
                <p class="policy-text">{{ currentDest.policy }}</p>
              </div>
              
              <div class="modal-section">
                <h3>Reviews</h3>
                <div class="review-item" v-for="i in 2" :key="i">
                  <div class="review-header">
                    <div class="reviewer-avatar">{{ i === 1 ? 'J' : 'M' }}</div>
                    <div class="reviewer-info">
                      <strong>{{ i === 1 ? 'John Doe' : 'Maria Garcia' }}</strong>
                      <span>March 2026</span>
                    </div>
                  </div>
                  <p>Absolutely amazing experience! The view was breathtaking and the local guides were so helpful. Highly recommend to anyone visiting China.</p>
                </div>
              </div>
            </div>
            
            <div class="modal-right-col">
              <div class="booking-card">
                <div class="booking-header">
                  <span class="price">¥{{ currentDest.price }} <span>/ night</span></span>
                  <span class="rating">★ {{ currentDest.rating }}</span>
                </div>
                <div class="booking-form">
                  <div class="form-row">
                    <div class="form-group">
                      <label>CHECK-IN</label>
                      <input type="text" value="03/15/2026" readonly />
                    </div>
                    <div class="form-group">
                      <label>CHECK-OUT</label>
                      <input type="text" value="03/20/2026" readonly />
                    </div>
                  </div>
                  <div class="form-group full">
                    <label>GUESTS</label>
                    <select><option>1 guest</option><option selected>2 guests</option></select>
                  </div>
                  <button class="reserve-btn">Reserve</button>
                  <p class="reserve-hint">You won't be charged yet</p>
                </div>
                <div class="booking-total">
                  <div class="total-row"><span>¥{{ currentDest.price }} x 5 nights</span> <span>¥{{ currentDest.price * 5 }}</span></div>
                  <div class="total-row"><span>Service fee</span> <span>¥80</span></div>
                  <hr />
                  <div class="total-row grand"><span>Total</span> <span>¥{{ currentDest.price * 5 + 80 }}</span></div>
                </div>
              </div>
              
              <div class="map-widget">
                <h3>Location</h3>
                <div class="mini-map-placeholder">
                  <span>📍 {{ currentDest.lat }}, {{ currentDest.lng }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuth } from './composables/useAuth'

const { locale } = useI18n()
const { token, user, isLoggedIn, setAuth, clearAuth, authHeaders } = useAuth()

const showAuthModal = ref(null)
const authEmail = ref('')
const authPassword = ref('')
const authConfirmPassword = ref('')
const authResetToken = ref('')
const authError = ref('')
const authSuccess = ref('')

async function doLogin() {
  authError.value = ''
  try {
    const res = await fetch(API + '/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: authEmail.value.trim().toLowerCase(), password: authPassword.value }),
    })
    const data = await res.json()
    if (!res.ok) {
      authError.value = data.error === 'invalid_credentials' ? 'Invalid email or password.' : (data.error || 'Login failed')
      return
    }
    setAuth(data.token, data.user)
    showAuthModal.value = null
    fetchHomePage()
  } catch (e) {
    authError.value = 'Network error'
  }
}

async function doRegister() {
  authError.value = ''
  if (authPassword.value !== authConfirmPassword.value) {
    authError.value = 'Passwords do not match'
    return
  }
  try {
    const res = await fetch(API + '/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: authEmail.value.trim().toLowerCase(), password: authPassword.value }),
    })
    const data = await res.json()
    if (!res.ok) {
      authError.value = data.error === 'email_already_registered' ? 'Email already registered.' : (data.error || 'Registration failed')
      return
    }
    authSuccess.value = 'Account created. Sign in below.'
    showAuthModal.value = 'login'
  } catch (e) {
    authError.value = 'Network error'
  }
}

async function doForgotPassword() {
  authError.value = ''
  authSuccess.value = ''
  try {
    const res = await fetch(API + '/auth/forgot-password', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: authEmail.value.trim().toLowerCase() }),
    })
    const data = await res.json()
    if (!res.ok) {
      authError.value = data.error === 'user_not_found' ? 'No account with this email.' : (data.error || 'Request failed')
      return
    }
    authSuccess.value = 'Check your email for reset link. (Demo: use reset_token from response if needed.)'
    if (data.reset_token) authResetToken.value = data.reset_token
    showAuthModal.value = 'reset'
  } catch (e) {
    authError.value = 'Network error'
  }
}

async function doResetPassword() {
  authError.value = ''
  if (authPassword.value !== authConfirmPassword.value) {
    authError.value = 'Passwords do not match'
    return
  }
  try {
    const res = await fetch(API + '/auth/reset-password', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ reset_token: authResetToken.value, new_password: authPassword.value }),
    })
    const data = await res.json()
    if (!res.ok) {
      authError.value = data.error === 'invalid_or_expired_token' ? 'Invalid or expired reset token.' : (data.error || 'Reset failed')
      return
    }
    authSuccess.value = 'Password reset. Sign in below.'
    showAuthModal.value = 'login'
  } catch (e) {
    authError.value = 'Network error'
  }
}

function logout() {
  fetch(API + '/auth/logout', { method: 'POST', headers: authHeaders() }).catch(() => {})
  clearAuth()
  fetchHomePage()
}

function toggleLang() {
  locale.value = locale.value === 'en' ? 'zh' : 'en'
}

const API = '/api/v1'

// 图片加载失败时的兜底图（通用旅行占位）
const FALLBACK_IMAGE =
  'https://images.unsplash.com/photo-1488646953014-85cb44e25828?auto=format&fit=crop&w=800&q=80'

function onImgError(e) {
  if (e?.target && e.target.src !== FALLBACK_IMAGE) {
    e.target.src = FALLBACK_IMAGE
  }
}

const keyword = ref('')

const searchPlaceholderIndex = ref(0)
const searchPlaceholders = {
  en: [
    'Search destinations, attractions, activities',
    'Explore by interest or theme',
    'City, experience, date',
    'Where to? Try Beijing, Hangzhou, Shanghai',
    'Attractions, tours, local experiences',
  ],
  zh: [
    '搜索目的地、景点、活动',
    '按兴趣或主题探索',
    '城市、体验、日期',
    '想去哪？试试北京、杭州、上海',
    '景点、一日游、当地体验',
  ],
}
const searchPlaceholder = computed(() => {
  const list = searchPlaceholders[locale.value] || searchPlaceholders.en
  return list[searchPlaceholderIndex.value % list.length]
})
let searchPlaceholderTimer = null
function startSearchPlaceholderRotation() {
  searchPlaceholderTimer = setInterval(() => {
    const list = searchPlaceholders[locale.value] || searchPlaceholders.en
    searchPlaceholderIndex.value = (searchPlaceholderIndex.value + 1) % list.length
  }, 2800)
}

function onSearch() {
  const k = keyword.value.trim()
  if (!k) return
  // 这里先简单打印，后续可接搜索结果页
  console.log('search:', k)
}

const history = ref([])
const wishlist = ref([])
const activeSidebarTab = ref('history')
const showDetailModal = ref(false)
const currentDest = ref(null)
const categoryTree = ref([
  { id: 'all', icon: '🔥', label: 'All', children: [] },
  { 
    id: 'theme-parks', icon: '🎢', label: 'Theme Parks', 
    children: [
      { id: 'disney', label: 'Disney Resort' },
      { id: 'universal', label: 'Universal Studios' },
      { id: 'happy-valley', label: 'Happy Valley' }
    ] 
  },
  { 
    id: 'museums', icon: '🏛️', label: 'Museums',
    children: [
      { id: 'history', label: 'History Museums' },
      { id: 'art', label: 'Art Galleries' },
      { id: 'science', label: 'Science Centers' }
    ]
  },
  { 
    id: 'camping', icon: '🏕️', label: 'Camping',
    children: [
      { id: 'glamping', label: 'Glamping' },
      { id: 'rv', label: 'RV Parks' }
    ]
  },
  { 
    id: 'trains', icon: '🚄', label: 'Trains',
    children: [
      { id: 'high-speed', label: 'High Speed Rail' },
      { id: 'scenic', label: 'Scenic Routes' }
    ]
  },
  { 
    id: 'food', icon: '🍜', label: 'Food Tours',
    children: [
      { id: 'street', label: 'Street Food' },
      { id: 'fine-dining', label: 'Fine Dining' }
    ]
  },
  { 
    id: 'spas', icon: '💆', label: 'Spas',
    children: [
      { id: 'massage', label: 'Massage' },
      { id: 'onsen', label: 'Onsen / Hot Springs' }
    ]
  },
  { id: 'nature', icon: '🏔️', label: 'Nature', children: [] },
  { id: 'shows', icon: '🎭', label: 'Shows', children: [] },
])

const expandedCats = ref(new Set(['all']))
const showAllCats = ref(false)

function toggleCat(id) {
  if (expandedCats.value.has(id)) {
    expandedCats.value.delete(id)
  } else {
    expandedCats.value.add(id)
  }
}

const visibleCategories = computed(() => {
  return showAllCats.value ? categoryTree.value : categoryTree.value.slice(0, 6)
})

const nearbyActivitiesRef = ref(null)

function scrollNearbyActivities(dir) {
  if (nearbyActivitiesRef.value) {
    const scrollAmount = 300 + 16 // itemWidth + gap
    nearbyActivitiesRef.value.scrollBy({ left: scrollAmount * dir, behavior: 'smooth' })
  }
}

function openDetail(d) {
  if (!isLoggedIn.value) {
    showAuthModal.value = 'login'
    return
  }
  currentDest.value = d
  showDetailModal.value = true
  goDest(d) // Record view
}

function scrollToWishlist() {
  if (!isLoggedIn.value) {
    showAuthModal.value = 'login'
    return
  }
  activeSidebarTab.value = 'wishlist'
  const sidebar = document.querySelector('.sidebar-right')
  if (sidebar) {
    sidebar.scrollIntoView({ behavior: 'smooth', block: 'center' })
  }
}

function scrollToHistory() {
  if (!isLoggedIn.value) {
    showAuthModal.value = 'login'
    return
  }
  activeSidebarTab.value = 'history'
  const sidebar = document.querySelector('.sidebar-right')
  if (sidebar) {
    sidebar.scrollIntoView({ behavior: 'smooth', block: 'center' })
  }
}

const recommendations = ref([])
const brokenRecommendationCoverIDs = ref(new Set())

function onRecCoverError(id, e) {
  if (e?.target && e.target.src !== FALLBACK_IMAGE) {
    e.target.src = FALLBACK_IMAGE
    brokenRecommendationCoverIDs.value = new Set([...brokenRecommendationCoverIDs.value, id])
  }
}

const filteredRecommendations = computed(() => {
  return (recommendations.value || []).filter((d) => d?.id && d?.cover)
})

const displayRecommendations = computed(() => {
  return filteredRecommendations.value.length ? [...filteredRecommendations.value, ...filteredRecommendations.value] : []
})
const recLoading = ref(true)
const recError = ref('')

const nearby = ref([])
const nearbyLoading = ref(true)
const nearbyError = ref('')

const deals = ref([])
const trendingThisWeek = ref([])
const mostViewedNearby = ref([])

async function fetchHomePage() {
  recLoading.value = true
  nearbyLoading.value = true
  try {
    const res = await fetch(API + '/home', {
      headers: { 'Accept-Language': locale.value, ...authHeaders() },
    })
    const data = await res.json()
    recommendations.value = data.recommendations || []
    nearby.value = data.nearby || []
    deals.value = data.deals || []
    history.value = data.history || []
    wishlist.value = data.wishlist || []
    trendingThisWeek.value = data.trending_this_week || []
    mostViewedNearby.value = data.most_viewed_nearby || []
  } catch (e) {
    recError.value = 'Failed to load home page'
  } finally {
    recLoading.value = false
    nearbyLoading.value = false
  }
}

function goDest(d) {
  if (!isLoggedIn.value) {
    showAuthModal.value = 'login'
    return
  }
  fetch(`${API}/destinations/${d.id}/view`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...authHeaders() },
  }).then((res) => {
    if (res.status === 401) {
      clearAuth()
      showAuthModal.value = 'login'
      return
    }
    fetchHomePage()
  })
}

async function toggleFav(d) {
  if (!isLoggedIn.value) {
    showAuthModal.value = 'login'
    return
  }
  try {
    const res = await fetch(`${API}/destinations/${d.id}/favorite`, {
      method: 'POST',
      headers: authHeaders(),
    })
    const data = await res.json()
    if (res.status === 401) {
      clearAuth()
      showAuthModal.value = 'login'
      return
    }
    if (data.ok) {
      d.is_favorite = data.is_favorite
      fetchHomePage()
    }
  } catch (e) {
    console.error(e)
  }
}

const heroImages = [
  'https://images.unsplash.com/photo-1547981609-4b6bfe67ca0b?w=1920', // West Lake
  'https://images.unsplash.com/photo-1548115184-bc65ee498ad0?w=1920', // Shanghai
  'https://images.unsplash.com/photo-1508804185872-d7badad00f7d?w=1920', // Great Wall
  'https://images.unsplash.com/photo-1525113190471-9969be29263a?w=1920', // Yellow Mountain
  'https://images.unsplash.com/photo-1523482580672-f109ba8cb9be?w=1920', // Xi'an
]
const heroIndex = ref(0)
const heroCollapsed = ref(false)
const headerVisible = ref(false)
const showAiHint = ref(false)
const pauseAiHint = ref(false)
const aiPulse = ref(false)
let aiHintTimer = null
let aiPulseTimer = null

function startAiHintLoop() {
  function scheduleShow() {
    if (aiHintTimer) clearTimeout(aiHintTimer)
    aiHintTimer = setTimeout(() => {
      if (pauseAiHint.value) { scheduleShow(); return }
      showAiHint.value = true
      aiPulse.value = true
      if (aiPulseTimer) clearTimeout(aiPulseTimer)
      aiPulseTimer = setTimeout(() => { aiPulse.value = false }, 600)
      aiHintTimer = setTimeout(() => {
        showAiHint.value = false
        aiHintTimer = setTimeout(scheduleShow, 8000)
      }, 4500)
    }, 1500)
  }
  scheduleShow()
}

function onAiFloatClick() {
  showAiHint.value = !showAiHint.value
  if (showAiHint.value) console.log('AI assistant coming soon')
}
let heroTimer = null
let scrollListener = null
let headerTimer = null

function handleMouseMove() {
  headerVisible.value = true
  if (headerTimer) clearTimeout(headerTimer)
  headerTimer = setTimeout(() => {
    headerVisible.value = false
  }, 2000) // Hide after 2s of no movement
}

watch(locale, () => {
  fetchHomePage()
})
watch(showAuthModal, () => {
  authError.value = ''
  authSuccess.value = ''
})
watch(locale, () => {
  searchPlaceholderIndex.value = 0
})

onMounted(() => {
  fetchHomePage()
  startAiHintLoop()
  startSearchPlaceholderRotation()
  heroTimer = setInterval(() => {
    heroIndex.value = (heroIndex.value + 1) % heroImages.length
  }, 5000)
  scrollListener = () => {
    heroCollapsed.value = window.scrollY > 120
  }
  window.addEventListener('scroll', scrollListener, { passive: true })
})

onUnmounted(() => {
  if (heroTimer) clearInterval(heroTimer)
  if (searchPlaceholderTimer) clearInterval(searchPlaceholderTimer)
  if (scrollListener) window.removeEventListener('scroll', scrollListener)
  if (aiHintTimer) clearTimeout(aiHintTimer)
  if (aiPulseTimer) clearTimeout(aiPulseTimer)
})
</script>

<style scoped>
.meta-hint {
  font-size: 0.9rem;
  color: var(--text-muted);
  margin: -8px 0 16px;
}
</style>
