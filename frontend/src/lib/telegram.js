/**
 * @returns {string}
 */
export function getInitData() {
  const real = window.Telegram?.WebApp?.initData
  if (real) return real

  const urlParams = new URLSearchParams(window.location.search)
  const mode = urlParams.get('mode')
  const startapp = urlParams.get('startapp')
  if (mode === 'master' || startapp === 'master') return 'test-master'
  return 'test-vasya'
}

export function initTelegramWebApp() {
  const webApp = window.Telegram?.WebApp
  if (!webApp) return
  webApp.ready()
  webApp.expand()
  webApp.enableClosingConfirmation()
}

/**
 * @returns {Identity | null}
 */
export function getBrowserTestIdentity() {
  if (window.Telegram?.WebApp?.initData) return null

  const urlParams = new URLSearchParams(window.location.search)
  const mode = urlParams.get('mode')
  const startapp = urlParams.get('startapp')

  if (mode === 'master') {
    return {
      role: 'master',
      master: { telegram_id: 999999, name: 'Тестовый Мастер' },
    }
  }

  if (startapp && startapp !== 'master') {
    return {
      role: 'client',
      invite_link: startapp,
      user: { telegram_id: 777111222 },
    }
  }

  return null
}

/**
 * @returns {string}
 */
export function getClientDisplayName() {
  const user = window.Telegram?.WebApp?.initDataUnsafe?.user
  if (!user?.first_name) return 'Клиент'
  return [user.first_name, user.last_name].filter(Boolean).join(' ')
}

export function showTelegramPopup(message) {
  if (window.Telegram?.WebApp) window.Telegram.WebApp.showAlert(message)
  else window.alert(message)
}

/**
 * @param {'light' | 'medium' | 'heavy' | 'rigid' | 'soft'} [type]
 */
export function hapticFeedback(type = 'light') {
  window.Telegram?.WebApp?.HapticFeedback?.impactOccurred(type)
}

export const TelegramWebApp = window.Telegram?.WebApp || null
