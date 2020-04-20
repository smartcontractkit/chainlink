export interface FeedConfig {
  contractAddress: string
  contractType: 'aggregator'
  valuePrefix: '$' | 'Ξ' | '£' | '¥'
  name: string
  pair: string[]
  path: string
  networkId: number
  history: boolean
  threshold: number
  listing: boolean

  heartbeat?: number
  compareOffchain?: string
  healthPrice?: string
  multiply?: string
  sponsored?: string[]
  decimalPlaces?: number
  contractVersion?: 1 | 2 | 3
}

export interface OracleNode {
  address: string
  name: string
  networkId: number
}

export interface Config {
  feedsJson: string
  nodesJson: string
}

export const DefaultConfig: Config = {
  feedsJson:
    process.env.REACT_APP_FEEDS_JSON ?? 'https://feeds.chain.link/feeds.json',
  nodesJson:
    process.env.REACT_APP_NODES_JSON ?? 'https://feeds.chain.link/nodes.json',
}
