<template>
  <div class="travel-home" @mousemove="handleMouseMove">
    <header class="site-header" :class="{ 'header--visible': headerVisible }">
      <router-link to="/" class="header-logo">
        <span class="logo-icon">✈️</span>
        <span>ChinaTravel</span>
      </router-link>

      <nav class="header-nav">
        <a href="#travel-guide" class="header-nav-link" @click.prevent="scrollToGuide">{{ $t('nav.guides') }}</a>
        <router-link to="/trips" class="header-nav-link">{{ $t('nav.myTrips') }}</router-link>
        <a href="#" class="header-nav-link" @click.prevent="scrollToHistory">{{ $t('nav.history') }}</a>
        <a href="#" class="header-nav-link" @click.prevent="scrollToWishlist">{{ $t('nav.wishlist') }}</a>
      </nav>

      <div class="header-actions">
        <button class="map-toggle-header" @click="showMapModal = true">
          <span class="map-icon">🗺️</span>
          <span>{{ $t('ui.map') }}</span>
        </button>
        <button class="action-btn" @click="toggleLang" :title="$t('ui.switchLanguageCurrency')">🌐 {{ locale.toUpperCase() }}</button>
        <div class="currency-dropdown">
          <button class="currency-btn" @click="showCurrencyMenu = !showCurrencyMenu">
            {{ currencySymbol }} {{ currency }}
            <span class="dropdown-arrow">▼</span>
          </button>
          <div class="currency-menu" :class="{ show: showCurrencyMenu }">
            <button v-for="c in currencies" :key="c.code" :class="{ active: currency === c.code }" @click="selectCurrency(c.code)">
              {{ c.symbol }} {{ c.code }} - {{ c.name }}
            </button>
          </div>
        </div>

        <div class="user-profile" v-if="isLoggedIn">
          <router-link to="/account" class="user-name">{{ user?.email }}</router-link>
          <div class="user-avatar">{{ (user?.email || '?')[0].toUpperCase() }}</div>
          <button class="logout-btn" @click="logout">{{ $t('auth.logOut') }}</button>
        </div>
        <button v-else class="signin-btn" @click="showAuthModal = 'login'">{{ $t('auth.signIn') }}</button>
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
              <input v-model="keyword" type="text" class="omnibox-input" :placeholder="searchPlaceholder" @keyup.enter="onSearch" />
            </div>
            <button class="search-submit" @click="onSearch">
              <span>{{ $t('nav.search') }}</span>
            </button>
          </div>
          <div class="hero-popular-tags">
            <span>{{ $t('nav.popular') }}:</span>
            <router-link v-for="city in popularCities" :key="city.id" :to="'/city/' + city.id">{{ $t(city.nameKey) }}</router-link>
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
              <router-link
                :to="'/category/' + cat.id"
                class="sidebar-cat-item"
                :class="{ active: expandedCats.has(cat.id) }"
              >
                <div class="cat-main">
                  <span class="cat-icon">{{ cat.icon }}</span>
                  <span>{{ $t(cat.labelKey) }}</span>
                </div>
                <span v-if="cat.children && cat.children.length" class="cat-chevron" :class="{ rotated: expandedCats.has(cat.id) }">›</span>
              </router-link>

              <div v-if="expandedCats.has(cat.id) && cat.children && cat.children.length" class="cat-children">
                <router-link v-for="child in cat.children" :key="child.id" :to="'/category/' + child.id" class="sidebar-cat-child">
                  {{ $t(child.labelKey) }}
                </router-link>
              </div>
            </div>

            <button v-if="categoryTree.length > 6" class="cat-show-more" :class="{ expanded: showAllCats }" @click="showAllCats = !showAllCats">
              {{ showAllCats ? $t('ui.showLess') : $t('ui.showMore') }}
            </button>
          </div>
        </div>

        <div class="sidebar-widget nearby-widget">
          <h3 class="filter-title">{{ $t('nearby.title') }}</h3>
          <div v-if="nearbyLoading" class="loading">{{ $t('ui.loading') }}</div>
          <div v-else-if="nearbyError" class="error">{{ nearbyError }}</div>
          <div v-else class="nearby-list-unified">
            <router-link v-for="d in nearby.slice(0, 5)" :key="d.id" :to="'/destination/' + d.id" class="nearby-item-unified">
              <div class="name">
                <div class="icon-box">📍</div>
                <span>{{ localizeDestination(d) }}</span>
              </div>
              <span class="dist">{{ d.distance_km }}km</span>
            </router-link>
          </div>
        </div>
      </aside>

      <!-- 主内容区 -->
      <main class="page-main">
        <div class="content-wrap">
          <!-- 首页推荐：一整块卡片区，单卡更友好 -->
          <section class="section recommendations-section">
            <div class="recommendations-block">
              <div class="section-header">
                <div class="header-left">
                  <h2 class="section-title">{{ $t('recommendations.title') }}</h2>
                  <p v-if="isLoggedIn" class="section-subtitle">{{ $t('recommendations.locationHint') }}</p>
                  <p v-else class="section-subtitle section-subtitle--muted">{{ $t('auth.syncHint') }}</p>
                </div>
                <router-link to="/search" class="view-all-link">{{ $t('homeContent.viewAll') }}</router-link>
              </div>
              <div v-if="recLoading" class="loading">{{ $t('ui.loading') }}</div>
              <div v-else-if="recError" class="error">{{ recError }}</div>
              <div
                v-else
                class="card-carousel card-carousel--horizontal card-carousel--auto"
              >
                <div class="carousel-track">
                  <router-link
                    v-for="item in scrollingRecommendations"
                    :key="item.key"
                    :to="'/destination/' + item.destination.id"
                    class="dest-card carousel-item"
                  >
                    <div class="cover-wrap">
                      <img :src="item.destination.cover" :alt="item.destination.name" class="cover" loading="lazy" @error="onRecCoverError(item.destination.id, $event)" />
                      <button type="button" class="fav-btn" :class="{ favorited: item.destination.is_favorite && isLoggedIn }" @click.prevent.stop="toggleFav(item.destination)">{{ (item.destination.is_favorite && isLoggedIn) ? '♥' : '♡' }}</button>
                      <div class="card-badge" v-if="item.visualIndex % 5 === 0">{{ $t('common.rareFind') }}</div>
                    </div>
                    <div class="body">
                      <div class="card-header">
                        <div class="name">{{ localizeDestination(item.destination) }}</div>
                        <div class="rating">★ {{ item.destination.rating }}</div>
                      </div>
                      <div class="meta">{{ localizeCity(item.destination) }}</div>
                      <div class="tags">
                        <span v-for="t in (item.destination.tags || []).slice(0, 2)" :key="t" class="tag">{{ t }}</span>
                      </div>
                      <div class="price">
                        <span class="amount">¥{{ 168 + item.visualIndex * 10 }}</span>
                        <span class="unit">{{ $t('common.night') }}</span>
                      </div>
                      <div class="trust-signal">
                        <span class="reviews">{{ $t('common.reviews', { count: 100 + item.visualIndex * 50 }) }}</span>
                        <span class="booked">{{ $t('common.booked', { count: 14 }) }}</span>
                      </div>
                    </div>
                  </router-link>
                </div>
              </div>
            </div>
          </section>

          <section class="section product-home-section">
            <div class="section-header">
              <div class="header-left">
                <h2 class="section-title">{{ $t('auto.auto_d6d57239') }}</h2>
                <p class="section-subtitle section-subtitle--muted">{{ $t('auto.auto_e6f5b5b1') }}</p>
              </div>
              <router-link :to="activeProductChannel.search" class="view-all-link">{{ $t('auto.auto_ea841d99') }}</router-link>
            </div>
            <div class="product-channel-tabs">
              <button
                v-for="channel in productChannels"
                :key="channel.id"
                type="button"
                :class="{ active: activeProductChannelId === channel.id }"
                @click="activeProductChannelId = channel.id"
              >
                <span>{{ channel.icon }}</span>
                {{ $t(channel.labelKey) }}
              </button>
            </div>
            <div v-if="productsLoading" class="loading">{{ $t('ui.loading') }}</div>
            <div v-else-if="activeChannelProducts.length" class="product-home-grid">
              <ProductCard v-for="product in activeChannelProducts" :key="product.id" :product="product" />
            </div>
            <div v-else class="product-home-empty">
              {{ $t('auto.auto_b4f57855') }}
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
                <router-link v-for="(d, idx) in nearby" :key="'nearby-' + idx" :to="'/destination/' + d.id" class="dest-card carousel-item">
                  <div class="cover-wrap">
                    <img :src="d.cover" :alt="localizeDestination(d)" class="cover" loading="lazy" @error="onImgError" />
                    <div class="card-badge">{{ $t('homeContent.hotBadge') }}</div>
                  </div>
                  <div class="body">
                    <div class="name">{{ localizeDestination(d) }}</div>
                    <div class="meta">{{ localizeCity(d) }} · {{ d.distance_km }}km</div>
                  </div>
                </router-link>
              </div>
            </div>
          </section>

          <!-- 排行榜：最近一周喜欢最多 / 周边点击榜 -->
          <section v-if="trendingThisWeek.length" class="section leaderboard-section">
            <h2 class="section-title">{{ $t('auto.auto_f23dc30e') }}</h2>
            <div class="leaderboard-list">
              <router-link v-for="(d, idx) in trendingThisWeek" :key="'trend-' + d.id" :to="'/destination/' + d.id" class="leaderboard-row">
                <span class="leaderboard-rank" :class="{ 'leaderboard-rank--top': idx < 3 }">{{ idx + 1 }}</span>
                <img :src="d.cover" :alt="localizeDestination(d)" class="leaderboard-thumb" @error="onImgError" />
                <div class="leaderboard-info">
                  <span class="leaderboard-name">{{ localizeDestination(d) }}</span>
                  <span class="leaderboard-meta">{{ localizeCity(d) }}</span>
                </div>
                <button type="button" class="fav-btn fav-btn--small" :class="{ favorited: d.is_favorite && isLoggedIn }" @click.prevent.stop="toggleFav(d)">{{ (d.is_favorite && isLoggedIn) ? '♥' : '♡' }}</button>
              </router-link>
            </div>
          </section>
          <section v-if="mostViewedNearby.length" class="section leaderboard-section">
            <h2 class="section-title">{{ $t('auto.auto_bac57312') }}</h2>
            <div class="leaderboard-list">
              <router-link v-for="(d, idx) in mostViewedNearby" :key="'view-' + d.id" :to="'/destination/' + d.id" class="leaderboard-row">
                <span class="leaderboard-rank" :class="{ 'leaderboard-rank--top': idx < 3 }">{{ idx + 1 }}</span>
                <img :src="d.cover" :alt="localizeDestination(d)" class="leaderboard-thumb" @error="onImgError" />
                <div class="leaderboard-info">
                  <span class="leaderboard-name">{{ localizeDestination(d) }}</span>
                  <span class="leaderboard-meta">{{ localizeCity(d) }}</span>
                </div>
                <button type="button" class="fav-btn fav-btn--small" :class="{ favorited: d.is_favorite && isLoggedIn }" @click.prevent.stop="toggleFav(d)">{{ (d.is_favorite && isLoggedIn) ? '♥' : '♡' }}</button>
              </router-link>
            </div>
          </section>

          <section id="travel-guide" class="section travel-guide">
            <h2 class="section-title">{{ $t('common.travelGuide') }}</h2>
            <div class="guide-grid">
              <div class="guide-card">
                <img src="https://images.unsplash.com/photo-1540959733332-eab4deabeeaf?w=400" alt="Guide 1" class="guide-img" @error="onImgError" />
                <div class="guide-info">
                  <h3>{{ $t('homeContent.streetFoodsTitle') }}</h3>
                  <span>{{ $t('common.readMore') }} →</span>
                </div>
              </div>
              <div class="guide-card">
                <img src="https://images.unsplash.com/photo-1518548419970-58e3b4079ab2?w=400" alt="Guide 2" class="guide-img" @error="onImgError" />
                <div class="guide-info">
                  <h3>{{ $t('homeContent.beijingGemsTitle') }}</h3>
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
                  <span class="blog-tag">{{ $t('homeContent.roadTripTag') }}</span>
                  <h3>{{ $t('homeContent.tibetTitle') }}</h3>
                  <p>{{ $t('homeContent.himalayaDesc') }}</p>
                </div>
              </div>
              <div class="blog-col">
                <div class="blog-card blog-card-sm">
                  <img src="https://images.unsplash.com/photo-1528127269322-539801943592?w=600" class="blog-img" @error="onImgError" />
                  <div class="blog-content">
                    <span class="blog-tag">{{ $t('homeContent.teaTag') }}</span>
                    <h3>{{ $t('homeContent.teaTitle') }}</h3>
                  </div>
                </div>
                <div class="blog-card blog-card-sm">
                  <img src="https://images.unsplash.com/photo-1506377247377-2a5b3b417ebb?w=600" class="blog-img" @error="onImgError" />
                  <div class="blog-content">
                    <span class="blog-tag">{{ $t('homeContent.foodTag') }}</span>
                    <h3>{{ $t('homeContent.foodTitle') }}</h3>
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
              <p class="auth-required-hint">{{ $t('auth.authRequiredHistory') }}</p>
              <button class="auth-required-btn" @click="showAuthModal = 'login'">{{ $t('auth.signIn') }}</button>
            </template>
            <template v-else>
              <div class="tabs tabs-compact">
                <button :class="{ active: activeSidebarTab === 'history' }" @click="activeSidebarTab = 'history'">{{ $t('ui.recent') }}</button>
                <button :class="{ active: activeSidebarTab === 'wishlist' }" @click="activeSidebarTab = 'wishlist'">{{ $t('nav.wishlist') }}</button>
              </div>
              <div v-if="activeSidebarTab === 'history'" class="sidebar-dest-list">
                <router-link v-for="d in history" :key="d.id" :to="'/destination/' + d.id" class="sidebar-dest-row">
                  <img :src="d.cover" :alt="localizeDestination(d)" class="sidebar-dest-thumb" @error="onImgError" />
                  <span class="sidebar-dest-name">{{ localizeDestination(d) }}</span>
                </router-link>
                <p v-if="history.length === 0" class="empty-hint-mini">{{ $t('common.noRecent') }}</p>
              </div>
              <div v-else class="sidebar-dest-list">
                <router-link v-for="d in wishlist" :key="d.id" :to="'/destination/' + d.id" class="sidebar-dest-row">
                  <img :src="d.cover" :alt="localizeDestination(d)" class="sidebar-dest-thumb" @error="onImgError" />
                  <span class="sidebar-dest-name">{{ localizeDestination(d) }}</span>
                </router-link>
                <p v-if="wishlist.length === 0" class="empty-hint-mini">{{ $t('ui.noWishlist') }}</p>
              </div>
            </template>
          </div>
          <!-- Deals Widget -->
          <div class="sidebar-widget deals-widget">
            <h3 class="widget-title">🔥 {{ $t('deals.title') }}</h3>
            <div class="sidebar-deals-list">
              <div v-for="deal in deals" :key="deal.id" class="sidebar-deal-card" :class="'deal-' + deal.type">
                <div class="deal-content-mini">
                  <h4>{{ localizeText(deal.title) }}</h4>
                  <p>{{ localizeText(deal.description) }}</p>
                  <button class="deal-btn-mini">{{ $t('deals.explore') }}</button>
                </div>
              </div>
            </div>
          </div>

        <div class="sidebar-widget inspiration-widget">
          <h3 class="widget-title">{{ $t('common.inspiration') }}</h3>
          <div class="inspiration-list">
            <article class="mini-inspiration">
              <span class="ins-badge">{{ $t('ui.topic') }}</span>
              <h4>{{ $t('homeContent.hangzhou48Title') }}</h4>
              <p>{{ $t('homeContent.hangzhou48Desc') }}</p>
            </article>
            <article class="mini-inspiration">
              <span class="ins-badge hot">{{ $t('ui.hot') }}</span>
              <h4>{{ $t('homeContent.familyFunTitle') }}</h4>
              <p>{{ $t('homeContent.familyFunDesc') }}</p>
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

    <!-- 右侧浮动 AI 小助手 -->
    <AiAssistantBubble
      container-class="ai-float-wrap"
      :hint-text="t('auto.auto_b64c841b')"
      :show-hint="showAiHint"
      :open="aiAssistantOpen"
      :pulse="aiPulse"
      @toggle="toggleAiAssistant"
      @hover-change="pauseAiHint = $event"
    >
      <section class="ai-chat-panel" @click.stop>
        <header class="ai-chat-header">
          <div>
            <span class="ai-chat-kicker">{{ t('auto.auto_49a73ec6') }}</span>
            <h3>{{ t('auto.auto_011b534d') }}</h3>
          </div>
          <button type="button" class="ai-chat-close" @click="closeAiAssistant">×</button>
        </header>

        <div class="ai-chat-messages">
          <article
            v-for="message in aiMessages"
            :key="message.id"
            class="ai-chat-message"
            :class="'ai-chat-message--' + message.role"
          >
            <p>{{ message.text }}</p>
            <div v-if="message.destinations?.length" class="ai-chat-destinations">
              <router-link
                v-for="destination in message.destinations"
                :key="destination.id"
                :to="'/destination/' + destination.id"
                class="ai-chat-destination"
                @click="closeAiAssistant"
              >
                <img :src="destination.cover || FALLBACK_IMAGE" :alt="destination.name" @error="onImgError" />
                <span>
                  <strong>{{ destination.name }}</strong>
                  <small>{{ destination.city }} · {{ destination.rating }}★ · {{ formatPrice(destination.price || 0) }}</small>
                </span>
              </router-link>
            </div>
          </article>
        </div>

        <div class="ai-chat-prompts">
          <button v-for="prompt in aiQuickPrompts" :key="prompt" type="button" @click="sendAiPrompt(prompt)">
            {{ prompt }}
          </button>
        </div>

        <form class="ai-chat-form" @submit.prevent="sendAiPrompt()">
          <input
            v-model="aiQuestion"
            type="text"
            :placeholder="t('auto.auto_d4f59bd8')"
          />
          <button type="submit">{{ t('auto.auto_b87eb8d9') }}</button>
        </form>
      </section>
    </AiAssistantBubble>

    <!-- 信任信号页脚 -->
    <footer class="site-footer">
      <div class="footer-trust-bar">
        <span>✅ {{ $t('trust.verifiedReviews') }}</span>
        <span>🛡️ {{ $t('trust.secureBooking') }}</span>
        <span>🌍 {{ $t('trust.globalSupport') }}</span>
      </div>
      <div class="footer-links">
        <p>© 2026 ChinaTravel, Inc. · Created by Alan Wang · <a href="#">{{ $t('homeContent.privacy') }}</a> · <a href="#">{{ $t('homeContent.terms') }}</a></p>
      </div>
    </footer>

    <!-- Auth Modal -->
    <div v-if="showAuthModal" class="modal-overlay auth-modal-overlay" @click.self="showAuthModal = null">
      <div class="auth-modal-card">
        <button class="modal-close" @click="showAuthModal = null">×</button>
        <template v-if="showAuthModal === 'login'">
          <h2 class="auth-modal-title">{{ $t('auth.signIn') }}</h2>
          <form @submit.prevent="doLogin" class="auth-form">
            <input v-model="authEmail" type="email" :placeholder="$t('auth.email')" required class="auth-input" />
            <input v-model="authPassword" type="password" :placeholder="$t('auth.password')" required class="auth-input" />
            <p v-if="authError" class="auth-error">{{ authError }}</p>
            <button type="submit" class="auth-submit">{{ $t('auth.signIn') }}</button>
            <button type="button" class="auth-link" @click="showAuthModal = 'forgot'">{{ $t('auth.forgotPasswordQuestion') }}</button>
            <button type="button" class="auth-link" @click="showAuthModal = 'register'">{{ $t('auth.createAccount') }}</button>
          </form>
        </template>
        <template v-else-if="showAuthModal === 'register'">
          <h2 class="auth-modal-title">{{ $t('auth.createAccount') }}</h2>
          <form @submit.prevent="doRegister" class="auth-form">
            <input v-model="authEmail" type="email" :placeholder="$t('auth.email')" required class="auth-input" />
            <input v-model="authPassword" type="password" :placeholder="$t('auth.passwordMin')" required minlength="6" class="auth-input" />
            <input v-model="authConfirmPassword" type="password" :placeholder="$t('auth.confirmPassword')" class="auth-input" />
            <p v-if="authError" class="auth-error">{{ authError }}</p>
            <button type="submit" class="auth-submit">{{ $t('auth.register') }}</button>
            <button type="button" class="auth-link" @click="showAuthModal = 'login'">{{ $t('auth.alreadyHaveAccount') }}</button>
          </form>
        </template>
        <template v-else-if="showAuthModal === 'forgot'">
          <h2 class="auth-modal-title">{{ $t('auth.forgotPassword') }}</h2>
          <form @submit.prevent="doForgotPassword" class="auth-form">
            <input v-model="authEmail" type="email" :placeholder="$t('auth.email')" required class="auth-input" />
            <p v-if="authError" class="auth-error">{{ authError }}</p>
            <p v-if="authSuccess" class="auth-success">{{ authSuccess }}</p>
            <button type="submit" class="auth-submit">{{ $t('auth.sendResetLink') }}</button>
            <button type="button" class="auth-link" @click="showAuthModal = 'login'">{{ $t('auth.backToSignIn') }}</button>
          </form>
        </template>
        <template v-else-if="showAuthModal === 'reset'">
          <h2 class="auth-modal-title">{{ $t('auth.resetPassword') }}</h2>
          <form @submit.prevent="doResetPassword" class="auth-form">
            <input v-model="authResetToken" type="text" :placeholder="$t('auth.resetToken')" class="auth-input" />
            <input v-model="authPassword" type="password" :placeholder="$t('auth.newPassword')" required minlength="6" class="auth-input" />
            <input v-model="authConfirmPassword" type="password" :placeholder="$t('auth.confirmNewPassword')" class="auth-input" />
            <p v-if="authError" class="auth-error">{{ authError }}</p>
            <button type="submit" class="auth-submit">{{ $t('auth.resetPassword') }}</button>
            <button type="button" class="auth-link" @click="showAuthModal = 'login'">{{ $t('auth.backToSignIn') }}</button>
          </form>
        </template>
      </div>
    </div>

    <!-- Map Modal -->
    <div v-if="showMapModal" class="modal-overlay" @click.self="showMapModal = false">
      <div class="map-modal">
        <button class="modal-close" @click="showMapModal = false">×</button>
        <h2>{{ t('auto.auto_f30cad85') }}</h2>
        <div class="map-container">
          <div class="map-placeholder">
            <div class="map-markers">
              <div v-for="d in recommendations.slice(0, 5)" :key="d.id" class="map-marker" :style="{ top: (30 + d.id * 10) + '%', left: (20 + d.id * 15) + '%' }">
                <span class="marker-icon">📍</span>
                <span class="marker-label">{{ localizeDestination(d) }}</span>
              </div>
            </div>
            <p class="map-hint">{{ t('auto.auto_b7acfe0e') }}</p>
          </div>
        </div>
        <div class="map-destinations">
          <h3>{{ t('auto.auto_a97839c5') }}</h3>
          <div class="map-dest-grid">
            <router-link v-for="d in recommendations.slice(0, 5)" :key="d.id" :to="'/destination/' + d.id" class="map-dest-card" @click="showMapModal = false">
              <img :src="d.cover" :alt="localizeDestination(d)" />
              <span>{{ localizeDestination(d) }}</span>
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter, useRoute } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import { useCurrency } from '../composables/useCurrency'
import { useTravelAssistant } from '../composables/useTravelAssistant'
import { fetchProducts } from '../composables/useProducts'
import AiAssistantBubble from '../components/AiAssistantBubble.vue'
import ProductCard from '../components/ProductCard.vue'
import { useLocalization } from '../composables/useLocalization'

