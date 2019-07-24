import moment from 'moment'

export default (creationDate: string): string => creationDate && moment(creationDate).format()
