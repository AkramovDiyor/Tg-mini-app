import { useEffect, useRef } from 'react'

export function useDragScroll() {
  const ref = useRef(null)
  const state = useRef({ down: false, moved: 0, startX: 0, startLeft: 0 })

  useEffect(() => {
    const el = ref.current
    if (!el) return

    const onDown = (e) => {
      if (e.pointerType !== 'mouse') return
      state.current = { down: true, moved: 0, startX: e.clientX, startLeft: el.scrollLeft }
    }
    const onMove = (e) => {
      const s = state.current
      if (!s.down) return
      const dx = e.clientX - s.startX
      s.moved = Math.max(s.moved, Math.abs(dx))
      if (s.moved > 4) el.scrollLeft = s.startLeft - dx
    }
    const onUp = () => {
      const s = state.current
      if (s.down && s.moved > 4) {
        el.addEventListener(
          'click',
          (e) => { e.preventDefault(); e.stopPropagation() },
          { capture: true, once: true },
        )
      }
      s.down = false
    }

    el.addEventListener('pointerdown', onDown)
    window.addEventListener('pointermove', onMove)
    window.addEventListener('pointerup', onUp)
    return () => {
      el.removeEventListener('pointerdown', onDown)
      window.removeEventListener('pointermove', onMove)
      window.removeEventListener('pointerup', onUp)
    }
  }, [])

  return ref
}