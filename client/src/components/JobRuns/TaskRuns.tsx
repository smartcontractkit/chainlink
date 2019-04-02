import React from 'react'
import Typography from '@material-ui/core/Typography'

interface IProps {
  taskRuns?: ITaskRun[]
}

const TaskRuns = ({ taskRuns }: IProps) => {
  return (
    <ul>
      {taskRuns &&
        taskRuns.map((run: ITaskRun) => {
          return (
            <li key={run.id}>
              <Typography variant="body1">{run.type}</Typography>
            </li>
          )
        })}
    </ul>
  )
}

export default TaskRuns
