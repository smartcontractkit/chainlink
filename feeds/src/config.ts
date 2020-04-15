import feeds from './feeds.json'
import nodes from './nodes.json'

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

export function getFeedsConfig(): FeedConfig[] {
  return feeds as FeedConfig[]
}

export function getOracleNodes(): OracleNode[] {
  return nodes as OracleNode[]
}
