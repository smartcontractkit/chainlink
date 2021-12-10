import React from 'react'
import { useHistory, useLocation } from 'react-router-dom'

import { Card, TablePagination } from '@material-ui/core'

import { JobRunsTable } from 'src/components/Table/JobRunsTable'

interface Props {
  job: JobPayload_Fields
  fetchMore: (page: number, per: number) => void
}

export const TabRuns: React.FC<Props> = ({ fetchMore, job }) => {
  const history = useHistory()
  const location = useLocation()
  const params = new URLSearchParams(location.search)

  const [{ page, pageSize }, setPagination] = React.useState<{
    page: number
    pageSize: number
  }>({
    page: parseInt(params.get('page') || '1', 10),
    pageSize: parseInt(params.get('per') || '10', 10),
  })

  React.useEffect(() => {
    fetchMore(page, pageSize)
  }, [fetchMore, page, pageSize])

  React.useEffect(() => {
    history.replace({
      search: `?page=${page}&size=${pageSize}`,
    })
  }, [history, page, pageSize])

  const runs = React.useMemo(() => {
    return job.runs.results.map(
      ({ allErrors, id, createdAt, finishedAt, status }) => ({
        id,
        createdAt,
        errors: allErrors,
        finishedAt,
        status,
      }),
    )
  }, [job.runs])

  return (
    <>
      <Card>
        <JobRunsTable runs={runs} />

        <TablePagination
          component="div"
          count={job.runs.metadata.total}
          rowsPerPage={pageSize}
          rowsPerPageOptions={[10, 25, 50, 100]}
          page={page - 1}
          onChangePage={(_, p) => {
            setPagination({ page: p + 1, pageSize })
          }}
          onChangeRowsPerPage={(e) => {
            setPagination({ page: 1, pageSize: parseInt(e.target.value, 10) })
          }}
          backIconButtonProps={{ 'aria-label': 'prev-page' }}
          nextIconButtonProps={{ 'aria-label': 'next-page' }}
        />
      </Card>
    </>
  )
}
