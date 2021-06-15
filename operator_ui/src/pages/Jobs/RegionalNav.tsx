import React, { useMemo, useCallback, useState, useRef } from 'react'
import { connect } from 'react-redux'
import { Redirect, useLocation } from 'react-router-dom'

import { localizedTimestamp, TimeAgo } from 'components/TimeAgo'
import Card from '@material-ui/core/Card'
import Dialog from '@material-ui/core/Dialog'
import Grid from '@material-ui/core/Grid'
import List from '@material-ui/core/List'
import ListItem from '@material-ui/core/ListItem'
import Badge from '@material-ui/core/Badge'
import TextField from '@material-ui/core/TextField'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { ApiResponse } from 'utils/json-api-client'
import { JobSpec } from 'core/store/models'
import classNames from 'classnames'

import { createJobRun, createJobRunV2, deleteJobSpec } from 'actionCreators'
import BaseLink from 'components/BaseLink'
import Button from 'components/Button'
import CopyJobSpec from 'components/CopyJobSpec'
import Close from 'components/Icons/Close'
import Link from 'components/Link'
import ErrorMessage from 'components/Notifications/DefaultError'
import { JobData } from './sharedTypes'
import { isJobV2 } from 'pages/Jobs/utils'

const styles = (theme: Theme) =>
  createStyles({
    badgePadding: {
      paddingLeft: theme.spacing.unit * 2,
      paddingRight: theme.spacing.unit * 2,
      marginLeft: theme.spacing.unit * -2,
      marginRight: theme.spacing.unit * -2,
      lineHeight: '1rem',
    },
    container: {
      backgroundColor: theme.palette.common.white,
      padding: theme.spacing.unit * 5,
      paddingBottom: 0,
    },
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
    horizontalNav: {
      paddingBottom: 0,
    },
    horizontalNavItem: {
      display: 'inline',
      paddingLeft: 0,
      paddingRight: 0,
    },
    horizontalNavLink: {
      padding: `${theme.spacing.unit * 4}px ${theme.spacing.unit * 4}px`,
      textDecoration: 'none',
      display: 'inline-block',
      borderBottom: 'solid 1px',
      borderBottomColor: theme.palette.common.white,
      '&:hover': {
        borderBottomColor: theme.palette.primary.main,
      },
    },
    activeNavLink: {
      color: theme.palette.primary.main,
      borderBottomColor: theme.palette.primary.main,
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

const isWebInitiator = (
  initiators: ApiResponse<JobSpec>['data']['attributes']['initiators'],
) => initiators.find((initiator) => initiator.type === 'web')

const CreateRunSuccessNotification = ({ data }: any) => (
  <React.Fragment>
    Successfully created job run{' '}
    <BaseLink href={`/jobs/${data.attributes.jobId}/runs/${data.id}`}>
      {data.id}
    </BaseLink>
  </React.Fragment>
)

const DeleteSuccessNotification = ({ id }: any) => (
  <React.Fragment>Successfully deleted job {id}</React.Fragment>
)

interface Props extends WithStyles<typeof styles> {
  createJobRun: Function
  createJobRunV2: Function
  deleteJobSpec: Function
  jobSpecId: string
  externalJobID?: string
  job: JobData['job']
  runsCount: JobData['recentRunsCount']
  getJobSpecRuns: (props: { page?: number; size?: number }) => Promise<void>
}

const RegionalNavComponent = ({
  classes,
  createJobRun,
  createJobRunV2,
  jobSpecId,
  job,
  deleteJobSpec,
  getJobSpecRuns,
  runsCount,
  externalJobID,
}: Props) => {
  const location = useLocation()
  const navErrorsActive = location.pathname.endsWith('/errors')
  const navDefinitionActive = location.pathname.endsWith('/definition')
  const navRunsActive = location.pathname.endsWith('/runs')
  const navOverviewActive =
    !navDefinitionActive && !navErrorsActive && !navRunsActive
  const [modalOpen, setModalOpen] = useState(false)
  const [deleted, setDeleted] = useState(false)
  const [runJobModalOpen, setRunJobModalOpen] = useState(false)

  const handleRun = async (pipelineInput: string) => {
    const params = new URLSearchParams(location.search)
    const page = params.get('page')
    const size = params.get('size')

    if (job?.id && isJobV2(job.id)) {
      await createJobRunV2(
        externalJobID || jobSpecId,
        pipelineInput,
        CreateRunSuccessNotification,
        ErrorMessage,
      )
    } else {
      await createJobRun(jobSpecId, CreateRunSuccessNotification, ErrorMessage)
    }
    await getJobSpecRuns({
      page: page ? parseInt(page, 10) : undefined,
      size: size ? parseInt(size, 10) : undefined,
    })
  }
  const handleDelete = (id: string) => {
    deleteJobSpec(
      id,
      () => DeleteSuccessNotification({ id }),
      ErrorMessage,
      job?.type,
    )
    setDeleted(true)
  }

  const typeDetail = useMemo(() => {
    let type = 'unknown'

    if (!job) {
      return 'Unknown job type'
    }

    if (job.type === 'v2') {
      switch (job.specType) {
        case 'offchainreporting':
          type = 'Off-chain reporting'
          break
        case 'keeper':
          type = 'Keeper'
          break
        case 'cron':
          type = 'Cron'
          break
        case 'webhook':
          type = 'Webhook'
          break
        default:
          type = 'Direct request'
      }
    } else {
      type = job.type
    }

    return `${type} job spec detail`
  }, [job])

  const toggleRunJobModal = useCallback(() => {
    setRunJobModalOpen(!runJobModalOpen)
  }, [runJobModalOpen, setRunJobModalOpen])

  return (
    <>
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
                  <Button
                    variant="danger"
                    onClick={() => handleDelete(jobSpecId)}
                  >
                    Delete {jobSpecId}
                    {deleted && <Redirect to="/" />}
                  </Button>
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </Dialog>

      <Card className={classes.container}>
        <Grid container spacing={0}>
          <Grid item xs={12}>
            {job && (
              <Typography variant="subtitle2" color="secondary" gutterBottom>
                {typeDetail}
              </Typography>
            )}
          </Grid>
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
                    variant="h3"
                    color="secondary"
                    className={classes.jobSpecId}
                  >
                    {job.name || jobSpecId}
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
                    {((job.type === 'Direct request' &&
                      job.initiators &&
                      isWebInitiator(job.initiators)) ||
                      (job.type == 'v2' && job.specType == 'webhook')) && (
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
          <Grid item xs={12}>
            {job?.name && (
              <Typography variant="subtitle2" color="secondary" gutterBottom>
                {job.id}
              </Typography>
            )}
            {job?.createdAt && (
              <Typography variant="subtitle2" color="textSecondary">
                Created <TimeAgo tooltip={false}>{job.createdAt}</TimeAgo> (
                {localizedTimestamp(job.createdAt)})
              </Typography>
            )}
          </Grid>
          <Grid item xs={12}>
            <List className={classes.horizontalNav}>
              <ListItem className={classes.horizontalNavItem}>
                <Link
                  href={`/jobs/${jobSpecId}`}
                  className={classNames(
                    classes.horizontalNavLink,
                    navOverviewActive && classes.activeNavLink,
                  )}
                >
                  Overview
                </Link>
              </ListItem>
              <ListItem className={classes.horizontalNavItem}>
                <Link
                  href={`/jobs/${jobSpecId}/definition`}
                  className={classNames(
                    classes.horizontalNavLink,
                    navDefinitionActive && classes.activeNavLink,
                  )}
                >
                  Definition
                </Link>
              </ListItem>
              <ListItem className={classes.horizontalNavItem}>
                <Link
                  href={`/jobs/${jobSpecId}/errors`}
                  className={classNames(
                    classes.horizontalNavLink,
                    navErrorsActive && classes.activeNavLink,
                  )}
                >
                  {job?.errors && job.errors.length > 0 ? (
                    <Badge
                      badgeContent={job.errors.length}
                      color="error"
                      className={classes.badgePadding}
                    >
                      Errors
                    </Badge>
                  ) : (
                    'Errors'
                  )}
                </Link>
              </ListItem>
              <ListItem className={classes.horizontalNavItem}>
                <Link
                  href={`/jobs/${jobSpecId}/runs`}
                  className={classNames(
                    classes.horizontalNavLink,
                    navRunsActive && classes.activeNavLink,
                  )}
                >
                  <Badge
                    badgeContent={runsCount || 0}
                    color="primary"
                    className={classes.badgePadding}
                    max={99999}
                  >
                    Runs
                  </Badge>
                </Link>
              </ListItem>
            </List>
          </Grid>
        </Grid>
      </Card>
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
  createJobRun,
  createJobRunV2,
  deleteJobSpec,
})(RegionalNavComponent)

export const RegionalNav = withStyles(styles)(ConnectedRegionalNav)

export default RegionalNav