const { locale, t } = useI18n()
const { localizeText, localizeField, localizeList, localizeDestination, localizeCity } = useLocalization()
const router = useRouter()
const route = useRoute()
const { token, user, isLoggedIn, setAuth, clearAuth, authHeaders } = useAuth()

const showAuthModal = ref(null)
const showMapModal = ref(false)
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
      authError.value = data.error === 'invalid_credentials' ? t('auth.invalidCredentials') : (data.error || t('auth.loginFailed'))
      return
    }
    setAuth(data.token, data.user)
    showAuthModal.value = null
    fetchHomePage()
  } catch (e) {
    authError.value = t('auth.networkError')
  }
}

async function doRegister() {
  authError.value = ''
  if (authPassword.value !== authConfirmPassword.value) {
    authError.value = t('auth.passwordsDoNotMatch')
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
      authError.value = data.error === 'email_already_registered' ? t('auth.emailAlreadyRegistered') : (data.error || t('auth.registrationFailed'))
      return
    }
    authSuccess.value = t('auth.accountCreated')
    showAuthModal.value = 'login'
  } catch (e) {
    authError.value = t('auth.networkError')
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
      authError.value = data.error === 'user_not_found' ? t('auth.noAccount') : (data.error || t('auth.requestFailed'))
      return
    }
    authSuccess.value = t('auth.checkEmailReset')
    if (data.reset_token) authResetToken.value = data.reset_token
    showAuthModal.value = 'reset'
  } catch (e) {
    authError.value = t('auth.networkError')
  }
}

