import React from 'react'
import { useHistory, useLocation } from 'react-router-dom'
import { Card, TablePagination } from '@material-ui/core'
import Content from 'components/Content'
import NodesList from './NodesList'
import TableButtons from 'components/TableButtons'

interface Props {
  //   ErrorComponent: React.FC
  //   LoadingPlaceholder: React.FC
  //   error: unknown
  //   getJobRuns: (props: { page: number; size: number }) => Promise<void>
}

const NodesIndex = () =>
  //   {
  //       ErrorComponent,
  //       LoadingPlaceholder,
  //       error,
  //       getJobRuns,
  //     }: Props) => {
  //   },
  {
    const location = useLocation()
    const params = new URLSearchParams(location.search)
    const [{ page, pageSize }, setPagination] = React.useState<{
      page: number
      pageSize: number
    }>({
      page: parseInt(params.get('page') || '1', 10),
      pageSize: parseInt(params.get('size') || '10', 10),
    })

    const history = useHistory()

    //   React.useEffect(() => {
    //     getJobRuns({ page, size: pageSize })
    //     history.replace({
    //       search: `?page=${page}&size=${pageSize}`,
    //     })
    //   }, [getJobRuns, history, page, pageSize])

    return (
      <Content>
        {/* <ErrorComponent />
        <LoadingPlaceholder /> */}
        <Card>
          {/* {!error && <JobRunsList runs={recentRuns} />} */}
          <TablePagination
            component="div"
            count={25}
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
                count={25}
                onChangePage={(_: React.SyntheticEvent, newPage: number) =>
                  setPagination({ page: newPage, pageSize })
                }
                page={page}
                rowsPerPage={pageSize}
              />
            )}
          />
        </Card>
      </Content>
    )
  }

export default NodesIndex
