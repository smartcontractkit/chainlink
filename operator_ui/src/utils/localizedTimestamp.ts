import moment from 'moment'

export default (creationDate: Date) => creationDate && moment(creationDate).format()
