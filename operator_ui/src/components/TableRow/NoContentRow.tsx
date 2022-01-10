import React from 'react'

import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'

export const NoContentRow: React.FC<{ visible: boolean }> = ({
  children,
  visible,
}) => {
  if (!visible) {
    return null
  }

  return (
    <TableRow>
      <TableCell component="th" scope="row">
        {children ? children : 'No entries to show'}
      </TableCell>
    </TableRow>
  )
}
