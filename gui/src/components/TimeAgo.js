import React from 'react'
import TimeAgo from 'react-time-ago/no-tooltip'
import Tooltip from '@material-ui/core/Tooltip'

export default ({children}) => (
  <Tooltip title={new Date(children).toISOString()}>
    <span>
      <TimeAgo>{children}</TimeAgo>
    </span>
  </Tooltip>
)
