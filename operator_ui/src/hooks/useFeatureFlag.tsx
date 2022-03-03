import React from 'react'

import { gql, useQuery } from '@apollo/client'

export enum Feature {
  CSA = 'csa',
  FeedsManager = 'feeds_manager',
}

export const FEATURES_QUERY = gql`
  query FetchFeatures {
    features {
      ... on Features {
        __typename
        csa
        feedsManager
      }
    }
  }
`

// useFeature is a hook which returns whether a feature is enabled/disabled.
export function useFeatureFlag(feature: Feature): boolean {
  const { data } = useQuery<FetchFeatures, FetchFeaturesVariables>(
    FEATURES_QUERY,
    { fetchPolicy: 'network-only' },
  )

  const [isEnabled, setIsEnabled] = React.useState(false)

  React.useEffect(() => {
    if (data?.features.__typename == 'Features') {
      switch (feature) {
        case Feature.CSA:
          setIsEnabled(data?.features.csa)

          break
        case Feature.FeedsManager:
          setIsEnabled(data?.features.feedsManager)

          break
        default:
          setIsEnabled(false)
      }
    }
  }, [data, feature])

  return isEnabled
}
