import { Clock } from 'lucide-react'
import { useBookingStore } from '../../store/bookingStore'
import { WAITLIST, AVATAR_GRADIENTS } from '../../mock/masterData'

export function QueueTab() {
  const showToast = useBookingStore((s) => s.showToast)

  const handleOffer = (name) => {
    showToast(`Предложение отправлено ${name}. Он получит уведомление в Telegram 🔔`)
  }

  return (
    <div>
      <div className="mb-5">
        <p className="text-xs font-semibold uppercase tracking-widest text-emerald-600">
          Очередь
        </p>
        <h1 className="mt-1 text-2xl font-extrabold text-slate-900">Лист ожидания</h1>
        <p className="mt-1 text-sm text-slate-400">
          {WAITLIST.length} клиента готовы приехать
        </p>
      </div>

      <div className="space-y-3">
        {WAITLIST.map((item, idx) => (
          <div key={item.id} className="rounded-2xl bg-white p-4 shadow-sm">
            <div className="flex items-center gap-3">
              <div
                className={`flex h-11 w-11 shrink-0 items-center justify-center rounded-full bg-gradient-to-br ${
                  AVATAR_GRADIENTS[idx % AVATAR_GRADIENTS.length]
                } text-sm font-bold text-white shadow-md`}
              >
                {item.initials}
              </div>
              <div className="min-w-0 flex-1">
                <p className="font-bold text-slate-900">{item.name}</p>
                <p className="mt-0.5 flex items-center gap-1 text-sm text-slate-400">
                  <Clock className="h-3.5 w-3.5" />
                  Хочет в {item.time} · встал в {item.joinedAt}
                </p>
              </div>
            </div>
            <button
              onClick={() => handleOffer(item.name)}
              className="mt-3 w-full rounded-xl bg-emerald-50 py-2.5 text-sm font-bold text-emerald-700 transition active:scale-[0.98]"
            >
              Предложить это окно
            </button>
          </div>
        ))}
      </div>
    </div>
  )
}