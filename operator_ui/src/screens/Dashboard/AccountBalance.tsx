import React from 'react'

import { gql, useQuery } from '@apollo/client'

import {
  AccountBalanceCard,
  ACCOUNT_BALANCES_PAYLOAD__RESULTS_FIELDS,
} from './AccountBalanceCard'

export const ACCOUNT_BALANCES_QUERY = gql`
  ${ACCOUNT_BALANCES_PAYLOAD__RESULTS_FIELDS}
  query FetchAccountBalances {
    ethKeys {
      results {
        ...AccountBalancesPayload_ResultsFields
      }
    }
  }
`

export const AccountBalance = () => {
  const { data, loading, error } = useQuery<
    FetchAccountBalances,
    FetchAccountBalancesVariables
  >(ACCOUNT_BALANCES_QUERY, {
    fetchPolicy: 'cache-and-network',
  })

  return (
    <AccountBalanceCard
      data={data}
      errorMsg={error?.message}
      loading={loading}
    />
  )
}
