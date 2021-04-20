import { ApiResponse } from 'utils/json-api-client'
import { JobRun } from 'core/store/models'
import jsonApiJobSpecRun from 'factories/jsonApiJobSpecRun'
import { jsonApiOcrJobRun } from 'factories/jsonApiOcrJobRun'
import {
  transformDirectRequestJobRun,
  transformPipelineJobRun,
} from './transformJobRuns'

describe('transformPipelineJobRun', () => {
  it('transforms api response to PipelineJobRun', () => {
    const apiResponse = jsonApiOcrJobRun({
      id: '1',
    })

    expect(transformPipelineJobRun('1')(apiResponse.data)).toEqual({
      createdAt: '2020-09-22T11:48:20.410Z',
      errors: [],
      finishedAt: '2020-09-22T11:48:20.410Z',
      id: '1',
      jobId: '1',
      outputs: [null],
      pipelineSpec: {
        CreatedAt: '2020-11-19T14:01:24.989522Z',
        dotDagSource: `   fetch    [type=http method=POST url="http://localhost:8001" requestData="{\\"hi\\": \\"hello\\"}"];
    parse    [type=jsonparse path="data,result"];
    multiply [type=multiply times=100];
    fetch -> parse -> multiply;
`,
        ID: 1,
      },
      status: 'errored',
      taskRuns: [
        {
          createdAt: '2020-11-19T14:01:24.989522Z',
          error:
            'error making http request: Post "http://localhost:8001": dial tcp 127.0.0.1:8001: connect: connection refused',
          finishedAt: '2020-11-19T14:01:25.015681Z',
          output: null,
          status: 'not_run',
          dotId: 'multiply',
          type: 'multiply',
        },
        {
          createdAt: '2020-11-19T14:01:24.989522Z',
          error:
            'error making http request: Post "http://localhost:8001": dial tcp 127.0.0.1:8001: connect: connection refused',
          finishedAt: '2020-11-19T14:01:25.005568Z',
          output: null,
          status: 'not_run',
          dotId: 'parse',
          type: 'jsonparse',
        },
        {
          createdAt: '2020-11-19T14:01:24.989522Z',
          error:
            'error making http request: Post "http://localhost:8001": dial tcp 127.0.0.1:8001: connect: connection refused',
          finishedAt: '2020-11-19T14:01:24.997068Z',
          output: null,
          status: 'errored',
          dotId: 'fetch',
          type: 'http',
        },
      ],
      type: 'Pipeline job run',
    })
  })
})

describe('transformDirectRequestJobRun', () => {
  it('transforms api response to DirectRequestJobRun', () => {
    const apiResponse = jsonApiJobSpecRun({
      id: '1',
    }) as ApiResponse<JobRun>

    expect(transformDirectRequestJobRun('1')(apiResponse.data)).toEqual({
      createdAt: '2018-06-19T15:39:53.315919143-07:00',
      id: '1',
      initiator: { params: {}, type: 'web' },
      jobId: '1',
      result: {
        data: {
          value:
            '0x05070f7f6a40e4ce43be01fa607577432c68730c2cb89a0f50b665e980d926b5',
        },
      },
      status: 'completed',
      taskRuns: [],
      type: 'Direct request job run',
    })
  })
})
