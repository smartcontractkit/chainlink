import React from 'react'
import TimeAgoNoTooltip from 'react-time-ago/no-tooltip'
import Tooltip from './Tooltip'

interface IProps {
  children: string
  tooltip?: boolean
}

const TimeAgo = ({ children, tooltip }: IProps) => {
  const date = Date.parse(children)
  const ago = <TimeAgoNoTooltip date={date} tooltip={false} />

  if (tooltip) {
    return (
      <Tooltip title={children}>
        <span>{ago}</span>
      </Tooltip>
    )
  }

  return ago
}

export default TimeAgo
