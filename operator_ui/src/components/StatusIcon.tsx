import React from 'react'

import SuccessIcon from 'components/Icons/Success'
import ErrorIcon from 'components/Icons/Error'
import PendingIcon from 'components/Icons/Pending'
import ListIcon from 'components/Icons/ListIcon'

interface Props {
  children: React.ReactNode
  width?: number
  height?: number
}

const StatusIcon = ({ children, width, height }: Props) => {
  if (children === 'completed') {
    return <SuccessIcon width={width} height={height} data-testid="completed" />
  } else if (children === 'errored') {
    return <ErrorIcon width={width} height={height} data-testid="errored" />
  } else if (children === 'not_run') {
    return <ListIcon width={width} height={height} data-testid="not_run" />
  }

  return <PendingIcon width={width} height={height} data-testid="pending" />
}

export default StatusIcon
