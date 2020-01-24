import * as moment from 'moment'

export function humanizeUnixTimestamp(
  timestamp: number,
  format = 'ddd h:mm A',
) {
  return timestamp && moment.unix(timestamp).format(format)
}
