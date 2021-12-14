import React from 'react'

import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'

import { CopyIconButton } from 'src/components/Copy/CopyIconButton'
import { fromJuels } from 'src/utils/tokens/link'
import { shortenHex } from 'src/utils/shortenHex'
import { TimeAgo } from 'src/components/TimeAgo'

interface Props {
  ethKey: EthKeysPayload_ResultsFields
}

export const EVMAccountRow: React.FC<Props> = ({ ethKey }) => {
  return (
    <TableRow hover>
      <TableCell>
        <Typography variant="body1">
          {shortenHex(ethKey.address, { start: 6, end: 6 })}{' '}
          <CopyIconButton data={ethKey.address} />
        </Typography>
      </TableCell>
      <TableCell>
        <Typography variant="body1">{ethKey.chain.id}</Typography>
      </TableCell>
      <TableCell>
        <Typography variant="body1">
          {ethKey.isFunding ? 'Emergency funding' : 'Regular'}
        </Typography>
      </TableCell>
      <TableCell>
        <Typography variant="body1">
          {ethKey.linkBalance && fromJuels(ethKey.linkBalance)}
        </Typography>
      </TableCell>
      <TableCell>
        <Typography variant="body1">{ethKey.ethBalance}</Typography>
      </TableCell>
      <TableCell>
        <Typography variant="body1">
          <TimeAgo tooltip>{ethKey.createdAt}</TimeAgo>
        </Typography>
      </TableCell>
    </TableRow>
  )
}
