import React from 'react'

import { JobResource } from './JobsIndex'
import { tableStyles } from 'components/Table'
import { TimeAgo } from 'components/TimeAgo'
import Link from 'components/Link'

import { withStyles, WithStyles } from '@material-ui/core/styles'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'

interface Props extends WithStyles<typeof tableStyles> {
  job: JobResource
}

export const JobRow = withStyles(tableStyles)(({ job, classes }: Props) => {
  const createdAt = React.useMemo(() => {
    switch (job.attributes.type) {
      case 'directrequest':
        return job.attributes.directRequestSpec.createdAt
      case 'fluxmonitor':
        return job.attributes.fluxMonitorSpec.createdAt
      case 'offchainreporting':
        return job.attributes.offChainReportingOracleSpec.createdAt
      case 'keeper':
        return job.attributes.keeperSpec.createdAt
      case 'cron':
        return job.attributes.cronSpec.createdAt
      case 'webhook':
        return job.attributes.webhookSpec.createdAt
      case 'vrf':
        return job.attributes.vrfSpec.createdAt
    }
  }, [job])

  return (
    <TableRow className={classes.row} hover>
      <TableCell className={classes.cell} component="th" scope="row">
        <Link className={classes.link} href={`/jobs/${job.id}`}>
          {job.id}
        </Link>
      </TableCell>
      <TableCell>{job.attributes.name ? job.attributes.name : '-'}</TableCell>

      <TableCell>{job.attributes.type}</TableCell>
      <TableCell>
        <TimeAgo tooltip>{createdAt}</TimeAgo>
      </TableCell>
    </TableRow>
  )
})
