import { Clock, ChevronRight } from 'lucide-react'
import { rub } from '../../lib/currency'

export function ServiceCard({ service, onClick }) {
  const Icon = service.icon
  return (
    <button
      onClick={onClick}
      className="group relative mb-4 w-full rounded-2xl bg-white p-5 pb-16 text-left shadow-md transition active:scale-[0.98]"
    >
      <div className="flex items-center justify-between gap-4">
        <div className="min-w-0 flex-1">
          <p className="text-lg font-bold text-slate-900">{service.name}</p>
          <p className="mt-1 flex items-center gap-1.5 text-sm text-slate-400">
            <Clock className="h-4 w-4" />
            {service.duration} мин
          </p>
        </div>
        <p className="text-xl font-extrabold text-slate-900">{rub(service.price)}</p>
      </div>
      <span className="absolute bottom-4 right-4 flex h-10 w-10 items-center justify-center rounded-full bg-emerald-500 text-white shadow-lg shadow-emerald-500/30 transition-transform duration-200 group-active:translate-x-0.5">
        <ChevronRight className="h-5 w-5" />
      </span>
    </button>
  )
}