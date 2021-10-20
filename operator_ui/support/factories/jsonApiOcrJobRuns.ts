import { ApiResponse } from 'utils/json-api-client'
import { OcrJobRun } from 'core/store/models'
import { jobRunV2 } from './jobRunV2'

function getRandomInt(max: number) {
  return Math.floor(Math.random() * Math.floor(max))
}

export const jsonApiOcrJobRuns = (
  jobs: Partial<OcrJobRun & { id?: string }>[] = [],
  count?: number,
) => {
  const rc = count || jobs.length

  return {
    data: jobs.map((config) => {
      const id = config.id || getRandomInt(1_000_000).toString()

      return {
        type: 'runs',
        id,
        attributes: jobRunV2(config),
      }
    }),
    meta: { count: rc },
  } as ApiResponse<OcrJobRun[]>
}
