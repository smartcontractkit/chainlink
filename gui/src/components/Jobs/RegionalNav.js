import React from 'react'
import { connect } from 'react-redux'
import { Link as BaseLink } from 'react-static'
import { withStyles } from '@material-ui/core/styles'
import Card from '@material-ui/core/Card'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import List from '@material-ui/core/List'
import ListItem from '@material-ui/core/ListItem'
import Button from '@material-ui/core/Button'
import Link from 'components/Link'
import ReactStaticLinkComponent from 'components/ReactStaticLinkComponent'
import CopyJobSpec from 'components/CopyJobSpec'
import ErrorMessage from 'components/Notifications/DefaultError'
import TimeAgo from 'components/TimeAgo'
import jobSpecDefinition from 'utils/jobSpecDefinition'
import { isWebInitiator } from 'utils/jobSpecInitiators'
import { fetchJob, createJobRun } from 'actions'

const styles = theme => {
  return ({
    container: {
      backgroundColor: theme.palette.common.white,
      padding: theme.spacing.unit * 5,
      paddingBottom: 0
    },
    duplicate: {
      margin: theme.spacing.unit
    },
    horizontalNav: {
      paddingBottom: 0
    },
    horizontalNavItem: {
      display: 'inline'
    },
    horizontalNavLink: {
      color: theme.palette.secondary.main,
      paddingTop: theme.spacing.unit * 4,
      paddingBottom: theme.spacing.unit * 4,
      textDecoration: 'none',
      display: 'inline-block',
      borderBottom: 'solid 1px',
      borderBottomColor: theme.palette.common.white,
      '&:hover': {
        color: theme.palette.primary.main,
        borderBottomColor: theme.palette.primary.main
      }
    }
  })
}

const SuccessNotification = ({ data }) => (
  <React.Fragment>
    Successfully created job run <BaseLink to={`/jobs/${data.attributes.jobId}/runs/id/${data.id}`}>{data.id}</BaseLink>
  </React.Fragment>
)

const RegionalNav = ({ classes, createJobRun, fetchJob, jobSpecId, job }) => {
  const definition = job && jobSpecDefinition(job)
  const handleClick = () => {
    createJobRun(job.id, SuccessNotification, ErrorMessage)
      .then(() => fetchJob(job.id))
  }

  return (
    <Card className={classes.container}>
      <Grid container spacing={0}>
        <Grid item xs={12}>
          <Typography variant='subtitle2' color='secondary' gutterBottom>
            Job Spec Detail
          </Typography>
        </Grid>
        <Grid item xs={12}>
          <Grid container spacing={0} alignItems='center'>
            <Grid item xs={7}>
              <Typography variant='h3' color='secondary' gutterBottom>
                @{jobSpecId}
              </Typography>
            </Grid>
            <Grid item align='right' xs={5}>
              {job && isWebInitiator(job.initiators) && (
                <Button variant='outlined' color='primary' onClick={handleClick}>
                  Run
                </Button>
              )}
              {definition &&
                <Button
                  to={{ pathname: '/jobs/new', state: { definition: definition } }}
                  component={ReactStaticLinkComponent}
                  color='primary'
                  className={classes.duplicate}
                  variant='outlined'>
                  Duplicate
                </Button>
              }
              {definition &&
                <CopyJobSpec JobSpec={definition} />
              }
            </Grid>
          </Grid>
        </Grid>
        <Grid item xs={12}>
          <Typography variant='subtitle2' color='textSecondary'>
            {job && <React.Fragment>Created <TimeAgo>{job.createdAt}</TimeAgo></React.Fragment>}
          </Typography>
        </Grid>
        <Grid item xs={12}>
          <List className={classes.horizontalNav}>
            <ListItem className={classes.horizontalNavItem}>
              <Link to={`/jobs/${jobSpecId}`} className={classes.horizontalNavLink}>
                Overview
              </Link>
            </ListItem>
            <ListItem className={classes.horizontalNavItem}>
              <Link to={`/jobs/${jobSpecId}/definition`} className={classes.horizontalNavLink}>
                JSON
              </Link>
            </ListItem>
          </List>
        </Grid>
      </Grid>
    </Card>
  )
}

export const ConnectedRegionalNav = connect(
  null,
  { fetchJob, createJobRun }
)(RegionalNav)

export default withStyles(styles)(ConnectedRegionalNav)
