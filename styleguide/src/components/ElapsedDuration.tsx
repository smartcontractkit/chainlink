import React from 'react'
import { Typography } from '@material-ui/core'
import { elapsedDuration, FinishedAt } from '../utils/elapsedDuration'

interface Props {
  start: string
  end: FinishedAt
  className?: string
}

export const ElapsedDuration: React.FC<Props> = ({ start, end, className }) => {
  return (
    <Typography variant="h6" className={className}>
      {elapsedDuration(start, end)}
    </Typography>
  )
}
