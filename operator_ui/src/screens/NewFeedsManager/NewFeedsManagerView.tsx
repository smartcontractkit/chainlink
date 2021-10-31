import React from 'react'

import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Grid from '@material-ui/core/Grid'

import {
  FeedsManagerForm,
  Props as FormProps,
} from 'components/Forms/FeedsManagerForm'

const initialValues = {
  name: 'Chainlink Feeds Manager',
  uri: 'localhost:8080',
  jobTypes: [],
  publicKey: '1bf8600b908b1411bcef8fc12f6268d9cd4bfb67981895e26502d6870042406e',
  isBootstrapPeer: false,
  bootstrapPeerMultiaddr: undefined,
}

type Props = Pick<FormProps, 'onSubmit'>

export const NewFeedsManagerView: React.FC<Props> = ({ onSubmit }) => {
  return (
    <Grid container>
      <Grid item xs={12} md={11} lg={9}>
        <Card>
          <CardHeader title="Register Feeds Manager" />
          <CardContent>
            <FeedsManagerForm
              initialValues={initialValues}
              onSubmit={onSubmit}
            />
          </CardContent>
        </Card>
      </Grid>
    </Grid>
  )
}
