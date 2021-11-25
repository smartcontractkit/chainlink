import React from 'react'

import { gql, useQuery } from '@apollo/client'

import { ChainsView, CHAINS_PAYLOAD__RESULTS_FIELDS } from './ChainsView'
import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'
import { Loading } from 'src/components/Feedback/Loading'
import { useQueryParams } from 'src/hooks/useQueryParams'

export const CHAINS_QUERY = gql`
  ${CHAINS_PAYLOAD__RESULTS_FIELDS}
  query FetchChains($offset: Int, $limit: Int) {
    chains(offset: $offset, limit: $limit) {
      results {
        ...ChainsPayload_ResultsFields
      }
      metadata {
        total
      }
    }
  }
`

export const ChainsScreen = () => {
  const qp = useQueryParams()
  const page = parseInt(qp.get('page') || '1', 10)
  const pageSize = parseInt(qp.get('per') || '50', 10)

  const { data, loading, error } = useQuery<FetchChains, FetchChainsVariables>(
    CHAINS_QUERY,
    {
      variables: { offset: (page - 1) * pageSize, limit: pageSize },
      fetchPolicy: 'network-only',
    },
  )

  if (loading) {
    return <Loading />
  }

  if (error) {
    return <GraphqlErrorHandler error={error} />
  }

  if (data) {
    return (
      <ChainsView
        chains={data.chains.results}
        page={page}
        pageSize={pageSize}
        total={data.chains.metadata.total}
      />
    )
  }

  return null
}
