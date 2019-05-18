import React from 'react'
import SuccessIcon from './Success'
import ErrorIcon from './Error'
import PendingIcon from './Pending'

const StatusIcon = ({ className, children, width, height }) => {
  if (children === 'completed') {
    return <SuccessIcon className={className} width={width} height={height} />
  } else if (children === 'errored') {
    return <ErrorIcon className={className} width={width} height={height} />
  }

  return <PendingIcon className={className} width={width} height={height} />
}

export default StatusIcon
