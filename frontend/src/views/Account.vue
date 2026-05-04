<template>
  <div class="account-page">
    <header class="site-header">
      <router-link to="/" class="header-logo">
        <span class="logo-icon">✈️</span>
        <span>ChinaTravel</span>
      </router-link>
      <nav class="header-nav">
        <router-link to="/" class="header-nav-link" :class="{ 'is-active': route.path === '/' }">{{ $t('nav.home') }}</router-link>
        <router-link to="/search" class="header-nav-link" :class="{ 'is-active': route.path.startsWith('/search') }">{{ $t('nav.search') }}</router-link>
        <router-link to="/trips" class="header-nav-link" :class="{ 'is-active': route.path.startsWith('/trips') }">{{ $t('nav.myTrips') }}</router-link>
        <router-link to="/account" class="header-nav-link" :class="{ 'is-active': route.path.startsWith('/account') }">{{ $t('nav.account') }}</router-link>
      </nav>
      <div class="header-actions">
        <button class="action-btn" @click="toggleLang">🌐 {{ locale.toUpperCase() }}</button>
        <div class="user-profile" v-if="isLoggedIn">
          <router-link to="/account" class="user-name">{{ user?.email }}</router-link>
          <div class="user-avatar">{{ accountInitial }}</div>
        </div>
        <button v-else class="signin-btn" @click="showAuthModal = 'login'">{{ $t('auth.signIn') }}</button>
      </div>
    </header>

    <main class="account-content">
      <div class="page-header">
        <span class="section-kicker">{{ t('userContext.sectionLabel') }}</span>
        <h1>{{ t('userContext.title') }}</h1>
        <p>{{ t('userContext.subtitle') }}</p>
      </div>

      <div v-if="!isLoggedIn" class="auth-prompt">
        <p>{{ $t('auto.auto_6bb5f99a') }}</p>
        <button class="auth-btn" @click="showAuthModal = 'login'">{{ $t('auto.auto_a4f0af04') }}</button>
      </div>

      <template v-else>
        <section class="account-hero account-card">
          <div class="hero-main">
            <div class="hero-user">
              <div class="user-avatar-large">{{ accountInitial }}</div>
              <div>
                <span class="section-kicker">{{ t('userContext.profileSnapshot') }}</span>
                <h2>{{ profile.display_name || user?.email }}</h2>
                <p>{{ user?.email }} · ID {{ user?.id }}</p>
                <div class="hero-badges">
                  <span>{{ membership.level || 'Silver' }}</span>
                  <span>{{ profileCompletion }}% {{ t('userContext.profileComplete') }}</span>
                </div>
              </div>
            </div>
            <div class="hero-actions">
              <router-link to="/trips" class="primary-cta">{{ t('userContext.viewTrips') }}</router-link>
              <button class="secondary-cta" type="button" @click="openTravelerModal()">{{ t('userContext.addTraveler') }}</button>
            </div>
          </div>
          <aside class="readiness-card">
            <span class="section-kicker">{{ t('userContext.travelReady') }}</span>
            <h3>{{ t('userContext.readyTitle') }}</h3>
            <div class="readiness-list">
              <div class="readiness-item">
                <span class="status-dot" :class="{ active: defaultTraveler }"></span>
                <div>
                  <strong>{{ t('userContext.readyDefaultTraveler') }}</strong>
                  <small>{{ primaryTraveler ? primaryTraveler.name : t('userContext.readyNoTraveler') }}</small>
                </div>
              </div>
              <div class="readiness-item">
                <span class="status-dot active"></span>
                <div>
                  <strong>{{ t('userContext.readyLanguageCurrency') }}</strong>
                  <small>{{ locale.toUpperCase() }} · {{ currency }} {{ currencySymbol }}</small>
                </div>
              </div>
              <div class="readiness-item">
                <span class="status-dot" :class="{ active: emailNotifications || pushNotifications }"></span>
                <div>
                  <strong>{{ t('userContext.readyNotifications') }}</strong>
                  <small>{{ emailNotifications || pushNotifications ? t('userContext.readyOn') : t('userContext.readyOff') }}</small>
                </div>
              </div>
            </div>
          </aside>
        </section>

        <section class="quick-actions">
          <router-link to="/trips" class="quick-action-card">
            <span class="quick-icon">🧳</span>
            <strong>{{ t('userContext.quickTripsTitle') }}</strong>
            <small>{{ t('userContext.quickTripsSubtitle') }}</small>
          </router-link>
          <button class="quick-action-card" type="button" @click="openTravelerModal()">
            <span class="quick-icon">🪪</span>
            <strong>{{ t('userContext.quickTravelersTitle') }}</strong>
            <small>{{ t('userContext.quickTravelersSubtitle') }}</small>
          </button>
          <a class="quick-action-card" href="#membership-settings">
            <span class="quick-icon">💎</span>
            <strong>{{ t('userContext.quickMembershipTitle') }}</strong>
            <small>{{ t('userContext.quickMembershipSubtitle') }}</small>
          </a>
          <a class="quick-action-card" href="#profile-settings">
            <span class="quick-icon">⚙️</span>
            <strong>{{ t('userContext.quickPreferencesTitle') }}</strong>
            <small>{{ t('userContext.quickPreferencesSubtitle') }}</small>
          </a>
        </section>

        <section class="account-card">
          <div class="card-head">
            <div>
              <span class="section-kicker">{{ t('userContext.travelerProfiles') }}</span>
              <h2>{{ t('userContext.travelerProfilesTitle') }}</h2>
              <p>{{ t('userContext.travelerSecurityHint') }}</p>
            </div>
            <button class="auth-submit add-traveler-btn" type="button" @click="openTravelerModal()">{{ t('userContext.addTraveler') }}</button>
          </div>
          <div v-if="travelersLoading" class="loading-state"><div class="spinner"></div></div>
          <div v-else-if="!travelers.length" class="traveler-empty">
            <div>
              <strong>{{ t('userContext.noTravelers') }}</strong>
              <p>{{ t('userContext.noTravelersHint') }}</p>
            </div>
            <button class="secondary-cta" type="button" @click="openTravelerModal()">{{ t('userContext.addTraveler') }}</button>
          </div>
          <div v-else class="traveler-grid">
            <article v-for="traveler in travelers" :key="traveler.id" class="traveler-card">
              <div>
                <div class="traveler-title">
                  <strong>{{ traveler.name }}</strong>
                  <span v-if="traveler.is_default">{{ t('userContext.defaultTraveler') }}</span>
                </div>
                <p>{{ traveler.document_type }} · {{ traveler.document_no_masked }}</p>
                <small>{{ traveler.nationality || t('userContext.nationalityFallback') }} · {{ traveler.birth_date || t('userContext.birthDateFallback') }}</small>
              </div>
              <div class="traveler-actions">
                <button type="button" @click="openTravelerModal(traveler)">{{ t('userContext.edit') }}</button>
                <button type="button" :disabled="traveler.is_default" @click="makeDefaultTraveler(traveler)">{{ t('userContext.setDefault') }}</button>
                <button type="button" class="danger-link" @click="removeTraveler(traveler)">{{ t('userContext.delete') }}</button>
              </div>
            </article>
          </div>
          <p v-if="travelerMessage" class="success-text">{{ travelerMessage }}</p>
          <p v-if="travelerError" class="auth-error">{{ travelerError }}</p>
        </section>

        <div id="membership-settings" class="account-grid">
          <section class="account-card">
            <div class="card-head">
              <div>
                <span class="section-kicker">{{ t('userContext.membership') }}</span>
                <h2>{{ t('userContext.membershipBenefits') }}</h2>
              </div>
              <strong>{{ membership.level || 'Silver' }}</strong>
            </div>
            <div class="membership-panel">
              <div class="points-ring">{{ membership.points_balance || 0 }}</div>
              <div>
                <p>{{ t('userContext.validUntil', { date: membership.valid_until || '—' }) }}</p>
                <div class="tag-row">
                  <span v-for="benefit in membership.benefits || []" :key="benefit">{{ benefit }}</span>
                </div>
              </div>
            </div>
          </section>

          <section class="account-card">
            <div class="card-head">
              <div>
                <span class="section-kicker">{{ t('userContext.preferences') }}</span>
                <h2>{{ t('userContext.fastSettings') }}</h2>
              </div>
            </div>
            <div class="settings-list compact">
              <div class="settings-item" @click="toggleLang">
                <span>{{ $t('auto.auto_729bdc9e') }}</span>
                <span class="arrow">{{ locale.toUpperCase() }} ↝</span>
              </div>
              <div class="settings-item" @click="showCurrencyModal = true">
                <span>{{ $t('auto.auto_f2bafa5b') }}</span>
                <span class="arrow">{{ currency }} {{ currencySymbol }} ↝</span>
              </div>
              <div class="settings-item">
                <span>{{ $t('auto.auto_cb67c40f') }}</span>
                <label class="toggle"><input v-model="emailNotifications" type="checkbox" /><span class="slider"></span></label>
              </div>
              <div class="settings-item">
                <span>{{ $t('auto.auto_4396fb56') }}</span>
                <label class="toggle"><input v-model="pushNotifications" type="checkbox" /><span class="slider"></span></label>
              </div>
            </div>
          </section>
        </div>

        <section id="profile-settings" class="account-card profile-card">
          <div class="profile-summary">
            <div>
              <span class="section-kicker">{{ t('userContext.userProfile') }}</span>
              <h2>{{ t('userContext.profileTitle') }}</h2>
              <p>{{ t('userContext.profileEditHint') }}</p>
            </div>
            <div class="completion-pill">{{ profileCompletion }}% {{ t('userContext.profileComplete') }}</div>
          </div>
          <details class="profile-details">
            <summary>{{ t('userContext.editProfile') }}</summary>
            <div class="card-head">
              <div>
                <span class="section-kicker">{{ t('userContext.preferences') }}</span>
              </div>
              <button class="mini-btn" type="button" :disabled="profileSaving" @click="saveProfile">
                {{ profileSaving ? t('userContext.saving') : t('userContext.saveProfile') }}
              </button>
            </div>
            <form class="profile-form" @submit.prevent="saveProfile">
              <input v-model="profileForm.display_name" class="auth-input" :placeholder="t('userContext.displayName')" />
              <input v-model="profileForm.avatar" class="auth-input" :placeholder="t('userContext.avatarUrl')" />
              <input v-model="profileForm.phone" class="auth-input" :placeholder="t('userContext.phone')" />
              <input v-model="profileForm.nationality" class="auth-input" :placeholder="t('userContext.nationality')" />
              <select v-model="profileForm.language" class="auth-input">
                <option value="en">English</option>
                <option value="zh">{{ t('auto.auto_27588cf2') }}</option>
              </select>
              <select v-model="profileForm.currency" class="auth-input">
                <option v-for="(_, code) in currencySymbols" :key="code" :value="code">{{ code }}</option>
              </select>
              <input v-model="profileForm.travel_style" class="auth-input" :placeholder="t('userContext.travelStyle')" />
              <input v-model="profileForm.budget_level" class="auth-input" :placeholder="t('userContext.budgetLevel')" />
              <input v-model="profileForm.family_type" class="auth-input" :placeholder="t('userContext.familyType')" />
              <input v-model="dietaryText" class="auth-input" :placeholder="t('userContext.dietaryRestrictions')" />
              <input v-model="accessibilityText" class="auth-input" :placeholder="t('userContext.accessibilityNeeds')" />
            </form>
            <p v-if="profileMessage" class="success-text">{{ profileMessage }}</p>
          </details>
        </section>

        <section class="account-card">
          <h2>{{ $t('auto.auto_3ed40bf3') }}</h2>
          <div v-if="notificationsLoading" class="loading-state"><div class="spinner"></div></div>
          <div v-else-if="notifications.length === 0" class="empty-state"><p>{{ $t('auto.auto_d8298dc2') }}</p></div>
          <div v-else class="notifications-list">
            <div v-for="n in notifications" :key="n.id" class="notification-item" :class="{ unread: !n.read }" @click="markNotificationAsRead(n)">
              <div class="notification-icon">{{ n.type === 'booking_confirmed' ? '✓' : '🔔' }}</div>
              <div class="notification-content">
                <strong>{{ n.title }}</strong>
                <p>{{ n.message }}</p>
                <span class="notification-time">{{ n.created_at }}</span>
              </div>
            </div>
          </div>
        </section>

        <section class="account-card danger-zone">
          <h2>{{ $t('auto.auto_38f922a3') }}</h2>
          <button class="btn-danger" @click="doLogout">{{ $t('auto.auto_c44a476c') }}</button>
        </section>
      </template>
    </main>

    <div v-if="showTravelerModal" class="modal-overlay" @click.self="closeTravelerModal">
      <div class="auth-modal-card traveler-modal-card">
        <button class="modal-close" @click="closeTravelerModal">×</button>
        <h2 class="auth-modal-title">{{ travelerForm.id ? t('userContext.editTravelerTitle') : t('userContext.addTravelerTitle') }}</h2>
        <form class="auth-form" @submit.prevent="saveTraveler">
          <input v-model="travelerForm.name" class="auth-input" :placeholder="t('userContext.name')" required />
          <select v-model="travelerForm.gender" class="auth-input">
            <option value="">{{ t('userContext.gender') }}</option>
            <option value="male">{{ t('userContext.male') }}</option>
            <option value="female">{{ t('userContext.female') }}</option>
            <option value="other">{{ t('userContext.other') }}</option>
          </select>
          <input v-model="travelerForm.birth_date" class="auth-input" type="date" />
          <select v-model="travelerForm.document_type" class="auth-input" required>
            <option value="PASSPORT">Passport</option>
            <option value="ID_CARD">ID Card</option>
            <option value="OTHER">Other</option>
          </select>
          <input v-model="travelerForm.document_no" class="auth-input" :placeholder="travelerForm.id ? t('userContext.documentNoOptional') : t('userContext.documentNo')" :required="!travelerForm.id" />
          <input v-model="travelerForm.nationality" class="auth-input" :placeholder="t('userContext.nationality')" />
          <input v-model="travelerForm.phone" class="auth-input" :placeholder="t('userContext.phone')" />
          <input v-model="travelerForm.email" class="auth-input" type="email" :placeholder="t('userContext.email')" />
          <label class="default-check"><input v-model="travelerForm.is_default" type="checkbox" /> {{ t('userContext.markDefault') }}</label>
          <p v-if="travelerError" class="auth-error">{{ travelerError }}</p>
          <button class="auth-submit" type="submit" :disabled="travelerSaving">{{ travelerSaving ? t('userContext.saving') : t('userContext.saveTraveler') }}</button>
        </form>
      </div>
    </div>

    <div v-if="showCurrencyModal" class="modal-overlay" @click.self="showCurrencyModal = false">
      <div class="auth-modal-card">
        <button class="modal-close" @click="showCurrencyModal = false">×</button>
        <h2 class="auth-modal-title">{{ t('auto.auto_61d71b04') }}</h2>
        <div class="currency-options">
          <button v-for="(symbol, code) in currencySymbols" :key="code" type="button" class="currency-option" :class="{ active: currency === code }" @click="selectCurrency(code)">
            <span>{{ code }}</span>
            <span>{{ symbol }}</span>
          </button>
        </div>
      </div>
    </div>

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
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import { useCurrency } from '../composables/useCurrency'
import { createTraveler, deleteTraveler, fetchMembership, fetchProfile, fetchRoles, fetchTravelers, setDefaultTraveler, updateProfile, updateTraveler } from '../composables/useUser'

