import React from 'react'
import ReactTimeAgoNoTooltip from 'react-time-ago/no-tooltip'
import moment from 'moment'
import { Tooltip } from './Tooltip'

interface Props {
  children: string
  tooltip: boolean
}

export const localizedTimestamp = (creationDate: string): string =>
  creationDate && moment(creationDate).format()

export const TimeAgo: React.FC<Props> = ({ children, tooltip = false }) => {
  const date = Date.parse(children)
  const ago = <ReactTimeAgoNoTooltip tooltip={false} date={date} />

  if (tooltip) {
    return (
      <Tooltip title={localizedTimestamp(children)}>
        <span>{ago}</span>
      </Tooltip>
    )
  }

  return ago
}
