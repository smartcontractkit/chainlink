import React from 'react'
import Grid from '@material-ui/core/Grid'
import Content from 'components/Content'
import { CSAKeys } from 'src/screens/KeyManagement/CSAKeys'
import { OcrKeys } from './OcrKeys'
import { P2PKeys } from './P2PKeys'
import { AccountAddresses } from './AccountAddresses'
import { Feature, useFeature } from 'src/hooks/useFeatureFlag'

export const KeysIndex = () => {
  const isCSAKeysFeatureEnabled = useFeature(Feature.CSA)

  React.useEffect(() => {
    document.title = 'Keys and account addresses'
  }, [])

  return (
    <Content>
      <Grid container>
        <OcrKeys />
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
