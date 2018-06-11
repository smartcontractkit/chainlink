import uuid from 'uuid/v4'

const defaults = () => (
  {
    id: uuid().replace(/-/g, ''),
    initiators: [{'type': 'web'}],
    tasks: [{confirmations: 0, type: 'httpget', url: 'https://bitstamp.net/api/ticker/'}],
    createdAt: (new Date()).toISOString()
  }
)

export default (attrs) => Object.assign(defaults(), attrs)
