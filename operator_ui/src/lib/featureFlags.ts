import React from 'react'

const FeatureFlagKey = 'features'

export enum Feature {
  CSAKeys = 'csa_keys',
  FeedsManager = 'feeds_manager',
}

type FeatureFlags = {
  [Property in Feature]: boolean
}

function getFeatureFlags(): FeatureFlags | null {
  // Retrieve the object from storage
  var retrievedObject = localStorage.getItem(FeatureFlagKey)

  if (retrievedObject) {
    return JSON.parse(retrievedObject) as FeatureFlags
  }

  return null
}

// useFeature is a hook which returns whether a feature is enabled/disabled.
//
// The flags are stored in localstorage. If you want to enable a feature you
// will have to manually update local storage with a stringified JSON object
// under the 'features' key.
//
// {"csa_keys": true, "feeds_manager": true}
export function useFeature(feature: Feature): boolean {
  const [isEnabled, setIsEnabled] = React.useState(false)

  function checkFlag(): boolean {
    const flags = getFeatureFlags()

    if (flags && flags.hasOwnProperty(feature)) {
      return flags[feature]
    }

    return false
  }

  React.useEffect(() => {
    // when app loaded
    setIsEnabled(checkFlag())

    // when storage updated
    window.addEventListener('storage', () => setIsEnabled(checkFlag()))
    return () =>
      window.removeEventListener('storage', () => setIsEnabled(checkFlag()))
  }, [])

  return isEnabled
}
