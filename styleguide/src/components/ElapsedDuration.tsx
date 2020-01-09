import React from 'react'
import { Typography } from '@material-ui/core'
import { elapsedDuration } from '../utils/elapsedDuration'

interface Props {
  start: string
  end: string
  className?: string
}

export const ElapsedDuration: React.FC<Props> = ({ start, end, className }) => {
  return (
    <Typography role="elapsedduration" variant="h6" className={className}>
      {elapsedDuration(start, end)}
    </Typography>
  )
}
