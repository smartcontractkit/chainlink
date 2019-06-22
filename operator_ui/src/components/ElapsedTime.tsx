import React from 'react'
import { Typography } from '@material-ui/core'
import elapsedTimeHHMMSS from '../utils/elapsedTimeHHMMSS'

const ElapsedTime = ({ start, end, className }) => {
  return (
    <Typography variant="h6" className={className}>
      {elapsedTimeHHMMSS(start, end)}
    </Typography>
  )
}

export default ElapsedTime
