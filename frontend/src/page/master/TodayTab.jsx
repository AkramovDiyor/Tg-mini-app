import { useEffect, useState } from 'react'
import { TrendingUp } from 'lucide-react'
import { rub } from '../../lib/currency'
import { fetchTodaySchedule, fetchMasterServices } from '../../services/api'

const STATUS_CONFIG = {
  active:    { dot: 'bg-emerald-500', label: 'Подтверждено',       text: 'text-emerald-700' },
  confirmed: { dot: 'bg-emerald-500', label: 'Подтверждено',       text: 'text-emerald-700' },
  pending:   { dot: 'bg-amber-400',   label: 'Ждёт подтверждения', text: 'text-amber-700'   },
}

const MASTER_TIME_ZONE = 'Europe/Moscow'

export function TodayTab() {
  const [data, setData] = useState(null)
  const [servicesMap, setServicesMap] = useState(new Map())
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    loadTodaySchedule()
  }, [])

  const loadTodaySchedule = async () => {
    try {
      setLoading(true)
      setError(null)

      const [scheduleResult, servicesResult] = await Promise.all([
        fetchTodaySchedule(),
        fetchMasterServices().catch(() => []),
      ])

      setData(scheduleResult)

      const map = new Map()
      if (Array.isArray(servicesResult)) {
        servicesResult.forEach((s) => map.set(s.id, s.name))
      }
      setServicesMap(map)
    } catch (err) {
      console.error('Failed to load today schedule:', err)
      setError('Не удалось загрузить расписание')
    } finally {
      setLoading(false)
    }
  }

  const today = new Date()
  const dateLabel = today.toLocaleDateString('ru-RU', {
    day: 'numeric', month: 'long',
  })

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <p className="text-sm text-slate-400">Загрузка расписания...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center py-20">
        <p className="text-sm text-red-500">{error}</p>
        <button
          onClick={loadTodaySchedule}
          className="mt-4 rounded-xl bg-emerald-500 px-6 py-2.5 text-sm font-bold text-white transition active:scale-95"
        >
          Повторить
        </button>
      </div>
    )
  }

  const stats = data?.stats || { count: 0, total: 0 }
  const schedule = data?.schedule || []

  // ✅ Функция форматирования — вынесена отдельно для чистоты
  const formatTimeInMasterTz = (isoString) => {
    const date = new Date(isoString)
    console.log('Browser TZ:', Intl.DateTimeFormat().resolvedOptions().timeZone)
    console.log('Using TZ:', MASTER_TIME_ZONE)
    return date.toLocaleTimeString('ru-RU', {
      hour: '2-digit',
      minute: '2-digit',
      timeZone: MASTER_TIME_ZONE,
    })
  }
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
          <p className="mt-1.5 text-2xl font-extrabold">{rub(stats.total)}</p>
          <p className="mt-0.5 text-xs text-slate-400">
            {stats.count} записей запланировано
          </p>
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
              // ✅ Используем функцию — больше нет ошибки с startTime
              const timeStr = formatTimeInMasterTz(appt.start_time)
              const statusCfg = STATUS_CONFIG[appt.status] || STATUS_CONFIG.pending
              const serviceName = servicesMap.get(appt.service_id) || 'Услуга'
              const price = appt.price_locked ?? appt.service_price ?? appt.price ?? 0

              return (
                <div key={appt.id} className="relative pl-9">
                  <span className={`absolute left-1.5 top-5 h-3 w-3 rounded-full border-[3px] border-[#F9FAFB] ${statusCfg.dot}`} />

                  <div className="rounded-2xl bg-white p-4 shadow-sm">
                    <div className="flex items-start justify-between gap-3">
                      <div className="min-w-0">
                        <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
                          {timeStr}
                        </p>
                        <p className="mt-1 text-[15px] font-bold text-slate-900">
                          {appt.client_name}
                        </p>
                        <p className="mt-0.5 text-sm text-slate-400">{serviceName}</p>
                      </div>
                      <p className="shrink-0 text-[15px] font-bold text-slate-900">
                        {rub(price)}
                      </p>
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