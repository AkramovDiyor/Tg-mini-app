import axios from 'axios'

const API_BASE = 'http://10.21.33.135/api/v1'
const INVITE_LINK = 'LINK123243'

const api = axios.create({
  baseURL: API_BASE,
  headers: {
    'Content-Type': 'application/json',
  },
})

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
  const { data } = await api.post(`/client/bookings/${bookingId}/cancel`, {},
    {
    headers: {
      'X-Telegram-Init-Data': 'test-vasya',
    },
  })
  return data
}