import { useState } from 'react'
import { TodayTab } from './master/TodayTab'
import { QueueTab } from './master/QueueTab'
import { ProfileTab } from './master/ProfileTab'
import { MasterFloatingNav } from './master/MasterFloatingNav'
// import { ErrorFallback } from './master/ErrorFallback'

export function MasterPage() {
  const [tab, setTab] = useState('today')
  const [errorKey, setErrorKey] = useState(0)

  // const handleRetry = () => setErrorKey((k) => k + 1)

  return (
    <div className="animate-fade-up pb-28">
      <div className="px-5 pt-6">
        {tab === 'today'   && <TodayTab   key={`today-${errorKey}`}   />}
        {tab === 'queue'   && <QueueTab   key={`queue-${errorKey}`}   />}
        {tab === 'profile' && <ProfileTab key={`profile-${errorKey}`} />}
      </div>
      <MasterFloatingNav active={tab} onChange={setTab} />
    </div>
  )
}