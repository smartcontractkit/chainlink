import React, { Fragment } from 'react'
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

const renderConfigs = entries => (
  entries.map(([k, v]) => (
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

const renderBody = (entries, error) => {
  if (error) {
    return renderError(error)
  } else if (entries.length === 0) {
    return renderFetching()
  } else {
    return renderConfigs(entries)
  }
}

const KeyValueList = ({entries, error, showHead}) => (
  <Card>
    <Table>
      {showHead &&
        <TableHead>
          <TableRow>
            <TableCell>
              <Typography variant='body1' color='textSecondary'>Key</Typography>
            </TableCell>
            <TableCell>
              <Typography variant='body1' color='textSecondary'>Value</Typography>
            </TableCell>
          </TableRow>
        </TableHead>}
      <TableBody>
        {renderBody(entries, error)}
      </TableBody>
    </Table>
  </Card>
)

KeyValueList.propTypes = {
  showHead: PropTypes.bool.isRequired,
  entries: PropTypes.array.isRequired,
  error: PropTypes.string
}

KeyValueList.defaultProps = {
  showHead: false
}

export default KeyValueList
