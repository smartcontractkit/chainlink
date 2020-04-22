import React from 'react'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import MuiTable from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import MuiTableCell from '@material-ui/core/TableCell'
import TableCell, { Column } from './Table/TableCell'
import TablePagination from '@material-ui/core/TablePagination'
import Paper from '@material-ui/core/Paper'
import PaginationActions from './Table/PaginationActions'

export const DEFAULT_ROWS_PER_PAGE = 10
export const DEFAULT_CURRENT_PAGE = 1

interface MsgProps {
  colCount: number
  msg?: string
}

const Loading: React.FC<MsgProps> = ({ colCount, msg }) => {
  return (
    <TableRow>
      <MuiTableCell component="th" scope="row" colSpan={colCount}>
        {msg ?? 'Loading...'}
      </MuiTableCell>
    </TableRow>
  )
}

const Empty: React.FC<MsgProps> = ({ colCount, msg }) => {
  return (
    <TableRow>
      <MuiTableCell component="th" scope="row" colSpan={colCount}>
        {msg ?? 'There are no results added to the Explorer yet.'}
      </MuiTableCell>
    </TableRow>
  )
}

const Error: React.FC<MsgProps> = ({ colCount, msg }) => {
  return (
    <TableRow>
      <MuiTableCell component="th" scope="row" colSpan={colCount}>
        {msg ?? 'Error loading resources.'}
      </MuiTableCell>
    </TableRow>
  )
}

const styles = (theme: Theme) =>
  createStyles({
    header: {
      backgroundColor: theme.palette.grey['50'],
    },
    table: {
      minHeight: 150,
      whiteSpace: 'nowrap',
    },
    root: {
      width: '100%',
      overflowX: 'auto',
    },
  })

export type ChangePageEvent = React.MouseEvent<HTMLButtonElement> | null

interface Props extends WithStyles<typeof styles> {
  headers: readonly string[]
  rowsPerPage: number
  currentPage: number
  onChangePage: (event: ChangePageEvent, page: number) => void
  loaded?: boolean
  rows?: Column[][]
  count?: number
  loadingMsg?: string
  emptyMsg?: string
  errorMsg?: string
}

const renderRows = ({
  loaded,
  headers,
  rows,
  loadingMsg,
  emptyMsg,
  errorMsg,
}: Props) => {
  if (loaded && !rows) {
    return <Error colCount={headers.length} msg={errorMsg} />
  } else if (!rows) {
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

const Table = (props: Props) => {
  return (
    <Paper className={props.classes.root}>
      <MuiTable className={props.classes.table}>
        <TableHead>
          <TableRow className={props.classes.header}>
            {props.headers.map((h: string) => (
              <MuiTableCell key={h}>{h}</MuiTableCell>
            ))}
          </TableRow>
        </TableHead>
        <TableBody>{renderRows(props)}</TableBody>
      </MuiTable>
      <TablePagination
        component="div"
        count={props.count ?? 0}
        rowsPerPageOptions={[]}
        rowsPerPage={props.rowsPerPage}
        page={props.currentPage - 1}
        SelectProps={{
          native: true,
        }}
        onChangePage={props.onChangePage}
        ActionsComponent={PaginationActions}
      />
    </Paper>
  )
}

Table.defaultProps = {
  rowsPerPage: DEFAULT_ROWS_PER_PAGE,
  currentPage: DEFAULT_CURRENT_PAGE,
}

export default withStyles(styles)(Table)
