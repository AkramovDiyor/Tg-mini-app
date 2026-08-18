import { useEffect, useState } from 'react'
import { Scissors, ChevronRight } from 'lucide-react'
import { getBookingTimeLabel } from '../lib/dates'
import { fetchClientBookings } from '../services/api'

export function ActiveBookingBar({ onOpenDetails }) {
  const [activeBooking, setActiveBooking] = useState(null)
  const [loading, setLoading] = useState(true)
  const [now, setNow] = useState(() => new Date())

  useEffect(() => {
    loadBookings()
  }, [])

  // Обновляем таймер каждую минуту
  useEffect(() => {
    const t = setInterval(() => {
      setNow(new Date())
    }, 60000)
    return () => clearInterval(t)
  }, [])

  // Каждую минуту перезапрашиваем записи — чтобы убрать те, что прошли
  useEffect(() => {
    const t = setInterval(() => {
      loadBookings()
    }, 60000)
    return () => clearInterval(t)
  }, [])

  const loadBookings = async () => {
    try {
      setLoading(true)
      const bookings = await fetchClientBookings()

      if (!Array.isArray(bookings)) {
        setActiveBooking(null)
        return
      }

      const nowTime = new Date()

      // СТРОГИЙ ФИЛЬТР: только будущие активные записи
      const futureBookings = bookings
        .filter((b) => {
          // Исключаем отменённые
          if (b.status === 'cancelled' || b.status === 'canceled') return false
          // Только будущие
          return new Date(b.start_time) > nowTime
        })
        .sort((a, b) => new Date(a.start_time) - new Date(b.start_time))

      setActiveBooking(futureBookings[0] || null)
    } catch (err) {
      console.error('Failed to load bookings:', err)
      setActiveBooking(null)
    } finally {
      setLoading(false)
    }
  }

  // Нет будущих записей — бар не показывается
  if (loading || !activeBooking) {
    return null
  }

  const bookingDate = new Date(activeBooking.start_time)
  const timeLabel = getBookingTimeLabel(bookingDate, now)

  return (
    <button
      onClick={() => onOpenDetails(activeBooking)}
      className="fixed bottom-4 left-1/2 z-40 w-[calc(100%-2.5rem)] max-w-[380px] -translate-x-1/2 animate-fade-up"
    >
      <div className="flex items-center justify-between rounded-2xl bg-slate-900 p-3.5 text-white shadow-2xl">
        <span className="flex shrink-0 items-center justify-center rounded-full bg-white/10 p-2">
          <Scissors className="h-5 w-5 text-emerald-400" />
        </span>
        <span className="mx-3 flex min-w-0 flex-1 flex-col items-start text-left">
          <span className="text-xs font-bold text-emerald-400">
            {timeLabel}
          </span>
          <span className="w-full truncate text-sm font-semibold">
            {activeBooking.service_name || 'Услуга'}
          </span>
        </span>
        <ChevronRight className="h-5 w-5 shrink-0 text-slate-400" />
      </div>
    </button>
  )
}