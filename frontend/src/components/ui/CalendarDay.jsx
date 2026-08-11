export function CalendarDay({ date, selected, isToday, disabled, onClick }) {
    return (
      <button
        onClick={onClick}
        disabled={disabled}
        className={`flex aspect-square items-center justify-center rounded-xl text-sm font-semibold transition ${
          disabled
            ? 'cursor-not-allowed text-slate-300'
            : selected
              ? 'bg-emerald-600 text-white shadow-lg shadow-emerald-600/30'
              : isToday
                ? 'bg-emerald-50 text-emerald-700 ring-1 ring-inset ring-emerald-200 active:scale-95'
                : 'border border-slate-200/70 bg-white text-slate-700 shadow-sm hover:border-emerald-300 active:scale-95'
        }`}
      >
        {date.getDate()}
      </button>
    )
  }