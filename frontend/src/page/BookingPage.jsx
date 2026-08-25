import { useCallback, useMemo, useState } from 'react'
import { ArrowLeft, Calendar, Clock } from 'lucide-react'
import { QueryRetry } from '../components/AppFrame'
import { SlotsSkeleton } from '../components/skeletons/Skeletons'
import { CalendarDay } from '../components/ui/CalendarDay'
import { DayPill } from '../components/ui/DayPill'
import { NavArrow } from '../components/ui/NavArrow'
import { SlotButton } from '../components/ui/SlotButton'
import { SlotsLegend } from '../components/ui/SlotsLegend'
import { useClientBookingsQuery, useSlotsQuery } from '../hooks/useClientQueries'
import {
  TODAY,
  MONTHS_TITLE,
  WEEKDAYS_GRID,
  buildQuickDays,
  buildMonthGrid,
  fromISO,
  toISO,
  formatTimeInMasterTz,
} from '../lib/dates'
import { rub } from '../lib/currency'
import { useDragScroll } from '../lib/useDragScroll'
import { useBookingStore } from '../store/bookingStore'
import { WaitlistBlock } from '../widgets/WaitlistBlock'

const QUICK_DAYS = buildQuickDays()
const CANCELLED = new Set(['cancelled', 'canceled', 'cancelled_by_client', 'cancelled_no_show'])

