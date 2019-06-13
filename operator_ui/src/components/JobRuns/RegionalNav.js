import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { withStyles } from '@material-ui/core/styles'
import Card from '@material-ui/core/Card'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import List from '@material-ui/core/List'
import ListItem from '@material-ui/core/ListItem'
import { fetchJob, createJobRun } from 'actions'
import Link from 'components/Link'
import TimeAgo from 'components/TimeAgo'
import classNames from 'classnames'

const styles = theme => {
  return {
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
    }
  }
}

const RegionalNav = ({ classes, jobSpecId, jobRunId, jobRun, url }) => {
  const navOverviewActive = url && !url.includes('json')
  const navDefinitionACtive = !navOverviewActive
  return (
    <Card className={classes.container}>
      <Grid container spacing={0}>
        <Grid item xs={12}>
          <Typography variant="subtitle2" color="secondary" gutterBottom>
            Job Run Detail
          </Typography>
          <Link to={`/jobs/${jobSpecId}`} variant="subtitle1" color="primary">
            {jobSpecId}
          </Link>
        </Grid>
        <Grid item xs={12}>
          <Typography variant="h3" color="secondary" gutterBottom>
            {jobRunId}
          </Typography>
        </Grid>
        <Grid item xs={12}>
          <Typography variant="subtitle2" color="textSecondary">
            {jobRun && (
              <React.Fragment>
                Started <TimeAgo>{jobRun.createdAt}</TimeAgo>
              </React.Fragment>
            )}
          </Typography>
        </Grid>
        <Grid item xs={12}>
          <List className={classes.horizontalNav}>
            <ListItem className={classes.horizontalNavItem}>
              <Link
                to={`/jobs/${jobSpecId}/runs/id/${jobRunId}`}
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
                to={`/jobs/${jobSpecId}/runs/id/${jobRunId}/json`}
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

RegionalNav.propTypes = {
  jobSpecId: PropTypes.string.isRequired,
  jobRunId: PropTypes.string.isRequired,
  jobRun: PropTypes.object
}

const mapStateToProps = state => ({
  url: state.notifications.currentUrl
})

export const ConnectedRegionalNav = connect(
  mapStateToProps,
  { fetchJob, createJobRun }
)(RegionalNav)

export default withStyles(styles)(ConnectedRegionalNav)
