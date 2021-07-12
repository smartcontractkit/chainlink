import { partialAsFull } from 'support/test-helpers/partialAsFull'
import { OcrJobRun } from 'core/store/models'

export function jobRunV2(
  config: Partial<OcrJobRun & { id?: string; dotDagSource?: string }> = {},
): OcrJobRun {
  return partialAsFull<OcrJobRun>({
    outputs: config.outputs || [null],
    errors: config.errors || [],
    pipelineSpec: {
      ID: 1,
      dotDagSource:
        config?.dotDagSource ||
        '   fetch    [type=http method=POST url="http://localhost:8001" requestData="{\\"hi\\": \\"hello\\"}"];\n    parse    [type=jsonparse path="data,result"];\n    multiply [type=multiply times=100];\n    fetch -\u003e parse -\u003e multiply;\n',
      CreatedAt: '2020-11-19T14:01:24.989522Z',
    },
    taskRuns: config.taskRuns || [
      {
        createdAt: '2020-11-19T14:01:24.989522Z',
        error: `error making http request: Post "http://localhost:8001": dial tcp 127.0.0.1:8001: connect: connection refused`,
        finishedAt: '2020-11-19T14:01:25.015681Z',
        output: null,
        dotId: 'multiply',
        type: 'multiply',
      },

      {
        createdAt: '2020-11-19T14:01:24.989522Z',
        error: `error making http request: Post "http://localhost:8001": dial tcp 127.0.0.1:8001: connect: connection refused`,
        finishedAt: '2020-11-19T14:01:25.005568Z',
        output: null,
        dotId: 'parse',
        type: 'jsonparse',
      },
      {
        createdAt: '2020-11-19T14:01:24.989522Z',
        error: `error making http request: Post "http://localhost:8001": dial tcp 127.0.0.1:8001: connect: connection refused`,
        finishedAt: '2020-11-19T14:01:24.997068Z',
        output: null,
        dotId: 'fetch',
        type: 'http',
      },
    ],
    createdAt: config.createdAt || new Date(1600775300410).toISOString(),
    finishedAt: config.finishedAt || new Date(1600775300410).toISOString(),
  })
}
