import { useState } from 'react'
import { MapPin } from 'lucide-react'
import { SheetShell } from './SheetShell'
import { useBookingStore } from '../../store/bookingStore'

export function EditProfileSheet({ onClose }) {
  const showToast = useBookingStore((s) => s.showToast)

  const [name, setName] = useState('Педро Барбер')
  const [specialty, setSpecialty] = useState('Барбер · 6 лет опыта')
  const [address, setAddress] = useState('ул. Центральная, 1')

  const handleSave = () => {
    onClose()
    showToast('Профиль обновлён ✨')
  }

  return (
    <SheetShell onClose={onClose}>
      <h2 className="mb-5 text-xl font-bold text-slate-900">Редактировать профиль</h2>

      {/* Имя */}
      <div>
        <label className="mb-2 block text-sm font-semibold text-slate-500">Имя</label>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Имя и фамилия"
          className="w-full rounded-xl border-0 bg-slate-100 px-4 py-3.5 text-base font-medium text-slate-900 outline-none transition placeholder:text-slate-400 focus:bg-emerald-50 focus:ring-2 focus:ring-emerald-500"
        />
      </div>

      {/* Специализация */}
      <div className="mt-3">
        <label className="mb-2 block text-sm font-semibold text-slate-500">Специализация</label>
        <input
          type="text"
          value={specialty}
          onChange={(e) => setSpecialty(e.target.value)}
          placeholder="Например, Барбер · 6 лет опыта"
          className="w-full rounded-xl border-0 bg-slate-100 px-4 py-3.5 text-base font-medium text-slate-900 outline-none transition placeholder:text-slate-400 focus:bg-emerald-50 focus:ring-2 focus:ring-emerald-500"
        />
      </div>

      {/* Адрес с иконкой MapPin внутри */}
      <div className="mt-3">
        <label className="mb-2 block text-sm font-semibold text-slate-500">Адрес студии</label>
        <div className="relative">
          <MapPin className="pointer-events-none absolute left-4 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
          <input
            type="text"
            value={address}
            onChange={(e) => setAddress(e.target.value)}
            placeholder="Улица, дом"
            className="w-full rounded-xl border-0 bg-slate-100 py-3.5 pl-11 pr-4 text-base font-medium text-slate-900 outline-none transition placeholder:text-slate-400 focus:bg-emerald-50 focus:ring-2 focus:ring-emerald-500"
          />
        </div>
      </div>

      {/* Кнопка сохранения */}
      <button
        onClick={handleSave}
        className="mt-6 w-full rounded-xl bg-emerald-500 py-4 font-bold text-white shadow-lg shadow-emerald-500/25 transition active:scale-[0.98]"
      >
        Сохранить
      </button>
    </SheetShell>
  )
}