import React, { Component } from 'react'
import PropTypes from 'prop-types'
import Card from '@material-ui/core/Card'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import TablePagination from '@material-ui/core/TablePagination'

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

const formatInitiators = (initiators) => (initiators.map(i => i.type).join(', '))
const renderJobs = jobs => (
  jobs.map(j => (
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
)

const renderBody = (jobs, fetching, error) => {
  if (fetching) {
    return renderFetching()
  } else if (error) {
    return renderError(error)
  } else {
    return renderJobs(jobs)
  }
}

const handleChangeRowsPerPage = () => {
  console.log('handleChangeRowsPerPage')
}

export class JobList extends Component {
  constructor (props) {
    super(props)
    this.state = {
      page: 0
    }
    this.handleChangePage = this.handleChangePage.bind(this)
  }

  handleChangePage (e, page) {
    const {fetchJobs, pageSize} = this.props

    fetchJobs(page + 1, pageSize)
    this.setState({page})
  }

  render () {
    const {jobs, jobCount, pageSize, fetching, error} = this.props

    return (
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
            {renderBody(jobs, fetching, error)}
          </TableBody>
        </Table>
        <TablePagination
          component='div'
          count={jobCount}
          rowsPerPage={pageSize}
          rowsPerPageOptions={[pageSize]}
          page={this.state.page}
          backIconButtonProps={{'aria-label': 'Previous Page'}}
          nextIconButtonProps={{'aria-label': 'Next Page'}}
          onChangePage={this.handleChangePage}
          onChangeRowsPerPage={handleChangeRowsPerPage}
        />
      </Card>
    )
  }
}

JobList.propTypes = {
  jobs: PropTypes.array.isRequired,
  jobCount: PropTypes.number.isRequired,
  pageSize: PropTypes.number.isRequired,
  fetching: PropTypes.bool,
  error: PropTypes.string,
  fetchJobs: PropTypes.func.isRequired
}

export default JobList
