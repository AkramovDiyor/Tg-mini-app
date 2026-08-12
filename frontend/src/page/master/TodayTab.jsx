import { TrendingUp } from 'lucide-react'
import { rub } from '../../lib/currency'
import { TODAY_APPOINTMENTS, DOT_COLORS } from '../../mock/masterData'

export function TodayTab() {
  const today = new Date()
  const dateLabel = today.toLocaleDateString('ru-RU', {
    day: 'numeric', month: 'long',
  })

  const total = TODAY_APPOINTMENTS.reduce(
    (sum, a) => sum + (a.price || 0), 0,
  )
  const booked = TODAY_APPOINTMENTS.filter((a) => a.status !== 'free').length

  return (
    <div>
      <div className="mb-5">
        <p className="text-xs font-semibold uppercase tracking-widest text-emerald-600">
          Сегодня
        </p>
        <h1 className="mt-1 text-2xl font-extrabold text-slate-900 capitalize">
          Сегодня, {dateLabel}
        </h1>
      </div>

      <div className="mb-6 flex items-center justify-between rounded-3xl bg-slate-900 p-5 text-white shadow-xl shadow-slate-900/20">
        <div>
          <p className="text-[11px] font-semibold uppercase tracking-widest text-slate-400">
            Ожидаемая выручка
          </p>
          <p className="mt-1.5 text-2xl font-extrabold">{rub(total)}</p>
          <p className="mt-0.5 text-xs text-slate-400">
            {booked} записей запланировано
          </p>
        </div>
        <div className="flex h-12 w-12 items-center justify-center rounded-full bg-white/10 backdrop-blur">
          <TrendingUp className="h-5 w-5 text-emerald-400" />
        </div>
      </div>

      <div className="relative">
        <div className="absolute left-3 top-3 bottom-3 w-0.5 bg-gradient-to-b from-slate-200 via-slate-200 to-transparent" />

        <div className="space-y-3">
          {TODAY_APPOINTMENTS.map((appt) => (
            <div key={appt.time} className="relative pl-9">
              <span
                className={`absolute left-1.5 top-5 h-3 w-3 rounded-full border-[3px] border-[#F9FAFB] ${DOT_COLORS[appt.status]}`}
              />

              {appt.status === 'free' ? (
                <div className="rounded-2xl border-2 border-dashed border-slate-200 bg-white/60 p-4">
                  <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
                    {appt.time}
                  </p>
                  <p className="mt-1 text-sm font-medium text-slate-400">
                    Свободное окно
                  </p>
                </div>
              ) : (
                <div className="rounded-2xl bg-white p-4 shadow-sm">
                  <div className="flex items-start justify-between gap-3">
                    <div className="min-w-0">
                      <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
                        {appt.time}
                      </p>
                      <p className="mt-1 text-[15px] font-bold text-slate-900">
                        {appt.client}
                      </p>
                      <p className="mt-0.5 text-sm text-slate-400">{appt.service}</p>
                    </div>
                    <p className="shrink-0 text-[15px] font-bold text-slate-900">
                      {rub(appt.price)}
                    </p>
                  </div>

                  <div className="mt-3 flex items-center gap-1.5">
                    <i className={`h-1.5 w-1.5 rounded-full ${DOT_COLORS[appt.status]}`} />
                    <span
                      className={`text-[11px] font-semibold ${
                        appt.status === 'confirmed' ? 'text-emerald-700' : 'text-amber-700'
                      }`}
                    >
                      {appt.status === 'confirmed' ? 'Подтверждено' : 'Ждёт подтверждения'}
                    </span>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}