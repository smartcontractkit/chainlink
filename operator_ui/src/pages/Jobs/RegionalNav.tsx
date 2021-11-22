import React, { useCallback, useState, useRef } from 'react'
import { connect, useDispatch } from 'react-redux'
import { Redirect, useLocation } from 'react-router-dom'

import Dialog from '@material-ui/core/Dialog'
import Grid from '@material-ui/core/Grid'
import TextField from '@material-ui/core/TextField'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'

import * as api from 'api'
import { createJobRunV2 } from 'actionCreators'
import BaseLink from 'components/BaseLink'
import Button from 'components/Button'
import CopyJobSpec from 'components/CopyJobSpec'
import Close from 'components/Icons/Close'
import ErrorMessage from 'components/Notifications/DefaultError'
import { notifySuccess, notifyError } from 'actionCreators'
import { JobData } from './sharedTypes'

const styles = (theme: Theme) =>
  createStyles({
    mainRow: {
      marginBottom: theme.spacing.unit * 2,
    },
    actions: {
      textAlign: 'right',
    },
    regionalNavButton: {
      marginLeft: theme.spacing.unit,
      marginRight: theme.spacing.unit,
    },
    jobSpecId: {
      overflow: 'hidden',
      textOverflow: 'ellipsis',
    },
    dialogPaper: {
      minHeight: '240px',
      maxHeight: '240px',
      minWidth: '670px',
      maxWidth: '670px',
      overflow: 'hidden',
      borderRadius: theme.spacing.unit * 3,
    },
    warningText: {
      fontWeight: 500,
      marginLeft: theme.spacing.unit * 3,
      marginTop: theme.spacing.unit * 3,
      marginBottom: theme.spacing.unit,
    },
    closeButton: {
      marginRight: theme.spacing.unit * 3,
      marginTop: theme.spacing.unit * 3,
    },
    infoText: {
      fontSize: theme.spacing.unit * 2,
      fontWeight: 450,
      marginLeft: theme.spacing.unit * 6,
    },
    modalTextarea: {
      marginLeft: theme.spacing.unit * 2,
    },
    modalContent: {
      width: 'inherit',
    },
    deleteButton: {
      marginTop: theme.spacing.unit * 4,
    },
    runJobButton: {
      marginBottom: theme.spacing.unit * 3,
    },
    runJobModalContent: {
      overflow: 'hidden',
    },
  })

const CreateRunSuccessNotification = ({ data }: any) => (
  <React.Fragment>
    Successfully created job run{' '}
    <BaseLink href={`/jobs/${data.attributes.jobId}/runs/${data.id}`}>
      {data.id}
    </BaseLink>
  </React.Fragment>
)

const DeleteSuccessNotification = ({ id }: { id: string }) => (
  <React.Fragment>Successfully deleted job {id}</React.Fragment>
)
interface Props extends WithStyles<typeof styles> {
  createJobRunV2: Function
  jobId: string
  externalJobID?: string
  job: JobData['job']
  getJobSpecRuns: (props: { page?: number; size?: number }) => Promise<void>
}

