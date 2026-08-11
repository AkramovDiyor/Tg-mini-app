import { Check, Clock, Bell } from 'lucide-react'
import { SheetShell } from './SheetShell'
import { useBookingStore } from '../../store/bookingStore'
import { fromISO, MONTHS_GEN, SLOT_TIMES } from '../../lib/dates'
import { rub } from '../../lib/currency'

export function ConfirmBookingSheet({ slot, confirming, onClose, onConfirm }) {
  const service = useBookingStore((s) => s.service)
  const date = fromISO(slot.iso)
  const time = SLOT_TIMES[slot.slotIndex]

  return (
    <SheetShell onClose={onClose} closableOnBackdrop={!confirming}>
      {confirming ? (
        <div className="flex animate-pop-in flex-col items-center py-10 text-center">
          <div className="mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-emerald-100 text-emerald-600">
            <Check className="h-8 w-8" strokeWidth={3} />
          </div>
          <p className="text-lg font-bold">Запись подтверждена!</p>
          <p className="mt-1 text-sm text-slate-400">Напомним за 2 часа до начала</p>
        </div>
      ) : (
        <>
          <p className="text-center text-[11px] font-bold uppercase tracking-widest text-slate-400">
            Подтверждение записи
          </p>
          <h3 className="mt-3 text-center text-2xl font-extrabold leading-snug">
            {service.name},<br />
            {date.getDate()} {MONTHS_GEN[date.getMonth()]}, {time}
          </h3>
          <div className="mt-4 flex justify-center gap-2">
            <span className="inline-flex items-center gap-1 rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold text-slate-600">
              <Clock className="h-3.5 w-3.5" />
              {service.duration} мин
            </span>
            <span className="inline-flex items-center rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold text-emerald-600">
              {rub(service.price)}
            </span>
          </div>
          <button
            onClick={onConfirm}
            className="mt-6 flex w-full items-center justify-center gap-2 rounded-xl bg-emerald-600 py-4 text-[15px] font-bold text-white shadow-lg shadow-emerald-600/25 transition active:scale-[0.98]"
          >
            <Check className="h-5 w-5" strokeWidth={3} />
            Подтвердить запись
          </button>
          <p className="mt-3 flex items-start gap-1.5 px-2 text-[11px] leading-relaxed text-slate-400">
            <Bell className="mt-0.5 h-3.5 w-3.5 shrink-0" />
            За 2 часа до начала придет кнопка подтверждения. Если не нажать — запись может быть передана другому.
          </p>
        </>
      )}
    </SheetShell>
  )
}