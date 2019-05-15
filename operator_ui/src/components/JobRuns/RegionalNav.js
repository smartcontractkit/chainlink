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
  }
}

const RegionalNav = ({ classes, jobSpecId, jobRunId }) => {
  return (
    <Card className={classes.container}>
      <Grid container spacing={0}>
        <Grid item xs={12}>
          <Link to={`/jobs/${jobSpecId}`}>
            <Typography variant="subtitle1" color="primary">
              @{jobSpecId}
            </Typography>
          </Link>
          <Typography variant="h3" color="secondary">
            Job Run Detail
          </Typography>
        </Grid>
        <Grid item xs={12}>
          <Typography variant="subtitle1" color="textSecondary">
            #{jobRunId}
          </Typography>
        </Grid>
        <Grid item xs={12}>
          <List className={classes.horizontalNav}>
            <ListItem className={classes.horizontalNavItem}>
              <Link
                to={`/jobs/${jobSpecId}/runs/id/${jobRunId}`}
                className={classes.horizontalNavLink}
              >
                Overview
              </Link>
            </ListItem>
            <ListItem className={classes.horizontalNavItem}>
              <Link
                to={`/jobs/${jobSpecId}/runs/id/${jobRunId}/json`}
                className={classes.horizontalNavLink}
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

export const ConnectedRegionalNav = connect(
  null,
  { fetchJob, createJobRun }
)(RegionalNav)

export default withStyles(styles)(ConnectedRegionalNav)
