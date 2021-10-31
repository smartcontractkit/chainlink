import { gql, useQuery } from '@apollo/client'

export const FETCH_FEEDS_MANAGERS = gql`
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

export const useFetchFeedsManagers = () => {
  const response = useQuery<FetchFeedsManagers, FetchFeedsManagersVariables>(
    FETCH_FEEDS_MANAGERS,
  )

  return response
}
