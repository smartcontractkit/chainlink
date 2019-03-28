import React from 'react'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableRow from '@material-ui/core/TableRow'
import TableCell from '@material-ui/core/TableCell'
import Typography from '@material-ui/core/Typography'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import TaskRuns from './TaskRuns'

interface IColProps {
  children: React.ReactNode
  className?: string
}

const BaseCol = ({ children, className }: IColProps) => (
  <TableCell className={className}>{children}</TableCell>
)

const Col = ({ children, className }: IColProps) => (
  <BaseCol className={className}>
    <Typography variant="body1">{children}</Typography>
  </BaseCol>
)

const styles = ({ spacing }: Theme) =>
  createStyles({
    bottomCol: {
      borderBottom: 'none'
    }
  })

interface IProps extends WithStyles<typeof styles> {
  jobRun: IJobRun
}

const Details = ({ classes, jobRun }: IProps) => {
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
        <TableRow>
          <Col className={classes.bottomCol}>Tasks</Col>
          <BaseCol className={classes.bottomCol}>
            <TaskRuns taskRuns={jobRun.taskRuns} />
          </BaseCol>
        </TableRow>
      </TableBody>
    </Table>
  )
}

export default withStyles(styles)(Details)
