import React from 'react'

import { tableStyles } from 'components/Table'
import { TimeAgo } from 'components/TimeAgo'
import Link from 'components/Link'

import { withStyles, WithStyles } from '@material-ui/core/styles'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'

interface Props extends WithStyles<typeof tableStyles> {
  node: NodesPayload_ResultsFields
}

export const NodeRow = withStyles(tableStyles)(({ node, classes }: Props) => {
  return (
    <TableRow className={classes.row} hover>
      <TableCell className={classes.cell} component="th" scope="row">
        <Link className={classes.link} href={`/nodes/${node.id}`}>
          {node.id}
        </Link>
      </TableCell>

      <TableCell>{node.name}</TableCell>
      <TableCell>{node.chain.id}</TableCell>
      <TableCell>
        <TimeAgo tooltip>{node.createdAt}</TimeAgo>
      </TableCell>
      <TableCell>{node.state}</TableCell>
    </TableRow>
  )
})
