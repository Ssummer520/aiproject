import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const destinationMap = {
  'West Lake': '西湖',
  'The Bund': '外滩',
  'Great Wall': '万里长城',
  'Yellow Mountain': '黄山',
  'Terracotta Army': '兵马俑',
  'Lingyin Temple': '灵隐寺',
  'Shanghai Disney Resort': '上海迪士尼度假区',
  'Wuzhen Water Town': '乌镇水乡',
  'Summer Palace': '颐和园',
  'Forbidden City': '故宫',
  'Chengdu Panda Base': '成都大熊猫基地',
  'Jiuzhaigou Valley': '九寨沟',
  'Zhangjiajie National Forest Park': '张家界国家森林公园',
  'Li River Guilin': '桂林漓江',
  'Potala Palace': '布达拉宫',
  'Suzhou Gardens': '苏州园林',
  'Mount Emei': '峨眉山',
  'Sanya Yalong Bay': '三亚亚龙湾',
  'Dali Old Town': '大理古城',
  'Lijiang Old Town': '丽江古城',
  'Hong Kong Victoria Harbour': '香港维多利亚港',
}

const cityMap = {
  Hangzhou: '杭州',
  Shanghai: '上海',
  Beijing: '北京',
  Huangshan: '黄山',
  "Xi'an": '西安',
  Chengdu: '成都',
  Sichuan: '四川',
  Hunan: '湖南',
  Guilin: '桂林',
  Lhasa: '拉萨',
  Suzhou: '苏州',
  Sanya: '三亚',
  Dali: '大理',
  Lijiang: '丽江',
  'Hong Kong': '香港',
  China: '中国',
}

const termMap = {
  Tickets: '门票', Tours: '一日游', Experiences: '体验', Transport: '交通', Nature: '自然', Culture: '文化',
  'City Pass': '城市通票', Boat: '游船', 'Mobile voucher': '手机凭证', 'English guide': '英文导游', Tea: '茶文化',
  'Theme Park': '主题乐园', Family: '亲子', Bestseller: '热卖', 'Night View': '夜景', Cruise: '游船', Couples: '情侣',
  Museum: '博物馆', History: '历史', Hiking: '徒步', Pickup: '接送', 'Airport transfer': '机场接送', Inbound: '入境游',
  'English support': '英文支持', Beijing: '北京', Hangzhou: '杭州', Shanghai: '上海', Night: '夜游', 'AI itinerary': 'AI 行程',
  'English service': '英文服务', eSIM: 'eSIM', Rail: '高铁', 'Digital pass': '电子通票',
  'Spring Break Deals': '春日大促', 'New User Gift': '新人礼包', 'Weekend Getaway': '周末出逃',
  'Family Fun Pass': '亲子畅玩套票', 'Museum Night Special': '博物馆夜场特惠', 'Foodie Trail Package': '城市觅食路线包',
  'Save 80¥ on bookings over 500¥': '满 500 减 80', '30¥ OFF your first trip in China': '首单立减 30 元',
  'Up to 50% off for local experiences': '本地体验低至 5 折', 'Bundle 2 adult + 1 kid tickets and save 20%': '2 大 1 小组合立减 20%',
  'Late-entry museum packages from 99¥': '夜场联票 99 元起', 'Street food + local guide combo from 129¥': '美食 + 地陪组合 129 元起',
}

