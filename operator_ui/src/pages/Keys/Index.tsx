import React from 'react'

import Grid from '@material-ui/core/Grid'

import { AccountAddresses } from './AccountAddresses'
import Content from 'components/Content'
import { CSAKeys } from 'src/screens/KeyManagement/CSAKeys'
import { OCRKeys } from 'src/screens/KeyManagement/OCRKeys'
import { P2PKeys } from 'src/screens/KeyManagement/P2PKeys'
import { Feature, useFeatureFlag } from 'src/hooks/useFeatureFlag'

export const KeysIndex = () => {
  const isCSAKeysFeatureEnabled = useFeatureFlag(Feature.CSA)

  return (
    <Content>
      <Grid container>
        <Grid item xs={12}>
          <OCRKeys />
        </Grid>

        <Grid item xs={12}>
          <P2PKeys />
        </Grid>

        <AccountAddresses />

        <Grid item xs={12}>
          {isCSAKeysFeatureEnabled && <CSAKeys />}
        </Grid>
      </Grid>
    </Content>
  )
}

export default KeysIndex
