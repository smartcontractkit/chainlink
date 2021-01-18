import React from 'react'
import { Typography } from '@material-ui/core'

const MS_IN_SECOND = 1000
const SECONDS_IN_MINUTE = 60
const SECONDS_IN_HOUR = 3600

/**
 * A string, or a number that can be converted to a Date using `new Date(...)` e.g.
 *
 * string: 2020-01-03T22:45:00.166261Z
 * number: 1578091500166
 *
 * @typedef {(number|string)} DateValue
 */
export type DateValue = string | number

export function elapsedDuration(
  createdAt: DateValue,
  finishedAt: DateValue,
): string {
  if (createdAt === '' && finishedAt === '') {
    return ''
  }

  const es = elapsedSeconds(new Date(createdAt), new Date(finishedAt))
  const hours = Math.floor(es / SECONDS_IN_HOUR)
  const minutes = Math.floor((es % SECONDS_IN_HOUR) / SECONDS_IN_MINUTE)
  const seconds = Math.ceil((es % SECONDS_IN_HOUR) % SECONDS_IN_MINUTE)

  return format(hours, minutes, seconds)
}

function elapsedSeconds(from: Date, to: Date): number {
  return to.getTime() / MS_IN_SECOND - from.getTime() / MS_IN_SECOND
}

function format(hours: number, minutes: number, seconds: number): string {
  return `${hours ? `${hours}h` : ''}${minutes ? `${minutes}m` : ''}${seconds}s`
}

interface Props {
  start: string
  end: DateValue
  className?: string
}

export const ElapsedDuration: React.FC<Props> = ({ start, end, className }) => {
  return (
    <Typography variant="h6" className={className}>
      {elapsedDuration(start, end)}
    </Typography>
  )
}
