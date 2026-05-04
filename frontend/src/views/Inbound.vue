<template>
  <div class="inbound-page">
    <SiteHeader />
    <main class="inbound-content">
      <section class="inbound-hero">
        <div>
          <span class="section-kicker">{{ locale === 'zh' ? '中国旅行工具包' : 'China Inbound Kit' }}</span>
          <h1>{{ locale === 'zh' ? '中国入境游工具包' : 'China travel, made arrival-ready' }}</h1>
          <p>{{ locale === 'zh' ? '把 eSIM、接送、高铁、城市通票、支付预约攻略和中文沟通一次准备好。' : 'Prepare eSIM, transfers, high-speed rail, city passes, payment/reservation guides, and Chinese communication in one place.' }}</p>
        </div>
        <router-link class="hero-cta" to="/search?category=City%20Pass">{{ locale === 'zh' ? '查看通票商品' : 'Explore passes' }}</router-link>
      </section>

      <div v-if="loading" class="inbound-card loading-card"><div class="spinner"></div><p>{{ locale === 'zh' ? '加载入境服务...' : 'Loading inbound services...' }}</p></div>

      <template v-else>
        <section class="toolkit-grid">
          <article v-for="item in toolkit" :key="item.id" class="inbound-card toolkit-card">
            <span class="section-kicker">{{ item.category }}</span>
            <h2>{{ locale === 'zh' ? item.title_zh : item.title }}</h2>
            <p>{{ item.description }}</p>
            <ul><li v-for="step in item.steps" :key="step">{{ step }}</li></ul>
            <router-link v-if="item.product_id" class="secondary-btn" :to="`/product/${item.product_id}`">{{ item.cta }}</router-link>
          </article>
        </section>

        <section class="inbound-grid">
          <article class="inbound-card">
            <div class="card-head"><div><span class="section-kicker">{{ locale === 'zh' ? '高铁' : 'Rail' }}</span><h2>{{ locale === 'zh' ? '高铁建议' : 'High-speed rail' }}</h2></div></div>
            <div class="route-list">
              <div v-for="route in rails" :key="route.id" class="route-row">
                <strong>{{ route.from }} → {{ route.to }}</strong>
                <p>{{ route.duration }} · {{ route.frequency }} · ¥{{ route.price_from }}</p>
                <small>{{ route.tip }}</small>
                <router-link v-if="route.product_id" :to="`/product/${route.product_id}`">{{ locale === 'zh' ? '查看服务' : 'View helper' }}</router-link>
              </div>
            </div>
          </article>

          <article class="inbound-card">
            <div class="card-head"><div><span class="section-kicker">{{ locale === 'zh' ? '接送' : 'Transfer' }}</span><h2>{{ locale === 'zh' ? '机场接送' : 'Airport transfers' }}</h2></div></div>
            <div class="route-list">
              <div v-for="transfer in transfers" :key="transfer.id" class="route-row">
                <strong>{{ transfer.city }} · {{ transfer.from }} → {{ transfer.to }}</strong>
                <p>{{ transfer.vehicle }} · ¥{{ transfer.price_from }}</p>
                <small>{{ transfer.driver_tip }}</small>
                <router-link v-if="transfer.product_id" :to="`/product/${transfer.product_id}`">{{ locale === 'zh' ? '预订接送' : 'Book transfer' }}</router-link>
              </div>
            </div>
          </article>
        </section>

        <section class="inbound-card">
          <div class="card-head"><div><span class="section-kicker">{{ locale === 'zh' ? '城市通票' : 'City Pass' }}</span><h2>{{ locale === 'zh' ? '城市通票' : 'City passes' }}</h2></div></div>
          <div class="pass-grid">
            <article v-for="pass in passes" :key="pass.id" class="pass-card">
              <strong>{{ pass.name }}</strong>
              <p>{{ pass.city }} · {{ pass.duration }} · ¥{{ pass.price_from }}</p>
              <div class="pass-tags"><span v-for="included in pass.includes" :key="included">{{ included }}</span></div>
              <router-link v-if="pass.product_id" class="primary-btn" :to="`/product/${pass.product_id}`">{{ locale === 'zh' ? '购买通票' : 'Book pass' }}</router-link>
            </article>
          </div>
        </section>

        <section class="inbound-grid">
          <article class="inbound-card concierge-card">
            <div><span class="section-kicker">{{ locale === 'zh' ? 'AI 管家' : 'AI Concierge' }}</span><h2>{{ locale === 'zh' ? '跨语言旅行管家' : 'Cross-language concierge' }}</h2></div>
            <form class="concierge-form" @submit.prevent="askConcierge">
              <input v-model="conciergePrompt" class="auth-input" :placeholder="locale === 'zh' ? '帮我生成给司机的中文，并安排杭州2天低预算' : 'Generate Chinese for driver and plan Hangzhou 2 days low budget'" />
              <button class="primary-btn" type="submit" :disabled="conciergeLoading">{{ conciergeLoading ? (locale === 'zh' ? '生成中...' : 'Generating...') : (locale === 'zh' ? '生成建议' : 'Generate') }}</button>
            </form>
            <div v-if="conciergeAnswer" class="answer-box">
              <strong>{{ conciergeAnswer.city }} · {{ conciergeAnswer.summary }}</strong>
              <p>🇨🇳 {{ conciergeAnswer.chinese_message }}</p>
              <p>☁️ {{ conciergeAnswer.weather_adjustment }}</p>
              <p>💰 {{ conciergeAnswer.budget_suggestion }}</p>
              <p>🚄 {{ conciergeAnswer.transport_plan }}</p>
              <div class="pass-tags"><span v-for="item in conciergeAnswer.practical_checklist" :key="item">{{ item }}</span></div>
            </div>
          </article>

          <article class="inbound-card">
            <div><span class="section-kicker">{{ locale === 'zh' ? '攻略' : 'Guides' }}</span><h2>{{ locale === 'zh' ? '海外游客城市攻略' : 'City practical guides' }}</h2></div>
            <select v-model="selectedCity" class="auth-input"><option v-for="guide in guides" :key="guide.city" :value="guide.city">{{ guide.city }}</option></select>
            <div v-if="selectedGuide" class="guide-detail">
              <p><strong>{{ locale === 'zh' ? '最佳季节：' : 'Best season: ' }}</strong>{{ selectedGuide.best_season }}</p>
              <p><strong>{{ locale === 'zh' ? '交通：' : 'Transport: ' }}</strong>{{ selectedGuide.transport }}</p>
              <p><strong>{{ locale === 'zh' ? '支付：' : 'Payment: ' }}</strong>{{ selectedGuide.payment }}</p>
              <p><strong>{{ locale === 'zh' ? '网络：' : 'Connectivity: ' }}</strong>{{ selectedGuide.connectivity }}</p>
              <p><strong>{{ locale === 'zh' ? '预约：' : 'Reservation: ' }}</strong>{{ selectedGuide.reservation }}</p>
              <div class="pass-tags"><span v-for="tip in selectedGuide.language_tips" :key="tip">{{ tip }}</span></div>
            </div>
          </article>
        </section>
      </template>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import SiteHeader from '../components/SiteHeader.vue'
