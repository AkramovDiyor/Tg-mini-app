import { useEffect, useState } from 'react'
import { Clock } from 'lucide-react'
import { useBookingStore } from '../../store/bookingStore'
import { fetchWaitlist } from '../../services/api'

const AVATAR_GRADIENTS = [
  'from-amber-400 to-rose-500',
  'from-sky-400 to-indigo-600',
  'from-emerald-400 to-teal-600',
]

export function QueueTab() {
  const showToast = useBookingStore((s) => s.showToast)
  const [waitlist, setWaitlist] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    loadWaitlist()
  }, [])

  const loadWaitlist = async () => {
    try {
      setLoading(true)
      setError(null)
      const data = await fetchWaitlist()
      setWaitlist(Array.isArray(data) ? data : [])
    } catch (err) {
      console.error('Failed to load waitlist:', err)
      setError('Не удалось загрузить лист ожидания')
    } finally {
      setLoading(false)
    }
  }

  const handleOffer = (name) => {
    showToast(`Предложение отправлено ${name}. Он получит уведомление в Telegram 🔔`)
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <p className="text-sm text-slate-400">Загрузка листа ожидания...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center py-20">
        <p className="text-sm text-red-500">{error}</p>
        <button
          onClick={loadWaitlist}
          className="mt-4 rounded-xl bg-emerald-500 px-6 py-2.5 text-sm font-bold text-white transition active:scale-95"
        >
          Повторить
        </button>
      </div>
    )
  }

  return (
    <div>
      <div className="mb-5">
        <p className="text-xs font-semibold uppercase tracking-widest text-emerald-600">
          Очередь
        </p>
        <h1 className="mt-1 text-2xl font-extrabold text-slate-900">Лист ожидания</h1>
        <p className="mt-1 text-sm text-slate-400">
          {waitlist.length} {waitlist.length === 1 ? 'клиент готов' : 'клиентов готовы'} приехать
        </p>
      </div>

      {waitlist.length === 0 ? (
        <div className="rounded-2xl border-2 border-dashed border-slate-200 bg-white/60 p-6 text-center">
          <p className="text-sm text-slate-400">Лист ожидания пуст</p>
        </div>
      ) : (
        <div className="space-y-3">
          {waitlist.map((item, idx) => {
            const initials = item.client_name
              .split(' ')
              .map((n) => n[0])
              .join('')
              .toUpperCase()
              .slice(0, 2)

            const createdAt = new Date(item.created_at)
            const createdAtStr = `${String(createdAt.getHours()).padStart(2, '0')}:${String(
              createdAt.getMinutes()
            ).padStart(2, '0')}`

            return (
              <div key={item.id} className="rounded-2xl bg-white p-4 shadow-sm">
                <div className="flex items-center gap-3">
                  <div
                    className={`flex h-11 w-11 shrink-0 items-center justify-center rounded-full bg-gradient-to-br ${
                      AVATAR_GRADIENTS[idx % AVATAR_GRADIENTS.length]
                    } text-sm font-bold text-white shadow-md`}
                  >
                    {initials}
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="font-bold text-slate-900">{item.client_name}</p>
                    <p className="mt-0.5 flex items-center gap-1 text-sm text-slate-400">
                      <Clock className="h-3.5 w-3.5" />
                      Хочет {item.preferred_date} · встал в {createdAtStr}
                    </p>
                  </div>
                </div>
                <button
                  onClick={() => handleOffer(item.client_name)}
                  className="mt-3 w-full rounded-xl bg-emerald-50 py-2.5 text-sm font-bold text-emerald-700 transition active:scale-[0.98]"
                >
                  Предложить это окно
                </button>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}