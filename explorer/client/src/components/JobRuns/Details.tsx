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

const itemContentStyles = ({ spacing, breakpoints }: Theme) =>
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
}

const Key = withStyles(itemContentStyles)(
  ({ children, classes }: IItemProps) => (
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
  ({ children, classes }: IItemProps) => (
    <BaseItem sm={12} md={8}>
      <Typography
        variant="body1"
        className={classNames(classes.value, classes.text)}>
        {children}
      </Typography>
    </BaseItem>
  )
)

const rowStyles = ({ palette }: Theme) =>
  createStyles({
    row: {
      borderBottom: 'solid 1px',
      borderBottomColor: palette.divider,
      display: 'block',
      width: '100%'
    }
  })

interface IRowProps extends WithStyles<typeof rowStyles> {
  children: React.ReactNode
  className?: string
}

const Row = withStyles(rowStyles)(
  ({ children, classes, className }: IRowProps) => (
    <div className={classNames(classes.row, className)}>
      <Grid container spacing={0}>
        {children}
      </Grid>
    </div>
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
  etherscanHost: string
}

const buildSearchQuery = (id: string) => `/job-runs?search=${id}`

const Details = ({ classes, jobRun, etherscanHost }: IProps) => {
  const nodeHasUrl = jobRun.chainlinkNode.url !== ''
  const nodeName = jobRun.chainlinkNode.name
  return (
    <div>
      <Row>
        <Key>Job ID</Key>
        <Value>
          <a href={buildSearchQuery(jobRun.jobId)}>{jobRun.jobId}</a>
        </Value>
      </Row>

      <Row>
        <Key>Node</Key>
        <Value>
          {nodeHasUrl ? (
            <a
              href={jobRun.chainlinkNode.url}
              target="_blank"
              rel="noopener noreferrer">
              {nodeName}
            </a>
          ) : (
            nodeName
          )}
        </Value>
      </Row>

      <Row>
        <Key>Initiator</Key>
        <Value>{jobRun.type}</Value>
      </Row>

      <Row>
        <Key>Requester</Key>
        <Value>
          <a href={buildSearchQuery(jobRun.requester)}>{jobRun.requester}</a>
        </Value>
      </Row>

      <Row>
        <Key>Request ID</Key>
        <Value>
          <a href={buildSearchQuery(jobRun.requestId)}>{jobRun.requestId}</a>
        </Value>
      </Row>

      <Row>
        <Key>Request Transaction Hash</Key>
        <Value>{jobRun.txHash}</Value>
      </Row>

      <Row>
        <Key>Finished At</Key>
        <Value>{jobRun.finishedAt && moment(jobRun.finishedAt).format()}</Value>
      </Row>

      {jobRun.error && (
        <Row>
          <Key>Error</Key>
          <Value>{jobRun.error}</Value>
        </Row>
      )}

      <Row className={classes.bottomRow}>
        <Key>Tasks</Key>
        <BaseItem sm={12} md={8} className={classes.task}>
          <TaskRuns taskRuns={jobRun.taskRuns} etherscanHost={etherscanHost} />
        </BaseItem>
      </Row>
    </div>
  )
}

export default withStyles(styles)(Details)
