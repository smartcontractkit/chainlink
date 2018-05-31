import uuid from 'uuid/v4'
import { decamelizeKeys } from 'humps'

export default (jobs) => {
  const j = jobs || []

  return decamelizeKeys({
    meta: { count: j.length },
    data: j.map((c) => {
      const config = c || {}
      const id = config.id || uuid().replace(/-/g, '')
      const initiators = config.initiators || [{'type': 'web'}]
      const tasks = config.tasks || [{confirmations: 0, type: 'httpget', url: 'https://bitstamp.net/api/ticker/'}]
      const createdAt = config.createdAt || (new Date()).toISOString()

      return {
        type: 'specs',
        id: id,
        attributes: {
          initiators: initiators,
          id: id,
          tasks: tasks,
          createdAt: createdAt
        }
      }
    })
  })
}
