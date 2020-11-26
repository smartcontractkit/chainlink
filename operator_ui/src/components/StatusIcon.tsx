import React from 'react'
import SuccessIcon from 'components/Icons/Success'
import ErrorIcon from 'components/Icons/Error'
import PendingIcon from 'components/Icons/Pending'

interface Props {
  children: React.ReactNode
  width?: number
  height?: number
}

const StatusIcon = ({ children, width, height }: Props) => {
  if (children === 'completed') {
    return <SuccessIcon width={width} height={height} />
  } else if (children === 'errored') {
    return <ErrorIcon width={width} height={height} />
  }

  return <PendingIcon width={width} height={height} />
}

export default StatusIcon
