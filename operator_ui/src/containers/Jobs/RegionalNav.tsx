import React from 'react'
import { Redirect } from 'react-router-dom'
import { connect } from 'react-redux'
import {
  createStyles,
  withStyles,
  WithStyles,
  Theme
} from '@material-ui/core/styles'
import Card from '@material-ui/core/Card'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import List from '@material-ui/core/List'
import ListItem from '@material-ui/core/ListItem'
import Dialog from '@material-ui/core/Dialog'
import TimeAgo from '@chainlink/styleguide/src/components/TimeAgo'
import localizedTimestamp from '@chainlink/styleguide/src/utils/localizedTimestamp'
import Button from '../../components/Button'
import BaseLink from '../../components/BaseLink'
import Link from '../../components/Link'
import CopyJobSpec from '../../components/CopyJobSpec'
import ErrorMessage from '../../components/Notifications/DefaultError'
import jobSpecDefinition from '../../utils/jobSpecDefinition'
import { isWebInitiator } from '../../utils/jobSpecInitiators'
import { fetchJobRuns, createJobRun } from '../../actions'
import classNames from 'classnames'
import { deleteJobSpec } from '../../actions'
import Close from '../../components/Icons/Close'
import { JobSpec } from '../../../@types/operator_ui'

const styles = (theme: Theme) =>
  createStyles({
    container: {
      backgroundColor: theme.palette.common.white,
      padding: theme.spacing.unit * 5,
      paddingBottom: 0
    },
    mainRow: {
      marginBottom: theme.spacing.unit * 2
    },
    actions: {
      textAlign: 'right'
    },
    duplicate: {
      marginLeft: theme.spacing.unit,
      marginRight: theme.spacing.unit
    },
    horizontalNav: {
      paddingBottom: 0
    },
    horizontalNavItem: {
      display: 'inline'
    },
    horizontalNavLink: {
      paddingTop: theme.spacing.unit * 4,
      paddingBottom: theme.spacing.unit * 4,
      textDecoration: 'none',
      display: 'inline-block',
      borderBottom: 'solid 1px',
      borderBottomColor: theme.palette.common.white,
      '&:hover': {
        borderBottomColor: theme.palette.primary.main
      }
    },
    activeNavLink: {
      color: theme.palette.primary.main,
      borderBottomColor: theme.palette.primary.main
    },
    jobSpecId: {
      overflow: 'hidden',
      textOverflow: 'ellipsis'
    },
    dialogPaper: {
      minHeight: '240px',
      maxHeight: '240px',
      minWidth: '670px',
      maxWidth: '670px',
      overflow: 'hidden',
      borderRadius: theme.spacing.unit * 3
    },
    warningText: {
      fontWeight: 500,
      marginLeft: theme.spacing.unit * 3,
      marginTop: theme.spacing.unit * 3,
      marginBottom: theme.spacing.unit
    },
    closeButton: {
      marginRight: theme.spacing.unit * 3,
      marginTop: theme.spacing.unit * 3
    },
    infoText: {
      fontSize: theme.spacing.unit * 2,
      fontWeight: 450,
      marginLeft: theme.spacing.unit * 6
    },
    modalContent: {
      width: 'inherit'
    },
    archiveButton: {
      marginTop: theme.spacing.unit * 4
    }
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

interface IProps extends WithStyles<typeof styles> {
  fetchJobRuns: Function
  createJobRun: Function
  deleteJobSpec: Function
  jobSpecId: string
  job: JobSpec
  url: String
}

const RegionalNav = ({
  classes,
  createJobRun,
  fetchJobRuns,
  jobSpecId,
  job,
  deleteJobSpec,
  url
}: IProps) => {
  const navOverviewActive = url && !url.includes('json')
  const navDefinitionACtive = !navOverviewActive
  const definition = job && jobSpecDefinition(job)
  const [modalOpen, setModalOpen] = React.useState(false)
  const [archived, setArchived] = React.useState(false)
  const handleRun = () => {
    createJobRun(job.id, CreateRunSuccessNotification, ErrorMessage).then(() =>
      fetchJobRuns({
        jobSpecId: job.id,
        page: DEFAULT_PAGE,
        size: RECENT_RUNS_COUNT
      })
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
              <Grid item xs={7}>
                <Typography
                  variant="h3"
                  color="secondary"
                  className={classes.jobSpecId}
                >
                  {jobSpecId}
                </Typography>
              </Grid>
              <Grid item xs={5} className={classes.actions}>
                <Button
                  className={classes.duplicate}
                  onClick={() => setModalOpen(true)}
                >
                  Archive
                </Button>
                {job && isWebInitiator(job.initiators) && (
                  <Button onClick={handleRun}>Run</Button>
                )}
                {definition && (
                  <Button
                    href={{
                      pathname: '/jobs/new',
                      state: { definition: definition }
                    }}
                    component={BaseLink}
                    className={classes.duplicate}
                  >
                    Duplicate
                  </Button>
                )}
                {definition && <CopyJobSpec JobSpec={definition} />}
              </Grid>
            </Grid>
          </Grid>
          <Grid item xs={12}>
            <Typography variant="subtitle2" color="textSecondary">
              {job && (
                <React.Fragment>
                  Created <TimeAgo tooltip={false}>{job.createdAt}</TimeAgo> (
                  {localizedTimestamp(job.createdAt)})
                </React.Fragment>
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
                    navOverviewActive && classes.activeNavLink
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
                    navDefinitionACtive && classes.activeNavLink
                  )}
                >
                  JSON
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
  url: state.notifications.currentUrl
})

export const ConnectedRegionalNav = connect(
  mapStateToProps,
  { fetchJobRuns, createJobRun, deleteJobSpec }
)(RegionalNav)

export default withStyles(styles)(ConnectedRegionalNav)
