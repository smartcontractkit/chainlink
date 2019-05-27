import React from 'react'
import Link from 'components/Link'
import PropTypes from 'prop-types'
import Card from '@material-ui/core/Card'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import TablePagination from '@material-ui/core/TablePagination'
import Typography from '@material-ui/core/Typography'
import TableButtons, { FIRST_PAGE } from 'components/TableButtons'

const renderFetching = () => (
  <TableRow>
    <TableCell component="th" scope="row" colSpan={4}>
      ...
    </TableCell>
  </TableRow>
)

const renderError = error => (
  <TableRow>
    <TableCell component="th" scope="row" colSpan={4}>
      {error}
    </TableCell>
  </TableRow>
)

const renderBridges = bridges =>
  bridges.map(bridge => (
    <TableRow key={bridge.name}>
      <TableCell scope="row" component="th">
        <Link to={`/bridges/${bridge.name}`}>{bridge.name}</Link>
      </TableCell>
      <TableCell>
        <Typography variant="body1">{bridge.url}</Typography>
      </TableCell>
      <TableCell>
        <Typography variant="body1">{bridge.confirmations}</Typography>
      </TableCell>
      <TableCell>
        <Typography variant="body1">{bridge.minimumContractPayment}</Typography>
      </TableCell>
    </TableRow>
  ))

const renderBody = (bridges, fetching, error) => {
  if (fetching) {
    return renderFetching()
  } else if (error) {
    return renderError(error)
  } else {
    return renderBridges(bridges)
  }
}

export const BridgeList = props => {
  const [page, setPage] = React.useState(FIRST_PAGE)
  React.useEffect(() => {
    const queryPage =
      (props.match && parseInt(props.match.params.bridgePage, 10)) || FIRST_PAGE
    setPage(queryPage)
    fetchBridges(queryPage, pageSize)
  }, [])

  const {
    bridges,
    bridgeCount,
    fetchBridges,
    pageSize,
    fetching,
    error
  } = props
  const TableButtonsWithProps = () => (
    <TableButtons
      history={props.history}
      count={bridgeCount}
      onChangePage={handleChangePage}
      rowsPerPage={pageSize}
      page={page}
      replaceWith={`/bridges/page`}
    />
  )
  const handleChangePage = (e, page) => {
    fetchBridges(page, pageSize)
    setPage(page)
  }

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

BridgeList.propTypes = {
  bridges: PropTypes.array.isRequired,
  bridgeCount: PropTypes.number.isRequired,
  pageSize: PropTypes.number.isRequired,
  fetching: PropTypes.bool,
  error: PropTypes.string,
  fetchBridges: PropTypes.func.isRequired
}

export default BridgeList