const productMap = {
  'West Lake Scenic Boat Ticket': '西湖景区游船票',
  'Classic lake cruise with mobile voucher': '经典湖上游船，支持手机凭证',
  'Cruise across West Lake, pass Su Causeway and Three Pools Mirroring the Moon, and enjoy a flexible entry window designed for first-time visitors.': '泛舟西湖，途经苏堤与三潭印月，灵活入园时段适合首次到访游客。',
  'Lingyin Temple Half-Day Guided Tour': '灵隐寺半日导览游',
  'English guide, tea village stop, small group': '英文导游、茶村停留、小团出行',
  'Explore Lingyin Temple with an English-speaking guide, learn Buddhist culture, and finish with a relaxed Longjing tea village walk.': '跟随英文导游游览灵隐寺，了解佛教文化，并在龙井茶村轻松漫步。',
  'Shanghai Disney Resort 1-Day Ticket': '上海迪士尼度假区一日票',
  'Instant confirmation for family trips': '亲子出行即时确认',
  'Spend a magical day at Shanghai Disney Resort with mobile entry and optional priority add-ons for families.': '在上海迪士尼度假区度过梦幻一天，支持手机入园和亲子优先体验加购。',
  'Huangpu River Night Cruise': '黄浦江夜游船票',
  'Bund skyline, mobile voucher, evening sail': '外滩天际线、手机凭证、夜间航班',
  'See Shanghai from the water with a night cruise along the Huangpu River, passing the Bund and Pudong skyline.': '乘夜游船沿黄浦江观赏上海，从水上经过外滩与浦东天际线。',
  'Forbidden City English Guided Walk': '故宫英文导览步行游',
  'Verified guide with palace highlights': '认证导游讲解宫殿精华',
  'Walk through the Forbidden City with a licensed English-speaking guide and understand imperial stories, architecture, and daily court life.': '跟随持证英文导游步行游览故宫，了解皇家故事、建筑与宫廷日常。',
  'Mutianyu Great Wall Private Transfer': '慕田峪长城私人接送',
  'Hotel pickup, flexible return, optional ticket': '酒店接送、灵活返程、可选门票',
  'Book a private transfer from downtown Beijing to Mutianyu Great Wall with a bilingual driver and flexible waiting time.': '预订北京市区至慕田峪长城私人接送，双语司机与灵活等待时间。',
  'China eSIM 7-Day Data Pack': '中国 eSIM 7 日流量包',
  'QR activation, VPN tips and bilingual setup guide': '二维码激活、网络提示与双语设置指南',
  'Stay connected after landing with a China-compatible eSIM data pack, setup checklist and support note for common travel apps.': '落地后使用兼容中国网络的 eSIM 流量包，附设置清单和常用旅行 App 支持说明。',
  'Shanghai-Hangzhou High-Speed Rail Helper': '上海至杭州高铁助手',
  'Route planning, station guide and passport ticketing tips': '路线规划、车站指南和护照购票提示',
  'A rail planning service that explains station choice, passport ticketing, transfer buffers and arrival tips between Shanghai and Hangzhou.': '高铁规划服务，说明上海与杭州间车站选择、护照购票、换乘预留和到达提示。',
  'Hangzhou Airport to West Lake Transfer': '杭州机场至西湖接送',
  'Private arrival transfer with Chinese hotel confirmation': '私人抵达接送，协助确认中文酒店地址',
  'Licensed car transfer for overseas travellers arriving at Hangzhou Xiaoshan Airport, with English pickup instructions and Chinese address confirmation.': '为抵达杭州萧山机场的游客提供合规车辆接送，含英文接机说明与中文地址确认。',
  'Beijing Airport Downtown Transfer': '北京机场至市区接送',
  'Business van arrival service for long-haul flights': '适合长途航班的商务车抵达服务',
  'Airport-to-hotel transfer with terminal confirmation, driver contact and bilingual arrival message.': '机场至酒店接送，含航站楼确认、司机联系和双语抵达信息。',
  'Hangzhou 2-Day Culture Pass': '杭州 2 日文化通票',
  'West Lake, Lingyin route, transfer coupon and AI plan': '西湖、灵隐路线、接送券和 AI 行程',
  'A China inbound city pass combining scenic experiences, arrival essentials and AI route planning for Hangzhou.': '杭州城市通票，组合景区体验、抵达服务与 AI 路线规划。',
  'Shanghai Night Pass': '上海夜游通票',
  'River cruise, night transfer and food street guide': '江景游船、夜间接送和美食街指南',
  'A night-focused Shanghai pass for overseas travellers with transfer support and bilingual food guidance.': '面向夜游上海的通票，包含接送支持和双语美食指南。',
}

const guideMap = {
  'April-May and September-November': '4-5 月与 9-11 月',
  'Mild and humid; summers can be hot with sudden rain.': '气候温和湿润，夏季较热且可能有阵雨。',
  'Use metro for city center, taxi or private transfer for temples and tea villages.': '市中心建议乘地铁，寺庙和茶村建议打车或预约接送。',
  'Mobile payment is common, but carry some cash for small vendors.': '移动支付很普遍，但小商户建议准备少量现金。',
  'eSIM or roaming is useful; save hotel address offline.': '建议使用 eSIM 或漫游，并离线保存酒店地址。',
  'Popular temples and museums may need passport-based reservation.': '热门寺庙和博物馆可能需要使用护照预约。',
  'April-May and October-November': '4-5 月与 10-11 月',
  'Windy riverside nights; bring a light jacket.': '江边夜晚有风，建议带薄外套。',
  'Metro is efficient; taxis need Chinese destination names.': '地铁高效便捷，打车建议准备中文目的地。',
  'Cards are easier in hotels/malls than small vendors.': '酒店和商场刷卡更方便，小商户不一定支持。',
  'eSIM works well; keep VPN needs checked before travel.': 'eSIM 使用方便，如需 VPN 请出行前确认。',
  'Museums and observation decks may require timed slots.': '博物馆和观景台可能需要预约时段。',
}

export function useLocalization() {
  const { locale } = useI18n()
  const isZh = computed(() => locale.value === 'zh')
  function localizeText(value) {
    if (!isZh.value || value === undefined || value === null) return value || ''
    return productMap[value] || destinationMap[value] || cityMap[value] || termMap[value] || guideMap[value] || value || ''
  }
  function localizeField(item, field) {
    if (!item) return ''
    if (isZh.value) {
      const direct = item[`${field}_zh`] || item[`${field}Zh`]
      if (direct) return direct
    }
    return localizeText(item[field])
  }
  function localizeList(values = []) {
    return (values || []).map(localizeText)
  }
  function localizeDestination(item) { return localizeField(item, 'name') }
  function localizeCity(itemOrValue) { return typeof itemOrValue === 'string' ? localizeText(itemOrValue) : localizeField(itemOrValue, 'city') }
  return { localizeText, localizeField, localizeList, localizeDestination, localizeCity }
}
