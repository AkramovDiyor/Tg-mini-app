import { ChevronLeft, ChevronRight } from 'lucide-react'

export function NavArrow({ direction, onClick, disabled }) {
  const Icon = direction === 'prev' ? ChevronLeft : ChevronRight
  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className="flex h-9 w-9 items-center justify-center rounded-xl border border-slate-200 bg-white text-slate-600 shadow-sm transition active:scale-90 disabled:cursor-not-allowed disabled:opacity-40"
    >
      <Icon className="h-5 w-5" />
    </button>
  )
}