import React from 'react'

import Grid from '@material-ui/core/Grid'

import Content from 'components/Content'
import { EVMAccounts } from './EVMAccounts'
import { CSAKeys } from './CSAKeys'
import { OCRKeys } from './OCRKeys'
import { P2PKeys } from './P2PKeys'

interface Props {
  isCSAKeysFeatureEnabled: boolean
}

export const KeyManagementView: React.FC<Props> = ({
  isCSAKeysFeatureEnabled,
}) => {
  return (
    <Content>
      <Grid container>
        <Grid item xs={12}>
          <OCRKeys />
        </Grid>

        <Grid item xs={12}>
          <P2PKeys />
        </Grid>

        <Grid item xs={12}>
          <EVMAccounts />
        </Grid>

        <Grid item xs={12}>
          {isCSAKeysFeatureEnabled && <CSAKeys />}
        </Grid>
      </Grid>
    </Content>
  )
}
