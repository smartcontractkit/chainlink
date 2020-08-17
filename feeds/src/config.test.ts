import { partialAsFull } from '@chainlink/ts-helpers'
import { Config } from './config'

describe('config', () => {
  const originalEnv = { ...process.env }
  beforeEach(() => {
    process.env = originalEnv
  })

  it('returns infura key from the process env', () => {
    process.env.REACT_APP_INFURA_KEY = 'infuraKey'
    expect(Config.infuraKey()).toEqual('infuraKey')
  })

  it('returns the google analytics id from the process env', () => {
    process.env.REACT_APP_GA_ID = 'ga-id'
    expect(Config.gaId()).toEqual('ga-id')
  })

  it('returns an array of trimmed whitelist hostnames', () => {
    process.env.REACT_APP_HOSTNAME_WHITELIST = 'domain-a , domain-b'
    expect(Config.hostnameWhitelist()).toEqual(['domain-a', 'domain-b'])
  })

  it('returns an array of trimmed dev whitelist hostnames', () => {
    process.env.REACT_APP_DEV_HOSTNAME_WHITELIST = 'dev-domain-a , dev-domain-b'
    expect(Config.devHostnameWhitelist()).toEqual([
      'dev-domain-a',
      'dev-domain-b',
    ])
  })

  it('returns the dev provider from the process env', () => {
    process.env.REACT_APP_DEV_PROVIDER = 'http://mock-provider'
    expect(Config.devProvider()).toEqual('http://mock-provider')
  })

  describe('feeds json', () => {
    it('returns feeds json from the process env', () => {
      process.env.REACT_APP_FEEDS_JSON = 'https://test.dev/feeds.json'
      expect(Config.feedsJson()).toEqual('https://test.dev/feeds.json')
    })

    it('returns a default feeds json when undefined', () => {
      process.env.REACT_APP_FEEDS_JSON = undefined
      expect(Config.feedsJson()).toEqual('/feeds.json')
    })

    it('can override feeds json with a query parameter from the whitelisted hostnames', () => {
      process.env.REACT_APP_FEEDS_JSON = 'https://env-configured.dev/feeds.json'
      process.env.REACT_APP_HOSTNAME_WHITELIST = 'override.dev'

      const whitelistLocation = partialAsFull<Location>({
        search: '?feeds-json=https%3A%2F%2Foverride.dev%2Ffeeds.json',
      })
      expect(Config.feedsJson(process.env, whitelistLocation)).toEqual(
        'https://override.dev/feeds.json',
      )

      const unknownLocation = partialAsFull<Location>({
        search: '?feeds-json=https%3A%2F%2Funknown.dev%5Ffeeds.json',
      })
      expect(Config.feedsJson(process.env, unknownLocation)).toEqual(
        'https://env-configured.dev/feeds.json',
      )
    })
  })

  describe('nodes json', () => {
    it('returns nodes json from the process env', () => {
      process.env.REACT_APP_NODES_JSON = 'https://test.dev/nodes.json'
      expect(Config.nodesJson()).toEqual('https://test.dev/nodes.json')
    })

    it('returns a default nodes json when undefined', () => {
      process.env.REACT_APP_NODES_JSON = undefined
      expect(Config.nodesJson()).toEqual('/nodes.json')
    })

    it('can override nodes json with a query parameter from the whitelisted hostnames', () => {
      process.env.REACT_APP_NODES_JSON = 'https://env-configured.dev/nodes.json'
      process.env.REACT_APP_HOSTNAME_WHITELIST = 'override.dev'

      const whitelistLocation = partialAsFull<Location>({
        search: '?nodes-json=https%3A%2F%2Foverride.dev%2Fnodes.json',
      })
      expect(Config.nodesJson(process.env, whitelistLocation)).toEqual(
        'https://override.dev/nodes.json',
      )

      const unknownLocation = partialAsFull<Location>({
        search: '?nodes-json=https%3A%2F%2Funknown.dev%5Fnodes.json',
      })
      expect(Config.nodesJson(process.env, unknownLocation)).toEqual(
        'https://env-configured.dev/nodes.json',
      )
    })
  })
})
