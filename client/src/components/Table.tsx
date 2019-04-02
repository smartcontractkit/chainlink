import React from 'react'
import MuiTable from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import MuiTableCell from '@material-ui/core/TableCell'
import TableCell, { Column } from './Table/TableCell'

interface ILoadingProps {
  colCount: number
  msg?: string
}

const Loading = ({ colCount, msg }: ILoadingProps) => {
  return (
    <TableRow>
      <MuiTableCell component="th" scope="row" colSpan={colCount}>
        {msg || 'Loading...'}
      </MuiTableCell>
    </TableRow>
  )
}

interface IEmptyProps {
  colCount: number
  msg?: string
}

const Empty = ({ colCount, msg }: IEmptyProps) => {
  return (
    <TableRow>
      <MuiTableCell component="th" scope="row" colSpan={colCount}>
        {msg || 'No results'}
      </MuiTableCell>
    </TableRow>
  )
}

interface IProps {
  headers: string[]
  rows?: any[][]
  loadingMsg?: string
  emptyMsg?: string
}

const renderRows = ({ headers, rows, loadingMsg, emptyMsg }: IProps) => {
  if (!rows) {
    return <Loading colCount={headers.length} msg={loadingMsg} />
  } else if (rows.length === 0) {
    return <Empty colCount={headers.length} msg={emptyMsg} />
  } else {
    return rows.map((r: any[], idx: number) => (
      <TableRow key={idx}>
        {r.map((col: Column, idx: number) => (
          <TableCell key={idx} column={col} />
        ))}
      </TableRow>
    ))
  }
}

const Table = (props: IProps) => {
  return (
    <MuiTable>
      <TableHead>
        <TableRow>
          {props.headers.map((h: string) => (
            <MuiTableCell key={h}>{h}</MuiTableCell>
          ))}
        </TableRow>
      </TableHead>
      <TableBody>{renderRows(props)}</TableBody>
    </MuiTable>
  )
}

export default Table
