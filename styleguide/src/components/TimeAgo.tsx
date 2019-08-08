import React from 'react'
import TimeAgoNoTooltip from 'react-time-ago/no-tooltip'
import { localizedTimestamp } from '../utils/localizedTimestamp'
import { Tooltip } from './Tooltip'

interface Props {
  children: string
  tooltip: boolean
}

export const TimeAgo: React.FC<Props> = ({ children, tooltip = false }) => {
  const date = Date.parse(children)
  const ago = <TimeAgoNoTooltip tooltip={false}>{date}</TimeAgoNoTooltip>

  if (tooltip) {
    return (
      <Tooltip title={localizedTimestamp(children)}>
        <span>{ago}</span>
      </Tooltip>
    )
  }

  return ago
}
