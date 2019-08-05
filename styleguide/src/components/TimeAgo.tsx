import React from 'react'
import TimeAgoNoTooltip from 'react-time-ago/no-tooltip'
import { localizedTimestamp } from '../utils/localizedTimestamp'
import { Tooltip } from './Tooltip'

interface IProps {
  children: string
  tooltip: boolean
}

export const TimeAgo = ({ children, tooltip = false }: IProps) => {
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
