import React from 'react'

import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'

export const ErrorRow = ({ msg }: { msg?: string }) => {
  if (!msg) {
    return null
  }

  return (
    <TableRow>
      <TableCell component="th" scope="row">
        {msg}
      </TableCell>
    </TableRow>
  )
}
