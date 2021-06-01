import React from 'react'
import { TableCell, TableRow, Typography } from '@material-ui/core'
import { TimeAgo } from 'components/TimeAgo'
import Link from 'components/Link'
import { JobSpecV2 } from './JobsIndex'
import { withStyles, WithStyles } from '@material-ui/core/styles'
import { tableStyles } from 'components/Table'

interface Props extends WithStyles<typeof tableStyles> {
  job: JobSpecV2
}

export const JobV2Row = withStyles(tableStyles)(({ job, classes }: Props) => {
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

  const type = React.useMemo(() => {
    switch (job.attributes.type) {
      case 'directrequest':
      case 'fluxmonitor':
        return 'Direct Request'
      case 'offchainreporting':
        return 'Off-chain reporting'
      case 'keeper':
        return 'Keeper'
      case 'cron':
        return 'Cron'
      case 'webhook':
        return 'Webhook'
      case 'vrf':
        return 'VRF'
      default:
        return ''
    }
  }, [job])

  const initiator = React.useMemo(() => {
    switch (job.attributes.type) {
      case 'fluxmonitor':
        return 'fluxmonitor'
      case 'directrequest':
        return job.attributes.directRequestSpec.initiator
      case 'vrf':
      case 'keeper':
      case 'cron':
      case 'webhook':
      case 'offchainreporting':
        return 'N/A'
      default:
        return ''
    }
  }, [job])

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
