import React from 'react'

import { gql, useQuery } from '@apollo/client'
import { Redirect, useLocation } from 'react-router-dom'

import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'
import {
  FeedsManagerView,
  FEEDS_MANAGERS_PAYLOAD__RESULTS_FIELDS,
} from './FeedsManagerView'
import { Loading } from 'src/components/Feedback/Loading'

export const FEEDS_MANAGERS_WITH_PROPOSALS_QUERY = gql`
  ${FEEDS_MANAGERS_PAYLOAD__RESULTS_FIELDS}
  query FetchFeedManagersWithProposals {
    feedsManagers {
      results {
        ...FeedsManagerPayload_ResultsFields
      }
    }
  }
`

export const FeedsManagerScreen: React.FC = () => {
  const location = useLocation()

  const { data, loading, error } = useQuery<
    FetchFeedManagersWithProposals,
    FetchFeedManagersWithProposalsVariables
  >(FEEDS_MANAGERS_WITH_PROPOSALS_QUERY)

  if (loading) {
    return <Loading />
  }

  if (error) {
    return <GraphqlErrorHandler error={error} />
  }

  // We currently only support a single feeds manager, but plan to support more
  // in the future.
  const manager =
    data != undefined && data.feedsManagers.results.length > 0
      ? data.feedsManagers.results[0]
      : undefined

  if (data && manager) {
    return <FeedsManagerView manager={manager} />
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
