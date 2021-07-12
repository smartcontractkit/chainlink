import React from 'react'
import { TableCell, TableRow, Typography } from '@material-ui/core'
import { TimeAgo } from 'components/TimeAgo'
import Link from 'components/Link'
import { formatInitiators } from 'utils/jobSpecInitiators'
import { DirectRequest } from './JobsIndex'
import { withStyles, WithStyles } from '@material-ui/core/styles'
import { tableStyles } from 'components/Table'

interface Props extends WithStyles<typeof tableStyles> {
  job: DirectRequest
}

export const DirectRequestRow = withStyles(tableStyles)(
  ({ job, classes }: Props) => {
    return (
      <TableRow className={classes.row} hover>
        <TableCell className={classes.cell} component="th" scope="row">
          <Link className={classes.link} href={`/jobs/${job.id}`}>
            {job.attributes.name || job.id}
            {job.attributes.name && (
              <>
                <br />
                <Typography
                  variant="subtitle2"
                  color="textSecondary"
                  component="span"
                >
                  {job.id}
                </Typography>
              </>
            )}
          </Link>
        </TableCell>
        <TableCell>
          <Typography variant="body1">
            <TimeAgo tooltip>{job.attributes.createdAt}</TimeAgo>
          </Typography>
        </TableCell>
        <TableCell>
          <Typography variant="body1">Direct request</Typography>
        </TableCell>
        <TableCell>
          <Typography variant="body1">
            {formatInitiators(job.attributes.initiators)}
          </Typography>
        </TableCell>
      </TableRow>
    )
  },
)
