import { TimeAgo } from 'components/TimeAgo'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Grid from '@material-ui/core/Grid'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableFooter from '@material-ui/core/TableFooter'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'
import { JobRun, JobRuns } from 'operator_ui'
import React from 'react'
import BaseLink from '../BaseLink'
import Button from '../Button'
import StatusIcon from 'components/StatusIcon'
import Link from '../Link'
import NoContentLogo from '../Logos/NoContent'

const noActivityStyles = ({ palette, spacing }: Theme) =>
  createStyles({
    noActivity: {
      backgroundColor: palette.primary.light,
      padding: spacing.unit * 3,
    },
  })

type NoActivityProps = WithStyles<typeof noActivityStyles>

const NoActivity = withStyles(noActivityStyles)(
  ({ classes }: NoActivityProps) => (
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
  ),
)

const Fetching = () => {
  return (
    <CardContent>
      <Typography variant="body1" color="textSecondary">
        Loading ...
      </Typography>
    </CardContent>
  )
}

const styles = ({ palette, spacing }: Theme) =>
  createStyles({
    cell: {
      borderColor: palette.divider,
      borderTop: `1px solid`,
      borderBottom: 'none',
      padding: 0,
    },
    footer: {
      borderColor: palette.divider,
      borderTop: `1px solid`,
      paddingTop: spacing.unit * 2,
      paddingBottom: spacing.unit * 2,
    },
    content: {
      position: 'relative',
      paddingLeft: 50,
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
      height: '100%',
    },
    runDetails: {
      paddingTop: spacing.unit * 3,
      paddingBottom: spacing.unit * 3,
      paddingLeft: spacing.unit * 4,
      paddingRight: spacing.unit * 4,
    },
  })

interface Props extends WithStyles<typeof styles> {
  pageSize: number
  runs?: JobRuns
  count?: number
}

const Activity = ({ classes, runs, count, pageSize }: Props) => {
  let activity

  if (!runs) {
    activity = <Fetching />
  } else if (runs.length === 0) {
    activity = <NoActivity />
  } else {
    activity = (
      <Table>
        <TableBody>
          {runs.map((r: JobRun) => (
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
                        <Link href={`/jobs/${r.jobId}/runs/${r.id}`}>
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
        {count && count > pageSize && (
          <TableFooter>
            <TableRow>
              <TableCell scope="row" className={classes.footer}>
                <Button href={'/runs'} component={BaseLink}>
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
          <Grid item xs={12} sm={4}>
            <Grid container spacing={0} justify="flex-end">
              <Grid item>
                <Button href={'/jobs/new'} component={BaseLink}>
                  New Job
                </Button>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </CardContent>

      {activity}
    </Card>
  )
}

export default withStyles(styles)(Activity)
