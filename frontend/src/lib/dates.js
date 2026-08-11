// Все helpers по датам и календарю. Без React-зависимостей.

export const SLOT_TIMES = ['10:00', '11:30', '13:00', '14:30', '16:00']

export const MONTHS_GEN = [
  'января', 'февраля', 'марта', 'апреля', 'мая', 'июня',
  'июля', 'августа', 'сентября', 'октября', 'ноября', 'декабря',
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
  2: ['busy', 'busy', 'busy', 'busy', 'busy'], // ← демо листа ожидания
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
  if (diff > 24 * HOUR) {
    return `${target.getDate()} ${MONTHS_GEN[target.getMonth()]}, ${formatTime(target)}`
  }
  if (diff > 3 * HOUR) {
    return `Сегодня, ${formatTime(target)}`
  }

  const h = Math.floor(diff / HOUR)
  const m = Math.floor((diff % HOUR) / 60000)
  return h > 0 ? `Через ${h} ч ${m} мин` : `Через ${m} мин`
}