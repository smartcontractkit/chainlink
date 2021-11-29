import React from 'react'

import { tableStyles } from 'components/Table'
import { TimeAgo } from 'components/TimeAgo'
import Link from 'components/Link'

import { withStyles, WithStyles } from '@material-ui/core/styles'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'

interface Props extends WithStyles<typeof tableStyles> {
  chain: ChainsPayload_ResultsFields
}

export const ChainRow = withStyles(tableStyles)(({ chain, classes }: Props) => {
  return (
    <TableRow className={classes.row} hover>
      <TableCell className={classes.cell} component="th" scope="row">
        <Link className={classes.link} href={`/chains/${chain.id}`}>
          {chain.id}
        </Link>
      </TableCell>

      <TableCell>{chain.enabled.toString()}</TableCell>

      <TableCell>
        <TimeAgo tooltip>{chain.createdAt}</TimeAgo>
      </TableCell>
    </TableRow>
  )
})
