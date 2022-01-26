import React from 'react'

import { withStyles, WithStyles } from '@material-ui/core/styles'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'

import { tableStyles } from 'components/Table'
import { TimeAgo } from 'components/TimeAgo'
import Link from 'components/Link'

interface Props extends WithStyles<typeof tableStyles> {
  job: JobsPayload_ResultsFields
}

export const JobRow = withStyles(tableStyles)(({ job, classes }: Props) => {
  const type = React.useMemo(() => {
    const typename = job.spec.__typename as string

    switch (typename) {
      case 'DirectRequestSpec':
        return 'Direct Request'
      case 'FluxMonitorSpec':
        return 'Flux Monitor'
      default:
        return typename.replace(/Spec$/, '')
    }
  }, [job.spec.__typename])

  return (
    <TableRow className={classes.row} hover>
      <TableCell className={classes.cell} component="th" scope="row">
        <Link className={classes.link} href={`/jobs/${job.id}`}>
          {job.id}
        </Link>
      </TableCell>
      <TableCell>{job.name != 'undefined' ? job.name : '--'}</TableCell>
      <TableCell>{type}</TableCell>
      <TableCell>{job.externalJobID}</TableCell>
      <TableCell>
        <TimeAgo tooltip>{job.createdAt}</TimeAgo>
      </TableCell>
    </TableRow>
  )
})
