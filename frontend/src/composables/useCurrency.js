import { ref } from 'vue'

const currency = ref('CNY')
const currencySymbols = {
  CNY: '¥',
  USD: '$',
  EUR: '€',
  GBP: '£',
  JPY: '¥',
  KRW: '₩',
  THB: '฿',
  SGD: 'S$',
  AUD: 'A$',
  HKD: 'HK$'
}

const exchangeRates = {
  CNY: 1,
  USD: 0.14,
  EUR: 0.13,
  GBP: 0.11,
  JPY: 21,
  KRW: 190,
  THB: 5,
  SGD: 0.19,
  AUD: 0.21,
  HKD: 1.1
}

export function useCurrency() {
  function setCurrency(code) {
    if (currencySymbols[code]) {
      currency.value = code
    }
  }

  function formatPrice(cnyPrice) {
    const rate = exchangeRates[currency.value] || 1
    const converted = cnyPrice * rate
    const symbol = currencySymbols[currency.value] || '¥'

    if (currency.value === 'JPY' || currency.value === 'KRW') {
      return symbol + Math.round(converted).toLocaleString()
    }
    return symbol + converted.toFixed(2)
  }

  function getSymbol() {
    return currencySymbols[currency.value] || '¥'
  }

  return {
    currency,
    setCurrency,
    formatPrice,
    getSymbol,
    currencySymbols
  }
}
