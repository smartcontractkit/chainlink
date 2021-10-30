// TO BE REMOVED ONCE GENERATION WORKS
export interface FeedsManager {
  id: string
  name: string
  uri: string
  publicKey: string
  jobTypes: string[]
  isBootstrapPeer: boolean
  isConnectionActive: boolean
  bootstrapPeerMultiaddr: string
}

interface FetchFeeds {
  feedsManagers: {
    results: FeedsManager[]
  }
}

interface FetchFeedsVariables {}
