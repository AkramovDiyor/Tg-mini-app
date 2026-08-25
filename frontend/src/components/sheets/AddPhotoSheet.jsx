import { Upload, Loader2 } from 'lucide-react'
import { SheetShell } from './SheetShell'
import { useUploadPhotoMutation } from '../../hooks/useMasterMutations'

export function AddPhotoSheet({ file, preview, onClose }) {
  const uploadMutation = useUploadPhotoMutation()

  const handleUpload = () => {
    if (!file || uploadMutation.isPending) return
    uploadMutation.mutate(file, { onSuccess: () => onClose() })
  }

  return (
    <SheetShell
      onClose={uploadMutation.isPending ? () => {} : onClose}
      closableOnBackdrop={!uploadMutation.isPending}
    >
      <h2 className="mb-5 text-xl font-bold text-slate-900">Новое фото</h2>
      <div className="overflow-hidden rounded-2xl bg-slate-100">
        <img src={preview} alt="Предпросмотр" className="h-64 w-full object-cover" />
      </div>
      <button
        type="button"
        onClick={handleUpload}
        disabled={uploadMutation.isPending}
        className="mt-6 flex w-full items-center justify-center gap-2 rounded-xl bg-emerald-500 py-4 font-bold text-white shadow-lg shadow-emerald-500/25 transition active:scale-[0.98] disabled:opacity-50"
      >
        {uploadMutation.isPending ? (
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
      <button
        type="button"
        onClick={onClose}
        disabled={uploadMutation.isPending}
        className="mt-2 w-full py-4 font-bold text-slate-500 transition active:scale-[0.98] disabled:opacity-50"
      >
        Отмена
      </button>
    </SheetShell>
  )
}
