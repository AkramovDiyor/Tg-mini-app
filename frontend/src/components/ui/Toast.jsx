import { useEffect } from 'react'
import { Bell } from 'lucide-react'
import { useBookingStore } from '../../store/bookingStore'

export function Toast() {
  const toast = useBookingStore((s) => s.toast)
  const hideToast = useBookingStore((s) => s.hideToast)

  useEffect(() => {
    if (!toast) return
    const t = setTimeout(hideToast, 3500)
    return () => clearTimeout(t)
  }, [toast, hideToast])

  if (!toast) return null

  return (
    <div className="pointer-events-none fixed inset-x-0 bottom-6 z-[60] flex justify-center px-5">
      <div
        key={toast.id}
        className="pointer-events-auto flex w-full max-w-[380px] animate-toast-in items-start gap-2.5 rounded-2xl bg-slate-900/95 px-4 py-3.5 text-white shadow-2xl"
      >
        <Bell className="mt-0.5 h-4 w-4 shrink-0 text-emerald-400" />
        <p className="text-[13px] font-medium leading-snug">{toast.text}</p>
      </div>
    </div>
  )
}