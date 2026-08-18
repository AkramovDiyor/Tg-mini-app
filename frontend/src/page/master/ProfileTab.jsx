import { useState, useEffect } from 'react'
import {
  MapPin, Clock, Plus, Bell, Moon, LogOut, User, Pencil, ChevronRight,
} from 'lucide-react'
import { useBookingStore } from '../../store/bookingStore'
import { rub } from '../../lib/currency'
import { fetchMasterProfile, fetchMasterServices, updateSettings } from '../../services/api'
import { Toggle } from '../../components/ui/Toggle'
import { AddServiceSheet } from '../../components/sheets/AddServiceSheet'
import { EditProfileSheet } from '../../components/sheets/EditProfileSheet'

const APP_SETTINGS = [
  { id: 'notifications', label: 'Уведомления', icon: Bell, kind: 'toggle' },
  { id: 'theme', label: 'Тёмная тема', icon: Moon, kind: 'toggle' },
  { id: 'logout', label: 'Выйти из аккаунта', icon: LogOut, kind: 'danger' },
]

const WEEK_TAGS = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']

export function ProfileTab() {
  const showToast = useBookingStore((s) => s.showToast)

  const [profile, setProfile] = useState(null)
  const [services, setServices] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const [isAddServiceOpen, setIsAddServiceOpen] = useState(false)
  const [isEditProfileOpen, setIsEditProfileOpen] = useState(false)

  // График работы (загружается из profile.work_hours)
  const [workDays, setWorkDays] = useState([true, true, true, true, true, false, false])
  const [hours, setHours] = useState({
    start: '09:00', end: '20:00', lunchFrom: '13:00', lunchTo: '14:00',
  })

  // Состояния настроек (загружаются из profile.settings)
  const [autoCancel, setAutoCancel] = useState(true)
  const [cancelHours, setCancelHours] = useState('2')
  const [offerWaitlist, setOfferWaitlist] = useState(true)
  const [notifications, setNotifications] = useState(true)
  const [darkTheme, setDarkTheme] = useState(false)

  useEffect(() => {
    loadProfileData()
  }, [])

  const loadProfileData = async () => {
    try {
      setLoading(true)
      setError(null)

      const [profileData, servicesData] = await Promise.all([
        fetchMasterProfile(),
        fetchMasterServices(),
      ])

      setProfile(profileData)
      setServices(Array.isArray(servicesData) ? servicesData : [])

      // Загружаем график работы из profile.work_hours
      if (profileData?.work_hours) {
        const wh = profileData.work_hours
        // work_days может быть массивом чисел (1=Пн .. 7=Вс) или массивом булевых
        if (Array.isArray(wh.work_days) && wh.work_days.length === 7) {
          if (typeof wh.work_days[0] === 'boolean') {
            setWorkDays(wh.work_days)
          } else {
            setWorkDays([1, 2, 3, 4, 5, 6, 7].map((d) => wh.work_days.includes(d)))
          }
        }
        setHours({
          start: wh.start_time || wh.start || '09:00',
          end: wh.end_time || wh.end || '20:00',
          lunchFrom: wh.lunch_start || wh.lunchFrom || '13:00',
          lunchTo: wh.lunch_end || wh.lunchTo || '14:00',
        })
      }

      // Загружаем настройки из profile.settings
      if (profileData?.settings) {
        setAutoCancel(profileData.settings.auto_cancel ?? true)
        setCancelHours(profileData.settings.cancel_hours || '2')
        setOfferWaitlist(profileData.settings.offer_waitlist ?? true)
      }
    } catch (err) {
      console.error('Failed to load profile:', err)
      setError('Не удалось загрузить профиль')
    } finally {
      setLoading(false)
    }
  }

  const toggleDay = (i) =>
    setWorkDays((days) => days.map((d, idx) => (idx === i ? !d : d)))

  const handleCopyLink = () => {
    if (!profile?.invite_link) return
    try {
      navigator.clipboard?.writeText(profile.invite_link)
      showToast('Ссылка скопирована в буфер обмена')
    } catch {
      showToast('Не удалось скопировать. Скопируйте вручную')
    }
  }

  const handleLogout = () => {
    showToast('Выход выполнен. До встречи! 👋')
  }

  // Сохраняем и график, и настройки одним запросом PUT /settings
  const handleSaveSettings = async () => {
    try {
      await updateSettings({
        work_hours: {
          work_days: workDays,
          start_time: hours.start,
          end_time: hours.end,
          lunch_start: hours.lunchFrom,
          lunch_end: hours.lunchTo,
        },
        settings: {
          auto_cancel: autoCancel,
          cancel_hours: cancelHours,
          offer_waitlist: offerWaitlist,
        },
      })
      showToast('Настройки сохранены ✨')
    } catch (err) {
      showToast('Ошибка при сохранении настроек')
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <p className="text-sm text-slate-400">Загрузка профиля...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center py-20">
        <p className="text-sm text-red-500">{error}</p>
        <button
          onClick={loadProfileData}
          className="mt-4 rounded-xl bg-emerald-500 px-6 py-2.5 text-sm font-bold text-white transition active:scale-95"
        >
          Повторить
        </button>
      </div>
    )
  }

  return (
    <div className="pb-24">
      {/* ШАПКА ПРОФИЛЯ */}
      <button
        onClick={() => setIsEditProfileOpen(true)}
        className="mb-6 flex w-full items-center gap-4 rounded-2xl text-left transition active:scale-[0.98]"
      >
        <div className="relative shrink-0">
          <div className="flex h-20 w-20 items-center justify-center rounded-full bg-slate-200 text-slate-400">
            <User className="h-10 w-10" strokeWidth={1.5} />
          </div>
        </div>

        <div className="min-w-0 flex-1">
          <h1 className="text-xl font-extrabold text-slate-900">{profile?.name || 'Мастер'}</h1>
          <p className="mt-0.5 text-sm text-slate-400">{profile?.bio || 'Барбер'}</p>
          <div className="mt-1.5 flex items-center gap-1">
            <span className="h-2 w-2 rounded-full bg-emerald-500" />
            <span className="text-xs font-semibold text-emerald-600">Онлайн</span>
          </div>
        </div>

        <Pencil className="h-5 w-5 shrink-0 text-slate-300" />
      </button>

      {/* КАРТОЧКА АДРЕСА */}
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
          <p className="mt-0.5 truncate font-bold text-slate-900">{profile?.address || 'Адрес не указан'}</p>
        </div>
        <ChevronRight className="h-5 w-5 shrink-0 text-slate-300" />
      </button>

      <div className="space-y-6">
        {/* УСЛУГИ */}
        <section>
          <h2 className="mb-2 text-base font-bold text-slate-800">Услуги</h2>
          <div className="rounded-2xl bg-white p-4 shadow-sm">
            <div className="space-y-2">
              {services.length === 0 ? (
                <p className="py-4 text-center text-sm text-slate-400">Услуги не добавлены</p>
              ) : (
                services.map((service) => (
                  <div
                    key={service.id}
                    className="flex w-full items-center gap-3 rounded-xl bg-slate-50 p-3 text-left"
                  >
                    <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-emerald-50 text-emerald-600">
                      <User className="h-5 w-5" />
                    </span>
                    <span className="min-w-0 flex-1">
                      <span className="block text-sm font-bold text-slate-900">{service.name}</span>
                      <span className="mt-0.5 flex items-center gap-1 text-xs text-slate-400">
                        <Clock className="h-3 w-3" />
                        {service.duration_min} мин
                      </span>
                    </span>
                    <span className="text-sm font-bold text-slate-900">{rub(service.price)}</span>
                  </div>
                ))
              )}
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

        {/* ===== ГРАФИК РАБОТЫ (возвращён и подключён к API) ===== */}
        <section>
          <h2 className="mb-2 text-base font-bold text-slate-800">График работы</h2>
          <div className="rounded-2xl bg-white p-4 shadow-sm">
            {/* Дни недели — круглые теги */}
            <div className="flex flex-wrap gap-2">
              {WEEK_TAGS.map((tag, i) => {
                const isWork = workDays[i]
                return (
                  <button
                    key={tag}
                    onClick={() => toggleDay(i)}
                    className={`h-10 rounded-full px-4 text-sm font-bold transition ${
                      isWork
                        ? 'bg-emerald-500 text-white shadow-md shadow-emerald-500/25 active:scale-95'
                        : 'bg-slate-100 text-slate-400 active:scale-95'
                    }`}
                  >
                    {tag}
                  </button>
                )
              })}
            </div>

            {/* Часы работы */}
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

            {/* Обед */}
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

        {/* ПРАВИЛА ЗАПИСИ */}
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
                  {['1', '2', '4'].map((h) => (
                    <button
                      key={h}
                      onClick={() => setCancelHours(h)}
                      className={`flex-1 rounded-xl py-2 text-xs font-bold transition ${cancelHours === h
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

            {/* Кнопка сохранения теперь сохраняет и график, и правила */}
            <button
              onClick={handleSaveSettings}
              className="mt-4 w-full rounded-xl bg-emerald-500 py-3 font-bold text-white transition active:scale-[0.98]"
            >
              Сохранить настройки
            </button>
          </div>
        </section>

        {/* ПРИГЛАШЕНИЕ */}
        <section className="rounded-2xl bg-slate-900 p-4 text-white shadow-xl shadow-slate-900/20">
          <p className="text-sm font-bold">Твоя персональная ссылка</p>
          <p className="mt-1.5 truncate font-mono text-xs text-emerald-400">
            {profile?.invite_link || 'Ссылка не создана'}
          </p>
          <button
            onClick={handleCopyLink}
            className="mt-3 w-full rounded-xl bg-white py-3 font-bold text-slate-900 transition active:scale-[0.98]"
          >
            Скопировать ссылку
          </button>
        </section>

        {/* НАСТРОЙКИ ПРИЛОЖЕНИЯ */}
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
                    className={`flex h-10 w-10 shrink-0 items-center justify-center rounded-xl ${item.kind === 'danger'
                        ? 'bg-red-50 text-red-500'
                        : 'bg-slate-100 text-slate-600'
                      }`}
                  >
                    <Icon className="h-5 w-5" />
                  </span>

                  <span
                    className={`min-w-0 flex-1 text-[15px] font-semibold ${item.kind === 'danger' ? 'text-red-600' : 'text-slate-900'
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

      {isAddServiceOpen && (
        <AddServiceSheet
          onClose={() => setIsAddServiceOpen(false)}
          onCreated={loadProfileData}
        />
      )}
      {isEditProfileOpen && (
        <EditProfileSheet
          profile={profile}
          onClose={() => setIsEditProfileOpen(false)}
          onSaved={loadProfileData}
        />
      )}
    </div>
  )
}