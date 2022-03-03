import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'

export function buildRun(
  overrides?: Partial<JobRunPayload_Fields>,
): JobRunPayload_Fields {
  const minuteAgo = isoDate(Date.now() - MINUTE_MS)
  const twoMinutesAgo = isoDate(Date.now() - MINUTE_MS * 2)

  return {
    __typename: 'JobRun',
    id: '1',
    allErrors: [],
    createdAt: twoMinutesAgo,
    fatalErrors: [],
    finishedAt: minuteAgo,
    inputs: '',
    job: {
      id: '10',
      name: 'job 1',
      observationSource: '',
    },
    outputs: [],
    status: 'COMPLETED',
    taskRuns: [],
    ...overrides,
  }
}

export function buildTaskRun(
  overrides?: Partial<JobRunPayload_TaskRunsFields>,
): JobRunPayload_TaskRunsFields {
  const minuteAgo = isoDate(Date.now() - MINUTE_MS)
  const twoMinutesAgo = isoDate(Date.now() - MINUTE_MS * 2)

  return {
    __typename: 'TaskRun',
    id: '00000000-0000-0000-0000-000000000001',
    createdAt: twoMinutesAgo,
    dotID: 'parse_request',
    error: 'data: parameter is empty',
    finishedAt: minuteAgo,
    output: 'null',
    type: 'jsonparse',
    ...overrides,
  }
}
