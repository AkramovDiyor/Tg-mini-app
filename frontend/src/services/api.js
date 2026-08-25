import { http, STATIC_BASE_URL, getApiErrorMessage } from '../lib/http'
import { getBrowserTestIdentity } from '../lib/telegram'

export { STATIC_BASE_URL, getApiErrorMessage }
export {
  showTelegramPopup,
  hapticFeedback,
  TelegramWebApp,
  getClientDisplayName,
} from '../lib/telegram'

/**
 * @returns {Promise<Identity>}
 */
export async function fetchUserIdentity() {
  const test = getBrowserTestIdentity()
  if (test) return test
  const { data } = await http.get('/me')
  return data
}

/**
 * @param {string} inviteLink
 * @returns {Promise<Service[]>}
 */
export async function fetchServices(inviteLink) {
  if (!inviteLink) throw new Error('invite_link required')
  const { data } = await http.get(`/invite/${inviteLink}/services`)
  return Array.isArray(data) ? data : []
}

/**
 * @param {string} date
 * @param {number} serviceId
 * @param {string} inviteLink
 * @returns {Promise<Slot[]>}
 */
export async function fetchSlots(date, serviceId, inviteLink) {
  if (!inviteLink) throw new Error('invite_link required')
  const { data } = await http.get(`/invite/${inviteLink}/slots`, {
    params: { date, service_id: serviceId },
  })
  return Array.isArray(data) ? data : []
}

/**
 * @param {string} inviteLink
 * @returns {Promise<MasterInfo>}
 */
export async function fetchMasterInfo(inviteLink) {
  if (!inviteLink) throw new Error('invite_link required')
  const { data } = await http.get(`/invite/${inviteLink}/info`)
  return data
}

/**
 * @param {BookSlotRequest} params
 */
export async function bookSlot(params) {
  const { data } = await http.post('/book', {
    start_time: params.start_time,
    service_id: params.service_id,
    name: params.name,
    price: params.price,
  })
  return data
}

/**
 * @returns {Promise<ClientBooking[]>}
 */
export async function fetchClientBookings() {
  const { data } = await http.get('/client/bookings')
  return Array.isArray(data) ? data : []
}

/**
 * @param {number} bookingId
 */
export async function cancelBooking(bookingId) {
  const { data } = await http.post(`/client/bookings/${bookingId}/cancel`, {})
  return data
}

export async function fetchTodaySchedule() {
  const { data } = await http.get('/master/today')
  return data
}

export async function fetchWaitlist() {
  const { data } = await http.get('/master/waitlist')
  return Array.isArray(data) ? data : []
}

export async function fetchMasterProfile() {
  const { data } = await http.get('/master/profile')
  return data
}

/**
 * @param {{ name: string, bio: string, address: string }} profileData
 */
export async function updateMasterProfile(profileData) {
  const { data } = await http.put('/master/profile', profileData)
  return data
}

export async function fetchMasterServices() {
  const { data } = await http.get('/master/services')
  return Array.isArray(data) ? data : []
}

/**
 * @param {{ name: string, duration_min: number, price: number }} serviceData
 */
export async function createService(serviceData) {
  const { data } = await http.post('/master/services', serviceData)
  return data
}

/**
 * @param {number} serviceId
 * @param {{ name: string, duration_min: number, price: number }} serviceData
 */
export async function updateService(serviceId, serviceData) {
  const { data } = await http.put(`/master/services/${serviceId}`, serviceData)
  return data
}

/**
 * @param {number} serviceId
 */
export async function deleteService(serviceId) {
  const { data } = await http.delete(`/master/services/${serviceId}`)
  return data
}

/**
 * @param {UpdateSettingsRequest} settingsData
 */
export async function updateSettings(settingsData) {
  const { data } = await http.put('/master/settings', settingsData)
  return data
}

/**
 * @param {string} url
 */
export function normalizePhotoUrl(url) {
  if (!url) return ''
  if (url.startsWith('http')) return url
  return `${STATIC_BASE_URL}${url.startsWith('/') ? '' : '/'}${url}`
}

export async function fetchPhotos() {
  const { data } = await http.get('/master/photos')
  return Array.isArray(data) ? data : data?.photos || []
}

/**
 * @param {File} file
 */
export async function uploadPhoto(file) {
  const formData = new FormData()
  formData.append('file', file)
  const { data } = await http.post('/master/photos', formData)
  return data
}

/**
 * @param {number} photoId
 */
export async function deletePhoto(photoId) {
  const { data } = await http.delete(`/master/photos/${photoId}`)
  return data
}
