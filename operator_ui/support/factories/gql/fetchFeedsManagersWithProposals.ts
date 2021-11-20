import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'

// buildFeedsManagerResultFields is a convenience function to construct a result
// with default build values.
export function buildFeedsManagerResultFields(): FeedsManagerPayload_ResultsFields {
  return {
    ...buildFeedsManagerFields(),
    jobProposals: buildJobProposals(),
  }
}

// buildFeedsManagerFields builds the feeds manager fields for a
// FetchFeedsManagersWithProposals query.
export function buildFeedsManagerFields(
  overrides?: Partial<FeedsManagerFields>,
): FeedsManagerFields {
  return {
    __typename: 'FeedsManager',
    id: '1',
    name: 'Chainlink Feeds Manager',
    uri: 'localhost:8080',
    publicKey: '1111',
    jobTypes: ['FLUX_MONITOR'],
    isConnectionActive: false,
    isBootstrapPeer: false,
    bootstrapPeerMultiaddr: null,
    ...overrides,
  }
}

// buildJobProposals builds a list of job proposals each containing a different
// status for a FetchFeedsManagersWithProposals query
export function buildJobProposals(): FeedsManager_JobProposalsFields[] {
  const minuteAgo = isoDate(Date.now() - MINUTE_MS)

  return [
    {
      id: '1',
      proposedAt: minuteAgo,
      status: 'PENDING',
    },
    {
      id: '2',
      proposedAt: minuteAgo,
      externalJobID: '00000000-0000-0000-0000-000000000002',
      status: 'APPROVED',
    },
    {
      id: '3',
      proposedAt: minuteAgo,
      status: 'REJECTED',
    },
    {
      id: '4',
      proposedAt: minuteAgo,
      status: 'CANCELLED',
    },
  ]
}
