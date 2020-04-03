import * as moment from 'moment'

export function humanizeUnixTimestamp(
  timestamp: number,
  format = 'ddd h:mm A',
): string {
  return moment.unix(timestamp).format(format)
}
