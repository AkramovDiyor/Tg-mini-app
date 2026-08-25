import { MASTER_TABS } from '../../mock/masterData'

export function MasterFloatingNav({ active, onChange }) {
  return (
    <div className="fixed bottom-4 left-1/2 z-30 w-[calc(100%-2.5rem)] max-w-[380px] -translate-x-1/2">
      <div className="flex items-center justify-around rounded-full border border-white/50 bg-white/80 px-4 py-2.5 shadow-lg backdrop-blur-xl">
        {MASTER_TABS.map((tab) => {
          const Icon = tab.icon
          const isActive = active === tab.id
          return (
            <button
              key={tab.id}
              onClick={() => onChange(tab.id)}
              className={`flex items-center justify-center p-1 transition-all duration-200 ease-out active:scale-90 ${
                isActive ? '' : 'hover:-translate-y-0.5'
              }`}
            >
              {isActive ? (
                // Ключ меняется при смене активной вкладки → срабатывает animate-pop-in
                <span
                  key={`${tab.id}-active`}
                  className="flex h-10 w-10 items-center justify-center rounded-full bg-emerald-500 text-white shadow-md shadow-emerald-500/30 animate-pop-in"
                >
                  <Icon className="h-5 w-5" strokeWidth={2.5} />
                </span>
              ) : (
                <span className="flex h-10 w-10 items-center justify-center text-slate-400 transition-all duration-200 ease-out hover:text-slate-600">
                  <Icon className="h-5 w-5" />
                </span>
              )}
            </button>
          )
        })}
      </div>
    </div>
  )
}