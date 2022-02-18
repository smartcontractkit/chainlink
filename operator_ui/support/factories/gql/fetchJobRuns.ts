import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'

export function buildRun(
  overrides?: Partial<JobRunsPayload_ResultsFields>,
): JobRunsPayload_ResultsFields {
  const minuteAgo = isoDate(Date.now() - MINUTE_MS)

  return {
    __typename: 'JobRun',
    id: '1',
    allErrors: [],
    createdAt: minuteAgo,
    finishedAt: minuteAgo,
    status: 'COMPLETED',
    job: {
      id: '10',
    },
    ...overrides,
  }
}

export function buildRuns(): ReadonlyArray<JobRunsPayload_ResultsFields> {
  const twoMinutesAgo = isoDate(Date.now() - 2 * MINUTE_MS)

  return [
    buildRun(),
    buildRun({
      id: '2',
      allErrors: ['error'],
      createdAt: twoMinutesAgo,
      finishedAt: twoMinutesAgo,
      status: 'ERRORED',
      job: {
        id: '20',
      },
    }),
  ]
}
