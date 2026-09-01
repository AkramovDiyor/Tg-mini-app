/**
 * @returns {string}
 */

export function getInitData() {
  // 🔥 Пытаемся получить реальный initData от Telegram
  const realInitData = window.Telegram?.WebApp?.initData
  
  if (realInitData && realInitData.trim() !== '') {
    console.log('✅ Используем РЕАЛЬНЫЙ initData от Telegram')
    return realInitData
  }
  
  // Fallback для тестов в браузере
  console.log('⚠️ Telegram initData отсутствует, тестовый режим')
  const urlParams = new URLSearchParams(window.location.search)
  const mode = urlParams.get('mode')
  const startapp = urlParams.get('startapp')
  
  if (mode === 'master' || startapp === 'master') {
    return 'test-master'
  }
  
  return 'test-vasya'
}

export function getBrowserTestIdentity() {
  // 🔥 Если есть реальный initData — не используем тестовые данные
  if (window.Telegram?.WebApp?.initData && window.Telegram.WebApp.initData.trim() !== '') {
    return null // Пусть бэкенд определяет роль
  }

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


/**
 * @returns {string}
 */
export function getClientDisplayName() {
  try {
    // Пытаемся достать данные пользователя из Telegram WebApp
    const user = window.Telegram?.WebApp?.initDataUnsafe?.user;
    
    if (user) {
      // Собираем имя и фамилию, отфильтровывая пустые значения
      const nameParts = [user.first_name, user.last_name].filter(Boolean);
      
      if (nameParts.length > 0) {
        return nameParts.join(' '); // Например: "Akramov Diyor"
      }
    }
  } catch (error) {
    console.warn('⚠️ Не удалось получить имя из Telegram:', error);
  }
  
  // Фолбэк для тестов в браузере или если имя не удалось получить
  return 'Клиент';
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
