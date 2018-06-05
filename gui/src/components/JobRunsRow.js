import React from 'react'
import PropTypes from 'prop-types'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'

const statusColor = (status) => {
  if (status === 'error') {
    return 'error'
  }
}

const JobRunsRow = ({id, status, createdAt, result}) => (
  <TableRow>
    <TableCell component='th' scope='row'>
      <Typography variant='body1'>{id}</Typography>
    </TableCell>
    <TableCell component='th' scope='row'>
      <Typography variant='body1' color={statusColor(status)}>{status}</Typography>
    </TableCell>
    <TableCell>
      <Typography variant='body1'>{createdAt}</Typography>
    </TableCell>
    <TableCell>
      <Typography variant='body1'>{JSON.stringify(result.data)}</Typography>
    </TableCell>
    <TableCell>
      <Typography variant='body1' color='error'>
        {result.error && JSON.stringify(result.error)}
      </Typography>
    </TableCell>
  </TableRow>
)

JobRunsRow.propTypes = {
  id: PropTypes.string,
  status: PropTypes.string,
  createdAt: PropTypes.string,
  result: PropTypes.object
}

export default JobRunsRow
