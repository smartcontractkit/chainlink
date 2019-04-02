import React from 'react'
import MuiTable from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'

interface ILoadingProps {
  colCount: number
  msg?: string
}

const Loading = ({ colCount, msg }: ILoadingProps) => {
  return (
    <TableRow>
      <TableCell component="th" scope="row" colSpan={colCount}>
        {msg || 'Loading...'}
      </TableCell>
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
      <TableCell component="th" scope="row" colSpan={colCount}>
        {msg || 'No results'}
      </TableCell>
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
        {r.map((c: any, idx: number) => (
          <TableCell key={idx}>{c}</TableCell>
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
            <TableCell key={h}>{h}</TableCell>
          ))}
        </TableRow>
      </TableHead>
      <TableBody>{renderRows(props)}</TableBody>
    </MuiTable>
  )
}

export default Table
