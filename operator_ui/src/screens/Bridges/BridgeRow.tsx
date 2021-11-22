import React from 'react'

import { withStyles, WithStyles } from '@material-ui/core/styles'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'

import Link from 'components/Link'
import { tableStyles } from 'components/Table'

interface Props extends WithStyles<typeof tableStyles> {
  bridge: BridgesPayload_ResultsFields
}

export const BridgeRow = withStyles(tableStyles)(
  ({ bridge, classes }: Props) => {
    return (
      <TableRow className={classes.row} key={bridge.name} hover>
        <TableCell scope="row" component="th">
          <Link className={classes.link} href={`/bridges/${bridge.name}`}>
            {bridge.name}
          </Link>
        </TableCell>
        <TableCell>
          <Typography variant="body1">{bridge.url}</Typography>
        </TableCell>
        <TableCell>
          <Typography variant="body1">{bridge.confirmations}</Typography>
        </TableCell>
        <TableCell>
          <Typography variant="body1">
            {bridge.minimumContractPayment}
          </Typography>
        </TableCell>
      </TableRow>
    )
  },
)
