import { useMemo, useState, useEffect } from 'react'
import { Scissors, ChevronRight } from 'lucide-react'
import { useClientBookingsQuery } from '../hooks/useClientQueries'
import { formatTimeInMasterTz } from '../lib/dates'

const CANCELLED = new Set([
  'cancelled',
  'canceled',
  'rejected',
  'declined',
  'cancelled_by_client',
  'cancelled_no_show',
])

export function ActiveBookingBar({ onOpenDetails }) {
  const { data: bookings = [], isPending } = useClientBookingsQuery()
  const [now, setNow] = useState(() => Date.now())

  useEffect(() => {
    const t = setInterval(() => setNow(Date.now()), 60_000)
    return () => clearInterval(t)
  }, [])

  const activeBooking = useMemo(() => {
    return bookings
      .filter((b) => {
        if (CANCELLED.has(b.status?.toLowerCase())) return false
        return new Date(b.start_time).getTime() > now
      })
      .sort((a, b) => new Date(a.start_time) - new Date(b.start_time))[0]
  }, [bookings, now])

  if (isPending || !activeBooking) return null

  return (
    <button
      type="button"
      onClick={() => onOpenDetails(activeBooking)}
      className="fixed bottom-4 left-1/2 z-40 w-[calc(100%-2.5rem)] max-w-[380px] -translate-x-1/2 animate-fade-up"
    >
      <div className="flex items-center justify-between rounded-2xl bg-slate-900 p-3.5 text-white shadow-2xl">
        <span className="flex shrink-0 items-center justify-center rounded-full bg-white/10 p-2">
          <Scissors className="h-5 w-5 text-emerald-400" />
        </span>
        <span className="mx-3 flex min-w-0 flex-1 flex-col items-start text-left">
          <span className="text-xs font-bold text-emerald-400">
            {formatTimeInMasterTz(activeBooking.start_time)}
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
