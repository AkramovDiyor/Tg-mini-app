import { useState } from 'react'
import { SheetShell } from './SheetShell'
import { useUpdateServiceMutation, useDeleteServiceMutation } from '../../hooks/useMasterMutations'

const DURATION_OPTIONS = [
  { label: '15 мин', minutes: 15 },
  { label: '30 мин', minutes: 30 },
  { label: '45 мин', minutes: 45 },
  { label: '1 ч', minutes: 60 },
  { label: '1.5 ч', minutes: 90 },
  { label: '2 ч', minutes: 120 },
]

export function EditServiceSheet({ service, onClose }) {
  const updateMutation = useUpdateServiceMutation()
  const deleteMutation = useDeleteServiceMutation()
  const [name, setName] = useState(service?.name || '')
  const [minutes, setMinutes] = useState(service?.duration_min || 45)
  const [price, setPrice] = useState(service?.price != null ? String(service.price) : '')

  const isValid = name.trim().length > 0 && price.trim().length > 0 && Number(price) > 0
  const isBusy = updateMutation.isPending || deleteMutation.isPending

  const handleSave = () => {
    if (!isValid || isBusy) return
    updateMutation.mutate(
      {
        serviceId: service.id,
        payload: { name: name.trim(), duration_min: minutes, price: Number(price) },
      },
      { onSuccess: () => onClose() },
    )
  }

  const handleDelete = () => {
    if (isBusy) return
    deleteMutation.mutate(
      { serviceId: service.id, name: service.name },
      { onSuccess: () => onClose() },
    )
  }

  return (
    <SheetShell onClose={isBusy ? () => {} : onClose} closableOnBackdrop={!isBusy}>
      <h2 className="mb-5 text-xl font-bold text-slate-900">Редактировать услугу</h2>
      <div>
        <label className="mb-2 block text-sm font-semibold text-slate-500">Название</label>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          disabled={isBusy}
          placeholder="Например, Стрижка ножницами"
          className="w-full rounded-xl border-0 bg-slate-100 px-4 py-3.5 text-base font-medium text-slate-900 outline-none transition placeholder:text-slate-400 focus:bg-emerald-50 focus:ring-2 focus:ring-emerald-500 disabled:opacity-50"
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
                type="button"
                onClick={() => setMinutes(opt.minutes)}
                disabled={isBusy}
                className={`rounded-xl px-4 py-2 text-sm transition ${
                  isActive
                    ? 'bg-emerald-500 font-bold text-white shadow-md shadow-emerald-500/25'
                    : 'bg-slate-100 font-medium text-slate-600 active:scale-95'
                } disabled:opacity-50`}
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
            disabled={isBusy}
            placeholder="0"
            className="w-full rounded-xl border-0 bg-slate-100 px-4 py-3.5 pr-10 text-base font-bold text-slate-900 outline-none transition placeholder:text-slate-400 focus:bg-emerald-50 focus:ring-2 focus:ring-emerald-500 disabled:opacity-50"
          />
          <span className="pointer-events-none absolute right-4 top-1/2 -translate-y-1/2 font-bold text-slate-400">
            ₽
          </span>
        </div>
      </div>
      <button
        type="button"
        onClick={handleSave}
        disabled={!isValid || isBusy}
        className={`mt-6 w-full rounded-xl py-4 font-bold transition ${
          isValid && !isBusy
            ? 'bg-emerald-500 text-white shadow-lg shadow-emerald-500/25 active:scale-[0.98]'
            : 'bg-slate-200 text-slate-400'
        }`}
      >
        {updateMutation.isPending ? 'Сохранение...' : 'Сохранить'}
      </button>
      <button
        type="button"
        onClick={handleDelete}
        disabled={isBusy}
        className="mt-2 w-full py-4 font-bold text-red-500 transition active:scale-[0.98] disabled:opacity-50"
      >
        {deleteMutation.isPending ? 'Удаление...' : 'Удалить услугу'}
      </button>
    </SheetShell>
  )
}
