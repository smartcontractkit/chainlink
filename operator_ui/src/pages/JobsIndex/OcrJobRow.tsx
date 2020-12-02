import React from 'react'
import { TableCell, TableRow, Typography } from '@material-ui/core'
import { TimeAgo } from '@chainlink/styleguide'
import Link from 'components/Link'
import { OffChainReporting } from './JobsIndex'
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
  job: OffChainReporting
}

export const OcrJobRow = withStyles(styles)(({ job, classes }: Props) => {
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
          <TimeAgo tooltip>
            {job.attributes.offChainReportingOracleSpec.createdAt}
          </TimeAgo>
        </Typography>
      </TableCell>
      <TableCell>
        <Typography variant="body1">Off-chain reporting</Typography>
      </TableCell>
      <TableCell>
        <Typography variant="body1" color="textSecondary">
          n/a
        </Typography>
      </TableCell>
    </TableRow>
  )
})
