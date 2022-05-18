import { gql, useQuery } from '@apollo/client'

export const FEEDS_MANAGERS_QUERY = gql`
  query FetchFeedsManagers {
    feedsManagers {
      results {
        __typename
        id
        name
        uri
        publicKey
        isConnectionActive
        createdAt
      }
    }
  }
`

export const useFeedsManagersQuery = () => {
  return useQuery<FetchFeedsManagers, FetchFeedsManagersVariables>(
    FEEDS_MANAGERS_QUERY,
  )
}
