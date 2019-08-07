import moment from 'moment'

export const localizedTimestamp = (creationDate: string): string =>
  creationDate && moment(creationDate).format()
