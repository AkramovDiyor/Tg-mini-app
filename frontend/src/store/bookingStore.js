import { create } from 'zustand'
import { toISO, TODAY } from '../lib/dates'

// 🔥 OPTIMIZED: только client UI-state. Серверные данные — в TanStack Query.
export const useBookingStore = create((set) => ({
  service: null,
  selectedDate: toISO(TODAY),
  screen: 'services',
  sheetSlot: null,
  detailsBooking: null,
  toast: null,

  pickService: (service) =>
    set({ service, selectedDate: toISO(TODAY), screen: 'booking' }),
  setDate: (iso) => set({ selectedDate: iso }),
  goServices: () => set({ screen: 'services' }),
  openSlotSheet: (slot) => set({ sheetSlot: slot }),
  closeSlotSheet: () => set({ sheetSlot: null }),
  openDetails: (booking) => set({ detailsBooking: booking }),
  closeDetails: () => set({ detailsBooking: null }),
  showToast: (text) => set({ toast: { id: Date.now(), text } }),
  hideToast: () => set({ toast: null }),
}))
