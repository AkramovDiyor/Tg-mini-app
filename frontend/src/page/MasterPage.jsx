import { useState } from 'react'
import { TodayTab } from './master/TodayTab'
import { QueueTab } from './master/QueueTab'
import { StudioTab } from './master/StudioTab'
import { ProfileTab } from './master/ProfileTab'
import { MasterFloatingNav } from './master/MasterFloatingNav'

export function MasterPage() {
  const [tab, setTab] = useState('today')

  return (
    <div className="animate-fade-up pb-28">
      <div className="px-5 pt-6">
        {tab === 'today'   && <TodayTab />}
        {tab === 'queue'   && <QueueTab />}
        {tab === 'studio'  && <StudioTab />}
        {tab === 'profile' && <ProfileTab />}
      </div>
      <MasterFloatingNav active={tab} onChange={setTab} />
    </div>
  )
}