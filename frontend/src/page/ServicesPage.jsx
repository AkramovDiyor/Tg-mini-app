import { useCallback } from 'react'
import { Image, Star, User } from 'lucide-react'
import { QueryRetry } from '../components/AppFrame'
import { ServicesSkeleton } from '../components/skeletons/Skeletons'
import { ServiceCard } from '../components/ui/ServiceCard'
import { useMasterInfoQuery, useServicesQuery } from '../hooks/useClientQueries'
import { useDragScroll } from '../lib/useDragScroll'
import { normalizePhotoUrl } from '../services/api'
import { useBookingStore } from '../store/bookingStore'
import { ActiveBookingBar } from '../widgets/ActiveBookingBar'

export function ServicesPage({ inviteLink }) {
  const pickService = useBookingStore((s) => s.pickService)
  const openDetails = useBookingStore((s) => s.openDetails)
  const galleryRef = useDragScroll()

  const servicesQuery = useServicesQuery(inviteLink)
  const infoQuery = useMasterInfoQuery(inviteLink)

  const handlePick = useCallback(
    (service) => {
      pickService({
        id: service.id,
        name: service.name,
        price: service.price,
        duration: service.duration_min,
        icon: User,
      })
    },
    [pickService],
  )

  if (servicesQuery.isPending || infoQuery.isPending) {
    return <ServicesSkeleton />
  }

  if (servicesQuery.isError && infoQuery.isError) {
    return (
      <QueryRetry
        message="Не удалось загрузить страницу мастера"
        onRetry={() => {
          servicesQuery.refetch()
          infoQuery.refetch()
        }}
      />
    )
  }

  const services = servicesQuery.data || []
  const profile = infoQuery.data
  const masterName = profile?.name || 'Мастер'
  const masterBio = profile?.bio || 'Барбер'
  const masterRating = profile?.rating || null
  const photos = profile?.photos || []

  return (
    <div className="animate-fade-up">
      <div className="relative h-[280px] overflow-hidden bg-gradient-to-b from-slate-900 to-slate-700">
        <div className="pointer-events-none absolute -top-20 left-1/2 h-56 w-56 -translate-x-1/2 rounded-full bg-emerald-500/20 blur-3xl" />
        <div className="pointer-events-none absolute -bottom-24 -right-10 h-48 w-48 rounded-full bg-teal-400/10 blur-3xl" />

        <div className="relative flex h-full flex-col items-center justify-center px-5">
          <div className="relative">
            <div className="flex h-24 w-24 items-center justify-center rounded-full border-4 border-white bg-gradient-to-br from-emerald-400 to-teal-600 text-white shadow-xl">
              <User className="h-12 w-12" />
            </div>
            <span className="absolute bottom-1 right-1 h-4 w-4 rounded-full border-[3px] border-slate-900 bg-emerald-400" />
          </div>

          <h1 className="mt-4 text-2xl font-bold text-white">{masterName}</h1>

          <span className="mt-2.5 flex items-center gap-1.5 rounded-full bg-white/10 px-3 py-1 text-xs font-medium text-white backdrop-blur">
            <Star className="h-3.5 w-3.5 fill-amber-400 text-amber-400" />
            {masterRating ? `${masterRating} · ${masterBio}` : masterBio}
          </span>
        </div>
      </div>

      <div className="relative -mt-8 rounded-t-3xl bg-[#F9FAFB] pb-10">
        <div className="mx-auto mt-3 h-1.5 w-12 rounded-full bg-slate-200" />

        <h2 className="mt-6 px-5 text-lg font-bold text-slate-900">Примеры работ</h2>
        <div
          ref={galleryRef}
          className="thin-scrollbar mt-3 flex cursor-grab snap-x snap-mandatory select-none gap-3 scroll-pl-5 overflow-x-auto px-5 pb-2 active:cursor-grabbing"
        >
          {photos.length === 0 ? (
            <div className="flex h-32 w-40 shrink-0 snap-start flex-col items-center justify-center gap-2 rounded-2xl bg-gradient-to-br from-slate-200 to-slate-300">
              <Image className="h-8 w-8 text-slate-400/70" strokeWidth={1.5} />
              <span className="text-[11px] font-medium text-slate-400">Фото скоро появятся</span>
            </div>
          ) : (
            photos.map((photo) => (
              <div
                key={photo.id}
                className="h-32 w-40 shrink-0 snap-start overflow-hidden rounded-2xl bg-slate-200"
              >
                <img
                  src={normalizePhotoUrl(photo.url)}
                  alt="Работа мастера"
                  className="h-full w-full object-cover"
                  draggable={false}
                />
              </div>
            ))
          )}
        </div>

        <h2 className="mt-6 px-5 text-lg font-bold text-slate-900">Услуги</h2>
        <div className="mt-3 px-5">
          {servicesQuery.isError ? (
            <QueryRetry message="Не удалось загрузить услуги" onRetry={servicesQuery.refetch} />
          ) : (
            services.map((service) => (
              <ServiceCard
                key={service.id}
                service={{
                  id: service.id,
                  name: service.name,
                  price: service.price,
                  duration: service.duration_min,
                  icon: User,
                }}
                onClick={() => handlePick(service)}
              />
            ))
          )}
        </div>
      </div>

      <ActiveBookingBar onOpenDetails={openDetails} />
    </div>
  )
}
