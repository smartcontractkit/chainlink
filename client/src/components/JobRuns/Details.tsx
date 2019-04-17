import React from 'react'
import classNames from 'classnames'
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

const colStyles = ({ spacing }: Theme) =>
  createStyles({
    col: {
      verticalAlign: 'top',
      paddingTop: spacing.unit * 2,
      paddingBottom: spacing.unit * 2
    }
  })

interface IBaseColProps extends WithStyles<typeof colStyles> {
  children: React.ReactNode
  className?: string
}

const BaseCol = withStyles(colStyles)(
  ({ children, className, classes }: IBaseColProps) => {
    return (
      <TableCell className={classNames(className, classes.col)}>
        {children}
      </TableCell>
    )
  }
)

interface IColProps {
  children: React.ReactNode
  className?: string
}

const Col = ({ children, className }: IColProps) => (
  <BaseCol className={className}>
    <Typography variant="body1">{children}</Typography>
  </BaseCol>
)

const KeyCol = ({ children, className }: IColProps) => (
  <BaseCol className={className}>
    <Typography variant="body1" color="textPrimary">
      {children}
    </Typography>
  </BaseCol>
)

const styles = () =>
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
          <KeyCol>Job ID</KeyCol>
          <Col>{jobRun.jobId}</Col>
        </TableRow>
        <TableRow>
          <KeyCol>Initiator</KeyCol>
          <Col>{jobRun.type}</Col>
        </TableRow>
        <TableRow>
          <KeyCol>Requester</KeyCol>
          <Col>{jobRun.requester}</Col>
        </TableRow>
        <TableRow>
          <KeyCol>Request ID</KeyCol>
          <Col>{jobRun.requestId}</Col>
        </TableRow>
        <TableRow>
          <KeyCol>Completed At</KeyCol>
          <Col>{jobRun.completedAt}</Col>
        </TableRow>
        {jobRun.error && (
          <TableRow>
            <KeyCol>Error</KeyCol>
            <Col>{jobRun.error}</Col>
          </TableRow>
        )}
        <TableRow>
          <KeyCol className={classes.bottomCol}>Tasks</KeyCol>
          <BaseCol className={classes.bottomCol}>
            <TaskRuns taskRuns={jobRun.taskRuns} />
          </BaseCol>
        </TableRow>
      </TableBody>
    </Table>
  )
}

export default withStyles(styles)(Details)
