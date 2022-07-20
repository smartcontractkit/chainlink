import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'

// buildJob builds a job for the FetchJobs query.
export function buildJob(
  overrides?: Partial<JobsPayload_ResultsFields>,
): JobsPayload_ResultsFields {
  const minuteAgo = isoDate(Date.now() - MINUTE_MS)

  return {
    __typename: 'Job',
    id: '1',
    name: 'job 1',
    externalJobID: '00000000-0000-0000-0000-000000000001',
    spec: {
      __typename: 'FluxMonitorSpec',
    },
    createdAt: minuteAgo,
    ...overrides,
  }
}

// buildsJobs builds a list of jobs.
export function buildJobs(): ReadonlyArray<JobsPayload_ResultsFields> {
  return [
    buildJob({
      id: '1',
      name: 'job 1',
      externalJobID: '00000000-0000-0000-0000-000000000001',
      spec: {
        __typename: 'FluxMonitorSpec',
      },
    }),
    buildJob({
      id: '2',
      name: 'job 2',
      externalJobID: '00000000-0000-0000-0000-000000000002',
      spec: {
        __typename: 'OCRSpec',
        contractAddress: '0x0000000000000000000000000000000000000001',
        keyBundleID: 'keybundleid',
        transmitterAddress: 'transmitteraddress',
      },
    }),
  ]
}
