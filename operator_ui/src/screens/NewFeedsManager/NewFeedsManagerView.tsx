import React from 'react'

import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Grid from '@material-ui/core/Grid'

import { FeedsManagerForm, FormValues } from 'components/Forms/FeedsManagerForm'

const initialValues = {
  name: 'Chainlink Feeds Manager',
  uri: '',
  jobTypes: [],
  publicKey: '',
  isBootstrapPeer: false,
  bootstrapPeerMultiaddr: undefined,
}
interface Props {
  onSubmit: (values: FormValues) => void
}

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
