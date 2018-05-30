import React from 'react'
import PropTypes from 'prop-types'
import Card from '@material-ui/core/Card'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'

const formatInitiators = (initiators) => (initiators.map(i => i.type).join(', '))

const renderJobs = (jobs, fetching, error) => {
  if (fetching) {
    return (
      <TableRow>
        <TableCell component='th' scope='row' colSpan={3}>
          ...
        </TableCell>
      </TableRow>
    )
  } else if (error) {
    return (
      <TableRow>
        <TableCell component='th' scope='row' colSpan={3}>
          {error}
        </TableCell>
      </TableRow>
    )
  } else {
    return jobs.map(j => (
      <TableRow key={j.id}>
        <TableCell component='th' scope='row'>
          {j.id}
        </TableCell>
        <TableCell>{j.createdAt}</TableCell>
        <TableCell>
          {formatInitiators(j.initiators)}
        </TableCell>
      </TableRow>
    ))
  }
}

export const JobList = ({jobs, fetching, error}) => (
  <Card>
    <Table>
      <TableHead>
        <TableRow>
          <TableCell>ID</TableCell>
          <TableCell>Created</TableCell>
          <TableCell>Initiator</TableCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {renderJobs(jobs, fetching, error)}
      </TableBody>
    </Table>
  </Card>
)

JobList.propTypes = {
  jobs: PropTypes.array.isRequired,
  fetching: PropTypes.bool,
  error: PropTypes.string
}

export default JobList
