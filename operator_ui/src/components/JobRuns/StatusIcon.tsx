import React from 'react'
import SuccessIcon from '../Icons/Success'
import ErrorIcon from '../Icons/Error'
import PendingIcon from '../Icons/Pending'

interface IProps {
  children: React.ReactNode
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
