export function DayPill({ day, selected, onClick }) {
    return (
      <button
        onClick={onClick}
        className={`flex w-[72px] shrink-0 snap-start flex-col items-center gap-1 rounded-xl border py-2.5 transition-all ${
          selected
            ? 'border-emerald-600 bg-emerald-600 text-white shadow-lg shadow-emerald-600/30'
            : 'border-slate-200 bg-white text-slate-500 shadow-sm active:scale-95'
        }`}
      >
        <span className={`text-[10px] font-semibold uppercase ${selected ? 'text-emerald-100' : 'text-slate-400'}`}>
          {day.label}
        </span>
        <span className="text-lg font-bold leading-none">{day.date.getDate()}</span>
      </button>
    )
  }