import { useState } from 'react'
import {
  MapPin, Clock, Plus, Bell, Moon, LogOut, User, Pencil, ChevronRight,
} from 'lucide-react'
import { useBookingStore } from '../../store/bookingStore'
import { rub } from '../../lib/currency'
import { SERVICES_LIST } from '../../mock/data'
import { WEEK_TAGS, CANCEL_HOURS } from '../../mock/masterData'
import { Toggle } from '../../components/ui/Toggle'
import { AddServiceSheet } from '../../components/sheets/AddServiceSheet'
import { EditProfileSheet } from '../../components/sheets/EditProfileSheet'

const APP_SETTINGS = [
  { id: 'notifications', label: 'Уведомления',       icon: Bell,   kind: 'toggle' },
  { id: 'theme',         label: 'Тёмная тема',       icon: Moon,   kind: 'toggle' },
  { id: 'logout',        label: 'Выйти из аккаунта', icon: LogOut, kind: 'danger' },
]

export function ProfileTab() {
  const showToast = useBookingStore((s) => s.showToast)

  const [workDays, setWorkDays] = useState([true, true, true, true, true, false, false])
  const [hours, setHours] = useState({
    start: '09:00', end: '20:00', lunchFrom: '13:00', lunchTo: '14:00',
  })

  const [autoCancel, setAutoCancel] = useState(true)
  const [cancelHours, setCancelHours] = useState('2')
  const [offerWaitlist, setOfferWaitlist] = useState(true)

  const [notifications, setNotifications] = useState(true)
  const [darkTheme, setDarkTheme] = useState(false)

  const [isAddServiceOpen, setIsAddServiceOpen] = useState(false)
  const [isEditProfileOpen, setIsEditProfileOpen] = useState(false)

  const [link] = useState(
    () => `t.me/antizabivator/invite/${Math.random().toString(36).slice(2, 10)}`,
  )

  const toggleDay = (i) =>
    setWorkDays((days) => days.map((d, idx) => (idx === i ? !d : d)))

  const handleCopyLink = () => {
    try {
      navigator.clipboard?.writeText(link)
      showToast('Ссылка скопирована в буфер обмена')
    } catch {
      showToast('Не удалось скопировать. Скопируйте вручную')
    }
  }

  const handleLogout = () => {
    showToast('Выход выполнен. До встречи, Педро 👋')
  }

  return (
    <div className="pb-24">
      {/* ===== ШАПКА ПРОФИЛЯ (кликабельный блок) ===== */}
      <button
        onClick={() => setIsEditProfileOpen(true)}
        className="mb-6 flex w-full items-center gap-4 rounded-2xl text-left transition active:scale-[0.98]"
      >
        <div className="relative shrink-0">
          <div className="flex h-20 w-20 items-center justify-center rounded-full bg-slate-200 text-slate-400">
            <User className="h-10 w-10" strokeWidth={1.5} />
          </div>
          {/* Убрали кнопку Camera */}
        </div>

        <div className="min-w-0 flex-1">
          <h1 className="text-xl font-extrabold text-slate-900">Педро Барбер</h1>
          <p className="mt-0.5 text-sm text-slate-400">Барбер · 6 лет опыта</p>
          <div className="mt-1.5 flex items-center gap-1">
            <span className="h-2 w-2 rounded-full bg-emerald-500" />
            <span className="text-xs font-semibold text-emerald-600">Онлайн</span>
          </div>
        </div>

        {/* Индикатор кликабельности */}
        <Pencil className="h-5 w-5 shrink-0 text-slate-300" />
      </button>

      {/* ===== КАРТОЧКА АДРЕСА (кликабельный блок) ===== */}
      <button
        onClick={() => setIsEditProfileOpen(true)}
        className="mb-6 flex w-full items-center gap-3 rounded-2xl bg-white p-4 text-left shadow-sm transition active:scale-[0.98]"
      >
        <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-emerald-50 text-emerald-600">
          <MapPin className="h-5 w-5" />
        </div>
        <div className="min-w-0 flex-1">
          <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">
            Адрес студии
          </p>
          <p className="mt-0.5 truncate font-bold text-slate-900">ул. Центральная, 1</p>
        </div>
        {/* Убрали текстовую кнопку "Изменить", поставили ChevronRight */}
        <ChevronRight className="h-5 w-5 shrink-0 text-slate-300" />
      </button>

      <div className="space-y-6">
        {/* ===== УСЛУГИ ===== */}
        <section>
          <h2 className="mb-2 text-base font-bold text-slate-800">Услуги</h2>
          <div className="rounded-2xl bg-white p-4 shadow-sm">
            <div className="space-y-2">
              {SERVICES_LIST.map((service) => {
                const Icon = service.icon
                return (
                  <button
                    key={service.id}
                    className="flex w-full items-center gap-3 rounded-xl bg-slate-50 p-3 text-left transition active:scale-[0.98]"
                  >
                    <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-emerald-50 text-emerald-600">
                      <Icon className="h-5 w-5" />
                    </span>
                    <span className="min-w-0 flex-1">
                      <span className="block text-sm font-bold text-slate-900">{service.name}</span>
                      <span className="mt-0.5 flex items-center gap-1 text-xs text-slate-400">
                        <Clock className="h-3 w-3" />
                        {service.duration} мин
                      </span>
                    </span>
                    <span className="text-sm font-bold text-slate-900">{rub(service.price)}</span>
                  </button>
                )
              })}
            </div>
            <button
              onClick={() => setIsAddServiceOpen(true)}
              className="mt-3 flex w-full items-center justify-center gap-2 rounded-xl border-2 border-dashed border-slate-200 py-3 font-semibold text-slate-500 transition active:scale-[0.98]"
            >
              <Plus className="h-4 w-4" />
              Добавить услугу
            </button>
          </div>
        </section>

        {/* ===== ГРАФИК РАБОТЫ ===== */}
        <section>
          <h2 className="mb-2 text-base font-bold text-slate-800">График работы</h2>
          <div className="rounded-2xl bg-white p-4 shadow-sm">
            <div className="flex flex-wrap gap-2">
              {WEEK_TAGS.map((tag, i) => {
                const isWeekend = i >= 7
                const isWork = workDays[i]
                return (
                  <button
                    key={tag}
                    disabled={isWeekend}
                    onClick={() => toggleDay(i)}
                    className={`h-10 rounded-full px-4 text-sm font-bold transition ${
                      isWeekend
                        ? 'cursor-not-allowed bg-slate-100 text-slate-300'
                        : isWork
                          ? 'bg-emerald-500 text-white shadow-md shadow-emerald-500/25 active:scale-95'
                          : 'bg-slate-100 text-slate-400 active:scale-95'
                    }`}
                  >
                    {tag}
                  </button>
                )
              })}
            </div>

            <div className="mt-4 grid grid-cols-2 gap-3">
              <label className="block">
                <span className="mb-1.5 block text-xs font-semibold text-slate-400">Начало дня</span>
                <input
                  type="time"
                  value={hours.start}
                  onChange={(e) => setHours((h) => ({ ...h, start: e.target.value }))}
                  className="w-full rounded-xl border-0 bg-slate-100 px-4 py-3 text-center font-semibold text-slate-900 outline-none transition focus:bg-emerald-50 focus:ring-2 focus:ring-emerald-500"
                />
              </label>
              <label className="block">
                <span className="mb-1.5 block text-xs font-semibold text-slate-400">Конец дня</span>
                <input
                  type="time"
                  value={hours.end}
                  onChange={(e) => setHours((h) => ({ ...h, end: e.target.value }))}
                  className="w-full rounded-xl border-0 bg-slate-100 px-4 py-3 text-center font-semibold text-slate-900 outline-none transition focus:bg-emerald-50 focus:ring-2 focus:ring-emerald-500"
                />
              </label>
            </div>

            <div className="mt-3">
              <span className="mb-1.5 block text-xs font-semibold text-slate-400">Обед</span>
              <div className="flex items-center gap-3">
                <input
                  type="time"
                  value={hours.lunchFrom}
                  onChange={(e) => setHours((h) => ({ ...h, lunchFrom: e.target.value }))}
                  className="flex-1 rounded-xl border-0 bg-slate-100 px-4 py-3 text-center font-semibold text-slate-900 outline-none transition focus:bg-emerald-50 focus:ring-2 focus:ring-emerald-500"
                />
                <span className="text-slate-300">—</span>
                <input
                  type="time"
                  value={hours.lunchTo}
                  onChange={(e) => setHours((h) => ({ ...h, lunchTo: e.target.value }))}
                  className="flex-1 rounded-xl border-0 bg-slate-100 px-4 py-3 text-center font-semibold text-slate-900 outline-none transition focus:bg-emerald-50 focus:ring-2 focus:ring-emerald-500"
                />
              </div>
            </div>
          </div>
        </section>

        {/* ===== ПРАВИЛА ЗАПИСИ ===== */}
        <section>
          <h2 className="mb-2 text-base font-bold text-slate-800">Правила записи</h2>
          <div className="rounded-2xl bg-white p-4 shadow-sm">
            <div className="flex items-center justify-between gap-3">
              <p className="text-sm font-semibold text-slate-900">Авто-отмена без подтверждения</p>
              <Toggle checked={autoCancel} onChange={setAutoCancel} />
            </div>

            {autoCancel && (
              <div className="mt-3 animate-fade-in rounded-xl bg-slate-50 p-3">
                <p className="text-xs leading-relaxed text-slate-500">
                  Отменять запись, если клиент не подтвердил её за{' '}
                  <b className="text-slate-700">{cancelHours} {cancelHours === '1' ? 'час' : 'часа'}</b>
                </p>
                <div className="mt-2 flex gap-2">
                  {CANCEL_HOURS.map((h) => (
                    <button
                      key={h}
                      onClick={() => setCancelHours(h)}
                      className={`flex-1 rounded-xl py-2 text-xs font-bold transition ${
                        cancelHours === h
                          ? 'bg-emerald-600 text-white shadow-md shadow-emerald-600/20'
                          : 'bg-white text-slate-500 active:scale-95'
                      }`}
                    >
                      {h} {h === '1' ? 'час' : 'часа'}
                    </button>
                  ))}
                </div>
              </div>
            )}

            <div className="mt-4 border-t border-slate-100 pt-4">
              <div className="flex items-center justify-between gap-3">
                <div className="min-w-0">
                  <p className="text-sm font-semibold text-slate-900">Предлагать окно в лист ожидания</p>
                  <p className="mt-0.5 text-xs text-slate-400">
                    Бот сам предложит освободившееся время клиентам из очереди
                  </p>
                </div>
                <Toggle checked={offerWaitlist} onChange={setOfferWaitlist} />
              </div>
            </div>
          </div>
        </section>

        {/* ===== ПРИГЛАШЕНИЕ ===== */}
        <section className="rounded-2xl bg-slate-900 p-4 text-white shadow-xl shadow-slate-900/20">
          <p className="text-sm font-bold">Твоя персональная ссылка</p>
          <p className="mt-1.5 truncate font-mono text-xs text-emerald-400">{link}</p>
          <button
            onClick={handleCopyLink}
            className="mt-3 w-full rounded-xl bg-white py-3 font-bold text-slate-900 transition active:scale-[0.98]"
          >
            Скопировать ссылку
          </button>
        </section>

        {/* ===== НАСТРОЙКИ ПРИЛОЖЕНИЯ ===== */}
        <section>
          <h2 className="mb-2 text-base font-bold text-slate-800">Настройки</h2>
          <div className="divide-y divide-slate-100 rounded-2xl bg-white shadow-sm">
            {APP_SETTINGS.map((item) => {
              const Icon = item.icon
              const isNotifications = item.id === 'notifications'
              const isTheme = item.id === 'theme'

              return (
                <div key={item.id} className="flex items-center gap-3 p-4">
                  <span
                    className={`flex h-10 w-10 shrink-0 items-center justify-center rounded-xl ${
                      item.kind === 'danger'
                        ? 'bg-red-50 text-red-500'
                        : 'bg-slate-100 text-slate-600'
                    }`}
                  >
                    <Icon className="h-5 w-5" />
                  </span>

                  <span
                    className={`min-w-0 flex-1 text-[15px] font-semibold ${
                      item.kind === 'danger' ? 'text-red-600' : 'text-slate-900'
                    }`}
                  >
                    {item.label}
                  </span>

                  {isNotifications && (
                    <Toggle checked={notifications} onChange={setNotifications} />
                  )}

                  {isTheme && (
                    <Toggle checked={darkTheme} onChange={setDarkTheme} />
                  )}

                  {item.kind === 'danger' && (
                    <button
                      onClick={handleLogout}
                      className="shrink-0 rounded-lg bg-red-50 px-3 py-1.5 text-xs font-bold text-red-600 transition active:scale-95"
                    >
                      Выйти
                    </button>
                  )}
                </div>
              )
            })}
          </div>
        </section>
      </div>

      {/* Модальные окна */}
      {isAddServiceOpen && (
        <AddServiceSheet onClose={() => setIsAddServiceOpen(false)} />
      )}
      {isEditProfileOpen && (
        <EditProfileSheet onClose={() => setIsEditProfileOpen(false)} />
      )}
    </div>
  )
}