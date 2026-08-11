export function SlotsLegend() {
    return (
      <div className="mt-4 flex items-center gap-4 text-[11px] text-slate-400">
        <span className="flex items-center gap-1.5">
          <i className="h-2 w-2 rounded-full bg-emerald-500" />
          Свободно
        </span>
        <span className="flex items-center gap-1.5">
          <i className="h-2 w-2 rounded-full bg-amber-400" />
          Подтверждается
        </span>
        <span className="flex items-center gap-1.5">
          <i className="h-2 w-2 rounded-full bg-slate-300" />
          Занято
        </span>
      </div>
    )
  }