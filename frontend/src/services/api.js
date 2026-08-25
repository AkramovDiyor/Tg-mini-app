import axios from 'axios'

// Инициализация Telegram WebApp
if (window.Telegram?.WebApp) {
  window.Telegram.WebApp.ready()
  window.Telegram.WebApp.expand()
  window.Telegram.WebApp.enableClosingConfirmation()
}

// Получение initData
function getInitData() {
  if (window.Telegram?.WebApp?.initData && window.Telegram.WebApp.initData !== '') {
    return window.Telegram.WebApp.initData
  }
  
  // ТЕСТОВЫЙ РЕЖИМ
  const urlParams = new URLSearchParams(window.location.search)
  const mode = urlParams.get('mode')
  const startapp = urlParams.get('startapp')
  
  if (mode === 'master' || startapp === 'master') {
    return 'test-master'
  }
  return 'test-vasya'
}

const INIT_DATA = getInitData()

const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080/api/v1'
export const STATIC_BASE_URL = import.meta.env.VITE_STATIC_BASE || 'http://localhost:8080'

const api = axios.create({
  baseURL: API_BASE,
  headers: { 'Content-Type': 'application/json' },
})

api.interceptors.request.use((config) => {
  if (!(config.data instanceof FormData)) {
    config.headers['Content-Type'] = 'application/json'
  }
  config.headers['X-Telegram-Init-Data'] = INIT_DATA
  config.headers['ngrok-skip-browser-warning'] = 'true'
  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      console.error('❌ Ошибка авторизации')
    }
    return Promise.reject(error)
  }
)

// 🔥 ОПРЕДЕЛЕНИЕ РОЛИ ЧЕРЕЗ БЭКЕНД
let cachedIdentity = null

export async function fetchUserIdentity() {
  if (cachedIdentity) return cachedIdentity
  
  try {
    const { data } = await api.get('/me')
    cachedIdentity = data
    console.log('🎯 Роль определена:', data.role, data)
    return data
  } catch (error) {
    console.error('❌ Не удалось определить роль:', error)
    throw error
  }
}

export function getUserIdentity() {
  return cachedIdentity
}

// 🔓 ПУБЛИЧНЫЕ КЛИЕНТСКИЕ ЭНДПОИНТЫ
export async function fetchServices(inviteLink) {
  if (!inviteLink) throw new Error('invite_link required')
  const { data } = await api.get(`/invite/${inviteLink}/services`)
  return data
}

export async function fetchSlots(date, serviceId, inviteLink) {
  if (!inviteLink) throw new Error('invite_link required')
  const { data } = await api.get(`/invite/${inviteLink}/slots`, {
    params: { date, service_id: serviceId },
  })
  return data
}

export async function fetchMasterInfo(inviteLink) {
  if (!inviteLink) throw new Error('invite_link required')
  const { data } = await api.get(`/invite/${inviteLink}/info`)
  return data
}

// 🔐 ЗАЩИЩЁННЫЕ КЛИЕНТСКИЕ ЭНДПОИНТЫ
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

// 🔐 МАСТЕРСКИЕ ЭНДПОИНТЫ
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

// 📸 ФОТО
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
  formData.append('file', file)
  const { data } = await api.post('/master/photos', formData)
  return data
}

export async function deletePhoto(photoId) {
  const { data } = await api.delete(`/master/photos/${photoId}`)
  return data
}

// 🎨 TELEGRAM ФИЧИ
export function showTelegramPopup(message) {
  if (window.Telegram?.WebApp) window.Telegram.WebApp.showAlert(message)
  else alert(message)
}

export function hapticFeedback(type = 'light') {
  if (window.Telegram?.WebApp?.HapticFeedback) {
    window.Telegram.WebApp.HapticFeedback.impactOccurred(type)
  }
}

export const TelegramWebApp = window.Telegram?.WebApp || null