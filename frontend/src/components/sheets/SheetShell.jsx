export function SheetShell({ children, onClose, closableOnBackdrop = true }) {
    return (
      <div className="fixed inset-0 z-50 flex items-end justify-center">
        <div
          className="absolute inset-0 animate-fade-in bg-slate-900/60 backdrop-blur-[2px]"
          onClick={closableOnBackdrop ? onClose : undefined}
        />
        <div className="relative w-full max-w-[420px] animate-slide-up rounded-t-3xl bg-white px-5 pb-8 pt-3 shadow-2xl">
          <div className="mx-auto mb-6 h-1.5 w-12 rounded-full bg-slate-200" />
          {children}
        </div>
      </div>
    )
  }