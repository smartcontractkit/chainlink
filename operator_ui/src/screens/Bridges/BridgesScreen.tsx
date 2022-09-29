import React from 'react'

import { gql, useQuery } from '@apollo/client'

import { BridgesView, BRIDGES_PAYLOAD__RESULTS_FIELDS } from './BridgesView'
import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'
import { Loading } from 'src/components/Feedback/Loading'
import { useQueryParams } from 'src/hooks/useQueryParams'

export const BRIDGES_QUERY = gql`
  ${BRIDGES_PAYLOAD__RESULTS_FIELDS}
  query FetchBridges($offset: Int, $limit: Int) {
    bridges(offset: $offset, limit: $limit) {
      results {
        ...BridgesPayload_ResultsFields
      }
      metadata {
        total
      }
    }
  }
`

export const BridgesScreen: React.FC = () => {
  const qp = useQueryParams()
  const page = parseInt(qp.get('page') || '1', 10)
  const pageSize = parseInt(qp.get('per') || '10', 10)

  const { data, loading, error } = useQuery<
    FetchBridges,
    FetchBridgesVariables
  >(BRIDGES_QUERY, {
    variables: { offset: (page - 1) * pageSize, limit: pageSize },
    fetchPolicy: 'network-only',
  })

  if (loading) {
    return <Loading />
  }

  if (error) {
    return <GraphqlErrorHandler error={error} />
  }

  if (data) {
    return (
      <BridgesView
        bridges={data.bridges.results}
        page={page}
        pageSize={pageSize}
        total={data.bridges.metadata.total}
      />
    )
  }

  return null
}
