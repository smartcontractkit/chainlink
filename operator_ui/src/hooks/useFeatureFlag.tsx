import React from 'react'

import { v2 } from 'api'

export enum Feature {
  CSA = 'csa',
  FeedsManager = 'feeds_manager',
}

// useFeature is a hook which returns whether a feature is enabled/disabled.
export function useFeature(feature: Feature): boolean {
  const [isEnabled, setIsEnabled] = React.useState(false)

  React.useEffect(() => {
    const fetch = () => {
      v2.features.getFeatureFlags().then((res) => {
        for (const flag of res.data) {
          if (flag.id == feature && flag.attributes.enabled) {
            setIsEnabled(true)
          }
        }
      })
    }

    fetch()
  }, [feature])

  return isEnabled
}
