import React, { useRef, useEffect, useState } from 'react'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import { CSSProperties } from '@material-ui/core/styles/withStyles'
import RootRef from '@material-ui/core/RootRef'
import MuiTable from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import MuiTableCell from '@material-ui/core/TableCell'
import TableCell, { Column } from './Table/TableCell'
import TablePagination, {
  TablePaginationProps,
} from '@material-ui/core/TablePagination'
import Paper from '@material-ui/core/Paper'
import PaginationActions from './Table/PaginationActions'

export const DEFAULT_ROWS_PER_PAGE = 10
export const DEFAULT_CURRENT_PAGE = 1

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

export interface Props extends WithStyles<typeof styles> {
  headers: readonly string[]
  rowsPerPage?: number
  currentPage?: number
  onChangePage: TablePaginationProps['onChangePage']
  loading: boolean
  error: boolean
  rows?: Column[][]
  count?: number
  loadingMsg?: string
  emptyMsg?: string
  errorMsg?: string
}

const Table: React.FC<Props> = ({
  headers,
  rowsPerPage = DEFAULT_ROWS_PER_PAGE,
  currentPage = DEFAULT_CURRENT_PAGE,
  onChangePage,
  loading,
  error,
  rows,
  count,
  loadingMsg,
  emptyMsg,
  errorMsg,
  classes,
}) => {
  const tableRef = useRef<HTMLElement>(null)
  const [lastLoadedTableHeight, setLastLoadedTableHeight] = useState<
    number | undefined
  >(undefined)
  useEffect(() => {
    if (tableRef.current !== null && !loading) {
      setLastLoadedTableHeight(tableRef.current.offsetHeight)
    }
  }, [loading, setLastLoadedTableHeight, tableRef])

  const heightCss: CSSProperties = { height: lastLoadedTableHeight }

  return (
    <Paper className={classes.root}>
      <RootRef rootRef={tableRef}>
        <div style={heightCss}>
          <MuiTable className={classes.table}>
            <TableHead>
              <TableRow className={classes.header}>
                {headers.map((h: string) => (
                  <MuiTableCell key={h}>{h}</MuiTableCell>
                ))}
              </TableRow>
            </TableHead>
            <TableBody>
              <Rows
                loading={loading}
                error={error}
                headers={headers}
                rows={rows}
                count={count}
                loadingMsg={loadingMsg}
                emptyMsg={emptyMsg}
                errorMsg={errorMsg}
              />
            </TableBody>
          </MuiTable>
        </div>
      </RootRef>
      <TablePagination
        component="div"
        count={count ?? 0}
        rowsPerPageOptions={[]}
        rowsPerPage={rowsPerPage}
        page={currentPage}
        SelectProps={{
          native: true,
        }}
        onChangePage={onChangePage}
        ActionsComponent={PaginationActions}
      />
    </Paper>
  )
}

type RowProps = Pick<
  Props,
  | 'loading'
  | 'error'
  | 'headers'
  | 'rows'
  | 'count'
  | 'loadingMsg'
  | 'emptyMsg'
  | 'errorMsg'
>

function Rows(props: RowProps) {
  if (props.loading) {
    return <Loading colCount={props.headers.length} msg={props.loadingMsg} />
  } else if (props.error) {
    return <Error colCount={props.headers.length} msg={props.errorMsg} />
  } else if (props.count === 0) {
    return <Empty colCount={props.headers.length} msg={props.emptyMsg} />
  } else {
    return (
      <>
        {props.rows?.map((r: any[], idx: number) => (
          <TableRow key={idx}>
            {r.map((col: Column, idx: number) => (
              <TableCell key={idx} column={col} />
            ))}
          </TableRow>
        ))}
      </>
    )
  }
}

interface MsgProps {
  colCount: number
  msg?: string
}

function Loading({ colCount, msg }: MsgProps) {
  return (
    <TableRow>
      <MuiTableCell component="th" scope="row" colSpan={colCount}>
        {msg ?? 'Loading...'}
      </MuiTableCell>
    </TableRow>
  )
}

function Empty({ colCount, msg }: MsgProps) {
  return (
    <TableRow>
      <MuiTableCell component="th" scope="row" colSpan={colCount}>
        {msg ?? 'There are no results added to the Explorer yet.'}
      </MuiTableCell>
    </TableRow>
  )
}

function Error({ colCount, msg }: MsgProps) {
  return (
    <TableRow>
      <MuiTableCell component="th" scope="row" colSpan={colCount}>
        {msg ?? 'Error loading resources.'}
      </MuiTableCell>
    </TableRow>
  )
}

export default withStyles(styles)(Table)
