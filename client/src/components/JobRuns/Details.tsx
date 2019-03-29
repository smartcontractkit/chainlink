import React from 'react'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableRow from '@material-ui/core/TableRow'
import TableCell from '@material-ui/core/TableCell'
import Typography from '@material-ui/core/Typography'

interface IColProps {
  children: React.ReactNode
}

const Col = ({ children }: IColProps) => (
  <TableCell>
    <Typography variant="body1">
      <span>{children}</span>
    </Typography>
  </TableCell>
)

interface IProps {
  jobRun: IJobRun
}

const Details = ({ jobRun }: IProps) => {
  return (
    <Table>
      <TableBody>
        <TableRow>
          <Col>Job ID</Col>
          <Col>{jobRun.jobId}</Col>
        </TableRow>
        <TableRow>
          <Col>Status</Col>
          <Col>{jobRun.status}</Col>
        </TableRow>
        <TableRow>
          <Col>Initiator</Col>
          <Col>{jobRun.initiatorType}</Col>
        </TableRow>
        <TableRow>
          <Col>Completed At</Col>
          <Col>{jobRun.completedAt}</Col>
        </TableRow>
        {jobRun.error && (
          <TableRow>
            <Col>Error</Col>
            <Col>{jobRun.error}</Col>
          </TableRow>
        )}
      </TableBody>
    </Table>
  )
}

export default Details
