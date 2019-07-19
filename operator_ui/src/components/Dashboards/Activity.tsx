import React from 'react'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import TableFooter from '@material-ui/core/TableFooter'
import CardContent from '@material-ui/core/CardContent'
import Card from '@material-ui/core/Card'
import TimeAgo from '@chainlink/styleguide/components/TimeAgo'
import Button from '../Button'
import BaseLink from '../BaseLink'
import Link from '../Link'
import StatusIcon from '../JobRuns/StatusIcon'
import NoContentLogo from '../Logos/NoContent'
import { IJobRuns } from '../../../@types/operator_ui'

const styles = ({ palette, spacing }: Theme) =>
  createStyles({
    cell: {
      borderColor: palette.divider,
      borderTop: `1px solid`,
      borderBottom: 'none',
      padding: 0
    },
    footer: {
      borderColor: palette.divider,
      borderTop: `1px solid`,
      paddingTop: spacing.unit * 2,
      paddingBottom: spacing.unit * 2
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
      borderRightColor: palette.divider,
      width: 50,
      height: '100%'
    },
    runDetails: {
      paddingTop: spacing.unit * 3,
      paddingBottom: spacing.unit * 3,
      paddingLeft: spacing.unit * 4,
      paddingRight: spacing.unit * 4
    },
    noActivity: {
      backgroundColor: palette.primary.light,
      padding: spacing.unit * 3
    }
  })

const NoActivity = ({ classes }) => (
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

interface IProps extends WithStyles<typeof styles> {
  pageSize: number
  runs?: IJobRuns[]
  count?: number
}

const Activity = ({ classes, runs, count, pageSize }: IProps) => {
  const loading = !runs
  let activity

  if (loading) {
    activity = (
      <CardContent>
        <Typography variant="body1" color="textSecondary">
          Loading ...
        </Typography>
      </CardContent>
    )
  } else if (runs.length === 0) {
    activity = <NoActivity classes={classes} />
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
                          <TimeAgo tooltip>{r.createdAt}</TimeAgo>
                        </Typography>
                      </Grid>
                      <Grid item xs={12}>
                        <Link href={`/jobs/${r.jobId}`}>
                          <Typography
                            variant="h5"
                            color="primary"
                            component="span"
                          >
                            Job: {r.jobId}
                          </Typography>
                        </Link>
                      </Grid>
                      <Grid item xs={12}>
                        <Link href={`/jobs/${r.jobId}/runs/id/${r.id}`}>
                          <Typography
                            variant="subtitle1"
                            color="textSecondary"
                            component="span"
                          >
                            Run: {r.id}
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
        {count > pageSize && (
          <TableFooter>
            <TableRow>
              <TableCell scope="row" className={classes.footer}>
                <Button href={`/runs`} component={BaseLink}>
                  View More
                </Button>
              </TableCell>
            </TableRow>
          </TableFooter>
        )}
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
          <Grid item container xs={12} sm={4} justify="flex-end">
            <Button component={BaseLink} href={'/jobs/new'}>
              New Job
            </Button>
          </Grid>
        </Grid>
      </CardContent>

      {activity}
    </Card>
  )
}

export default withStyles(styles)(Activity)
