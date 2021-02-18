import React from 'react'
import { TableCell, TableRow, Typography } from '@material-ui/core'
import { TimeAgo } from 'components/TimeAgo'
import Link from 'components/Link'
import { JobSpecV2 } from './JobsIndex'
import {
  createStyles,
  withStyles,
  WithStyles,
  Theme,
} from '@material-ui/core/styles'

const styles = (theme: Theme) =>
  createStyles({
    cell: {
      paddingTop: theme.spacing.unit * 2,
      paddingBottom: theme.spacing.unit * 2,
    },
    link: {
      '&::before': {
        content: "''",
        position: 'absolute',
        top: 0,
        left: 0,
        width: '100%',
        height: '100%',
      },
    },
  })

interface Props extends WithStyles<typeof styles> {
  job: JobSpecV2
}

const JobV2Row = withStyles(styles)(({ job, classes }: Props) => {
  const createdAt = React.useMemo(() => {
    const { fluxMonitorSpec, offChainReportingOracleSpec } = job.attributes

    switch (job.attributes.type) {
      case 'fluxmonitor':
        return fluxMonitorSpec ? fluxMonitorSpec.createdAt : undefined
      case 'offchainreporting':
        return offChainReportingOracleSpec
          ? offChainReportingOracleSpec.createdAt
          : undefined
      default:
        return undefined
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
        {createdAt && (
          <Typography variant="body1">
            <TimeAgo tooltip>{createdAt}</TimeAgo>
          </Typography>
        )}
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

export default JobV2Row
