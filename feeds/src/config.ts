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

export class Config {
  static infuraKey(env = process.env): string {
    if (env.REACT_APP_INFURA_KEY === undefined) {
      return ''
    }
    return env.REACT_APP_INFURA_KEY
  }

  static gaId(env = process.env): string {
    if (env.REACT_APP_GA_ID === undefined) {
      return ''
    }
    return env.REACT_APP_GA_ID
  }

  static feedsJson(env = process.env): string {
    if (env.REACT_APP_FEEDS_JSON === undefined) {
      return 'https://weiwatchers.com/feeds.json'
    }
    return env.REACT_APP_FEEDS_JSON
  }

  static nodesJson(env = process.env): string {
    if (env.REACT_APP_NODES_JSON === undefined) {
      return 'https://weiwatchers.com/nodes.json'
    }
    return env.REACT_APP_NODES_JSON
  }
}
