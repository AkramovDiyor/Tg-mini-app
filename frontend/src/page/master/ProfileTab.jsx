import { useState, useEffect } from 'react'
import {
  MapPin, Bell, Moon, LogOut, User, Pencil, ChevronRight,
} from 'lucide-react'
import { useBookingStore } from '../../store/bookingStore'
import { fetchMasterProfile } from '../../services/api'
import { Toggle } from '../../components/ui/Toggle'
import { EditProfileSheet } from '../../components/sheets/EditProfileSheet'

const APP_SETTINGS = [
  { id: 'notifications', label: 'Уведомления',       icon: Bell,   kind: 'toggle' },
  { id: 'theme',         label: 'Тёмная тема',       icon: Moon,   kind: 'toggle' },
  { id: 'logout',        label: 'Выйти из аккаунта', icon: LogOut, kind: 'danger' },
]

export function ProfileTab() {
  const showToast = useBookingStore((s) => s.showToast)

  const [profile, setProfile] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [isEditProfileOpen, setIsEditProfileOpen] = useState(false)

  // Мелкие настройки приложения
  const [notifications, setNotifications] = useState(true)
  const [darkTheme, setDarkTheme] = useState(false)

  useEffect(() => {
    loadProfileData()
  }, [])

  const loadProfileData = async () => {
    try {
      setLoading(true)
      setError(null)
      const profileData = await fetchMasterProfile()
      setProfile(profileData)
    } catch (err) {
      console.error('Failed to load profile:', err)
      setError('Не удалось загрузить профиль')
    } finally {
      setLoading(false)
    }
  }

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
      {/* ===== ШАПКА ПРОФИЛЯ ===== */}
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

      {/* ===== КАРТОЧКА АДРЕСА ===== */}
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
        {/* ===== ПЕРСОНАЛЬНАЯ ССЫЛКА ===== */}
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

      {/* Модалка редактирования профиля */}
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