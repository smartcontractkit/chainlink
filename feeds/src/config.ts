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
  contractVersion: 1 | 2 | 3
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

  static hostnameWhitelist(env = process.env): string[] {
    return splitHostnames(env.REACT_APP_HOSTNAME_WHITELIST)
  }

  static devHostnameWhitelist(env = process.env): string[] {
    return splitHostnames(env.REACT_APP_DEV_HOSTNAME_WHITELIST)
  }

  static devProvider(env = process.env): string | undefined {
    return env.REACT_APP_DEV_PROVIDER
  }

  static feedsJson(env = process.env, location = window.location): string {
    const queryOverride = UrlConfig.feedsJson(location)
    if (queryOverride) {
      const overrideUrl = new URL(queryOverride)
      if (Config.hostnameWhitelist().includes(overrideUrl.hostname)) {
        return queryOverride
      }
    }
    return env.REACT_APP_FEEDS_JSON ?? '/feeds.json'
  }

  static nodesJson(env = process.env, location = window.location): string {
    const queryOverride = UrlConfig.nodesJson(location)
    if (queryOverride) {
      const overrideUrl = new URL(queryOverride)
      if (Config.hostnameWhitelist().includes(overrideUrl.hostname)) {
        return queryOverride
      }
    }
    return env.REACT_APP_NODES_JSON ?? '/nodes.json'
  }
}

function splitHostnames(list?: string): string[] {
  return (list ?? '').split(',').map(s => s.trim())
}