const { locale, t } = useI18n()
const route = useRoute()
const router = useRouter()
const { isLoggedIn, user, setAuth, clearAuth, authHeaders } = useAuth()
const { currency, setCurrency, getSymbol: currencySymbol, currencySymbols } = useCurrency()

const showAuthModal = ref(null)
const showCurrencyModal = ref(false)
const showTravelerModal = ref(false)
const authEmail = ref('')
const authPassword = ref('')
const authConfirmPassword = ref('')
const authError = ref('')

const profile = ref({})
const profileForm = ref(defaultProfile())
const dietaryText = ref('')
const accessibilityText = ref('')
const profileSaving = ref(false)
const profileMessage = ref('')

const travelers = ref([])
const travelersLoading = ref(false)
const travelerSaving = ref(false)
const travelerMessage = ref('')
const travelerError = ref('')
const travelerForm = ref(defaultTraveler())

const membership = ref({})
const roles = ref([])
const notifications = ref([])
const notificationsLoading = ref(false)
const emailNotifications = ref(true)
const pushNotifications = ref(true)

const API = '/api/v1'
const PREFERENCES_KEY = 'travel_preferences'

const accountInitial = computed(() => (profile.value.display_name || user.value?.email || '?')[0].toUpperCase())
const roleCodes = computed(() => roles.value.map(role => role.code).join(', ') || 'customer')
const primaryTraveler = computed(() => travelers.value.find(traveler => traveler.is_default) || travelers.value[0] || null)
const profileCompletion = computed(() => {
  const fields = [
    profile.value.display_name,
    profile.value.phone,
    profile.value.nationality,
    profile.value.language,
    profile.value.currency,
    profile.value.travel_style,
    profile.value.budget_level,
    profile.value.family_type,
    ...(profile.value.dietary_restrictions || []),
    ...(profile.value.accessibility_needs || []),
  ]
  const completed = fields.filter(Boolean).length
  return Math.min(100, Math.round((completed / fields.length) * 100))
})

