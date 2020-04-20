import { partialAsFull } from '@chainlink/ts-helpers'
import { Config } from './config'

describe('config', () => {
  it('returns infura key from the process env', () => {
    const processEnv = {
      NODE_ENV: 'test' as const,
      PUBLIC_URL: 'https://test.dev',
      REACT_APP_INFURA_KEY: 'infuraKey',
    }

    expect(Config.infuraKey(processEnv)).toEqual('infuraKey')
  })

  it('returns the google analytics id from the process env', () => {
    const processEnv = {
      NODE_ENV: 'test' as const,
      PUBLIC_URL: 'https://test.dev',
      REACT_APP_GA_ID: 'ga-id',
    }

    expect(Config.gaId(processEnv)).toEqual('ga-id')
  })

  describe('feeds json', () => {
    it('returns feeds json from the process env', () => {
      const processEnv = {
        NODE_ENV: 'test' as const,
        PUBLIC_URL: 'https://test.dev',
        REACT_APP_FEEDS_JSON: 'https://test.dev/feeds.json',
      }

      expect(Config.feedsJson(processEnv)).toEqual(
        'https://test.dev/feeds.json',
      )
    })

    it('returns a default feeds json when not provided', () => {
      const processEnv = {
        NODE_ENV: 'test' as const,
        PUBLIC_URL: 'https://test.dev',
      }

      expect(Config.feedsJson(processEnv)).toEqual(
        'https://weiwatchers.com/feeds.json',
      )
    })

    it('can override feeds json with a query parameter', () => {
      const location = partialAsFull<Location>({
        search: '?feeds-json=https%3A%2F%2Foverride.dev%2Ffeeds.json',
      })
      const processEnv = {
        NODE_ENV: 'test' as const,
        PUBLIC_URL: 'https://test.dev',
        REACT_APP_FEEDS_JSON: 'https://test.dev/feeds.json',
      }

      expect(Config.feedsJson(processEnv, location)).toEqual(
        'https://override.dev/feeds.json',
      )
    })
  })

  describe('nodes json', () => {
    it('returns nodes json from the process env', () => {
      const processEnv = {
        NODE_ENV: 'test' as const,
        PUBLIC_URL: 'https://test.dev',
        REACT_APP_NODES_JSON: 'https://test.dev/nodes.json',
      }

      expect(Config.nodesJson(processEnv)).toEqual(
        'https://test.dev/nodes.json',
      )
    })

    it('returns a default nodes json when not provided', () => {
      const processEnv = {
        NODE_ENV: 'test' as const,
        PUBLIC_URL: 'https://test.dev',
      }

      expect(Config.nodesJson(processEnv)).toEqual(
        'https://weiwatchers.com/nodes.json',
      )
    })

    it('can override nodes json with a query parameter', () => {
      const location = partialAsFull<Location>({
        search: '?nodes-json=https%3A%2F%2Foverride.dev%2Fnodes.json',
      })
      const processEnv = {
        NODE_ENV: 'test' as const,
        PUBLIC_URL: 'https://test.dev',
        REACT_APP_FEEDS_JSON: 'https://test.dev/nodes.json',
      }

      expect(Config.nodesJson(processEnv, location)).toEqual(
        'https://override.dev/nodes.json',
      )
    })
  })
})
