import { CalendarDays, Users, Store, User } from 'lucide-react'

export const TODAY_APPOINTMENTS = [
  { time: '10:00', client: 'Алексей П.', service: 'Мужская стрижка',  price: 1500, status: 'confirmed' },
  { time: '13:00', client: 'Мария К.',  service: 'Стрижка машинкой', price: 800,  status: 'pending'   },
  { time: '15:30', client: null,        service: null,               price: 0,    status: 'free'      },
]

export const WAITLIST = [
  { id: 1, name: 'Алексей М.', initials: 'АМ', time: '13:00–14:00', joinedAt: '09:42' },
  { id: 2, name: 'Дмитрий С.', initials: 'ДС', time: '16:00–17:00', joinedAt: '10:15' },
]

export const DOT_COLORS = {
  confirmed: 'bg-emerald-500',
  pending:   'bg-amber-400',
  free:      'bg-slate-200',
}

export const AVATAR_GRADIENTS = [
  'from-amber-400 to-rose-500',
  'from-sky-400 to-indigo-600',
  'from-emerald-400 to-teal-600',
]

// 4 вкладки вместо 3
export const MASTER_TABS = [
  { id: 'today',   icon: CalendarDays },
  { id: 'queue',   icon: Users        },
  { id: 'studio',  icon: Store        }, // ← НОВАЯ
  { id: 'profile', icon: User         },
]

export const WEEK_TAGS = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']
export const CANCEL_HOURS = ['1', '2', '4']