function defaultProfile() {
  return {
    display_name: '',
    avatar: '',
    phone: '',
    nationality: '',
    language: 'en',
    currency: 'CNY',
    travel_style: '',
    budget_level: '',
    family_type: '',
    dietary_restrictions: [],
    accessibility_needs: [],
  }
}

function defaultTraveler() {
  return {
    id: 0,
    name: '',
    gender: '',
    birth_date: '',
    document_type: 'PASSPORT',
    document_no: '',
    nationality: '',
    phone: '',
    email: '',
    is_default: false,
  }
}

function splitTags(value) {
  return value.split(',').map(item => item.trim()).filter(Boolean)
}

function joinTags(value) {
  return (value || []).join(', ')
}

function syncProfileForm(data) {
  profile.value = data || {}
  profileForm.value = { ...defaultProfile(), ...data }
  dietaryText.value = joinTags(data?.dietary_restrictions)
  accessibilityText.value = joinTags(data?.accessibility_needs)
  if (['en', 'zh'].includes(data?.language)) locale.value = data.language
  if (data?.currency) setCurrency(data.currency)
}

function toggleLang() {
  locale.value = locale.value === 'en' ? 'zh' : 'en'
  profileForm.value.language = locale.value
}

function loadPreferences() {
  try {
    const raw = localStorage.getItem(PREFERENCES_KEY)
    if (!raw) return
    const saved = JSON.parse(raw)
    if (['en', 'zh'].includes(saved.locale)) locale.value = saved.locale
    if (saved.currency) setCurrency(saved.currency)
    emailNotifications.value = saved.emailNotifications ?? true
    pushNotifications.value = saved.pushNotifications ?? true
  } catch (_) {}
}

