export function Pulse({ className }) {
  return <div className={`animate-pulse rounded-xl bg-slate-200 ${className}`} />
}

export function IdentitySkeleton() {
  return (
    <div className="flex min-h-screen items-center justify-center px-6">
      <div className="w-full max-w-xs text-center">
        <Pulse className="mx-auto mb-4 h-12 w-12 rounded-full" />
        <Pulse className="mx-auto h-3 w-40" />
      </div>
    </div>
  )
}

export function PageFallbackSkeleton() {
  return (
    <div className="px-5 pt-6">
      <Pulse className="h-8 w-48" />
      <Pulse className="mt-4 h-24 w-full rounded-2xl" />
      <Pulse className="mt-3 h-24 w-full rounded-2xl" />
      <Pulse className="mt-3 h-24 w-full rounded-2xl" />
    </div>
  )
}

export function ServicesSkeleton() {
  return (
    <div>
      <div className="h-[280px] bg-slate-800 px-5 pt-16">
        <Pulse className="mx-auto h-24 w-24 rounded-full bg-slate-600" />
        <Pulse className="mx-auto mt-4 h-7 w-40 bg-slate-600" />
        <Pulse className="mx-auto mt-3 h-5 w-32 rounded-full bg-slate-600" />
      </div>
      <div className="-mt-8 rounded-t-3xl bg-[#F9FAFB] px-5 pt-8 pb-10">
        <Pulse className="h-5 w-36" />
        <div className="mt-3 flex gap-3">
          <Pulse className="h-32 w-40 shrink-0 rounded-2xl" />
          <Pulse className="h-32 w-40 shrink-0 rounded-2xl" />
        </div>
        <Pulse className="mt-6 h-5 w-24" />
        <Pulse className="mt-3 h-20 w-full rounded-2xl" />
        <Pulse className="mt-2 h-20 w-full rounded-2xl" />
        <Pulse className="mt-2 h-20 w-full rounded-2xl" />
      </div>
    </div>
  )
}

export function SlotsSkeleton() {
  return (
    <div className="grid grid-cols-3 gap-2.5">
      {Array.from({ length: 9 }, (_, i) => (
        <Pulse key={i} className="h-14 rounded-xl" />
      ))}
    </div>
  )
}

export function MasterListSkeleton() {
  return (
    <div className="space-y-3 pt-2">
      <Pulse className="h-8 w-40" />
      <Pulse className="h-28 w-full rounded-3xl" />
      <Pulse className="h-24 w-full rounded-2xl" />
      <Pulse className="h-24 w-full rounded-2xl" />
    </div>
  )
}
