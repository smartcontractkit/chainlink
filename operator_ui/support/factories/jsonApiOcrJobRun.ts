import { ApiResponse } from 'utils/json-api-client'
import { JobRunV2 } from 'core/store/models'
import { jobRunV2 } from './jobRunV2'

function getRandomInt(max: number) {
  return Math.floor(Math.random() * Math.floor(max))
}

export const jobRunAPIResponse = (
  config: Partial<JobRunV2 & { id?: string; dotDagSource?: string }> = {},
) => {
  const id = config.id || getRandomInt(1_000_000).toString()

  return {
    data: {
      type: 'run',
      id,
      attributes: jobRunV2(config),
    },
  } as ApiResponse<JobRunV2>
}

export const jobRunsAPIResponse = (
  configs: Partial<JobRunV2 & { id?: string }>[] = [],
  count?: number,
) => {
  const rc = count || configs.length

  return {
    data: configs.map((config) => {
      const id = config.id || getRandomInt(1_000_000).toString()

      return {
        type: 'runs',
        id,
        attributes: jobRunV2(config),
      }
    }),
    meta: { count: rc },
  } as ApiResponse<JobRunV2[]>
}
