import React from 'react'
import TimeAgo from 'react-time-ago/no-tooltip'
import StyledTooltip from 'components/Tooltip'

export default ({children}) => (
  <StyledTooltip title={children}>
    <span>
      <TimeAgo tooltip={false}>{Date.parse(children)}</TimeAgo>
    </span>
  </StyledTooltip>
)
