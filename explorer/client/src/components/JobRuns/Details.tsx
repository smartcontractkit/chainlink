import React from 'react'
import classNames from 'classnames'
import moment from 'moment'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import Grid, { GridSize } from '@material-ui/core/Grid'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableRow from '@material-ui/core/TableRow'
import TableCell from '@material-ui/core/TableCell'
import Typography from '@material-ui/core/Typography'
import TaskRuns from './TaskRuns'

interface IBaseItemProps {
  children: React.ReactNode
  className?: string
  sm: GridSize
  md: GridSize
}

const BaseItem = ({ children, className, sm, md }: IBaseItemProps) => {
  return (
    <Grid item xs={sm} sm={sm} md={md} className={className}>
      {children}
    </Grid>
  )
}

const itemContentStyles = ({ spacing, breakpoints, palette }: Theme) =>
  createStyles({
    text: {
      paddingLeft: spacing.unit * 2,
      paddingRight: spacing.unit * 2,
      paddingBottom: spacing.unit
    },
    key: {
      paddingTop: spacing.unit * 2
    },
    value: {
      paddingTop: 0,
      [breakpoints.up('md')]: {
        paddingTop: spacing.unit * 2,
        paddingBottom: spacing.unit
      }
    }
  })

interface IItemProps extends WithStyles<typeof itemContentStyles> {
  children: React.ReactNode
  className?: string
}

const Key = withStyles(itemContentStyles)(
  ({ children, className, classes }: IItemProps) => (
    <BaseItem sm={12} md={4}>
      <Typography
        variant="body1"
        color="textPrimary"
        className={classNames(classes.key, classes.text)}>
        {children}
      </Typography>
    </BaseItem>
  )
)

const Value = withStyles(itemContentStyles)(
  ({ children, className, classes }: IItemProps) => (
    <BaseItem sm={12} md={8}>
      <Typography
        variant="body1"
        className={classNames(classes.value, classes.text)}>
        {children}
      </Typography>
    </BaseItem>
  )
)

const styles = ({ spacing, palette }: Theme) =>
  createStyles({
    row: {
      borderBottom: 'solid 1px',
      borderBottomColor: palette.divider,
      display: 'block',
      width: '100%'
    },
    bottomRow: {
      borderBottom: 'none'
    },
    task: {
      paddingLeft: spacing.unit * 2,
      paddingRight: spacing.unit * 2,
      paddingTop: spacing.unit
    }
  })

interface IProps extends WithStyles<typeof styles> {
  jobRun: IJobRun
}

const Details = ({ classes, jobRun }: IProps) => {
  return (
    <div>
      <div className={classes.row}>
        <Grid container spacing={0}>
          <Key>Job ID</Key>
          <Value>{jobRun.jobId}</Value>
        </Grid>
      </div>

      <div className={classes.row}>
        <Grid container spacing={0}>
          <Key>Node</Key>
          <Value>{jobRun.chainlinkNode.name}</Value>
        </Grid>
      </div>

      <div className={classes.row}>
        <Grid container spacing={0}>
          <Key>Initiator</Key>
          <Value>{jobRun.type}</Value>
        </Grid>
      </div>

      <div className={classes.row}>
        <Grid container spacing={0}>
          <Key>Requester</Key>
          <Value>{jobRun.requester}</Value>
        </Grid>
      </div>

      <div className={classes.row}>
        <Grid container spacing={0}>
          <Key>Request ID</Key>
          <Value>{jobRun.requestId}</Value>
        </Grid>
      </div>

      <div className={classes.row}>
        <Grid container spacing={0}>
          <Key>Request Transaction Hash</Key>
          <Value>{jobRun.txHash}</Value>
        </Grid>
      </div>

      <div className={classes.row}>
        <Grid container spacing={0}>
          <Key>Finished At</Key>
          <Value>
            {jobRun.finishedAt && moment(jobRun.finishedAt).format()}
          </Value>
        </Grid>
      </div>

      {jobRun.error && (
        <Grid container spacing={0}>
          <div className={classes.row}>
            <Key>Error</Key>
            <Value>{jobRun.error}</Value>
          </div>
        </Grid>
      )}

      <div className={classNames(classes.row, classes.bottomRow)}>
        <Grid container spacing={0}>
          <Key>Tasks</Key>
          <BaseItem sm={12} md={8} className={classes.task}>
            <TaskRuns taskRuns={jobRun.taskRuns} />
          </BaseItem>
        </Grid>
      </div>
    </div>
  )
}

export default withStyles(styles)(Details)
