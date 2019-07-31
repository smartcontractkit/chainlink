import React from 'react'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import classNames from 'classnames'
import PaddedCard from '@chainlink/styleguide/src/components/PaddedCard'
import { titleCase } from 'change-case'
import StatusIcon from '../JobRuns/StatusIcon'
import { Grid } from '@material-ui/core'
import { IJobRun } from '../../../@types/operator_ui'
import ElapsedTime from '../ElapsedTime'

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
  const linkEarned = jobRun && jobRun.result && jobRun.result.amount
  return (
    <Typography className={classes.earnedLink} variant="h6">
      +{linkEarned ? selectLink(linkEarned) : 0} Link
    </Typography>
  )
}

const StatusCard = ({ title, classes, children, jobRun }: IProps) => {
  const statusClass = classes[title as keyof typeof classes] || classes.pending
  const { status, createdAt, finishedAt } = jobRun || {
    status: '',
    createdAt: '',
    finishedAt: ''
  }
  return (
    <PaddedCard className={classNames(classes.statusCard, statusClass)}>
      <div className={classes.head}>
        <StatusIcon width={80}>{title}</StatusIcon>
        <Grid container alignItems="center" className={classes.statusRoot}>
          <Grid item xs={9}>
            <Typography className={classes.statusText} variant="h5">
              {titleCase(title)}
            </Typography>
            {status !== 'pending' && (
              <ElapsedTime
                start={createdAt}
                end={finishedAt}
                className={classes.elapsedText}
              />
            )}
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
}

export default withStyles(styles)(StatusCard)
