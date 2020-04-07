import { FeedConfig, getFeedsConfig } from 'config'
import React from 'react'
import { match } from 'react-router'
import { Redirect } from 'react-router-dom'
import { Networks } from '../utils'

interface Params {
  pair?: string
  address?: string
}

interface Props {
  render: (config: FeedConfig) => any
  match: match<Params>
  networkId?: Networks
}

/**
 * WithFeedConfig enhancer
 *
 * Find a FeedConfig that satisfies the URL match params and inject it
 * into the rendered component. If a FeedConfig doesn't satisy the match,
 * this component redirects to the root of the application '/'
 */
const WithFeedConfig: React.FC<Props> = ({ render, match, networkId }) => {
  const config = getFeedsConfig().find(feedConfig => {
    if (match.params.pair) {
      return (
        compareInsensitive(feedConfig.path, match.params.pair) &&
        feedConfig.networkId === networkId
      )
    } else if (match.params.address) {
      return compareInsensitive(
        feedConfig.contractAddress,
        match.params.address,
      )
    } else {
      return false
    }
  })

  return config ? render(config) : <Redirect to={'/'} />
}

function compareInsensitive(a: string, b: string): boolean {
  return a.toLowerCase() === b.toLowerCase()
}

export default WithFeedConfig
