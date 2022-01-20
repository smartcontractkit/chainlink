import React from 'react'

import ErrorIcon from './Error'
import PendingIcon from './Pending'
import ListIcon from './ListIcon'
import SuccessIcon from './Success'
import { TaskRunStatus } from 'src/utils/taskRunStatus'

interface Props {
  status: TaskRunStatus
  width?: number
  height?: number
}

export const TaskRunStatusIcon = ({ status, width, height }: Props) => {
  switch (status) {
    case TaskRunStatus.COMPLETE:
      return (
        <SuccessIcon
          width={width}
          height={height}
          data-testid="complete-run-icon"
        />
      )
    case TaskRunStatus.ERROR:
      return (
        <ErrorIcon width={width} height={height} data-testid="error-run-icon" />
      )
    case TaskRunStatus.PENDING:
      return (
        <PendingIcon
          width={width}
          height={height}
          data-testid="pending-run-icon"
        />
      )
    default:
      return (
        <ListIcon
          width={width}
          height={height}
          data-testid="default-run-icon"
        />
      )
  }
}
