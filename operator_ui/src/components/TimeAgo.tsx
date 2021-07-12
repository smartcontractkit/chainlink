import React from 'react'
import ReactTimeAgo from 'react-time-ago'
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
  const ago = <ReactTimeAgo tooltip={false} date={date} />

  if (tooltip) {
    return (
      <Tooltip title={localizedTimestamp(children)}>
        <span>{ago}</span>
      </Tooltip>
    )
  }

  return ago
}
