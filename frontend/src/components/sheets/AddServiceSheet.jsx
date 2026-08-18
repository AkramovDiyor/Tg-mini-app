import { useState } from 'react'
import { SheetShell } from './SheetShell'
import { useBookingStore } from '../../store/bookingStore'
import { createService } from '../../services/api'

const DURATION_OPTIONS = [
  { label: '15 мин', minutes: 15 },
  { label: '30 мин', minutes: 30 },
  { label: '45 мин', minutes: 45 },
  { label: '1 ч',    minutes: 60 },
  { label: '1.5 ч',  minutes: 90 },
  { label: '2 ч',    minutes: 120 },
]

export function AddServiceSheet({ onClose, onCreated }) {
  const showToast = useBookingStore((s) => s.showToast)
  const [saving, setSaving] = useState(false)

  const [name, setName] = useState('')
  const [minutes, setMinutes] = useState(45)
  const [price, setPrice] = useState('')

  const isValid = name.trim().length > 0 && price.trim().length > 0 && Number(price) > 0

  const handleSave = async () => {
    if (!isValid) return
    try {
      setSaving(true)
      await createService({
        name: name.trim(),
        duration_min: minutes,
        price: Number(price),
      })
      onClose()
      showToast(`Услуга «${name}» добавлена 🎉`)
      onCreated?.() // Перезагружаем список услуг
    } catch (err) {
      showToast('Ошибка при создании услуги')
    } finally {
      setSaving(false)
    }
  }

  return (
    <SheetShell onClose={onClose}>
      <h2 className="mb-5 text-xl font-bold text-slate-900">Новая услуга</h2>

      <div>
        <label className="mb-2 block text-sm font-semibold text-slate-500">Название</label>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Например, Стрижка ножницами"
          className="w-full rounded-xl border-0 bg-slate-100 px-4 py-3.5 text-base font-medium text-slate-900 outline-none transition placeholder:text-slate-400 focus:bg-emerald-50 focus:ring-2 focus:ring-emerald-500"
        />
      </div>

      <div className="mt-4">
        <label className="mb-2 block text-sm font-semibold text-slate-500">Длительность</label>
        <div className="flex flex-wrap gap-2">
          {DURATION_OPTIONS.map((opt) => {
            const isActive = minutes === opt.minutes
            return (
              <button
                key={opt.minutes}
                onClick={() => setMinutes(opt.minutes)}
                className={`rounded-xl px-4 py-2 text-sm transition ${
                  isActive
                    ? 'bg-emerald-500 font-bold text-white shadow-md shadow-emerald-500/25'
                    : 'bg-slate-100 font-medium text-slate-600 active:scale-95'
                }`}
              >
                {opt.label}
              </button>
            )
          })}
        </div>
      </div>

      <div className="mt-4">
        <label className="mb-2 block text-sm font-semibold text-slate-500">Стоимость</label>
        <div className="relative">
          <input
            type="text"
            inputMode="numeric"
            value={price}
            onChange={(e) => setPrice(e.target.value.replace(/\D/g, ''))}
            placeholder="0"
            className="w-full rounded-xl border-0 bg-slate-100 px-4 py-3.5 pr-10 text-base font-bold text-slate-900 outline-none transition placeholder:text-slate-400 focus:bg-emerald-50 focus:ring-2 focus:ring-emerald-500"
          />
          <span className="pointer-events-none absolute right-4 top-1/2 -translate-y-1/2 font-bold text-slate-400">
            ₽
          </span>
        </div>
      </div>

      <button
        onClick={handleSave}
        disabled={!isValid || saving}
        className={`mt-6 w-full rounded-xl py-4 font-bold transition ${
          isValid && !saving
            ? 'bg-emerald-500 text-white shadow-lg shadow-emerald-500/25 active:scale-[0.98]'
            : 'bg-slate-200 text-slate-400'
        }`}
      >
        {saving ? 'Сохранение...' : 'Сохранить услугу'}
      </button>
    </SheetShell>
  )
}