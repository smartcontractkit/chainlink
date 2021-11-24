import React from 'react'

import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Grid from '@material-ui/core/Grid'

import BaseLink from 'src/components/BaseLink'
import { BridgeForm, Props as FormProps } from 'src/components/Form/BridgeForm'
import Button from 'src/components/Button'
import Content from 'components/Content'

type Props = Pick<FormProps, 'onSubmit'> & {
  bridge: BridgePayload_Fields
}

export const EditBridgeView: React.FC<Props> = ({ bridge, onSubmit }) => {
  const initialValues = {
    name: bridge.name,
    url: bridge.url,
    minimumContractPayment: bridge.minimumContractPayment,
    confirmations: bridge.confirmations,
  }

  return (
    <Content>
      <Grid container spacing={40}>
        <Grid item xs={12} md={11} lg={9}>
          <Card>
            <CardHeader
              title="Edit Bridge"
              action={
                <Button component={BaseLink} href={`/bridges/${bridge.id}`}>
                  Cancel
                </Button>
              }
            />

            <CardContent>
              <BridgeForm
                nameDisabled
                initialValues={initialValues}
                onSubmit={onSubmit}
                submitButtonText="Save Bridge"
              />
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Content>
  )
}