import { askInboundConcierge, fetchInboundSnapshot } from '../composables/useInbound'

const { locale } = useI18n()
const loading = ref(true)
const toolkit = ref([])
const rails = ref([])
const transfers = ref([])
const passes = ref([])
const guides = ref([])
const selectedCity = ref('Hangzhou')
const conciergePrompt = ref('Need eSIM, driver Chinese phrase and Hangzhou 2 days low budget')
const conciergeLoading = ref(false)
const conciergeAnswer = ref(null)
const selectedGuide = computed(() => guides.value.find(item => item.city === selectedCity.value) || guides.value[0])

async function loadInbound() {
  loading.value = true
  try {
    const data = await fetchInboundSnapshot()
    toolkit.value = data.toolkit || []
    rails.value = data.rails || []
    transfers.value = data.transfers || []
    passes.value = data.passes || []
    guides.value = data.guides || []
    selectedCity.value = guides.value[0]?.city || 'Hangzhou'
  } finally {
    loading.value = false
  }
}

async function askConcierge() {
  const prompt = conciergePrompt.value.trim()
  if (!prompt) return
  conciergeLoading.value = true
  try {
    conciergeAnswer.value = await askInboundConcierge({ prompt, city: selectedCity.value, budget: 1000, days: 2 })
  } finally {
    conciergeLoading.value = false
  }
}

