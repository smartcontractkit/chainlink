import React from 'react'
import SuccessIcon from './Success'
import ErrorIcon from './Error'
import PendingIcon from './Pending'

interface IProps {
  children: string
  width?: number
  height?: number
}

const StatusIcon = ({ children, width, height }: IProps) => {
  if (children === 'completed') {
    return <SuccessIcon width={width} height={height} />
  } else if (children === 'errored') {
    return <ErrorIcon width={width} height={height} />
  }

  return <PendingIcon width={width} height={height} />
}

export default StatusIcon