function savePreferences() {
  localStorage.setItem(PREFERENCES_KEY, JSON.stringify({
    locale: locale.value,
    currency: currency.value,
    emailNotifications: emailNotifications.value,
    pushNotifications: pushNotifications.value,
  }))
}

async function loadUserContext() {
  if (!isLoggedIn.value) return
  travelersLoading.value = true
  try {
    const headers = authHeaders()
    const [profileData, travelerData, membershipData, roleData] = await Promise.all([
      fetchProfile(headers),
      fetchTravelers(headers),
      fetchMembership(headers),
      fetchRoles(headers),
    ])
    syncProfileForm(profileData)
    travelers.value = travelerData
    membership.value = membershipData
    roles.value = roleData
  } catch (e) {
    travelerError.value = t('userContext.userLoadFailed')
  } finally {
    travelersLoading.value = false
  }
}

async function saveProfile() {
  profileSaving.value = true
  profileMessage.value = ''
  try {
    if (['en', 'zh'].includes(profileForm.value.language)) locale.value = profileForm.value.language
    if (profileForm.value.currency) setCurrency(profileForm.value.currency)
    const saved = await updateProfile({
      ...profileForm.value,
      language: profileForm.value.language,
      currency: profileForm.value.currency,
      dietary_restrictions: splitTags(dietaryText.value),
      accessibility_needs: splitTags(accessibilityText.value),
    }, authHeaders())
    syncProfileForm(saved)
    profileMessage.value = t('userContext.profileSaved')
  } catch (e) {
    profileMessage.value = t('userContext.profileSaveFailed')
  } finally {
    profileSaving.value = false
  }
}

