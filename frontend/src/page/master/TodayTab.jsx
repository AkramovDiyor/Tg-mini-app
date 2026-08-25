import { useMemo } from 'react'
import { TrendingUp } from 'lucide-react'
import { QueryRetry } from '../../components/AppFrame'
import { MasterListSkeleton } from '../../components/skeletons/Skeletons'
import { useTodayScheduleQuery, useMasterServicesQuery } from '../../hooks/useMasterQueries'
import { formatTimeInMasterTz } from '../../lib/dates'
import { rub } from '../../lib/currency'

const STATUS_CONFIG = {
  active: { dot: 'bg-emerald-500', label: 'Подтверждено', text: 'text-emerald-700' },
  confirmed: { dot: 'bg-emerald-500', label: 'Подтверждено', text: 'text-emerald-700' },
  pending: { dot: 'bg-amber-400', label: 'Ждёт подтверждения', text: 'text-amber-700' },
}

export function TodayTab() {
  const scheduleQuery = useTodayScheduleQuery()
  const servicesQuery = useMasterServicesQuery()

  const servicesMap = useMemo(() => {
    const map = new Map()
    for (const s of servicesQuery.data || []) map.set(s.id, s.name)
    return map
  }, [servicesQuery.data])

  const today = new Date()
  const dateLabel = today.toLocaleDateString('ru-RU', { day: 'numeric', month: 'long' })

  if (scheduleQuery.isPending) return <MasterListSkeleton />

  if (scheduleQuery.isError) {
    return (
      <QueryRetry
        message="Не удалось загрузить расписание"
        onRetry={scheduleQuery.refetch}
      />
    )
  }

  const stats = scheduleQuery.data?.stats || { count: 0, total: 0 }
  const schedule = scheduleQuery.data?.schedule || []

  return (
    <div>
      <div className="mb-5">
        <p className="text-xs font-semibold uppercase tracking-widest text-emerald-600">Сегодня</p>
        <h1 className="mt-1 text-2xl font-extrabold text-slate-900 capitalize">
          Сегодня, {dateLabel}
        </h1>
      </div>

      <div className="mb-6 flex items-center justify-between rounded-3xl bg-slate-900 p-5 text-white shadow-xl shadow-slate-900/20">
        <div>
          <p className="text-[11px] font-semibold uppercase tracking-widest text-slate-400">
            Ожидаемая выручка
          </p>
          <p className="mt-1.5 text-2xl font-extrabold">{rub(stats.total)}</p>
          <p className="mt-0.5 text-xs text-slate-400">{stats.count} записей запланировано</p>
        </div>
        <div className="flex h-12 w-12 items-center justify-center rounded-full bg-white/10 backdrop-blur">
          <TrendingUp className="h-5 w-5 text-emerald-400" />
        </div>
      </div>

      <div className="relative">
        <div className="absolute left-3 top-3 bottom-3 w-0.5 bg-gradient-to-b from-slate-200 via-slate-200 to-transparent" />

        <div className="space-y-3">
          {schedule.length === 0 ? (
            <div className="rounded-2xl border-2 border-dashed border-slate-200 bg-white/60 p-6 text-center">
              <p className="text-sm text-slate-400">На сегодня записей нет</p>
            </div>
          ) : (
            schedule.map((appt) => {
              const timeStr = formatTimeInMasterTz(appt.start_time)
              const statusCfg = STATUS_CONFIG[appt.status] || STATUS_CONFIG.pending
              const serviceName = servicesMap.get(appt.service_id) || 'Услуга'
              const price = appt.price_locked ?? appt.service_price ?? appt.price ?? 0

              return (
                <div key={appt.id} className="relative pl-9">
                  <span
                    className={`absolute left-1.5 top-5 h-3 w-3 rounded-full border-[3px] border-[#F9FAFB] ${statusCfg.dot}`}
                  />
                  <div className="rounded-2xl bg-white p-4 shadow-sm">
                    <div className="flex items-start justify-between gap-3">
                      <div className="min-w-0">
                        <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
                          {timeStr}
                        </p>
                        <p className="mt-1 text-[15px] font-bold text-slate-900">{appt.client_name}</p>
                        <p className="mt-0.5 text-sm text-slate-400">{serviceName}</p>
                      </div>
                      <p className="shrink-0 text-[15px] font-bold text-slate-900">{rub(price)}</p>
                    </div>
                    <div className="mt-3 flex items-center gap-1.5">
                      <i className={`h-1.5 w-1.5 rounded-full ${statusCfg.dot}`} />
                      <span className={`text-[11px] font-semibold ${statusCfg.text}`}>
                        {statusCfg.label}
                      </span>
                    </div>
                  </div>
                </div>
              )
            })
          )}
        </div>
      </div>
    </div>
  )
}
