import React, { useState, useEffect } from 'react'
import { useBookingStore } from './store/bookingStore'
import { ConfirmBookingSheet } from './components/sheets/ConfirmBookingSheet'
import { BookingDetailsSheet } from './components/sheets/BookingDetailsSheet'
import { Toast } from './components/ui/Toast'
import { bookSlot, cancelBooking, fetchUserIdentity, getInviteLinkFromUrl } from './services/api'
import { ServicesPage } from './page/ServicesPage'
import { BookingPage } from './page/BookingPage'
import { MasterPage } from './page/MasterPage'

export default function App() {
  const [role, setRole] = useState(null) // 🔥 null = ещё не определили
  const [identity, setIdentity] = useState(null) // 🔥 полные данные из /me
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [screen, setScreen] = useState('services')
  const [sheetSlot, setSheetSlot] = useState(null)
  const [confirming, setConfirming] = useState(false)
  const [isDetailsOpen, setIsDetailsOpen] = useState(false)
  const [activeBookingForSheet, setActiveBookingForSheet] = useState(null)
  const [bookingsVersion, setBookingsVersion] = useState(0)

  const bookSlotAction = useBookingStore((s) => s.bookSlot)
  const showToast = useBookingStore((s) => s.showToast)
  const service = useBookingStore((s) => s.service)

  // 🔥 Автоматическое определение роли через бэкенд
  useEffect(() => {
    async function determineRole() {
      try {
        const me = await fetchUserIdentity()
        setIdentity(me)
        setRole(me.role)
        console.log('🎯 Роль определена:', me.role)
      } catch (err) {
        console.error('❌ Ошибка определения роли:', err)
        setError('Не удалось определить пользователя. Попробуйте перезапустить приложение.')
      } finally {
        setLoading(false)
      }
    }

    determineRole()
  }, [])

  const handleConfirm = async () => {
    if (!sheetSlot) return
    setConfirming(true)

    try {
      await bookSlot({
        start_time: sheetSlot.startTime,
        service_id: service.id,
        name: 'Имя клиента',
        price: service.price,
      })

      const key = `${sheetSlot.iso}-${sheetSlot.startTime}`
      bookSlotAction(key)
      setConfirming(false)
      setSheetSlot(null)
      setBookingsVersion((v) => v + 1)
      showToast('Запись подтверждена! За 2 часа до визита пришлём напоминание 🤝')
    } catch (err) {
      setConfirming(false)
      const message = err.response?.data?.message || err.message || 'Ошибка при бронировании'
      showToast(message)
    }
  }

  const handleOpenDetails = (booking) => {
    setActiveBookingForSheet(booking)
    setIsDetailsOpen(true)
  }

  const handleCancelBooking = async (bookingId) => {
    try {
      await cancelBooking(bookingId)
      setIsDetailsOpen(false)
      setActiveBookingForSheet(null)
      setBookingsVersion((v) => v + 1)
      showToast('Запись отменена. Время освобождено для других')
    } catch (err) {
      const message = err.response?.data?.message || err.message || 'Ошибка при отмене записи'
      showToast(message)
    }
  }

  // 🔥 Экран загрузки
  if (loading) {
    return (
      <div className="flex min-h-screen justify-center font-sans">
        <div className="relative min-h-screen w-full max-w-[420px] bg-[#F9FAFB] shadow-2xl">
          <div className="flex min-h-screen items-center justify-center px-6">
            <div className="text-center">
              <div className="mx-auto mb-4 h-12 w-12 animate-spin rounded-full border-4 border-emerald-200 border-t-emerald-600" />
              <p className="text-sm text-slate-500">Определяем пользователя...</p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  // 🔥 Экран ошибки
  if (error) {
    return (
      <div className="flex min-h-screen justify-center font-sans">
        <div className="relative min-h-screen w-full max-w-[420px] bg-[#F9FAFB] shadow-2xl">
          <div className="flex min-h-screen items-center justify-center px-6">
            <div className="text-center">
              <div className="mb-3 text-4xl">⚠️</div>
              <h2 className="mb-2 text-lg font-bold text-slate-800">Ошибка загрузки</h2>
              <p className="text-sm text-slate-500">{error}</p>
              <button
                onClick={() => window.location.reload()}
                className="mt-4 rounded-xl bg-emerald-600 px-6 py-2.5 text-sm font-bold text-white shadow-lg shadow-emerald-600/25 transition hover:bg-emerald-700"
              >
                Попробовать снова
              </button>
            </div>
          </div>
        </div>
      </div>
    )
  }

  // 🔥 Мастер
  if (role === 'master') {
    return (
      <div className="flex min-h-screen justify-center font-sans">
        <div className="relative min-h-screen w-full max-w-[420px] bg-[#F9FAFB] shadow-2xl">
          <MasterPage master={identity?.master} />
          <Toast />
        </div>
      </div>
    )
  }

  // 🔥 Клиент — проверка invite_link
  const inviteLink = getInviteLinkFromUrl()

  if (!inviteLink) {
    return (
      <div className="flex min-h-screen justify-center font-sans">
        <div className="relative min-h-screen w-full max-w-[420px] bg-[#F9FAFB] shadow-2xl">
          <div className="flex min-h-screen items-center justify-center px-6">
            <div className="text-center">
              <div className="mb-3 text-4xl">🔗</div>
              <h2 className="mb-2 text-lg font-bold text-slate-800">Нужна ссылка мастера</h2>
              <p className="mb-1 text-sm text-slate-500">
                Откройте Mini App по ссылке, которую вам отправил мастер.
              </p>
              <p className="text-xs text-slate-400">
                Например: <code className="rounded bg-slate-100 px-1.5 py-0.5">t.me/bot?startapp=abc123</code>
              </p>
            </div>
          </div>
          <Toast />
        </div>
      </div>
    )
  }

  // 🔥 Клиент с invite_link
  return (
    <div className="flex min-h-screen justify-center font-sans">
      <div className="relative min-h-screen w-full max-w-[420px] bg-[#F9FAFB] shadow-2xl">
        {screen === 'services' ? (
          <ServicesPage
            onPick={() => setScreen('booking')}
            onOpenDetails={handleOpenDetails}
            bookingsVersion={bookingsVersion}
          />
        ) : (
          <BookingPage
            onBack={() => setScreen('services')}
            onSlotClick={(startTime, iso) => setSheetSlot({ startTime, iso })}
          />
        )}

        {sheetSlot && (
          <ConfirmBookingSheet
            slot={sheetSlot}
            confirming={confirming}
            onClose={() => setSheetSlot(null)}
            onConfirm={handleConfirm}
          />
        )}

        {isDetailsOpen && activeBookingForSheet && (
          <BookingDetailsSheet
            booking={activeBookingForSheet}
            onClose={() => setIsDetailsOpen(false)}
            onCancel={handleCancelBooking}
          />
        )}

        <Toast />
      </div>
    </div>
  )
}