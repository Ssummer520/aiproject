import { computed, ref } from 'vue'

const interestKeywords = {
  nature: ['nature', 'lake', 'mountain', 'hiking', 'park', 'forest', 'outdoor', '自然', '山', '湖', '徒步'],
  culture: ['culture', 'history', 'museum', 'heritage', 'temple', 'archaeology', '古', '文化', '历史', '博物馆', '寺'],
  family: ['theme', 'disney', 'resort', 'family', 'park', 'kids', '亲子', '乐园', '迪士尼'],
  food: ['food', 'tea', 'snack', 'restaurant', 'dining', '美食', '茶', '小吃'],
  city: ['city', 'night', 'skyline', 'bund', 'shopping', '城市', '夜景', '外滩'],
}

const promptMeta = {
  en: {
    greeting: 'Hi! Tell me your vibe, budget, city, or trip length — I can shortlist places and explain why they fit.',
    fallback: 'I found a few balanced picks that match popular traveler preferences.',
    empty: 'I do not have enough destination data yet. Try refreshing the page or searching by city.',
    quickPrompts: [
      'Plan a 2-day Hangzhou trip',
      'Best family-friendly places',
      'Nature and hiking ideas',
      'Low budget weekend picks',
    ],
    recommendPrefix: 'Here are my best matches:',
    cityPrefix: 'City match',
    budgetPrefix: 'Budget fit',
    interestPrefix: 'Interest fit',
    highlyRated: 'high rating',
    popular: 'popular with travelers',
  },
  zh: {
    greeting: '你好！告诉我你的城市、预算、出行天数或偏好，我可以帮你快速筛选目的地并说明理由。',
    fallback: '我先按热门度、评分和旅行体验平衡推荐这几处。',
    empty: '当前目的地数据还不够，请刷新页面或先按城市搜索。',
    quickPrompts: [
      '帮我规划杭州2日游',
      '推荐适合亲子的地方',
      '想看自然风光和徒步',
      '低预算周末去哪玩',
    ],
    recommendPrefix: '我推荐这些目的地：',
    cityPrefix: '城市匹配',
    budgetPrefix: '预算合适',
    interestPrefix: '兴趣匹配',
    highlyRated: '评分高',
    popular: '近期热门',
  },
}

function normalize(value) {
  return String(value || '').trim().toLowerCase()
}

function uniqueDestinations(items) {
  const seen = new Set()
  const list = []
  for (const item of items || []) {
    if (!item?.id || seen.has(item.id)) continue
    seen.add(item.id)
    list.push(item)
  }
  return list
}

function getMeta(lang) {
  return promptMeta[lang === 'zh' ? 'zh' : 'en']
}

function detectBudget(text) {
  const lower = normalize(text)
  const numberMatches = lower.match(/\d+/g) || []
  const maxBudget = numberMatches.length ? Math.max(...numberMatches.map(Number)) : 0

  if (maxBudget > 0) return maxBudget
  if (/cheap|budget|low|affordable|省钱|便宜|低预算|穷游/.test(lower)) return 180
  if (/premium|luxury|high end|高端|豪华|品质/.test(lower)) return 9999
  return 0
}

function detectTripLength(text) {
  const lower = normalize(text)
  const match = lower.match(/(\d+)\s*(day|days|d|天|日)/)
  if (match) return Number(match[1])
  if (/weekend|周末/.test(lower)) return 2
  return 0
}

function detectInterests(text) {
  const lower = normalize(text)
  return Object.entries(interestKeywords)
    .filter(([, words]) => words.some((word) => lower.includes(word)))
    .map(([key]) => key)
}

function destinationText(destination) {
  return [
    destination?.name,
    destination?.city,
    destination?.description,
    destination?.policy,
    ...(destination?.tags || []),
    ...(destination?.amenities || []),
  ].join(' ').toLowerCase()
}

function scoreDestination(destination, text, lang) {
  const query = normalize(text)
  const haystack = destinationText(destination)
  const reasons = []
  let score = 0

  if (!destination?.id) return { score, reasons }

  const city = normalize(destination.city)
  if (city && query.includes(city)) {
    score += 5
    reasons.push(getMeta(lang).cityPrefix)
  }

  const tokens = query.split(/[\s,，。.!?;；、]+/).filter((word) => word.length > 1)
  for (const token of tokens) {
    if (haystack.includes(token)) score += 1.4
  }

  const budget = detectBudget(text)
  if (budget > 0 && Number(destination.price || 0) <= budget) {
    score += 2.6
    reasons.push(getMeta(lang).budgetPrefix)
  }

  for (const interest of detectInterests(text)) {
    if ((interestKeywords[interest] || []).some((word) => haystack.includes(word))) {
      score += 3
      reasons.push(getMeta(lang).interestPrefix)
    }
  }

  if (Number(destination.rating || 0) >= 4.8) {
    score += 1.2
    reasons.push(getMeta(lang).highlyRated)
  }

  if (Number(destination.booked_count || 0) >= 30 || Number(destination.review_count || 0) >= 3000) {
    score += 0.8
    reasons.push(getMeta(lang).popular)
  }

  const tripLength = detectTripLength(text)
  if (tripLength >= 2 && (destination.tags || []).length > 1) score += 0.5

  return { score, reasons: [...new Set(reasons)].slice(0, 3) }
}

function buildReply(text, destinations, lang) {
  const meta = getMeta(lang)
  const candidates = uniqueDestinations(destinations)

  if (!candidates.length) {
    return { text: meta.empty, destinations: [] }
  }

  const ranked = candidates
    .map((destination) => ({ destination, ...scoreDestination(destination, text, lang) }))
    .sort((a, b) => b.score - a.score || (b.destination.rating || 0) - (a.destination.rating || 0))

  const top = ranked.filter((item) => item.score > 0).slice(0, 3)
  const fallback = ranked.slice(0, 3)
  const selected = top.length ? top : fallback

  const intro = top.length ? meta.recommendPrefix : meta.fallback
  const detail = selected
    .map(({ destination, reasons }, index) => {
      const reasonText = reasons.length ? reasons.join(' · ') : `${destination.city || ''} · ${destination.rating || ''}★`
      return `${index + 1}. ${destination.name} — ${reasonText}`
    })
    .join('\n')

  return {
    text: `${intro}\n${detail}`,
    destinations: selected.map((item) => item.destination),
  }
}

export function useTravelAssistant(options = {}) {
  const locale = options.locale
  const sourceDestinations = options.destinations
  const messages = ref([])

  const language = computed(() => locale?.value === 'zh' ? 'zh' : 'en')
  const meta = computed(() => getMeta(language.value))
  const quickPrompts = computed(() => meta.value.quickPrompts)

  function resetGreeting() {
    messages.value = [
      {
        id: `assistant-${Date.now()}`,
        role: 'assistant',
        text: meta.value.greeting,
        destinations: [],
      },
    ]
  }

  function ask(text) {
    const trimmed = String(text || '').trim()
    if (!trimmed) return null

    const userMessage = {
      id: `user-${Date.now()}-${messages.value.length}`,
      role: 'user',
      text: trimmed,
      destinations: [],
    }
    const answer = buildReply(trimmed, sourceDestinations?.value || [], language.value)
    const assistantMessage = {
      id: `assistant-${Date.now()}-${messages.value.length}`,
      role: 'assistant',
      text: answer.text,
      destinations: answer.destinations,
    }

    messages.value = [...messages.value, userMessage, assistantMessage]
    return assistantMessage
  }

  resetGreeting()

  return {
    messages,
    quickPrompts,
    ask,
    resetGreeting,
  }
}
