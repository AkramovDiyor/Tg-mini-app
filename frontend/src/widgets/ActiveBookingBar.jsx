import { useEffect, useState } from 'react'
import { Scissors, ChevronRight } from 'lucide-react'
import { useBookingStore } from '../store/bookingStore'
import { getBookingTimeLabel } from '../lib/dates'

export function ActiveBookingBar({ onOpenDetails }) {
  const activeBooking = useBookingStore((s) => s.activeBooking)
  const [now, setNow] = useState(() => new Date())

  useEffect(() => {
    const t = setInterval(() => setNow(new Date()), 60000)
    return () => clearInterval(t)
  }, [])

  if (!activeBooking) return null

  return (
    <button
      onClick={onOpenDetails}
      className="fixed bottom-4 left-1/2 z-40 w-[calc(100%-2.5rem)] max-w-[380px] -translate-x-1/2 animate-fade-up"
    >
      <div className="flex items-center justify-between rounded-2xl bg-slate-900 p-3.5 text-white shadow-2xl">
        <span className="flex shrink-0 items-center justify-center rounded-full bg-white/10 p-2">
          <Scissors className="h-5 w-5 text-emerald-400" />
        </span>
        <span className="mx-3 flex min-w-0 flex-1 flex-col items-start text-left">
          <span className="text-xs font-bold text-emerald-400">
            {getBookingTimeLabel(activeBooking.date, now)}
          </span>
          <span className="w-full truncate text-sm font-semibold">
            {activeBooking.service}, {activeBooking.time}
          </span>
        </span>
        <ChevronRight className="h-5 w-5 shrink-0 text-slate-400" />
      </div>
    </button>
  )
}