const RegionalNavComponent = ({
  classes,
  createJobRunV2,
  jobId,
  job,
  getJobSpecRuns,
  externalJobID,
}: Props) => {
  const dispatch = useDispatch()
  const location = useLocation()
  const [modalOpen, setModalOpen] = useState(false)
  const [deleted, setDeleted] = useState(false)
  const [runJobModalOpen, setRunJobModalOpen] = useState(false)

  const handleRun = async (pipelineInput: string) => {
    const params = new URLSearchParams(location.search)
    const page = params.get('page')
    const size = params.get('size')

    await createJobRunV2(
      externalJobID || jobId,
      pipelineInput,
      CreateRunSuccessNotification,
      ErrorMessage,
    )

    await getJobSpecRuns({
      page: page ? parseInt(page, 10) : undefined,
      size: size ? parseInt(size, 10) : undefined,
    })
  }

  const handleDelete = async (id: string) => {
    try {
      await api.v2.jobs.destroyJobSpec(id)
      setDeleted(true)
      dispatch(notifySuccess(DeleteSuccessNotification, { id }))
    } catch (e: any) {
      dispatch(notifyError(ErrorMessage, e))

      setModalOpen(false)
    }
  }

  const toggleRunJobModal = useCallback(() => {
    setRunJobModalOpen(!runJobModalOpen)
  }, [runJobModalOpen, setRunJobModalOpen])

  return (
    <>
      <Grid container>
        <Grid item xs={12}>
          <Grid
            container
            spacing={0}
            alignItems="center"
            className={classes.mainRow}
          >
            <Grid item xs={6}>
              {job && (
                <Typography
                  variant="h5"
                  color="secondary"
                  className={classes.jobSpecId}
                >
                  {job.name || '--'}
                </Typography>
              )}
            </Grid>
            <Grid item xs={6} className={classes.actions}>
              {job && (
                <>
                  <Button
                    className={classes.regionalNavButton}
                    onClick={() => setModalOpen(true)}
                  >
                    Delete
                  </Button>
                  {job.specType == 'webhook' && (
                    <React.Fragment>
                      <Button
                        onClick={toggleRunJobModal}
                        className={classes.regionalNavButton}
                      >
                        Run
                      </Button>
                      <RunJobModal
                        open={runJobModalOpen}
                        onClose={toggleRunJobModal}
                        run={handleRun}
                        classes={classes}
                      />
                    </React.Fragment>
                  )}
                  {job.definition && (
                    <>
                      <Button
                        href={`/jobs/new?definition=${encodeURIComponent(
                          job.definition,
                        )}`}
                        component={BaseLink}
                        className={classes.regionalNavButton}
                      >
                        Duplicate
                      </Button>
                      <CopyJobSpec
                        JobSpec={job.definition}
                        className={classes.regionalNavButton}
                      />
                    </>
                  )}
                </>
              )}
            </Grid>
          </Grid>
        </Grid>
      </Grid>

      <Dialog
        open={modalOpen}
        classes={{ paper: classes.dialogPaper }}
        onClose={() => setModalOpen(false)}
      >
        <Grid container spacing={0}>
          <Grid item className={classes.modalContent}>
            <Grid container alignItems="baseline" justify="space-between">
              <Grid item>
                <Typography
                  variant="h5"
                  color="secondary"
                  className={classes.warningText}
                >
                  Warning: This Action Cannot Be Undone
                </Typography>
              </Grid>
              <Grid item>
                <Close
                  className={classes.closeButton}
                  onClick={() => setModalOpen(false)}
                />
              </Grid>
            </Grid>
            <Grid container direction="column">
              <Grid item>
                <Grid item>
                  <Typography
                    className={classes.infoText}
                    variant="h5"
                    color="secondary"
                  >
                    - All associated job runs will be deleted
                  </Typography>
                  <Typography
                    className={classes.infoText}
                    variant="h5"
                    color="secondary"
                  >
                    - Access to this page will be lost
                  </Typography>
                </Grid>
              </Grid>
              <Grid container spacing={0} alignItems="center" justify="center">
                <Grid item className={classes.deleteButton}>
                  <Button variant="danger" onClick={() => handleDelete(jobId)}>
                    Delete {jobId}
                    {deleted && <Redirect to="/" />}
                  </Button>
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </Dialog>
    </>
  )
}

const RunJobModal = (props: {
  open: boolean
  onClose: () => void
  run: (pipelineInput: string) => void
  classes: any
}) => {
  const { open, onClose, run, classes } = props

  const textarea = useRef<HTMLTextAreaElement>(null)

  const onClickRun = useCallback(() => {
    if (!textarea.current) {
      return
    }
    run(textarea.current.value)
    textarea.current.value = ''
    onClose()
  }, [run, onClose, textarea])

  return (
    <Dialog onClose={onClose} open={open}>
      <Grid container spacing={0} className={classes.runJobModalContent}>
        <Grid item className={classes.modalContent}>
          <Grid container alignItems="baseline" justify="space-between">
            <Grid item>
              <Typography
                variant="h5"
                color="secondary"
                className={classes.warningText}
              >
                Pipeline input
              </Typography>
            </Grid>
            <Grid item>
              <Close className={classes.closeButton} onClick={onClose} />
            </Grid>
          </Grid>
          <Grid container direction="column">
            <Grid item>
              <Grid item className={classes.modalTextarea}>
                <TextField
                  label="Multiline"
                  multiline
                  rows={4}
                  variant="outlined"
                  inputRef={textarea}
                />
              </Grid>
            </Grid>
            <Grid container spacing={0} alignItems="center" justify="center">
              <Grid item className={classes.runJobButton}>
                <Button variant="danger" onClick={onClickRun}>
                  Run job
                </Button>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </Dialog>
  )
}

export const ConnectedRegionalNav = connect(null, {
  createJobRunV2,
})(RegionalNavComponent)

export const RegionalNav = withStyles(styles)(ConnectedRegionalNav)

export default RegionalNav
