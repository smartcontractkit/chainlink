import React from 'react'
import { gql, useQuery } from '@apollo/client'
import { Redirect, Route, useRouteMatch } from 'react-router-dom'

import Content from 'components/Content'
import { useErrorHandler } from 'hooks/useErrorHandler'

import { EditFeedsManagerView } from './EditFeedsManagerView'
import { RegisterFeedsManagerView } from './RegisterFeedsManagerView'
import { FeedsManagerView } from './FeedsManagerView'

export interface FeedsManagerGQL {
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
  results: FeedsManagerGQL[]
}

interface FetchFeedsVariables {}

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
  const { path } = useRouteMatch()
  const { ErrorComponent } = useErrorHandler()

  const {
    data,
    loading,
    error: gqlError,
    refetch,
  } = useQuery<FetchFeeds, FetchFeedsVariables>(FETCH_FEEDS)

  // We currently only support a single feeds manager, but plan to support more
  // in the future.
  const manager =
    data != undefined && data.feedsManagers.results[0]
      ? data.feedsManagers.results[0]
      : undefined

  if (loading) {
    return null
  }

  if (gqlError) {
    return <ErrorComponent />
  }

  return (
    <Content>
      <Route
        exact
        path={`${path}/new`}
        render={({ location }) =>
          manager ? (
            <Redirect
              to={{
                pathname: '/feeds_manager',
                state: { from: location },
              }}
            />
          ) : (
            <RegisterFeedsManagerView onSuccess={refetch} />
          )
        }
      />

      <Route
        exact
        path={path}
        render={({ location }) =>
          manager ? (
            <FeedsManagerView manager={manager} />
          ) : (
            <Redirect
              to={{
                pathname: '/feeds_manager/new',
                state: { from: location },
              }}
            />
          )
        }
      />

      {/* <Route
        exact
        path={`${path}/edit`}
        render={({ location }) =>
          manager ? (
            <EditFeedsManagerView manager={manager} onSuccess={refetch} />
          ) : (
            <Redirect
              to={{
                pathname: '/feeds_manager',
                state: { from: location },
              }}
            />
          )
        }
      /> */}
    </Content>
  )
}
