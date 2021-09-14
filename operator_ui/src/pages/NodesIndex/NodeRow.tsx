import React from 'react'

import { NodeSpecV2 } from './NodesIndex'
import { tableStyles } from 'components/Table'
import { TimeAgo } from 'components/TimeAgo'

import { withStyles, WithStyles } from '@material-ui/core/styles'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'

interface Props extends WithStyles<typeof tableStyles> {
  node: NodeSpecV2
}

export const NodeRow = withStyles(tableStyles)(({ node, classes }: Props) => {
  const createdAt = node.attributes.createdAt

  return (
    <TableRow className={classes.row} hover>
      <TableCell className={classes.cell} component="th" scope="row">
        {node.id}
      </TableCell>

      <TableCell>{node.attributes.name}</TableCell>

      <TableCell>{node.attributes.evmChainID}</TableCell>

      <TableCell>
        <TimeAgo tooltip>{createdAt}</TimeAgo>
      </TableCell>
    </TableRow>
  )
})
