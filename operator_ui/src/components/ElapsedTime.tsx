import React from 'react'
import { Typography } from '@material-ui/core'
import elapsedTimeHHMMSS from '../utils/elapsedTimeHHMMSS'

interface IProps {
  start: string
  end: string
  className: string
}

const ElapsedTime =
({ start, end, className }: IProps) => {
  return (
    <Typography id='elapsedTime' variant="h6" className={className}>
      {elapsedTimeHHMMSS(start, end)}
    </Typography>
  )
}

export default ElapsedTime
