import uuid from 'uuid/v4'

export default (o) => {
  const opts = o || {}
  const id = opts.id || uuid().replace(/-/g, '')
  const initiators = opts.initiators || [{'type': 'web'}]
  const tasks = opts.tasks || [{confirmations: 0, type: 'httpget', url: 'https://bitstamp.net/api/ticker/'}]
  const createdAt = opts.createdAt || (new Date()).toISOString()

  return {
    data: [
      {
        type: 'specs',
        id: id,
        attributes: {
          initiators: initiators,
          id: id,
          tasks: tasks,
          createdAt: createdAt
        }
      }
    ]
  }
}
