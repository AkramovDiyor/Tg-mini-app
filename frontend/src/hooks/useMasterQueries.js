import { useQuery } from '@tanstack/react-query'
import { queryKeys } from '../lib/queryKeys'
import {
  fetchTodaySchedule,
  fetchWaitlist,
  fetchMasterProfile,
  fetchMasterServices,
  fetchPhotos,
} from '../services/api'

export function useTodayScheduleQuery() {
  return useQuery({
    queryKey: queryKeys.todaySchedule,
    queryFn: fetchTodaySchedule,
    staleTime: 20_000,
  })
}

export function useWaitlistQuery() {
  return useQuery({
    queryKey: queryKeys.waitlist,
    queryFn: fetchWaitlist,
  })
}

export function useMasterProfileQuery() {
  return useQuery({
    queryKey: queryKeys.masterProfile,
    queryFn: fetchMasterProfile,
  })
}

export function useMasterServicesQuery() {
  return useQuery({
    queryKey: queryKeys.masterServices,
    queryFn: fetchMasterServices,
  })
}

export function usePhotosQuery() {
  return useQuery({
    queryKey: queryKeys.photos,
    queryFn: fetchPhotos,
  })
}
