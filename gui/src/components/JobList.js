import React, { Component } from 'react'
import PropTypes from 'prop-types'
import Card from '@material-ui/core/Card'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import TablePagination from '@material-ui/core/TablePagination'
import Typography from '@material-ui/core/Typography'
import { formatInitiators } from 'utils/jobSpecInitiators'
import TableButtons, { FIRST_PAGE } from 'components/TableButtons'
import Link from 'components/Link'
import TimeAgo from 'components/TimeAgo'

const renderBody = (jobs, error) => {
  if (error) {
    return (
      <TableRow>
        <TableCell component='th' scope='row' colSpan={3}>
          {error}
        </TableCell>
      </TableRow>
    )
  } else if (jobs && jobs.length === 0) {
    return (
      <TableRow>
        <TableCell component='th' scope='row' colSpan={3}>
          You haven't created any jobs yet. Create a new job <Link to={`/jobs/new`}>here</Link>
        </TableCell>
      </TableRow>
    )
  } else if (jobs) {
    return jobs.map(j => (
      <TableRow key={j.id}>
        <TableCell component='th' scope='row'>
          <Link to={`/jobs/${j.id}`}>{j.id}</Link>
        </TableCell>
        <TableCell>
          <Typography variant='body1'>
            <TimeAgo>{j.createdAt}</TimeAgo>
          </Typography>
        </TableCell>
        <TableCell>
          <Typography variant='body1'>
            {formatInitiators(j.initiators)}
          </Typography>
        </TableCell>
      </TableRow>
    ))
  }

  return (
    <TableRow>
      <TableCell component='th' scope='row' colSpan={3}>
        Loading...
      </TableCell>
    </TableRow>
  )
}

export class JobList extends Component {
  constructor (props) {
    super(props)
    this.state = { page: FIRST_PAGE }
    this.handleChangePage = this.handleChangePage.bind(this)
  }

  componentDidMount () {
    const { pageSize, fetchJobs } = this.props
    const queryPage = this.props.match ? (parseInt(this.props.match.params.jobPage, 10) || FIRST_PAGE) : FIRST_PAGE
    this.setState({ page: queryPage })
    fetchJobs(queryPage, pageSize)
  }

  componentDidUpdate (prevProps) {
    if (prevProps.match && this.props.match) {
      const prevJobPage = prevProps.match.params.jobPage
      const currentJobPage = this.props.match.params.jobPage

      if (prevJobPage !== currentJobPage) {
        const { pageSize, fetchJobs } = this.props
        this.setState({page: parseInt(currentJobPage, 10) || FIRST_PAGE})
        fetchJobs(parseInt(currentJobPage, 10) || FIRST_PAGE, pageSize)
      }
    }
  }

  handleChangePage (e, page) {
    const { fetchJobs, pageSize } = this.props
    fetchJobs(page, pageSize)
    this.setState({ page })
  }

  render () {
    const { jobs, jobCount, pageSize, error } = this.props
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
            {renderBody(jobs, error)}
          </TableBody>
        </Table>
        <TablePagination
          component='div'
          count={jobCount}
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

JobList.propTypes = {
  jobs: PropTypes.array,
  jobCount: PropTypes.number.isRequired,
  pageSize: PropTypes.number.isRequired,
  error: PropTypes.string,
  fetchJobs: PropTypes.func.isRequired
}

export default JobList
