import { useQuery } from '@tanstack/react-query'
import { queryKeys } from '../lib/queryKeys'
import { fetchUserIdentity } from '../services/api'

export function useMeQuery() {
  return useQuery({
    queryKey: queryKeys.me,
    queryFn: fetchUserIdentity,
    staleTime: 5 * 60_000,
  })
}
