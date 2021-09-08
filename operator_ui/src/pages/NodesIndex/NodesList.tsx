import React from 'react'
import {
  Table,
  TableHead,
  TableBody,
  TableCell,
  TableRow,
  Typography,
} from '@material-ui/core'
import { NodeRow } from './NodeRow'
import { NodeSpecV2 } from './NodesIndex'

interface Props {
  nodes: NodeSpecV2[]
  nodeFilter: (node: NodeSpecV2) => boolean
}

const List = ({ nodes, nodeFilter }: Props) => {
  return (
    <Table>
      <TableHead>
        <TableRow>
          <TableCell>
            <Typography variant="body1" color="textSecondary">
              ID
            </Typography>
          </TableCell>

          <TableCell>
            <Typography variant="body1" color="textSecondary">
              Name
            </Typography>
          </TableCell>

          <TableCell>
            <Typography variant="body1" color="textSecondary">
              EVM Chain ID
            </Typography>
          </TableCell>

          <TableCell>
            <Typography variant="body1" color="textSecondary">
              Created
            </Typography>
          </TableCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {nodes.filter(nodeFilter).map((node) => (
          <NodeRow key={node.id} node={node} />
        ))}
      </TableBody>
    </Table>
  )
}

export default List
