import { lazy, Suspense } from 'react'
import { AppFrame } from './components/AppFrame'
import { BookingDetailsSheet } from './components/sheets/BookingDetailsSheet'
import { ConfirmBookingSheet } from './components/sheets/ConfirmBookingSheet'
import { IdentitySkeleton, PageFallbackSkeleton } from './components/skeletons/Skeletons'
import { Toast } from './components/ui/Toast'
import { useBookSlotMutation, useCancelBookingMutation } from './hooks/useBookingMutations'
import { useMeQuery } from './hooks/useMeQuery'
import { getClientDisplayName } from './lib/telegram'
import { useBookingStore } from './store/bookingStore'


export default function App() {
  const { data: identity, isPending, isError, error, refetch } = useMeQuery()

  const screen = useBookingStore((s) => s.screen)
  const service = useBookingStore((s) => s.service)
  const sheetSlot = useBookingStore((s) => s.sheetSlot)
  const detailsBooking = useBookingStore((s) => s.detailsBooking)
  const closeSlotSheet = useBookingStore((s) => s.closeSlotSheet)
  const closeDetails = useBookingStore((s) => s.closeDetails)

  const inviteLink =
    identity?.invite_link ||
    (typeof window !== 'undefined'
      ? new URLSearchParams(window.location.search).get('startapp')
      : null)

  const bookMutation = useBookSlotMutation({ inviteLink: inviteLink || '' })
  const cancelMutation = useCancelBookingMutation()

  if (isPending) {
    return (
      <AppFrame>
        <IdentitySkeleton />
      </AppFrame>
    )
  }

  if (isError) {
    return (
      <AppFrame>
        <div className="flex min-h-screen items-center justify-center px-6">
          <div className="text-center">
            <h2 className="mb-2 text-lg font-bold text-slate-800">Ошибка загрузки</h2>
            <p className="mb-4 text-sm text-slate-500">
              {error?.message || 'Откройте приложение через Telegram.'}
            </p>
            <button
              type="button"
              onClick={() => refetch()}
              className="rounded-xl bg-emerald-600 px-5 py-2.5 text-sm font-bold text-white"
            >
              Повторить
            </button>
          </div>
        </div>
      </AppFrame>
    )
  }

  if (identity?.role === 'master') {
    return (
      <AppFrame>
        <Suspense fallback={<PageFallbackSkeleton />}>
          <MasterPage />
        </Suspense>
        <Toast />
      </AppFrame>
    )
  }

  if (!inviteLink) {
    return (
      <AppFrame>
        <div className="flex min-h-screen items-center justify-center px-6">
          <div className="text-center">
            <h2 className="mb-2 text-lg font-bold text-slate-800">Нужна ссылка мастера</h2>
            <p className="mb-1 text-sm text-slate-500">
              Откройте Mini App по ссылке, которую вам отправил мастер.
            </p>
            <p className="text-xs text-slate-400">
              Например:{' '}
              <code className="rounded bg-slate-100 px-1.5 py-0.5">t.me/bot?startapp=abc123</code>
            </p>
          </div>
        </div>
        <Toast />
      </AppFrame>
    )
  }

  const showBooking = screen === 'booking' && service

  const handleConfirm = () => {
    if (!sheetSlot || !service) return
    bookMutation.mutate({
      start_time: sheetSlot.startTime,
      date: sheetSlot.iso,
      service_id: service.id,
      name: getClientDisplayName(),
      price: service.price,
      service_name: service.name,
      service_duration: service.duration,
    })
  }

  return (
    <AppFrame>
      <Suspense fallback={<PageFallbackSkeleton />}>
        {showBooking ? (
          <BookingPage inviteLink={inviteLink} />
        ) : (
          <ServicesPage inviteLink={inviteLink} />
        )}
      </Suspense>

      {sheetSlot && (
        <ConfirmBookingSheet
          slot={sheetSlot}
          confirming={bookMutation.isPending}
          onClose={closeSlotSheet}
          onConfirm={handleConfirm}
        />
      )}

      {detailsBooking && (
        <BookingDetailsSheet
          booking={detailsBooking}
          onClose={closeDetails}
          onCancel={() =>
            cancelMutation.mutate({
              bookingId: detailsBooking.booking_id,
              booking: detailsBooking,
            })
          }
        />
      )}

      <Toast />
    </AppFrame>
  )
}
