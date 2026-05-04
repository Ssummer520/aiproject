import { createI18n } from 'vue-i18n'

const messages = {
  en: {
    nav: {
      destinations: 'Destinations',
      experiences: 'Experiences',
      guides: 'Guides',
      myTrips: 'My Trips',
      wishlist: 'Wishlist',
      history: 'History',
      search: 'Search',
      hangzhouHot: 'Popular Destinations around Hangzhou',
      where: 'Where',
      when: 'When',
      who: 'Who',
      popular: 'Popular',
      inbound: 'Inbound Kit'
    },
    auth: {
      signIn: 'Sign in',
      logOut: 'Log out',
      createAccount: 'Create account',
      register: 'Register',
      forgotPassword: 'Forgot password',
      forgotPasswordQuestion: 'Forgot password?',
      resetPassword: 'Reset password',
      sendResetLink: 'Send reset link',
      backToSignIn: 'Back to Sign in',
      alreadyHaveAccount: 'Already have an account? Sign in',
      email: 'Email',
      password: 'Password',
      passwordMin: 'Password (min 6)',
      confirmPassword: 'Confirm password',
      resetToken: 'Reset token (from email)',
      newPassword: 'New password (min 6)',
      confirmNewPassword: 'Confirm new password',
      passwordsDoNotMatch: 'Passwords do not match',
      emailAlreadyRegistered: 'Email already registered.',
      registrationFailed: 'Registration failed',
      accountCreated: 'Account created. Sign in below.',
      noAccount: 'No account with this email.',
      requestFailed: 'Request failed',
      passwordReset: 'Password reset. Sign in below.',
      authRequiredHistory: 'Sign in to view your history and wishlist.',
      syncHint: 'Sign in to sync wishlist & history'
    },
    ui: {
      map: 'Map',
      switchLanguageCurrency: 'Switch Language/Currency',
      showLess: 'Show Less',
      showMore: 'Show More',
      loading: 'Loading...',
      recent: 'Recent',
      noWishlist: 'No wishlist items',
      topic: 'Topic',
      hot: 'Hot'
    },
    homeContent: {
      teaTag: 'Culture',
      teaTitle: 'The Ancient Art of Tea Making in Hangzhou',
      foodTag: 'Food',
      foodTitle: 'Why Spicy Food is Life in Sichuan',
      himalayaDesc: 'Discovering the hidden valleys and sacred mountains of the Himalayas.',
      hangzhou48Title: '48h in Hangzhou',
      hangzhou48Desc: 'Cycling West Lake & Lingyin Temple',
      familyFunTitle: 'Family Fun List',
      familyFunDesc: 'Top theme parks in Shanghai & beyond',
      streetFoodsTitle: 'Top 10 Street Foods in Chengdu',
      beijingGemsTitle: 'Hidden Gems of Beijing',
      roadTripTag: 'Road Trip',
      tibetTitle: 'Switzerland of the East: A Week in Tibet',
      viewAll: 'View all',
      hotBadge: 'HOT',
      privacy: 'Privacy',
      terms: 'Terms'
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
      title: 'Popular Destinations around Hangzhou',
      viewAll: 'View all',
      away: '{dist}km away'
    },
    hotActivities: {
      title: 'Nearby Hangzhou Hot Activities'
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
      hotActivities: 'Hot City Activities',
      frequentSearch: 'Frequently searched cities',
      inspiration: 'Inspiration',
      categories: 'Explore by Category',
      travelGuide: 'Travel Guides',
      travelBlog: 'Travel Stories & Tips',
      readMore: 'Read Story'
    },
    activities: {
      shanghai: {
        title: 'Shanghai · Night Cruise & Skyline',
        desc: 'Slowly sail along the Huangpu River, with skyscrapers as your backdrop, and let the city lights tell the story of the day.'
      },
      beijing: {
        title: 'Beijing · Morning on the City Walls',
        desc: 'Ascend the city tower before the crowds arrive, and watch the first light of dawn touch the eaves and bricks, as the city suddenly falls quiet.'
      },
      chengdu: {
        title: 'Chengdu · Midnight in a Bowl of Noodles',
        desc: 'Start from a small shop at the alley entrance, follow the aroma through the entire night market, and let the local life be the destination of your journey.'
      }
    }
  },
  zh: {
    nav: {
      destinations: '目的地',
      experiences: '体验',
      guides: '攻略',
      myTrips: '我的旅行',
      wishlist: '收藏夹',
      history: '最近浏览',
      search: '搜索',
      hangzhouHot: '杭州周边的热门目的地',
      where: '地点',
      when: '时间',
      who: '人数',
      popular: '热门',
      inbound: '入境工具包'
    },
    auth: {
      signIn: '登录',
      logOut: '退出登录',
      createAccount: '创建账号',
      register: '注册',
      forgotPassword: '忘记密码',
      forgotPasswordQuestion: '忘记密码？',
      resetPassword: '重置密码',
      sendResetLink: '发送重置链接',
      backToSignIn: '返回登录',
      alreadyHaveAccount: '已有账号？去登录',
      email: '邮箱',
      password: '密码',
      passwordMin: '密码（至少 6 位）',
      confirmPassword: '确认密码',
      resetToken: '重置令牌（来自邮件）',
      newPassword: '新密码（至少 6 位）',
      confirmNewPassword: '确认新密码',
      passwordsDoNotMatch: '两次输入的密码不一致',
      emailAlreadyRegistered: '该邮箱已注册。',
      registrationFailed: '注册失败',
      accountCreated: '账号已创建，请在下方登录。',
      noAccount: '没有找到该邮箱对应的账号。',
      requestFailed: '请求失败',
      passwordReset: '密码已重置，请在下方登录。',
      authRequiredHistory: '登录后可查看浏览记录与收藏。',
      syncHint: '登录后可同步收藏与浏览记录'
    },
    ui: {
      map: '地图',
      switchLanguageCurrency: '切换语言/货币',
      showLess: '收起',
      showMore: '查看更多',
      loading: '加载中...',
      recent: '最近浏览',
      noWishlist: '暂无收藏',
      topic: '专题',
      hot: '热门'
    },
    homeContent: {
      teaTag: '文化',
      teaTitle: '杭州古法制茶之美',
      foodTag: '美食',
      foodTitle: '为什么川味麻辣让人上头',
      himalayaDesc: '探索喜马拉雅隐秘山谷与神圣雪山。',
      hangzhou48Title: '杭州 48 小时',
      hangzhou48Desc: '骑行西湖，漫游灵隐寺',
      familyFunTitle: '亲子玩乐清单',
      familyFunDesc: '上海及周边主题乐园精选',
      streetFoodsTitle: '成都十大街头美食',
      beijingGemsTitle: '北京隐藏玩法',
      roadTripTag: '自驾',
      tibetTitle: '东方瑞士：西藏一周旅行',
      viewAll: '查看全部',
      hotBadge: '热门',
      privacy: '隐私',
      terms: '条款'
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
      title: '杭州周边的热门目的地',
      viewAll: '查看全部',
      away: '距离 {dist}km'
    },
    hotActivities: {
      title: '杭州周边热门活动'
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
      hotActivities: '热门城市活动',
      frequentSearch: '这些城市正在被频繁搜索',
      inspiration: '灵感推荐',
      categories: '按分类探索',
      travelGuide: '旅行指南',
      travelBlog: '旅行故事与灵感',
      readMore: '阅读全文'
    },
    activities: {
      shanghai: {
        title: '上海 · 夜色游船与天际线',
        desc: '沿着黄浦江缓慢驶过，把摩天大楼当作背景板，让城市灯光替你讲完这一天。'
      },
      beijing: {
        title: '北京 · 城墙之上的清晨',
        desc: '在人潮之前登上城楼，看第一缕阳光落在屋檐和城砖上，城市忽然安静下来。'
      },
      chengdu: {
        title: '成都 · 一碗面里的深夜',
        desc: '从巷子口的小馆开始，顺着香味走完一整条夜市，把烟火气当作行程的终点。'
      }
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
