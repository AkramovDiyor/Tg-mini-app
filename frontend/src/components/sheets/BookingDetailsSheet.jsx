import { useEffect, useState } from 'react'
import { User, MapPin, Clock } from 'lucide-react'
import { SheetShell } from './SheetShell'
import { getBookingTimeLabelFromISO, WEEKDAYS_SHORT, MONTHS_GEN } from '../../lib/dates'
import { rub } from '../../lib/currency'

export function BookingDetailsSheet({ booking, onClose, onCancel }) {
  const [now, setNow] = useState(() => new Date())

  // Обновляем таймер каждую минуту для обратного отсчёта
  useEffect(() => {
    const t = setInterval(() => setNow(new Date()), 60000)
    return () => clearInterval(t)
  }, [])

  if (!booking) return null

  // Парсим дату из строки
  // Просто парсим ISO строку, браузер сам даст локальную дату
  const bookingDate = new Date(booking.start_time)
  const weekday = WEEKDAYS_SHORT[bookingDate.getDay()] // getDay() без UTC
  const day = bookingDate.getDate()
  const month = bookingDate.getMonth() + 1
  const fullDate = `${weekday}, ${day} ${MONTHS_GEN[month - 1]}`

  // Умный таймер
  const timeLabel = getBookingTimeLabelFromISO(booking.start_time, now)
  
  return (
    <SheetShell onClose={onClose}>
      {/* Заголовок */}
      <p className="text-center text-[11px] font-bold uppercase tracking-widest text-slate-400">
        Моя запись
      </p>

      {/* Услуга */}
      <h3 className="mt-3 text-center text-2xl font-extrabold leading-snug text-slate-900">
        {booking.service_name}
      </h3>

      {/* Дата */}
      <p className="mt-1 text-center text-sm capitalize text-slate-400">
        {fullDate}
      </p>

      {/* Умный таймер */}
      <p className="mt-1 text-center text-2xl font-extrabold leading-snug text-emerald-600">
        {timeLabel}
      </p>

      {/* Детали: цена и длительность */}
      <div className="mt-5 grid grid-cols-2 gap-2">
        <div className="rounded-2xl bg-slate-50 p-3 text-center">
          <p className="text-[11px] font-semibold uppercase tracking-widest text-slate-400">
            Стоимость
          </p>
          <p className="mt-1 text-lg font-extrabold text-slate-900">
            {rub(booking.service_price)}
          </p>
        </div>
        <div className="rounded-2xl bg-slate-50 p-3 text-center">
          <p className="text-[11px] font-semibold uppercase tracking-widest text-slate-400">
            Длительность
          </p>
          <p className="mt-1 flex items-center justify-center gap-1 text-lg font-extrabold text-slate-900">
            <Clock className="h-4 w-4" />
            {booking.service_duration || 0} мин
          </p>
        </div>
      </div>

      {/* Мастер */}
      <div className="mt-4 flex items-center gap-3 rounded-2xl border border-slate-100 bg-slate-50 p-4 text-left">
        <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-emerald-400 to-teal-600 text-white">
          <User className="h-6 w-6" />
        </div>
        <div className="min-w-0 flex-1">
          <p className="font-bold text-slate-900">{booking.master_name}</p>
          <p className="mt-0.5 flex items-center gap-1 text-xs text-slate-400">
            <MapPin className="h-3.5 w-3.5 shrink-0" />
            {booking.master_address}
          </p>
        </div>
      </div>

      {/* Информация о подтверждении */}
      <p className="mt-3 flex items-start gap-1.5 px-2 text-[11px] leading-relaxed text-slate-400">
        <Clock className="mt-0.5 h-3.5 w-3.5 shrink-0" />
        За 2 часа до начала придёт кнопка подтверждения. Если не нажать — запись может быть передана другому.
      </p>

      {/* Отмена */}
      <button
        onClick={() => onCancel(booking.booking_id)}
        className="mt-4 w-full rounded-xl bg-red-50 py-4 font-bold text-red-600 transition active:scale-[0.98]"
      >
        Отменить запись
      </button>
    </SheetShell>
  )
}