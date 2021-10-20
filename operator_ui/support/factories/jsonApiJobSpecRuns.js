import uuid from 'uuid/v4'

export default (runs, count) => {
  const r = runs || []
  const rc = count || r.length

  return {
    meta: { count: rc },
    data: r.map((c) => {
      const config = c || {}
      const id = config.id || uuid().replace(/-/g, '')
      const jobId = config.jobId || uuid().replace(/-/g, '')
      const status = config.status || 'completed'

      return {
        id,
        type: 'runs',
        attributes: {
          id,
          jobId,
          result: {
            jobRunId: id,
            data: {
              value: { result: 'value' },
            },
            status,
            error: null,
          },
          status,
          createdAt: '2018-06-18T15:49:33.015913563-04:00',
          finishedAt: '2018-06-18T15:49:33.023078819-04:00',
        },
      }
    }),
  }
}
