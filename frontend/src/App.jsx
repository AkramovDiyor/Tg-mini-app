import React, { useState } from 'react'
import { useBookingStore } from './store/bookingStore'
import { fromISO, SLOT_TIMES } from './lib/dates'
// import { ServicesPage } from './pages/ServicesPage'
// import { BookingPage } from './pages/BookingPage'
import { ConfirmBookingSheet } from './components/sheets/ConfirmBookingSheet'
import { BookingDetailsSheet } from './components/sheets/BookingDetailsSheet'
import { Toast } from './components/ui/Toast'
import { ServicesPage } from './page/ServicesPage'
import { BookingPage } from './page/BookingPage'
import { MasterPage } from './page/MasterPage'

export default function App() {
  const [role, setRole] = useState('master') // ← НОВОЕ

  const [screen, setScreen] = useState('services')
  const [sheetSlot, setSheetSlot] = useState(null)
  const [confirming, setConfirming] = useState(false)
  const [isDetailsOpen, setIsDetailsOpen] = useState(false)

  const bookSlot = useBookingStore((s) => s.bookSlot)
  const showToast = useBookingStore((s) => s.showToast)
  const activeBooking = useBookingStore((s) => s.activeBooking)
  const cancelActiveBooking = useBookingStore((s) => s.cancelActiveBooking)

  const handleConfirm = () => {
    if (!sheetSlot) return
    setConfirming(true)
    const key = `${sheetSlot.iso}-${sheetSlot.slotIndex}`
    setTimeout(() => {
      bookSlot(key)
      setConfirming(false)
      setSheetSlot(null)
      showToast('Запись подтверждена! За 2 часа до визита пришлём напоминание 🤝')
    }, 1300)
  }

  const handleCancelBooking = () => {
    setIsDetailsOpen(false)
    cancelActiveBooking()
    showToast('Запись отменена. Время освобождено для других')
  }

  return (
    <div className="flex min-h-screen justify-center font-sans">
      <div className="relative min-h-screen w-full max-w-[420px] bg-[#F9FAFB] shadow-2xl">

        {/* ====== ПЕРЕКЛЮЧАТЕЛЬ РОЛЕЙ (mock) ====== */}
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

        {/* ====== КОНТЕНТ ====== */}
        {role === 'client' ? (
          <>
            {screen === 'services' ? (
              <ServicesPage
                onPick={() => setScreen('booking')}
                onOpenDetails={() => setIsDetailsOpen(true)}
              />
            ) : (
              <BookingPage
                onBack={() => setScreen('services')}
                onSlotClick={(iso, slotIndex) => setSheetSlot({ iso, slotIndex })}
              />
            )}
          </>
        ) : (
          <MasterPage />
        )}

        {/* Шиты показываем только для клиента */}
        {role === 'client' && sheetSlot && (
          <ConfirmBookingSheet
          slot={sheetSlot}
          confirming={confirming}
          onClose={() => setSheetSlot(null)}
          onConfirm={handleConfirm}
        />
      )}

      {role === 'client' && isDetailsOpen && activeBooking && (
        <BookingDetailsSheet
          onClose={() => setIsDetailsOpen(false)}
          onCancel={handleCancelBooking}
        />
      )}

      <Toast />
    </div>
  </div>
)
}