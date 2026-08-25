import { useMutation, useQueryClient } from '@tanstack/react-query'
import { queryKeys } from '../lib/queryKeys'
import { bookSlot, cancelBooking, getApiErrorMessage } from '../services/api'
import { useBookingStore } from '../store/bookingStore'

/**
 * @param {{ inviteLink: string }} opts
 */
export function useBookSlotMutation({ inviteLink }) {
  const queryClient = useQueryClient()
  const showToast = useBookingStore((s) => s.showToast)
  const closeSlotSheet = useBookingStore((s) => s.closeSlotSheet)

  return useMutation({
    mutationFn: bookSlot,
    // 🔥 OPTIMIZED: слот сразу «занят», бар записи появляется без ожидания сети
    onMutate: async (vars) => {
      const slotsKey = queryKeys.slots(inviteLink, vars.date, vars.service_id)
      await queryClient.cancelQueries({ queryKey: slotsKey })
      await queryClient.cancelQueries({ queryKey: queryKeys.clientBookings })

      const prevSlots = queryClient.getQueryData(slotsKey)
      const prevBookings = queryClient.getQueryData(queryKeys.clientBookings)

      queryClient.setQueryData(slotsKey, (old) => {
        if (!Array.isArray(old)) return old
        return old.map((slot) =>
          slot.start_time === vars.start_time ? { ...slot, status: 'booked' } : slot,
        )
      })

      const optimistic = {
        booking_id: -Date.now(),
        client_name: vars.name,
        service_name: vars.service_name || 'Услуга',
        service_duration: vars.service_duration || 0,
        service_price: vars.price,
        master_name: '',
        master_address: '',
        start_time: vars.start_time,
        end_time: vars.start_time,
        status: 'active',
        created_at: new Date().toISOString(),
      }

      queryClient.setQueryData(queryKeys.clientBookings, (old) =>
        Array.isArray(old) ? [optimistic, ...old] : [optimistic],
      )

      return { prevSlots, prevBookings, slotsKey }
    },
    onError: (err, _vars, ctx) => {
      if (ctx?.slotsKey && ctx.prevSlots !== undefined) {
        queryClient.setQueryData(ctx.slotsKey, ctx.prevSlots)
      }
      if (ctx?.prevBookings !== undefined) {
        queryClient.setQueryData(queryKeys.clientBookings, ctx.prevBookings)
      }
      showToast(getApiErrorMessage(err, 'Ошибка при бронировании'))
    },
    onSuccess: () => {
      closeSlotSheet()
      showToast('Запись подтверждена')
    },
    onSettled: (_data, _err, vars) => {
      queryClient.invalidateQueries({
        queryKey: queryKeys.slots(inviteLink, vars.date, vars.service_id),
      })
      queryClient.invalidateQueries({ queryKey: queryKeys.clientBookings })
    },
  })
}

export function useCancelBookingMutation() {
  const queryClient = useQueryClient()
  const showToast = useBookingStore((s) => s.showToast)
  const closeDetails = useBookingStore((s) => s.closeDetails)

  return useMutation({
    mutationFn: ({ bookingId }) => cancelBooking(bookingId),
    onMutate: async ({ booking }) => {
      await queryClient.cancelQueries({ queryKey: queryKeys.clientBookings })
      const prevBookings = queryClient.getQueryData(queryKeys.clientBookings)

      queryClient.setQueryData(queryKeys.clientBookings, (old) => {
        if (!Array.isArray(old)) return old
        return old.filter((b) => b.booking_id !== booking.booking_id)
      })

      return { prevBookings }
    },
    onError: (err, _vars, ctx) => {
      if (ctx?.prevBookings !== undefined) {
        queryClient.setQueryData(queryKeys.clientBookings, ctx.prevBookings)
      }
      showToast(getApiErrorMessage(err, 'Ошибка при отмене записи'))
    },
    onSuccess: () => {
      closeDetails()
      showToast('Запись отменена. Время освобождено для других')
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.clientBookings })
      queryClient.invalidateQueries({ queryKey: ['slots'] })
    },
  })
}
