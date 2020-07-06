import uuid from 'uuid/v4'
import isoDate from 'test-helpers/isoDate'

export default (jobs, count) => {
  const j = jobs || []
  const jc = count || j.length

  return {
    meta: { count: jc },
    data: j.map(c => {
      const config = c || {}
      const id = config.id || uuid().replace(/-/g, '')
      const initiators = config.initiators || [{ type: 'web' }]
      const earnings = config.earnings
      const minPay = config.minPayment
      const tasks = config.tasks || [
        {
          confirmations: 0,
          type: 'httpget',
          url: 'https://bitstamp.net/api/ticker/',
        },
      ]
      const createdAt = config.createdAt || new Date().toISOString()
      const errors = c.errors || []
      let runs = c.runs || []

      runs = runs.map(r =>
        Object.assign(
          {},
          { createdAt: isoDate(Date.now()) },
          { result: {} },
          { jobId: id },
          r,
        ),
      )

      return {
        type: 'specs',
        id,
        attributes: {
          initiators,
          id,
          tasks,
          minPayment: minPay,
          createdAt,
          earnings,
          runs,
          errors,
        },
      }
    }),
  }
}
