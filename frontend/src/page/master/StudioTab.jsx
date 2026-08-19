import { useState, useEffect, useRef } from 'react'
import {
  Clock, Plus, User, Trash2,
} from 'lucide-react'
import { useBookingStore } from '../../store/bookingStore'
import { rub } from '../../lib/currency'
import {
  fetchMasterProfile, fetchMasterServices, updateSettings,
  fetchPhotos, deletePhoto, normalizePhotoUrl,
} from '../../services/api'
import { Toggle } from '../../components/ui/Toggle'
import { AddServiceSheet } from '../../components/sheets/AddServiceSheet'
import { EditServiceSheet } from '../../components/sheets/EditServiceSheet'
import { AddPhotoSheet } from '../../components/sheets/AddPhotoSheet'

const WEEK_TAGS = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']

export function StudioTab() {
  const showToast = useBookingStore((s) => s.showToast)

  // === Услуги ===
  const [servicesVersion, setServicesVersion] = useState(0)
  const [services, setServices] = useState([])
  const [isAddServiceOpen, setIsAddServiceOpen] = useState(false)
  const [editingService, setEditingService] = useState(null)

  // === График работы + настройки (из профиля) ===
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [workDays, setWorkDays] = useState([true, true, true, true, true, false, false])
  const [hours, setHours] = useState({
    start: '09:00', end: '20:00', lunchFrom: '13:00', lunchTo: '14:00',
  })

  // === Правила записи ===
  const [autoCancel, setAutoCancel] = useState(true)
  const [cancelHours, setCancelHours] = useState('2')
  const [offerWaitlist, setOfferWaitlist] = useState(true)

  // === Фото работ ===
  const [photos, setPhotos] = useState([])
  const [photosVersion, setPhotosVersion] = useState(0)
  const [selectedFile, setSelectedFile] = useState(null)
  const [previewUrl, setPreviewUrl] = useState(null)
  const [isAddPhotoOpen, setIsAddPhotoOpen] = useState(false)
  const fileInputRef = useRef(null)

  useEffect(() => {
    loadStudioData()
  }, [])

  useEffect(() => {
    loadPhotos()
  }, [])

  const loadStudioData = async () => {
    try {
      setLoading(true)
      setError(null)
      const [profileData, servicesData] = await Promise.all([
        fetchMasterProfile(),
        fetchMasterServices(),
      ])

      setServices(Array.isArray(servicesData) ? servicesData : [])

      if (profileData?.work_hours) {
        const wh = profileData.work_hours
        if (Array.isArray(wh.work_days) && wh.work_days.length === 7) {
          if (typeof wh.work_days[0] === 'boolean') {
            setWorkDays(wh.work_days)
          } else {
            setWorkDays([1, 2, 3, 4, 5, 6, 7].map((d) => wh.work_days.includes(d)))
          }
        }
        setHours({
          start: wh.start_time || wh.start || '09:00',
          end: wh.end_time || wh.end || '20:00',
          lunchFrom: wh.lunch_start || wh.lunchFrom || '13:00',
          lunchTo: wh.lunch_end || wh.lunchTo || '14:00',
        })
      }

      if (profileData?.settings) {
        setAutoCancel(profileData.settings.auto_cancel ?? true)
        setCancelHours(profileData.settings.cancel_hours || '2')
        setOfferWaitlist(profileData.settings.offer_waitlist ?? true)
      }
    } catch (err) {
      console.error('Failed to load studio data:', err)
      setError('Не удалось загрузить данные')
    } finally {
      setLoading(false)
    }
  }

  const loadPhotos = async () => {
    try {
      const data = await fetchPhotos()
      setPhotos(Array.isArray(data) ? data : [])
    } catch (err) {
      console.error('Failed to load photos:', err)
      setPhotos([])
    }
  }

  const refreshServices = () => {
    setServicesVersion((v) => v + 1)
    loadStudioData()
  }

  const refreshPhotos = () => {
    setPhotosVersion((v) => v + 1)
    loadPhotos()
  }

  const toggleDay = (i) =>
    setWorkDays((days) => days.map((d, idx) => (idx === i ? !d : d)))

  const handleSaveSettings = async () => {
    try {
      await updateSettings({
        work_hours: {
          work_days: workDays,
          start_time: hours.start,
          end_time: hours.end,
          lunch_start: hours.lunchFrom,
          lunch_end: hours.lunchTo,
        },
        settings: {
          auto_cancel: autoCancel,
          cancel_hours: cancelHours,
          offer_waitlist: offerWaitlist,
        },
      })
      showToast('Настройки сохранены ✨')
    } catch (err) {
      showToast('Ошибка при сохранении настроек')
    }
  }

  const handlePickPhoto = () => {
    fileInputRef.current?.click()
  }

  const handleFileChange = (e) => {
    const file = e.target.files?.[0]
    if (!file) return
    if (!file.type.startsWith('image/')) {
      showToast('Выберите изображение')
      return
    }
    const reader = new FileReader()
    reader.onload = () => {
      setSelectedFile(file)
      setPreviewUrl(reader.result)
      setIsAddPhotoOpen(true)
    }
    reader.readAsDataURL(file)
    e.target.value = ''
  }

  const handleDeletePhoto = async (photoId) => {
    try {
      await deletePhoto(photoId)
      showToast('Фото удалено 🗑')
      refreshPhotos()
    } catch (err) {
      const message = err.response?.data?.message || err.message || 'Ошибка при удалении'
      showToast(message)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20">
        <p className="text-sm text-slate-400">Загрузка...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center py-20">
        <p className="text-sm text-red-500">{error}</p>
        <button
          onClick={loadStudioData}
          className="mt-4 rounded-xl bg-emerald-500 px-6 py-2.5 text-sm font-bold text-white transition active:scale-95"
        >
          Повторить
        </button>
      </div>
    )
  }

  return (
    <div className="pb-24">
      {/* ===== ПРИМЕРЫ РАБОТ ===== */}
      <section key={`photos-${photosVersion}`} className="mb-6">
        <h2 className="mb-2 px-1 text-base font-bold text-slate-800">Примеры работ</h2>

        <input
          ref={fileInputRef}
          type="file"
          accept="image/*"
          onChange={handleFileChange}
          className="hidden"
        />

        <div className="flex gap-3 overflow-x-auto pb-2">
          {photos.map((photo) => (
            <div
              key={photo.id}
              className="group relative h-32 w-40 shrink-0 overflow-hidden rounded-2xl bg-slate-200"
            >
              <img
                src={normalizePhotoUrl(photo.url)}
                alt="Работа мастера"
                className="h-full w-full object-cover"
              />
              <div className="pointer-events-none absolute inset-x-0 bottom-0 h-12 bg-gradient-to-t from-black/50 to-transparent" />
              <button
                onClick={() => handleDeletePhoto(photo.id)}
                className="absolute right-2 top-2 flex h-8 w-8 items-center justify-center rounded-full bg-black/60 text-white backdrop-blur-sm transition active:scale-90"
                aria-label="Удалить фото"
              >
                <Trash2 className="h-4 w-4" />
              </button>
            </div>
          ))}

          <button
            onClick={handlePickPhoto}
            className="flex h-32 w-40 shrink-0 flex-col items-center justify-center gap-2 rounded-2xl border-2 border-dashed border-slate-300 bg-white/50 text-slate-400 transition active:scale-95"
          >
            {photos.length === 0 ? (
              <>
                <Plus className="h-7 w-7" />
                <span className="text-xs font-semibold">Добавить фото</span>
              </>
            ) : (
              <>
                <Plus className="h-6 w-6" />
                <span className="text-[11px] font-semibold">Ещё</span>
              </>
            )}
          </button>
        </div>
      </section>

      <div className="space-y-6">
        {/* ===== УСЛУГИ ===== */}
        <section key={servicesVersion}>
          <h2 className="mb-2 text-base font-bold text-slate-800">Услуги</h2>
          <div className="rounded-2xl bg-white p-4 shadow-sm">
            <div className="space-y-2">
              {services.length === 0 ? (
                <p className="py-4 text-center text-sm text-slate-400">Услуги не добавлены</p>
              ) : (
                services.map((service) => (
                  <button
                    key={service.id}
                    onClick={() => setEditingService(service)}
                    className="flex w-full items-center gap-3 rounded-xl bg-slate-50 p-3 text-left transition active:scale-[0.98]"
                  >
                    <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-emerald-50 text-emerald-600">
                      <User className="h-5 w-5" />
                    </span>
                    <span className="min-w-0 flex-1">
                      <span className="block text-sm font-bold text-slate-900">{service.name}</span>
                      <span className="mt-0.5 flex items-center gap-1 text-xs text-slate-400">
                        <Clock className="h-3 w-3" />
                        {service.duration_min} мин
                      </span>
                    </span>
                    <span className="shrink-0 text-sm font-bold text-slate-900">
                      {rub(service.price)}
                    </span>
                  </button>
                ))
              )}
            </div>
            <button
              onClick={() => setIsAddServiceOpen(true)}
              className="mt-3 flex w-full items-center justify-center gap-2 rounded-xl border-2 border-dashed border-slate-200 py-3 font-semibold text-slate-500 transition active:scale-[0.98]"
            >
              <Plus className="h-4 w-4" />
              Добавить услугу
            </button>
          </div>
        </section>

        {/* ===== ГРАФИК РАБОТЫ ===== */}
        <section>
          <h2 className="mb-2 text-base font-bold text-slate-800">График работы</h2>
          <div className="rounded-2xl bg-white p-4 shadow-sm">
            <div className="flex flex-wrap gap-2">
              {WEEK_TAGS.map((tag, i) => {
                const isWork = workDays[i]
                return (
                  <button
                    key={tag}
                    onClick={() => toggleDay(i)}
                    className={`h-10 rounded-full px-4 text-sm font-bold transition ${
                      isWork
                        ? 'bg-emerald-500 text-white shadow-md shadow-emerald-500/25 active:scale-95'
                        : 'bg-slate-100 text-slate-400 active:scale-95'
                    }`}
                  >
                    {tag}
                  </button>
                )
              })}
            </div>

            <div className="mt-4 grid grid-cols-2 gap-3">
              <label className="block">
                <span className="mb-1.5 block text-xs font-semibold text-slate-400">Начало дня</span>
                <input
                  type="time"
                  value={hours.start}
                  onChange={(e) => setHours((h) => ({ ...h, start: e.target.value }))}
                  className="w-full rounded-xl border-0 bg-slate-100 px-4 py-3 text-center font-semibold text-slate-900 outline-none transition focus:bg-emerald-50 focus:ring-2 focus:ring-emerald-500"
                />
              </label>
              <label className="block">
                <span className="mb-1.5 block text-xs font-semibold text-slate-400">Конец дня</span>
                <input
                  type="time"
                  value={hours.end}
                  onChange={(e) => setHours((h) => ({ ...h, end: e.target.value }))}
                  className="w-full rounded-xl border-0 bg-slate-100 px-4 py-3 text-center font-semibold text-slate-900 outline-none transition focus:bg-emerald-50 focus:ring-2 focus:ring-emerald-500"
                />
              </label>
            </div>

            <div className="mt-3">
              <span className="mb-1.5 block text-xs font-semibold text-slate-400">Обед</span>
              <div className="flex items-center gap-3">
                <input
                  type="time"
                  value={hours.lunchFrom}
                  onChange={(e) => setHours((h) => ({ ...h, lunchFrom: e.target.value }))}
                  className="flex-1 rounded-xl border-0 bg-slate-100 px-4 py-3 text-center font-semibold text-slate-900 outline-none transition focus:bg-emerald-50 focus:ring-2 focus:ring-emerald-500"
                />
                <span className="text-slate-300">—</span>
                <input
                  type="time"
                  value={hours.lunchTo}
                  onChange={(e) => setHours((h) => ({ ...h, lunchTo: e.target.value }))}
                  className="flex-1 rounded-xl border-0 bg-slate-100 px-4 py-3 text-center font-semibold text-slate-900 outline-none transition focus:bg-emerald-50 focus:ring-2 focus:ring-emerald-500"
                />
              </div>
            </div>
          </div>
        </section>

        {/* ===== ПРАВИЛА ЗАПИСИ ===== */}
        <section>
          <h2 className="mb-2 text-base font-bold text-slate-800">Правила записи</h2>
          <div className="rounded-2xl bg-white p-4 shadow-sm">
            <div className="flex items-center justify-between gap-3">
              <p className="text-sm font-semibold text-slate-900">Авто-отмена без подтверждения</p>
              <Toggle checked={autoCancel} onChange={setAutoCancel} />
            </div>

            {autoCancel && (
              <div className="mt-3 animate-fade-in rounded-xl bg-slate-50 p-3">
                <p className="text-xs leading-relaxed text-slate-500">
                  Отменять запись, если клиент не подтвердил её за{' '}
                  <b className="text-slate-700">{cancelHours} {cancelHours === '1' ? 'час' : 'часа'}</b>
                </p>
                <div className="mt-2 flex gap-2">
                  {['1', '2', '4'].map((h) => (
                    <button
                      key={h}
                      onClick={() => setCancelHours(h)}
                      className={`flex-1 rounded-xl py-2 text-xs font-bold transition ${
                        cancelHours === h
                          ? 'bg-emerald-600 text-white shadow-md shadow-emerald-600/20'
                          : 'bg-white text-slate-500 active:scale-95'
                      }`}
                    >
                      {h} {h === '1' ? 'час' : 'часа'}
                    </button>
                  ))}
                </div>
              </div>
            )}

            <div className="mt-4 border-t border-slate-100 pt-4">
              <div className="flex items-center justify-between gap-3">
                <div className="min-w-0">
                  <p className="text-sm font-semibold text-slate-900">Предлагать окно в лист ожидания</p>
                  <p className="mt-0.5 text-xs text-slate-400">
                    Бот сам предложит освободившееся время клиентам из очереди
                  </p>
                </div>
                <Toggle checked={offerWaitlist} onChange={setOfferWaitlist} />
              </div>
            </div>

            <button
              onClick={handleSaveSettings}
              className="mt-4 w-full rounded-xl bg-emerald-500 py-3 font-bold text-white transition active:scale-[0.98]"
            >
              Сохранить настройки
            </button>
          </div>
        </section>
      </div>

      {/* Модалки */}
      {isAddServiceOpen && (
        <AddServiceSheet
          onClose={() => setIsAddServiceOpen(false)}
          onCreated={refreshServices}
        />
      )}
      {editingService && (
        <EditServiceSheet
          service={editingService}
          onClose={() => setEditingService(null)}
          onSaved={refreshServices}
        />
      )}
      {isAddPhotoOpen && (
        <AddPhotoSheet
          file={selectedFile}
          preview={previewUrl}
          onClose={() => {
            setIsAddPhotoOpen(false)
            setSelectedFile(null)
            setPreviewUrl(null)
          }}
          onUploaded={refreshPhotos}
        />
      )}
    </div>
  )
}