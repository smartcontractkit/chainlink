import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import Button from 'components/Button'
import CardContent from '@material-ui/core/CardContent'
import Card from '@material-ui/core/Card'
import ReactStaticLinkComponent from 'components/ReactStaticLinkComponent'
import Link from 'components/Link'
import TimeAgo from 'components/TimeAgo'
import StatusIcon from 'components/JobRuns/StatusIcon'
import NoContentLogo from 'components/Logos/NoContent'

const styles = theme => {
  return {
    cell: {
      borderColor: theme.palette.divider,
      borderTop: `1px solid`,
      borderBottom: 'none',
      padding: 0
    },
    content: {
      position: 'relative',
      paddingLeft: 50
    },
    status: {
      position: 'absolute',
      top: 0,
      left: 0,
      paddingTop: 18,
      paddingLeft: 30,
      borderRight: 'solid 1px',
      borderRightColor: theme.palette.divider,
      width: 50,
      height: '100%'
    },
    runDetails: {
      paddingTop: theme.spacing(3),
      paddingBottom: theme.spacing(3),
      paddingLeft: theme.spacing(4),
      paddingRight: theme.spacing(4)
    },
    noActivity: {
      backgroundColor: theme.palette.primary.light,
      padding: theme.spacing(3)
    }
  }
}

const NoRecentActivity = ({ classes }) => (
  <CardContent>
    <Card elevation={0} className={classes.noActivity}>
      <Grid container alignItems="center" spacing={16}>
        <Grid item>
          <NoContentLogo width={40} />
        </Grid>
        <Grid item>
          <Typography variant="body1" color="textPrimary" inline>
            No recent activity
          </Typography>
        </Grid>
      </Grid>
    </Card>
  </CardContent>
)

const RecentActivity = ({ classes, runs }) => {
  const loading = !runs
  let activity

  if (loading) {
    activity = (
      <CardContent>
        <Typography variant="body1" color="textSecondary">
          ...
        </Typography>
      </CardContent>
    )
  } else if (runs.length === 0) {
    activity = <NoRecentActivity classes={classes} />
  } else {
    activity = (
      <Table>
        <TableBody>
          {runs.map(r => (
            <TableRow key={r.id}>
              <TableCell scope="row" className={classes.cell}>
                <div className={classes.content}>
                  <div className={classes.status}>
                    <StatusIcon width={38}>{r.status}</StatusIcon>
                  </div>
                  <div className={classes.runDetails}>
                    <Grid container spacing={0}>
                      <Grid item xs={12}>
                        <Typography variant="body1" color="textSecondary">
                          <TimeAgo>{r.createdAt}</TimeAgo>
                        </Typography>
                      </Grid>
                      <Grid item xs={12}>
                        <Link to={`/jobs/${r.jobId}`}>
                          <Typography
                            variant="h5"
                            color="primary"
                            component="span"
                          >
                            {r.jobId}
                          </Typography>
                        </Link>
                      </Grid>
                      <Grid item xs={12}>
                        <Link to={`/jobs/${r.jobId}/runs/id/${r.id}`}>
                          <Typography
                            variant="subtitle1"
                            color="textSecondary"
                            component="span"
                          >
                            {r.id}
                          </Typography>
                        </Link>
                      </Grid>
                    </Grid>
                  </div>
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    )
  }

  return (
    <Card>
      <CardContent>
        <Grid container spacing={0}>
          <Grid item xs={12} sm={8}>
            <Typography variant="h5" color="secondary">
              Activity
            </Typography>
          </Grid>
          <Grid item xs={12} sm={4} align="right">
            <Button component={ReactStaticLinkComponent} to={'/jobs/new'}>
              New Job
            </Button>
          </Grid>
        </Grid>
      </CardContent>

      {activity}
    </Card>
  )
}

RecentActivity.propTypes = {
  runs: PropTypes.array
}

export default withStyles(styles)(RecentActivity)
