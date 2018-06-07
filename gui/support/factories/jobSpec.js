import uuid from 'uuid/v4'

export default ({id, initiators, createdAt, tasks}) => {
  return {
    id: id || uuid().replace(/-/g, ''),
    initiators: initiators || [{'type': 'web'}],
    tasks: tasks || [{confirmations: 0, type: 'httpget', url: 'https://bitstamp.net/api/ticker/'}],
    createdAt: createdAt || (new Date()).toISOString()
  }
}
