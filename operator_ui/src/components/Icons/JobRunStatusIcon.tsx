import React from 'react'

import ErrorIcon from './Error'
import { JobRunStatus } from 'src/types/generated/graphql'
import ListIcon from 'components/Icons/ListIcon'
import PendingIcon from './Pending'
import SuccessIcon from './Success'

interface Props {
  status: JobRunStatus
  width?: number
  height?: number
}

export const JobRunStatusIcon = ({ status, width, height }: Props) => {
  switch (status) {
    case 'COMPLETED':
      return (
        <SuccessIcon width={width} height={height} data-testid="completed" />
      )
    case 'ERRORED':
      return <ErrorIcon width={width} height={height} data-testid="errored" />
    case 'RUNNING':
      return <PendingIcon width={width} height={height} data-testid="running" />
    case 'SUSPENDED':
      return <ListIcon width={width} height={height} data-testid="suspended" />
    default:
      return null
  }
}
