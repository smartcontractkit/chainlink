import React from 'react'
import { RouteComponentProps } from 'react-router-dom'
import Grid from '@material-ui/core/Grid'
import Button from 'components/Button'
import { Title } from 'components/Title'
import Content from 'components/Content'
import BaseLink from 'components/BaseLink'
import { v2 } from 'api'
import * as jsonapi from '@chainlink/json-api-client'
import * as models from 'core/store/models'

import { TimeAgo } from '@chainlink/styleguide'
import Card from '@material-ui/core/Card'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TablePagination from '@material-ui/core/TablePagination'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'
import TableButtons, { FIRST_PAGE } from 'components/TableButtons'
import { useHistory } from 'react-router-dom'
import { formatInitiators } from 'utils/jobSpecInitiators'
import Link from 'components/Link'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'

type IndexProps = {
  pageSize?: number
} & RouteComponentProps<{
  pageNumber?: string
}>

export const JobsIndex = ({ pageSize = 10, match }: IndexProps) => {
  React.useEffect(() => {
    document.title = 'Jobs'
  }, [])

  const [jobs, setJobs] = React.useState<
    jsonapi.PaginatedApiResponse<models.JobSpec[]>['data']
  >()
  const [jobsCount, setJobsCount] = React.useState(0)
  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !jobs)

  const history = useHistory()
  const pageNumber = match.params.pageNumber
    ? parseInt(match.params.pageNumber, 10)
    : FIRST_PAGE

  React.useEffect(() => {
    v2.specs
      .getJobSpecs(pageNumber, pageSize)
      .then(({ data, meta }) => {
        setJobs(data)
        setJobsCount(meta.count)
      })
      .catch(setError)
  }, [pageNumber, pageSize, setError])

  const TableButtonsWithProps = () => (
    <TableButtons
      count={jobsCount}
      onChangePage={(_event: React.SyntheticEvent, page: number) => {
        history.push(`/jobs/page/${page}`)
      }}
      rowsPerPage={pageSize}
      page={pageNumber}
      replaceWith={`/jobs/page`}
    />
  )

  return (
    <Content>
      <Grid container>
        <Grid item xs={9}>
          <Title>Jobs</Title>
        </Grid>
        <Grid item xs={3}>
          <Grid container justify="flex-end">
            <Grid item>
              <Button
                variant="secondary"
                component={BaseLink}
                href={'/jobs/new'}
              >
                New Job
              </Button>
            </Grid>
          </Grid>
        </Grid>
        <Grid item xs={12}>
          <ErrorComponent />
          <LoadingPlaceholder />
          {!error && jobs && (
            <Card>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>
                      <Typography variant="body1" color="textSecondary">
                        Name
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
                <TableBody>
                  {jobs && jobsCount === 0 && (
                    <TableRow>
                      <TableCell component="th" scope="row" colSpan={3}>
                        You haven’t created any jobs yet. Create a new job{' '}
                        <Link href={`/jobs/new`}>here</Link>
                      </TableCell>
                    </TableRow>
                  )}
                  {jobs &&
                    jobsCount > 0 &&
                    jobs.map((job) => (
                      <TableRow key={job.id}>
                        <TableCell component="th" scope="row">
                          <Link href={`/jobs/${job.id}`}>
                            {job.attributes.name || '-'}
                            <br />
                            <Typography
                              variant="subtitle2"
                              color="textSecondary"
                              component="span"
                            >
                              {job.id}
                            </Typography>
                          </Link>
                        </TableCell>
                        <TableCell>
                          <Typography variant="body1">
                            <TimeAgo tooltip>
                              {job.attributes.createdAt || ''}
                            </TimeAgo>
                          </Typography>
                        </TableCell>
                        <TableCell>
                          <Typography variant="body1">
                            {formatInitiators(job.attributes.initiators)}
                          </Typography>
                        </TableCell>
                      </TableRow>
                    ))}
                </TableBody>
              </Table>
              <TablePagination
                component="div"
                count={jobsCount}
                rowsPerPage={pageSize}
                rowsPerPageOptions={[pageSize]}
                page={pageNumber - 1}
                onChangePage={
                  () => {} /* handler required by component, so make it a no-op */
                }
                onChangeRowsPerPage={
                  () => {} /* handler required by component, so make it a no-op */
                }
                ActionsComponent={TableButtonsWithProps}
              />
            </Card>
          )}
        </Grid>
      </Grid>
    </Content>
  )
}

export default JobsIndex
