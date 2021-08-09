import React from 'react'

import { v2 } from 'api'
import { FeedsManagerForm, FormValues } from './FeedsManagerForm'
import { FeedsManager, Resource } from 'core/store/models'

import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Grid from '@material-ui/core/Grid'

const initialValues = {
  name: 'Chainlink Feeds Manager',
  uri: '',
  jobTypes: [],
  publicKey: '',
  isBootstrapPeer: false,
  bootstrapPeerMultiaddr: undefined,
}
interface Props {
  onSuccess?: (manager: Resource<FeedsManager>) => void
}

export const RegisterFeedsManagerView: React.FC<Props> = ({ onSuccess }) => {
  const handleSubmit = async (values: FormValues) => {
    try {
      const res = await v2.feedsManagers.createFeedsManager(values)

      if (onSuccess) {
        onSuccess(res.data)
      }
    } catch (e) {
      console.log(e)
    }
  }

  return (
    <Grid container>
      <Grid item xs={12} md={11} lg={9}>
        <Card>
          <CardHeader title="Register Feeds Manager" />
          <CardContent>
            <FeedsManagerForm
              initialValues={initialValues}
              onSubmit={handleSubmit}
            />
          </CardContent>
        </Card>
      </Grid>
    </Grid>
  )
}
