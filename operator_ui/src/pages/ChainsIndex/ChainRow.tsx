import React from 'react'

import { ChainSpecV2 } from './ChainsIndex'
import { tableStyles } from 'components/Table'
import { TimeAgo } from 'components/TimeAgo'

import { withStyles, WithStyles } from '@material-ui/core/styles'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'

interface Props extends WithStyles<typeof tableStyles> {
  chain: ChainSpecV2
}

export const ChainRow = withStyles(tableStyles)(({ chain, classes }: Props) => {
  const createdAt = chain.attributes.createdAt

  const configOverrides = Object.fromEntries(
    Object.entries(chain.attributes.config).filter(
      ([_key, value]) => value !== null,
    ),
  )

  return (
    <TableRow className={classes.row} hover>
      <TableCell className={classes.cell} component="th" scope="row">
        {chain.id}
      </TableCell>

      <TableCell>
        <pre>{JSON.stringify(configOverrides, null, 2)}</pre>
      </TableCell>

      <TableCell>
        <TimeAgo tooltip>{createdAt}</TimeAgo>
      </TableCell>
    </TableRow>
  )
})
