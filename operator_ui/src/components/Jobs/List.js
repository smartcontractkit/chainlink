import React from 'react'
import PropTypes from 'prop-types'
import { useEffect, useHooks, useState } from 'use-react-hooks'
import { formatInitiators } from 'utils/jobSpecInitiators'
import { FIRST_PAGE, GenericList } from '../GenericList'

const buildItems = jobs =>
  jobs.map(j => [
    { type: 'link', text: j.id, to: `/jobs/${j.id}` },
    { type: 'time_ago', text: j.createdAt },
    { type: 'text', text: formatInitiators(j.initiators) },
  ])

export const List = useHooks(props => {
  const { jobs, jobCount, fetchJobs, pageSize } = props
  const [page, setPage] = useState(FIRST_PAGE)
  useEffect(() => {
    const queryPage =
      (props.match && parseInt(props.match.params.jobPage, 10)) || FIRST_PAGE
    setPage(queryPage - 1)
    fetchJobs(queryPage, pageSize)
  }, [])
  const handleChangePage = (e, page) => {
    if (e) {
      setPage(page)
      fetchJobs(page + 1, pageSize)
      if (props.history) props.history.push(`/jobs/page/${page + 1}`)
    }
  }
  return (
    <GenericList
      emptyMsg="You haven’t created any jobs yet. Create a new job"
      headers={['ID', 'Created', 'Initiator']}
      items={jobs && buildItems(jobs)}
      onChangePage={handleChangePage}
      count={jobCount}
      rowsPerPage={pageSize}
      currentPage={page}
    />
  )
})

List.propTypes = {
  jobs: PropTypes.array,
  jobCount: PropTypes.number.isRequired,
  pageSize: PropTypes.number.isRequired,
  error: PropTypes.string,
  fetchJobs: PropTypes.func.isRequired,
}

export default List
