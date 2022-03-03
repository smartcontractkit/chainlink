import React from 'react'

import Grid from '@material-ui/core/Grid'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'

import Link from 'src/components/Link'
import { TimeAgo } from 'src/components/TimeAgo'

const styles = (theme: Theme) =>
  createStyles({
    cell: {
      borderColor: theme.palette.divider,
      borderTop: `1px solid`,
      borderBottom: 'none',
      paddingTop: theme.spacing.unit * 2,
      paddingBottom: theme.spacing.unit * 2,
      paddingLeft: theme.spacing.unit * 2,
    },
    block: { display: 'block' },
    overflowEllipsis: { textOverflow: 'ellipsis', overflow: 'hidden' },
  })

interface Props extends WithStyles<typeof styles> {
  job: RecentJobsPayload_ResultsFields
}

export const RecentJobRow = withStyles(styles)(({ classes, job }: Props) => {
  return (
    <TableRow>
      <TableCell scope="row" className={classes.cell}>
        <Grid container spacing={0}>
          <Grid item xs={12}>
            <Link
              href={`/jobs/${job.id}`}
              classes={{ linkContent: classes.block }}
            >
              <Typography
                className={classes.overflowEllipsis}
                variant="body1"
                component="span"
                color="primary"
              >
                {job.name || job.id}
              </Typography>
            </Link>
          </Grid>
          <Grid item xs={12}>
            <Typography variant="body1" color="textSecondary">
              Created <TimeAgo tooltip>{job.createdAt}</TimeAgo>
            </Typography>
          </Grid>
        </Grid>
      </TableCell>
    </TableRow>
  )
})