function openTravelerModal(traveler = null) {
  travelerError.value = ''
  travelerForm.value = traveler
    ? { ...defaultTraveler(), ...traveler, document_no: '' }
    : defaultTraveler()
  showTravelerModal.value = true
}

function closeTravelerModal() {
  showTravelerModal.value = false
  travelerForm.value = defaultTraveler()
}

async function saveTraveler() {
  travelerSaving.value = true
  travelerError.value = ''
  travelerMessage.value = ''
  try {
    const payload = { ...travelerForm.value }
    const saved = payload.id
      ? await updateTraveler(payload.id, payload, authHeaders())
      : await createTraveler(payload, authHeaders())
    const index = travelers.value.findIndex(item => item.id === saved.id)
    if (index >= 0) travelers.value[index] = saved
    else travelers.value.unshift(saved)
    if (saved.is_default) {
      travelers.value = travelers.value.map(item => ({ ...item, is_default: item.id === saved.id }))
    }
    travelerMessage.value = t('userContext.travelerSaved')
    closeTravelerModal()
  } catch (e) {
    travelerError.value = e.message === 'document_duplicate' ? t('userContext.travelerDuplicate') : t('userContext.travelerSaveFailed')
  } finally {
    travelerSaving.value = false
  }
}

async function makeDefaultTraveler(traveler) {
  try {
    const saved = await setDefaultTraveler(traveler.id, authHeaders())
    travelers.value = travelers.value.map(item => ({ ...item, is_default: item.id === saved.id }))
    travelerMessage.value = t('userContext.defaultUpdated')
  } catch (e) {
    travelerError.value = t('userContext.defaultUpdateFailed')
  }
}

