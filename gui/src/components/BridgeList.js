import React, { Component } from 'react'
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
    <TableCell component='th' scope='row' colSpan={4}>...</TableCell>
  </TableRow>
)

const renderError = error => (
  <TableRow>
    <TableCell component='th' scope='row' colSpan={4}>
      {error}
    </TableCell>
  </TableRow>
)

const renderBridges = bridges => (
  bridges.map(bridge => (
    <TableRow key={bridge.name}>
      <TableCell scope='row' component='th'>
        <Link to={`/bridges/${bridge.name}`}>
          {bridge.name}
        </Link>
      </TableCell>
      <TableCell>
        <Typography variant='body1'>{bridge.url}</Typography>
      </TableCell>
      <TableCell>
        <Typography variant='body1'>{bridge.confirmations}</Typography>
      </TableCell>
      <TableCell>
        <Typography variant='body1'>{bridge.minimumContractPayment}</Typography>
      </TableCell>
    </TableRow>
  ))
)

const renderBody = (bridges, fetching, error) => {
  if (fetching) {
    return renderFetching()
  } else if (error) {
    return renderError(error)
  } else {
    return renderBridges(bridges)
  }
}

export class BridgeList extends Component {
  constructor (props) {
    super(props)
    this.state = { page: 1 }
    this.handleChangePage = this.handleChangePage.bind(this)
  }

  handleChangePage (e, page) {
    const {fetchBridges, pageSize} = this.props
    fetchBridges(page, pageSize)
    this.setState({page})
  }

  componentDidMount () {
    const { pageSize, fetchBridges } = this.props
    const queryPage = this.props.match ? (parseInt(this.props.match.params.bridgePage, 10) || FIRST_PAGE) : FIRST_PAGE
    this.setState({ page: queryPage })
    fetchBridges(queryPage, pageSize)
  }

  componentDidUpdate (prevProps) {
    if (prevProps.match && this.props.match) {
      const prevBridgePage = prevProps.match.params.bridgePage
      const currentBridgePage = this.props.match.params.bridgePage

      if (prevBridgePage !== currentBridgePage) {
        const { pageSize, fetchBridges } = this.props
        this.setState({page: parseInt(currentBridgePage, 10) || FIRST_PAGE})
        fetchBridges(parseInt(currentBridgePage, 10) || FIRST_PAGE, pageSize)
      }
    }
  }

  render () {
    const {bridges, bridgeCount, pageSize, fetching, error} = this.props
    const TableButtonsWithProps = () => (
      <TableButtons
        {...this.props}
        count={bridgeCount}
        onChangePage={this.handleChangePage}
        rowsPerPage={pageSize}
        page={this.state.page}
        replaceWith={`/bridges/page`}
      />
    )
    return (
      <Card>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>
                <Typography variant='body1' color='textSecondary'>
                  Name
                </Typography>
              </TableCell>
              <TableCell>
                <Typography variant='body1' color='textSecondary'>
                  URL
                </Typography>
              </TableCell>
              <TableCell>
                <Typography variant='body1' color='textSecondary'>
                  Default Confirmations
                </Typography>
              </TableCell>
              <TableCell>
                <Typography variant='body1' color='textSecondary'>
                  Minimum Contract Payment
                </Typography>
              </TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {renderBody(bridges, fetching, error)}
          </TableBody>
        </Table>
        <TablePagination
          component='div'
          count={bridgeCount}
          rowsPerPage={pageSize}
          rowsPerPageOptions={[pageSize]}
          page={this.state.page - 1}
          onChangePage={() => { } /* handler required by component, so make it a no-op */}
          onChangeRowsPerPage={() => { } /* handler required by component, so make it a no-op */}
          ActionsComponent={TableButtonsWithProps}
        />
      </Card>
    )
  }
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
