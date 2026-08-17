import React, { useState } from 'react'
import { useBookingStore } from './store/bookingStore'
import { ConfirmBookingSheet } from './components/sheets/ConfirmBookingSheet'
import { BookingDetailsSheet } from './components/sheets/BookingDetailsSheet'
import { Toast } from './components/ui/Toast'
import { bookSlot, cancelBooking } from './services/api'
import { ServicesPage } from './page/ServicesPage'
import { BookingPage } from './page/BookingPage'
import { MasterPage } from './page/MasterPage'

export default function App() {
  const [role, setRole] = useState('client')
  const [screen, setScreen] = useState('services')
  const [sheetSlot, setSheetSlot] = useState(null)
  const [confirming, setConfirming] = useState(false)
  const [isDetailsOpen, setIsDetailsOpen] = useState(false)
  const [activeBookingForSheet, setActiveBookingForSheet] = useState(null)

  const bookSlotAction = useBookingStore((s) => s.bookSlot)
  const showToast = useBookingStore((s) => s.showToast)
  const service = useBookingStore((s) => s.service)

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
      showToast('Запись отменена. Время освобождено для других')
    } catch (err) {
      const message = err.response?.data?.message || err.message || 'Ошибка при отмене записи'
      showToast(message)
    }
  }

  return (
    <div className="flex min-h-screen justify-center font-sans">
      <div className="relative min-h-screen w-full max-w-[420px] bg-[#F9FAFB] shadow-2xl">
        <div className="sticky top-0 z-30 border-b border-slate-200 bg-[#F9FAFB]/90 px-4 py-3 backdrop-blur">
          <div className="flex gap-2">
            <button
              onClick={() => setRole('client')}
              className={`flex-1 rounded-xl py-2.5 text-sm font-bold transition ${
                role === 'client'
                  ? 'bg-emerald-600 text-white shadow-lg shadow-emerald-600/25'
                  : 'border border-slate-200 bg-white text-slate-500'
              }`}
            >
              Я клиент
            </button>
            <button
              onClick={() => setRole('master')}
              className={`flex-1 rounded-xl py-2.5 text-sm font-bold transition ${
                role === 'master'
                  ? 'bg-emerald-600 text-white shadow-lg shadow-emerald-600/25'
                  : 'border border-slate-200 bg-white text-slate-500'
              }`}
            >
              Я мастер
            </button>
          </div>
        </div>

        {role === 'client' ? (
          <>
            {screen === 'services' ? (
              <ServicesPage
                onPick={() => setScreen('booking')}
                onOpenDetails={handleOpenDetails}
              />
            ) : (
              <BookingPage
                onBack={() => setScreen('services')}
                onSlotClick={(startTime, iso) => setSheetSlot({ startTime, iso })}
              />
            )}
          </>
        ) : (
          <MasterPage />
        )}

        {role === 'client' && sheetSlot && (
          <ConfirmBookingSheet
            slot={sheetSlot}
            confirming={confirming}
            onClose={() => setSheetSlot(null)}
            onConfirm={handleConfirm}
          />
        )}

        {role === 'client' && isDetailsOpen && activeBookingForSheet && (
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