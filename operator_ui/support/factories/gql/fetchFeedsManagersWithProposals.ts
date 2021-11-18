import { FeedsManager, JobProposals } from 'screens/FeedsManager/types'

import isoDate, { MINUTE_MS } from 'test-helpers/isoDate'

// buildFeedsManager builds a feeds manager for a FetchFeedsManagersWithProposals
// query
export function buildFeedsManager(
  overrides?: Partial<FeedsManager>,
): FeedsManager {
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
    createdAt: new Date(),
    jobProposals: [],
    ...overrides,
  }
}

// buildJobProposals builds a list of job proposals each containing a different
// status for a FetchFeedsManagersWithProposals query
export function buildJobProposals(): JobProposals {
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
