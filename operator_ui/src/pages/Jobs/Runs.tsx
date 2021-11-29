import React from 'react'
import { useHistory, useLocation } from 'react-router-dom'
import { Card, TablePagination } from '@material-ui/core'
import JobRunsList from './JobRunsList'
import TableButtons from 'components/TableButtons'
import { JobData } from './sharedTypes'

interface Props {
  ErrorComponent: React.FC
  LoadingPlaceholder: React.FC
  error: unknown
  getJobRuns: (props: { page: number; size: number }) => Promise<void>
  job?: JobData['job']
  recentRuns?: JobData['recentRuns']
  recentRunsCount: JobData['recentRunsCount']
}

export const Runs = ({
  ErrorComponent,
  LoadingPlaceholder,
  error,
  job,
  getJobRuns,
  recentRuns,
  recentRunsCount,
}: Props) => {
  const location = useLocation()
  const params = new URLSearchParams(location.search)
  const [{ page, pageSize }, setPagination] = React.useState<{
    page: number
    pageSize: number
  }>({
    page: parseInt(params.get('page') || '1', 10),
    pageSize: parseInt(params.get('size') || '10', 10),
  })

  React.useEffect(() => {
    document.title = job?.name ? `${job.name} | Job runs` : 'Job runs'
  }, [job])

  const history = useHistory()

  React.useEffect(() => {
    getJobRuns({ page, size: pageSize })
    history.replace({
      search: `?page=${page}&size=${pageSize}`,
    })
  }, [getJobRuns, history, page, pageSize])

  return (
    <>
      <ErrorComponent />
      <LoadingPlaceholder />
      <Card>
        {!error && recentRuns && <JobRunsList runs={recentRuns} />}
        <TablePagination
          component="div"
          count={recentRunsCount}
          rowsPerPage={pageSize}
          rowsPerPageOptions={[10, 25, 50, 100]}
          page={page - 1}
          onChangePage={
            () => {} /* handler required by component, so make it a no-op */
          }
          onChangeRowsPerPage={(e) => {
            setPagination({ page: 1, pageSize: parseInt(e.target.value, 10) })
          }}
          ActionsComponent={() => (
            <TableButtons
              count={recentRunsCount}
              onChangePage={(_: React.SyntheticEvent, newPage: number) =>
                setPagination({ page: newPage, pageSize })
              }
              page={page}
              rowsPerPage={pageSize}
            />
          )}
        />
      </Card>
    </>
  )
}
