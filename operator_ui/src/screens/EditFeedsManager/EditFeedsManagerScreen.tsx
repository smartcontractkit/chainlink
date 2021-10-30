import React from 'react'
import { Redirect, useLocation, useHistory } from 'react-router-dom'
import { gql, useQuery } from '@apollo/client'

import { v2 } from 'api'
import { EditFeedsManagerView } from './EditFeedsManagerView'
import { FormValues } from 'components/Forms/FeedsManagerForm'

import { FetchFeeds, FetchFeedsVariables } from 'types/feeds_manager'

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

export const EditFeedsManagerScreen: React.FC = () => {
  const history = useHistory()
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

  if (!manager) {
    return (
      <Redirect
        to={{
          pathname: '/feeds_manager/new',
          state: { from: location },
        }}
      />
    )
  }

  const handleSubmit = async (values: FormValues) => {
    try {
      await v2.feedsManagers.updateFeedsManager(manager.id, values)

      history.push('/feeds_manager')
    } catch (e) {
      console.log(e)
    }
  }

  return <EditFeedsManagerView data={manager} onSubmit={handleSubmit} />
}