export function BookingPage({ inviteLink }) {
  const service = useBookingStore((s) => s.service)
  const selectedDate = useBookingStore((s) => s.selectedDate)
  const setDate = useBookingStore((s) => s.setDate)
  const goServices = useBookingStore((s) => s.goServices)
  const openSlotSheet = useBookingStore((s) => s.openSlotSheet)

  const daysScrollRef = useDragScroll()
  const [view, setView] = useState({ year: TODAY.getFullYear(), month: TODAY.getMonth() })

  const slotsQuery = useSlotsQuery(inviteLink, selectedDate, service?.id)
  const bookingsQuery = useClientBookingsQuery()

  const diffMonths = (view.year - TODAY.getFullYear()) * 12 + (view.month - TODAY.getMonth())
  const canPrev = diffMonths > 0
  const canNext = diffMonths < 1

  const monthGrid = useMemo(
    () => buildMonthGrid(view.year, view.month),
    [view.year, view.month],
  )

  const mineTimes = useMemo(() => {
    const set = new Set()
    for (const b of bookingsQuery.data || []) {
      if (CANCELLED.has(b.status?.toLowerCase())) continue
      set.add(new Date(b.start_time).getTime())
    }
    return set
  }, [bookingsQuery.data])

  const selectedDateObj = useMemo(() => fromISO(selectedDate), [selectedDate])
  const isToday =
    selectedDateObj.getFullYear() === TODAY.getFullYear() &&
    selectedDateObj.getMonth() === TODAY.getMonth() &&
    selectedDateObj.getDate() === TODAY.getDate()

  const formattedSlots = useMemo(() => {
    const now = Date.now()
    return (slotsQuery.data || [])
      .map((slot) => {
        const startTime = new Date(slot.start_time)
        return {
          start_time: slot.start_time,
          startMs: startTime.getTime(),
          time: formatTimeInMasterTz(slot.start_time),
          status: slot.status === 'booked' ? 'busy' : slot.status,
        }
      })
      .filter((slot) => !isToday || slot.startMs > now)
  }, [slotsQuery.data, isToday])

  const handleSelect = useCallback(
    (iso) => {
      setDate(iso)
      const d = fromISO(iso)
      setView({ year: d.getFullYear(), month: d.getMonth() })
    },
    [setDate],
  )

  const prevMonth = useCallback(() => {
    setView((v) => (v.month === 0 ? { year: v.year - 1, month: 11 } : { ...v, month: v.month - 1 }))
  }, [])

  const nextMonth = useCallback(() => {
    setView((v) => (v.month === 11 ? { year: v.year + 1, month: 0 } : { ...v, month: v.month + 1 }))
  }, [])

  const handleSlotClick = useCallback(
    (startTime) => {
      openSlotSheet({ startTime, iso: selectedDate })
    },
    [openSlotSheet, selectedDate],
  )

  const slots = slotsQuery.data || []
  const isDayOff = !slotsQuery.isPending && !slotsQuery.isError && slots.length === 0
  const isFull = formattedSlots.length > 0 && formattedSlots.every((s) => s.status !== 'free')
  const allPast = slots.length > 0 && formattedSlots.length === 0

  return (
    <div className="animate-fade-up pb-16">
      <header className="flex items-center gap-3 px-5 pt-6">
        <button
          type="button"
          onClick={goServices}
          className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-slate-200 bg-white text-slate-600 shadow-sm transition active:scale-90"
        >
          <ArrowLeft className="h-5 w-5" />
        </button>
        <div className="min-w-0">
          <h1 className="text-lg font-bold leading-tight">Выбор времени</h1>
          <p className="truncate text-xs text-slate-400">
            {service?.name} · {service?.duration} мин · {rub(service?.price || 0)}
          </p>
        </div>
      </header>

      <section className="mt-6 px-5">
        <div className="mb-2.5 flex items-center gap-1.5 text-slate-500">
          <Calendar className="h-4 w-4" />
          <span className="text-sm font-semibold">Дата</span>
        </div>

        <div
          ref={daysScrollRef}
          className="thin-scrollbar -mx-5 flex cursor-grab snap-x snap-mandatory select-none gap-2 scroll-pl-5 overflow-x-auto px-5 pb-2 active:cursor-grabbing"
        >
          {QUICK_DAYS.map((d) => (
            <DayPill
              key={d.iso}
              day={d}
              selected={d.iso === selectedDate}
              onClick={() => handleSelect(d.iso)}
            />
          ))}
        </div>

        <div className="mb-3 mt-5 flex items-center justify-between">
          <NavArrow direction="prev" onClick={prevMonth} disabled={!canPrev} />
          <p className="font-bold">
            {MONTHS_TITLE[view.month]} {view.year}
          </p>
          <NavArrow direction="next" onClick={nextMonth} disabled={!canNext} />
        </div>

        <div className="grid grid-cols-7 gap-1.5 text-center text-[11px] font-semibold uppercase text-slate-400">
          {WEEKDAYS_GRID.map((w) => (
            <span key={w} className="py-1">
              {w}
            </span>
          ))}
        </div>

        <div key={`${view.year}-${view.month}`} className="mt-1.5 grid animate-fade-in grid-cols-7 gap-1.5">
          {monthGrid.map((d, i) =>
            d === null ? (
              <div key={`empty-${i}`} />
            ) : (
              <CalendarDay
                key={d.toISOString()}
                date={d}
                selected={d.getTime() === selectedDateObj.getTime()}
                isToday={d.getTime() === TODAY.getTime()}
                disabled={d < TODAY}
                onClick={() => handleSelect(toISO(d))}
              />
            ),
          )}
        </div>
      </section>

      <section className="mt-6 px-5">
        <div className="mb-3 flex items-center gap-1.5 text-slate-500">
          <Clock className="h-4 w-4" />
          <span className="text-sm font-semibold">Время</span>
        </div>

        {slotsQuery.isPending ? (
          <SlotsSkeleton />
        ) : slotsQuery.isError ? (
          <QueryRetry message="Не удалось загрузить слоты" onRetry={slotsQuery.refetch} />
        ) : allPast ? (
          <div className="animate-fade-in rounded-2xl border border-slate-100 bg-white p-6 text-center shadow-sm">
            <p className="text-lg font-bold text-slate-900">Все слоты на сегодня прошли</p>
            <p className="mt-1.5 text-sm text-slate-400">Выберите другой день для записи</p>
          </div>
        ) : isDayOff ? (
          <div className="animate-fade-in rounded-2xl border border-slate-100 bg-white p-6 text-center shadow-sm">
            <p className="text-lg font-bold text-slate-900">Мастер не работает в этот день</p>
            <p className="mt-1.5 text-sm text-slate-400">Выберите другую дату</p>
          </div>
        ) : isFull ? (
          <WaitlistBlock iso={selectedDate} />
        ) : (
          <div key={selectedDate} className="animate-fade-in">
            <div className="grid grid-cols-3 gap-2.5">
              {formattedSlots.map((slot) => (
                <SlotButton
                  key={slot.start_time}
                  time={slot.time}
                  status={mineTimes.has(slot.startMs) ? 'mine' : slot.status}
                  onClick={() => handleSlotClick(slot.start_time)}
                />
              ))}
            </div>
            <SlotsLegend />
          </div>
        )}
      </section>
    </div>
  )
}