onMounted(loadInbound)
</script>

<style scoped>
.inbound-page { min-height: 100vh; background: var(--bg); }
.inbound-content { max-width: 1180px; margin: 0 auto; padding: 118px 28px 44px; }
.inbound-hero, .inbound-card, .pass-card { background: var(--surface); border: 1px solid var(--surface-border); border-radius: var(--radius-lg); box-shadow: var(--shadow-sm); }
.inbound-hero { display: flex; justify-content: space-between; gap: 20px; align-items: center; padding: 30px; margin-bottom: 22px; background: linear-gradient(135deg, rgba(255,56,92,.09), rgba(0,122,255,.09)), var(--surface); }
.inbound-hero h1 { margin: 6px 0; font-size: 2.35rem; }
.inbound-hero p, .toolkit-card p, .route-row p, .guide-detail p, .answer-box p { color: var(--text-muted); }
.section-kicker { color: var(--primary); font-size: .76rem; font-weight: 950; letter-spacing: .08em; text-transform: uppercase; }
.hero-cta, .primary-btn { display: inline-flex; border: none; border-radius: 12px; background: linear-gradient(135deg, var(--primary), var(--primary-dark)); color: #fff; font-weight: 950; padding: 12px 18px; text-decoration: none; cursor: pointer; }
.secondary-btn { display: inline-flex; width: fit-content; border: 1px solid var(--surface-border); border-radius: 10px; background: #fff; color: var(--text); font-weight: 850; padding: 10px 14px; text-decoration: none; }
.toolkit-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 16px; margin-bottom: 18px; }
.inbound-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 18px; margin-bottom: 18px; }
.inbound-card { padding: 22px; margin-bottom: 18px; }
.card-head { display: flex; justify-content: space-between; gap: 12px; align-items: center; margin-bottom: 16px; }
.card-head h2, .toolkit-card h2 { margin: 4px 0; }
ul { padding-left: 18px; color: var(--text-muted); }
.route-list, .concierge-form, .guide-detail { display: grid; gap: 12px; }
.route-row { padding: 13px; border: 1px solid var(--surface-border); border-radius: 14px; background: var(--bg-soft); }
.route-row a { color: var(--primary); font-weight: 900; text-decoration: none; }
.route-row p, .route-row small { display: block; margin: 5px 0; }
.pass-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 14px; }
.pass-card { padding: 18px; display: grid; gap: 10px; }
.pass-tags { display: flex; flex-wrap: wrap; gap: 8px; }
.pass-tags span { padding: 6px 9px; border-radius: 999px; background: var(--accent-soft); color: var(--primary); font-size: .78rem; font-weight: 900; }
.auth-input { width: 100%; padding: 12px; border: 1px solid #ddd; border-radius: 10px; font: inherit; }
.answer-box { display: grid; gap: 8px; margin-top: 14px; padding: 14px; border-radius: 14px; background: var(--bg-soft); }
.loading-card { text-align: center; padding: 50px 20px; }
.spinner { width: 34px; height: 34px; border: 3px solid var(--bg-soft); border-top-color: var(--primary); border-radius: 50%; animation: spin 1s linear infinite; margin: 0 auto 14px; }
@keyframes spin { to { transform: rotate(360deg); } }
@media (max-width: 980px) { .toolkit-grid, .inbound-grid, .pass-grid, .inbound-hero { grid-template-columns: 1fr; display: grid; } }
</style>
