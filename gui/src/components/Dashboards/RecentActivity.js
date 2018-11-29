import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'
import { fade } from '@material-ui/core/styles/colorManipulator'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import CardContent from '@material-ui/core/CardContent'
import Button from '@material-ui/core/Button'
import Card from '@material-ui/core/Card'
import ReactStaticLinkComponent from 'components/ReactStaticLinkComponent'
import Link from 'components/Link'
import TimeAgo from 'components/TimeAgo'
import StatusIcon from 'components/JobRuns/StatusIcon'

const styles = theme => {
  return {
    cell: {
      borderColor: theme.palette.divider,
      borderTop: `1px solid`,
      borderBottom: 'none',
      padding: 0,
      '&:hover': {
        backgroundColor: fade(theme.palette.grey.A100, 0.2)
      }
    },
    status: {
      position: 'absolute',
      top: 0,
      left: 0,
      paddingTop: 35,
      paddingLeft: 35,
      borderRight: 'solid 1px',
      borderRightColor: theme.palette.divider,
      width: 50,
      height: '100%'
    },
    runDetails: {
      padding: theme.spacing.unit * 4
    }
  }
}

const RecentActivity = ({classes, runs}) => {
  const loading = !runs
  let activity

  if (loading) {
    activity = (
      <CardContent>
        <Typography variant='body1' color='textSecondary'>...</Typography>
      </CardContent>
    )
  } else if (runs.length === 0) {
    activity = (
      <CardContent>
        <Typography variant='body1' color='textSecondary'>
          No recently activity :(
        </Typography>
      </CardContent>
    )
  } else {
    activity = (
      <Table>
        <TableBody>
          {runs.map(r => (
            <TableRow key={r.id}>
              <TableCell scope='row' className={classes.cell}>
                <div style={{position: 'relative', paddingLeft: '50px'}}>
                  <div className={classes.status}>
                    <StatusIcon>{r.status}</StatusIcon>
                  </div>
                  <div className={classes.runDetails}>
                    <Grid container>
                      <Grid item xs={12}>
                        <Typography variant='body1' color='textSecondary'>
                          <TimeAgo>{r.createdAt}</TimeAgo>
                        </Typography>
                      </Grid>
                      <Grid item xs={12}>
                        <Link to={`/jobs/${r.jobId}`}>
                          <Typography variant='h6' color='textPrimary' component='span'>
                            @{r.jobId}
                          </Typography>
                        </Link>
                      </Grid>
                      <Grid item xs={12}>
                        <Link to={`/jobs/${r.jobId}/runs/id/${r.id}`}>
                          <Typography variant='subtitle1' color='textSecondary' component='span'>
                            #{r.id}
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
        <Grid container>
          <Grid item xs={12} sm={8}>
            <Typography variant='h5' color='secondary'>
              Activity
            </Typography>
          </Grid>
          <Grid item xs={12} sm={4} align='right'>
            <Button
              variant='outlined'
              color='primary'
              component={ReactStaticLinkComponent}
              to={'/jobs/new'}
            >
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
