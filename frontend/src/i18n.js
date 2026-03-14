import { createI18n } from 'vue-i18n'

const messages = {
  en: {
    nav: {
      destinations: 'Destinations',
      experiences: 'Experiences',
      guides: 'Guides',
      myTrips: 'My Trips',
      search: 'Search',
      where: 'Where',
      when: 'When',
      who: 'Who',
      popular: 'Popular'
    },
    hero: {
      title: "Discover China's Best",
      subtitle: 'From Great Wall to West Lake, find your next adventure',
      searchPlaceholder: 'Search destinations'
    },
    deals: {
      title: 'Exclusive Deals',
      endsIn: 'Ends in',
      claimNow: 'Claim Now',
      getCoupon: 'Get Coupon',
      explore: 'Explore'
    },
    recommendations: {
      title: 'Recommended for you',
      locationHint: 'Based on your location'
    },
    nearby: {
      title: 'Explore Nearby {city}',
      viewAll: 'View all destinations',
      away: '{dist} km away'
    },
    trust: {
      securePayment: 'Secure Payment',
      securePaymentDesc: 'All transactions are encrypted',
      support: '24/7 Support',
      supportDesc: "We're here to help anytime",
      verifiedReviews: 'Verified Reviews',
      secureBooking: 'Secure Booking',
      globalSupport: '24/7 Global Support'
    },
    common: {
      night: '/ night',
      reviews: '{count}+ reviews',
      booked: 'Booked {count} times in 24h',
      rareFind: 'Rare find',
      noRecent: 'No recent views',
      recentlyViewed: 'Recently Viewed',
      getApp: 'Get the App',
      getAppDesc: 'Download for exclusive mobile-only deals and offline maps.'
    }
  },
  zh: {
    nav: {
      destinations: '目的地',
      experiences: '体验',
      guides: '攻略',
      myTrips: '我的旅行',
      search: '搜索',
      where: '地点',
      when: '时间',
      who: '人数',
      popular: '热门'
    },
    hero: {
      title: '探索大美中国',
      subtitle: '从长城到西湖，开启你的下一段旅程',
      searchPlaceholder: '搜索目的地'
    },
    deals: {
      title: '特惠活动',
      endsIn: '距结束',
      claimNow: '立即领取',
      getCoupon: '领优惠券',
      explore: '去看看'
    },
    recommendations: {
      title: '为你推荐',
      locationHint: '基于你的位置'
    },
    nearby: {
      title: '{city}周边探索',
      viewAll: '查看全部目的地',
      away: '距离 {dist} km'
    },
    trust: {
      securePayment: '安全支付',
      securePaymentDesc: '所有交易均加密处理',
      support: '24/7 客服',
      supportDesc: '随时为你提供帮助',
      verifiedReviews: '真实评价',
      secureBooking: '安全预订',
      globalSupport: '24/7 全球支持'
    },
    common: {
      night: '/ 晚',
      reviews: '{count}+ 条评价',
      booked: '过去 24 小时预订 {count} 次',
      rareFind: '稀缺好房',
      noRecent: '暂无最近浏览',
      recentlyViewed: '最近浏览',
      getApp: '下载 App',
      getAppDesc: '下载 App 享受移动端专属优惠及离线地图。'
    }
  }
}

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages
})

export default i18n
