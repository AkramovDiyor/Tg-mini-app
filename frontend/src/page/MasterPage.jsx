import { useState } from 'react'
import { TodayTab } from './master/TodayTab'
import { QueueTab } from './master/QueueTab'
import { StudioTab } from './master/StudioTab'
import { ProfileTab } from './master/ProfileTab'
import { MasterFloatingNav } from './master/MasterFloatingNav'

export function MasterPage() {
  const [tab, setTab] = useState('today')

  return (
    <div className="pb-28">
      <div className="px-5 pt-6">
        <div className="animate-fade-up">
          {tab === 'today' && <TodayTab />}
          {tab === 'queue' && <QueueTab />}
          {tab === 'studio' && <StudioTab />}
          {tab === 'profile' && <ProfileTab />}
        </div>
      </div>
      <MasterFloatingNav active={tab} onChange={setTab} />
    </div>
  )
}
