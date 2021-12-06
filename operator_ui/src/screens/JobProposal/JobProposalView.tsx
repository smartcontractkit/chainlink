import React from 'react'

import { gql } from '@apollo/client'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { prism } from 'react-syntax-highlighter/dist/esm/styles/prism'

import Button from '@material-ui/core/Button'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Grid from '@material-ui/core/Grid'

import { ConfirmationDialog } from 'src/components/Dialogs/ConfirmationDialog'
import Content from 'src/components/Content'
import {
  EditJobSpecDialog,
  Props as EditJobSpecDialogProps,
} from './EditJobSpecDialog'
import titleize from 'src/utils/titleize'

export const JOB_PROPOSAL_PAYLOAD_FIELDS = gql`
  fragment JobProposalPayloadFields on JobProposal {
    id
    externalJobID
    spec
    status
  }
`

interface Props {
  proposal: JobProposalPayloadFields
  onUpdateSpec: EditJobSpecDialogProps['onSubmit']
  onApprove: () => void
  onCancel: () => void
  onReject: () => void
}

export const JobProposalView: React.FC<Props> = ({
  onApprove,
  onCancel,
  onReject,
  onUpdateSpec,
  proposal,
}) => {
  const [isEditing, setIsEditing] = React.useState(false)
  const [confirmApprove, setConfirmApprove] = React.useState(false)
  const [confirmReject, setConfirmReject] = React.useState(false)
  const [confirmCancel, setConfirmCancel] = React.useState(false)

  const canEdit =
    proposal.status === 'PENDING' || proposal.status === 'CANCELLED'

  const renderActions = () => {
    switch (proposal.status) {
      case 'PENDING':
        return (
          <>
            <Button
              variant="text"
              color="secondary"
              onClick={() => setConfirmReject(true)}
            >
              Reject
            </Button>
            <Button
              variant="contained"
              color="primary"
              onClick={() => setConfirmApprove(true)}
            >
              Approve
            </Button>
          </>
        )
      case 'APPROVED':
        return (
          <>
            <Button variant="contained" onClick={() => setConfirmCancel(true)}>
              Cancel
            </Button>
          </>
        )
      case 'CANCELLED':
        return (
          <>
            <Button
              variant="contained"
              color="primary"
              onClick={() => setConfirmApprove(true)}
            >
              Approve
            </Button>
          </>
        )

      default:
        return null
    }
  }

  return (
    <Content>
      <Grid container>
        <Grid item xs={12}>
          <Card>
            <CardHeader
              title={`Job proposal #${proposal.id}`}
              subheader={`Status: ${titleize(proposal.status)}`}
              action={renderActions()}
            />

            <CardContent>
              <SyntaxHighlighter
                language="toml"
                style={prism}
                data-testid="codeblock"
              >
                {proposal.spec}
              </SyntaxHighlighter>

              {canEdit && (
                <Button variant="contained" onClick={() => setIsEditing(true)}>
                  Edit job spec
                </Button>
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      <ConfirmationDialog
        open={confirmApprove}
        title="Approve Job Proposal"
        body="Approving this job proposal will start running a new job"
        confirmButtonText="Confirm"
        onConfirm={() => {
          onApprove()
          setConfirmApprove(false)
        }}
        cancelButtonText="Cancel"
        onCancel={() => setConfirmApprove(false)}
      />

      <ConfirmationDialog
        open={confirmReject}
        title="Reject Job Proposal"
        body="Are you sure you want to reject this job proposal?"
        onConfirm={() => {
          onReject()
          setConfirmReject(false)
        }}
        cancelButtonText="Cancel"
        onCancel={() => setConfirmReject(false)}
      />

      <ConfirmationDialog
        open={confirmCancel}
        title="Cancel Job Proposal"
        body="Cancelling this job proposal will delete the running job. Are you sure you want to cancel this job proposal?"
        onConfirm={() => {
          onCancel()
          setConfirmCancel(false)
        }}
        cancelButtonText="Cancel"
        onCancel={() => setConfirmCancel(false)}
      />

      {canEdit && (
        <EditJobSpecDialog
          open={isEditing}
          onClose={() => setIsEditing(false)}
          initialValues={{ spec: proposal.spec }}
          onSubmit={onUpdateSpec}
        />
      )}
    </Content>
  )
}
