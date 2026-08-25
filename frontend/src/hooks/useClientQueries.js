import { keepPreviousData, useQuery } from '@tanstack/react-query'
import { queryKeys } from '../lib/queryKeys'
import {
  fetchMasterInfo,
  fetchServices,
  fetchSlots,
  fetchClientBookings,
} from '../services/api'

/**
 * @param {string | undefined} inviteLink
 */
export function useMasterInfoQuery(inviteLink) {
  return useQuery({
    queryKey: queryKeys.masterInfo(inviteLink),
    queryFn: () => fetchMasterInfo(inviteLink),
    enabled: Boolean(inviteLink),
  })
}

/**
 * @param {string | undefined} inviteLink
 */
export function useServicesQuery(inviteLink) {
  return useQuery({
    queryKey: queryKeys.services(inviteLink),
    queryFn: () => fetchServices(inviteLink),
    enabled: Boolean(inviteLink),
  })
}

/**
 * @param {string | undefined} inviteLink
 * @param {string | undefined} date
 * @param {number | undefined} serviceId
 */
export function useSlotsQuery(inviteLink, date, serviceId) {
  return useQuery({
    queryKey: queryKeys.slots(inviteLink, date, serviceId),
    queryFn: () => fetchSlots(date, serviceId, inviteLink),
    enabled: Boolean(inviteLink && date && serviceId),
    staleTime: 15_000,
    placeholderData: keepPreviousData,
  })
}

export function useClientBookingsQuery() {
  return useQuery({
    queryKey: queryKeys.clientBookings,
    queryFn: fetchClientBookings,
    staleTime: 20_000,
  })
}
