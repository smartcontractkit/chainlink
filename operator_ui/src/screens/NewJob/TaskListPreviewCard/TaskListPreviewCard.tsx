import React from 'react'

import TOML from '@iarna/toml'
import { TaskListCard } from 'src/components/Cards/TaskListCard'

interface Props {
  toml: string
}

export const TaskListPreviewCard: React.FC<Props> = ({ toml }) => {
  const [observationSource, setObservationSource] = React.useState('')

  React.useEffect(() => {
    try {
      const spec = TOML.parse(toml)

      if (spec.observationSource) {
        setObservationSource((spec.observationSource as string).trim())
      } else {
        setObservationSource('')
      }
    } catch (e) {
      setObservationSource('')
    }
  }, [toml])

  return <TaskListCard observationSource={observationSource} />
}
