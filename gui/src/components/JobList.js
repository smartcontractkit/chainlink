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
import formatInitiators from 'utils/formatInitiators'

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

const renderJobs = jobs => (
  jobs.map(j => (
    <TableRow key={j.id}>
      <TableCell component='th' scope='row'>
        <Link to={`/job_specs/${j.id}`}>{j.id}</Link>
      </TableCell>
      <TableCell>
        <Typography variant='body1'>{j.createdAt}</Typography>
      </TableCell>
      <TableCell>
        <Typography variant='body1'>{formatInitiators(j.initiators)}</Typography>
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
              <TableCell>
                <Typography variant='body1' color='textSecondary'>ID</Typography>
              </TableCell>
              <TableCell>
                <Typography variant='body1' color='textSecondary'>Created</Typography>
              </TableCell>
              <TableCell>
                <Typography variant='body1' color='textSecondary'>Initiator</Typography>
              </TableCell>
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
          onChangeRowsPerPage={() => {} /* handler required by component, so make it a no-op */}
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