async function doResetPassword() {
  authError.value = ''
  if (authPassword.value !== authConfirmPassword.value) {
    authError.value = t('auth.passwordsDoNotMatch')
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
      authError.value = data.error === 'invalid_or_expired_token' ? t('auth.invalidResetToken') : (data.error || t('auth.resetFailed'))
      return
    }
    authSuccess.value = t('auth.passwordReset')
    showAuthModal.value = 'login'
  } catch (e) {
    authError.value = t('auth.networkError')
  }
}

function logout() {
  fetch(API + '/auth/logout', { method: 'POST', headers: authHeaders() }).catch(() => {})
  clearAuth()
  fetchHomePage()
}

const { currency, setCurrency, formatPrice, getSymbol, currencySymbols } = useCurrency()
const currencySymbol = computed(() => getSymbol())
const showCurrencyMenu = ref(false)
const currencies = [
  { code: 'CNY', symbol: '¥', name: 'Chinese Yuan' },
  { code: 'USD', symbol: '$', name: 'US Dollar' },
  { code: 'EUR', symbol: '€', name: 'Euro' },
  { code: 'GBP', symbol: '£', name: 'British Pound' },
  { code: 'JPY', symbol: '¥', name: 'Japanese Yen' },
  { code: 'KRW', symbol: '₩', name: 'Korean Won' },
  { code: 'THB', symbol: '฿', name: 'Thai Baht' },
  { code: 'SGD', symbol: 'S$', name: 'Singapore Dollar' },
  { code: 'AUD', symbol: 'A$', name: 'Australian Dollar' },
  { code: 'HKD', symbol: 'HK$', name: 'Hong Kong Dollar' }
]

