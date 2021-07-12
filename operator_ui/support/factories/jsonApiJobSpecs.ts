import { v4 as uuid } from 'uuid'
import { PaginatedApiResponse } from 'utils/json-api-client'
import { partialAsFull } from 'support/test-helpers/partialAsFull'
import { JobSpec, TaskSpec, InitiatorType } from 'core/store/models'

export const jsonApiJobSpecs = (
  jobs: Partial<JobSpec>[] = [],
  count?: number,
): PaginatedApiResponse<JobSpec[]> => {
  const jobsCount = count || jobs.length

  return {
    meta: { count: jobsCount },
    links: {},
    data: jobs.map((config, index) => {
      const id = config.id || uuid().replace(/-/g, '')
      const initiators = config.initiators || [
        {
          id: 1,
          type: 'web' as InitiatorType.WEB,
          jobSpecId: id,
          CreatedAt: new Date(1600775300410).toISOString(),
        },
      ]
      const earnings = config.earnings
      const minPay = config.minPayment
      const name = config.name || `Job ${index + 1}`
      const tasks = config.tasks || [
        partialAsFull<TaskSpec>({
          confirmations: 0,
          type: 'httpget',
          params: {
            get: 'https://bitstamp.net/api/ticker/',
          },
        }),
      ]
      const createdAt =
        config.createdAt || new Date(1600775300410).toISOString()
      const startAt = config.startAt || new Date(1600775390410).toISOString()
      const endAt = config.endAt || new Date(1600775990410).toISOString()
      const errors = config.errors || []

      const attributes = partialAsFull<JobSpec>({
        createdAt,
        earnings,
        endAt,
        errors,
        id,
        initiators,
        minPayment: minPay,
        name,
        startAt,
        tasks,
      })

      return {
        type: 'specs',
        id,
        attributes,
        relationships: {} as never,
        links: {} as never,
        meta: {} as never,
      }
    }),
  }
}

export default jsonApiJobSpecs
