import { ApiResponse } from 'utils/json-api-client'
import { JobSpecV2 } from 'core/store/models'
import { ocrJobSpecV2 } from './jobSpecV2'

function getRandomInt(max: number) {
  return Math.floor(Math.random() * Math.floor(max))
}

export const jsonApiOcrJobSpec = (
  config: Partial<
    JobSpecV2['offChainReportingOracleSpec'] & { id?: string } & {
      dotDagSource?: string
    }
  > = {},
) => {
  const id = config.id || getRandomInt(1_000_000).toString()

  return {
    data: {
      type: 'jobSpecV2',
      id,
      attributes: ocrJobSpecV2(config),
    },
  } as ApiResponse<JobSpecV2>
}
