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
import { TimeAgo } from 'components/TimeAgo'
import { JobData } from './sharedTypes'

export const JobsErrors: React.FC<{
  error: unknown
  ErrorComponent: React.FC
  LoadingPlaceholder: React.FC
  job?: JobData['job']
}> = ({ error, ErrorComponent, LoadingPlaceholder, job }) => {
  React.useEffect(() => {
    document.title = job?.name ? `${job?.name} | Job errors` : 'Job errors'
  }, [job])

  const tableHeaders = ['Occurrences', 'Created', 'Last Seen', 'Message']

  return (
    <>
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
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </Card>
      )}
    </>
  )
}

export default JobsErrors
