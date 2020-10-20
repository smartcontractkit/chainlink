import { localizedTimestamp, TimeAgo } from '@chainlink/styleguide'
import Card from '@material-ui/core/Card'
import Dialog from '@material-ui/core/Dialog'
import Grid from '@material-ui/core/Grid'
import List from '@material-ui/core/List'
import ListItem from '@material-ui/core/ListItem'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import classNames from 'classnames'
import React from 'react'
import { connect } from 'react-redux'
import { Redirect, useLocation } from 'react-router-dom'
import { createJobRun, deleteJobSpec, fetchJobRuns } from 'actionCreators'
import BaseLink from 'components/BaseLink'
import Button from 'components/Button'
import CopyJobSpec from 'components/CopyJobSpec'
import Close from 'components/Icons/Close'
import Link from 'components/Link'
import ErrorMessage from 'components/Notifications/DefaultError'
import jobSpecDefinition from 'utils/jobSpecDefinition'
import { isWebInitiator } from 'utils/jobSpecInitiators'
import { JobData } from './Show'

const styles = (theme: Theme) =>
  createStyles({
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
    },
    horizontalNavLink: {
      paddingTop: theme.spacing.unit * 4,
      paddingBottom: theme.spacing.unit * 4,
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
    modalContent: {
      width: 'inherit',
    },
    archiveButton: {
      marginTop: theme.spacing.unit * 4,
    },
  })

const CreateRunSuccessNotification = ({ data }: any) => (
  <React.Fragment>
    Successfully created job run{' '}
    <BaseLink href={`/jobs/${data.attributes.jobId}/runs/id/${data.id}`}>
      {data.id}
    </BaseLink>
  </React.Fragment>
)

const DeleteSuccessNotification = ({ id }: any) => (
  <React.Fragment>Successfully archived job {id}</React.Fragment>
)

const DEFAULT_PAGE = 1
const RECENT_RUNS_COUNT = 5

interface Props extends WithStyles<typeof styles> {
  fetchJobRuns: Function
  createJobRun: Function
  deleteJobSpec: Function
  jobSpecId: string
  job: JobData['jobSpec']
  url: string
}

const RegionalNavComponent = ({
  classes,
  createJobRun,
  fetchJobRuns,
  jobSpecId,
  job,
  deleteJobSpec,
}: Props) => {
  const location = useLocation()
  const navErrorsActive = location.pathname.endsWith('/errors')
  const navDefinitionActive = location.pathname.endsWith('/json')
  const navOverviewActive = !navDefinitionActive && !navErrorsActive
  const definition = job && jobSpecDefinition({ ...job, ...job.attributes })
  const [modalOpen, setModalOpen] = React.useState(false)
  const [archived, setArchived] = React.useState(false)
  const errorsTabText =
    job?.attributes.errors && job.attributes.errors.length > 0
      ? `Errors (${job.attributes.errors.length})`
      : 'Errors'
  const handleRun = () => {
    createJobRun(jobSpecId, CreateRunSuccessNotification, ErrorMessage).then(
      () =>
        fetchJobRuns({
          jobSpecId,
          page: DEFAULT_PAGE,
          size: RECENT_RUNS_COUNT,
        }),
    )
  }
  const handleDelete = (id: string) => {
    deleteJobSpec(id, () => DeleteSuccessNotification({ id }), ErrorMessage)
    setArchived(true)
  }
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
                <Grid item className={classes.archiveButton}>
                  <Button
                    variant="danger"
                    onClick={() => handleDelete(jobSpecId)}
                  >
                    Archive {jobSpecId}
                    {archived && <Redirect to="/" />}
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
            <Typography variant="subtitle2" color="secondary" gutterBottom>
              Job Spec Detail
            </Typography>
          </Grid>
          <Grid item xs={12}>
            <Grid
              container
              spacing={0}
              alignItems="center"
              className={classes.mainRow}
            >
              <Grid item xs={6}>
                <Typography
                  variant="h3"
                  color="secondary"
                  className={classes.jobSpecId}
                >
                  {jobSpecId}
                </Typography>
              </Grid>
              <Grid item xs={6} className={classes.actions}>
                <Button
                  className={classes.regionalNavButton}
                  onClick={() => setModalOpen(true)}
                >
                  Archive
                </Button>
                {job && isWebInitiator(job.attributes.initiators) && (
                  <Button
                    onClick={handleRun}
                    className={classes.regionalNavButton}
                  >
                    Run
                  </Button>
                )}
                {definition && (
                  <Button
                    href={{
                      pathname: '/jobs/new',
                      state: { definition },
                    }}
                    component={BaseLink}
                    className={classes.regionalNavButton}
                  >
                    Duplicate
                  </Button>
                )}
                {definition && (
                  <CopyJobSpec
                    JobSpec={definition}
                    className={classes.regionalNavButton}
                  />
                )}
              </Grid>
            </Grid>
          </Grid>
          <Grid item xs={12}>
            <Typography variant="subtitle2" color="textSecondary">
              Created{' '}
              {job && job.attributes.createdAt && (
                <>
                  <TimeAgo tooltip={false}>{job.attributes.createdAt}</TimeAgo>{' '}
                  ({localizedTimestamp(job.attributes.createdAt)})
                </>
              )}
            </Typography>
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
                  href={`/jobs/${jobSpecId}/json`}
                  className={classNames(
                    classes.horizontalNavLink,
                    navDefinitionActive && classes.activeNavLink,
                  )}
                >
                  JSON
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
                  {errorsTabText}
                </Link>
              </ListItem>
            </List>
          </Grid>
        </Grid>
      </Card>
    </>
  )
}

const mapStateToProps = (state: any) => ({
  url: state.notifications.currentUrl,
})

export const ConnectedRegionalNav = connect(mapStateToProps, {
  fetchJobRuns,
  createJobRun,
  deleteJobSpec,
})(RegionalNavComponent)

export const RegionalNav = withStyles(styles)(ConnectedRegionalNav)

export default RegionalNav
