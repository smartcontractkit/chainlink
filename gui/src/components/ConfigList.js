import React, { Component, Fragment } from 'react'
import PropTypes from 'prop-types'
import Card from '@material-ui/core/Card'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'

const renderFetching = () => (
  <TableRow>
    <TableCell component='th' scope='row' colSpan={3}>...</TableCell>
  </TableRow>
)

const renderError = error => (
  <TableRow>
    <TableCell component='th' scope='row' colSpan={3}>
      {error}
    </TableCell>
  </TableRow>
)

const renderConfigs = configs => (
  configs.map(([k, v]) => (
    <TableRow key={k}>
      <TableCell>
        <Typography variant='body1'>
          <Fragment>
            {k}
          </Fragment>
        </Typography>
      </TableCell>
      <TableCell>
        <Typography variant='body1'>
          <Fragment>
            {v}
          </Fragment>
        </Typography>
      </TableCell>
    </TableRow>
  ))
)

const renderBody = (configs, error) => {
  if (error) {
    return renderError(error)
  } else if (configs.length === 0) {
    return renderFetching()
  } else {
    return renderConfigs(configs)
  }
}

export class ConfigList extends Component {
  constructor (props) {
    super(props)
    this.state = {}
  }

  render () {
    const {configs, error} = this.props

    return (
      <Card>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>
                <Typography variant='body1' color='textSecondary'>Key</Typography>
              </TableCell>
              <TableCell>
                <Typography variant='body1' color='textSecondary'>Value</Typography>
              </TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {renderBody(configs, error)}
          </TableBody>
        </Table>
      </Card>
    )
  }
}

ConfigList.propTypes = {
  configs: PropTypes.array.isRequired,
  error: PropTypes.string
}

export default ConfigList
