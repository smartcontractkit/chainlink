import React from 'react'

import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'

import { Loading } from 'src/components/Feedback/Loading'

export const LoadingRow = ({ visible }: { visible: boolean }) => {
  if (!visible) {
    return null
  }

  return (
    <TableRow>
      <TableCell component="th" scope="row">
        <Loading />
      </TableCell>
    </TableRow>
  )
}

export const NoContentRow = ({ visible }: { visible: boolean }) => {
  if (!visible) {
    return null
  }

  return (
    <TableRow>
      <TableCell component="th" scope="row">
        No entries to show
      </TableCell>
    </TableRow>
  )
}

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
