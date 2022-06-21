import React from 'react'

import { gql } from '@apollo/client'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { prism } from 'react-syntax-highlighter/dist/esm/styles/prism'

import Button from '@material-ui/core/Button'
import Chip from '@material-ui/core/Chip'
import ExpandMoreIcon from '@material-ui/icons/ExpandMore'
import ExpansionPanel from '@material-ui/core/ExpansionPanel'
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary'
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails'
import {
  createStyles,
  Theme,
  WithStyles,
  withStyles,
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'

import { ConfirmationDialog } from 'src/components/Dialogs/ConfirmationDialog'
import {
  EditJobSpecDialog,
  Props as EditJobSpecDialogProps,
} from './EditJobSpecDialog'
import { SpecStatus } from 'src/types/generated/graphql'
import { TimeAgo } from 'src/components/TimeAgo'

export const JOB_PROPOSAL__SPECS_FIELDS = gql`
  fragment JobProposal_SpecsFields on JobProposalSpec {
    id
    definition
    status
    version
    createdAt
  }
`

const styles = (theme: Theme) => {
  return createStyles({
    versionText: {
      marginRight: theme.spacing.unit * 2,
    },
    proposedAtContainer: {
      flex: 1,
      textAlign: 'right',
    },
    expansionPanelDetails: {
      display: 'block',
    },
    actions: {
      display: 'flex',
      marginBottom: theme.spacing.unit * 2,
    },
    editContainer: {
      flex: 1,
    },
    actionsContainer: {},
  })
}

interface Props extends WithStyles<typeof styles> {
  specs: ReadonlyArray<JobProposal_SpecsFields>
  onApprove: (specID: string) => void
  onReject: (specID: string) => void
  onCancel: (specID: string) => void
  onUpdateSpec: EditJobSpecDialogProps['onSubmit']
}

interface ConfirmationDialogArgs {
  action: 'reject' | 'approve' | 'cancel'
  id: string
}

const confirmationDialogText = {
  approve: {
    title: 'Approve Job Proposal',
    body: 'Approving this job proposal will start running a new job. WARNING: If a job using the same contract address already exists, it will be deleted before running the new one.',
  },
  cancel: {
    title: 'Cancel Job Proposal',
    body: 'Cancelling this job proposal will delete the running job. Are you sure you want to cancel this job proposal?',
  },
  reject: {
    title: 'Reject Job Proposal',
    body: 'Are you sure you want to reject this job proposal?',
  },
}

export const SpecsView = withStyles(styles)(
  ({ classes, onApprove, onCancel, onReject, onUpdateSpec, specs }: Props) => {
    const [confirmationDialog, setConfirmationDialog] =
      React.useState<ConfirmationDialogArgs | null>(null)
    const [isEditing, setIsEditing] = React.useState(false)

    const openConfirmationDialog = (
      action: ConfirmationDialogArgs['action'],
      id: string,
    ) => {
      setConfirmationDialog({ action, id })
    }

    const latestSpec = React.useMemo(() => {
      return specs.reduce((max, spec) =>
        max.version > spec.version ? max : spec,
      )
    }, [specs])

    const sortedSpecs = React.useMemo(() => {
      const sorted = [...specs]

      return sorted.sort((a, b) => b.version - a.version)
    }, [specs])

    const renderActions = (status: SpecStatus, specID: string) => {
      switch (status) {
        case 'PENDING':
          return (
            <>
              <Button
                variant="text"
                color="secondary"
                onClick={() => openConfirmationDialog('reject', specID)}
              >
                Reject
              </Button>

              {latestSpec.id === specID && (
                <Button
                  variant="contained"
                  color="primary"
                  onClick={() => openConfirmationDialog('approve', specID)}
                >
                  Approve
                </Button>
              )}
            </>
          )
        case 'APPROVED':
          return (
            <Button
              variant="contained"
              onClick={() => openConfirmationDialog('cancel', specID)}
            >
              Cancel
            </Button>
          )
        case 'CANCELLED':
          if (latestSpec.id !== specID) {
            return null
          }

          return (
            <Button
              variant="contained"
              color="primary"
              onClick={() => openConfirmationDialog('approve', specID)}
            >
              Approve
            </Button>
          )

        default:
          return null
      }
    }

    return (
      <div>
        {sortedSpecs.map((spec, idx) => (
          <ExpansionPanel defaultExpanded={idx === 0} key={idx}>
            <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
              <Typography className={classes.versionText}>
                Version {spec.version}
              </Typography>
              <Chip
                label={spec.status}
                color={spec.status === 'APPROVED' ? 'primary' : 'default'}
                variant={
                  spec.status === 'REJECTED' || spec.status === 'CANCELLED'
                    ? 'outlined'
                    : 'default'
                }
              />
              <div className={classes.proposedAtContainer}>
                <Typography>
                  Proposed <TimeAgo tooltip>{spec.createdAt}</TimeAgo>
                </Typography>
              </div>
            </ExpansionPanelSummary>
            <ExpansionPanelDetails className={classes.expansionPanelDetails}>
              <div className={classes.actions}>
                <div className={classes.editContainer}>
                  {idx === 0 &&
                    (spec.status === 'PENDING' ||
                      spec.status === 'CANCELLED') && (
                      <Button
                        variant="contained"
                        onClick={() => setIsEditing(true)}
                      >
                        Edit
                      </Button>
                    )}
                </div>
                <div className={classes.actionsContainer}>
                  {renderActions(spec.status, spec.id)}
                </div>
              </div>

              <SyntaxHighlighter
                language="toml"
                style={prism}
                data-testid="codeblock"
              >
                {spec.definition}
              </SyntaxHighlighter>
            </ExpansionPanelDetails>
          </ExpansionPanel>
        ))}

        <ConfirmationDialog
          open={confirmationDialog != null}
          title={
            confirmationDialog
              ? confirmationDialogText[confirmationDialog.action].title
              : ''
          }
          body={
            confirmationDialog
              ? confirmationDialogText[confirmationDialog.action].body
              : ''
          }
          onConfirm={() => {
            if (confirmationDialog) {
              switch (confirmationDialog.action) {
                case 'approve':
                  onApprove(confirmationDialog.id)

                  break
                case 'cancel':
                  onCancel(confirmationDialog.id)

                  break
                case 'reject':
                  onReject(confirmationDialog.id)

                  break
                default:
                // NOOP
              }

              setConfirmationDialog(null)
            }
          }}
          cancelButtonText="Cancel"
          onCancel={() => setConfirmationDialog(null)}
        />

        <EditJobSpecDialog
          open={isEditing}
          onClose={() => setIsEditing(false)}
          initialValues={{
            definition: latestSpec.definition,
            id: latestSpec.id,
          }}
          onSubmit={onUpdateSpec}
        />
      </div>
    )
  },
)
