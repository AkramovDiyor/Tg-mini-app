import { useState } from 'react'
import { Upload, Loader2 } from 'lucide-react'
import { SheetShell } from './SheetShell'
import { useBookingStore } from '../../store/bookingStore'
import { uploadPhoto } from '../../services/api'

export function AddPhotoSheet({ file, preview, onClose, onUploaded }) {
  const showToast = useBookingStore((s) => s.showToast)
  const [uploading, setUploading] = useState(false)

  const handleUpload = async () => {
    if (!file || uploading) return
    try {
      setUploading(true)
      await uploadPhoto(file)
      showToast('Фото загружено 📸')
      onUploaded?.()
      onClose()
    } catch (err) {
      const message = err.response?.data?.message || err.message || 'Ошибка при загрузке'
      showToast(message)
    } finally {
      setUploading(false)
    }
  }

  return (
    <SheetShell onClose={uploading ? () => {} : onClose} closableOnBackdrop={!uploading}>
      <h2 className="mb-5 text-xl font-bold text-slate-900">Новое фото</h2>

      {/* Предпросмотр */}
      <div className="overflow-hidden rounded-2xl bg-slate-100">
        <img
          src={preview}
          alt="Предпросмотр"
          className="h-64 w-full object-cover"
        />
      </div>

      {/* Кнопка загрузки */}
      <button
        onClick={handleUpload}
        disabled={uploading}
        className="mt-6 flex w-full items-center justify-center gap-2 rounded-xl bg-emerald-500 py-4 font-bold text-white shadow-lg shadow-emerald-500/25 transition active:scale-[0.98] disabled:opacity-50"
      >
        {uploading ? (
          <>
            <Loader2 className="h-5 w-5 animate-spin" />
            Загрузка...
          </>
        ) : (
          <>
            <Upload className="h-5 w-5" />
            Загрузить
          </>
        )}
      </button>

      {/* Отмена */}
      <button
        onClick={onClose}
        disabled={uploading}
        className="mt-2 w-full py-4 font-bold text-slate-500 transition active:scale-[0.98] disabled:opacity-50"
      >
        Отмена
      </button>
    </SheetShell>
  )
}