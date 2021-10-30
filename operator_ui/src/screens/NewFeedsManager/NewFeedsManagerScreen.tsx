import React from 'react'
import { Redirect, useHistory, useLocation } from 'react-router-dom'
import { gql, useQuery } from '@apollo/client'

import { v2 } from 'api'
import { FormValues } from 'components/Forms/FeedsManagerForm'
import { NewFeedsManagerView } from './NewFeedsManagerView'

import { FetchFeeds, FetchFeedsVariables } from 'types/feeds_manager'

export const FETCH_FEEDS = gql`
  query FetchFeeds {
    feedsManagers {
      results {
        __typename
        id
      }
    }
  }
`

export const NewFeedsManagerScreen: React.FC = () => {
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

  // TODO - Use GQL API
  const handleSubmit = async (values: FormValues) => {
    try {
      await v2.feedsManagers.createFeedsManager(values)

      history.push('/feeds_manager')
    } catch (e) {
      console.log(e)
    }
  }

  if (manager) {
    return (
      <Redirect
        to={{
          pathname: '/feeds_manager',
          state: { from: location },
        }}
      />
    )
  }

  return <NewFeedsManagerView onSubmit={handleSubmit} />
}