async function removeTraveler(traveler) {
  try {
    await deleteTraveler(traveler.id, authHeaders())
    travelers.value = travelers.value.filter(item => item.id !== traveler.id)
    travelerMessage.value = t('userContext.travelerDeleted')
  } catch (e) {
    travelerError.value = t('userContext.travelerDeleteFailed')
  }
}

function selectCurrency(code) {
  setCurrency(code)
  profileForm.value.currency = code
  showCurrencyModal.value = false
}

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
    loadUserContext()
    fetchNotifications()
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
      authError.value = data.error || t('auth.registrationFailed')
      return
    }
    showAuthModal.value = 'login'
  } catch (e) {
    authError.value = t('auth.networkError')
  }
}

async function doLogout() {
  await fetch(API + '/auth/logout', { method: 'POST', headers: authHeaders() })
  clearAuth()
  router.push('/')
}

async function fetchNotifications() {
  if (!isLoggedIn.value) return
  notificationsLoading.value = true
  try {
    const res = await fetch(API + '/notifications', { headers: authHeaders() })
    if (res.ok) {
      const data = await res.json()
      notifications.value = data.notifications || []
    }
  } catch (e) {
    console.error(e)
  } finally {
    notificationsLoading.value = false
  }
}

async function markNotificationAsRead(notification) {
  if (!notification || notification.read || !isLoggedIn.value) return
  notification.read = true
  try {
    await fetch(API + '/notifications', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...authHeaders() },
      body: JSON.stringify({ notification_id: notification.id }),
    })
  } catch (e) {
    notification.read = false
  }
}

loadPreferences()

watch(isLoggedIn, (newVal) => {
  if (newVal) {
    loadUserContext()
    fetchNotifications()
  } else {
    notifications.value = []
    travelers.value = []
    profile.value = {}
  }
}, { immediate: true })

watch([locale, currency, emailNotifications, pushNotifications], savePreferences)
</script>

<style scoped>
.account-page {
  min-height: 100vh;
  background: var(--bg);
}

.account-content {
  padding: 120px 40px 40px;
  max-width: 1120px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 28px;
}

.page-header h1 {
  font-size: 2.25rem;
  margin: 6px 0 8px;
}

.page-header p,
.card-head p {
  color: var(--text-muted);
  margin: 0;
}

.auth-prompt,
.account-card {
  background: var(--surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--surface-border);
  box-shadow: var(--shadow-sm);
}

.auth-prompt {
  text-align: center;
  padding: 60px 20px;
}