function selectCurrency(code) {
  setCurrency(code)
  showCurrencyMenu.value = false
}

function toggleLang() {
  locale.value = locale.value === 'en' ? 'zh' : 'en'
}

const API = '/api/v1'

const FALLBACK_IMAGE =
  'https://images.unsplash.com/photo-1488646953014-85cb44e25828?auto=format&fit=crop&w=800&q=80'

function onImgError(e) {
  if (e?.target && e.target.src !== FALLBACK_IMAGE) {
    e.target.src = FALLBACK_IMAGE
  }
}

const keyword = ref('')
const popularCities = [
  { id: 'hangzhou', nameKey: 'cityNames.hangzhou' },
  { id: 'shanghai', nameKey: 'cityNames.shanghai' },
  { id: 'beijing', nameKey: 'cityNames.beijing' },
  { id: 'xian', nameKey: 'cityNames.xian' },
  { id: 'chengdu', nameKey: 'cityNames.chengdu' },
]

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
  router.push({ path: '/search', query: { q: k } })
}

const history = ref([])
const wishlist = ref([])
const activeSidebarTab = ref('history')

const categoryTree = ref([
  { id: 'all', icon: '🔥', labelKey: 'auto.auto_e24412d7', children: [] },
  {
    id: 'theme-parks', icon: '🎢', labelKey: 'auto.auto_c644051b',
    children: [
      { id: 'disney', labelKey: 'auto.auto_b70fea5f' },
      { id: 'universal', labelKey: 'auto.auto_97377c59' },
      { id: 'happy-valley', labelKey: 'auto.auto_92cb5a65' }
    ]
  },
  {
    id: 'museums', icon: '🏛️', labelKey: 'auto.auto_c95e9619',
    children: [
      { id: 'history', labelKey: 'auto.auto_a70064bf' },
      { id: 'art', labelKey: 'auto.auto_b149f4d1' },
      { id: 'science', labelKey: 'auto.auto_b1a99c0d' }
    ]
  },
  {
    id: 'camping', icon: '🏕️', labelKey: 'auto.auto_0af4e014',
    children: [
      { id: 'glamping', labelKey: 'auto.auto_1e2153e2' },
      { id: 'rv', labelKey: 'auto.auto_d41f4fba' }
    ]
  },
  {
    id: 'trains', icon: '🚄', labelKey: 'auto.auto_6058d182',
    children: [
      { id: 'high-speed', labelKey: 'auto.auto_4b87c9c1' },
      { id: 'scenic', labelKey: 'auto.auto_fdec6201' }
    ]
  },
  {
    id: 'food', icon: '🍜', labelKey: 'auto.auto_a587f6d2',
    children: [
      { id: 'street', labelKey: 'auto.auto_05606ebc' },
      { id: 'fine-dining', labelKey: 'auto.auto_e2a1e281' }
    ]
  },
  {
    id: 'spas', icon: '💆', labelKey: 'auto.auto_dcc60d90',
    children: [
      { id: 'massage', labelKey: 'auto.auto_771c6989' },
      { id: 'onsen', labelKey: 'auto.auto_5b4ba8e9' }
    ]
  },
  { id: 'nature', icon: '🏔️', labelKey: 'auto.auto_8cbebb8a', children: [] },
  { id: 'shows', icon: '🎭', labelKey: 'auto.auto_aa5020cc', children: [] },
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
    const scrollAmount = 300 + 16
    nearbyActivitiesRef.value.scrollBy({ left: scrollAmount * dir, behavior: 'smooth' })
  }
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

function scrollToGuide() {
  const guide = document.querySelector('#travel-guide, .travel-guide')
  if (guide) {
    guide.scrollIntoView({ behavior: 'smooth', block: 'start' })
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

const homeRecommendations = computed(() => {
  return filteredRecommendations.value.slice(0, 8)
})

const scrollingRecommendations = computed(() => {
  const items = homeRecommendations.value
  if (!items.length) return []

  return [...items, ...items].map((destination, idx) => ({
    destination,
    visualIndex: idx % items.length,
    key: `rec-${destination.id}-${Math.floor(idx / items.length)}`,
  }))
})

const recLoading = ref(true)
const recError = ref('')

const nearby = ref([])
const nearbyLoading = ref(true)
const nearbyError = ref('')

const deals = ref([])
const trendingThisWeek = ref([])
const mostViewedNearby = ref([])
const featuredProducts = ref([])
const productsLoading = ref(true)
const activeProductChannelId = ref('things')
const productChannels = [
  { id: 'stays', labelKey: 'auto.auto_965c784c', icon: '🏨', search: '/search?mode=products&type=stay' },
  { id: 'things', labelKey: 'auto.auto_9595b2ba', icon: '✨', search: '/search?mode=products' },
  { id: 'tickets', labelKey: 'auto.auto_ece59dbb', icon: '🎫', search: '/search?mode=products&type=ticket' },
  { id: 'tours', labelKey: 'auto.auto_dbdd76d8', icon: '🧭', search: '/search?mode=products&type=tour' },
  { id: 'transport', labelKey: 'auto.auto_c0f33179', icon: '🚗', search: '/search?mode=products&type=transport' },
  { id: 'deals', labelKey: 'auto.auto_99d1be73', icon: '🔥', search: '/search?mode=products&free_cancel=true' },
]
const activeProductChannel = computed(() => productChannels.find(channel => channel.id === activeProductChannelId.value) || productChannels[1])
const activeChannelProducts = computed(() => {
  const products = featuredProducts.value
  switch (activeProductChannelId.value) {
    case 'stays':
      return products.filter(product => product.type === 'stay')
    case 'tickets':
      return products.filter(product => product.type === 'ticket')
    case 'tours':
      return products.filter(product => product.type === 'tour')
    case 'transport':
      return products.filter(product => product.type === 'transport')
    case 'deals':
      return products.filter(product => product.free_cancel || product.instant_confirm).slice(0, 6)
    default:
      return products.filter(product => ['ticket', 'tour', 'experience', 'transport'].includes(product.type)).slice(0, 6)
  }
})
const assistantDestinations = computed(() => uniqueDestinations([
  ...recommendations.value,
  ...nearby.value,
  ...trendingThisWeek.value,
  ...mostViewedNearby.value,
  ...history.value,
  ...wishlist.value,
]))
const aiQuestion = ref('')
const aiAssistantOpen = ref(false)
const {
  messages: aiMessages,
  quickPrompts: aiQuickPrompts,
  ask: askTravelAssistant,
  resetGreeting: resetTravelAssistantGreeting,
} = useTravelAssistant({ locale, destinations: assistantDestinations })

function uniqueDestinations(items) {
  const seen = new Set()
  return (items || []).filter((item) => {
    if (!item?.id || seen.has(item.id)) return false
    seen.add(item.id)
    return true
  })
}

function toggleAiAssistant() {
  aiAssistantOpen.value = !aiAssistantOpen.value
  showAiHint.value = false
}

function closeAiAssistant() {
  aiAssistantOpen.value = false
}

function sendAiPrompt(prompt = '') {
  const text = String(prompt || aiQuestion.value).trim()
  if (!text) return
  askTravelAssistant(text)
  aiQuestion.value = ''
  aiAssistantOpen.value = true
  showAiHint.value = false
}

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

async function fetchFeaturedProducts() {
  productsLoading.value = true
  try {
    const data = await fetchProducts({ sort: 'booked' })
    featuredProducts.value = data.results || []
  } catch (e) {
    featuredProducts.value = []
  } finally {
    productsLoading.value = false
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
  'https://images.unsplash.com/photo-1547981609-4b6bfe67ca0b?w=1920',
  'https://images.unsplash.com/photo-1548115184-bc65ee498ad0?w=1920',
  'https://images.unsplash.com/photo-1508804185872-d7badad00f7d?w=1920',
  'https://images.unsplash.com/photo-1525113190471-9969be29263a?w=1920',
  'https://images.unsplash.com/photo-1523482580672-f109ba8cb9be?w=1920',
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
      if (pauseAiHint.value || aiAssistantOpen.value) { scheduleShow(); return }
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
let heroTimer = null
let scrollListener = null
let headerTimer = null

function handleMouseMove() {
  headerVisible.value = true
  if (headerTimer) clearTimeout(headerTimer)
  headerTimer = setTimeout(() => {
    headerVisible.value = false
  }, 2000)
}

async function handleRouteFocus(focus) {
  if (!focus) return
  await nextTick()
  if (focus === 'guide') scrollToGuide()
  if (focus === 'history') scrollToHistory()
  if (focus === 'wishlist') scrollToWishlist()

  const nextQuery = { ...route.query }
  delete nextQuery.focus
  router.replace({ query: nextQuery })
}

watch(locale, () => {
  fetchHomePage()
  resetTravelAssistantGreeting()
})
watch(showAuthModal, () => {
  authError.value = ''
  authSuccess.value = ''
})
watch(locale, () => {
  searchPlaceholderIndex.value = 0
})
watch(() => route.query.focus, (focus) => {
  if (focus) handleRouteFocus(focus)
})

onMounted(() => {
  fetchHomePage()
  fetchFeaturedProducts()
  startAiHintLoop()
  startSearchPlaceholderRotation()
  heroTimer = setInterval(() => {
    heroIndex.value = (heroIndex.value + 1) % heroImages.length
  }, 5000)
  scrollListener = () => {
    heroCollapsed.value = window.scrollY > 120
  }
  window.addEventListener('scroll', scrollListener, { passive: true })
  
  // 点击外部关闭货币菜单
  document.addEventListener('click', handleClickOutside)

  if (route.query.focus) {
    handleRouteFocus(route.query.focus)
  }
})

onUnmounted(() => {
  if (heroTimer) clearInterval(heroTimer)
  if (searchPlaceholderTimer) clearInterval(searchPlaceholderTimer)
  if (scrollListener) window.removeEventListener('scroll', scrollListener)
  if (aiHintTimer) clearTimeout(aiHintTimer)
  if (aiPulseTimer) clearTimeout(aiPulseTimer)
  document.removeEventListener('click', handleClickOutside)
})

function handleClickOutside(e) {
  const dropdown = document.querySelector('.currency-dropdown')
  if (dropdown && !dropdown.contains(e.target)) {
    showCurrencyMenu.value = false
  }
}
</script>
