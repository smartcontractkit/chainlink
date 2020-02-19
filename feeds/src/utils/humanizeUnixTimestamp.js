import moment from 'moment'

export const humanizeUnixTimestamp = (timestamp, format = 'ddd h:mm A') =>
  timestamp && moment.unix(timestamp).format(format)
