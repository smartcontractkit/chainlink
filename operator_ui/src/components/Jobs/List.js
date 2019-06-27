import React from 'react'
import PropTypes from 'prop-types'
import { useHooks, useState, useEffect } from 'use-react-hooks'
import Card from '@material-ui/core/Card'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import TablePagination from '@material-ui/core/TablePagination'
import Typography from '@material-ui/core/Typography'
import TimeAgo from '@chainlink/styleguide/components/TimeAgo'
import { formatInitiators } from 'utils/jobSpecInitiators'
import TableButtons, { FIRST_PAGE } from 'components/TableButtons'
import Link from 'components/Link'

const renderBody = (jobs, error) => {
  if (error) {
    return (
      <TableRow>
        <TableCell component="th" scope="row" colSpan={3}>
          {error}
        </TableCell>
      </TableRow>
    )
  } else if (jobs && jobs.length === 0) {
    return (
      <TableRow>
        <TableCell component="th" scope="row" colSpan={3}>
          You haven’t created any jobs yet. Create a new job{' '}
          <Link to={`/jobs/new`}>here</Link>
        </TableCell>
      </TableRow>
    )
  } else if (jobs) {
    return jobs.map(j => (
      <TableRow key={j.id}>
        <TableCell component="th" scope="row">
          <Link to={`/jobs/${j.id}`}>{j.id}</Link>
        </TableCell>
        <TableCell>
          <Typography variant="body1">
            <TimeAgo tooltip>{j.createdAt}</TimeAgo>
          </Typography>
        </TableCell>
        <TableCell>
          <Typography variant="body1">
            {formatInitiators(j.initiators)}
          </Typography>
        </TableCell>
      </TableRow>
    ))
  }

  return (
    <TableRow>
      <TableCell component="th" scope="row" colSpan={3}>
        Loading...
      </TableCell>
    </TableRow>
  )
}

export const List = useHooks(props => {
  const [page, setPage] = useState(FIRST_PAGE)
  useEffect(() => {
    const queryPage =
      (props.match && parseInt(props.match.params.jobPage, 10)) || FIRST_PAGE
    setPage(queryPage)
    fetchJobs(queryPage, pageSize)
  }, [])
  const { jobs, jobCount, fetchJobs, pageSize, error } = props
  const handleChangePage = (e, page) => {
    fetchJobs(page, pageSize)
    setPage(page)
  }
  const TableButtonsWithProps = () => (
    <TableButtons
      {...props}
      count={jobCount}
      onChangePage={handleChangePage}
      rowsPerPage={pageSize}
      page={page}
      replaceWith={`/jobs/page`}
    />
  )

  return (
    <Card>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>
              <Typography variant="body1" color="textSecondary">
                ID
              </Typography>
            </TableCell>
            <TableCell>
              <Typography variant="body1" color="textSecondary">
                Created
              </Typography>
            </TableCell>
            <TableCell>
              <Typography variant="body1" color="textSecondary">
                Initiator
              </Typography>
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody>{renderBody(jobs, error)}</TableBody>
      </Table>
      <TablePagination
        component="div"
        count={jobCount}
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
})

List.propTypes = {
  jobs: PropTypes.array,
  jobCount: PropTypes.number.isRequired,
  pageSize: PropTypes.number.isRequired,
  error: PropTypes.string,
  fetchJobs: PropTypes.func.isRequired
}

export default List
