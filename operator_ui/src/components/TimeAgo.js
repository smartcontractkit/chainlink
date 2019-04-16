import React from 'react'
import TimeAgoNoTooltip from 'react-time-ago/no-tooltip'
import StyledTooltip from 'components/Tooltip'

const TimeAgo = ({ children }) => (
  <StyledTooltip title={children}>
    <span>
      <TimeAgoNoTooltip tooltip={false}>
        {Date.parse(children)}
      </TimeAgoNoTooltip>
    </span>
  </StyledTooltip>
)

export default TimeAgo
