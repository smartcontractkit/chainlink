import React from 'react'
import { useHistory } from 'react-router-dom'

import { v2 } from 'api'
import { FeedsManagerForm, FormValues } from './FeedsManagerForm'
import { FeedsManager, Resource } from 'core/store/models'

import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Grid from '@material-ui/core/Grid'

interface Props {
  manager: Resource<FeedsManager>
  onSuccess?: (manager: Resource<FeedsManager>) => void
}

export const EditFeedsManagerView: React.FC<Props> = ({
  manager,
  onSuccess,
}) => {
  const history = useHistory()
  const initialValues: FormValues = {
    name: manager.attributes.name,
    uri: manager.attributes.uri,
    jobTypes: manager.attributes.jobTypes,
    publicKey: manager.attributes.publicKey,
    isBootstrapPeer: manager.attributes.isBootstrapPeer,
    bootstrapPeerMultiaddr: manager.attributes.bootstrapPeerMultiaddr,
  }

  const handleSubmit = async (values: FormValues) => {
    try {
      const res = await v2.feedsManagers.updateFeedsManager(manager.id, values)

      if (onSuccess) {
        history.push('/feeds_manager')
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
          <CardHeader title="Edit Feeds Manager" />
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
