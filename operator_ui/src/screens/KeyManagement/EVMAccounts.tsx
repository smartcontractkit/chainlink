import React from 'react'

import { gql, useQuery } from '@apollo/client'

import {
  EVMAccountsCard,
  ETH_KEYS_PAYLOAD__RESULTS_FIELDS,
} from './EVMAccountsCard'

export const ETH_KEYS_QUERY = gql`
  ${ETH_KEYS_PAYLOAD__RESULTS_FIELDS}
  query FetchETHKeys {
    ethKeys {
      results {
        ...ETHKeysPayload_ResultsFields
      }
    }
  }
`

export const EVMAccounts = () => {
  const { data, loading, error } = useQuery<
    FetchEthKeys,
    FetchEthKeysVariables
  >(ETH_KEYS_QUERY, {
    fetchPolicy: 'cache-and-network',
  })

  return (
    <EVMAccountsCard loading={loading} data={data} errorMsg={error?.message} />
  )
}
