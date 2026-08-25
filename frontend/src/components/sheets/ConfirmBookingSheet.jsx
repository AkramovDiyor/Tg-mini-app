import { Check, Clock, Bell, Loader2 } from 'lucide-react'
import { SheetShell } from './SheetShell'
import { useBookingStore } from '../../store/bookingStore'
import { fromISO, MONTHS_GEN, formatTimeInMasterTz } from '../../lib/dates'
import { rub } from '../../lib/currency'

export function ConfirmBookingSheet({ slot, confirming, onClose, onConfirm }) {
  const service = useBookingStore((s) => s.service)
  const date = fromISO(slot.iso)
  const timeStr = formatTimeInMasterTz(slot.startTime)

  return (
    <SheetShell onClose={onClose} closableOnBackdrop={!confirming}>
      {confirming ? (
        <div className="flex flex-col items-center py-10 text-center">
          <Loader2 className="mb-4 h-10 w-10 animate-spin text-emerald-600" />
          <p className="text-lg font-bold">Бронируем слот…</p>
          <p className="mt-1 text-sm text-slate-400">Не закрывайте приложение</p>
        </div>
      ) : (
        <>
          <p className="text-center text-[11px] font-bold uppercase tracking-widest text-slate-400">
            Подтверждение записи
          </p>
          <h3 className="mt-3 text-center text-2xl font-extrabold leading-snug">
            {service.name},<br />
            {date.getDate()} {MONTHS_GEN[date.getMonth()]}, {timeStr}
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
            type="button"
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
