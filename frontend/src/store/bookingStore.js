import { create } from 'zustand'
import { toISO, TODAY } from '../lib/dates'

export const useBookingStore = create((set) => ({
  service: null,
  selectedDate: toISO(TODAY),
  myBookings: {},
  waitlist: {},
  toast: null,

  pickService: (service) => set({ service, selectedDate: toISO(TODAY) }),
  setDate: (iso) => set({ selectedDate: iso }),
  bookSlot: (key) => set((s) => ({ myBookings: { ...s.myBookings, [key]: true } })),
  joinWaitlist: (iso) => set((s) => ({ waitlist: { ...s.waitlist, [iso]: true } })),
  showToast: (text) => set({ toast: { id: Date.now(), text } }),
  hideToast: () => set({ toast: null }),
}))