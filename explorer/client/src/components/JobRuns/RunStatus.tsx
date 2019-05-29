import React from 'react'
import classNames from 'classnames'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import Card from '@material-ui/core/Card'
import Typography from '@material-ui/core/Typography'
import yellow from '@material-ui/core/colors/yellow'
import StatusIcon from '../Icons/Status'
import green from '../../colors/green'
import status from '../../utils/status'

const styles = ({ palette, spacing }: Theme) =>
  createStyles({
    completed: {
      backgroundColor: green['50']
    },
    errored: {
      backgroundColor: palette.error.light
    },
    pending: {
      backgroundColor: yellow['50']
    },
    statusCard: {
      display: 'flex',
      alignItems: 'center',
      '&:last-child': {
        padding: spacing.unit * 2
      }
    },
    statusText: {
      marginLeft: spacing.unit * 2
    }
  })

interface IProps extends WithStyles<typeof styles> {
  jobRun: IJobRun
}

const StatusCard = ({ classes, jobRun }: IProps) => {
  const [text, unfulfilled] = status(jobRun)
  const key = statusKey(jobRun, unfulfilled)
  const statusClass = classes[key] || classes.pending

  return (
    <Card className={classNames(classes.statusCard, statusClass)}>
      <StatusIcon width={80}>{jobRun.status}</StatusIcon>
      <Typography
        className={classes.statusText}
        variant="h5"
        color="textPrimary">
        {text} {unfulfilled && '(Receiver Reverted)'}
      </Typography>
    </Card>
  )
}

type Status = 'completed' | 'errored' | 'pending'

const statusKey = (jobRun: IJobRun, unfulfilled: boolean): Status => {
  return (unfulfilled ? 'errored' : jobRun.status) as Status
}

export default withStyles(styles)(StatusCard)
