import axios from 'axios'

const API_BASE = 'http://localhost:8080/api/v1'
const INVITE_LINK = 'LINK123243'

const api = axios.create({
  baseURL: API_BASE,
  headers: {
    'Content-Type': 'application/json',
  },
})

// ===== КЛИЕНТСКИЕ ЭНДПОИНТЫ =====

export async function fetchServices() {
  const { data } = await api.get(`/master/${INVITE_LINK}/services`)
  return data
}

export async function fetchSlots(date, serviceId) {
  const { data } = await api.get(`/master/${INVITE_LINK}/slots`, {
    params: { date, service_id: serviceId },
  })
  return data
}

export async function bookSlot(params) {
  const { data } = await api.post('/book', params, {
    headers: {
      'X-Telegram-Init-Data': 'test-vasya',
    },
  })
  return data
}

export async function fetchClientBookings() {
  const { data } = await api.get('/client/bookings', {
    headers: {
      'X-Telegram-Init-Data': 'test-vasya',
    },
  })
  return data
}

export async function cancelBooking(bookingId) {
  const { data } = await api.post(`/client/bookings/${bookingId}/cancel`, {}, {
    headers: {
      'X-Telegram-Init-Data': 'test-vasya',
    },
  })
  return data
}

// ===== МАСТЕРСКИЕ ЭНДПОИНТЫ =====

export async function fetchTodaySchedule() {
  const { data } = await api.get('master/today', {
    headers: {
      'X-Telegram-Init-Data': 'test-vasya',
    },
  })
  return data
}

export async function fetchWaitlist() {
  const { data } = await api.get('master/waitlist', {
    headers: {
      'X-Telegram-Init-Data': 'test-vasya',
    },
  })
  return data
}

export async function fetchMasterProfile() {
  const { data } = await api.get('master/profile', {
    headers: {
      'X-Telegram-Init-Data': 'test-vasya',
    },
  })
  return data
}

export async function updateMasterProfile(profileData) {
  const { data } = await api.put('master/profile', profileData, {
    headers: {
      'X-Telegram-Init-Data': 'test-vasya',
    },
  })
  return data
}

export async function fetchMasterServices() {
  const { data } = await api.get('master/services', {
    headers: {
      'X-Telegram-Init-Data': 'test-vasya',
    },
  })
  return data
}

export async function createService(serviceData) {
  const { data } = await api.post('master/services', serviceData, {
    headers: {
      'X-Telegram-Init-Data': 'test-vasya',
    },
  })
  return data
}

export async function updateSettings(settingsData) {
  const { data } = await api.put('master/settings', settingsData, {
    headers: {
      'X-Telegram-Init-Data': 'test-vasya',
    },
  })
  return data
}