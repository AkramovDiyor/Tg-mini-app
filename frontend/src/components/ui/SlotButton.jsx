export function SlotButton({ time, status, onClick }) {
    const base = 'flex items-center justify-center rounded-xl border-2 py-4 text-[15px] font-bold select-none transition'
  
    if (status === 'busy' || status === 'mine') {
      return (
        <button disabled className={`${base} cursor-not-allowed border-transparent bg-slate-100 text-slate-300`}>
          {time}
        </button>
      )
    }
  
    if (status === 'pending') {
      return (
        <button disabled className={`${base} cursor-not-allowed border-amber-200 bg-amber-100 text-amber-700`}>
          {time}
        </button>
      )
    }
  
    return (
      <button
        onClick={onClick}
        className={`${base} border-emerald-500 bg-white text-slate-900 shadow-sm hover:bg-emerald-50 active:scale-95`}
      >
        {time}
      </button>
    )
  }