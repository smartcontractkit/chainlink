import { jobRunAPIResponse } from 'factories/jsonApiOcrJobRun'
import { transformPipelineJobRun } from './transformJobRuns'

describe('transformPipelineJobRun', () => {
  it('transforms api response to PipelineJobRun', () => {
    const apiResponse = jobRunAPIResponse({
      id: '1',
      errors: ['task inputs: too many errors'],
    })

    expect(transformPipelineJobRun('1')(apiResponse.data)).toEqual({
      createdAt: '2020-09-22T11:48:20.410Z',
      errors: ['task inputs: too many errors'],
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
        jobID: '1',
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