.auth-btn,
.mini-btn {
  background: var(--primary);
  color: #fff;
  border: none;
  padding: 12px 28px;
  border-radius: 999px;
  font-weight: 800;
  cursor: pointer;
}

.account-card {
  padding: 24px;
  margin-bottom: 24px;
}

.account-card h2 {
  font-size: 1.2rem;
  margin: 0;
}

.account-hero {
  display: grid;
  grid-template-columns: minmax(0, 1.35fr) minmax(280px, 0.65fr);
  align-items: stretch;
  gap: 18px;
  background: linear-gradient(135deg, rgba(255, 56, 92, 0.12), rgba(0, 102, 204, 0.1)), var(--surface);
}

.hero-main {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 18px;
}

.hero-user {
  display: flex;
  align-items: center;
  gap: 18px;
}

.hero-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.hero-badges span,
.completion-pill,
.status-dot {
  border-radius: 999px;
  font-weight: 900;
}

.hero-badges span {
  background: rgba(0, 102, 204, 0.1);
  color: var(--secondary);
  padding: 5px 10px;
  font-size: 0.8rem;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.primary-cta,
.secondary-cta {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 12px 18px;
  border-radius: 999px;
  font-weight: 900;
  text-decoration: none;
  border: 1px solid transparent;
  cursor: pointer;
}

.primary-cta {
  color: #fff;
  background: linear-gradient(135deg, var(--primary), var(--secondary));
}

.secondary-cta {
  background: rgba(255, 255, 255, 0.82);
  border-color: var(--surface-border);
  color: var(--text);
}

.readiness-card {
  border-radius: 24px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.75);
  border: 1px solid var(--surface-border);
}

.readiness-card h3 {
  margin: 6px 0 16px;
  font-size: 1.05rem;
}

.readiness-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.readiness-item {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.readiness-item strong {
  display: block;
  margin-bottom: 4px;
}

.readiness-item small {
  color: var(--text-muted);
}

.status-dot {
  width: 12px;
  height: 12px;
  margin-top: 5px;
  background: rgba(255, 56, 92, 0.16);
  flex-shrink: 0;
}

.status-dot.active {
  background: var(--primary);
}

.quick-actions {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  margin: 18px 0 24px;
}

.quick-action-card {
  border-radius: 20px;
  border: 1px solid var(--surface-border);
  background: var(--surface);
  box-shadow: var(--shadow-sm);
  padding: 18px;
  text-decoration: none;
  color: var(--text);
  display: flex;
  flex-direction: column;
  gap: 8px;
  text-align: left;
  cursor: pointer;
}

.quick-icon {
  font-size: 1.3rem;
}

.quick-action-card small {
  color: var(--text-muted);
}

.user-avatar-large {
  width: 72px;
  height: 72px;
  background: linear-gradient(135deg, var(--primary), var(--secondary));
  color: #fff;
  border-radius: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.7rem;
  font-weight: 950;
}

.hero-user h2 {
  font-size: 1.8rem;
  margin: 4px 0;
}

.hero-user p {
  color: var(--text-muted);
  margin: 0;
}

.hero-stats {
  display: grid;
  grid-template-columns: repeat(4, minmax(90px, 1fr));
  gap: 12px;
}

.compact-stats {
  margin-top: 16px;
}

.hero-stats div,
.membership-panel,
.traveler-card {
  background: rgba(255, 255, 255, 0.76);
  border: 1px solid var(--surface-border);
  border-radius: 18px;
  padding: 16px;
}

.hero-stats small {
  color: var(--text-muted);
  display: block;
  margin-bottom: 6px;
}

.hero-stats strong {
  font-size: 1.15rem;
}

.account-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.25fr) minmax(320px, 0.75fr);
  gap: 24px;
}

.card-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 18px;
}

