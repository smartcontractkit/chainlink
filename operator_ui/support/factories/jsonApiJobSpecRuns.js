import uuid from 'uuid/v4'
import { decamelizeKeys } from 'humps'

export default (jobs, jobSpecId, count) => {
  const j = jobs || []
  const jc = count || j.length

  return decamelizeKeys({
    meta: { count: jc },
    data: j.map(c => {
      const config = c || {}
      const id = config.id || uuid().replace(/-/g, '')

      return {
        id: id,
        type: 'runs',
        attributes: {
          id: id,
          jobId: jobSpecId,
          result: {
            jobRunId: id,
            data: {
              value: { result: 'value' }
            },
            status: 'completed',
            error: null
          },
          status: 'completed',
          createdAt: '2018-06-18T15:49:33.015913563-04:00',
          completedAt: '2018-06-18T15:49:33.023078819-04:00'
        }
      }
    })
  })
}
