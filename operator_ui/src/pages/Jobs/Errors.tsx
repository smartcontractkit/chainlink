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
import { TimeAgo } from 'components/TimeAgo'
import { JobData } from './sharedTypes'

export const JobsErrors: React.FC<{
  error: unknown
  ErrorComponent: React.FC
  LoadingPlaceholder: React.FC
  job?: JobData['job']
  setState: React.Dispatch<React.SetStateAction<JobData>>
  getJobSpec: () => Promise<void>
}> = ({
  error,
  ErrorComponent,
  getJobSpec,
  LoadingPlaceholder,
  job,
  setState,
}) => {
  React.useEffect(() => {
    document.title = job?.name ? `${job?.name} | Job errors` : 'Job errors'
  }, [job])

  const handleDismiss = async (jobSpecErrorId: string) => {
    // Optimistic delete
    const jobCopy: NonNullable<JobData['job']> = JSON.parse(JSON.stringify(job))
    jobCopy.errors = jobCopy.errors.filter((e) => e.id !== jobSpecErrorId)
    setState((state) => ({ ...state, job: jobCopy }))

    await v2.jobSpecErrors.destroyJobSpecError(jobSpecErrorId)
    await getJobSpec()
  }

  const tableHeaders = ['Occurrences', 'Created', 'Last Seen', 'Message']

  if (job?.type === 'Direct request') {
    tableHeaders.push('Actions')
  }

  return (
    <Content>
      <ErrorComponent />
      <LoadingPlaceholder />
      {!error && job && (
        <Card>
          <Table>
            <TableHead>
              <TableRow>
                {tableHeaders.map((header) => (
                  <TableCell key={header}>
                    <Typography variant="body1" color="textSecondary">
                      {header}
                    </Typography>
                  </TableCell>
                ))}
              </TableRow>
            </TableHead>
            <TableBody>
              {job.errors.length === 0 ? (
                <TableRow>
                  <TableCell component="th" scope="row" colSpan={5}>
                    No errors
                  </TableCell>
                </TableRow>
              ) : (
                job.errors.map((jobSpecError) => (
                  <TableRow key={jobSpecError.id}>
                    <TableCell>
                      <Typography variant="body1">
                        {jobSpecError.occurrences}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body1">
                        <TimeAgo tooltip>{jobSpecError.createdAt}</TimeAgo>
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body1">
                        <TimeAgo tooltip>{jobSpecError.updatedAt}</TimeAgo>
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body1">
                        {jobSpecError.description}
                      </Typography>
                    </TableCell>
                    {job?.type === 'Direct request' && (
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
                    )}
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
