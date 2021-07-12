import React from 'react'
import Grid from '@material-ui/core/Grid'
import Content from 'components/Content'
import { OcrKeys } from './OcrKeys'
import { P2PKeys } from './P2PKeys'
import { AccountAddresses } from './AccountAddresses'

export const KeysIndex = () => {
  React.useEffect(() => {
    document.title = 'Keys and account addresses'
  }, [])
  return (
    <Content>
      <Grid container>
        <OcrKeys />
        <P2PKeys />
        <AccountAddresses />
      </Grid>
    </Content>
  )
}

export default KeysIndex
