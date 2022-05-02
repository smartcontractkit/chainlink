import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'

// buildJob builds the job payload fields for the FetchJob query.
export function buildJob(
  overrides?: Partial<JobPayload_Fields>,
): JobPayload_Fields {
  const minuteAgo = isoDate(Date.now() - MINUTE_MS)

  return {
    __typename: 'Job',
    id: '1',
    type: 'directrequest',
    schemaVersion: 1,
    name: 'direct request job',
    externalJobID: '00000000-0000-0000-0000-0000000000001',
    maxTaskDuration: '10s',
    spec: {
      __typename: 'DirectRequestSpec',
      contractAddress: '0x0000000000000000000000000000000000000000',
      evmChainID: '42',
      minIncomingConfirmations: 3,
      minIncomingConfirmationsEnv: false,
      minContractPaymentLinkJuels: '100000000000000',
      requesters: ['0x59bbE8CFC79c76857fE0eC27e67E4957370d72B5'],
    },
    runs: {
      results: [
        {
          __typename: 'JobRun',
          id: '1',
          allErrors: [],
          status: 'COMPLETED',
          createdAt: minuteAgo,
          finishedAt: minuteAgo,
        },
      ],
      metadata: { total: 0 },
    },
    observationSource:
      '    fetch    [type=http method=POST url="http://localhost:8001" requestData="{\\"hi\\": \\"hello\\"}"];\n    parse    [type=jsonparse path="data,result"];\n    multiply [type=multiply times=100];\n    fetch -> parse -> multiply;\n',
    createdAt: minuteAgo,
    errors: [],
    ...overrides,
  }
}

export function buildRun(
  overrides?: Partial<JobPayload_RunsFields>,
): JobPayload_RunsFields {
  const minuteAgo = isoDate(Date.now() - MINUTE_MS)

  return {
    __typename: 'JobRun',
    id: '1',
    allErrors: [],
    status: 'COMPLETED',
    createdAt: minuteAgo,
    finishedAt: minuteAgo,
    ...overrides,
  }
}

export function buildError(
  overrides?: Partial<JobPayload_ErrorsFields>,
): JobPayload_ErrorsFields {
  const minuteAgo = isoDate(Date.now() - MINUTE_MS)
  const twoMinutesAgo = isoDate(Date.now() - MINUTE_MS * 2)

  return {
    __typename: 'JobError',
    description: 'no contract at address',
    occurrences: 1,
    createdAt: twoMinutesAgo,
    updatedAt: minuteAgo,
    ...overrides,
  }
}
