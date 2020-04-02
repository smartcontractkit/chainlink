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

class UrlConfig {
  static feedsJson(location: Location): string | undefined {
    const query = new URLSearchParams(location.search)
    const feedsJson = query.get('feeds-json')

    if (feedsJson) {
      return decodeURIComponent(feedsJson)
    }
    return undefined
  }

  static nodesJson(location: Location): string | undefined {
    const query = new URLSearchParams(location.search)
    const nodesJson = query.get('nodes-json')

    if (nodesJson) {
      return decodeURIComponent(nodesJson)
    }
    return undefined
  }
}

export class Config {
  static infuraKey(env = process.env): string | undefined {
    return env.REACT_APP_INFURA_KEY
  }

  static gaId(env = process.env): string | undefined {
    return env.REACT_APP_GA_ID
  }

  static devProvider(env = process.env): string | undefined {
    return env.REACT_APP_DEV_PROVIDER
  }

  static feedsJson(env = process.env, location = window.location): string {
    const urlFeedsJson = UrlConfig.feedsJson(location)
    if (urlFeedsJson) {
      return urlFeedsJson
    }
    return env.REACT_APP_FEEDS_JSON ?? '/feeds.json'
  }

  static nodesJson(env = process.env, location = window.location): string {
    const urlNodesJson = UrlConfig.nodesJson(location)
    if (urlNodesJson) {
      return urlNodesJson
    }
    return env.REACT_APP_NODES_JSON ?? '/nodes.json'
  }
}
