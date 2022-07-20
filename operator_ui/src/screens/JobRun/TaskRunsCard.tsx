import React from 'react'

import Card from '@material-ui/core/Card'
import { parseDot } from 'src/utils/parseDot'
import { TaskRunItem } from './TaskRunItem'

interface Props {
  taskRuns: ReadonlyArray<JobRunPayload_TaskRunsFields>
  observationSource?: string
}

export const TaskRunsCard = ({ taskRuns, observationSource }: Props) => {
  const items = React.useMemo(() => {
    const graph = parseDot(`digraph {${observationSource}}`)

    return graph.map((node) => {
      const taskRun = taskRuns.find(({ dotID }) => {
        return dotID === node.id
      })

      if (!taskRun) {
        return undefined
      }

      return {
        ...taskRun,
        attrs: node.attributes,
      }
    })
  }, [observationSource, taskRuns])

  return (
    <Card>
      {items.map((item, idx) => {
        return item ? <TaskRunItem {...item} key={idx} /> : null
      })}
    </Card>
  )
}
