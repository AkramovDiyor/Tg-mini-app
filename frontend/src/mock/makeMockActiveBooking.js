import { formatTime, toISO } from '../lib/dates'

export const makeMockActiveBooking = () => {
  const date = new Date(Date.now() + 2 * 60 * 60 * 1000)
  return {
    iso: toISO(date),
    time: formatTime(date),
    service: 'Мужская стрижка',
    date,
  }
}