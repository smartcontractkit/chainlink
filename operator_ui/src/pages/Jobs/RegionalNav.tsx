import { localizedTimestamp, TimeAgo } from '@chainlink/styleguide'
import Card from '@material-ui/core/Card'
import Dialog from '@material-ui/core/Dialog'
import Grid from '@material-ui/core/Grid'
import List from '@material-ui/core/List'
import ListItem from '@material-ui/core/ListItem'
import Badge from '@material-ui/core/Badge'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { ApiResponse } from '@chainlink/json-api-client'
import { JobSpec } from 'core/store/models'
import classNames from 'classnames'
import React from 'react'
import { connect } from 'react-redux'
import { Redirect, useLocation } from 'react-router-dom'
import { createJobRun, deleteJobSpec } from 'actionCreators'
import BaseLink from 'components/BaseLink'
import Button from 'components/Button'
import CopyJobSpec from 'components/CopyJobSpec'
import Close from 'components/Icons/Close'
import Link from 'components/Link'
import ErrorMessage from 'components/Notifications/DefaultError'
import { JobData } from './sharedTypes'

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
    modalContent: {
      width: 'inherit',
    },
    archiveButton: {
      marginTop: theme.spacing.unit * 4,
    },
  })

const isWebInitiator = (
  initiators: ApiResponse<JobSpec>['data']['attributes']['initiators'],
) => initiators.find((initiator) => initiator.type === 'web')

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

interface Props extends WithStyles<typeof styles> {
  createJobRun: Function
  deleteJobSpec: Function
  jobSpecId: string
  job: JobData['job']
  runsCount: JobData['recentRunsCount']
  url: string
  getJobSpecRuns: (props: { page?: number; size?: number }) => Promise<void>
}

const RegionalNavComponent = ({
  classes,
  createJobRun,
  jobSpecId,
  job,
  deleteJobSpec,
  getJobSpecRuns,
  runsCount,
}: Props) => {
  const location = useLocation()
  const navErrorsActive = location.pathname.endsWith('/errors')
  const navDefinitionActive = location.pathname.endsWith('/definition')
  const navRunsActive = location.pathname.endsWith('/runs')
  const navOverviewActive =
    !navDefinitionActive && !navErrorsActive && !navRunsActive
  const [modalOpen, setModalOpen] = React.useState(false)
  const [archived, setArchived] = React.useState(false)

  const handleRun = () => {
    const params = new URLSearchParams(location.search)
    const page = params.get('page')
    const size = params.get('size')

    createJobRun(jobSpecId, CreateRunSuccessNotification, ErrorMessage).then(
      () =>
        getJobSpecRuns({
          page: page ? parseInt(page, 10) : undefined,
          size: size ? parseInt(size, 10) : undefined,
        }),
    )
  }
  const handleDelete = (id: string) => {
    deleteJobSpec(
      id,
      () => DeleteSuccessNotification({ id }),
      ErrorMessage,
      job?.type,
    )
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
            {job && (
              <Typography variant="subtitle2" color="secondary" gutterBottom>
                {job?.type} job spec detail
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
                      Archive
                    </Button>
                    {job.type === 'Direct request' &&
                      job.initiators &&
                      isWebInitiator(job.initiators) && (
                        <Button
                          onClick={handleRun}
                          className={classes.regionalNavButton}
                        >
                          Run
                        </Button>
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

const mapStateToProps = (state: any) => ({
  url: state.notifications.currentUrl,
})

export const ConnectedRegionalNav = connect(mapStateToProps, {
  createJobRun,
  deleteJobSpec,
})(RegionalNavComponent)

export const RegionalNav = withStyles(styles)(ConnectedRegionalNav)

export default RegionalNav
