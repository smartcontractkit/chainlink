import React from 'react'
import { useDispatch } from 'react-redux'
import { gql, useQuery } from '@apollo/client'
import { useHistory, useLocation, useParams } from 'react-router'

import NotFound from 'src/pages/NotFound'
import { Loading } from 'src/components/Feedback/Loading'
import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'
import { ChainView, CHAIN_PAYLOAD__FIELDS } from './ChainView'

export const CHAIN_QUERY = gql`
  ${CHAIN_PAYLOAD__FIELDS}
  query FetchChain($id: ID!) {
    chain(id: $id) {
      __typename
      ... on Chain {
        ...ChainPayload_Fields
      }
      ... on NotFoundError {
        message
      }
    }
  }
`

interface RouteParams {
  id: string
}

export const ChainScreen = () => {
  const { id } = useParams<RouteParams>()
  const dispatch = useDispatch()
  const history = useHistory()
  const { search } = useLocation()

  const { data, loading, error, refetch } = useQuery<
    FetchChain,
    FetchChainVariables
  >(CHAIN_QUERY, {
    variables: { id },
    fetchPolicy: 'network-only',
  })

  if (loading) {
    return <Loading />
  }

  if (error) {
    return <GraphqlErrorHandler error={error} />
  }

  const payload = data?.chain
  switch (payload?.__typename) {
    case 'Chain':
      return (
        <ChainView
          chain={payload}
          //   onDelete={handleDelete}
          //   onRun={handleRun}
          //   refetch={handleRefetch}
          //   refetchRecentRuns={handleRefetchRecentRuns}
        />
      )
    case 'NotFoundError':
      return <NotFound />
    default:
      return null
  }
}
