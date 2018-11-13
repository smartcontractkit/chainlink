import React from 'react'
import TimeAgo from 'react-time-ago/no-tooltip'
import Tooltip from '@material-ui/core/Tooltip'

export default ({children}) => (
  <Tooltip title={children}>
    <span>
      <TimeAgo>{Date.parse(children)}</TimeAgo>
    </span>
  </Tooltip>
)
