import React from 'react'

import { gql, useQuery } from '@apollo/client'
import { useParams } from 'react-router-dom'

import {
  TransactionView,
  ETH_TRANSACTION_PAYLOAD_FIELDS,
} from './TransactionView'
import { Loading } from 'src/components/Feedback/Loading'
import { GraphqlErrorHandler } from 'src/components/ErrorHandler/GraphqlErrorHandler'
import NotFound from 'src/pages/NotFound'

export const ETH_TRANSACTION_QUERY = gql`
  ${ETH_TRANSACTION_PAYLOAD_FIELDS}
  query FetchEthTransaction($hash: ID!) {
    ethTransaction(hash: $hash) {
      __typename
      ... on EthTransaction {
        ...EthTransactionPayloadFields
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

export const TransactionScreen = () => {
  const { id } = useParams<RouteParams>()
  const { data, loading, error } = useQuery<
    FetchEthTransaction,
    FetchEthTransactionVariables
  >(ETH_TRANSACTION_QUERY, {
    variables: { hash: id },
  })

  if (loading) {
    return <Loading />
  }

  if (error) {
    return <GraphqlErrorHandler error={error} />
  }

  const payload = data?.ethTransaction
  switch (payload?.__typename) {
    case 'EthTransaction':
      return <TransactionView tx={payload} />
    case 'NotFoundError':
      return <NotFound />
    default:
      return null
  }
}
