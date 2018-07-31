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
import TableButtons from 'components/TableButtons'

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
    this.state = { page: 1 }
    this.handleChangePage = this.handleChangePage.bind(this)
  }

  componentDidMount () {
    const { pageSize, fetchJobs } = this.props
    const firstPage = 1
    if (this.props.match.params.page) {
      const queryPage = this.props.match.params.page
      this.setState({ page: queryPage })
      fetchJobs(queryPage, pageSize)
    } else {
      this.setState({ page: firstPage })
      fetchJobs(firstPage, pageSize)
    }
  }

  handleChangePage (e, page) {
    const {fetchJobs, pageSize} = this.props
    fetchJobs(page, pageSize)
    this.setState({ page })
  }

  render () {
    const {jobs, jobCount, pageSize, fetching, error} = this.props
    const TableButtonsWithProps = () => (
      <TableButtons
        {...this.props}
        count={jobCount}
        onChangePage={this.handleChangePage}
        rowsPerPage={pageSize}
        page={this.state.page}
        replaceWith={`/jobs/page`}
      />
    )

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
          page={this.state.page - 1}
          onChangePage={() => {} /* handler required by component, so make it a no-op */}
          onChangeRowsPerPage={() => {} /* handler required by component, so make it a no-op */}
          ActionsComponent={TableButtonsWithProps}
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
