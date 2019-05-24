import React from 'react'
import { connect } from 'react-redux'
import { Link as BaseLink } from 'react-router-dom'
import { withStyles } from '@material-ui/core/styles'
import Card from '@material-ui/core/Card'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import List from '@material-ui/core/List'
import ListItem from '@material-ui/core/ListItem'
import classNames from 'classnames'
import TimeAgo from '@chainlink/styleguide/components/TimeAgo'
import localizedTimestamp from '@chainlink/styleguide/utils/localizedTimestamp'
import Button from 'components/Button'
import Link from 'components/Link'
import ReactStaticLinkComponent from 'components/ReactStaticLinkComponent'
import CopyJobSpec from 'components/CopyJobSpec'
import ErrorMessage from 'components/Notifications/DefaultError'
import jobSpecDefinition from 'utils/jobSpecDefinition'
import { isWebInitiator } from 'utils/jobSpecInitiators'
import { fetchJobRuns, createJobRun } from 'actions'

const styles = theme => {
  return {
    container: {
      backgroundColor: theme.palette.common.white,
      padding: theme.spacing.unit * 5,
      paddingBottom: 0
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
    }
  }
}

const SuccessNotification = ({ data }) => (
  <React.Fragment>
    Successfully created job run{' '}
    <BaseLink to={`/jobs/${data.attributes.jobId}/runs/id/${data.id}`}>
      {data.id}
    </BaseLink>
  </React.Fragment>
)

const DEFAULT_PAGE = 1
const RECENT_RUNS_COUNT = 5

const RegionalNav = ({
  classes,
  createJobRun,
  fetchJobRuns,
  jobSpecId,
  job,
  url
}) => {
  const navOverviewActive = url && !url.includes('json')
  const navDefinitionACtive = !navOverviewActive
  const definition = job && jobSpecDefinition(job)
  const handleClick = () => {
    createJobRun(job.id, SuccessNotification, ErrorMessage).then(() =>
      fetchJobRuns({
        jobSpecId: job.id,
        page: DEFAULT_PAGE,
        size: RECENT_RUNS_COUNT
      })
    )
  }

  return (
    <Card className={classes.container}>
      <Grid container spacing={0}>
        <Grid item xs={12}>
          <Typography variant="subtitle2" color="secondary" gutterBottom>
            Job Spec Detail
          </Typography>
        </Grid>
        <Grid item xs={12}>
          <Grid container spacing={0} alignItems="center">
            <Grid item xs={7}>
              <Typography
                variant="h3"
                color="secondary"
                className={classes.jobSpecId}
                gutterBottom
              >
                {jobSpecId}
              </Typography>
            </Grid>
            <Grid item align="right" xs={5}>
              {job && isWebInitiator(job.initiators) && (
                <Button onClick={handleClick}>Run</Button>
              )}
              {definition && (
                <Button
                  to={{
                    pathname: '/jobs/new',
                    state: { definition: definition }
                  }}
                  component={ReactStaticLinkComponent}
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
                to={`/jobs/${jobSpecId}`}
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
                to={`/jobs/${jobSpecId}/json`}
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
  )
}

const mapStateToProps = state => ({
  url: state.notifications.currentUrl
})

export const ConnectedRegionalNav = connect(
  mapStateToProps,
  { fetchJobRuns, createJobRun }
)(RegionalNav)

export default withStyles(styles)(ConnectedRegionalNav)
