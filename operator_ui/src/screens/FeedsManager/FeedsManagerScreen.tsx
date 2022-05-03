import React from 'react'

import { Redirect, useLocation } from 'react-router-dom'

import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'
import { FeedsManagerView } from './FeedsManagerView'
import { Loading } from 'src/components/Feedback/Loading'
import { useFeedsManagersWithProposalsQuery } from 'src/hooks/queries/useFeedsManagersWithProposalsQuery'

export const FeedsManagerScreen: React.FC = () => {
  const location = useLocation()

  const { data, loading, error } = useFeedsManagersWithProposalsQuery({
    fetchPolicy: 'cache-and-network',
  })

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
