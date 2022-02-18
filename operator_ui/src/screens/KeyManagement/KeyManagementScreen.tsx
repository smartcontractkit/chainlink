import React from 'react'

import { KeyManagementView } from './KeyManagementView'
import { Feature, useFeatureFlag } from 'src/hooks/useFeatureFlag'

export const KeyManagementScreen = () => {
  const isCSAKeysFeatureEnabled = useFeatureFlag(Feature.CSA)

  return <KeyManagementView isCSAKeysFeatureEnabled={isCSAKeysFeatureEnabled} />
}
