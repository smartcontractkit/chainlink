import { gql, QueryHookOptions, useQuery } from '@apollo/client'

export const ETH_KEYS_PAYLOAD__RESULTS_FIELDS = gql`
  fragment ETHKeysPayload_ResultsFields on EthKey {
    address
    chain {
      id
    }
    createdAt
    ethBalance
    isFunding
    linkBalance
  }
`

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

// useEVMAccountsQuery fetches the EVM accounts.
export const useEVMAccountsQuery = (opts: QueryHookOptions = {}) => {
  return useQuery<FetchEthKeys>(ETH_KEYS_QUERY, opts)
}
