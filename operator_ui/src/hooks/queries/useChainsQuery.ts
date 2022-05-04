import { gql, QueryHookOptions, useQuery } from '@apollo/client'

export const CHAINS_PAYLOAD__RESULTS_FIELDS = gql`
  fragment ChainsPayload_ResultsFields on Chain {
    id
    enabled
    createdAt
  }
`

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

// useChainsQuery fetches the chains
export const useChainsQuery = (opts: QueryHookOptions = {}) => {
  return useQuery<FetchChains, FetchChainsVariables>(CHAINS_QUERY, opts)
}
