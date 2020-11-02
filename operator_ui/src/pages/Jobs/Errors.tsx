import React from 'react'
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
import { localizedTimestamp, TimeAgo } from '@chainlink/styleguide'
import { JobData } from './sharedTypes'

export const JobsErrors: React.FC<{
  error: unknown
  ErrorComponent: React.FC
  LoadingPlaceholder: React.FC
  jobSpec?: JobData['jobSpec']
  setState: React.Dispatch<React.SetStateAction<JobData>>
  getJobSpec: () => Promise<void>
}> = ({
  error,
  ErrorComponent,
  getJobSpec,
  LoadingPlaceholder,
  jobSpec,
  setState,
}) => {
  React.useEffect(() => {
    document.title =
      jobSpec && jobSpec.attributes.name
        ? `${jobSpec.attributes.name} | Job errors`
        : 'Job errors'
  }, [jobSpec])

  const handleDismiss = async (jobSpecErrorId: string) => {
    // Optimistic delete
    const jobSpecCopy: NonNullable<JobData['jobSpec']> = JSON.parse(
      JSON.stringify(jobSpec),
    )
    jobSpecCopy.attributes.errors = jobSpecCopy.attributes.errors.filter(
      (e) => e.id !== jobSpecErrorId,
    )
    setState((state) => ({ ...state, jobSpec: jobSpecCopy }))

    await v2.jobSpecErrors.destroyJobSpecError(jobSpecErrorId)
    await getJobSpec()
  }

  return (
    <Content>
      <ErrorComponent />
      <LoadingPlaceholder />
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
