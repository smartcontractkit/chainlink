import React from 'react'

import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Grid from '@material-ui/core/Grid'

import {
  FeedsManagerForm,
  Props as FormProps,
} from 'components/Form/FeedsManagerForm'

type Props = {
  data: FetchFeedsManagers['feedsManagers']['results'][number]
} & Pick<FormProps, 'onSubmit'>

export const EditFeedsManagerView: React.FC<Props> = ({ data, onSubmit }) => {
  const initialValues = {
    name: data.name,
    uri: data.uri,
    publicKey: data.publicKey,
  }

  return (
    <Grid container>
      <Grid item xs={12} md={11} lg={9}>
        <Card>
          <CardHeader title="Edit Feeds Manager" />
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
