import { ref, computed, onMounted } from 'vue'

const STORAGE_KEY = 'travel_auth'

const token = ref('')
const user = ref(null)

function loadFromStorage() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return
    const data = JSON.parse(raw)
    if (data.token && data.user) {
      token.value = data.token
      user.value = data.user
    }
  } catch (_) {}
}

export function useAuth() {
  // 立即从 localStorage 恢复登录态，避免首屏已登录却仍提示「需要登录」
  if (!token.value && typeof localStorage !== 'undefined') {
    loadFromStorage()
  }

  const isLoggedIn = computed(() => !!token.value)

  function saveToStorage() {
    if (token.value && user.value) {
      localStorage.setItem(STORAGE_KEY, JSON.stringify({ token: token.value, user: user.value }))
    } else {
      localStorage.removeItem(STORAGE_KEY)
    }
  }

  function setAuth(newToken, newUser) {
    token.value = newToken
    user.value = newUser
    saveToStorage()
  }

  function clearAuth() {
    token.value = ''
    user.value = null
    localStorage.removeItem(STORAGE_KEY)
  }

  function authHeaders() {
    const h = {}
    if (token.value) h['Authorization'] = `Bearer ${token.value}`
    return h
  }

  onMounted(() => {
    loadFromStorage()
  })

  return {
    token,
    user,
    isLoggedIn,
    loadFromStorage,
    setAuth,
    clearAuth,
    authHeaders,
  }
}
