import { FeedsManager } from 'src/types/generated/graphql'

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
