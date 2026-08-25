export function AppFrame({ children }) {
  return (
    <div className="flex min-h-screen justify-center font-sans">
      <div className="relative min-h-screen w-full max-w-[420px] bg-[#F9FAFB] shadow-2xl">
        {children}
      </div>
    </div>
  )
}

export function QueryRetry({ message, onRetry }) {
  return (
    <div className="flex flex-col items-center justify-center py-20 px-5">
      <p className="text-center text-sm text-red-500">{message}</p>
      {onRetry && (
        <button
          type="button"
          onClick={onRetry}
          className="mt-4 rounded-xl bg-emerald-500 px-6 py-2.5 text-sm font-bold text-white transition active:scale-95"
        >
          Повторить
        </button>
      )}
    </div>
  )
}
