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

export default function App() {
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

        {sheetSlot && (
          <ConfirmBookingSheet
            slot={sheetSlot}
            confirming={confirming}
            onClose={() => setSheetSlot(null)}
            onConfirm={handleConfirm}
          />
        )}

        {isDetailsOpen && activeBooking && (
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