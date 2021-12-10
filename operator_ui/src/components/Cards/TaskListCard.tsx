import React from 'react'

import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Typography from '@material-ui/core/Typography'

import { parseDot, Stratify } from 'utils/parseDot'
import TaskListDag from 'pages/Jobs/TaskListDag'

// TaskListCard renders a card which displays the DAG
export const TaskListCard: React.FC<{ observationSource?: string }> = ({
  observationSource,
}) => {
  const [stratify, setStratify] = React.useState<{
    errorMsg?: string
    stratify?: Stratify[]
  }>()

  React.useEffect(() => {
    if (observationSource && observationSource !== '') {
      try {
        const stratify = parseDot(`digraph {${observationSource}}`)
        setStratify({ stratify })
      } catch (e) {
        setStratify({ errorMsg: 'Failed to parse task graph' })
      }
    } else {
      setStratify({ errorMsg: 'No Task Graph Found' })
    }
  }, [observationSource, setStratify])

  return (
    <Card>
      <CardHeader title="Task List" />
      <CardContent>
        {stratify && stratify.errorMsg && (
          <Typography align="center" variant="subtitle1">
            {stratify.errorMsg}
          </Typography>
        )}

        {stratify && stratify.stratify && (
          <TaskListDag stratify={parseDot(`digraph {${observationSource}}`)} />
        )}
      </CardContent>
    </Card>
  )
}
