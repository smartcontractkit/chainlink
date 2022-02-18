import React from 'react'

import { gql, useQuery } from '@apollo/client'

import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'
import {
  TransactionsView,
  ETH_TRANSACTIONS_PAYLOAD__RESULTS_FIELDS,
} from './TransactionsView'
import { useQueryParams } from 'src/hooks/useQueryParams'

export const ETH_TRANSACTIONS_QUERY = gql`
  ${ETH_TRANSACTIONS_PAYLOAD__RESULTS_FIELDS}
  query FetchEthTransactions($offset: Int, $limit: Int) {
    ethTransactions(offset: $offset, limit: $limit) {
      results {
        ...EthTransactionsPayload_ResultsFields
      }
      metadata {
        total
      }
    }
  }
`

export const TransactionsScreen = () => {
  const qp = useQueryParams()
  const page = parseInt(qp.get('page') || '1', 10)
  // Default set to 1000 until we can implement a server side search
  const pageSize = parseInt(qp.get('per') || '25', 10)

  const { data, loading, error } = useQuery<
    FetchEthTransactions,
    FetchEthTransactionsVariables
  >(ETH_TRANSACTIONS_QUERY, {
    variables: { offset: (page - 1) * pageSize, limit: pageSize },
    fetchPolicy: 'cache-and-network',
  })

  if (error) {
    return <GraphqlErrorHandler error={error} />
  }

  return (
    <TransactionsView
      data={data}
      loading={loading}
      page={page}
      pageSize={pageSize}
    />
  )
}
