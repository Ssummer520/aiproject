<template>
  <div class="not-found-page">
    <div class="not-found-content">
      <div class="error-code">404</div>
      <h1>{{ $t('auto.auto_6250ccb0') }}</h1>
      <p>{{ $t('auto.auto_5c2a6acf') }}</p>
      <div class="suggestions">
        <h3>{{ $t('auto.auto_221d0f62') }}</h3>
        <ul>
          <li><router-link to="/">{{ $t('auto.auto_bc1da3ba') }}</router-link></li>
          <li><router-link to="/search">{{ $t('auto.auto_8ec8a53f') }}</router-link></li>
          <li><router-link to="/category/nature">{{ $t('auto.auto_59bd9557') }}</router-link></li>
        </ul>
      </div>
      <router-link to="/" class="back-home-btn">
        {{ $t('auto.auto_9f5b5e10') }}
      </router-link>
    </div>
    <div class="not-found-illustration">
      <div class="floating-element e1">✈️</div>
      <div class="floating-element e2">🗺️</div>
      <div class="floating-element e3">🌍</div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
const { locale } = useI18n()
</script>

<style scoped>
.not-found-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 40px 20px;
}

.not-found-content {
  background: #fff;
  border-radius: 24px;
  padding: 60px;
  max-width: 500px;
  text-align: center;
  box-shadow: 0 20px 60px rgba(0,0,0,0.3);
}

.error-code {
  font-size: 120px;
  font-weight: 800;
  color: #667eea;
  line-height: 1;
  margin-bottom: 20px;
}

h1 {
  font-size: 28px;
  margin-bottom: 16px;
  color: #333;
}

p {
  color: #666;
  margin-bottom: 30px;
}

.suggestions {
  text-align: left;
  background: #f8f9fa;
  padding: 20px;
  border-radius: 12px;
  margin-bottom: 30px;
}

.suggestions h3 {
  font-size: 16px;
  margin-bottom: 12px;
  color: #333;
}

.suggestions ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.suggestions li {
  margin-bottom: 8px;
}

.suggestions a {
  color: #667eea;
  text-decoration: none;
}

.suggestions a:hover {
  text-decoration: underline;
}

.back-home-btn {
  display: inline-block;
  background: #667eea;
  color: #fff;
  padding: 14px 32px;
  border-radius: 12px;
  text-decoration: none;
  font-weight: 600;
  transition: all 0.2s;
}

.back-home-btn:hover {
  background: #5a6fd6;
  transform: translateY(-2px);
}

.not-found-illustration {
  position: fixed;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  pointer-events: none;
  overflow: hidden;
}

.floating-element {
  position: absolute;
  font-size: 60px;
  animation: float 6s ease-in-out infinite;
}

.e1 { top: 10%; left: 10%; animation-delay: 0s; }
.e2 { top: 20%; right: 15%; animation-delay: 2s; }
.e3 { bottom: 15%; left: 20%; animation-delay: 4s; }

@keyframes float {
  0%, 100% { transform: translateY(0) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(10deg); }
}
</style>
