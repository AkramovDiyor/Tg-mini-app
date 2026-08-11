import { User, MapPin } from 'lucide-react'
import { SheetShell } from './SheetShell'
import { useBookingStore } from '../../store/bookingStore'
import { getBookingTimeLabel } from '../../lib/dates'

export function BookingDetailsSheet({ onClose, onCancel }) {
  const activeBooking = useBookingStore((s) => s.activeBooking)
  if (!activeBooking) return null

  return (
    <SheetShell onClose={onClose}>
      <p className="text-center text-[11px] font-bold uppercase tracking-widest text-slate-400">
        Моя запись
      </p>
      <h3 className="mt-3 text-center text-2xl font-extrabold leading-snug text-slate-900">
        {activeBooking.service}
      </h3>
      <p className="mt-1 text-center text-2xl font-extrabold leading-snug text-emerald-600">
        {getBookingTimeLabel(activeBooking.date)}
      </p>

      <div className="mt-6 flex items-center gap-3 rounded-2xl border border-slate-100 bg-slate-50 p-4 text-left">
        <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-emerald-400 to-teal-600 text-white">
          <User className="h-6 w-6" />
        </div>
        <div className="min-w-0 flex-1">
          <p className="font-bold text-slate-900">Педро Барбер</p>
          <p className="mt-0.5 flex items-center gap-1 text-xs text-slate-400">
            <MapPin className="h-3.5 w-3.5 shrink-0" />
            ул. Центральная, 1
          </p>
        </div>
      </div>

      <button
        onClick={onCancel}
        className="mt-6 w-full rounded-xl bg-red-50 py-4 font-bold text-red-600 transition active:scale-[0.98]"
      >
        Отменить запись
      </button>
    </SheetShell>
  )
}