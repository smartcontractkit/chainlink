import { ApiResponse } from 'utils/json-api-client'
import { Job } from 'core/store/models'
import { ocrJob } from './jobV2'

function getRandomInt(max: number) {
  return Math.floor(Math.random() * Math.floor(max))
}

export const jsonApiOcrJobSpec = (
  config: Partial<
    Job['offChainReportingOracleSpec'] & { id?: string } & {
      dotDagSource?: string
    }
  > = {},
) => {
  const id = config.id || getRandomInt(1_000_000).toString()

  return {
    data: {
      type: 'jobSpecV2',
      id,
      attributes: ocrJob(config),
    },
  } as ApiResponse<Job>
}
