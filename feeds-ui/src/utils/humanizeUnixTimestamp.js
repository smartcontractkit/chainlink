import moment from 'moment'

export const humanizeUnixTimestamp = timestamp =>
  timestamp && moment.unix(timestamp).format('ddd h:mm A')
