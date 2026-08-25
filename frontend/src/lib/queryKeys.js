// 🔥 OPTIMIZED: единый словарь ключей кэша — invalidate без магических строк
export const queryKeys = {
  me: ['me'],
  masterInfo: (inviteLink) => ['masterInfo', inviteLink],
  services: (inviteLink) => ['services', inviteLink],
  slots: (inviteLink, date, serviceId) => ['slots', inviteLink, date, serviceId],
  clientBookings: ['clientBookings'],
  todaySchedule: ['master', 'today'],
  waitlist: ['master', 'waitlist'],
  masterProfile: ['master', 'profile'],
  masterServices: ['master', 'services'],
  photos: ['master', 'photos'],
}
