import { ApiResponse } from '@chainlink/json-api-client'
import { OcrJobSpec } from 'core/store/models'
import { jobSpecV2 } from './jobSpecV2'

function getRandomInt(max: number) {
  return Math.floor(Math.random() * Math.floor(max))
}

export const jsonApiOcrJobSpecs = (
  jobs: Partial<
    OcrJobSpec['offChainReportingOracleSpec'] & { id?: string }
  >[] = [],
  count?: number,
) => {
  const rc = count || jobs.length

  return {
    data: jobs.map((config) => {
      const id = config.id || getRandomInt(1_000_000).toString()

      return {
        type: 'jobSpecV2s',
        id,
        attributes: jobSpecV2(config),
      }
    }),
    meta: { count: rc },
  } as ApiResponse<OcrJobSpec[]>
}
