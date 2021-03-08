import React from 'react'
import { TableCell, TableRow, Typography } from '@material-ui/core'
import { TimeAgo } from 'components/TimeAgo'
import Link from 'components/Link'
import { JobSpecV2 } from './JobsIndex'
import { withStyles, WithStyles } from '@material-ui/core/styles'
import { styles } from './sharedStyles'

interface Props extends WithStyles<typeof styles> {
  job: JobSpecV2
}

export const JobV2Row = withStyles(styles)(({ job, classes }: Props) => {
  const createdAt = React.useMemo(() => {
    switch (job.attributes.type) {
      case 'fluxmonitor':
        return job.attributes.fluxMonitorSpec.createdAt
      case 'offchainreporting':
        return job.attributes.offChainReportingOracleSpec.createdAt
      case 'directrequest':
        return job.attributes.directRequestSpec.createdAt
    }
  }, [job])

  const type = React.useMemo(() => {
    switch (job.attributes.type) {
      case 'fluxmonitor':
        return 'Direct Request'
      case 'offchainreporting':
        return 'Off-chain reporting'
      default:
        return ''
    }
  }, [job])

  const initiator = React.useMemo(() => {
    switch (job.attributes.type) {
      case 'fluxmonitor':
        return 'fluxmonitor'
      case 'offchainreporting':
        return 'N/A'
      default:
        return ''
    }
  }, [job])

  return (
    <TableRow style={{ transform: 'scale(1)' }} hover>
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
          <TimeAgo tooltip>{createdAt}</TimeAgo>
        </Typography>
      </TableCell>
      <TableCell>
        <Typography variant="body1">{type}</Typography>
      </TableCell>
      <TableCell>
        <Typography
          variant="body1"
          color={
            job.attributes.type === 'offchainreporting'
              ? 'textSecondary'
              : 'textPrimary'
          }
        >
          {initiator}
        </Typography>
      </TableCell>
    </TableRow>
  )
})
