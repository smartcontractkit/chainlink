import { gql, QueryHookOptions, useQuery } from '@apollo/client'

export const P2P_KEYS_PAYLOAD__RESULTS_FIELDS = gql`
  fragment P2PKeysPayload_ResultsFields on P2PKey {
    id
    peerID
    publicKey
  }
`

export const P2P_KEYS_QUERY = gql`
  ${P2P_KEYS_PAYLOAD__RESULTS_FIELDS}
  query FetchP2PKeys {
    p2pKeys {
      results {
        ...P2PKeysPayload_ResultsFields
      }
    }
  }
`

// useP2PKeysQuery fetches the chains
export const useP2PKeysQuery = (opts: QueryHookOptions = {}) => {
  return useQuery<FetchP2PKeys>(P2P_KEYS_QUERY, opts)
}
