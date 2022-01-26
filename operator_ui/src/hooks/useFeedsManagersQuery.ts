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
        jobTypes
        isBootstrapPeer
        isConnectionActive
        bootstrapPeerMultiaddr
        createdAt
      }
    }
  }
`

export const useFeedsManagersQuery = () => {
  const response = useQuery<FetchFeedsManagers, FetchFeedsManagersVariables>(
    FEEDS_MANAGERS_QUERY,
  )

  return response
}
