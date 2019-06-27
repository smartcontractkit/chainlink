import moment from 'moment'

export default creationDate => creationDate && moment(creationDate).format()
