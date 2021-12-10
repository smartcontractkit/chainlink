import React from 'react'

import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'

import { CopyIconButton } from 'src/components/Copy/CopyIconButton'

interface Props {
  csaKey: CsaKeysPayload_ResultsFields
}

export const CSAKeyRow: React.FC<Props> = ({ csaKey }) => {
  return (
    <TableRow hover>
      <TableCell>
        <Typography variant="body1">
          {csaKey.publicKey} <CopyIconButton data={csaKey.publicKey} />
        </Typography>
      </TableCell>
    </TableRow>
  )
}
