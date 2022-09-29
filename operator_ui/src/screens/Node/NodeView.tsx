import React from 'react'

import { gql } from '@apollo/client'

import Grid from '@material-ui/core/Grid'

import { ConfirmationDialog } from 'src/components/Dialogs/ConfirmationDialog'
import Content from 'components/Content'
import { NodeCard } from './NodeCard'
import { Heading1 } from 'src/components/Heading/Heading1'

export const NODE_PAYLOAD_FIELDS = gql`
  fragment NodePayload_Fields on Node {
    id
    name
    chain {
      id
    }
    httpURL
    wsURL
    createdAt
    state
  }
`

interface Props {
  node: NodePayload_Fields
  onDelete: () => void
}

export const NodeView = ({ node, onDelete }: Props) => {
  const [confirmDelete, setConfirmDelete] = React.useState(false)

  return (
    <>
      <Content>
        <Grid container spacing={16}>
          <Grid item xs={12}>
            <Heading1>{node.name}</Heading1>
          </Grid>

          <Grid item xs={12}>
            <NodeCard node={node} onDelete={() => setConfirmDelete(true)} />
          </Grid>
        </Grid>
      </Content>

      <ConfirmationDialog
        open={confirmDelete}
        title={`Delete ${node.name}`}
        body="This action cannot be undone and access to this page will be lost"
        confirmButtonText="Confirm"
        onConfirm={() => {
          onDelete()
          setConfirmDelete(false)
        }}
        cancelButtonText="Cancel"
        onCancel={() => setConfirmDelete(false)}
      />
    </>
  )
}
