import { Frown, CheckCircle2 } from 'lucide-react'
import { useBookingStore } from '../store/bookingStore'

export function WaitlistBlock({ iso }) {
  const joined = useBookingStore((s) => !!s.waitlist[iso])
  const joinWaitlist = useBookingStore((s) => s.joinWaitlist)
  const showToast = useBookingStore((s) => s.showToast)

  const handleJoin = () => {
    joinWaitlist(iso)
    showToast('Ты в листе ожидания! Бот напишет, как только появится окно 🔔')
  }

  return (
    <div className="animate-fade-in rounded-2xl border border-slate-100 bg-white p-6 shadow-card">
      <div className="flex flex-col items-center text-center">
        <div className="mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-slate-100">
          <Frown className="h-8 w-8 text-slate-400" />
        </div>
        <p className="text-lg font-bold">На этот день всё занято</p>
        <p className="mt-1.5 max-w-[260px] text-sm leading-relaxed text-slate-400">
          Оставь заявку — если место освободится, бот сразу напишет тебе
        </p>
      </div>

      {joined ? (
        <div className="mt-5 flex animate-pop-in items-start gap-2.5 rounded-xl border border-emerald-200 bg-emerald-50 px-4 py-3.5 text-sm font-medium leading-snug text-emerald-800">
          <CheckCircle2 className="mt-0.5 h-5 w-5 shrink-0 text-emerald-600" />
          Ты в списке! Если кто-то отменит, бот пришлет тебе сообщение
        </div>
      ) : (
        <button
          onClick={handleJoin}
          className="mt-5 w-full rounded-xl bg-emerald-600 px-4 py-4 text-[15px] font-bold leading-snug text-white shadow-lg shadow-emerald-600/25 transition active:scale-[0.98]"
        >
          Встать в лист ожидания на этот день
        </button>
      )}
    </div>
  )
}