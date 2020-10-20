import React from 'react'
import { RouteComponentProps } from 'react-router-dom'
import { ApiResponse } from '@chainlink/json-api-client'
import {
  Card,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Typography,
} from '@material-ui/core'
import { v2 } from 'api'
import Button from 'components/Button'
import Content from 'components/Content'
import { JobSpec } from 'core/store/models'
import { localizedTimestamp, TimeAgo } from '@chainlink/styleguide'

export const JobsErrors: React.FC<RouteComponentProps<{
  jobSpecId: string
}>> = ({ match }) => {
  const { jobSpecId } = match.params

  const [error, setError] = React.useState()
  const [jobSpec, setJobSpec] = React.useState<ApiResponse<JobSpec>['data']>()

  const fetchJobSpec = React.useCallback(
    async () =>
      v2.specs
        .getJobSpec(jobSpecId)
        .then((response) => setJobSpec(response.data))
        .catch(setError),
    [jobSpecId],
  )

  React.useEffect(() => {
    document.title = 'Job Errors'
  }, [])

  React.useEffect(() => {
    fetchJobSpec()
  }, [fetchJobSpec])

  const handleDismiss = async (jobSpecErrorId: string) => {
    // Optimistic delete
    const jobSpecCopy: ApiResponse<JobSpec>['data'] = JSON.parse(
      JSON.stringify(jobSpec),
    )
    jobSpecCopy.attributes.errors = jobSpecCopy.attributes.errors.filter(
      (e) => e.id !== jobSpecErrorId,
    )
    setJobSpec(jobSpecCopy)

    await v2.jobSpecErrors.destroyJobSpecError(jobSpecErrorId)
    fetchJobSpec()
  }

  return (
    <Content>
      {error && <div>Error while fetching data: {JSON.stringify(error)}</div>}
      {!error && !jobSpec && <div>Fetching...</div>}
      {!error && jobSpec && (
        <Card>
          <Table>
            <TableHead>
              <TableRow>
                {[
                  'Occurrences',
                  'Created',
                  'Last Seen',
                  'Message',
                  'Actions',
                ].map((header) => (
                  <TableCell key={header}>
                    <Typography variant="body1" color="textSecondary">
                      {header}
                    </Typography>
                  </TableCell>
                ))}
              </TableRow>
            </TableHead>
            <TableBody>
              {jobSpec.attributes.errors.length === 0 ? (
                <TableRow>
                  <TableCell component="th" scope="row" colSpan={5}>
                    No errors
                  </TableCell>
                </TableRow>
              ) : (
                jobSpec.attributes.errors.map((jobSpecError) => (
                  <TableRow key={jobSpecError.id}>
                    <TableCell>
                      <Typography variant="body1">
                        {jobSpecError.occurrences}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body1">
                        <TimeAgo tooltip>
                          {localizedTimestamp(
                            jobSpecError.createdAt.toString(),
                          )}
                        </TimeAgo>
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body1">
                        <TimeAgo tooltip>
                          {localizedTimestamp(
                            jobSpecError.updatedAt.toString(),
                          )}
                        </TimeAgo>
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body1">
                        {jobSpecError.description}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Button
                        variant="danger"
                        size="small"
                        onClick={() => {
                          handleDismiss(jobSpecError.id)
                        }}
                      >
                        Dismiss
                      </Button>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </Card>
      )}
    </Content>
  )
}

export default JobsErrors
