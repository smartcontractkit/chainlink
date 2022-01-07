import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'

export function buildRun(
  overrides?: Partial<RecentJobRunsPayload_ResultsFields>,
): RecentJobRunsPayload_ResultsFields {
  const minuteAgo = isoDate(Date.now() - MINUTE_MS)

  return {
    __typename: 'JobRun',
    id: '1',
    allErrors: [],
    createdAt: minuteAgo,
    finishedAt: minuteAgo,
    job: {
      id: '10',
    },
    status: 'COMPLETED',
    ...overrides,
  }
}

export function buildRuns(): ReadonlyArray<RecentJobRunsPayload_ResultsFields> {
  const twoMinutesAgo = isoDate(Date.now() - 2 * MINUTE_MS)

  return [
    buildRun(),
    buildRun({
      id: '2',
      createdAt: twoMinutesAgo,
      finishedAt: twoMinutesAgo,
      status: 'COMPLETED',
      job: {
        id: '20',
      },
    }),
  ]
}
