import React from 'react'

import { gql, useQuery } from '@apollo/client'

import { NodesView, NODES_PAYLOAD__RESULTS_FIELDS } from './NodesView'
import { useQueryParams } from 'src/hooks/useQueryParams'
import { Loading } from 'src/components/Feedback/Loading'
import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'

export const NODES_QUERY = gql`
  ${NODES_PAYLOAD__RESULTS_FIELDS}
  query FetchNodes($offset: Int, $limit: Int) {
    nodes(offset: $offset, limit: $limit) {
      results {
        ...NodesPayload_ResultsFields
      }
      metadata {
        total
      }
    }
  }
`

export const NodesScreen = () => {
  const qp = useQueryParams()
  const page = parseInt(qp.get('page') || '1', 10)
  const pageSize = parseInt(qp.get('per') || '100', 10)

  const { data, loading, error } = useQuery<FetchNodes, FetchNodesVariables>(
    NODES_QUERY,
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
      <NodesView
        nodes={data.nodes.results}
        page={page}
        pageSize={pageSize}
        total={data.nodes.metadata.total}
      />
    )
  }

  return null
}
