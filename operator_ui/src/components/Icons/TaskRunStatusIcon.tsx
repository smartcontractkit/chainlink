import React from 'react'

import ErrorIcon from './Error'
import SuccessIcon from './Success'

interface Props {
  status: 'completed' | 'errored'
  width?: number
  height?: number
}

export const TaskRunStatusIcon = ({ status, width, height }: Props) => {
  switch (status) {
    case 'completed':
      return (
        <SuccessIcon width={width} height={height} data-testid="completed" />
      )
    case 'errored':
      return <ErrorIcon width={width} height={height} data-testid="errored" />
    default:
      return null
  }
}
