import React from 'react'
import { useDispatch } from 'react-redux'
import { useParams } from 'react-router-dom'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { prism } from 'react-syntax-highlighter/dist/esm/styles/prism'

import { notifySuccess, notifyError } from 'actionCreators'
import { v2 } from 'api'
import { ConfirmationDialog } from 'src/components/Dialogs/ConfirmationDialog'
import Content from 'components/Content'
import ErrorMessage from 'components/Notifications/DefaultError'
import { Resource, JobProposal } from 'core/store/models'
import titleize from 'src/utils/titleize'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'

import Button from '@material-ui/core/Button'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Grid from '@material-ui/core/Grid'
import { EditJobSpecDialog, FormValues } from './EditJobSpecDialog'

interface RouteParams {
  id: string
}

export const JobProposalScreen = () => {
  const dispatch = useDispatch()
  const { id } = useParams<RouteParams>()
  const [proposal, setProposal] = React.useState<Resource<JobProposal>>()
  const [confirmApprove, setConfirmApprove] = React.useState(false)
  const [confirmReject, setConfirmReject] = React.useState(false)
  const [isEditing, setIsEditing] = React.useState(false)
  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !proposal)

  React.useEffect(() => {
    v2.jobProposals
      .getJobProposal(id)
      .then((res) => {
        setProposal(res.data)
      })
      .catch(setError)
  }, [id, setError])

  const handleReject = () => {
    v2.jobProposals
      .rejectJobProposal(id)
      .then((res) => {
        setProposal(res.data)
        dispatch(notifySuccess(() => <>Job Proposal was rejected</>, {}))
      })
      .catch((e) => {
        dispatch(notifyError(ErrorMessage, e))
      })
      .finally(() => setConfirmReject(false))
  }

  const handleApprove = () => {
    v2.jobProposals
      .approveJobProposal(id)
      .then((res) => {
        setProposal(res.data)
        dispatch(notifySuccess(() => <>Job Proposal was approved</>, {}))
      })
      .catch((e) => {
        dispatch(notifyError(ErrorMessage, e))
      })
      .finally(() => setConfirmApprove(false))
  }

  const handleUpdateJobSpecSubmit = ({ spec }: FormValues) => {
    return v2.jobProposals
      .updateJobProposalSpec(id, { spec })
      .then((res) => {
        setProposal(res.data)
        dispatch(notifySuccess(() => <>Spec was updated</>, {}))
        setIsEditing(false)
      })
      .catch((e) => {
        dispatch(notifyError(ErrorMessage, e))
      })
  }

  return (
    <Content>
      <Grid container>
        <Grid item xs={12}>
          <ErrorComponent />
          <LoadingPlaceholder />

          {proposal && (
            <Card>
              <CardHeader
                title={`Job proposal #${proposal?.id}`}
                subheader={`Status: ${titleize(proposal.attributes.status)}`}
                action={
                  proposal.attributes.status === 'pending' && (
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
                }
              />

              <CardContent>
                <SyntaxHighlighter language="toml" style={prism}>
                  {proposal.attributes.spec}
                </SyntaxHighlighter>

                {proposal.attributes.status === 'pending' && (
                  <Button
                    variant="contained"
                    onClick={() => setIsEditing(true)}
                  >
                    Edit job spec
                  </Button>
                )}
              </CardContent>
            </Card>
          )}
        </Grid>
      </Grid>

      <ConfirmationDialog
        open={confirmApprove}
        title="Approve Job Proposal"
        body="Approving this job proposal will start running a new job"
        confirmButtonText="Confirm"
        onConfirm={handleApprove}
        cancelButtonText="Cancel"
        onCancel={() => setConfirmApprove(false)}
      />

      <ConfirmationDialog
        open={confirmReject}
        title="Reject Job Proposal"
        body="Are you sure you want to reject this job proposal?"
        onConfirm={handleReject}
        cancelButtonText="Cancel"
        onCancel={() => setConfirmReject(false)}
      />

      {proposal && (
        <EditJobSpecDialog
          open={isEditing}
          onClose={() => setIsEditing(false)}
          initialValues={{ spec: proposal.attributes.spec }}
          onSubmit={handleUpdateJobSpecSubmit}
        />
      )}
    </Content>
  )
}
