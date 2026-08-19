import axios from 'axios'

// ============================================
// 🕵️ ЛОГИКА ОПРЕДЕЛЕНИЯ РОЛИ (Мастер/Клиент)
// ============================================
function getStartParam() {
  // Сначала пытаемся достать параметр из реального Telegram WebApp
  if (window.Telegram && window.Telegram.WebApp && window.Telegram.WebApp.initDataUnsafe) {
    if (window.Telegram.WebApp.initDataUnsafe.start_param) {
      return window.Telegram.WebApp.initDataUnsafe.start_param
    }
  }
  // Если мы в обычном браузере (для теста) — читаем из URL
  const urlParams = new URLSearchParams(window.location.search)
  return urlParams.get('startapp') || ''
}

function getInitData() {
  // Если мы в реальном Telegram, берем криптографическую подпись
  if (window.Telegram && window.Telegram.WebApp && window.Telegram.WebApp.initData) {
    return window.Telegram.WebApp.initData
  }
  
  // ВРЕМЕННО ДЛЯ ТЕСТА В БРАУЗЕРЕ:
  // Подставляем тестовые заглушки в зависимости от того, по какой ссылке мы открыли
  const startParam = getStartParam()
  // if (startParam === 'master') {
  //   return 'test-pedro'
  // }
  return 'test-pedro'
}

// ============================================
// ⚙️ КОНФИГУРАЦИЯ AXIOS
// ============================================
const API_BASE = 'http://localhost:8080/api/v1'
const START_PARAM = getStartParam()

// Если открыл клиентскую часть, подставляем ссылку мастера. Если мастер — оставляем тестовую.
const INVITE_LINK = START_PARAM !== 'master' && START_PARAM !== '' ? START_PARAM : '7u72y9b6'
const TEST_INIT_DATA = getInitData()

const api = axios.create({
  baseURL: API_BASE,
  headers: {
    'Content-Type': 'application/json',
    'X-Telegram-Init-Data': TEST_INIT_DATA, //Interceptor — добавляется ко ВСЕМ запросам
  },
})

// ============================================
// 🔓 ПУБЛИЧНЫЕ КЛИЕНТСКИЕ ЭНДПОИНТЫ
// ============================================

export async function fetchServices() {
  const { data } = await api.get(`/invite/${INVITE_LINK}/services`)
  return data
}

export async function fetchSlots(date, serviceId) {
  const { data } = await api.get(`/invite/${INVITE_LINK}/slots`, {
    params: { date, service_id: serviceId },
  })
  return data
}

export async function fetchMasterInfo() {
  const { data } = await api.get(`/invite/${INVITE_LINK}/info`)
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
// 🔐 МАСТЕРСКИЕ ЭНДПОИНТЫ (защищённые)
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
 * Нормализует URL фото: если относительный путь — добавляет базовый URL
 */
export function normalizePhotoUrl(url) {
  if (!url) return ''
  if (url.startsWith('http')) return url
  return `${STATIC_BASE_URL}${url.startsWith('/') ? '' : '/'}${url}`
}

export async function fetchPhotos() {
  const { data } = await api.get('/master/photos')
  // Поддержка как массива, так и { photos: [...] }
  return Array.isArray(data) ? data : (data?.photos || [])
}

export async function uploadPhoto(file) {
  const formData = new FormData()
  formData.append('file', file) // ⚠️ проверьте имя поля в бэкенде (может быть 'photo' или 'image')
  const { data } = await api.post('/master/photos', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
  return data
}

export async function deletePhoto(photoId) {
  const { data } = await api.delete(`/master/photos/${photoId}`)
  return data
}