.profile-form {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.profile-card {
  scroll-margin-top: 120px;
}

.profile-summary {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 12px;
}

.profile-summary h2 {
  margin-bottom: 6px;
}

.profile-summary p {
  margin: 0;
  color: var(--text-muted);
}

.completion-pill {
  background: rgba(255, 56, 92, 0.12);
  color: var(--primary);
  padding: 8px 12px;
}

.profile-details {
  margin-top: 12px;
}

.profile-details > summary {
  cursor: pointer;
  font-weight: 900;
  color: var(--secondary);
  margin-bottom: 16px;
}

.profile-details[open] > summary {
  margin-bottom: 18px;
}

.auth-input {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 12px;
  font: inherit;
}

.auth-submit {
  padding: 12px 18px;
  background: linear-gradient(135deg, var(--primary), var(--secondary));
  color: #fff;
  border: none;
  border-radius: 12px;
  font-weight: 900;
  cursor: pointer;
}

.add-traveler-btn {
  width: auto;
  white-space: nowrap;
}

.membership-panel {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.points-ring {
  width: 74px;
  height: 74px;
  border-radius: 50%;
  background: rgba(255, 56, 92, 0.12);
  color: var(--primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 950;
}

.tag-row {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.tag-row span,
.traveler-title span {
  background: rgba(0, 102, 204, 0.1);
  color: var(--secondary);
  border-radius: 999px;
  padding: 4px 10px;
  font-size: 0.78rem;
  font-weight: 900;
}

.settings-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.settings-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px;
  border-radius: 14px;
  background: var(--bg-soft);
  cursor: pointer;
}

.settings-item .arrow {
  color: var(--text-muted);
}

.traveler-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 14px;
}

.traveler-empty {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 18px;
  border-radius: 18px;
  background: var(--bg-soft);
}

.traveler-empty p {
  margin: 6px 0 0;
  color: var(--text-muted);
}

.traveler-card {
  display: flex;
  justify-content: space-between;
  gap: 14px;
}

.traveler-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.traveler-card p {
  margin: 0 0 6px;
  color: var(--text);
}

.traveler-card small {
  color: var(--text-muted);
}

.traveler-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.traveler-actions button {
  border: 1px solid var(--surface-border);
  background: #fff;
  border-radius: 999px;
  padding: 7px 10px;
  cursor: pointer;
  font-weight: 800;
}

.traveler-actions .danger-link {
  color: var(--danger);
}

.success-text {
  color: #0f766e;
  font-weight: 800;
  margin: 12px 0 0;
}

.danger-zone {
  border-color: #f5c6cb;
}

.btn-danger {
  background: var(--danger);
  color: #fff;
  border: none;
  padding: 10px 20px;
  border-radius: 12px;
  font-weight: 800;
  cursor: pointer;
  width: 100%;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  padding: 20px;
}

.auth-modal-card {
  background: var(--card);
  border-radius: var(--radius-xl);
  padding: 32px;
  max-width: 430px;
  width: 100%;
  position: relative;
  box-shadow: var(--shadow-lg);
}

.traveler-modal-card {
  max-width: 560px;
}

.modal-close {
  position: absolute;
  top: 16px;
  right: 16px;
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
}

.auth-modal-title {
  margin: 0 0 24px;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.auth-link {
  background: none;
  border: none;
  color: var(--secondary);
  cursor: pointer;
  padding: 8px 0;
}

.auth-error {
  color: var(--danger);
  font-size: 0.9rem;
  margin: 0;
}

.default-check {
  font-weight: 800;
  color: var(--text-muted);
}

.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.notification-item {
  display: flex;
  gap: 12px;
  padding: 16px;
  border-radius: 14px;
  background: var(--bg-soft);
  cursor: pointer;
  border: 1px solid transparent;
}

.notification-item.unread {
  border-color: var(--primary);
}

.notification-icon {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: rgba(255, 56, 92, 0.12);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.notification-content p,
.notification-time {
  margin: 0;
  color: var(--text-muted);
}

.currency-options {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.currency-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 12px 14px;
  border: 1px solid var(--surface-border);
  border-radius: 12px;
  background: var(--bg-soft);
  cursor: pointer;
  font-weight: 800;
}

.currency-option.active {
  border-color: var(--primary);
  background: rgba(255, 56, 92, 0.08);
}

@media (max-width: 860px) {
  .account-content {
    padding: 100px 18px 28px;
  }

  .account-hero,
  .quick-actions,
  .account-grid {
    grid-template-columns: 1fr;
  }

  .hero-user,
  .profile-summary,
  .traveler-empty {
    flex-direction: column;
    align-items: flex-start;
  }

  .hero-stats,
  .profile-form {
    grid-template-columns: 1fr;
  }

  .traveler-card {
    flex-direction: column;
  }
}
</style>
