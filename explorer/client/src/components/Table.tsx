import React from 'react'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import MuiTable from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import MuiTableCell from '@material-ui/core/TableCell'
import TableCell, { Column } from './Table/TableCell'
import TablePagination from '@material-ui/core/TablePagination'
import PaginationActions from './Table/PaginationActions'
import { Hidden, Paper } from '@material-ui/core'

export const DEFAULT_ROWS_PER_PAGE = 10
export const DEFAULT_CURRENT_PAGE = 0

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

const styles = (theme: Theme) =>
  createStyles({
    header: {
      backgroundColor: theme.palette.grey['50']
    },
    table: {
      minHeight: 150,
      whiteSpace: 'nowrap'
    },
    root: {
      width: '100%',
      overflowX: 'auto'
    }
  })

export type ChangePageEvent = React.MouseEvent<HTMLButtonElement> | null

interface IProps extends WithStyles<typeof styles> {
  headers: string[]
  mobileHeaders: string[]
  rowsPerPage: number
  currentPage: number
  onChangePage: (event: ChangePageEvent, page: number) => void
  rows?: any[][]
  count?: number
  loadingMsg?: string
  emptyMsg?: string
}

const renderCells = ({ r, idx }: { r: any[]; idx: number }) => (
  <TableRow key={idx}>
    {r &&
      r.map((col: Column, idx: number) => <TableCell key={idx} column={col} />)}
  </TableRow>
)

const renderRows = ({ headers, rows, loadingMsg, emptyMsg }: IProps) => {
  if (!rows) {
    return <Loading colCount={headers.length} msg={loadingMsg} />
  } else if (rows.length === 0) {
    return <Empty colCount={headers.length} msg={emptyMsg} />
  } else {
    return rows.map((r: any[], idx: number) => (
      <>
        <Hidden xsDown>{renderCells({ r, idx })}</Hidden>
        <Hidden smUp>
          {renderCells({ r: [r[2], r[1], r[3], r[0]], idx: idx })}
        </Hidden>
      </>
    ))
  }
}

const Table = (props: IProps) => {
  return (
    <Paper className={props.classes.root}>
      <MuiTable className={props.classes.table}>
        <TableHead>
          <Hidden xsDown>
            <TableRow className={props.classes.header}>
              {props.headers.map((h: string) => (
                <MuiTableCell key={h}>{h}</MuiTableCell>
              ))}
            </TableRow>
          </Hidden>
          <Hidden smUp>
            <TableRow className={props.classes.header}>
              {props.mobileHeaders &&
                props.mobileHeaders.map((h: string) => (
                  <MuiTableCell key={h}>{h}</MuiTableCell>
                ))}
            </TableRow>
          </Hidden>
        </TableHead>
        <TableBody>{renderRows(props)}</TableBody>
      </MuiTable>
      <TablePagination
        component="div"
        count={props.count || 0}
        rowsPerPageOptions={[]}
        rowsPerPage={props.rowsPerPage}
        page={props.currentPage}
        SelectProps={{
          native: true
        }}
        onChangePage={props.onChangePage}
        ActionsComponent={PaginationActions}
      />
    </Paper>
  )
}

Table.defaultProps = {
  rowsPerPage: DEFAULT_ROWS_PER_PAGE,
  currentPage: DEFAULT_CURRENT_PAGE
}

export default withStyles(styles)(Table)
