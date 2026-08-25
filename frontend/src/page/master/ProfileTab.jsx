import { useState } from 'react'
import { MapPin, Bell, Moon, LogOut, User, Pencil, ChevronRight } from 'lucide-react'
import { QueryRetry } from '../../components/AppFrame'
import { MasterListSkeleton } from '../../components/skeletons/Skeletons'
import { Toggle } from '../../components/ui/Toggle'
import { EditProfileSheet } from '../../components/sheets/EditProfileSheet'
import { useMasterProfileQuery } from '../../hooks/useMasterQueries'
import { useBookingStore } from '../../store/bookingStore'

const APP_SETTINGS = [
  { id: 'notifications', label: 'Уведомления', icon: Bell, kind: 'toggle' },
  { id: 'theme', label: 'Тёмная тема', icon: Moon, kind: 'toggle' },
  { id: 'logout', label: 'Выйти из аккаунта', icon: LogOut, kind: 'danger' },
]

export function ProfileTab() {
  const showToast = useBookingStore((s) => s.showToast)
  const { data: profile, isPending, isError, refetch } = useMasterProfileQuery()
  const [isEditProfileOpen, setIsEditProfileOpen] = useState(false)
  const [notifications, setNotifications] = useState(true)
  const [darkTheme, setDarkTheme] = useState(false)

  const handleCopyLink = () => {
    if (!profile?.invite_link) return
    try {
      navigator.clipboard?.writeText(profile.invite_link)
      showToast('Ссылка скопирована в буфер обмена')
    } catch {
      showToast('Не удалось скопировать. Скопируйте вручную')
    }
  }

  if (isPending) return <MasterListSkeleton />
  if (isError) {
    return <QueryRetry message="Не удалось загрузить профиль" onRetry={refetch} />
  }

  return (
    <div className="pb-24">
      <button
        type="button"
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

      <button
        type="button"
        onClick={() => setIsEditProfileOpen(true)}
        className="mb-6 flex w-full items-center gap-3 rounded-2xl bg-white p-4 text-left shadow-sm transition active:scale-[0.98]"
      >
        <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-emerald-50 text-emerald-600">
          <MapPin className="h-5 w-5" />
        </div>
        <div className="min-w-0 flex-1">
          <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">Адрес студии</p>
          <p className="mt-0.5 truncate font-bold text-slate-900">
            {profile?.address || 'Адрес не указан'}
          </p>
        </div>
        <ChevronRight className="h-5 w-5 shrink-0 text-slate-300" />
      </button>

      <div className="space-y-6">
        <section className="rounded-2xl bg-slate-900 p-4 text-white shadow-xl shadow-slate-900/20">
          <p className="text-sm font-bold">Твоя персональная ссылка</p>
          <p className="mt-1.5 truncate font-mono text-xs text-emerald-400">
            {profile?.invite_link || 'Ссылка не создана'}
          </p>
          <button
            type="button"
            onClick={handleCopyLink}
            className="mt-3 w-full rounded-xl bg-white py-3 font-bold text-slate-900 transition active:scale-[0.98]"
          >
            Скопировать ссылку
          </button>
        </section>

        <section>
          <h2 className="mb-2 text-base font-bold text-slate-800">Настройки</h2>
          <div className="divide-y divide-slate-100 rounded-2xl bg-white shadow-sm">
            {APP_SETTINGS.map((item) => {
              const Icon = item.icon
              return (
                <div key={item.id} className="flex items-center gap-3 p-4">
                  <span
                    className={`flex h-10 w-10 shrink-0 items-center justify-center rounded-xl ${
                      item.kind === 'danger' ? 'bg-red-50 text-red-500' : 'bg-slate-100 text-slate-600'
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
                  {item.id === 'notifications' && (
                    <Toggle checked={notifications} onChange={setNotifications} />
                  )}
                  {item.id === 'theme' && <Toggle checked={darkTheme} onChange={setDarkTheme} />}
                  {item.kind === 'danger' && (
                    <button
                      type="button"
                      onClick={() => showToast('Выход выполнен. До встречи!')}
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

      {isEditProfileOpen && (
        <EditProfileSheet
          profile={profile}
          onClose={() => setIsEditProfileOpen(false)}
        />
      )}
    </div>
  )
}
