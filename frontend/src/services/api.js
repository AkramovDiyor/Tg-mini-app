import axios from 'axios'

// ============================================
// 🤖 ИНИЦИАЛИЗАЦИЯ TELEGRAM WEBAPP
// ============================================
// Говорим Telegram, что WebApp готов к работе
if (window.Telegram?.WebApp) {
  window.Telegram.WebApp.ready()
  // Опционально: расширяем на весь экран
  window.Telegram.WebApp.expand()
  // Включаем вертикальные свайпы для закрытия
  window.Telegram.WebApp.enableClosingConfirmation()
}

// ============================================
// 🕵️ ЛОГИКА ОПРЕДЕЛЕНИЯ РОЛИ
// ============================================

/**
 * Достает start_param из Telegram WebApp или URL
 * - 'master' → панель мастера
 * - 'abc123xy' (invite_link) → клиентская часть мастера
 */
function getStartParam() {
  // 1. Пытаемся достать из реального Telegram WebApp
  if (window.Telegram?.WebApp?.initDataUnsafe?.start_param) {
    return window.Telegram.WebApp.initDataUnsafe.start_param
  }
  
  // 2. Fallback: читаем из URL (для тестов в браузере)
  const urlParams = new URLSearchParams(window.location.search)
  return urlParams.get('startapp') || ''
}

/**
 * Возвращает initData для авторизации
 * - В реальном Telegram: криптографическая подпись от Telegram
 * - В браузере: тестовые токены ('test-master' или 'test-vasya')
 */
function getInitData() {
  // 1. Если мы в настоящем Telegram — используем реальную подпись
  if (window.Telegram?.WebApp?.initData && window.Telegram.WebApp.initData !== '') {
    return window.Telegram.WebApp.initData
  }
  
  // 2. ТЕСТОВЫЙ РЕЖИМ В БРАУЗЕРЕ
  const startParam = getStartParam()
  
  if (startParam === 'master') {
    // Мастерская панель → тестовый мастер (ID: 999999)
    return 'test-master'
  }
  
  // Клиентская часть → тестовый клиент (ID: 777111222)
  return 'test-vasya'
}

/**
 * Определяет, кто мы: мастер или клиент
 */
export function getUserRole() {
  const startParam = getStartParam()
  return startParam === 'master' ? 'master' : 'client'
}

/**
 * Возвращает invite_link текущего мастера (для клиента)
 */
export function getInviteLink() {
  const startParam = getStartParam()
  
  // Если мы мастер — invite_link не нужен
  if (startParam === 'master' || startParam === '') {
    return null
  }
  
  // Иначе startParam — это invite_link мастера
  return startParam
}

// ============================================
// ⚙️ КОНФИГУРАЦИЯ AXIOS
// ============================================

const API_BASE = 'http://localhost:8080/api/v1'
const INIT_DATA = getInitData()

console.log('🎯 Режим работы:', window.Telegram?.WebApp ? 'Telegram WebApp' : 'Браузер (тест)')
console.log('👤 Роль:', getUserRole())
console.log('🔗 Invite link:', getInviteLink() || '(не нужен)')
console.log('🔐 InitData (первые 50 символов):', INIT_DATA.substring(0, 50) + '...')

const api = axios.create({
  baseURL: API_BASE,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 🔥 INTERCEPTOR: автоматически добавляет X-Telegram-Init-Data ко ВСЕМ запросам
api.interceptors.request.use((config) => {
  // Для multipart/form-data НЕ устанавливаем Content-Type — axios сделает это сам
  if (!(config.data instanceof FormData)) {
    config.headers['Content-Type'] = 'application/json'
  }
  
  config.headers['X-Telegram-Init-Data'] = INIT_DATA
  return config
})

// 🔥 INTERCEPTOR: обработка ошибок авторизации
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      console.error('❌ Ошибка авторизации. Возможно, initData невалиден.')
      // Можно показать alert или редирект
    }
    return Promise.reject(error)
  }
)

// ============================================
// 🔓 ПУБЛИЧНЫЕ КЛИЕНТСКИЕ ЭНДПОИНТЫ
// ============================================

export async function fetchServices() {
  const inviteLink = getInviteLink()
  if (!inviteLink) {
    throw new Error('invite_link не найден. Откройте Mini App по ссылке мастера.')
  }
  const { data } = await api.get(`/invite/${inviteLink}/services`)
  return data
}

