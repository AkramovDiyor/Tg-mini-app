import axios from 'axios'
import { getInitData } from './telegram'

const API_BASE = import.meta.env.VITE_API_BASE || 'https://tg-mini-app-j29w.onrender.com/api/v1'
export const STATIC_BASE_URL = import.meta.env.VITE_STATIC_BASE || 'https://tg-mini-app-j29w.onrender.com'

const INIT_DATA = getInitData()

export const http = axios.create({
  baseURL: API_BASE,
  headers: { 'Content-Type': 'application/json' },
})

http.interceptors.request.use((config) => {
  if (config.data instanceof FormData) {
    delete config.headers['Content-Type']
  } else {
    config.headers['Content-Type'] = 'application/json'
  }
  config.headers['X-Telegram-Init-Data'] = INIT_DATA
  config.headers['ngrok-skip-browser-warning'] = 'true'
  return config
})

/**
 * Go http.Error отдаёт text/plain, не { message }.
 * @param {unknown} err
 * @param {string} fallback
 * @returns {string}
 */
export function getApiErrorMessage(err, fallback = 'Что-то пошло не так') {
  const axiosErr = /** @type {{ response?: { data?: unknown }, message?: string }} */ (err)
  const data = axiosErr?.response?.data
  if (typeof data === 'string' && data.trim()) return data
  if (data && typeof data === 'object' && 'message' in data) {
    const msg = /** @type {{ message?: string }} */ (data).message
    if (msg) return msg
  }
  return axiosErr?.message || fallback
}
