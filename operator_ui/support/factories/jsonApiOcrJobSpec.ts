import { ApiResponse } from '@chainlink/json-api-client'
import { OcrJobSpec } from 'core/store/models'
import { jobSpecV2 } from './jobSpecV2'

function getRandomInt(max: number) {
  return Math.floor(Math.random() * Math.floor(max))
}

export const jsonApiOcrJobSpec = (
  config: Partial<
    OcrJobSpec['offChainReportingOracleSpec'] & { id?: string } & {
      dotDagSource?: string
    }
  > = {},
) => {
  const id = config.id || getRandomInt(1_000_000).toString()

  return {
    data: {
      type: 'jobSpecV2',
      id,
      attributes: jobSpecV2(config),
    },
  } as ApiResponse<OcrJobSpec>
}
