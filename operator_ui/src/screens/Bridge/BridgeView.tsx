import React from 'react'

import { gql } from '@apollo/client'

import Grid from '@material-ui/core/Grid'

import { BridgeCard } from './BridgeCard'
import Content from 'components/Content'
import { ConfirmationDialog } from 'src/components/Dialogs/ConfirmationDialog'
import { Heading1 } from 'src/components/Heading/Heading1'

export const BRIDGE_PAYLOAD_FIELDS = gql`
  fragment BridgePayload_Fields on Bridge {
    id
    name
    url
    confirmations
    outgoingToken
    minimumContractPayment
  }
`

interface Props {
  bridge: BridgePayload_Fields
  onDelete: () => void
}

export const BridgeView = ({ bridge, onDelete }: Props) => {
  const [isDialogOpen, setIsDialogOpen] = React.useState(false)

  return (
    <>
      <Content>
        <Grid container spacing={16}>
          <Grid item xs={12}>
            <Heading1>{bridge.name}</Heading1>
          </Grid>

          <Grid item xs={12}>
            <BridgeCard
              bridge={bridge}
              onDelete={() => setIsDialogOpen(true)}
            />
          </Grid>
        </Grid>
      </Content>

      <ConfirmationDialog
        open={isDialogOpen}
        title="Delete Bridge"
        body="Are you sure you want to delete this bridge?"
        confirmButtonText="Confirm"
        onConfirm={onDelete}
        cancelButtonText="Cancel"
        onCancel={() => setIsDialogOpen(false)}
      />
    </>
  )
}
