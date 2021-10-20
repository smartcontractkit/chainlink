import { ApiResponse } from 'utils/json-api-client'
import { OcrJobRun } from 'core/store/models'
import { jobRunV2 } from './jobRunV2'

function getRandomInt(max: number) {
  return Math.floor(Math.random() * Math.floor(max))
}

export const jsonApiOcrJobRun = (
  config: Partial<OcrJobRun & { id?: string; dotDagSource?: string }> = {},
) => {
  const id = config.id || getRandomInt(1_000_000).toString()

  return {
    data: {
      type: 'run',
      id,
      attributes: jobRunV2(config),
    },
  } as ApiResponse<OcrJobRun>
}
