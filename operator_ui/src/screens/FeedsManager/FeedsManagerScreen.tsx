import React from 'react'
import { Redirect, useLocation } from 'react-router-dom'
import { gql, useQuery } from '@apollo/client'

import { FetchFeeds, FetchFeedsVariables } from 'types/feeds_manager'
import { FeedsManagerView } from './FeedsManagerView'

export const FETCH_FEEDS = gql`
  query FetchFeeds {
    feedsManagers {
      results {
        __typename
        id
        name
        uri
        publicKey
        jobTypes
        isBootstrapPeer
        isConnectionActive
        bootstrapPeerMultiaddr
        createdAt
      }
    }
  }
`

export const FeedsManagerScreen: React.FC = () => {
  const location = useLocation()
  const { data, loading, error } = useQuery<FetchFeeds, FetchFeedsVariables>(
    FETCH_FEEDS,
  )

  if (loading) {
    return null
  }

  if (error) {
    return <div>error</div>
  }

  // We currently only support a single feeds manager, but plan to support more
  // in the future.
  const manager =
    data != undefined && data.feedsManagers.results[0]
      ? data.feedsManagers.results[0]
      : undefined

  if (manager) {
    return <FeedsManagerView data={manager} />
  }

  return (
    <Redirect
      to={{
        pathname: '/feeds_manager/new',
        state: { from: location },
      }}
    />
  )
}
