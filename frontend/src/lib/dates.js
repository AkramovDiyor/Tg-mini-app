// Все helpers по датам и календарю. Без React-зависимостей.

export const SLOT_TIMES = ['10:00', '11:30', '13:00', '14:30', '16:00']

export const MONTHS_GEN = [
  'января', 'февраля', 'марта', 'апреля', 'мая', 'июня',
  'июля', 'августа', 'сентября', 'октября', 'ноября', 'декабря',
]

// Короткие названия месяцев для таймера
export const MONTHS_SHORT = [
  'янв', 'фев', 'мар', 'апр', 'мая', 'июн',
  'июл', 'авг', 'сен', 'окт', 'ноя', 'дек',
]

export const MONTHS_TITLE = [
  'Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь',
  'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь',
]
export const WEEKDAYS_SHORT = ['Вс', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб']
export const WEEKDAYS_GRID = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']

const startOfDay = (d) => new Date(d.getFullYear(), d.getMonth(), d.getDate())
export const TODAY = startOfDay(new Date())

export const toISO = (d) =>
  `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`

export const fromISO = (iso) => {
  const [y, m, d] = iso.split('-').map(Number)
  return new Date(y, m - 1, d)
}

export const daysFromToday = (d) => Math.round((startOfDay(d) - TODAY) / 86400000)

export const formatTime = (d) =>
  `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`

export function buildQuickDays(count = 7) {
  const res = []
  for (let i = 0; i < count; i++) {
    const date = new Date(TODAY.getFullYear(), TODAY.getMonth(), TODAY.getDate() + i)
    res.push({
      date,
      iso: toISO(date),
      label: i === 0 ? 'Сегодня' : i === 1 ? 'Завтра' : WEEKDAYS_SHORT[date.getDay()],
    })
  }
  return res
}

export function buildMonthGrid(year, month) {
  const firstDay = new Date(year, month, 1)
  const daysInMonth = new Date(year, month + 1, 0).getDate()

  let offset = firstDay.getDay() - 1
  if (offset < 0) offset = 6

  const cells = Array(offset).fill(null)
  for (let d = 1; d <= daysInMonth; d++) cells.push(new Date(year, month, d))
  return cells
}

const DEMO_PATTERNS = {
  0: ['free', 'busy', 'free', 'pending', 'free'],
  1: ['free', 'free', 'busy', 'free', 'pending'],
  2: ['busy', 'busy', 'busy', 'busy', 'busy'],
}

export function getScheduleForDate(date) {
  const diff = daysFromToday(date)
  if (DEMO_PATTERNS[diff]) return DEMO_PATTERNS[diff]

  const seed = date.getFullYear() * 10000 + (date.getMonth() + 1) * 100 + date.getDate()
  return SLOT_TIMES.map((_, i) => {
    const v = (seed * 31 + i * 57) % 100
    if (v < 55) return 'free'
    if (v < 85) return 'busy'
    return 'pending'
  })
}

export function getBookingTimeLabel(target, now = new Date()) {
  const diff = target - now
  if (diff <= 0) return 'Скоро начнётся'

  const HOUR = 3600000
  
  // Больше 24 часов: "20 авг, 11:00"
  if (diff > 24 * HOUR) {
    return `${target.getDate()} ${MONTHS_SHORT[target.getMonth()]}, ${formatTime(target)}`
  }
  
  // 3-24 часа: "Сегодня, 11:00"
  if (diff > 3 * HOUR) {
    return `Сегодня, ${formatTime(target)}`
  }

  // Меньше 3 часов: обратный отсчёт
  const h = Math.floor(diff / HOUR)
  const m = Math.floor((diff % HOUR) / 60000)
  return h > 0 ? `Через ${h} ч ${m} мин` : `Через ${m} мин`
}



// export function formatTimeFromISO(isoString) {
//   if (!isoString) return ''
  
//   // UTC формат: "2026-08-28T18:00:00Z"
//   if (isoString.endsWith('Z')) {
//     const timePart = isoString.split('T')[1]
//     const [hours, minutes] = timePart.split(':')
//     return `${hours}:${minutes}`
//   }
  
//   // Формат с таймзоной: "2026-08-28T21:00:00+03:00"
//   // Конвертируем в UTC чтобы получить оригинальное время
//   const date = new Date(isoString)
//   const hours = String(date.getUTCHours()).padStart(2, '0')
//   const minutes = String(date.getUTCMinutes()).padStart(2, '0')
//   return `${hours}:${minutes}`
// }


// Правильная функция форматирования времени (Локальное время)
export function formatTimeFromISO(isoString) {
  if (!isoString) return ''
  const date = new Date(isoString)
  // Используем getHours(), а не getUTCHours()!
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${hours}:${minutes}`
}

// Правильная функция таймера
export function getBookingTimeLabelFromISO(isoString, now = new Date()) {
  if (!isoString) return ''
  
  // new Date(isoString) уже дает нам локальный timestamp
  const targetTime = new Date(isoString).getTime()
  const nowTime = now.getTime()
  const diff = targetTime - nowTime
  
  if (diff <= 0) return 'Запись прошла'
  
  const HOUR = 3600000
  const timeStr = formatTimeFromISO(isoString)
  
  // Больше 24 часов: "Запись 20 авг, 11:00"
  if (diff > 24 * HOUR) {
    const targetDate = new Date(isoString)
    return `Запись ${targetDate.getDate()} ${MONTHS_SHORT[targetDate.getMonth()]}, ${timeStr}`
  }
  
  // 3-24 часа: "Сегодня, 11:00"
  if (diff > 3 * HOUR) {
    return `Сегодня, ${timeStr}`
  }
  
  // Меньше 3 часов: обратный отсчёт
  const h = Math.floor(diff / HOUR)
  const m = Math.floor((diff % HOUR) / 60000)
  return h > 0 ? `Через ${h} ч ${m} мин` : `Через ${m} мин`
}
