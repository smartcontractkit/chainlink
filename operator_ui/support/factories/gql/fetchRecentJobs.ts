import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'

// buildRecentJob builds a job for the FetchRecentJobs query.
export function buildRecentJob(
  overrides?: Partial<RecentJobsPayload_ResultsFields>,
): RecentJobsPayload_ResultsFields {
  const minuteAgo = isoDate(Date.now() - MINUTE_MS)

  return {
    __typename: 'Job',
    id: '1',
    name: 'job 1',
    createdAt: minuteAgo,
    ...overrides,
  }
}

// buildRecentJobs builds a list of recent jobs.
export function buildRecentJobs(): ReadonlyArray<RecentJobsPayload_ResultsFields> {
  return [
    buildRecentJob({
      id: '1',
      name: 'job 1',
    }),
    buildRecentJob({
      id: '2',
      name: 'job 2',
    }),
  ]
}
