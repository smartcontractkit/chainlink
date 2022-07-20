import { gql, QueryHookOptions, useQuery } from '@apollo/client'

const FEEDS_MANAGER__CHAIN_CONFIG_FIELDS = gql`
  fragment FeedsManager_ChainConfigFields on FeedsManagerChainConfig {
    id
    chainID
    chainType
    accountAddr
    adminAddr
    fluxMonitorJobConfig {
      enabled
    }
    ocr1JobConfig {
      enabled
      isBootstrap
      multiaddr
      p2pPeerID
      keyBundleID
    }
    ocr2JobConfig {
      enabled
      isBootstrap
      multiaddr
      p2pPeerID
      keyBundleID
    }
  }
`

export const FEEDS_MANAGER_FIELDS = gql`
  ${FEEDS_MANAGER__CHAIN_CONFIG_FIELDS}
  fragment FeedsManagerFields on FeedsManager {
    id
    name
    uri
    publicKey
    isConnectionActive
    chainConfigs {
      ...FeedsManager_ChainConfigFields
    }
  }
`

export const FEEDS_MANAGER__JOB_PROPOSAL_FIELDS = gql`
  fragment FeedsManager_JobProposalsFields on JobProposal {
    id
    externalJobID
    remoteUUID
    status
    pendingUpdate
    latestSpec {
      createdAt
      version
    }
  }
`

export const FEEDS_MANAGERS_PAYLOAD__RESULTS_FIELDS = gql`
  ${FEEDS_MANAGER_FIELDS}
  ${FEEDS_MANAGER__JOB_PROPOSAL_FIELDS}
  fragment FeedsManagerPayload_ResultsFields on FeedsManager {
    ...FeedsManagerFields
    jobProposals {
      ...FeedsManager_JobProposalsFields
    }
  }
`

export const FEEDS_MANAGERS_WITH_PROPOSALS_QUERY = gql`
  ${FEEDS_MANAGERS_PAYLOAD__RESULTS_FIELDS}
  query FetchFeedManagersWithProposals {
    feedsManagers {
      results {
        ...FeedsManagerPayload_ResultsFields
      }
    }
  }
`

export const useFeedsManagersWithProposalsQuery = (
  opts: QueryHookOptions = {},
) => {
  return useQuery<FetchFeedManagersWithProposals>(
    FEEDS_MANAGERS_WITH_PROPOSALS_QUERY,
    opts,
  )
}
