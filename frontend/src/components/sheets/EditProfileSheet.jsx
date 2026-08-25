import { useState } from 'react'
import { MapPin } from 'lucide-react'
import { SheetShell } from './SheetShell'
import { useUpdateProfileMutation } from '../../hooks/useMasterMutations'

export function EditProfileSheet({ profile, onClose }) {
  const updateMutation = useUpdateProfileMutation()
  const [name, setName] = useState(profile?.name || '')
  const [bio, setBio] = useState(profile?.bio || '')
  const [address, setAddress] = useState(profile?.address || '')

  const handleSave = () => {
    updateMutation.mutate({ name, bio, address }, { onSuccess: () => onClose() })
  }

  return (
    <SheetShell onClose={onClose}>
      <h2 className="mb-5 text-xl font-bold text-slate-900">Редактировать профиль</h2>
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
      <div className="mt-3">
        <label className="mb-2 block text-sm font-semibold text-slate-500">Специализация</label>
        <input
          type="text"
          value={bio}
          onChange={(e) => setBio(e.target.value)}
          placeholder="Например, Барбер · 6 лет опыта"
          className="w-full rounded-xl border-0 bg-slate-100 px-4 py-3.5 text-base font-medium text-slate-900 outline-none transition placeholder:text-slate-400 focus:bg-emerald-50 focus:ring-2 focus:ring-emerald-500"
        />
      </div>
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
      <button
        type="button"
        onClick={handleSave}
        disabled={updateMutation.isPending}
        className="mt-6 w-full rounded-xl bg-emerald-500 py-4 font-bold text-white shadow-lg shadow-emerald-500/25 transition active:scale-[0.98] disabled:opacity-50"
      >
        {updateMutation.isPending ? 'Сохранение...' : 'Сохранить'}
      </button>
    </SheetShell>
  )
}