export async function fetchSlots(date, serviceId) {
  const inviteLink = getInviteLink()
  if (!inviteLink) {
    throw new Error('invite_link не найден.')
  }
  const { data } = await api.get(`/invite/${inviteLink}/slots`, {
    params: { date, service_id: serviceId },
  })
  return data
}

/**
 * Получает агрегированную инфу о мастере (с фото)
 */
export async function fetchMasterInfo() {
  const inviteLink = getInviteLink()
  if (!inviteLink) {
    throw new Error('invite_link не найден.')
  }
  const { data } = await api.get(`/invite/${inviteLink}/info`)
  return data
}

// ============================================
// 🔐 ЗАЩИЩЁННЫЕ КЛИЕНТСКИЕ ЭНДПОИНТЫ
// ============================================

export async function bookSlot(params) {
  const { data } = await api.post('/book', params)
  return data
}

export async function fetchClientBookings() {
  const { data } = await api.get('/client/bookings')
  return data
}

export async function cancelBooking(bookingId) {
  const { data } = await api.post(`/client/bookings/${bookingId}/cancel`, {})
  return data
}

// ============================================
// 🔐 МАСТЕРСКИЕ ЭНДПОИНТЫ
// ============================================

export async function fetchTodaySchedule() {
  const { data } = await api.get('/master/today')
  return data
}

export async function fetchWaitlist() {
  const { data } = await api.get('/master/waitlist')
  return data
}

export async function fetchMasterProfile() {
  const { data } = await api.get('/master/profile')
  return data
}

export async function updateMasterProfile(profileData) {
  const { data } = await api.put('/master/profile', profileData)
  return data
}

export async function fetchMasterServices() {
  const { data } = await api.get('/master/services')
  return data
}

export async function createService(serviceData) {
  const { data } = await api.post('/master/services', serviceData)
  return data
}

export async function updateService(serviceId, serviceData) {
  const { data } = await api.put(`/master/services/${serviceId}`, serviceData)
  return data
}

export async function deleteService(serviceId) {
  const { data } = await api.delete(`/master/services/${serviceId}`)
  return data
}

export async function updateSettings(settingsData) {
  const { data } = await api.put('/master/settings', settingsData)
  return data
}

// ============================================
// 📸 ФОТО РАБОТ МАСТЕРА
// ============================================

export const STATIC_BASE_URL = 'http://localhost:8080'

/**
 * Превращает относительный URL фото в абсолютный
 */
export function normalizePhotoUrl(url) {
  if (!url) return ''
  if (url.startsWith('http')) return url
  return `${STATIC_BASE_URL}${url.startsWith('/') ? '' : '/'}${url}`
}

export async function fetchPhotos() {
  const { data } = await api.get('/master/photos')
  return Array.isArray(data) ? data : (data?.photos || [])
}

export async function uploadPhoto(file) {
  const formData = new FormData()
  formData.append('file', file) // Имя поля 'file' — как в Go (r.FormFile("file"))
  
  const { data } = await api.post('/master/photos', formData, {
    // ⚠️ Content-Type НЕ устанавливаем вручную — axios сам поставит multipart/form-data с boundary
  })
  return data
}

export async function deletePhoto(photoId) {
  const { data } = await api.delete(`/master/photos/${photoId}`)
  return data
}

// ============================================
// 🎨 ТЕЛЕГРАМ-СПЕЦИФИЧНЫЕ ФИЧИ
// ============================================

/**
 * Показывает нативный popup Telegram
 */
export function showTelegramPopup(message) {
  if (window.Telegram?.WebApp) {
    window.Telegram.WebApp.showAlert(message)
  } else {
    alert(message)
  }
}

/**
 * Показывает нативный confirm Telegram
 */
export function showTelegramConfirm(message, callback) {
  if (window.Telegram?.WebApp) {
    window.Telegram.WebApp.showConfirm(message, callback)
  } else {
    callback(confirm(message))
  }
}

/**
 * Haptic feedback (вибрация)
 */
export function hapticFeedback(type = 'light') {
  if (window.Telegram?.WebApp?.HapticFeedback) {
    window.Telegram.WebApp.HapticFeedback.impactOccurred(type)
  }
}

/**
 * Закрыть Mini App
 */
export function closeWebApp() {
  if (window.Telegram?.WebApp) {
    window.Telegram.WebApp.close()
  }
}

// Экспортируем сам WebApp для прямого доступа
export const TelegramWebApp = window.Telegram?.WebApp || null
// ngrok http http://127.0.0.1:5173