import { useMutation, useQueryClient } from '@tanstack/react-query'
import { queryKeys } from '../lib/queryKeys'
import {
  updateSettings,
  createService,
  updateService,
  deleteService,
  updateMasterProfile,
  uploadPhoto,
  deletePhoto,
  getApiErrorMessage,
} from '../services/api'
import { useBookingStore } from '../store/bookingStore'

export function useUpdateSettingsMutation() {
  const queryClient = useQueryClient()
  const showToast = useBookingStore((s) => s.showToast)

  return useMutation({
    mutationFn: updateSettings,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.masterProfile })
      showToast('Настройки сохранены')
    },
    onError: (err) => {
      showToast(getApiErrorMessage(err, 'Ошибка при сохранении настроек'))
    },
  })
}

export function useCreateServiceMutation() {
  const queryClient = useQueryClient()
  const showToast = useBookingStore((s) => s.showToast)

  return useMutation({
    mutationFn: createService,
    onSuccess: (_data, vars) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.masterServices })
      showToast(`Услуга «${vars.name}» добавлена`)
    },
    onError: (err) => {
      showToast(getApiErrorMessage(err, 'Ошибка при создании услуги'))
    },
  })
}

export function useUpdateServiceMutation() {
  const queryClient = useQueryClient()
  const showToast = useBookingStore((s) => s.showToast)

  return useMutation({
    mutationFn: ({ serviceId, payload }) => updateService(serviceId, payload),
    onSuccess: (_data, vars) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.masterServices })
      queryClient.invalidateQueries({ queryKey: queryKeys.todaySchedule })
      showToast(`Услуга «${vars.payload.name}» обновлена`)
    },
    onError: (err) => {
      showToast(getApiErrorMessage(err, 'Ошибка при обновлении'))
    },
  })
}

export function useDeleteServiceMutation() {
  const queryClient = useQueryClient()
  const showToast = useBookingStore((s) => s.showToast)

  return useMutation({
    mutationFn: ({ serviceId }) => deleteService(serviceId),
    onSuccess: (_data, vars) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.masterServices })
      showToast(`Услуга «${vars.name}» удалена`)
    },
    onError: (err) => {
      showToast(getApiErrorMessage(err, 'Ошибка при удалении'))
    },
  })
}

export function useUpdateProfileMutation() {
  const queryClient = useQueryClient()
  const showToast = useBookingStore((s) => s.showToast)

  return useMutation({
    mutationFn: updateMasterProfile,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.masterProfile })
      queryClient.invalidateQueries({ queryKey: queryKeys.me })
      showToast('Профиль обновлён')
    },
    onError: (err) => {
      showToast(getApiErrorMessage(err, 'Ошибка при сохранении профиля'))
    },
  })
}

export function useUploadPhotoMutation() {
  const queryClient = useQueryClient()
  const showToast = useBookingStore((s) => s.showToast)

  return useMutation({
    mutationFn: uploadPhoto,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.photos })
      showToast('Фото загружено')
    },
    onError: (err) => {
      showToast(getApiErrorMessage(err, 'Ошибка при загрузке'))
    },
  })
}

export function useDeletePhotoMutation() {
  const queryClient = useQueryClient()
  const showToast = useBookingStore((s) => s.showToast)

  return useMutation({
    mutationFn: deletePhoto,
    onMutate: async (photoId) => {
      await queryClient.cancelQueries({ queryKey: queryKeys.photos })
      const prev = queryClient.getQueryData(queryKeys.photos)
      queryClient.setQueryData(queryKeys.photos, (old) =>
        Array.isArray(old) ? old.filter((p) => p.id !== photoId) : old,
      )
      return { prev }
    },
    onError: (err, _id, ctx) => {
      if (ctx?.prev) queryClient.setQueryData(queryKeys.photos, ctx.prev)
      showToast(getApiErrorMessage(err, 'Ошибка при удалении'))
    },
    onSuccess: () => {
      showToast('Фото удалено')
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.photos })
    },
  })
}
