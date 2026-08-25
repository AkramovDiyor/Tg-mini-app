import { QueryClient } from '@tanstack/react-query'

// 🔥 OPTIMIZED: Mini App часто теряет фокус — refetchOnWindowFocus даёт лишний трафик
export function createQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 60_000,
        gcTime: 5 * 60_000,
        retry: 2,
        refetchOnWindowFocus: false,
        refetchOnReconnect: true,
      },
      mutations: {
        retry: 0,
      },
    },
  })
}
