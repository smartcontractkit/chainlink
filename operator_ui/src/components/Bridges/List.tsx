import Card from '@material-ui/core/Card'
import { withStyles, WithStyles } from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TablePagination from '@material-ui/core/TablePagination'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'
import { tableStyles } from 'components/Table'
import Link from 'components/Link'
import TableButtons, { FIRST_PAGE } from 'components/TableButtons'
import React, { useState, useEffect } from 'react'
import { RouteComponentProps } from 'react-router-dom'

const renderFetching = () => (
  <TableRow>
    <TableCell component="th" scope="row" colSpan={4}>
      ...
    </TableCell>
  </TableRow>
)

const renderError = (error: string) => (
  <TableRow>
    <TableCell component="th" scope="row" colSpan={4}>
      {error}
    </TableCell>
  </TableRow>
)

interface BridgeRowsProps extends WithStyles<typeof tableStyles> {
  bridges: any[]
}

const BridgeRows = withStyles(tableStyles)(
  ({ bridges, classes }: BridgeRowsProps) => {
    return (
      <>
        {bridges.map((bridge) => (
          <TableRow className={classes.row} key={bridge.name} hover>
            <TableCell scope="row" component="th">
              <Link className={classes.link} href={`/bridges/${bridge.name}`}>
                {bridge.name}
              </Link>
            </TableCell>
            <TableCell>
              <Typography variant="body1">{bridge.url}</Typography>
            </TableCell>
            <TableCell>
              <Typography variant="body1">{bridge.confirmations}</Typography>
            </TableCell>
            <TableCell>
              <Typography variant="body1">
                {bridge.minimumContractPayment}
              </Typography>
            </TableCell>
          </TableRow>
        ))}
      </>
    )
  },
)

const renderBody = (bridges: any[], fetching: boolean, error: string) => {
  if (fetching) {
    return renderFetching()
  } else if (error) {
    return renderError(error)
  } else {
    return <BridgeRows bridges={bridges} />
  }
}

// CHECKME
interface OwnProps {
  bridges: any[]
  bridgeCount: number
  pageSize: number
  fetching: boolean
  error: string
  fetchBridges: (...args: any[]) => any
}

// CHECKME
type RouteProps = RouteComponentProps<{
  bridgePage: string
}>

type Props = OwnProps & RouteProps

export const BridgeList: React.FC<Props> = ({
  bridges,
  bridgeCount,
  error,
  fetchBridges,
  fetching,
  history,
  match,
  pageSize,
}) => {
  const [page, setPage] = useState(FIRST_PAGE)

  useEffect(() => {
    const queryPage = parseInt(match?.params.bridgePage, 10) || FIRST_PAGE
    setPage(queryPage)
    fetchBridges(queryPage, pageSize)
  }, [fetchBridges, pageSize, match])

  const handleChangePage = (_: never, page: React.SetStateAction<number>) => {
    fetchBridges(page, pageSize)
    setPage(page)
  }

  const TableButtonsWithProps = () => (
    <TableButtons
      history={history}
      count={bridgeCount}
      onChangePage={handleChangePage}
      rowsPerPage={pageSize}
      page={page}
      replaceWith={`/bridges/page`}
    />
  )

  return (
    <Card>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>
              <Typography variant="body1" color="textSecondary">
                Name
              </Typography>
            </TableCell>
            <TableCell>
              <Typography variant="body1" color="textSecondary">
                URL
              </Typography>
            </TableCell>
            <TableCell>
              <Typography variant="body1" color="textSecondary">
                Default Confirmations
              </Typography>
            </TableCell>
            <TableCell>
              <Typography variant="body1" color="textSecondary">
                Minimum Contract Payment
              </Typography>
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody>{renderBody(bridges, fetching, error)}</TableBody>
      </Table>
      <TablePagination
        component="div"
        count={bridgeCount}
        rowsPerPage={pageSize}
        rowsPerPageOptions={[pageSize]}
        page={page - 1}
        onChangePage={
          () => {} /* handler required by component, so make it a no-op */
        }
        onChangeRowsPerPage={
          () => {} /* handler required by component, so make it a no-op */
        }
        ActionsComponent={TableButtonsWithProps}
      />
    </Card>
  )
}

export default BridgeList
