import { useState } from 'react'
import { ArrowLeft, Calendar, Clock } from 'lucide-react'
import { DayPill } from '../components/ui/DayPill'
import { CalendarDay } from '../components/ui/CalendarDay'
import { NavArrow } from '../components/ui/NavArrow'
import { SlotButton } from '../components/ui/SlotButton'
import { SlotsLegend } from '../components/ui/SlotsLegend'
import { WaitlistBlock } from '../widgets/WaitlistBlock'
import { useBookingStore } from '../store/bookingStore'
import { useDragScroll } from '../lib/useDragScroll'
import {
  TODAY, MONTHS_TITLE, WEEKDAYS_GRID,
  buildQuickDays, buildMonthGrid, fromISO, getScheduleForDate,
} from '../lib/dates'
import { rub } from '../lib/currency'

const QUICK_DAYS = buildQuickDays()

export function BookingPage({ onBack, onSlotClick }) {
  const service = useBookingStore((s) => s.service)
  const selectedDate = useBookingStore((s) => s.selectedDate)
  const setDate = useBookingStore((s) => s.setDate)
  const myBookings = useBookingStore((s) => s.myBookings)

  const daysScrollRef = useDragScroll()

  const selected = fromISO(selectedDate)
  const schedule = getScheduleForDate(selected)
  const isFull = schedule.every((st) => st !== 'free')

  const [view, setView] = useState({ year: TODAY.getFullYear(), month: TODAY.getMonth() })
  const diffMonths = (view.year - TODAY.getFullYear()) * 12 + (view.month - TODAY.getMonth())
  const canPrev = diffMonths > 0
  const canNext = diffMonths < 1

  const prevMonth = () =>
    setView((v) => (v.month === 0 ? { year: v.year - 1, month: 11 } : { ...v, month: v.month - 1 }))
  const nextMonth = () =>
    setView((v) => (v.month === 11 ? { year: v.year + 1, month: 0 } : { ...v, month: v.month + 1 }))

  const handleSelect = (iso) => {
    setDate(iso)
    const d = fromISO(iso)
    setView({ year: d.getFullYear(), month: d.getMonth() })
  }

  return (
    <div className="animate-fade-up pb-16">
      <header className="flex items-center gap-3 px-5 pt-6">
        <button
          onClick={onBack}
          className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-slate-200 bg-white text-slate-600 shadow-sm transition active:scale-90"
        >
          <ArrowLeft className="h-5 w-5" />
        </button>
        <div className="min-w-0">
          <h1 className="text-lg font-bold leading-tight">Выбор времени</h1>
          <p className="truncate text-xs text-slate-400">
            {service.name} · {service.duration} мин · {rub(service.price)}
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
          {WEEKDAYS_GRID.map((w) => <span key={w} className="py-1">{w}</span>)}
        </div>

        <div key={`${view.year}-${view.month}`} className="mt-1.5 grid animate-fade-in grid-cols-7 gap-1.5">
          {buildMonthGrid(view.year, view.month).map((d, i) =>
            d === null ? (
              <div key={`empty-${i}`} />
            ) : (
              <CalendarDay
                key={d.toISOString()}
                date={d}
                selected={d.getTime() === fromISO(selectedDate).getTime()}
                isToday={d.getTime() === TODAY.getTime()}
                disabled={d < TODAY}
                onClick={() => handleSelect(`${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`)}
              />
            )
          )}
        </div>
      </section>

      <section className="mt-6 px-5">
        <div className="mb-3 flex items-center gap-1.5 text-slate-500">
          <Clock className="h-4 w-4" />
          <span className="text-sm font-semibold">Время</span>
        </div>

        {isFull ? (
          <WaitlistBlock iso={selectedDate} />
        ) : (
          <div key={selectedDate} className="animate-fade-in">
            <div className="grid grid-cols-3 gap-2.5">
              {schedule.map((status, slotIndex) => {
                const key = `${selectedDate}-${slotIndex}`
                return (
                  <SlotButton
                    key={key}
                    time={['10:00', '11:30', '13:00', '14:30', '16:00'][slotIndex]}
                    status={myBookings[key] ? 'mine' : status}
                    onClick={() => onSlotClick(selectedDate, slotIndex)}
                  />
                )
              })}
            </div>
            <SlotsLegend />
          </div>
        )}
      </section>
    </div>
  )
}