import { augmentOcrTasksList } from './augmentOcrTasksList'
import { PipelineJobRun } from '../sharedTypes'

describe('augmentOcrTasksList', () => {
  it('adds error, output and status attributes', () => {
    const jobRun: PipelineJobRun = {
      pipelineSpec: {
        DotDagSource:
          '    fetch   [type=http method=GET url="https://bitstamp.net/api/ticker/"];\n    parseLast    [type=jsonparse path="last"];\n    multiplyLast [type=multiply times=100];\n\n    fetch2    [type=http method=GET url="https://bitstamp.net/api/ticker/"];\n    parseOpen    [type=jsonparse path="open"];\n    multiplyOpen [type=multiply times=100];\n\n\n fetch -> parseLast  -> multiplyLast -> answer;\n fetch2 -> parseOpen  -> multiplyOpen -> answer;\n\nanswer [type=median                      index=0];\nanswer [type=median                      index=1];\n\n',
      },
      errors: [
        'majority of fetchers in median failed: error making http request: reason; error making http request: reason: bad input for task',
      ],
      outputs: [null],
      createdAt: '2020-11-24T11:38:36.100272Z',
      finishedAt: '2020-11-24T11:39:26.211725Z',
      taskRuns: [
        {
          type: 'median',
          output: null,
          error:
            'majority of fetchers in median failed: error making http request: reason; error making http request: reason: bad input for task',
          taskSpec: {
            dotId: 'answer',
          },
          createdAt: '2020-11-24T11:38:36.100272Z',
          finishedAt: '2020-11-24T11:39:26.19516Z',
          status: 'errored',
        },
        {
          type: 'multiply',
          output: null,
          error: 'error making http request: reason',
          taskSpec: {
            dotId: 'multiplyLast',
          },
          createdAt: '2020-11-24T11:38:36.100272Z',
          finishedAt: '2020-11-24T11:39:26.171678Z',
          status: 'aborted',
        },
        {
          type: 'multiply',
          output: null,
          error: 'error making http request: reason',
          taskSpec: {
            dotId: 'multiplyOpen',
          },
          createdAt: '2020-11-24T11:38:36.100272Z',
          finishedAt: '2020-11-24T11:39:26.176633Z',
          status: 'aborted',
        },
        {
          type: 'jsonparse',
          output: null,
          error: 'error making http request: reason',
          taskSpec: {
            dotId: 'parseLast',
          },
          createdAt: '2020-11-24T11:38:36.100272Z',
          finishedAt: '2020-11-24T11:39:26.154488Z',
          status: 'aborted',
        },
        {
          type: 'jsonparse',
          output: null,
          error: 'error making http request: reason',
          taskSpec: {
            dotId: 'parseOpen',
          },
          createdAt: '2020-11-24T11:38:36.100272Z',
          finishedAt: '2020-11-24T11:39:26.15558Z',
          status: 'aborted',
        },
        {
          type: 'http',
          output: null,
          error: 'error making http request: reason',
          taskSpec: {
            dotId: 'fetch',
          },
          createdAt: '2020-11-24T11:38:36.100272Z',
          finishedAt: '2020-11-24T11:39:26.12949Z',
          status: 'errored',
        },
        {
          type: 'http',
          output: null,
          error: 'error making http request: reason',
          taskSpec: {
            dotId: 'fetch2',
          },
          createdAt: '2020-11-24T11:38:36.100272Z',
          finishedAt: '2020-11-24T11:39:26.127941Z',
          status: 'errored',
        },
      ],
      id: '321',
      jobId: '2',
      status: 'errored',
      type: 'Off-chain reporting job run',
    }
    expect(augmentOcrTasksList({ jobRun })).toEqual([
      {
        attributes: {
          error: 'error making http request: reason',
          method: 'GET',
          output: null,
          status: 'errored',
          type: 'http',
          url: 'https://bitstamp.net/api/ticker/',
        },
        id: 'fetch',
        parentIds: [],
      },
      {
        attributes: {
          error: 'error making http request: reason',
          output: null,
          path: 'last',
          status: 'aborted',
          type: 'jsonparse',
        },
        id: 'parseLast',
        parentIds: ['fetch'],
      },
      {
        attributes: {
          error: 'error making http request: reason',
          output: null,
          status: 'aborted',
          times: '100',
          type: 'multiply',
        },
        id: 'multiplyLast',
        parentIds: ['parseLast'],
      },
      {
        attributes: {
          error: 'error making http request: reason',
          method: 'GET',
          output: null,
          status: 'errored',
          type: 'http',
          url: 'https://bitstamp.net/api/ticker/',
        },
        id: 'fetch2',
        parentIds: [],
      },
      {
        attributes: {
          error: 'error making http request: reason',
          output: null,
          path: 'open',
          status: 'aborted',
          type: 'jsonparse',
        },
        id: 'parseOpen',
        parentIds: ['fetch2'],
      },
      {
        attributes: {
          error: 'error making http request: reason',
          output: null,
          status: 'aborted',
          times: '100',
          type: 'multiply',
        },
        id: 'multiplyOpen',
        parentIds: ['parseOpen'],
      },
      {
        attributes: {
          error:
            'majority of fetchers in median failed: error making http request: reason; error making http request: reason: bad input for task',
          index: '1',
          output: null,
          status: 'errored',
          type: 'median',
        },
        id: 'answer',
        parentIds: ['multiplyLast', 'multiplyOpen'],
      },
    ])
  })
})
