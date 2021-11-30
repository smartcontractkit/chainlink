import React from 'react'

import Grid from '@material-ui/core/Grid'

import { AccountAddresses } from './AccountAddresses'
import Content from 'components/Content'
import { CSAKeys } from 'src/screens/KeyManagement/CSAKeys'
import { OCRKeys } from 'src/screens/KeyManagement/OCRKeys'
import { P2PKeys } from './P2PKeys'
import { Feature, useFeature } from 'src/hooks/useFeatureFlag'

export const KeysIndex = () => {
  const isCSAKeysFeatureEnabled = useFeature(Feature.CSA)

  return (
    <Content>
      <Grid container>
        <Grid item xs={12}>
          <OCRKeys />
        </Grid>

        <P2PKeys />
        <AccountAddresses />

        <Grid item xs={12}>
          {isCSAKeysFeatureEnabled && <CSAKeys />}
        </Grid>
      </Grid>
    </Content>
  )
}

export default KeysIndex
