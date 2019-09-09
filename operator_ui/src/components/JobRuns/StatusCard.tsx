import { PaddedCard } from '@chainlink/styleguide'
import { Grid } from '@material-ui/core'
import { createStyles, withStyles, WithStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { titleCase } from 'change-case'
import classNames from 'classnames'
import { IJobRun } from 'operator_ui'
import React from 'react'
import { useEffect, useHooks, useState } from 'use-react-hooks'
import ElapsedTime from '../ElapsedTime'
import StatusIcon from '../JobRuns/StatusIcon'

const styles = (theme: any) =>
  createStyles({
    completed: {
      backgroundColor: theme.palette.success.light
    },
    errored: {
      backgroundColor: theme.palette.error.light
    },
    pending: {
      backgroundColor: theme.palette.warning.light
    },
    statusCard: {
      '&:last-child': {
        paddingBottom: theme.spacing.unit * 2
      }
    },
    head: {
      display: 'flex',
      alignItems: 'center'
    },
    statusIcon: {
      display: 'inline-block'
    },
    statusText: {
      display: 'inline-block',
      textTransform: 'capitalize',
      color: theme.palette.secondary.main
    },
    statusRoot: {
      paddingLeft: theme.spacing.unit * 2
    },
    elapsedText: {
      color: theme.typography.display1.color
    },
    earnedLink: {
      color: theme.palette.success.main
    }
  })

interface IProps extends WithStyles<typeof styles> {
  title: string
  children?: React.ReactNode
  jobRun?: IJobRun
}

const selectLink = (inWei: number) => inWei / 1e18
const EarnedLink = ({
  classes,
  jobRun
}: {
  jobRun?: IJobRun
  classes: WithStyles<typeof styles>['classes']
}) => {
  const linkEarned = jobRun && jobRun.payment
  return (
    <Typography className={classes.earnedLink} variant="h6">
      +{linkEarned ? selectLink(linkEarned) : 0} Link
    </Typography>
  )
}

const StatusCard = useHooks(({ title, classes, children, jobRun }: IProps) => {
  const statusClass = classes[title as keyof typeof classes] || classes.pending
  const { status, createdAt, finishedAt } = jobRun || {
    status: '',
    createdAt: '',
    finishedAt: ''
  }
  const stillPending = status !== 'completed' && status !== 'errored'
  const [liveTime, setLiveTime] = useState(Date.now())
  useEffect(() => {
    if (stillPending) setInterval(() => setLiveTime(Date.now()), 1000)
  }, [])
  const endDate = stillPending ? liveTime.toString() : finishedAt

  return (
    <PaddedCard className={classNames(classes.statusCard, statusClass)}>
      <div className={classes.head}>
        <StatusIcon width={80}>{title}</StatusIcon>
        <Grid container alignItems="center" className={classes.statusRoot}>
          <Grid item xs={9}>
            <Typography className={classes.statusText} variant="h5">
              {titleCase(title)}
            </Typography>
            {
              <ElapsedTime
                start={createdAt}
                end={endDate}
                className={classes.elapsedText}
              />
            }
          </Grid>
          <Grid item xs={3}>
            {title === 'completed' && (
              <EarnedLink classes={classes} jobRun={jobRun} />
            )}
          </Grid>
        </Grid>
      </div>
      {children}
    </PaddedCard>
  )
})

export default withStyles(styles)(StatusCard)
